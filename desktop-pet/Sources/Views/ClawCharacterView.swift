import SwiftUI

struct ClawCharacterView: View {
    let mood: PetMood
    let theme: PetTheme
    @State private var blinking = false
    @State private var blinkTask: Task<Void, Never>?

    var body: some View {
        TimelineView(.animation(minimumInterval: 1.0 / 30.0)) { context in
            let pose = CharacterPose(mood: mood, time: context.date.timeIntervalSinceReferenceDate)

            ZStack {
                Ellipse()
                    .fill(.black.opacity(mood == .error ? 0.11 : 0.16))
                    .frame(width: 150, height: 24)
                    .offset(y: 88)
                    .scaleEffect(pose.shadowScale)

                stateEffects(pose, palette: palette)

                VStack(spacing: -8) {
                    antennae(pose, palette: palette)
                    face(pose, palette: palette)
                    claws(pose, palette: palette)
                }
                .scaleEffect(pose.bodyScale)
                .offset(x: pose.bodyX, y: pose.bodyY)
                .rotationEffect(.degrees(pose.bodyRotation))
            }
        }
        .onAppear {
            startNaturalBlinking()
        }
        .onDisappear {
            blinkTask?.cancel()
        }
    }

    private var palette: CharacterPalette {
        CharacterPalette(theme: theme)
    }

    @ViewBuilder
    private func stateEffects(_ pose: CharacterPose, palette: CharacterPalette) -> some View {
        switch mood {
        case .thinking:
            HStack(spacing: 5) {
                ForEach(0..<3) { index in
                    Circle()
                        .fill(palette.effectColor.opacity(0.85))
                        .frame(width: CGFloat(7 + index * 2), height: CGFloat(7 + index * 2))
                        .offset(y: pose.effectPulse * CGFloat(index.isMultiple(of: 2) ? -5 : 5))
                }
            }
            .offset(x: 78, y: -78)

        case .working:
            Image(systemName: "sparkles")
                .font(.system(size: 34, weight: .bold))
                .foregroundStyle(.orange)
                .scaleEffect(0.85 + abs(pose.effectPulse) * 0.35)
                .rotationEffect(.degrees(pose.effectRotation))
                .offset(x: pose.effectPulse > 0 ? 95 : -95, y: 14)

        case .responding:
            ZStack {
                ForEach(0..<6) { index in
                    Capsule()
                        .fill(.green.opacity(0.75))
                        .frame(width: 5, height: 20)
                        .offset(y: -88)
                        .rotationEffect(.degrees(Double(index) * 60 + pose.effectRotation))
                }
            }

        case .error:
            Image(systemName: "exclamationmark")
                .font(.system(size: 40, weight: .black))
                .foregroundStyle(.red)
                .offset(x: 82, y: -78 + pose.effectPulse * 3)

        case .listening:
            Image(systemName: "ear.fill")
                .font(.system(size: 28, weight: .bold))
                .foregroundStyle(.cyan)
                .scaleEffect(1 + abs(pose.effectPulse) * 0.18)
                .offset(x: 88, y: -38)

        default:
            EmptyView()
        }
    }

    private func antennae(_ pose: CharacterPose, palette: CharacterPalette) -> some View {
        HStack(spacing: theme == .classic ? 52 : (theme == .cat ? 88 : 96)) {
            antenna(palette)
                .rotationEffect(.degrees(-28 + pose.antennaRotation), anchor: .bottom)
            antenna(palette)
                .scaleEffect(x: -1, y: 1)
                .rotationEffect(.degrees(28 - pose.antennaRotation), anchor: .bottom)
        }
    }

    @ViewBuilder
    private func antenna(_ palette: CharacterPalette) -> some View {
        switch theme {
        case .classic:
            Capsule()
                .fill(palette.primary.opacity(0.9))
                .frame(width: 7, height: 52)
        case .cat:
            ZStack {
                Triangle()
                    .fill(palette.primary)
                    .frame(width: 54, height: 62)
                Triangle()
                    .fill(Color.pink.opacity(0.55))
                    .frame(width: 27, height: 34)
                    .offset(y: 9)
            }
        case .dog:
            ZStack {
                Capsule()
                    .fill(palette.primary)
                    .frame(width: 48, height: 72)
                Capsule()
                    .fill(palette.highlight.opacity(0.65))
                    .frame(width: 25, height: 48)
                    .offset(y: 8)
            }
        }
    }

    private func face(_ pose: CharacterPose, palette: CharacterPalette) -> some View {
        ZStack {
            faceShell(palette)

            HStack(spacing: theme == .dog ? 46 : 38) {
                eye(pupilOffset: pose.pupilOffset, palette: palette)
                eye(pupilOffset: pose.pupilOffset, palette: palette)
            }
            .offset(y: -16)

            facialDetails(palette)
        }
    }

    @ViewBuilder
    private func facialDetails(_ palette: CharacterPalette) -> some View {
        switch theme {
        case .classic:
            MouthShape(mood: mood)
                .stroke(.white.opacity(0.92), style: StrokeStyle(lineWidth: 4.5, lineCap: .round))
                .frame(width: 50, height: 25)
                .offset(y: 23)
        case .cat:
            ZStack {
                Triangle()
                    .fill(Color.pink.opacity(0.9))
                    .frame(width: 18, height: 13)
                    .rotationEffect(.degrees(180))
                    .offset(y: 14)
                MouthShape(mood: mood)
                    .stroke(palette.eyeColor.opacity(0.8), style: StrokeStyle(lineWidth: 3, lineCap: .round))
                    .frame(width: 38, height: 18)
                    .offset(y: 28)
                ForEach([-1.0, 1.0], id: \.self) { side in
                    VStack(spacing: 7) {
                        Capsule().frame(width: 42, height: 2).rotationEffect(.degrees(side * 8))
                        Capsule().frame(width: 42, height: 2).rotationEffect(.degrees(side * -8))
                    }
                    .foregroundStyle(palette.eyeColor.opacity(0.55))
                    .offset(x: side * 58, y: 25)
                }
            }
        case .dog:
            ZStack {
                Ellipse()
                    .fill(palette.muzzleColor)
                    .frame(width: 78, height: 52)
                    .offset(y: 25)
                RoundedRectangle(cornerRadius: 8)
                    .fill(palette.eyeColor)
                    .frame(width: 25, height: 16)
                    .offset(y: 12)
                MouthShape(mood: mood)
                    .stroke(palette.eyeColor.opacity(0.85), style: StrokeStyle(lineWidth: 3, lineCap: .round))
                    .frame(width: 42, height: 19)
                    .offset(y: 35)
                Circle().fill(palette.highlight).frame(width: 25, height: 25).offset(x: -62, y: 22)
                Circle().fill(palette.highlight).frame(width: 25, height: 25).offset(x: 62, y: 22)
            }
        }
    }

    @ViewBuilder
    private func faceShell(_ palette: CharacterPalette) -> some View {
        let fill = LinearGradient(
            colors: [palette.highlight, palette.primary],
            startPoint: .topLeading,
            endPoint: .bottomTrailing
        )

        switch theme {
        case .classic:
            RoundedRectangle(cornerRadius: 62, style: .continuous)
                .fill(fill)
                .frame(width: 166, height: 132)
        case .cat:
            ZStack {
                RoundedRectangle(cornerRadius: 66, style: .continuous)
                    .fill(fill)
                    .frame(width: 170, height: 138)
                VStack(spacing: 8) {
                    Capsule().fill(palette.eyeColor.opacity(0.22)).frame(width: 58, height: 7)
                    Capsule().fill(palette.eyeColor.opacity(0.18)).frame(width: 45, height: 6)
                }
                .offset(y: -49)
            }
        case .dog:
            ZStack {
                RoundedRectangle(cornerRadius: 58, style: .continuous)
                    .fill(fill)
                    .frame(width: 178, height: 142)
                Capsule()
                    .fill(palette.highlight.opacity(0.65))
                    .frame(width: 72, height: 58)
                    .offset(y: -39)
            }
        }
    }

    private func eye(pupilOffset: CGSize, palette: CharacterPalette) -> some View {
        ZStack {
            RoundedRectangle(cornerRadius: 16)
                .fill(.white.opacity(0.96))
                .frame(width: theme == .dog ? 28 : 32, height: blinking ? 5 : 29)

            RoundedRectangle(cornerRadius: 8)
                .fill(palette.eyeColor)
                .frame(width: theme == .cat ? 8 : 12, height: blinking ? 2 : 16)
                .offset(pupilOffset)

            Circle()
                .fill(.white.opacity(blinking ? 0 : 0.9))
                .frame(width: 4, height: 4)
                .offset(x: pupilOffset.width - 2.5, y: pupilOffset.height - 4)
        }
        .animation(.easeInOut(duration: 0.09), value: blinking)
    }

    private func claws(_ pose: CharacterPose, palette: CharacterPalette) -> some View {
        HStack(spacing: theme == .classic ? 94 : 76) {
            claw(palette).rotationEffect(.degrees(-18 + pose.leftClawRotation), anchor: .bottom)
            claw(palette)
                .scaleEffect(x: -1, y: 1)
                .rotationEffect(.degrees(18 + pose.rightClawRotation), anchor: .bottom)
        }
    }

    @ViewBuilder
    private func claw(_ palette: CharacterPalette) -> some View {
        switch theme {
        case .classic:
            ZStack {
                Circle().trim(from: 0.12, to: 0.88)
                    .stroke(palette.primary, style: StrokeStyle(lineWidth: 18, lineCap: .round))
                    .frame(width: 54, height: 54)
                Capsule().fill(palette.primary).frame(width: 17, height: 48).offset(y: 31)
            }
        case .cat:
            ZStack {
                Capsule().fill(palette.primary).frame(width: 40, height: 65)
                HStack(spacing: 4) {
                    ForEach(0..<3) { _ in
                        Capsule().fill(palette.highlight).frame(width: 7, height: 22)
                    }
                }
                .offset(y: -25)
            }
        case .dog:
            ZStack {
                Capsule().fill(palette.primary).frame(width: 46, height: 68)
                Ellipse().fill(palette.highlight).frame(width: 38, height: 30).offset(y: -25)
                HStack(spacing: 4) {
                    ForEach(0..<3) { _ in
                        Circle().fill(palette.eyeColor.opacity(0.55)).frame(width: 6, height: 6)
                    }
                }
                .offset(y: -28)
            }
        }
    }

    private func startNaturalBlinking() {
        blinkTask?.cancel()
        blinkTask = Task {
            while !Task.isCancelled {
                try? await Task.sleep(for: .seconds(Double.random(in: 4...7)))
                guard !Task.isCancelled else { return }
                blinking = true
                try? await Task.sleep(for: .milliseconds(130))
                blinking = false
            }
        }
    }
}

private struct CharacterPose {
    let bodyX: CGFloat
    let bodyY: CGFloat
    let bodyRotation: Double
    let bodyScale: CGFloat
    let shadowScale: CGFloat
    let antennaRotation: Double
    let leftClawRotation: Double
    let rightClawRotation: Double
    let pupilOffset: CGSize
    let effectPulse: CGFloat
    let effectRotation: Double

    init(mood: PetMood, time: TimeInterval) {
        let slow = sin(time * 2.2)
        let medium = sin(time * 5.0)
        let fast = sin(time * 10.0)
        let pulse = CGFloat(sin(time * 6.0))
        effectPulse = pulse
        effectRotation = time * 90

        switch mood {
        case .listening:
            bodyX = 0
            bodyY = -8 + CGFloat(slow * 4)
            bodyRotation = 0
            bodyScale = 1.08
            shadowScale = 0.84
            antennaRotation = 16
            leftClawRotation = 12
            rightClawRotation = -12
            pupilOffset = CGSize(width: 0, height: 4)
        case .thinking:
            bodyX = CGFloat(slow * 7)
            bodyY = CGFloat(-8 + abs(medium) * -3)
            bodyRotation = slow * 13
            bodyScale = 1.03
            shadowScale = 0.82 + CGFloat(abs(slow) * 0.08)
            antennaRotation = medium * 11
            leftClawRotation = slow * 10
            rightClawRotation = slow * 10
            pupilOffset = CGSize(width: slow * 7, height: -3)
        case .working:
            bodyX = CGFloat(fast * 6)
            bodyY = CGFloat(-8 + abs(fast) * -6)
            bodyRotation = fast * 5
            bodyScale = 1.04
            shadowScale = 0.76 + CGFloat(abs(fast) * 0.16)
            antennaRotation = fast * 9
            leftClawRotation = fast * 42
            rightClawRotation = -fast * 42
            pupilOffset = CGSize(width: fast * 4, height: 2)
        case .responding:
            bodyX = 0
            bodyY = CGFloat(-10 - abs(medium) * 13)
            bodyRotation = medium * 2
            bodyScale = 1.06 + CGFloat(abs(medium) * 0.06)
            shadowScale = 0.68 + CGFloat(abs(medium) * 0.18)
            antennaRotation = medium * 12
            leftClawRotation = 25 + medium * 12
            rightClawRotation = -25 - medium * 12
            pupilOffset = .zero
        case .error:
            bodyX = CGFloat(fast * 2)
            bodyY = 16
            bodyRotation = fast * 2.5
            bodyScale = 0.92
            shadowScale = 1.14
            antennaRotation = -22
            leftClawRotation = -22
            rightClawRotation = 22
            pupilOffset = CGSize(width: 0, height: 5)
        case .offline:
            bodyX = 0
            bodyY = 7
            bodyRotation = 0
            bodyScale = 0.98
            shadowScale = 1
            antennaRotation = -8
            leftClawRotation = -5
            rightClawRotation = 5
            pupilOffset = CGSize(width: 0, height: 2)
        case .idle:
            bodyX = 0
            bodyY = 0
            bodyRotation = 0
            bodyScale = 1
            shadowScale = 0.94
            antennaRotation = 0
            leftClawRotation = 0
            rightClawRotation = 0
            pupilOffset = .zero
        }
    }
}

private struct CharacterPalette {
    let primary: Color
    let highlight: Color
    let eyeColor: Color
    let effectColor: Color
    let muzzleColor: Color

    init(theme: PetTheme) {
        switch theme {
        case .classic:
            primary = Color(red: 0.75, green: 0.04, blue: 0.12)
            highlight = Color(red: 1, green: 0.32, blue: 0.25)
            eyeColor = Color(red: 0.18, green: 0.08, blue: 0.10)
            effectColor = .blue
            muzzleColor = .white.opacity(0.85)
        case .cat:
            primary = Color(red: 0.88, green: 0.42, blue: 0.08)
            highlight = Color(red: 1.0, green: 0.72, blue: 0.28)
            eyeColor = Color(red: 0.18, green: 0.26, blue: 0.08)
            effectColor = .mint
            muzzleColor = Color(red: 1.0, green: 0.82, blue: 0.55)
        case .dog:
            primary = Color(red: 0.62, green: 0.28, blue: 0.08)
            highlight = Color(red: 0.96, green: 0.70, blue: 0.38)
            eyeColor = Color(red: 0.18, green: 0.09, blue: 0.04)
            effectColor = .orange
            muzzleColor = Color(red: 1.0, green: 0.86, blue: 0.65)
        }
    }
}

private struct Triangle: Shape {
    func path(in rect: CGRect) -> Path {
        var path = Path()
        path.move(to: CGPoint(x: rect.midX, y: rect.minY))
        path.addLine(to: CGPoint(x: rect.maxX, y: rect.maxY))
        path.addLine(to: CGPoint(x: rect.minX, y: rect.maxY))
        path.closeSubpath()
        return path
    }
}

private struct MouthShape: Shape {
    let mood: PetMood

    func path(in rect: CGRect) -> Path {
        var path = Path()

        switch mood {
        case .error:
            path.move(to: CGPoint(x: rect.minX + 5, y: rect.maxY - 5))
            path.addQuadCurve(
                to: CGPoint(x: rect.maxX - 5, y: rect.maxY - 5),
                control: CGPoint(x: rect.midX, y: rect.minY + 2)
            )
        case .thinking, .working:
            path.move(to: CGPoint(x: rect.minX + 8, y: rect.midY))
            path.addQuadCurve(
                to: CGPoint(x: rect.maxX - 8, y: rect.midY),
                control: CGPoint(x: rect.midX, y: rect.midY + 4)
            )
        default:
            path.move(to: CGPoint(x: rect.minX + 5, y: rect.minY + 5))
            path.addQuadCurve(
                to: CGPoint(x: rect.maxX - 5, y: rect.minY + 5),
                control: CGPoint(x: rect.midX, y: rect.maxY - 2)
            )
        }

        return path
    }
}
