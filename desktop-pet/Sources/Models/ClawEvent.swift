import Foundation

struct ClawEvent: Decodable {
    struct Message: Decodable {
        struct Content: Decodable {
            let type: String
            let text: String?
            let thinking: String?
            let name: String?
            let isError: Bool?
        }

        let role: String?
        let content: [Content]
        let errorMessage: String?

        private enum CodingKeys: String, CodingKey {
            case role
            case content
            case errorMessage
        }

        init(from decoder: Decoder) throws {
            let container = try decoder.container(keyedBy: CodingKeys.self)
            role = try container.decodeIfPresent(String.self, forKey: .role)
            errorMessage = try container.decodeIfPresent(String.self, forKey: .errorMessage)
            if let items = try? container.decode([Content].self, forKey: .content) {
                content = items
            } else if let text = try? container.decode(String.self, forKey: .content) {
                content = [Content(type: "text", text: text, thinking: nil, name: nil, isError: nil)]
            } else {
                content = []
            }
        }
    }

    struct TraceData: Decodable {
        let status: String?
    }

    let type: String?
    let timestamp: String?
    let ts: String?
    let message: Message?
    let modelId: String?
    let thinkingLevel: String?
    let data: TraceData?

    var activity: PetActivity? {
        let date = (timestamp ?? ts).flatMap(Self.dateFormatter.date) ?? Date()

        if let error = message?.errorMessage {
            return PetActivity(mood: .error, title: "执行遇到问题", detail: error.clawPreview, timestamp: date)
        }

        switch type {
        case "session.started":
            return PetActivity(mood: .listening, title: "开始处理新任务", detail: "正在准备运行环境", timestamp: date)
        case "context.compiled":
            return PetActivity(mood: .thinking, title: "正在理解上下文", detail: "整理会话、记忆与可用工具", timestamp: date)
        case "prompt.submitted":
            return PetActivity(mood: .thinking, title: "模型正在思考", detail: "等待模型生成下一步行动", timestamp: date)
        case "model.completed":
            return PetActivity(mood: .responding, title: "模型完成响应", detail: "正在整理输出结果", timestamp: date)
        case "session.ended":
            let failed = data?.status != nil && data?.status != "success"
            return PetActivity(
                mood: failed ? .error : .idle,
                title: failed ? "任务执行失败" : "任务完成",
                detail: failed ? (data?.status ?? "未知错误") : "等待下一项任务",
                timestamp: date
            )
        default:
            break
        }

        if type == "model_change" {
            return PetActivity(mood: .idle, title: "切换模型", detail: modelId ?? "新的模型", timestamp: date)
        }

        if type == "thinking_level_change" {
            return PetActivity(mood: .thinking, title: "调整思考强度", detail: thinkingLevel ?? "正在准备", timestamp: date)
        }

        guard let message else { return nil }

        if message.role == "user" {
            return PetActivity(
                mood: .listening,
                title: "用户消息",
                detail: message.firstText?.userMessagePreview ?? "收到新任务",
                timestamp: date
            )
        }

        if let tool = message.content.first(where: { $0.type == "toolCall" }) {
            return PetActivity(
                mood: tool.isError == true ? .error : .working,
                title: "正在使用工具",
                detail: tool.name ?? "处理本地任务",
                timestamp: date
            )
        }

        if message.content.contains(where: { $0.type == "thinking" }) {
            return PetActivity(
                mood: .thinking,
                title: "正在思考下一步",
                detail: "分析上下文并规划行动",
                timestamp: date
            )
        }

        if message.role == "assistant", let text = message.firstText {
            return PetActivity(mood: .responding, title: "正在回复", detail: text.clawPreview, timestamp: date)
        }

        if message.role == "tool" || message.role == "toolResult" {
            return PetActivity(mood: .working, title: "工具执行完成", detail: "正在检查结果", timestamp: date)
        }

        return nil
    }

    private static let dateFormatter = ISO8601DateFormatter()
}

private extension ClawEvent.Message {
    var firstText: String? {
        content.first(where: { $0.type == "text" })?.text
    }
}

private extension String {
    var clawPreview: String {
        let compact = replacingOccurrences(of: #"\s+"#, with: " ", options: .regularExpression)
            .trimmingCharacters(in: .whitespacesAndNewlines)
        guard compact.count > 96 else { return compact }
        return String(compact.prefix(96)) + "…"
    }

    var userMessagePreview: String {
        let withoutChannelPrefix = replacingOccurrences(
            of: #"^\[[^\]]+\]\s*"#,
            with: "",
            options: .regularExpression
        )
        return withoutChannelPrefix.clawPreview
    }
}
