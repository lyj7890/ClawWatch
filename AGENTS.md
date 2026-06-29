# AGENTS.md - ClawWatch

> OpenClaw / Hermes Agent 实时会话监控器

## 项目概述

ClawWatch 是 OpenClaw/Hermes Agent 的实时会话监控工具，提供：
- Web UI 实时查看会话日志和对话历史
- HUB 端持久化存储 trajectory 数据（`.trajectory.jsonl`）
- Agent 端上报机制（主机IP、主机名、agent/session ID）
- Conversation 视图整合运行时信息（模型调用、Token 用量、上下文编译等）
- Runtime Timeline 展示完整事件流
- 支持 OpenClaw 和 Hermes Agent 两种模式

**部署版本**: v1.0.18 (tbj6 cluster, openclaw namespace)  
**镜像仓库**: `hub.intra.mlamp.cn/k8s/clawwatch-hub`

## 技术栈

✅ Go (HUB 后端: `main.go`, `trajectory.go`)  
✅ Node.js + Vue/Vite (前端 UI)  
✅ Docker + Docker Compose (本地开发)  
✅ Kubernetes manifests (`k8s/clawwatch-hub.yaml`)  
✅ Trajectory 数据采集（`.trajectory.jsonl`）

## 运行环境

**前置要求**:
- Node.js (前端构建)
- Go 1.20+ (HUB 后端)
- Docker / Podman (容器化)
- K8s 集群（生产部署）

**关键组件**:
- **HUB**: Go 后端服务，接收并持久化 agent 上报数据
- **Agent Pusher**: Agent 端采集 `.trajectory.jsonl` 并上报到 HUB
- **Web Viewer**: Vue 前端，展示实时监控界面

## 常用命令

```bash
# 本地启动（OpenClaw 模式）
npx @lyj7890/claw-watch

# 本地启动（Hermes 模式）
npx @lyj7890/claw-watch start --hermes

# Docker Compose 本地开发
docker-compose up -d

# 构建 HUB 镜像
cd hub
go build -o clawwatch-hub-linux-amd64 .
docker build -f Dockerfile.prebuilt -t hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.18 .

# K8s 部署
kubectl apply -f k8s/clawwatch-hub.yaml
kubectl rollout status deploy/clawwatch-hub -n openclaw

# 测试
cd hub && go test ./...
```

## SOP 流程

### Agent 端配置

1. **启用 Trajectory 采集**:
   - Agent 启动时添加参数启用轨迹记录
   - 生成 `.trajectory.jsonl` 文件（每个 session 一个文件）

2. **配置上报参数**:
   ```bash
   # ❌ 删除不必要的参数
   --no-proxy              # 本地环境特殊，不需要
   --insecure-skip-verify  # 生产环境不应跳过证书验证
   
   # ✅ 保留核心参数
   --hub-url <HUB地址>
   --agent-id <唯一标识>
   ```

3. **隐私保护**:
   - 启用数据脱敏参数控制敏感信息上报
   - 用户只能看到自己的 agent 内容

### HUB 部署流程

1. **构建镜像**:
   ```bash
   cd hub
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o clawwatch-hub-linux-amd64 .
   docker build --platform linux/amd64 -f Dockerfile.prebuilt -t hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.18 .
   ```

2. **推送镜像**:
   ```bash
   docker push hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.18
   ```

3. **K8s 部署**:
   ```bash
   kubectl --context tbj6 apply -f k8s/clawwatch-hub.yaml
   kubectl --context tbj6 -n openclaw rollout status deploy/clawwatch-hub --timeout=180s
   kubectl --context tbj6 -n openclaw get pods -l app=clawwatch-hub
   ```

4. **验证部署**:
   ```bash
   # 检查 Pod 状态
   kubectl get pods -n openclaw -l app=clawwatch-hub
   
   # 健康检查
   curl http://<pod-ip>:4848/api/health
   ```

### 前端优化策略

**多 Session 管理**:
- 点击主机或 agent 时折叠/展开（而非跳转空白页）
- 限制单主机展示的总 session 数量上限：**15 个**
- 最近活动时间预置：**4 小时**
- 右侧数字显示未读消息数（总消息数无意义）

**Conversation 视图**:
- Session 级信息只在顶部出现一次（避免重复）
- Assistant 卡片底部聚合：模型调用状态、耗时、Token/cache、上下文消息数、可用工具数
- 未关联事件只显示数量，引导到 Runtime Timeline

**Runtime Timeline**:
- 整合的事件类型：`session.started`, `trace.metadata`, `prompt.submitted`, `model.completed`, `context.compiled`, `session.ended`, `trace.artifacts`
- 按时间顺序展示所有事件
- 支持搜索过滤

## ⚠️ 注意事项

❌ **HUB 方案尚未完全使用** - 目前主要是本地 UI，远端 HUB 需进一步推广
❌ **多 Agent/Session 混淆** - 多个 agent 同时会话时需明确区分，通过主机名/agent ID 标识
⚠️ **隐私保护** - Trajectory 数据包含详细调用信息，必须启用脱敏控制
⚠️ **持久化缺失** - HUB 端早期版本刷新后消息丢失，v1.0.18 已修复
⚠️ **Session 排序** - 选择 session 后自动排到第一个的行为已调整，改为折叠/展开
⚠️ **数据脱敏** - 采集时需支持尽量多的指标，但通过参数控制是否脱敏

## 已知坑

1. **消息丢失** - 早期版本 HUB 无持久化，刷新页面后数据消失（已在 v1.0.18 修复）
2. **Session 混乱** - 多个 agent/session 同时会话时，点击主机/agent 会揉到一起（改为折叠/展开解决）
3. **未读消息显示** - 右侧数字原本显示总消息数，应改为未读消息数
4. **Agent 参数冗余** - `no-proxy` 和 `insecure-skip-verify` 在生产环境不应使用
5. **Trajectory 数据完整性** - `.trajectory.jsonl` 记录了完整的工具调用、模型切换、参数等，需确认是否全部上报
6. **页面展示优化** - 前端需避免重复展示信息，同时不能遗漏关键数据
