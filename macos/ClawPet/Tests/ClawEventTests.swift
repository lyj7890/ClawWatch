import XCTest
@testable import ClawPet

final class ClawEventTests: XCTestCase {
    func testThinkingIsSummarizedWithoutRawReasoning() throws {
        let json = """
        {
          "type": "message",
          "timestamp": "2026-06-09T00:00:00Z",
          "message": {
            "role": "assistant",
            "content": [{"type": "thinking", "thinking": "private chain of thought"}]
          }
        }
        """

        let event = try JSONDecoder().decode(ClawEvent.self, from: Data(json.utf8))

        XCTAssertEqual(event.activity?.mood, .thinking)
        XCTAssertEqual(event.activity?.detail, "分析上下文并规划行动")
        XCTAssertFalse(event.activity?.detail.contains("private") == true)
    }

    func testToolCallShowsToolName() throws {
        let json = """
        {
          "type": "message",
          "message": {
            "role": "assistant",
            "content": [{"type": "toolCall", "name": "exec_command"}]
          }
        }
        """

        let event = try JSONDecoder().decode(ClawEvent.self, from: Data(json.utf8))

        XCTAssertEqual(event.activity?.mood, .working)
        XCTAssertEqual(event.activity?.detail, "exec_command")
    }

    func testTrajectoryPromptShowsThinkingState() throws {
        let json = """
        {
          "traceSchema": "openclaw-trajectory",
          "type": "prompt.submitted",
          "ts": "2026-06-09T00:00:00Z"
        }
        """

        let event = try JSONDecoder().decode(ClawEvent.self, from: Data(json.utf8))

        XCTAssertEqual(event.activity?.mood, .thinking)
        XCTAssertEqual(event.activity?.title, "模型正在思考")
    }

    func testFinalReplyIsNotImmediatelyOverwrittenByTerminalTraceEvents() {
        let now = Date()
        let protectedUntil = now.addingTimeInterval(ActivityDisplayPolicy.finalReplyProtectionDuration)
        let completed = PetActivity(
            mood: .responding,
            title: "模型完成响应",
            detail: "正在整理输出结果",
            timestamp: now
        )
        let ended = PetActivity(
            mood: .idle,
            title: "任务完成",
            detail: "等待下一项任务",
            timestamp: now
        )

        XCTAssertFalse(ActivityDisplayPolicy.shouldPresent(completed, protectedUntil: protectedUntil, now: now))
        XCTAssertFalse(ActivityDisplayPolicy.shouldPresent(ended, protectedUntil: protectedUntil, now: now))
    }

    func testNewUserMessageCanReplaceProtectedReply() {
        let now = Date()
        let nextMessage = PetActivity(
            mood: .listening,
            title: "收到新任务",
            detail: "new message",
            timestamp: now
        )

        XCTAssertTrue(
            ActivityDisplayPolicy.shouldPresent(
                nextMessage,
                protectedUntil: now.addingTimeInterval(18),
                now: now
            )
        )
    }

    func testUserMessageRemovesChannelMetadataPrefix() throws {
        let json = """
        {
          "type": "message",
          "message": {
            "role": "user",
            "content": "[Octo user:123 +2m Tue 2026-06-09 21:24 GMT+8] 你好，帮我检查服务"
          }
        }
        """

        let event = try JSONDecoder().decode(ClawEvent.self, from: Data(json.utf8))

        XCTAssertEqual(event.activity?.title, "用户消息")
        XCTAssertEqual(event.activity?.detail, "你好，帮我检查服务")
    }

    func testFinalReplyIsCompleteAndDoesNotShowProcessingSpinner() throws {
        let json = """
        {
          "type": "message",
          "message": {
            "role": "assistant",
            "content": [{"type": "text", "text": "处理完成"}]
          }
        }
        """

        let activity = try JSONDecoder().decode(ClawEvent.self, from: Data(json.utf8)).activity

        XCTAssertTrue(activity?.isFinalReply == true)
        XCTAssertFalse(activity?.isProcessing == true)
    }

    func testThinkingAndToolCallsShowProcessingSpinner() throws {
        let thinking = PetActivity(mood: .thinking, title: "模型正在思考", detail: "处理中", timestamp: Date())
        let tool = PetActivity(mood: .working, title: "正在使用工具", detail: "exec", timestamp: Date())

        XCTAssertTrue(thinking.isProcessing)
        XCTAssertTrue(tool.isProcessing)
    }
}
