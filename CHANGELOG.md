# Changelog

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
