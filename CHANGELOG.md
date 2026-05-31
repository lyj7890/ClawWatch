# Changelog

## v1.0.8 (2026-06-01)

### Bug Fixes
- **fix: 兼容 OpenClaw user message content 字符串格式**
  - OpenClaw 新版（5月末起）session jsonl 中 user message 的 `content` 字段格式由数组变为纯字符串，导致用户消息在 ClawWatch 中无法显示
  - 新增 `normalizeContent()` 函数，自动检测 content 类型并统一转为 `[{type, text}]` 数组格式
  - 同时兼容旧格式（数组）和新格式（字符串），历史消息加载与实时 WebSocket 推送均已覆盖

## 2026-05-20
- Webhook integration: push events now trigger automated change summaries via Hermes Agent
