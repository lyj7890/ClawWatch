# ClawWatch Token 鉴权机制

## 设计目标

ClawWatch Hub 是一个多租户公共服务，需要解决：

1. **数据隔离** — 每台主机只能看到自己的数据
2. **管理员全局视图** — 运维人员能看到所有主机
3. **向后兼容** — 未配置 token 时保持开放模式
4. **零配置接入** — Agent 连接后自动获得 token，无需预分配

## 鉴权模型

```
                     ┌─────────────────────────────┐
                     │     HUB_CONSOLE_TOKEN 配置？  │
                     └──────────────┬──────────────┘
                                    │
                       ┌────────────┴────────────┐
                       │ 是                       │ 否
                       ▼                         ▼
              ┌────────────────┐        ┌──────────────┐
              │ 鉴权模式        │        │ 开放模式      │
              │ 必须提供 token  │        │ 所有人=admin  │
              └───────┬────────┘        └──────────────┘
                      │
         ┌────────────┼────────────┐
         │            │            │
         ▼            ▼            ▼
   ┌──────────┐ ┌──────────┐ ┌──────────┐
   │ 空/错误   │ │ 全局token │ │ Agent    │
   │ → 401    │ │ → Admin  │ │ token    │
   │ 拒绝访问  │ │ 看所有    │ │ → 只看   │
   └──────────┘ └──────────┘ │   自己   │
                             └──────────┘
```

## Token 类型

### 1. Agent Token (`HUB_AGENT_TOKEN`)

- **用途**: Agent 连接 Hub 时的身份校验
- **配置**: Hub 环境变量 `HUB_AGENT_TOKEN`
- **传递**: Agent 连接时 URL 参数 `?token=<agent-token>`
- **不配置时**: 任何人可注册为 Agent

### 2. Console Token (`HUB_CONSOLE_TOKEN`)

- **用途**: 全局管理员 token，拥有查看所有 Agent 数据的权限
- **配置**: Hub 环境变量 `HUB_CONSOLE_TOKEN`
- **不配置时**: 开放模式，所有请求视为 admin

### 3. Per-Agent Token (自动生成)

- **用途**: 每台主机的专属 token，只能查看该主机的数据
- **生成**: Agent 首次连接时 Hub 自动生成（16 字节随机 hex = 32 字符）
- **存储**: Hub 端 `data/tokens.json`
- **分发**: 通过 `agent_hello_ack` 消息回传给 Agent
- **持久化**: 重启后复用，token 不变

## 数据流

### Agent 注册获取 Token

```
Agent                           Hub
  │                              │
  │── ws connect ──────────────→ │
  │   ?agentId=my-host           │
  │                              │
  │                   tokens.GetOrCreate("my-host")
  │                   → 生成 "abc123..." 并持久化
  │                              │
  │←── agent_hello_ack ─────────│
  │    {consoleToken: "abc123"}  │
  │                              │
  │ [打印到日志供用户获取]         │
```

### Console 鉴权流程

```
Console                         Hub
  │                              │
  │── ws connect ──────────────→ │
  │   ?token=abc123              │
  │                              │
  │              resolveScope(token)
  │              → token == HUB_CONSOLE_TOKEN? → admin
  │              → tokens.AgentIDForToken(token)? → per-agent
  │              → 都不匹配? → 401 拒绝
  │                              │
  │←── agent_status (filtered) ──│
  │    只推送有权限的 Agent 数据    │
```

### REST API 鉴权

所有 REST API（`/api/agents`、`/api/history`、`/api/trajectory`）使用相同的 `resolveScope` 逻辑：

```
GET /api/agents?token=abc123

→ resolveScope: abc123 匹配 agentId "my-host"
→ scope = {admin: false, allowedAgentID: "my-host"}
→ 只返回 "my-host" 的数据
```

## 权限矩阵

| 操作 | Admin | Per-Agent | 无 Token |
|------|-------|-----------|----------|
| 查看所有 Agent | ✅ | ❌ | ❌ (401) |
| 查看自己的 Agent | ✅ | ✅ | ❌ (401) |
| 查看他人的 Agent | ✅ | ❌ (403) | ❌ (401) |
| 查看历史日志 | ✅ 全部 | ✅ 仅自己 | ❌ (401) |
| 查看轨迹事件 | ✅ 全部 | ✅ 仅自己 | ❌ (401) |
| 获取 session 历史 | ✅ 全部 | ✅ 仅自己 | ❌ (401) |
| WebSocket 订阅 | ✅ 任意 | 🔒 强制锁定自己 | ❌ (401) |

## 实现细节

### Token 存储 (`token_store.go`)

```go
type TokenStore struct {
    mu     sync.RWMutex
    path   string            // data/tokens.json
    tokens map[string]string // agentId -> token
}
```

- **并发安全**: RWMutex 读写分离
- **原子写入**: 先写 `.tmp` 再 `rename`，防止断电丢数据
- **失败回滚**: 磁盘写入失败时删除内存中的新 token
- **双重检查**: GetOrCreate 使用 double-check locking

### Console 权限锁定 (`console.go`)

非 admin 的 Console 连接会被强制锁定 `subscribe` 为自己的 agentId：
```go
if !admin {
    subscribe = allowedAgentID
}
```

### Fanout 过滤 (`hub.go`)

消息广播时按权限过滤：
```go
func (h *Hub) fanout(msg *broadcastMsg) {
    for _, c := range h.consoles {
        if !c.canSee(msg.agentID) {
            continue  // 跳过无权限的 console
        }
        // ...
    }
}
```

## 已知限制

1. **Token 反查 O(n)**: `AgentIDForToken()` 线性扫描所有 token。当前 Agent 数量少（<100）性能无影响，规模增大后需加 reverse map。

2. **Token 不可撤销**: 目前无删除 token 的接口。如需撤销，手动删除 `tokens.json` 中的条目并重启 Hub。

3. **无 PVC 时 Token 丢失**: Pod 重建后 `tokens.json` 丢失，Agent 重连会生成新 token（旧 token 失效）。
