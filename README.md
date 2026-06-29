# ClawWatch 🦞

AI Agent 实时监控平台 — 让你随时随地掌握 Agent 运行状态。

## 项目总览

ClawWatch 包含 Hub 服务、Agent 采集端、本地模式和桌面宠物：

| 模块 | 路径 | 语言 | 说明 |
|------|------|------|------|
| **Hub 服务** | `./` (根目录) | Go | 中心服务器，接收 Agent 数据并分发给 Console |
| **Agent 采集** | `agent/` | Go | 运行在每台被监控主机，推送日志到 Hub |
| **本地模式** | `local/` | Node.js | 单机 Web 监控，直接读本地 session 文件 |
| **官网** | `website/` | React + TS | 项目官方介绍页 |
| **桌面宠物** | `desktop-pet/` | Swift | macOS 桌面置顶宠物 |

---

## 1. Hub 服务（Go，根目录）

**项目主体**。中心服务器，接收所有 Agent 上报的数据，按权限分发给 Console。

```bash
go build -o clawwatch-hub .
./clawwatch-hub --port 4848 --console-token YOUR_SECRET
```

**特点**：
- 内置 Web UI（无需独立前端）
- Token 鉴权（per-host 隔离 + admin 全局）
- 消息持久化（JSONL）
- Session 历史回溯

---

## 2. Agent 采集端（`agent/`）

轻量级采集器，运行在每台被监控主机上：
- 扫描本地 OpenClaw session 文件
- 增量推送日志和 trajectory 到 Hub
- 自动获取 per-host Token

```bash
cd agent && go build -o clawwatch-agent .
./clawwatch-agent --hub wss://clawwatch.intra.mlamp.cn
```

---

## 3. 本地模式（`local/`，Node.js）

**适用场景**：单台机器上快速查看本地 OpenClaw/Hermes 的会话日志。

```bash
npx @lyj7890/claw-watch        # OpenClaw 模式
npx @lyj7890/claw-watch start --hermes  # Hermes 模式
```

**核心文件**：
- `local/web-viewer-server.js` — 后端，HTTP + WebSocket
- `local/frontend/` — Vue 3 前端
- `local/bin/clawwatch.js` — CLI 入口
- `local/Dockerfile` + `local/docker-compose.yml` — 容器化部署

**特点**：
- 端口 3939，前后端合一
- 实时 WebSocket 推送
- 直接读取 `~/.openclaw/agents/` 目录
- npm 包发布为 `@lyj7890/claw-watch`

---

---

## 4. 官网 (`website/`)

项目官方介绍页面，dark tech 风格。

- **技术栈**：Vite + React + TypeScript + Tailwind v4 + motion/react
- **构建**：`cd website && npm run build`
- **部署**：Nginx 静态托管（Dockerfile 已包含）

---

## 5. 桌面宠物（`desktop-pet/`）

macOS 桌面置顶宠物，连接本地 ClawWatch 服务，以动画形式展示 Agent 实时活动。

- **语言**：Swift 6 / SwiftPM
- **平台**：macOS 14+
- **连接**：`http://localhost:3939`（需先启动本地模式）

```bash
cd desktop-pet && ./script/build_and_run.sh
```

---

## 目录结构

```
ClawWatch/
├── main.go                 # Hub 入口
├── hub.go                  # 核心 fanout 逻辑
├── agent.go                # Agent 连接管理
├── console.go              # Console 连接管理
├── token_store.go          # Token 持久化
├── config.go               # Hub 配置
├── message_store.go        # 消息持久化
├── index_html.go           # Hub 内置前端
├── embed_assets.go         # 静态资源嵌入
├── util.go                 # 工具函数
├── go.mod / go.sum         # Go 依赖
├── Makefile                # 构建命令
├── Dockerfile              # Hub 镜像（多阶段构建）
│
├── agent/                  # Agent 采集端（独立 Go module）
│   ├── main.go
│   ├── trajectory.go
│   ├── go.mod / go.sum
│   └── Makefile
│
├── local/                  # 本地模式（Node.js 单机版）
│   ├── web-viewer-server.js
│   ├── bin/clawwatch.js
│   ├── lib/
│   ├── frontend/           # Vue 3 前端
│   ├── package.json
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── desktop-pet/            # macOS 桌面宠物（Swift）
│   ├── Package.swift
│   ├── Sources/
│   ├── Tests/
│   └── script/
│
├── website/                # 官网（React + Tailwind）
│   ├── src/
│   ├── Dockerfile
│   └── package.json
│
├── k8s/                    # K8s 部署清单
│   ├── clawwatch-hub.yaml
│   ├── clawwatch-hub-ingress.yaml
│   └── clawwatch-website.yaml
│
├── docs/                   # 文档
│   ├── usage.md
│   ├── integration.md
│   ├── architecture.md
│   ├── deployment.md
│   └── auth.md
│
├── install.sh              # Agent 一键安装脚本
├── .github/workflows/      # CI: tag 触发自动构建 Release
├── AGENTS.md               # Codex 协作指引
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

## 部署信息

| 服务 | 集群 | 域名 |
|------|------|------|
| Hub Console | tbj6 / openclaw | clawatch.intra.mlamp.cn |
| 官网 | tbj6 / openclaw | clawwatch.intra.mlamp.cn |

---

## 仓库地址

- **GitHub**: https://github.com/lyj7890/ClawWatch
- **GitLab**: https://codex.mlamp.cn/paas/ClawWatch
- **npm**: https://www.npmjs.com/package/@lyj7890/claw-watch

## 许可证

MIT
