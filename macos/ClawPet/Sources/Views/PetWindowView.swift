import AppKit
import SwiftUI

struct PetWindowView: View {
    @ObservedObject var store: PetStore
    @State private var isHovering = false
    @State private var visibleActivity: PetActivity?
    @State private var hideTask: Task<Void, Never>?
    @State private var characterMood: PetMood = .idle
    @State private var settleTask: Task<Void, Never>?
    @AppStorage("petScale") private var petScale = PetPreferences.defaultScale
    @AppStorage("petTheme") private var petTheme = PetTheme.classic.rawValue

    var body: some View {
        ZStack(alignment: .topTrailing) {
            VStack(spacing: 8) {
                if let activity = visibleActivity {
                    eventBubble(activity)
                        .transition(.move(edge: .bottom).combined(with: .opacity))
                }

                ClawCharacterView(
                    mood: characterMood,
                    theme: PetTheme(rawValue: petTheme) ?? .classic
                )
                    .scaleEffect(petScale)
                    .frame(width: 230 * petScale, height: 210 * petScale)
            }

            hoverControls
                .opacity(isHovering ? 1 : 0)
        }
        .padding(8)
        .frame(width: max(300, 340 * petScale))
        .animation(.spring(duration: 0.3), value: visibleActivity)
        .animation(.spring(duration: 0.25), value: petScale)
        .onHover { hovering in
            isHovering = hovering
        }
        .onChange(of: store.eventRevision) {
            present(store.current)
            animateCharacter(for: store.current)
        }
        .onAppear {
            characterMood = initialCharacterMood(for: store.current)
        }
        .onDisappear {
            hideTask?.cancel()
            settleTask?.cancel()
        }
    }

    private var hoverControls: some View {
        Button {
            NSApp.keyWindow?.orderOut(nil)
        } label: {
            Image(systemName: "xmark.circle.fill")
                .foregroundStyle(.secondary)
        }
        .buttonStyle(.plain)
        .help("隐藏宠物，可从菜单栏重新显示")
        .font(.system(size: 14, weight: .semibold))
        .padding(8)
    }

    private func eventBubble(_ activity: PetActivity) -> some View {
        VStack(alignment: .leading, spacing: 7) {
            if activity.showsUserContext, let userMessage = store.latestUserMessage {
                HStack(alignment: .top, spacing: 7) {
                    Image(systemName: "person.fill")
                        .foregroundStyle(.cyan)
                    Text(userMessage.detail)
                        .font(.callout)
                        .lineLimit(3)
                }
                Divider()
            }

            HStack(alignment: .firstTextBaseline, spacing: 7) {
                if activity.isProcessing {
                    SpinningProgressView(tint: activity.mood.tint)
                } else if activity.isFinalReply {
                    Image(systemName: "checkmark.circle.fill")
                        .foregroundStyle(.green)
                } else {
                    Image(systemName: activity.mood.symbol)
                        .foregroundStyle(activity.mood.tint)
                }

                if activity.isProcessing {
                    Text(activity.title)
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(activity.mood.tint)
                }

                Text(activity.detail)
                    .font(.callout)
                    .foregroundStyle(activity.isFinalReply ? .primary : .secondary)
                    .textSelection(.enabled)
                    .lineLimit(activity.isFinalReply ? 8 : 4)
                    .frame(maxWidth: .infinity, alignment: .leading)
            }
        }
        .padding(12)
        .background(.regularMaterial, in: RoundedRectangle(cornerRadius: 16, style: .continuous))
        .overlay {
            RoundedRectangle(cornerRadius: 16, style: .continuous)
                .stroke(activity.mood.tint.opacity(0.25), lineWidth: 1)
        }
        .shadow(color: .black.opacity(0.16), radius: 14, y: 6)
    }

    private func present(_ activity: PetActivity) {
        hideTask?.cancel()
        visibleActivity = activity

        let seconds: Double
        switch activity.mood {
        case .thinking, .working:
            seconds = 12
        case .responding:
            seconds = 18
        case .listening:
            seconds = 7
        case .error:
            seconds = 20
        default:
            seconds = 4
        }

        hideTask = Task {
            try? await Task.sleep(for: .seconds(seconds))
            guard !Task.isCancelled else { return }
            visibleActivity = nil
        }
    }

    private func animateCharacter(for activity: PetActivity) {
        settleTask?.cancel()
        characterMood = activity.mood

        let settleAfter: Double?
        switch activity.mood {
        case .responding:
            settleAfter = 1.5
        case .listening:
            settleAfter = 1.0
        default:
            settleAfter = nil
        }

        guard let settleAfter else { return }
        settleTask = Task {
            try? await Task.sleep(for: .seconds(settleAfter))
            guard !Task.isCancelled else { return }
            characterMood = .idle
        }
    }

    private func initialCharacterMood(for activity: PetActivity) -> PetMood {
        switch activity.mood {
        case .responding, .listening:
            return .idle
        default:
            return activity.mood
        }
    }
}

private struct SpinningProgressView: View {
    let tint: Color
    @State private var rotation = 0.0

    var body: some View {
        Image(systemName: "arrow.trianglehead.2.clockwise.rotate.90")
            .font(.system(size: 13, weight: .semibold))
            .foregroundStyle(tint)
            .rotationEffect(.degrees(rotation))
            .onAppear {
                withAnimation(.linear(duration: 1).repeatForever(autoreverses: false)) {
                    rotation = 360
                }
            }
    }
}
