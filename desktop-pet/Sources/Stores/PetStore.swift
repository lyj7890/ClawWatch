import Combine
import Foundation

@MainActor
final class PetStore: ObservableObject {
    @Published private(set) var connected = false
    @Published private(set) var current = PetActivity(
        mood: .offline,
        title: "正在寻找 OpenClaw",
        detail: "请启动本地 ClawWatch 服务",
        timestamp: Date()
    )
    @Published private(set) var latestReply: PetActivity?
    @Published private(set) var latestUserMessage: PetActivity?
    @Published private(set) var recent: [PetActivity] = []
    @Published private(set) var eventRevision = 0
    @Published private(set) var availableAgents: [String] = ["main"]
    @Published private(set) var selectedAgent: String
    @Published private(set) var connectedAgent: String?

    private let client = ClawWatchClient()
    private var connectionTask: Task<Void, Never>?
    private var finalReplyProtectedUntil: Date?

    init() {
        selectedAgent = UserDefaults.standard.string(forKey: "selectedAgent") ?? "main"
    }

    func start() {
        guard connectionTask == nil else { return }
        refreshAgents()
        connectionTask = Task { await connectLoop() }
    }

    func stop() {
        connectionTask?.cancel()
        connectionTask = nil
        connected = false
        connectedAgent = nil
        Task { await client.disconnect() }
    }

    func reconnect() {
        stop()
        start()
    }

    func selectAgent(_ agent: String) {
        guard selectedAgent != agent else { return }
        selectedAgent = agent
        UserDefaults.standard.set(agent, forKey: "selectedAgent")
        latestReply = nil
        latestUserMessage = nil
        recent = []
        finalReplyProtectedUntil = nil
        current = PetActivity(
            mood: .offline,
            title: "正在切换 Agent",
            detail: "准备监听 \(agent) 的最新会话",
            timestamp: Date()
        )
        eventRevision += 1
        reconnect()
    }

    func refreshAgents() {
        Task {
            let agents = await client.availableAgents()
            availableAgents = agents
            if !agents.contains(selectedAgent), let fallback = agents.first {
                selectAgent(fallback)
            }
        }
    }

    private func connectLoop() async {
        while !Task.isCancelled {
            let agent = selectedAgent
            do {
                let initial = try await client.initialActivities(agent: agent)
                connected = true
                connectedAgent = agent
                recent = Array(initial.activities.suffix(5).reversed())
                latestReply = initial.activities.last(where: { $0.mood == .responding && $0.title == "正在回复" })
                latestUserMessage = initial.activities.last(where: \.isUserMessage)
                if let latest = recent.first {
                    current = latest
                }

                let sessionMonitor = Task {
                    while !Task.isCancelled {
                        try await Task.sleep(for: .seconds(1))
                        do {
                            let latestPath = try await client.latestSessionPath(agent: agent)
                            if latestPath != initial.path {
                                await client.disconnect()
                                return
                            }
                        } catch {
                            // A transient health-check failure should not stop
                            // the active WebSocket subscription.
                        }
                    }
                }

                let trajectoryPath = initial.path.replacingOccurrences(of: ".jsonl", with: ".trajectory.jsonl")
                for try await activity in await client.updates(watching: [initial.path, trajectoryPath]) {
                    receive(activity)
                }
                sessionMonitor.cancel()

                // A closed socket is often an intentional session switch.
                connected = true
                try? await Task.sleep(for: .milliseconds(150))
                continue
            } catch {
                connected = false
                connectedAgent = nil
                current = PetActivity(
                    mood: .offline,
                    title: "等待 \(agent) 的会话",
                    detail: "未找到可监听的 session，将自动重试",
                    timestamp: Date()
                )
            }

            try? await Task.sleep(for: .seconds(3))
        }
    }

    private func receive(_ activity: PetActivity) {
        guard ActivityDisplayPolicy.shouldPresent(activity, protectedUntil: finalReplyProtectedUntil) else {
            return
        }

        if activity.isFinalReply {
            latestReply = activity
            finalReplyProtectedUntil = Date().addingTimeInterval(ActivityDisplayPolicy.finalReplyProtectionDuration)
        }
        if activity.isUserMessage {
            latestUserMessage = activity
        }
        current = activity
        recent.insert(activity, at: 0)
        recent = Array(recent.prefix(5))
        eventRevision += 1
    }
}
