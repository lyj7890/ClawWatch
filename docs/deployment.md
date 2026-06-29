# ClawWatch Hub 部署文档

## 概述

ClawWatch Hub 是一个轻量级 Go 服务，支持多种部署方式。

## 编译

```bash
cd hub/

# 本地编译
make build

# 交叉编译 Linux
make build-linux
# 产出: clawwatch-hub-linux-amd64, clawwatch-hub-linux-arm64
```

## 配置项

所有配置通过命令行参数或环境变量设置：

| 参数 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `--port` | `HUB_PORT` | `4848` | 监听端口 |
| `--agent-token` | `HUB_AGENT_TOKEN` | (空=开放) | Agent 连接鉴权 token |
| `--console-token` | `HUB_CONSOLE_TOKEN` | (空=开放模式) | 全局管理员 Console token |
| `--storage-path` | `HUB_STORAGE_PATH` | `data/messages.jsonl` | 日志持久化路径 |
| `--history-limit` | `HUB_HISTORY_LIMIT` | `2000` | 最大持久化日志条数 |
| `--metrics-path` | `HUB_METRICS_PATH` | `data/metrics.jsonl` | 轨迹持久化路径 |
| `--metrics-limit` | `HUB_METRICS_LIMIT` | `10000` | 最大持久化轨迹条数 |
| `--token-store-path` | `HUB_TOKEN_STORE_PATH` | `data/tokens.json` | Token 存储路径 |

### 安全建议

**生产环境必须配置：**
- `HUB_CONSOLE_TOKEN` — 否则任何人都可以查看所有 Agent 数据（开放模式）
- `HUB_AGENT_TOKEN` — 否则任何人都可以注册为 Agent 推送数据

## Docker 部署

```bash
# 构建镜像
docker build -t clawwatch-hub:latest -f hub/Dockerfile hub/

# 运行
docker run -d \
  --name clawwatch-hub \
  -p 4848:4848 \
  -v clawwatch-data:/app/data \
  -e HUB_AGENT_TOKEN=your-agent-secret \
  -e HUB_CONSOLE_TOKEN=your-console-secret \
  -e HUB_STORAGE_PATH=/app/data/messages.jsonl \
  -e HUB_METRICS_PATH=/app/data/metrics.jsonl \
  -e HUB_TOKEN_STORE_PATH=/app/data/tokens.json \
  clawwatch-hub:latest
```

## Kubernetes 部署

### 前置条件

- Harbor 镜像仓库：`hub.intra.mlamp.cn/k8s/clawwatch-hub`
- Namespace：`openclaw`
- Ingress：`clawhub.intra.mlamp.cn`

### 部署步骤

```bash
# 1. 构建并推送镜像
cd hub/
docker build --platform linux/amd64 \
  -t hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.29 \
  -f Dockerfile .
docker push hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.29

# 2. 部署（首次）
kubectl --context tbj6 apply -f k8s/clawwatch-hub.yaml
kubectl --context tbj6 apply -f k8s/clawwatch-hub-ingress.yaml

# 3. 配置环境变量
kubectl --context tbj6 -n openclaw set env deploy/clawwatch-hub \
  HUB_AGENT_TOKEN=<agent-token> \
  HUB_CONSOLE_TOKEN=<admin-token> \
  HUB_STORAGE_PATH=/app/data/messages.jsonl \
  HUB_HISTORY_LIMIT=2000 \
  HUB_METRICS_PATH=/app/data/metrics.jsonl \
  HUB_METRICS_LIMIT=10000

# 4. 更新镜像（后续版本升级）
kubectl --context tbj6 -n openclaw set image deploy/clawwatch-hub \
  clawwatch-hub=hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.29

# 5. 查看部署状态
kubectl --context tbj6 -n openclaw rollout status deploy/clawwatch-hub
```

### K8s 资源清单

```yaml
# k8s/clawwatch-hub.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clawwatch-hub
  namespace: openclaw
spec:
  replicas: 1
  selector:
    matchLabels:
      app: clawwatch-hub
  template:
    spec:
      containers:
        - name: clawwatch-hub
          image: hub.intra.mlamp.cn/k8s/clawwatch-hub:v1.0.29
          ports:
            - containerPort: 4848
          resources:
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 200m
              memory: 256Mi
          readinessProbe:
            httpGet:
              path: /api/health  # 注意：不要用 /api/agents，鉴权开启后会 401
              port: 4848
          livenessProbe:
            httpGet:
              path: /api/health
              port: 4848
```

### ⚠️ 注意事项

1. **readinessProbe 必须用 `/api/health`**，不能用 `/api/agents`。因为开启 `HUB_CONSOLE_TOKEN` 后，`/api/agents` 需要 token 才能访问，探针不带 token 会导致 Pod 永远不 Ready。

2. **数据持久化**：当前部署未挂载 PVC，Pod 重建后 `tokens.json` 会丢失（agent 重连后自动重新生成）。如需持久化，添加 PVC 挂载到 `/app/data`。

3. **Ingress WebSocket 支持**：Ingress 需要配置超长超时（已设置 3600s）以支持 WebSocket 长连接。

## 运维命令

```bash
# 查看 Pod 状态
kubectl --context tbj6 -n openclaw get pods -l app=clawwatch-hub

# 查看日志
kubectl --context tbj6 -n openclaw logs -f deploy/clawwatch-hub

# 查看 token 存储
kubectl --context tbj6 -n openclaw exec deploy/clawwatch-hub -- cat /app/data/tokens.json

# 健康检查
curl https://clawhub.intra.mlamp.cn/api/health
```

## 版本发布流程

1. 修改代码，测试通过
2. `make build-linux` 交叉编译
3. `docker build` + `docker push` 推镜像
4. `kubectl set image` 滚动更新
5. `kubectl rollout status` 确认部署成功
