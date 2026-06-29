import AppKit
import SwiftUI

@main
struct ClawPetApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) private var appDelegate
    @StateObject private var store = PetStore()
    @AppStorage("petScale") private var petScale = PetPreferences.defaultScale
    @AppStorage("petTheme") private var petTheme = PetTheme.classic.rawValue

    var body: some Scene {
        WindowGroup("ClawPet", id: "pet") {
            PetWindowView(store: store)
                .background(WindowConfigurator())
                .task {
                    store.start()
                }
        }
        .windowStyle(.hiddenTitleBar)
        .windowResizability(.contentSize)
        .commands {
            CommandGroup(replacing: .newItem) {}
            CommandMenu("ClawPet") {
                Button("重新连接") {
                    store.reconnect()
                }
                .keyboardShortcut("r", modifiers: [.command])

                Divider()

                Button("放大宠物") {
                    petScale = PetPreferences.clampedScale(petScale + PetPreferences.scaleStep)
                }
                .keyboardShortcut("+", modifiers: [.command])

                Button("缩小宠物") {
                    petScale = PetPreferences.clampedScale(petScale - PetPreferences.scaleStep)
                }
                .keyboardShortcut("-", modifiers: [.command])

                Button("恢复默认大小") {
                    petScale = PetPreferences.defaultScale
                }
                .keyboardShortcut("0", modifiers: [.command])

                Divider()

                Menu("切换宠物") {
                    ForEach(PetTheme.allCases) { theme in
                        Button {
                            petTheme = theme.rawValue
                        } label: {
                            if petTheme == theme.rawValue {
                                Label(theme.name, systemImage: "checkmark")
                            } else {
                                Text(theme.name)
                            }
                        }
                    }
                }
            }
        }

        Settings {
            SettingsView(store: store)
        }

        MenuBarExtra("ClawPet", systemImage: "pawprint.fill") {
            MenuBarView(store: store)
        }
    }
}

final class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
    }

    func applicationShouldHandleReopen(_ sender: NSApplication, hasVisibleWindows flag: Bool) -> Bool {
        if !flag {
            Self.showPetWindow()
        }
        return true
    }

    static func showPetWindow() {
        for window in NSApp.windows where window.title == "ClawPet" {
            window.orderFrontRegardless()
        }
    }
}

private struct MenuBarView: View {
    @ObservedObject var store: PetStore
    @Environment(\.openSettings) private var openSettings

    var body: some View {
        Button("显示宠物") {
            AppDelegate.showPetWindow()
        }

        Button("重新连接") {
            store.reconnect()
        }

        Label(agentStatusText, systemImage: store.connectedAgent == nil ? "arrow.trianglehead.2.clockwise.rotate.90" : "checkmark.circle.fill")
            .foregroundStyle(store.connectedAgent == nil ? Color.secondary : Color.green)

        Menu("选择 Agent · \(store.selectedAgent)") {
            ForEach(store.availableAgents, id: \.self) { agent in
                Button {
                    store.selectAgent(agent)
                } label: {
                    if store.connectedAgent == agent {
                        Label(agent, systemImage: "checkmark.circle.fill")
                    } else if store.selectedAgent == agent {
                        Label(agent, systemImage: "arrow.trianglehead.2.clockwise.rotate.90")
                    } else {
                        Text(agent)
                    }
                }
            }
        }

        Divider()

        Button("设置") {
            openSettings()
        }

        Divider()

        Button("退出 ClawPet") {
            NSApp.terminate(nil)
        }
    }

    private var agentStatusText: String {
        if let connectedAgent = store.connectedAgent {
            return "当前连接：\(connectedAgent)"
        }
        return "正在连接：\(store.selectedAgent)"
    }
}
