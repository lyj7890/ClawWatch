import Foundation
import SwiftUI

enum PetMood: Equatable {
    case offline
    case idle
    case listening
    case thinking
    case working
    case responding
    case error

    var label: String {
        switch self {
        case .offline: "等待 ClawWatch"
        case .idle: "待命中"
        case .listening: "收到新任务"
        case .thinking: "思考中"
        case .working: "执行工具"
        case .responding: "回复中"
        case .error: "遇到问题"
        }
    }

    var symbol: String {
        switch self {
        case .offline: "wifi.slash"
        case .idle: "moon.stars.fill"
        case .listening: "ear.fill"
        case .thinking: "ellipsis.bubble.fill"
        case .working: "hammer.fill"
        case .responding: "text.bubble.fill"
        case .error: "exclamationmark.triangle.fill"
        }
    }

    var tint: Color {
        switch self {
        case .offline: .secondary
        case .idle: .indigo
        case .listening: .cyan
        case .thinking: .blue
        case .working: .orange
        case .responding: .green
        case .error: .red
        }
    }

    var displayPriority: Int {
        switch self {
        case .error: 100
        case .responding: 80
        case .listening: 70
        case .working: 60
        case .thinking: 50
        case .offline: 40
        case .idle: 10
        }
    }
}

struct PetActivity: Identifiable, Equatable {
    let id = UUID()
    let mood: PetMood
    let title: String
    let detail: String
    let timestamp: Date

    var isFinalReply: Bool {
        mood == .responding && title == "正在回复"
    }

    var isUserMessage: Bool {
        mood == .listening && title == "用户消息"
    }

    var isProcessing: Bool {
        mood == .thinking || mood == .working
    }

    var showsUserContext: Bool {
        isProcessing || isFinalReply
    }
}

struct ActivityDisplayPolicy {
    static let finalReplyProtectionDuration: TimeInterval = 18

    static func shouldPresent(
        _ next: PetActivity,
        protectedUntil: Date?,
        now: Date = Date()
    ) -> Bool {
        guard let protectedUntil, protectedUntil > now else { return true }
        if next.isFinalReply || next.mood == .error || next.mood == .listening {
            return true
        }
        return false
    }
}
