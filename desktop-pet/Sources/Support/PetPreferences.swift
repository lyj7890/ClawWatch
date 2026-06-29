import Foundation

enum PetPreferences {
    static let minimumScale = 0.6
    static let maximumScale = 1.6
    static let scaleStep = 0.1
    static let defaultScale = 1.0

    static func clampedScale(_ value: Double) -> Double {
        min(max(value, minimumScale), maximumScale)
    }
}

enum PetTheme: String, CaseIterable, Identifiable {
    case classic
    case cat
    case dog

    var id: String { rawValue }

    var name: String {
        switch self {
        case .classic: "经典龙虾"
        case .cat: "橘猫"
        case .dog: "柴犬"
        }
    }
}
