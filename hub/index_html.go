package main

func buildIndexHTML(cfg *Config) string {
	part1 := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ClawWatch Hub</title>
<script>`
	part2 := `</script>
<script>`
	part3 := `</script>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html, body { height: 100%; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8fafc; color: #1e293b; }
button, input { font: inherit; }
.topbar { height:56px;background:#fff;border-bottom:1px solid #e2e8f0;padding:0 20px;display:flex;align-items:center;justify-content:space-between;flex-shrink:0; }
.brand-mark { width:30px;height:30px;border-radius:9px;background:#fff7ed;color:#c2410c;display:grid;place-items:center;font-size:17px; }
.status-pill { display:flex;align-items:center;gap:7px;padding:5px 9px;border-radius:999px;background:#f8fafc;border:1px solid #e2e8f0;font-size:12px;color:#64748b; }
.toolbar-btn { padding:7px 11px;background:#fff;border:1px solid #e2e8f0;color:#475569;border-radius:7px;font-size:12px;cursor:pointer;font-weight:600; }
.toolbar-btn:hover { background:#f8fafc;border-color:#cbd5e1; }
.sidebar { width:286px;background:#fff;border-right:1px solid #e2e8f0;display:flex;flex-direction:column;flex-shrink:0;overflow:hidden; }
.sidebar-section-title { padding:14px 16px 7px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:1px;font-weight:700; }
.nav-row { cursor:pointer;transition:background .15s;border-left:3px solid transparent; }
.nav-row:hover { background:#f8fafc; }
.nav-row.active { background:#eff6ff;border-left-color:#3b82f6; }
.count-badge { min-width:22px;text-align:center;font-size:10px;color:#64748b;background:#f1f5f9;padding:2px 6px;border-radius:999px; }
.count-badge.unread { color:#fff;background:#2563eb;font-weight:700; }
.session-row { padding:7px 10px 7px 42px;display:flex;align-items:center;gap:7px;border-top:1px solid #f8fafc;background:#fcfdff; }
.session-row.active { background:#eef6ff;border-left-color:#3b82f6; }
.session-id { font-family:'SF Mono','Fira Code',monospace;font-size:10px;color:#64748b;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1; }
.hidden-sessions { padding:7px 10px 8px 26px;border-top:1px solid #f1f5f9;font-size:9px;color:#94a3b8;background:#f8fafc; }
.content-toolbar { min-height:54px;background:rgba(255,255,255,.92);border-bottom:1px solid #e2e8f0;padding:9px 18px;display:flex;align-items:center;gap:14px;justify-content:space-between;flex-shrink:0; }
.search-input { width:260px;background:#f8fafc;border:1px solid #e2e8f0;border-radius:7px;padding:7px 10px;font-size:12px;color:#1e293b;outline:none; }
.search-input:focus { background:#fff;border-color:#93c5fd;box-shadow:0 0 0 3px #dbeafe; }
.message-list { max-width:1180px;margin:0 auto;display:flex;flex-direction:column;gap:10px; }
.message-card { transition:box-shadow .15s,border-color .15s; }
.message-card:hover { box-shadow:0 5px 18px rgba(15,23,42,.07) !important; }
.message-header { display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap; }
.message-meta { display:flex;align-items:center;gap:6px;flex-wrap:wrap;min-width:0; }
.message-error { padding:8px 14px;background:#fff1f2;border-bottom:1px solid #fecaca;color:#b91c1c;font-size:11px;display:flex;align-items:flex-start;gap:7px;overflow-wrap:anywhere; }
.meta-chip { font-size:10px;color:#64748b;background:rgba(255,255,255,.72);border:1px solid rgba(203,213,225,.75);padding:2px 6px;border-radius:5px;font-family:'SF Mono','Fira Code',monospace;max-width:260px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap; }
.runtime-overview { max-width:1180px;margin:0 auto 10px;padding:10px 12px;background:#f8fafc;border:1px solid #e2e8f0;border-radius:9px;display:flex;align-items:center;gap:8px;flex-wrap:wrap; }
.runtime-label { font-size:10px;color:#64748b;font-weight:700;text-transform:uppercase;letter-spacing:.4px; }
.runtime-chip { font-size:10px;color:#475569;background:#fff;border:1px solid #dbe3ec;padding:3px 7px;border-radius:999px; }
.runtime-chip.ok { color:#047857;border-color:#a7f3d0;background:#ecfdf5; }
.runtime-chip.warn { color:#b45309;border-color:#fde68a;background:#fffbeb; }
.runtime-chip.error { color:#b91c1c;border-color:#fecaca;background:#fff1f2; }
.runtime-footer { margin:10px -15px -13px;padding:8px 15px;background:#f8fafc;border-top:1px solid #e2e8f0;display:flex;align-items:center;gap:6px;flex-wrap:wrap; }
.unassociated-box { max-width:1180px;margin:0 auto 10px;padding:9px 11px;background:#fffbeb;border:1px solid #fde68a;border-radius:8px;font-size:11px;color:#92400e; }
.prompt-context { margin-bottom:10px;background:#faf5ff;border:1px solid #e9d5ff;border-radius:7px;overflow:hidden; }
.prompt-context-btn { width:100%;border:0;background:#faf5ff;color:#6b21a8;padding:8px 10px;display:flex;align-items:center;gap:7px;cursor:pointer;text-align:left;font-size:11px;font-weight:650; }
.prompt-context-body { border-top:1px solid #e9d5ff;padding:9px 10px;display:flex;flex-direction:column;gap:8px; }
.prompt-block-title { font-size:9px;color:#7e22ce;text-transform:uppercase;letter-spacing:.5px;font-weight:700;margin-bottom:4px; }
.prompt-block pre { background:#fff;color:#334155;border:1px solid #e9d5ff;border-radius:6px;padding:8px 9px;font-size:10px;line-height:1.5;white-space:pre-wrap;overflow-wrap:anywhere;max-height:300px;overflow:auto; }

.markdown-body { font-size: 13px; line-height: 1.65; }
.markdown-body h1,.markdown-body h2,.markdown-body h3 { font-weight: 600; margin: 0.5em 0 0.3em; }
.markdown-body h1 { font-size: 1.2em; }
.markdown-body h2 { font-size: 1.1em; }
.markdown-body h3 { font-size: 1em; }
.markdown-body p { margin: 0.3em 0; }
.markdown-body code { background: #f1f5f9; border: 1px solid #e2e8f0; border-radius: 3px; padding: 1px 5px; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 0.85em; color: #0f172a; }
.markdown-body pre { background: #0f172a; border-radius: 6px; padding: 10px 14px; overflow-x: auto; margin: 0.4em 0; }
.markdown-body pre code { background: none; border: none; padding: 0; color: #e2e8f0; font-size: 12px; }
.markdown-body ul,.markdown-body ol { padding-left: 1.4em; margin: 0.3em 0; }
.markdown-body li { margin: 0.1em 0; }
.markdown-body blockquote { border-left: 3px solid #94a3b8; padding-left: 0.75em; color: #64748b; margin: 0.3em 0; }
.markdown-body table { border-collapse: collapse; width: 100%; margin: 0.4em 0; font-size: 12px; }
.markdown-body th,.markdown-body td { border: 1px solid #e2e8f0; padding: 4px 8px; text-align: left; }
.markdown-body th { background: #f1f5f9; font-weight: 600; }
.markdown-body a { color: #3b82f6; text-decoration: none; }
.markdown-body a:hover { text-decoration: underline; }
.markdown-body hr { border: none; border-top: 1px solid #e2e8f0; margin: 0.5em 0; }

::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.4} }
.animate-pulse { animation: pulse 2s infinite; }

btn { display:inline-flex;align-items:center;gap:5px;padding:5px 12px;border-radius:6px;font-size:12px;font-weight:500;cursor:pointer;border:1px solid transparent;transition:background .15s; }
@media (max-width: 900px) {
  .sidebar { width:230px; }
  .content-toolbar { align-items:flex-start;flex-direction:column; }
  .search-input { width:100%; }
}
</style>
</head>
<body>
<div id="app" style="height:100vh;display:flex;flex-direction:column;overflow:hidden">

  <!-- Header -->
  <header class="topbar">
    <div style="display:flex;align-items:center;gap:12px">
      <span class="brand-mark">🦞</span>
      <div>
        <h1 style="font-size:16px;font-weight:700;color:#0f172a;line-height:1.15">ClawWatch Hub</h1>
        <div style="font-size:10px;color:#94a3b8;margin-top:2px">{{ hostCount }} hosts · {{ openClawAgentCount }} agents · {{ sessionCount }} sessions</div>
      </div>
      <div class="status-pill" style="margin-left:8px">
        <div :style="connected ? 'background:#22c55e' : 'background:#ef4444'" style="width:8px;height:8px;border-radius:50%;transition:background .3s"></div>
        <span style="font-size:12px;color:#64748b">{{ connectionStatusText }}</span>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:8px">
      <button @click="toggleAllThinking" class="toolbar-btn">
        {{ allThinkingExpanded ? 'Collapse' : 'Expand' }} thinking
      </button>
      <button @click="toggleAllArgs" class="toolbar-btn">
        {{ allArgsExpanded ? 'Hide' : 'Show' }} tool details
      </button>
      <button @click="showConnectPanel = !showConnectPanel" class="toolbar-btn">
        Settings
      </button>
    </div>
  </header>

  <!-- Connect Panel -->
  <div v-show="showConnectPanel" style="background:#f8fafc;border-bottom:1px solid #e2e8f0;padding:10px 20px;display:flex;align-items:center;gap:10px;flex-shrink:0">
    <label style="font-size:12px;color:#64748b;white-space:nowrap">Token:</label>
    <input v-model="tokenInput" type="password" placeholder="主机 Token (留空=开放模式)" title="输入该主机的 Console Token（一机一 token），或全局 admin token 查看全部" style="width:240px;background:#fff;border:1px solid #e2e8f0;border-radius:6px;padding:5px 10px;font-size:13px;color:#1e293b;outline:none">
    <label style="font-size:12px;color:#64748b;white-space:nowrap">Subscribe:</label>
    <input v-model="subscribeInput" type="text" placeholder="* = all agents" style="width:160px;background:#fff;border:1px solid #e2e8f0;border-radius:6px;padding:5px 10px;font-size:13px;color:#1e293b;outline:none">
    <button @click="reconnect" style="padding:5px 16px;background:#2563eb;color:#fff;border:none;border-radius:6px;font-size:13px;cursor:pointer;font-weight:500">
      {{ connected ? '🔄 Reconnect' : '🔌 Connect' }}
    </button>
    <button @click="clearMessages" style="padding:5px 12px;background:#fff;border:1px solid #fca5a5;color:#dc2626;border-radius:6px;font-size:12px;cursor:pointer;font-weight:500">
      🗑️ Clear
    </button>
  </div>

  <!-- Main -->
  <div style="display:flex;flex:1;overflow:hidden">

    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-section-title">Overview</div>
      <div style="flex:1;overflow-y:auto">
        <div @click="toggleUnreadOnly" class="nav-row" :class="{active:unreadOnly}"
          :style="{cursor: unreadFor('all') || unreadOnly ? 'pointer' : 'default', opacity: unreadFor('all') || unreadOnly ? 1 : .72}"
          style="margin:0 8px 10px;padding:10px 11px;display:flex;align-items:center;gap:9px;font-size:13px;border:1px solid #e2e8f0;border-radius:7px;background:#f8fafc">
          <span style="width:26px;height:26px;border-radius:7px;background:#dbeafe;color:#2563eb;display:grid;place-items:center;font-weight:700">A</span>
          <div style="flex:1">
            <div style="font-weight:650;color:#334155">{{ unreadOnly ? 'Showing unread' : (unreadFor('all') ? 'Show unread sessions' : 'All caught up') }}</div>
            <div style="font-size:10px;color:#94a3b8;margin-top:2px">{{ unreadOnly ? 'Click to show all sessions' : (unreadFor('all') ? 'Filter navigation to unread sessions' : 'No unread sessions') }}</div>
          </div>
          <span v-if="unreadFor('all')" class="count-badge unread">{{ unreadFor('all') }}</span>
        </div>

        <div class="sidebar-section-title" style="border-top:1px solid #f1f5f9">{{ unreadOnly ? 'Hosts with unread' : 'Hosts' }} ({{ hostCount }})</div>
        <div v-for="host in hostGroups" :key="host.id" style="margin:0 8px 8px;border:1px solid #e2e8f0;border-radius:8px;overflow:hidden">
          <div @click="toggleHost(host.id)" class="nav-row"
            style="padding:10px;display:flex;align-items:center;gap:8px">
            <span style="font-size:9px;color:#94a3b8;width:9px">{{isHostExpanded(host.id)?'▼':'▶'}}</span>
            <div :style="host.online ? 'background:#22c55e' : 'background:#94a3b8'" style="width:8px;height:8px;border-radius:50%;flex-shrink:0"></div>
            <div style="overflow:hidden;flex:1">
              <div style="font-size:12px;font-weight:650;color:#334155;overflow:hidden;text-overflow:ellipsis;white-space:nowrap" :title="host.id">{{host.label}}</div>
              <div style="font-size:10px;color:#94a3b8;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{host.subtitle || host.id}}</div>
            </div>
            <span v-if="unreadFor('host', host.id)" class="count-badge unread">{{ unreadFor('host', host.id) }}</span>
          </div>
          <div v-if="isHostExpanded(host.id)" v-for="child in host.children" :key="child.id">
            <div @click="toggleAgent(child.id)" class="nav-row"
              style="padding:8px 10px 8px 26px;display:flex;align-items:center;gap:8px;border-top:1px solid #f1f5f9">
              <span style="font-size:8px;color:#94a3b8;width:8px">{{isAgentExpanded(child.id)?'▼':'▶'}}</span>
              <span style="width:5px;height:5px;border-radius:50%;background:#94a3b8"></span>
              <span style="font-size:11px;font-weight:600;color:#475569;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1" :title="child.id">{{child.name}}</span>
              <span v-if="!isAgentExpanded(child.id) && child.sessions.length" style="font-size:9px;color:#94a3b8">{{child.sessions.length}}</span>
              <span v-if="unreadFor('agent', child.id)" class="count-badge unread">{{ unreadFor('agent', child.id) }}</span>
            </div>
            <div v-if="isAgentExpanded(child.id)" v-for="session in child.sessions" :key="session.id" @click="selectSession(child.id, session.id)"
              class="nav-row session-row" :class="{active:selectedAgent === child.id && selectedSession === session.id}" :title="session.id">
              <span :style="{background:isRecentSession(session)?'#22c55e':'#cbd5e1'}" style="width:5px;height:5px;border-radius:50%;flex-shrink:0"></span>
              <span class="session-id">{{shortSession(session.id)}}</span>
              <span style="font-size:9px;color:#94a3b8;white-space:nowrap">{{relativeTime(session.lastTime)}}</span>
              <span v-if="unreadFor('session', child.id, session.id)" class="count-badge unread">{{ unreadFor('session', child.id, session.id) }}</span>
            </div>
          </div>
          <div v-if="isHostExpanded(host.id) && host.hiddenSessionCount" class="hidden-sessions">
            另有 {{host.hiddenSessionCount}} 个较早会话未展示
          </div>
        </div>
        <div v-if="hostCount === 0" style="padding:24px 14px;font-size:12px;color:#94a3b8;text-align:center">
          <div style="margin-bottom:6px">No agents</div>
          <div style="font-size:11px">Waiting for connections...</div>
        </div>
      </div>
    </aside>

    <div style="flex:1;display:flex;flex-direction:column;min-width:0;position:relative">
    <div class="content-toolbar">
      <div style="min-width:0">
        <div style="font-size:13px;font-weight:700;color:#0f172a;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ selectedViewTitle }}</div>
        <div style="font-size:10px;color:#94a3b8;margin-top:3px">{{ selectedViewSubtitle }}</div>
      </div>
      <input v-model="searchQuery" type="text" placeholder="Search session content" class="search-input">
    </div>
    <main ref="messageArea" style="flex:1;overflow-y:auto;padding:18px 22px;position:relative" @scroll="onMessageScroll">
      <div v-if="filteredMessages.length === 0" style="display:flex;align-items:center;justify-content:center;height:100%;flex-direction:column;gap:12px;color:#94a3b8">
        <span style="font-size:40px">{{ connected ? '💬' : '🔌' }}</span>
        <span style="font-size:14px;font-weight:500">{{ connected ? (selectedSession ? 'No messages in this session' : 'Select a session') : 'Not connected' }}</span>
        <span v-if="connected && !selectedSession" style="font-size:12px">Expand a host and agent, then choose a session</span>
        <span v-if="!connected" style="font-size:12px">Click ⚙️ Settings to configure and connect</span>
        <span v-if="searchQuery && connected" style="font-size:12px">No messages match "{{ searchQuery }}"</span>
      </div>

      <div v-else class="message-list">
        <div v-if="selectedSession && selectedTrajectory.length" class="runtime-overview">
          <span class="runtime-label">Session runtime</span>
          <span v-for="chip in sessionRuntime.chips" :key="chip.label" class="runtime-chip" :class="chip.tone">{{chip.label}}</span>
          <span style="margin-left:auto;font-size:10px;color:#94a3b8">{{runtimeCoverage.integrated}}/{{runtimeCoverage.total}} events integrated</span>
        </div>
        <div v-if="runtimeCoverage.unassociated" class="unassociated-box">
          {{runtimeCoverage.unassociated}} runtime events could not be associated with a message. They remain preserved in the raw trajectory store.
        </div>
        <message-card
          v-for="msg in filteredMessages"
          :key="msg._id"
          :message="msg"
          :show-thinking="true"
          :force-expand-thinking="forceExpandThinking"
          :force-expand-args="forceExpandArgs"
        ></message-card>
        <div ref="bottomAnchor" style="height:1px"></div>
      </div>
    </main>
    <button v-if="newMessageCount > 0" @click="scrollToBottom(true)" style="position:absolute;bottom:12px;left:50%;transform:translateX(-50%);background:#2563eb;color:#fff;border:none;padding:8px 18px;border-radius:999px;font-size:13px;font-weight:600;cursor:pointer;box-shadow:0 2px 8px rgba(0,0,0,.15);z-index:10;display:inline-flex;align-items:center;gap:6px;white-space:nowrap;">↓ {{ newMessageCount }} 条新消息</button>
    </div>
  </div>
</div>

<script>
marked.setOptions({ breaks: true, gfm: true });

// ===== MessageCard Component =====
const MessageCard = {
  name: 'MessageCard',
  props: {
    message: Object,
    showThinking: Boolean,
    forceExpandThinking: Number,
    forceExpandArgs: Number
  },
  data() {
    return { thinkingExpanded: false, promptExpanded: false, argsExpanded: {}, resultExpanded: {} }
  },
  watch: {
    forceExpandThinking(n, o) { if (n < o) this.thinkingExpanded = false },
    forceExpandArgs(n, o) { if (n < o) { this.argsExpanded = {}; this.resultExpanded = {} } }
  },
  computed: {
    role() { return this.message.message?.role || 'unknown' },
    content() {
      const c = this.message.message?.content;
      if (!c) return [];
      const items = typeof c === 'string' ? [{ type: 'text', text: c }] : c;
      if (!this.errorMessage) return items;
      return items.filter(item => !(item.type === 'text' && item.text === '[assistant turn failed before producing content]'));
    },
    errorMessage() { return this.message.message?.errorMessage || null },
    roleIcon() { return {user:'👤',assistant:'🤖',tool:'⚙️',toolResult:'📤'}[this.role] || '❓' },
    roleText() { return {user:'User',assistant:'Assistant',tool:'Tool',toolResult:'Tool Result'}[this.role] || this.role },
    colors() {
      return {
        user:       { card:'#eff6ff', border:'#bfdbfe', hdr:'#dbeafe', lbl:'#1d4ed8' },
        assistant:  this.errorMessage
          ? { card:'#fffafa', border:'#fecaca', hdr:'#fff1f2', lbl:'#b91c1c' }
          : { card:'#ffffff', border:'#e2e8f0', hdr:'#f8fafc', lbl:'#059669' },
        tool:       { card:'#f8fafc', border:'#e2e8f0', hdr:'#f1f5f9', lbl:'#6b7280' },
        toolResult: { card:'#f8fafc', border:'#e2e8f0', hdr:'#f1f5f9', lbl:'#6b7280' }
      }[this.role] || { card:'#fff', border:'#e2e8f0', hdr:'#f8fafc', lbl:'#374151' }
    },
    agentShort() {
      const host = this.message.hostname || this.message.hostId || ''
      const agent = this.message.openclawAgentId || ''
      return agent ? host + ' / ' + agent : host
    },
    sessionShort() {
      const id = this.message.sessionId || ''
      return id ? id.slice(0, 8) : ''
    },
    timestamp() {
      const ts = this.message.timestamp
      if (!ts) return ''
      return new Date(ts).toLocaleString('zh-CN',{month:'short',day:'numeric',hour:'2-digit',minute:'2-digit',second:'2-digit'})
    },
    tokenInfo() {
      const u = this.message.message?.usage
      return (u && u.totalTokens) ? (u.totalTokens + ' tokens') : null
    },
    runtime() { return this.message._runtime || null },
    promptContext() { return this.runtime?.promptContext || null },
    displayModel() { return this.runtime?.model || this.message.message?.model || '' },
    isThinkingExpanded() { return this.thinkingExpanded || this.forceExpandThinking > 0 }
  },
  methods: {
    isArgExpanded(i) { return !!this.argsExpanded[i] || this.forceExpandArgs > 0 },
    isResultExpanded(i) { return !!this.resultExpanded[i] || this.forceExpandArgs > 0 },
    toggleThinking() { this.thinkingExpanded = !this.thinkingExpanded },
    togglePrompt() { this.promptExpanded = !this.promptExpanded },
    hasArgs(item) { return item.arguments && Object.keys(item.arguments).length > 0 },
    toggleArgs(i) { this.argsExpanded = {...this.argsExpanded, [i]: !this.argsExpanded[i]} },
    toggleResult(i) { this.resultExpanded = {...this.resultExpanded, [i]: !this.resultExpanded[i]} },
    renderMd(text) { return text ? marked.parse(text) : '' },
    toolStyle(item) {
      if (item.result && item.isError) return { bg:'#fff1f2', border:'#fca5a5', lbl:'#991b1b' }
      if (item.result) return { bg:'#f0fdf4', border:'#86efac', lbl:'#15803d' }
      return { bg:'#fffbeb', border:'#fde68a', lbl:'#92400e' }
    },
    toolIcon(item) {
      if (item.result && item.isError) return '❌'
      if (item.result) return '✅'
      return '🔧'
    }
  },
  template: ` + "`" + `
<div class="message-card" :style="{background:colors.card, border:'1px solid '+colors.border, borderRadius:'9px', boxShadow:'0 1px 3px rgba(15,23,42,.05)', overflow:'hidden'}">
  <!-- Header -->
  <div class="message-header" :style="{background:colors.hdr, borderBottom:'1px solid '+colors.border, padding:'8px 14px'}">
    <div class="message-meta" style="font-size:13px">
      <span>{{roleIcon}}</span>
      <span :style="{fontWeight:600,color:colors.lbl}">{{roleText}}</span>
      <span v-if="agentShort" class="meta-chip" :title="agentShort">{{agentShort}}</span>
      <span v-if="sessionShort" class="meta-chip">session {{sessionShort}}</span>
    </div>
    <div class="message-meta" style="justify-content:flex-end;font-size:10px;color:#94a3b8">
      <span v-if="displayModel">{{displayModel}}</span>
      <span v-if="tokenInfo">{{tokenInfo}}</span>
      <span>{{timestamp}}</span>
    </div>
  </div>
  <div v-if="errorMessage" class="message-error">
    <span>⚠️</span>
    <div>
      <strong>Model attempt failed</strong>
      <span v-if="message._recoveredFromTrajectory"> · recovered from trajectory</span>
      <span v-if="message.message.provider || message.message.model"> · {{message.message.provider || 'unknown provider'}} / {{message.message.model || 'unknown model'}}</span>
      <div style="margin-top:3px;font-family:'SF Mono','Fira Code',monospace">{{errorMessage}}</div>
    </div>
  </div>
  <!-- Content -->
  <div style="padding:13px 15px">
    <div v-if="promptContext" class="prompt-context">
      <button @click="togglePrompt" class="prompt-context-btn">
        <span>{{promptExpanded?'▼':'▶'}}</span>
        <span>Prompt context</span>
        <span v-for="label in promptContext.labels" :key="label" class="runtime-chip">{{label}}</span>
      </button>
      <div v-if="promptExpanded" class="prompt-context-body">
        <div v-if="promptContext.systemPrompt" class="prompt-block">
          <div class="prompt-block-title">System prompt</div>
          <pre>{{promptContext.systemPrompt}}</pre>
        </div>
        <div v-if="promptContext.prompt" class="prompt-block">
          <div class="prompt-block-title">Submitted prompt</div>
          <pre>{{promptContext.prompt}}</pre>
        </div>
        <div v-if="promptContext.messages" class="prompt-block">
          <div class="prompt-block-title">Messages submitted to model</div>
          <pre>{{JSON.stringify(promptContext.messages,null,2)}}</pre>
        </div>
      </div>
    </div>
    <div v-for="(item,idx) in content" :key="idx" style="margin-bottom:4px">
      <!-- Text -->
      <div v-if="item.type==='text' && role==='assistant'" class="markdown-body" v-html="renderMd(item.text)"></div>
      <div v-else-if="item.type==='text'" style="font-size:14px;line-height:1.7;white-space:pre-wrap;color:#334155">{{item.text}}</div>
      <!-- Thinking -->
      <div v-if="item.type==='thinking'" style="margin-top:4px">
        <button @click="toggleThinking" style="display:flex;align-items:center;gap:6px;font-size:12px;color:#2563eb;background:none;border:none;cursor:pointer;padding:2px 0">
          <span :style="{display:'inline-block',transform:isThinkingExpanded?'rotate(90deg)':'none',transition:'transform .2s'}">▶</span>
          <span>💭 Thinking</span>
        </button>
        <div v-if="isThinkingExpanded" style="margin-top:6px;padding:10px 12px;background:#eff6ff;border:1px solid #bfdbfe;border-radius:6px;font-size:12px;color:#374151;max-height:320px;overflow-y:auto;line-height:1.6;white-space:pre-wrap;font-family:'SF Mono','Fira Code',monospace">{{item.thinking||item.text}}</div>
      </div>
      <!-- ToolCall -->
      <div v-if="item.type==='toolCall'" :style="{marginTop:'4px',background:toolStyle(item).bg,border:'1px solid '+toolStyle(item).border,borderRadius:'6px',padding:'8px 10px'}">
        <div style="display:flex;align-items:center;justify-content:space-between">
          <div :style="{display:'flex',alignItems:'center',gap:'6px',fontSize:'13px',fontWeight:600,color:toolStyle(item).lbl}">
            <span>{{toolIcon(item)}}</span>
            <span style="font-family:monospace">{{item.name||'Tool'}}</span>
          </div>
          <div style="display:flex;align-items:center;gap:6px">
            <button v-if="hasArgs(item)" @click="toggleArgs(idx)" style="font-size:11px;color:#64748b;background:#fff;border:1px solid #e2e8f0;cursor:pointer;padding:2px 7px;border-radius:4px">
              {{isArgExpanded(idx)?'▼':'▶'}} args
            </button>
            <button v-if="item.result" @click="toggleResult(idx)" style="font-size:11px;color:#64748b;background:#fff;border:1px solid #e2e8f0;cursor:pointer;padding:2px 7px;border-radius:4px">
              {{isResultExpanded(idx)?'▼':'▶'}} result
            </button>
          </div>
        </div>
        <div v-if="isArgExpanded(idx) && hasArgs(item)" style="margin-top:6px;padding:8px;background:#fff;border-radius:4px;font-size:11px;font-family:'SF Mono','Fira Code',monospace;color:#1e293b;max-height:200px;overflow-y:auto;border:1px solid #e2e8f0">
          <pre style="white-space:pre-wrap;word-break:break-all">{{JSON.stringify(item.arguments,null,2)}}</pre>
        </div>
        <div v-else-if="!hasArgs(item)&&!item.result" style="margin-top:4px;font-size:11px;color:#94a3b8">(no parameters)</div>
        <div v-if="isResultExpanded(idx) && item.result" style="margin-top:6px">
          <div v-for="(ri,ri_idx) in item.result" :key="ri_idx"
            :style="{padding:'8px',background:'#fff',borderRadius:'4px',fontSize:'11px',fontFamily:'SF Mono,Fira Code,monospace',maxHeight:'200px',overflowY:'auto',border:'1px solid '+toolStyle(item).border,color:item.isError?'#991b1b':'#1e293b',marginBottom:'4px'}">
            <div v-if="ri.type==='text'" style="white-space:pre-wrap;word-break:break-all">{{ri.text}}</div>
          </div>
          <div v-if="item.resultDetails&&item.resultDetails.diff" style="margin-top:4px;padding:8px;background:#0f172a;color:#e2e8f0;border-radius:4px;font-size:11px;font-family:'SF Mono',monospace;max-height:130px;overflow-y:auto">
            <pre style="white-space:pre-wrap">{{item.resultDetails.diff}}</pre>
          </div>
        </div>
      </div>
    </div>
    <div v-if="runtime && runtime.chips && runtime.chips.length" class="runtime-footer">
      <span class="runtime-label">Runtime</span>
      <span v-for="chip in runtime.chips" :key="chip.label" class="runtime-chip" :class="chip.tone">{{chip.label}}</span>
    </div>
  </div>
</div>
  ` + "`" + `
};

// ===== Main App =====
const { createApp, ref, computed, nextTick } = Vue;

createApp({
  components: { MessageCard },
  setup() {
    const connected = ref(false);
    const reconnecting = ref(false);
    const showConnectPanel = ref(false);
    const tokenInput = ref((function(){ try { return localStorage.getItem('clawwatch_token') || ''; } catch { return ''; } })());
    const subscribeInput = ref('*');
    const searchQuery = ref('');
    const selectedAgent = ref(null);
    const selectedSession = ref(null);
    const trajectoryByHost = ref({});
    const unreadOnly = ref(false);
    const expandedHosts = ref({});
    const expandedAgents = ref({});
    const agents = ref({});
    const knownSessions = ref({});
    const messagesByAgent = ref({});
    const allMessages = ref([]);
    const unread = ref({});
    const forceExpandThinking = ref(0);
    const forceExpandArgs = ref(0);
    const allThinkingExpanded = ref(false);
    const allArgsExpanded = ref(false);
    const messageArea = ref(null);
    const bottomAnchor = ref(null);
    let ws = null;
    let reconnectTimer = null;
    let reconnectAttempts = 0;
    let msgIdCounter = 0;
    const recentActivityWindowMs = 4 * 60 * 60 * 1000;
    const maxSessionsPerHost = 15;

    const agentCount = computed(() => Object.keys(agents.value).length);
    const connectionStatusText = computed(() => connected.value ? 'Connected' : (reconnecting.value ? 'Reconnecting...' : 'Disconnected'));

    const hostGroups = computed(() => {
      const groups = new Map();
      Object.entries(agents.value).forEach(([id, info]) => {
        if (!info.openclawAgentId) {
          groups.set(id, {
            id,
            label: info.hostname || info.label || shortId(id),
            subtitle: info.subtitle || '',
            online: !!info.online,
            children: [],
            totalCount: msgCountFor(id)
          });
        }
      });
      Object.entries(agents.value).forEach(([id, info]) => {
        if (!info.openclawAgentId) return;
        const suffix = '/' + info.openclawAgentId;
        const hostId = id.endsWith(suffix) ? id.slice(0, -suffix.length) : id;
        if (!groups.has(hostId)) {
          groups.set(hostId, {
            id: hostId,
            label: info.hostname || shortId(hostId),
            subtitle: info.subtitle || '',
            online: !!info.online,
            children: [],
            totalCount: 0
          });
        }
        const group = groups.get(hostId);
        const count = msgCountFor(id);
        const sessions = sessionGroups.value
          .filter(session => session.agentId === id)
          .sort((a, b) => b.lastTime - a.lastTime);
        group.children.push({ id, name: info.openclawAgentId, count, sessions });
        group.totalCount += count;
        group.online = group.online || !!info.online;
      });
      groups.forEach(group => {
        if (unreadOnly.value) {
          group.children = group.children
            .map(child => ({
              ...child,
              sessions: child.sessions.filter(session => unreadFor('session', child.id, session.id) > 0)
            }))
            .filter(child => child.sessions.length > 0);
        }
        const allSessions = group.children.flatMap(child => child.sessions);
        const sortedSessions = [...allSessions].sort(compareSessions);
        let visible = sortedSessions.slice(0, maxSessionsPerHost);
        const selected = allSessions.find(session => selectedAgent.value === session.agentId && selectedSession.value === session.id);
        if (selected && !visible.some(session => session.key === selected.key)) {
          visible = [...visible.slice(0, maxSessionsPerHost - 1), selected];
        }
        const visibleKeys = new Set(visible.map(session => session.key));
        group.children = group.children
          .map(child => ({
            ...child,
            sessions: child.sessions.filter(session => visibleKeys.has(session.key)).sort(compareSessions)
          }))
          .sort(compareAgents);
        group.hiddenSessionCount = Math.max(0, allSessions.length - visible.length);
      });
      return [...groups.values()]
        .filter(group => !unreadOnly.value || group.children.length > 0)
        .sort((a, b) => a.label.localeCompare(b.label));
    });
    const hostCount = computed(() => hostGroups.value.length);
    const openClawAgentCount = computed(() => hostGroups.value.reduce((sum, host) => sum + host.children.length, 0));
    const sessionCount = computed(() => sessionGroups.value.length);

    function msgCountFor(agentId) {
      return (messagesByAgent.value[agentId] || []).length;
    }

    function sessionRank(session) {
      return Date.now() - session.lastTime <= recentActivityWindowMs ? 2 : 1;
    }

    function compareSessions(a, b) {
      return sessionRank(b) - sessionRank(a) || b.lastTime - a.lastTime;
    }

    function isRecentSession(session) {
      return !!session.lastTime && Date.now() - session.lastTime <= recentActivityWindowMs;
    }

    function compareAgents(a, b) {
      const aLast = a.sessions[0]?.lastTime || 0;
      const bLast = b.sessions[0]?.lastTime || 0;
      return bLast - aLast || a.name.localeCompare(b.name);
    }

    function unreadKey(scope, agentId = '', sessionId = '') {
      return scope + '|' + agentId + '|' + sessionId;
    }

    function unreadFor(scope, agentId = '', sessionId = '') {
      if (scope === 'session') return unread.value[unreadKey('session', agentId, sessionId)] || 0;
      return Object.entries(unread.value).reduce((sum, [key, count]) => {
        if (!key.startsWith('session|')) return sum;
        const source = key.slice('session|'.length).split('|')[0];
        if (scope === 'all') return sum + count;
        if (scope === 'agent' && source === agentId) return sum + count;
        if (scope === 'host' && hostIdFor(source) === agentId) return sum + count;
        return sum;
      }, 0);
    }

    function hostIdFor(agentId) {
      const info = agents.value[agentId];
      if (info?.openclawAgentId) {
        const suffix = '/' + info.openclawAgentId;
        if (agentId.endsWith(suffix)) return agentId.slice(0, -suffix.length);
      }
      return agentId;
    }

    function isViewing(agentId, sessionId) {
      return !!selectedSession.value && selectedAgent.value === agentId && selectedSession.value === sessionId;
    }

    function incrementUnread(agentId, sessionId) {
      if (isViewing(agentId, sessionId)) return;
      const next = { ...unread.value };
      const key = unreadKey('session', agentId, sessionId);
      next[key] = (next[key] || 0) + 1;
      unread.value = next;
    }

    function clearUnread(scope, agentId = '', sessionId = '') {
      const next = { ...unread.value };
      if (scope === 'all') {
        unread.value = {};
        return;
      }
      Object.keys(next).forEach(key => {
        const source = key.startsWith('session|') ? key.slice('session|'.length).split('|')[0] : '';
        if (scope === 'host' && hostIdFor(source) === agentId) {
          delete next[key];
        } else if (scope === 'agent' && source === agentId) {
          delete next[key];
        } else if (scope === 'session' && key === unreadKey('session', agentId, sessionId)) {
          delete next[key];
        }
      });
      unread.value = next;
    }

    const sessionGroups = computed(() => {
      const sessions = new Map();
      Object.values(knownSessions.value).forEach(session => {
        sessions.set(session.key, { ...session, count: 0, lastTime: session.lastTime || 0 });
      });
      allMessages.value.forEach(msg => {
        if (!msg.sessionId) return;
        const key = msg.agentId + '|' + msg.sessionId;
        const current = sessions.get(key) || { key, id: msg.sessionId, agentId: msg.agentId, count: 0, lastTime: 0 };
        current.count++;
        current.lastTime = Math.max(current.lastTime, new Date(msg.timestamp || 0).getTime());
        sessions.set(key, current);
      });
      return [...sessions.values()];
    });

    function relativeTime(ts) {
      if (!ts) return '';
      const diff = Date.now() - ts;
      if (diff < 60000) return '刚刚';
      if (diff < 3600000) return Math.floor(diff / 60000) + ' 分钟前';
      if (diff < 86400000) return Math.floor(diff / 3600000) + ' 小时前';
      return Math.floor(diff / 86400000) + ' 天前';
    }
    const displayMessages = computed(() => {
      if (!selectedAgent.value || !selectedSession.value) return [];
      const messages = messagesByAgent.value[selectedAgent.value] || [];
      return messages.filter(msg => msg.sessionId === selectedSession.value);
    });

    const selectedHostId = computed(() => selectedAgent.value ? hostIdFor(selectedAgent.value) : '');
    const selectedOpenClawAgentId = computed(() => agents.value[selectedAgent.value]?.openclawAgentId || '');
    const selectedTrajectory = computed(() => {
      if (!selectedSession.value || !selectedHostId.value) return [];
      return (trajectoryByHost.value[selectedHostId.value] || [])
        .filter(metric => metric.openclawAgentId === selectedOpenClawAgentId.value && metric.sessionId === selectedSession.value)
        .sort((a, b) => new Date(a.eventTimestamp || a.timestamp || 0) - new Date(b.eventTimestamp || b.timestamp || 0));
    });

    const selectedViewTitle = computed(() => {
      if (!selectedSession.value) return 'Select a session';
      const info = agents.value[selectedAgent.value];
      const agentTitle = info?.openclawAgentId || info?.label || shortId(selectedAgent.value);
      return agentTitle + ' / session ' + shortSession(selectedSession.value);
    });

    const selectedViewSubtitle = computed(() => {
      if (!selectedSession.value) return 'Host and agent rows only expand or collapse navigation';
      return filteredMessages.value.length + ' visible messages · ' + runtimeCoverage.value.integrated + '/' + runtimeCoverage.value.total + ' runtime events integrated';
    });

    const recoveredAttempts = computed(() => {
      const existing = new Set(displayMessages.value.map(message => messageSignature(message.message)));
      const recovered = [];
      selectedTrajectory.value.forEach((event, eventIndex) => {
        [event.data?.messagesSnapshot, event.data?.messages].forEach(source => {
          if (!Array.isArray(source)) return;
          source.forEach((message, messageIndex) => {
            if (!isFailedAssistantAttempt(message)) return;
            const signature = messageSignature(message);
            if (existing.has(signature)) return;
            existing.add(signature);
            recovered.push({
              _id: 'recovered|' + eventIndex + '|' + messageIndex + '|' + signature,
              agentId: selectedAgent.value,
              hostId: selectedHostId.value,
              hostname: agents.value[selectedAgent.value]?.hostname,
              hostIPs: agents.value[selectedAgent.value]?.hostIPs || [],
              openclawAgentId: selectedOpenClawAgentId.value,
              sessionId: selectedSession.value,
              timestamp: message.timestamp || event.eventTimestamp || event.timestamp,
              message: {
                ...message,
                provider: message.provider || event.provider,
                model: message.model || message.modelId || event.modelId,
                api: message.api || event.modelApi
              },
              _recoveredFromTrajectory: true
            });
          });
        });
      });
      return recovered;
    });

    const conversationMessages = computed(() => {
      return [...displayMessages.value, ...recoveredAttempts.value]
        .sort((a, b) => new Date(a.timestamp || 0).getTime() - new Date(b.timestamp || 0).getTime());
    });

    const runtimeAssignments = computed(() => {
      const assistants = conversationMessages.value
        .filter(message => message.message?.role === 'assistant')
        .map(message => ({ id: message._id, time: new Date(message.timestamp || 0).getTime(), events: [] }));
      const integratedSessionEvents = new Set(['session.started', 'trace.metadata']);
      let sessionIntegrated = 0;
      const unassociated = [];
      selectedTrajectory.value.forEach(event => {
        if (integratedSessionEvents.has(event.eventType)) {
          sessionIntegrated++;
          return;
        }
        const eventTime = new Date(event.eventTimestamp || event.timestamp || 0).getTime();
        let best = null;
        assistants.forEach(assistant => {
          const distance = Math.abs(assistant.time - eventTime);
          if (!best || distance < best.distance) best = { assistant, distance };
        });
        if (best && best.distance <= 10 * 60 * 1000) best.assistant.events.push(event);
        else unassociated.push(event);
      });
      return { assistants, sessionIntegrated, unassociated };
    });

    const runtimeCoverage = computed(() => {
      const assigned = runtimeAssignments.value.assistants.reduce((sum, assistant) => sum + assistant.events.length, 0);
      const integrated = assigned + runtimeAssignments.value.sessionIntegrated;
      return { total: selectedTrajectory.value.length, integrated, unassociated: runtimeAssignments.value.unassociated.length };
    });

    const sessionRuntime = computed(() => {
      const started = selectedTrajectory.value.find(event => event.eventType === 'session.started');
      const metadata = selectedTrajectory.value.find(event => event.eventType === 'trace.metadata');
      const models = [...new Set(selectedTrajectory.value.map(event => event.modelId).filter(Boolean))];
      const providers = [...new Set(selectedTrajectory.value.map(event => event.provider).filter(Boolean))];
      const chips = [];
      if (started?.data?.trigger) chips.push(runtimeChip('trigger: ' + started.data.trigger));
      if (providers.length) chips.push(runtimeChip(providers.join(' → ')));
      if (models.length) chips.push(runtimeChip(models.join(' → ')));
      if (started?.modelApi) chips.push(runtimeChip(started.modelApi));
      if (started?.data?.toolCount !== undefined) chips.push(runtimeChip(started.data.toolCount + ' tools available'));
      const skillCount = Array.isArray(metadata?.data?.skills) ? metadata.data.skills.length : null;
      const pluginCount = Array.isArray(metadata?.data?.plugins) ? metadata.data.plugins.length : null;
      if (skillCount !== null) chips.push(runtimeChip(skillCount + ' skills'));
      if (pluginCount !== null) chips.push(runtimeChip(pluginCount + ' plugins'));
      return { chips };
    });

    const filteredMessages = computed(() => {
      // Build toolResult map
      const toolResultMap = new Map();
      conversationMessages.value.forEach(msg => {
        if (msg.message?.role === 'toolResult') {
          const cid = msg.message.toolCallId;
          if (cid) toolResultMap.set(cid, msg.message);
        }
      });

      const merged = [];
      conversationMessages.value.forEach(msg => {
        if (msg.message?.role === 'toolResult') return;
        const cloned = JSON.parse(JSON.stringify(msg));
        if (cloned.message?.role === 'assistant' && Array.isArray(cloned.message.content)) {
          cloned.message.content = cloned.message.content.map(item => {
            if (item.type === 'toolCall' && item.id) {
              const res = toolResultMap.get(item.id);
              if (res) {
                let isError = !!res.isError;
                if (!isError && res.details?.status === 'error') isError = true;
                if (!isError && JSON.stringify(res.content || []).includes('"status":"error"')) isError = true;
                return { ...item, result: res.content, resultDetails: res.details, isError };
              }
            }
            return item;
          });
        }
        const assignment = runtimeAssignments.value.assistants.find(assistant => assistant.id === cloned._id);
        if (assignment?.events.length) cloned._runtime = summarizeRuntimeEvents(assignment.events);
        merged.push(cloned);
      });

      if (!searchQuery.value) return merged;
      const q = searchQuery.value.toLowerCase();
      return merged.filter(m => JSON.stringify(m.message?.content || []).toLowerCase().includes(q));
    });

    const seenMessages = new Set();

    function shortId(id) {
      if (!id) return '';
      const parts = id.split(/[/:]/);
      return parts[parts.length - 1] || id;
    }

    function isFailedAssistantAttempt(message) {
      if (message?.role !== 'assistant') return false;
      if (message.errorMessage) return true;
      return Array.isArray(message.content) && message.content.some(item => item?.type === 'text' && item.text === '[assistant turn failed before producing content]');
    }

    function messageSignature(message) {
      if (!message) return '';
      return [message.role || '', message.timestamp || '', message.errorMessage || '', JSON.stringify(message.content || '')].join('|');
    }

    function shortSession(id) {
      if (!id) return '';
      return id.length > 12 ? id.slice(0, 8) + '…' + id.slice(-4) : id;
    }

    function isHostExpanded(hostId) {
      return unreadOnly.value || !!expandedHosts.value[hostId];
    }

    function isAgentExpanded(agentId) {
      return unreadOnly.value || !!expandedAgents.value[agentId];
    }

    function toggleUnreadOnly() {
      if (!unreadOnly.value && unreadFor('all') === 0) return;
      unreadOnly.value = !unreadOnly.value;
    }

    function toggleHost(hostId) {
      expandedHosts.value = { ...expandedHosts.value, [hostId]: !isHostExpanded(hostId) };
    }

    function toggleAgent(agentId) {
      expandedAgents.value = { ...expandedAgents.value, [agentId]: !isAgentExpanded(agentId) };
    }

    function selectSession(agentId, sessionId) {
      selectedAgent.value = agentId;
      selectedSession.value = sessionId;
      expandedHosts.value = { ...expandedHosts.value, [hostIdFor(agentId)]: true };
      expandedAgents.value = { ...expandedAgents.value, [agentId]: true };
      clearUnread('session', agentId, sessionId);
      newMessageCount.value = 0;
      fetchTrajectory(hostIdFor(agentId));
      nextTick(() => scrollToBottom(true));
    }

    async function fetchSessionHistory(agentId, sessionId) {
      const hostId = hostIdFor(agentId);
      const openclawAgentId = agents.value[agentId]?.openclawAgentId || '';
      if (!hostId || !openclawAgentId || !sessionId) return;
      try {
        const res = await fetch('/api/session-history?agentId=' + encodeURIComponent(hostId) + '&openclawAgentId=' + encodeURIComponent(openclawAgentId) + '&sessionId=' + encodeURIComponent(sessionId) + apiTokenSuffix());
        if (!res.ok) return;
        const text = await res.text();
        if (!text) return;
        const data = JSON.parse(text);
        if (data.lines && Array.isArray(data.lines)) {
          data.lines.forEach(line => {
            const raw = typeof line === 'string' ? line : JSON.stringify(line);
            handleMsg(JSON.stringify({type:'log',agentId:hostId,openclawAgentId:openclawAgentId,sessionId:sessionId,lines:[JSON.parse(raw)],timestamp:Date.now()}), true);
          });
          nextTick(() => scrollToBottom(true));
        }
      } catch(e) { /* silent - session history is best-effort */ }
    }

    function runtimeChip(label, tone = '') {
      return { label, tone };
    }

    function compactNumber(value) {
      if (typeof value !== 'number') return value;
      if (value >= 1000000) return (value / 1000000).toFixed(1) + 'm';
      if (value >= 1000) return (value / 1000).toFixed(1) + 'k';
      return String(value);
    }

    function summarizeRuntimeEvents(events) {
      const ordered = [...events].sort((a, b) => new Date(a.eventTimestamp || a.timestamp || 0) - new Date(b.eventTimestamp || b.timestamp || 0));
      const completed = [...ordered].reverse().find(event => event.eventType === 'trace.artifacts' || event.eventType === 'model.completed');
      const ended = [...ordered].reverse().find(event => event.eventType === 'session.ended');
      const context = [...ordered].reverse().find(event => event.eventType === 'context.compiled');
      const promptEvent = [...ordered].reverse().find(event => event.eventType === 'prompt.submitted');
      const firstTime = new Date(ordered[0]?.eventTimestamp || ordered[0]?.timestamp || 0).getTime();
      const lastTime = new Date(ordered[ordered.length - 1]?.eventTimestamp || ordered[ordered.length - 1]?.timestamp || 0).getTime();
      const chips = [];
      const status = ended?.data?.status || completed?.data?.finalStatus;
      if (status) chips.push(runtimeChip(status, status === 'success' ? 'ok' : 'error'));
      if (lastTime > firstTime) chips.push(runtimeChip(((lastTime - firstTime) / 1000).toFixed(1) + 's'));
      const usage = completed?.data?.usage;
      if (usage?.total !== undefined) chips.push(runtimeChip(compactNumber(usage.total) + ' tokens'));
      if (usage?.input !== undefined || usage?.output !== undefined) chips.push(runtimeChip((usage?.input || 0) + ' in / ' + (usage?.output || 0) + ' out'));
      if (usage?.cacheRead) chips.push(runtimeChip(compactNumber(usage.cacheRead) + ' cache read'));
      if (usage?.cacheWrite) chips.push(runtimeChip(compactNumber(usage.cacheWrite) + ' cache write'));
      if (context?.data?.messages && Array.isArray(context.data.messages)) chips.push(runtimeChip(context.data.messages.length + ' context messages'));
      if (context?.data?.tools && Array.isArray(context.data.tools)) chips.push(runtimeChip(context.data.tools.length + ' tools'));
      if (context?.data?.imagesCount) chips.push(runtimeChip(context.data.imagesCount + ' images'));
      if (completed?.data?.compactionCount) chips.push(runtimeChip(completed.data.compactionCount + ' compactions', 'warn'));
      if (completed?.data?.timedOut || ended?.data?.timedOut) chips.push(runtimeChip('timed out', 'error'));
      if (completed?.data?.aborted || ended?.data?.aborted) chips.push(runtimeChip('aborted', 'error'));
      if (completed?.data?.externalAbort || ended?.data?.externalAbort) chips.push(runtimeChip('external abort', 'error'));
      const promptData = promptEvent?.data || context?.data;
      const promptLabels = [];
      if (typeof promptData?.systemPrompt === 'string' && promptData.systemPrompt) promptLabels.push('system prompt');
      if (typeof promptData?.prompt === 'string' && promptData.prompt) promptLabels.push('prompt');
      if (Array.isArray(promptData?.messages)) promptLabels.push(promptData.messages.length + ' messages');
      return {
        model: completed?.modelId || ordered[ordered.length - 1]?.modelId,
        provider: completed?.provider,
        api: completed?.modelApi,
        chips: [
          ...(completed?.provider ? [runtimeChip(completed.provider)] : []),
          ...(completed?.modelApi ? [runtimeChip(completed.modelApi)] : []),
          ...chips
        ],
        promptContext: promptLabels.length ? {
          labels: promptLabels,
          systemPrompt: promptData.systemPrompt,
          prompt: promptData.prompt,
          messages: promptData.messages
        } : null,
        eventCount: events.length
      };
    }

    function trajectoryKey(metric, index) {
      return [metric.agentId, metric.openclawAgentId, metric.sessionId, metric.eventType, metric.eventTimestamp || metric.timestamp, index].join('|');
    }

    function addMetric(metric) {
      const hostId = metric.agentId || metric.host?.id || 'unknown';
      const current = trajectoryByHost.value[hostId] || [];
      const key = trajectoryKey(metric, 0);
      if (current.some(item => trajectoryKey(item, 0) === key)) return;
      trajectoryByHost.value = { ...trajectoryByHost.value, [hostId]: [...current, metric].slice(-2000) };
    }

    function ensureAgent(agentId, online, meta = {}) {
      if (!agents.value[agentId]) {
        agents.value = { ...agents.value, [agentId]: { online: !!online, ...meta } };
      } else {
        agents.value[agentId] = { ...agents.value[agentId], online: !!online, ...meta };
      }
      if (!messagesByAgent.value[agentId]) {
        messagesByAgent.value = { ...messagesByAgent.value, [agentId]: [] };
      }
    }

    function sourceKey(msg) {
      const hostId = msg.agentId || msg.host?.id || 'unknown';
      return msg.openclawAgentId ? hostId + '/' + msg.openclawAgentId : hostId;
    }

    function sourceMeta(msg) {
      const hostname = msg.hostname || msg.host?.hostname || msg.agentId || 'unknown';
      const ips = msg.hostIPs || msg.host?.ips || [];
      return {
        label: msg.openclawAgentId ? hostname + ' / ' + msg.openclawAgentId : hostname,
        subtitle: ips.join(', '),
        hostname,
        hostIPs: ips,
        openclawAgentId: msg.openclawAgentId || ''
      };
    }

    function addLines(agentId, lines, meta, trackUnread = true) {
      ensureAgent(agentId, true, sourceMeta(meta || {}));
      let added = false;
      lines.forEach(line => {
        let parsed;
        try { parsed = typeof line === 'string' ? JSON.parse(line) : line; } catch { return; }
        if (!parsed || !parsed.message) return;
        const role = parsed.message?.role;
        if (!role || !['user','assistant','tool','toolResult'].includes(role)) return;
        const dedupeKey = agentId + '|' + (meta?.sessionId || meta?.session || '') + '|' +
          (parsed.id || '') + '|' + (parsed.timestamp || '') + '|' + role + '|' + JSON.stringify(parsed.message);
        if (seenMessages.has(dedupeKey)) return;
        seenMessages.add(dedupeKey);
        const msgObj = {
          _id: ++msgIdCounter,
          agentId,
          hostId: meta?.agentId || meta?.host?.id,
          hostname: meta?.hostname || meta?.host?.hostname,
          hostIPs: meta?.hostIPs || meta?.host?.ips || [],
          openclawAgentId: meta?.openclawAgentId,
          sessionId: meta?.sessionId || meta?.session,
          timestamp: parsed.timestamp || meta?.timestamp || new Date().toISOString(),
          message: parsed.message
        };
        messagesByAgent.value[agentId].push(msgObj);
        allMessages.value.push(msgObj);
        if (trackUnread) incrementUnread(agentId, msgObj.sessionId || '');
        added = true;
      });
      if (added) onNewMessages(!trackUnread);
    }

    const isAtBottom = ref(true);
    const newMessageCount = ref(0);

    function onMessageScroll() {
      const el = messageArea.value;
      if (!el) return;
      const threshold = 50;
      isAtBottom.value = el.scrollHeight - el.scrollTop - el.clientHeight < threshold;
      if (isAtBottom.value) newMessageCount.value = 0;
    }

    function onNewMessages(isHistory) {
      if (isHistory) {
        nextTick(() => scrollToBottom(true));
        return;
      }
      if (isAtBottom.value) {
        nextTick(() => scrollToBottom(false));
      } else {
        newMessageCount.value++;
      }
    }

    function scrollToBottom(force) {
      nextTick(() => {
        if (bottomAnchor.value) bottomAnchor.value.scrollIntoView({ behavior: force ? 'auto' : 'smooth' });
        newMessageCount.value = 0;
        isAtBottom.value = true;
      });
    }

    function handleMsg(data, isHistory = false) {
      let msg;
      try { msg = JSON.parse(data); } catch { return; }
      const agentId = msg.agentId || msg.agent || 'unknown';

      if (msg.type === 'agent_status') {
        const hostId = msg.id || agentId;
        ensureAgent(hostId, msg.online, sourceMeta(msg));
        if (Array.isArray(msg.agents)) msg.agents.forEach(a => {
          ensureAgent(hostId + '/' + a.id, msg.online, sourceMeta({...msg, openclawAgentId:a.id}));
        });
        return;
      }
      if (msg.type === 'agent_hello') {
        ensureAgent(agentId, true, sourceMeta(msg));
        if (Array.isArray(msg.agents)) msg.agents.forEach(a => {
          ensureAgent(agentId + '/' + a.id, true, sourceMeta({...msg, openclawAgentId:a.id}));
        });
        return;
      }
      if (msg.type === 'log' && Array.isArray(msg.lines)) {
        addLines(sourceKey(msg), msg.lines, msg, !isHistory);
        return;
      }
      if (msg.type === 'trajectory') {
        addMetric(msg);
        return;
      }
      if (msg.type === 'session_list') {
        ensureAgent(agentId, true, sourceMeta(msg));
        if (Array.isArray(msg.sessions)) {
          msg.sessions.forEach(session => {
            const openclawAgentId = session.agentId || session.agent;
            if (!openclawAgentId || !session.sessionId) return;
            const source = agentId + '/' + openclawAgentId;
            ensureAgent(source, true, sourceMeta({...msg, openclawAgentId}));
            const key = source + '|' + session.sessionId;
            knownSessions.value = {
              ...knownSessions.value,
              [key]: {
                key,
                id: session.sessionId,
                agentId: source,
                lastTime: session.mtime || 0
              }
            };
          });
        }
        return;
      }
    }

    function scheduleReconnect() {
      if (reconnectTimer) return;
      reconnecting.value = true;
      const delay = Math.min(30000, 1000 * Math.pow(2, reconnectAttempts));
      reconnectAttempts++;
      reconnectTimer = setTimeout(() => {
        reconnectTimer = null;
        connect(false);
      }, delay);
    }

    function connect(resetBackoff = true) {
      if (ws) { try { ws.close(); } catch {} ws = null; }
      if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null; }
      if (resetBackoff) reconnectAttempts = 0;
      connected.value = false;
      reconnecting.value = !resetBackoff;
      const token = tokenInput.value.trim();
      try { localStorage.setItem('clawwatch_token', token); } catch {}
      const sub = subscribeInput.value.trim() || '*';
      const proto = location.protocol === 'https:' ? 'wss' : 'ws';
      let url = proto + '://' + location.host + '/ws/console?subscribe=' + encodeURIComponent(sub);
      if (token) url += '&token=' + encodeURIComponent(token);

      ws = new WebSocket(url);
      const socket = ws;
      ws.onopen = () => {
        if (ws !== socket) return;
        connected.value = true;
        reconnecting.value = false;
        reconnectAttempts = 0;
        showConnectPanel.value = false;
        fetchAgents();
        fetchHistory(sub);
      };
      ws.onmessage = e => handleMsg(e.data);
      ws.onerror = () => {
        if (ws === socket) connected.value = false;
      };
      ws.onclose = () => {
        if (ws !== socket) return;
        connected.value = false;
        scheduleReconnect();
      };
    }
    function reconnect() {
      connect(true);
    }

    function clearMessages() {
      allMessages.value = [];
      messagesByAgent.value = {};
      unread.value = {};
      seenMessages.clear();
      Object.keys(agents.value).forEach(id => { messagesByAgent.value[id] = []; });
    }

    // apiTokenSuffix 返回当前 token 的查询串后缀（带 & 前缀），用于 REST API 鉴权
    function apiTokenSuffix() {
      const token = tokenInput.value.trim();
      return token ? '&token=' + encodeURIComponent(token) : '';
    }

    async function fetchHistory(agentId) {
      try {
        const res = await fetch('/api/history?agentId=' + encodeURIComponent(agentId || '*') + '&limit=1000' + apiTokenSuffix());
        const data = await res.json();
        if (Array.isArray(data.messages)) data.messages.forEach(msg => handleMsg(JSON.stringify(msg), true));
      } catch {}
    }

    async function fetchTrajectory(hostId) {
      if (!hostId) return;
      try {
        const res = await fetch('/api/trajectory?agentId=' + encodeURIComponent(hostId) + '&limit=2000' + apiTokenSuffix());
        const data = await res.json();
        if (Array.isArray(data.messages)) {
          trajectoryByHost.value = { ...trajectoryByHost.value, [hostId]: data.messages };
        }
      } catch {}
    }

    async function fetchAgents() {
      try {
        const res = await fetch('/api/agents' + (apiTokenSuffix() ? '?' + apiTokenSuffix().slice(1) : ''));
        const data = await res.json();
        if (Array.isArray(data.agents)) data.agents.forEach(a => {
          ensureAgent(a.id, true, sourceMeta({...a, agentId:a.id, hostIPs:a.hostIPs || []}));
          if (Array.isArray(a.openclawAgents)) a.openclawAgents.forEach(child => {
            ensureAgent(a.id + '/' + child.id, true, sourceMeta({...a, agentId:a.id, hostIPs:a.hostIPs || [], openclawAgentId:child.id}));
          });
          if (Array.isArray(a.sessions)) a.sessions.forEach(session => {
            const openclawAgentId = session.agentId || session.agent;
            if (!openclawAgentId || !session.sessionId) return;
            const source = a.id + '/' + openclawAgentId;
            ensureAgent(source, true, sourceMeta({...a, agentId:a.id, hostIPs:a.hostIPs || [], openclawAgentId}));
            const key = source + '|' + session.sessionId;
            knownSessions.value = {
              ...knownSessions.value,
              [key]: { key, id: session.sessionId, agentId: source, lastTime: session.mtime || 0 }
            };
          });
        });
      } catch {}
    }

    function toggleAllThinking() {
      if (allThinkingExpanded.value) { forceExpandThinking.value--; allThinkingExpanded.value = false; }
      else { forceExpandThinking.value++; allThinkingExpanded.value = true; }
    }

    function toggleAllArgs() {
      if (allArgsExpanded.value) { forceExpandArgs.value--; allArgsExpanded.value = false; }
      else { forceExpandArgs.value++; allArgsExpanded.value = true; }
    }

    // Auto-connect on load
    connect();

    return {
      connected, reconnecting, connectionStatusText, showConnectPanel, tokenInput, subscribeInput, searchQuery,
      selectedAgent, selectedSession, unreadOnly, agents, allMessages, agentCount, hostCount, openClawAgentCount,
      sessionCount, hostGroups, selectedViewTitle, selectedViewSubtitle, forceExpandThinking, forceExpandArgs,
      allThinkingExpanded, allArgsExpanded, messageArea, bottomAnchor,
      isAtBottom, newMessageCount, onMessageScroll, scrollToBottom,
      filteredMessages, recoveredAttempts, selectedTrajectory, sessionRuntime, runtimeCoverage, msgCountFor, unreadFor, shortId, shortSession, isRecentSession,
      isHostExpanded, isAgentExpanded, toggleUnreadOnly, toggleHost, toggleAgent, selectSession, reconnect, clearMessages,
      toggleAllThinking, toggleAllArgs,
      relativeTime
    };
  }
}).mount('#app');
</script>
</body>
</html>`
	return part1 + string(vue3JS) + part2 + string(markedJS) + part3
}
