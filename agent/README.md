# ClawWatch Go Agent

ClawWatch Agent scans local OpenClaw session files and reports host, OpenClaw Agent, Session, and real-time log events to ClawWatch Hub.

## Build

```bash
make build
```

当前开发阶段只构建并测试本机二进制。

## Run

```bash
./clawwatch-agent --hub wss://clawatch.intra.mlamp.cn
```

参数：

| Option | Description |
|---|---|
| `--hub` | Hub 地址，必需 |
| `--dir` | OpenClaw agents directory, defaults to `~/.openclaw/agents` |
| `--trajectory-redact` | 是否脱敏 trajectory 数据，默认 `true` |

主机 ID 会在首次运行时自动生成并保存到 `~/.clawwatch/host-id`。hostname、IP、OS 和架构会自动采集。

Agent 会增量解析 `.trajectory.jsonl`，将其作为 session 运行事件上报。默认脱敏模式只
上报事件类型、模型、Provider 和标识符哈希；不会上报 Prompt、系统 Prompt、完整上下文、
工具参数/结果、本地路径或原始 `data`。

仅在用户明确接受隐私风险时关闭脱敏：

```bash
./clawwatch-agent --hub wss://clawatch.intra.mlamp.cn --trajectory-redact=false
```

关闭后会额外上报原始 trajectory `data`、路径及原始标识符，用于还原完整 session 运行过程。

For the local macOS LaunchAgent:

```bash
make build
cp clawwatch-agent ~/.local/bin/clawwatch-agent
launchctl kickstart -k gui/$(id -u)/com.clawwatch.agent
```

Logs:

```bash
tail -f /tmp/clawwatch-agent.log
```
