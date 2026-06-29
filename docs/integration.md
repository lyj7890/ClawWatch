# ClawWatch 接入文档

## 概述

将一台新主机接入 ClawWatch Hub，只需在该主机上运行 ClawWatch Agent 即可。Agent 会自动：
1. 扫描本地 OpenClaw session 文件
2. 通过 WebSocket 连接到 Hub
3. 实时推送日志和运行轨迹
4. 获取并打印分配的 Console Token

## 前置条件

- 主机上已安装 OpenClaw（session 目录默认 `~/.openclaw/agents`）
- 网络可达 Hub 地址（`wss://clawhub.intra.mlamp.cn`）
- Hub 的 AgentToken（如果 Hub 配置了 `HUB_AGENT_TOKEN`）

## 安装 Agent

### 方式 1：下载预编译二进制

```bash
# 从 Hub 服务器或内部分发获取二进制
cp clawwatch-agent ~/.local/bin/clawwatch-agent
chmod +x ~/.local/bin/clawwatch-agent
```

### 方式 2：从源码编译

```bash
cd ClawWatch/agent
go build -o clawwatch-agent .
cp clawwatch-agent ~/.local/bin/
```

## 启动 Agent

```bash
# 基本用法
clawwatch-agent --hub wss://clawhub.intra.mlamp.cn

# 如果 Hub 配置了 AgentToken
clawwatch-agent --hub wss://clawhub.intra.mlamp.cn --token <agent-token>

# 指定自定义 OpenClaw agents 目录
clawwatch-agent --hub wss://clawhub.intra.mlamp.cn --dir /path/to/agents
```

### 参数说明

| 参数 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `--hub` | `CLAWWATCH_HUB` | (必需) | Hub WebSocket 地址 |
| `--dir` | `CLAWWATCH_DIR` | `~/.openclaw/agents` | OpenClaw agents 目录 |
| `--trajectory-redact` | — | `true` | 是否脱敏 trajectory 数据 |

### Trajectory 脱敏

默认开启脱敏模式，只上报：
- 事件类型（tool.call, model.request 等）
- 模型名称、Provider
- 标识符哈希

**不会上报**：Prompt、系统 Prompt、工具参数/结果、本地路径。

仅在接受隐私风险时关闭：
```bash
clawwatch-agent --hub wss://clawhub.intra.mlamp.cn --trajectory-redact=false
```

## 主机 ID

首次启动时自动在 `~/.clawwatch/host-id` 生成唯一且稳定的主机 ID。后续重启复用同一 ID，确保 token 不变。

## 获取 Console Token

Agent 连接成功后，Hub 会下发该主机专属的 Console Token，Agent 日志输出：

```
[agent] Console token for this host: d24f34d27d9859dc1fb6798360598f38
```

将此 token 填入 Console 页面即可只查看本主机数据。

## macOS 后台运行（LaunchAgent）

```bash
# 编译并安装
cd ClawWatch/agent
go build -o ~/.local/bin/clawwatch-agent .

# 创建 LaunchAgent plist
cat > ~/Library/LaunchAgents/com.clawwatch.agent.plist << 'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.clawwatch.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/YOU/.local/bin/clawwatch-agent</string>
        <string>--hub</string>
        <string>wss://clawhub.intra.mlamp.cn</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/clawwatch-agent.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/clawwatch-agent.log</string>
</dict>
</plist>
PLIST

# 加载服务
launchctl load ~/Library/LaunchAgents/com.clawwatch.agent.plist

# 查看日志
tail -f /tmp/clawwatch-agent.log
```

## Linux systemd 后台运行

```bash
sudo cat > /etc/systemd/system/clawwatch-agent.service << 'UNIT'
[Unit]
Description=ClawWatch Agent
After=network.target

[Service]
Type=simple
User=YOUR_USER
ExecStart=/usr/local/bin/clawwatch-agent --hub wss://clawhub.intra.mlamp.cn
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
UNIT

sudo systemctl daemon-reload
sudo systemctl enable --now clawwatch-agent
```

## 验证接入

1. Agent 启动后查看日志，确认 `Console token for this host: xxx`
2. 访问 Hub Console，用管理员 token 查看是否出现新主机
3. 用 per-agent token 验证权限隔离

## 断连与重连

Agent 内置自动重连机制（间隔 3 秒）。网络中断或 Hub 重启后会自动恢复连接，token 保持不变。
