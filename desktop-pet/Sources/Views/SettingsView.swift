import SwiftUI

struct SettingsView: View {
    @ObservedObject var store: PetStore
    @AppStorage("petScale") private var petScale = PetPreferences.defaultScale
    @AppStorage("petTheme") private var petTheme = PetTheme.classic.rawValue

    var body: some View {
        Form {
            Section("外观") {
                Picker("宠物外观", selection: $petTheme) {
                    ForEach(PetTheme.allCases) { theme in
                        Text(theme.name).tag(theme.rawValue)
                    }
                }

                HStack {
                    Text("宠物大小")
                    Slider(
                        value: $petScale,
                        in: PetPreferences.minimumScale...PetPreferences.maximumScale,
                        step: PetPreferences.scaleStep
                    )
                    Text("\(Int(petScale * 100))%")
                        .monospacedDigit()
                        .frame(width: 44, alignment: .trailing)
                }

                Button("恢复默认大小") {
                    petScale = PetPreferences.defaultScale
                }
            }

            Section("连接") {
                LabeledContent("当前连接") {
                    Label(
                        store.connectedAgent ?? "正在连接 \(store.selectedAgent)",
                        systemImage: store.connectedAgent == nil
                            ? "arrow.trianglehead.2.clockwise.rotate.90"
                            : "checkmark.circle.fill"
                    )
                    .foregroundStyle(store.connectedAgent == nil ? Color.secondary : Color.green)
                }

                Picker("Agent", selection: agentSelection) {
                    ForEach(store.availableAgents, id: \.self) { agent in
                        Text(agent).tag(agent)
                    }
                }

                LabeledContent(
                    "数据来源",
                    value: "~/.openclaw/agents/\(store.selectedAgent)/sessions"
                )
                Button("刷新 Agent 列表") {
                    store.refreshAgents()
                }
                Text("ClawPet 直接监听本地 session 与 trajectory 文件，不依赖 ClawWatch 后端。")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }

            Section("隐私") {
                Text("宠物仅展示活动摘要、工具名称与回复预览，不展示原始隐藏推理。")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
        }
        .formStyle(.grouped)
        .frame(width: 480, height: 380)
        .onAppear {
            store.refreshAgents()
        }
    }

    private var agentSelection: Binding<String> {
        Binding(
            get: { store.selectedAgent },
            set: { store.selectAgent($0) }
        )
    }
}
