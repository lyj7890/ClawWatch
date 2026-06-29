# ClawWatch 🦞

AI Agent 实时监控平台 — 让你随时随地掌握 Agent 运行状态。

## 项目总览

ClawWatch 包含三种运行模式和一个桌面宠物：

| 模块 | 路径 | 语言 | 说明 |
|------|------|------|------|
| **本地模式** | `./` (根目录) | Node.js | 单机 Web 监控，直接读本地 session 文件 |
| **Hub 模式** | `hub/` + `agent/` | Go | 多机远程监控，Agent 推送到中心 Hub |
| **官网** | `website/` | React + TS | 项目官方介绍页 |
| **桌面宠物** | `macos/ClawPet/` | Swift | macOS 桌面置顶宠物，展示实时活动 |

---

## 1. 本地模式（Node.js）

**适用场景**：单台机器上快速查看本地 OpenClaw/Hermes 的会话日志。

```bash
npx @lyj7890/claw-watch        # OpenClaw 模式
npx @lyj7890/claw-watch start --hermes  # Hermes 模式
```

**核心文件**：
- `web-viewer-server.js` — 后端，HTTP + WebSocket
- `frontend/` — Vue 3 前端
- `bin/clawwatch.js` — CLI 入口
- `lib/agent-pusher.js` — Agent 推送辅助
- `Dockerfile` + `docker-compose.yml` — 容器化部署

**特点**：
- 端口 3939，前后端合一
- 实时 WebSocket 推送
- 直接读取 `~/.openclaw/agents/` 目录
- npm 包发布为 `@lyj7890/claw-watch`

---

## 2. Hub 模式（Go）

**适用场景**：多台机器集中监控，手机/远程查看 Agent 状态。

### 架构

```
[主机 A] agent → 
[主机 B] agent →  ──WebSocket──→  [Hub 服务器]  ←──WebSocket──  [浏览器 Console]
[主机 C] agent → 
```

### Agent (`agent/`)

轻量级采集端，运行在每台被监控主机上：
- 扫描本地 OpenClaw session 文件
- 增量推送日志和 trajectory 到 Hub
- 自动获取 per-host Token

```bash
cd agent && go build -o clawwatch-agent .
./clawwatch-agent --hub wss://clawwatch.intra.mlamp.cn
```

### Hub (`hub/`)

中心服务器，接收所有 Agent 数据，按权限分发给 Console：
- 内置 Web UI（无需独立前端）
- Token 鉴权（per-host 隔离 + admin 全局）
- 消息持久化（JSONL）
- Session 历史回溯

```bash
cd hub && go build -o clawwatch-hub .
./clawwatch-hub --port 4848 --console-token YOUR_SECRET
```

### 部署信息

| 服务 | 集群 | 域名 |
|------|------|------|
| Hub | tbj6 / openclaw | clawatch.intra.mlamp.cn |
| 官网 | tbj6 / openclaw | clawwatch.intra.mlamp.cn |

---

## 3. 官网 (`website/`)

项目官方介绍页面，dark tech 风格。

- **技术栈**：Vite + React + TypeScript + Tailwind v4 + motion/react
- **构建**：`cd website && npm run build`
- **部署**：Nginx 静态托管（Dockerfile 已包含）

---

## 4. 桌面宠物 (`macos/ClawPet/`)

macOS 桌面置顶宠物，连接本地 ClawWatch 服务，以动画形式展示 Agent 实时活动。

- **语言**：Swift 6 / SwiftPM
- **平台**：macOS 14+
- **连接**：`http://localhost:3939`（需先启动本地模式）
- **功能**：活动摘要、工具调用展示、回复预览

```bash
./script/build_and_run.sh
```

---

## 目录结构

```
ClawWatch/
├── bin/                    # CLI 入口（Node.js 本地模式）
├── lib/                    # Node.js 辅助模块
├── frontend/               # Vue 3 前端（本地模式 UI）
├── web-viewer-server.js    # Node.js 后端（本地模式）
├── Dockerfile              # 本地模式容器镜像
├── docker-compose.yml      # 本地模式 compose
├── package.json            # npm 包定义
│
├── hub/                    # Go Hub 服务（远程模式）
│   ├── main.go
│   ├── hub.go              # 核心 fanout 逻辑
│   ├── agent.go            # Agent 连接管理
│   ├── console.go          # Console 连接管理
│   ├── token_store.go      # Token 持久化
│   ├── config.go           # 配置
│   ├── message_store.go    # 消息持久化
│   ├── index_html.go       # 内置前端
│   ├── Dockerfile          # Hub 镜像
│   └── Makefile
│
├── agent/                  # Go Agent 采集端（远程模式）
│   ├── main.go
│   ├── trajectory.go       # Trajectory 采集
│   └── Makefile
│
├── website/                # 官网（React）
│   ├── src/
│   ├── Dockerfile
│   └── package.json
│
├── macos/ClawPet/          # macOS 桌面宠物（Swift）
│   ├── Package.swift
│   ├── Sources/
│   └── Tests/
│
├── k8s/                    # K8s 部署清单
│   ├── clawwatch-hub.yaml
│   ├── clawwatch-hub-ingress.yaml
│   └── clawwatch-website.yaml
│
├── docs/                   # 文档
│   ├── usage.md            # 使用文档
│   ├── integration.md      # 接入文档
│   ├── architecture.md     # 架构与原理
│   ├── deployment.md       # 部署指南
│   └── auth.md             # Token 鉴权
│
├── script/                 # 构建脚本
├── tests/                  # 测试
├── install.sh              # Agent 一键安装脚本
├── .github/workflows/      # CI: tag 触发自动构建 Release
├── CHANGELOG.md            # 版本记录
└── README.md               # 本文件
```

---

## 两种模式对比

| | 本地模式 | Hub 模式 |
|---|---|---|
| 语言 | Node.js | Go |
| 部署 | `npx` 或 Docker | K8s / 二进制 |
| 数据源 | 直接读本地文件 | Agent 远程推送 |
| 多机支持 | ❌ 单机 | ✅ 无限 Agent |
| Token 鉴权 | ❌ | ✅ per-host 隔离 |
| Trajectory | ❌ | ✅ 运行轨迹 |
| npm 包 | `@lyj7890/claw-watch` | — |
| 适合场景 | 开发调试 | 生产/团队 |

---

## 快速开始

### 本地开发用

```bash
npx @lyj7890/claw-watch
# 访问 http://localhost:3939
```

### 团队/生产用

```bash
# 安装 Agent
curl -fsSL https://raw.githubusercontent.com/lyj7890/ClawWatch/main/install.sh | sh

# 连接到 Hub
clawwatch-agent --hub wss://clawwatch.intra.mlamp.cn
```

---

## 仓库地址

- **GitHub**: https://github.com/lyj7890/ClawWatch
- **GitLab**: https://codex.mlamp.cn/paas/ClawWatch
- **npm**: https://www.npmjs.com/package/@lyj7890/claw-watch

## 许可证

MIT
