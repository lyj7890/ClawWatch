# Changelog

## v1.0.12 (2026-06-05)

### New Features
- **最近访问栏**：Sidebar 顶部新增「🕐 最近访问」区块，显示最近有消息活动的 Agent（按最后消息时间排序，最多 5 个），支持彩色圆点区分、相对时间显示，点击直接切换


- **feat: Sidebar 两级 agent → session 结构**
  - Session tab 从 header 移入 sidebar，按 agent 主机名分组展示
  - 每个 agent 分组可折叠/展开，默认折叠
  - Session 列表无数量限制，可无限滚动（原来最多只显示 5 个）
  - 时间显示改为相对时间（just now / 5m ago / 2h ago / yesterday）
- **feat: Agent 级未读汇总**
  - Agent 分组行右侧显示该 agent 下所有 session 的未读消息总数（红色 badge）
  - 无未读时显示灰色 session 总数
- **feat: Header 精简**
  - 去掉 session tab 行，header 只保留工具按钮，更简洁

## v1.0.10 (2026-06-01)

### Bug Fixes
- **fix: 暂停按钮改用 fixed 定位**
  - 原 `sticky` 在滚动容器内定位不可靠，改为 `fixed bottom-6 left-1/2` 居中显示
  - z-index 提升至 50，确保始终显示在最顶层

## v1.0.9 (2026-06-01)

### New Features
- **feat: 日志自动滚动暂停**
  - 向上滚动查看历史日志时，自动暂停跟踪，不再被新日志顶走
  - 底部显示蓝色悬浮提示条「⏸ 自动滚动已暂停 · 点击恢复」
  - 点击悬浮条跳回最新位置并恢复自动滚动
  - 新消息到来时，若用户正在查看历史则不打断，保持在当前位置

## v1.0.8 (2026-06-01)

### Bug Fixes
- **fix: 兼容 OpenClaw user message content 字符串格式**
  - OpenClaw 新版（5月末起）session jsonl 中 user message 的 `content` 字段格式由数组变为纯字符串，导致用户消息在 ClawWatch 中无法显示
  - 新增 `normalizeContent()` 函数，自动检测 content 类型并统一转为 `[{type, text}]` 数组格式
  - 同时兼容旧格式（数组）和新格式（字符串），历史消息加载与实时 WebSocket 推送均已覆盖

## 2026-05-20
- Webhook integration: push events now trigger automated change summaries via Hermes Agent
