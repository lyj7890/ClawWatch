import Foundation

actor ClawWatchClient {
    private let fileManager = FileManager.default
    private let sessionsRoot: URL
    private var pollingTask: Task<Void, Never>?

    init(homeDirectory: URL = FileManager.default.homeDirectoryForCurrentUser) {
        sessionsRoot = homeDirectory
            .appending(path: ".openclaw")
            .appending(path: "agents")
    }

    func availableAgents() -> [String] {
        guard let directories = try? fileManager.contentsOfDirectory(
            at: sessionsRoot,
            includingPropertiesForKeys: [.isDirectoryKey],
            options: [.skipsHiddenFiles]
        ) else {
            return ["main"]
        }

        let agents = directories.compactMap { directory -> String? in
            let sessions = directory.appending(path: "sessions")
            var isDirectory: ObjCBool = false
            guard fileManager.fileExists(atPath: sessions.path, isDirectory: &isDirectory),
                  isDirectory.boolValue else {
                return nil
            }
            return directory.lastPathComponent
        }
        return agents.isEmpty ? ["main"] : agents.sorted()
    }

    func initialActivities(agent: String = "main") async throws -> (path: String, activities: [PetActivity]) {
        let path = try await latestSessionPath(agent: agent)
        let events = try readEvents(at: path)
        return (path, events.compactMap(\.activity))
    }

    func latestSessionPath(agent: String = "main") async throws -> String {
        let directory = sessionsRoot
            .appending(path: agent)
            .appending(path: "sessions")
        let files = try fileManager.contentsOfDirectory(
            at: directory,
            includingPropertiesForKeys: [.contentModificationDateKey],
            options: [.skipsHiddenFiles]
        )
        let sessions = files.filter {
            $0.pathExtension == "jsonl"
                && !$0.lastPathComponent.contains(".trajectory")
                && !$0.lastPathComponent.contains(".checkpoint")
        }
        guard let latest = try sessions.max(by: {
            try modificationDate($0) < modificationDate($1)
        }) else {
            throw CocoaError(.fileNoSuchFile)
        }
        return latest.path
    }

    func updates(watching paths: [String]) -> AsyncThrowingStream<PetActivity, Error> {
        AsyncThrowingStream { continuation in
            pollingTask?.cancel()
            pollingTask = Task {
                var offsets = Dictionary(uniqueKeysWithValues: paths.map { ($0, fileSize($0)) })

                while !Task.isCancelled {
                    for path in paths {
                        do {
                            let newEvents = try readNewEvents(at: path, offset: &offsets[path, default: 0])
                            for event in newEvents {
                                if let activity = event.activity {
                                    continuation.yield(activity)
                                }
                            }
                        } catch {
                            // Files may not exist until the trajectory starts.
                        }
                    }
                    try? await Task.sleep(for: .milliseconds(300))
                }

                continuation.finish()
            }

            continuation.onTermination = { _ in
                Task { await self.disconnect() }
            }
        }
    }

    func disconnect() {
        pollingTask?.cancel()
        pollingTask = nil
    }

    private func readEvents(at path: String) throws -> [ClawEvent] {
        let data = try Data(contentsOf: URL(filePath: path))
        return decodeLines(data)
    }

    private func readNewEvents(at path: String, offset: inout UInt64) throws -> [ClawEvent] {
        let handle = try FileHandle(forReadingFrom: URL(filePath: path))
        defer { try? handle.close() }

        let size = try handle.seekToEnd()
        guard size > offset else {
            offset = size
            return []
        }

        try handle.seek(toOffset: offset)
        let data = try handle.readToEnd() ?? Data()
        offset = size
        return decodeLines(data)
    }

    private func decodeLines(_ data: Data) -> [ClawEvent] {
        String(decoding: data, as: UTF8.self)
            .split(separator: "\n")
            .compactMap { try? JSONDecoder().decode(ClawEvent.self, from: Data($0.utf8)) }
    }

    private func fileSize(_ path: String) -> UInt64 {
        let attributes = try? fileManager.attributesOfItem(atPath: path)
        return attributes?[.size] as? UInt64 ?? 0
    }

    private func modificationDate(_ url: URL) throws -> Date {
        try url.resourceValues(forKeys: [.contentModificationDateKey]).contentModificationDate ?? .distantPast
    }
}
