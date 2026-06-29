import Foundation
import XCTest
@testable import ClawPet

final class ClawWatchClientTests: XCTestCase {
    func testDirectLocalWatcherReceivesChineseFinalReply() async throws {
        let root = FileManager.default.temporaryDirectory
            .appending(path: "clawpet-direct-\(UUID().uuidString)")
        let sessions = root
            .appending(path: ".openclaw")
            .appending(path: "agents")
            .appending(path: "main")
            .appending(path: "sessions")
        let session = sessions.appending(path: "test.jsonl")
        try FileManager.default.createDirectory(at: sessions, withIntermediateDirectories: true)
        try Data(#"{"type":"message","message":{"role":"user","content":"中文历史"}}"#.utf8)
            .write(to: session)
        defer { try? FileManager.default.removeItem(at: root) }

        let client = ClawWatchClient(homeDirectory: root)
        let initial = try await client.initialActivities()
        let updates = await client.updates(watching: [initial.path])
        let waiter = Task {
            for try await activity in updates where activity.isFinalReply {
                return activity
            }
            throw CocoaError(.fileReadUnknown)
        }

        try await Task.sleep(for: .milliseconds(500))
        let handle = try FileHandle(forWritingTo: session)
        try handle.seekToEnd()
        let line = "\n" + #"{"type":"message","message":{"role":"assistant","content":[{"type":"text","text":"最终回复：中文内容已显示"}]}}"# + "\n"
        try handle.write(contentsOf: Data(line.utf8))
        try handle.close()

        let activity = try await waiter.value
        XCTAssertEqual(activity.detail, "最终回复：中文内容已显示")
        await client.disconnect()
    }
}
