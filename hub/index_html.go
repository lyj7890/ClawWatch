package main

func buildIndexHTML(cfg *Config) string {
	part1 := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ClawWatch Hub</title>
<script>`
	part2 := `<\/script>
<script>`
	part3 := `<\/script>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html, body { height: 100%; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8fafc; color: #1e293b; }

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

.agent-item:hover { background: #f8fafc; }

btn { display:inline-flex;align-items:center;gap:5px;padding:5px 12px;border-radius:6px;font-size:12px;font-weight:500;cursor:pointer;border:1px solid transparent;transition:background .15s; }
</style>
</head>
<body>
<div id="app" style="height:100vh;display:flex;flex-direction:column;overflow:hidden">

  <!-- Header -->
  <header style="background:#fff;border-bottom:1px solid #e2e8f0;padding:10px 20px;display:flex;align-items:center;justify-content:space-between;flex-shrink:0;box-shadow:0 1px 3px rgba(0,0,0,.06)">
    <div style="display:flex;align-items:center;gap:12px">
      <span style="font-size:22px">🦞</span>
      <h1 style="font-size:17px;font-weight:700;color:#1e293b">ClawWatch Hub</h1>
      <div style="display:flex;align-items:center;gap:6px;margin-left:8px">
        <div :style="connected ? 'background:#22c55e' : 'background:#ef4444'" style="width:8px;height:8px;border-radius:50%;transition:background .3s"></div>
        <span style="font-size:12px;color:#64748b">{{ connected ? 'Connected' : 'Disconnected' }}</span>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:8px">
      <button @click="toggleAllThinking" style="padding:5px 12px;background:#eff6ff;border:1px solid #bfdbfe;color:#2563eb;border-radius:6px;font-size:12px;cursor:pointer;font-weight:500">
        💭 {{ allThinkingExpanded ? 'Collapse' : 'Expand' }} Thinking
      </button>
      <button @click="toggleAllArgs" style="padding:5px 12px;background:#fffbeb;border:1px solid #fde68a;color:#92400e;border-radius:6px;font-size:12px;cursor:pointer;font-weight:500">
        🔧 {{ allArgsExpanded ? 'Hide' : 'Show' }} Args
      </button>
      <div style="width:1px;height:20px;background:#e2e8f0;margin:0 4px"></div>
      <button @click="showConnectPanel = !showConnectPanel" style="padding:5px 12px;background:#f1f5f9;border:1px solid #e2e8f0;color:#475569;border-radius:6px;font-size:12px;cursor:pointer;font-weight:500">
        ⚙️ Settings
      </button>
    </div>
  </header>

  <!-- Connect Panel -->
  <div v-show="showConnectPanel" style="background:#f8fafc;border-bottom:1px solid #e2e8f0;padding:10px 20px;display:flex;align-items:center;gap:10px;flex-shrink:0">
    <label style="font-size:12px;color:#64748b;white-space:nowrap">Token:</label>
    <input v-model="tokenInput" type="password" placeholder="Console Token (空=跳过)" style="width:200px;background:#fff;border:1px solid #e2e8f0;border-radius:6px;padding:5px 10px;font-size:13px;color:#1e293b;outline:none">
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
    <aside style="width:220px;background:#fff;border-right:1px solid #e2e8f0;display:flex;flex-direction:column;flex-shrink:0;overflow:hidden">

      <!-- 最近访问 -->
      <div v-if="recentAgents.length > 0" style="flex-shrink:0;border-bottom:1px solid #e2e8f0">
        <div style="padding:8px 14px 4px;font-size:11px;color:#94a3b8;text-transform:uppercase;letter-spacing:.8px;font-weight:600;display:flex;align-items:center;gap:5px">
          <span>🕐</span><span>最近访问</span>
        </div>
        <div v-for="a in recentAgents" :key="'r-'+a.id" @click="selectAgent(a.id)" class="agent-item"
          :style="selectedAgent === a.id ? 'background:#eff6ff;border-left:3px solid #3b82f6;' : 'border-left:3px solid transparent;'"
          style="padding:7px 12px;cursor:pointer;transition:background .15s">
          <div style="display:flex;align-items:center;gap:7px">
            <div :style="{width:'8px',height:'8px',borderRadius:'50%',background:agentColor(a.id),flexShrink:0}"></div>
            <span :style="selectedAgent === a.id ? 'color:#1d4ed8;font-weight:600' : 'color:#374151'"
              style="font-size:12px;font-family:monospace;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:140px" :title="a.id">{{ shortId(a.id) }}</span>
            <span style="margin-left:auto;font-size:11px;color:#94a3b8;background:#f1f5f9;padding:1px 5px;border-radius:8px;flex-shrink:0">{{ a.count }}</span>
          </div>
          <div style="font-size:11px;color:#94a3b8;margin-top:2px;padding-left:15px">{{ relativeTime(a.lastTime) }}</div>
        </div>
      </div>

      <div style="padding:10px 14px;font-size:11px;color:#94a3b8;text-transform:uppercase;letter-spacing:.8px;font-weight:600;border-bottom:1px solid #f1f5f9">
        Agents ({{ agentCount }})
      </div>
      <div style="flex:1;overflow-y:auto">
        <!-- All -->
        <div @click="selectAgent(null)" class="agent-item"
          :style="selectedAgent === null ? 'background:#eff6ff;border-left:3px solid #3b82f6;' : 'border-left:3px solid transparent;'"
          style="padding:9px 14px;cursor:pointer;display:flex;align-items:center;gap:8px;font-size:13px;border-bottom:1px solid #f8fafc;transition:background .15s">
          <span style="font-size:15px">🌐</span>
          <span :style="selectedAgent === null ? 'color:#1d4ed8;font-weight:600' : 'color:#374151'">All Agents</span>
          <span style="margin-left:auto;font-size:11px;color:#94a3b8;background:#f1f5f9;padding:1px 5px;border-radius:8px">{{ allMessages.length }}</span>
        </div>
        <!-- Per-agent -->
        <div v-for="(info, agentId) in agents" :key="agentId" @click="selectAgent(agentId)" class="agent-item"
          :style="selectedAgent === agentId ? 'background:#eff6ff;border-left:3px solid #3b82f6;' : 'border-left:3px solid transparent;'"
          style="padding:9px 14px;cursor:pointer;display:flex;align-items:center;gap:8px;font-size:13px;border-bottom:1px solid #f8fafc;transition:background .15s">
          <div :style="info.online ? 'background:#22c55e' : 'background:#94a3b8'" style="width:8px;height:8px;border-radius:50%;flex-shrink:0;transition:background .3s"></div>
          <span :style="selectedAgent === agentId ? 'color:#1d4ed8;font-weight:600' : 'color:#374151'" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1" :title="agentId">{{ shortId(agentId) }}</span>
          <span style="margin-left:auto;font-size:11px;color:#94a3b8;background:#f1f5f9;padding:1px 5px;border-radius:8px">{{ msgCountFor(agentId) }}</span>
        </div>
        <div v-if="agentCount === 0" style="padding:24px 14px;font-size:12px;color:#94a3b8;text-align:center">
          <div style="margin-bottom:6px">No agents</div>
          <div style="font-size:11px">Waiting for connections...</div>
        </div>
      </div>

      <!-- Search box -->
      <div style="padding:10px 12px;border-top:1px solid #f1f5f9">
        <input v-model="searchQuery" type="text" placeholder="🔍 Search messages..." style="width:100%;background:#f8fafc;border:1px solid #e2e8f0;border-radius:6px;padding:6px 10px;font-size:12px;color:#1e293b;outline:none">
      </div>
    </aside>

    <!-- Message area -->
    <main ref="messageArea" style="flex:1;overflow-y:auto;padding:14px 16px">
      <div v-if="filteredMessages.length === 0" style="display:flex;align-items:center;justify-content:center;height:100%;flex-direction:column;gap:12px;color:#94a3b8">
        <span style="font-size:40px">{{ connected ? '📭' : '🔌' }}</span>
        <span style="font-size:14px;font-weight:500">{{ connected ? 'No messages yet' : 'Not connected' }}</span>
        <span v-if="!connected" style="font-size:12px">Click ⚙️ Settings to configure and connect</span>
        <span v-if="searchQuery && connected" style="font-size:12px">No messages match "{{ searchQuery }}"</span>
      </div>

      <div v-else style="max-width:900px;margin:0 auto;display:flex;flex-direction:column;gap:6px">
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
    return { thinkingExpanded: false, argsExpanded: {}, resultExpanded: {} }
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
      if (typeof c === 'string') return [{ type: 'text', text: c }];
      return c;
    },
    roleIcon() { return {user:'👤',assistant:'🤖',tool:'⚙️',toolResult:'📤'}[this.role] || '❓' },
    roleText() { return {user:'User',assistant:'Assistant',tool:'Tool',toolResult:'Tool Result'}[this.role] || this.role },
    colors() {
      return {
        user:       { card:'#eff6ff', border:'#bfdbfe', hdr:'#dbeafe', lbl:'#1d4ed8' },
        assistant:  { card:'#ffffff', border:'#e2e8f0', hdr:'#f8fafc', lbl:'#059669' },
        tool:       { card:'#f8fafc', border:'#e2e8f0', hdr:'#f1f5f9', lbl:'#6b7280' },
        toolResult: { card:'#f8fafc', border:'#e2e8f0', hdr:'#f1f5f9', lbl:'#6b7280' }
      }[this.role] || { card:'#fff', border:'#e2e8f0', hdr:'#f8fafc', lbl:'#374151' }
    },
    agentShort() {
      const id = this.message.agentId || ''
      if (!id) return ''
      const parts = id.split(':')
      return parts[parts.length-1] || id
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
    isThinkingExpanded() { return this.thinkingExpanded || this.forceExpandThinking > 0 }
  },
  methods: {
    isArgExpanded(i) { return !!this.argsExpanded[i] || this.forceExpandArgs > 0 },
    isResultExpanded(i) { return !!this.resultExpanded[i] || this.forceExpandArgs > 0 },
    toggleThinking() { this.thinkingExpanded = !this.thinkingExpanded },
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
<div :style="{background:colors.card, border:'1px solid '+colors.border, borderRadius:'8px', boxShadow:'0 1px 2px rgba(0,0,0,.04)', overflow:'hidden'}">
  <!-- Header -->
  <div :style="{background:colors.hdr, borderBottom:'1px solid '+colors.border, padding:'6px 14px', display:'flex', alignItems:'center', justifyContent:'space-between'}">
    <div style="display:flex;align-items:center;gap:8px;font-size:13px">
      <span>{{roleIcon}}</span>
      <span :style="{fontWeight:600,color:colors.lbl}">{{roleText}}</span>
      <span v-if="agentShort" style="font-size:11px;color:#94a3b8;background:#f1f5f9;padding:1px 6px;border-radius:4px;font-family:monospace">{{agentShort}}</span>
      <span style="color:#cbd5e1">·</span>
      <span style="font-size:11px;color:#94a3b8">{{timestamp}}</span>
    </div>
    <div style="display:flex;align-items:center;gap:10px;font-size:11px;color:#94a3b8">
      <span v-if="message.message && message.message.model">{{message.message.model}}</span>
      <span v-if="tokenInfo">📊 {{tokenInfo}}</span>
    </div>
  </div>
  <!-- Content -->
  <div style="padding:10px 14px">
    <div v-for="(item,idx) in content" :key="idx" style="margin-bottom:4px">
      <!-- Text -->
      <div v-if="item.type==='text' && role==='assistant'" class="markdown-body" v-html="renderMd(item.text)"></div>
      <div v-else-if="item.type==='text'" style="font-size:13px;line-height:1.65;white-space:pre-wrap;color:#374151">{{item.text}}</div>
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
    const showConnectPanel = ref(false);
    const tokenInput = ref('');
    const subscribeInput = ref('*');
    const searchQuery = ref('');
    const selectedAgent = ref(null);
    const agents = ref({});
    const messagesByAgent = ref({});
    const allMessages = ref([]);
    const forceExpandThinking = ref(0);
    const forceExpandArgs = ref(0);
    const allThinkingExpanded = ref(false);
    const allArgsExpanded = ref(false);
    const messageArea = ref(null);
    const bottomAnchor = ref(null);
    let ws = null;
    let msgIdCounter = 0;

    const agentCount = computed(() => Object.keys(agents.value).length);

    function msgCountFor(agentId) {
      return (messagesByAgent.value[agentId] || []).length;
    }

    // 最近访问：按最后一条消息时间排序，取 top 5
    const agentColors = ['#3b82f6','#10b981','#8b5cf6','#f59e0b','#ef4444','#06b6d4','#ec4899','#84cc16'];
    const agentColorMap = {};
    let colorIdx = 0;
    function agentColor(agentId) {
      if (!agentColorMap[agentId]) {
        agentColorMap[agentId] = agentColors[colorIdx % agentColors.length];
        colorIdx++;
      }
      return agentColorMap[agentId];
    }
    function lastMsgTime(agentId) {
      const msgs = messagesByAgent.value[agentId] || [];
      if (!msgs.length) return 0;
      return new Date(msgs[msgs.length - 1].timestamp || 0).getTime();
    }
    function relativeTime(ts) {
      const diff = Date.now() - ts;
      if (diff < 60000) return '刚刚';
      if (diff < 3600000) return Math.floor(diff / 60000) + ' 分钟前';
      if (diff < 86400000) return Math.floor(diff / 3600000) + ' 小时前';
      return Math.floor(diff / 86400000) + ' 天前';
    }
    const recentAgents = computed(() => {
      return Object.keys(agents.value)
        .map(id => ({ id, lastTime: lastMsgTime(id), count: msgCountFor(id) }))
        .filter(a => a.count > 0)
        .sort((a, b) => b.lastTime - a.lastTime)
        .slice(0, 5);
    });

    const displayMessages = computed(() => {
      if (selectedAgent.value === null) return allMessages.value;
      return messagesByAgent.value[selectedAgent.value] || [];
    });

    const filteredMessages = computed(() => {
      // Build toolResult map
      const toolResultMap = new Map();
      displayMessages.value.forEach(msg => {
        if (msg.message?.role === 'toolResult') {
          const cid = msg.message.toolCallId;
          if (cid) toolResultMap.set(cid, msg.message);
        }
      });

      const merged = [];
      displayMessages.value.forEach(msg => {
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
        merged.push(cloned);
      });

      if (!searchQuery.value) return merged;
      const q = searchQuery.value.toLowerCase();
      return merged.filter(m => JSON.stringify(m.message?.content || []).toLowerCase().includes(q));
    });

    function shortId(id) {
      if (!id) return '';
      const parts = id.split(':');
      return parts[parts.length - 1] || id;
    }

    function selectAgent(agentId) {
      selectedAgent.value = agentId;
    }

    function ensureAgent(agentId, online) {
      if (!agents.value[agentId]) {
        agents.value = { ...agents.value, [agentId]: { online: !!online } };
      } else {
        agents.value[agentId] = { ...agents.value[agentId], online: !!online };
      }
      if (!messagesByAgent.value[agentId]) {
        messagesByAgent.value = { ...messagesByAgent.value, [agentId]: [] };
      }
    }

    function addLines(agentId, lines, meta) {
      ensureAgent(agentId, true);
      let added = false;
      lines.forEach(line => {
        let parsed;
        try { parsed = typeof line === 'string' ? JSON.parse(line) : line; } catch { return; }
        if (!parsed || !parsed.message) return;
        const role = parsed.message?.role;
        if (!role || !['user','assistant','tool','toolResult'].includes(role)) return;
        const msgObj = {
          _id: ++msgIdCounter,
          agentId,
          timestamp: parsed.timestamp || meta?.timestamp || new Date().toISOString(),
          message: parsed.message
        };
        messagesByAgent.value[agentId].push(msgObj);
        allMessages.value.push(msgObj);
        added = true;
      });
      if (added) scrollToBottom();
    }

    function scrollToBottom() {
      nextTick(() => {
        if (bottomAnchor.value) bottomAnchor.value.scrollIntoView({ behavior: 'smooth' });
      });
    }

    function handleMsg(data) {
      let msg;
      try { msg = JSON.parse(data); } catch { return; }
      const agentId = msg.agentId || msg.agent || 'unknown';

      if (msg.type === 'agent_status') {
        ensureAgent(msg.id || agentId, msg.online);
        return;
      }
      if (msg.type === 'log' && Array.isArray(msg.lines)) {
        addLines(agentId, msg.lines, msg);
        return;
      }
      if (msg.type === 'session_list') {
        ensureAgent(agentId, true);
        return;
      }
    }

    function connect() {
      if (ws) { try { ws.close(); } catch {} ws = null; }
      connected.value = false;
      const token = tokenInput.value.trim();
      const sub = subscribeInput.value.trim() || '*';
      const proto = location.protocol === 'https:' ? 'wss' : 'ws';
      let url = proto + '://' + location.host + '/ws/console?subscribe=' + encodeURIComponent(sub);
      if (token) url += '&token=' + encodeURIComponent(token);

      ws = new WebSocket(url);
      ws.onopen = () => { connected.value = true; showConnectPanel.value = false; fetchAgents(); };
      ws.onmessage = e => handleMsg(e.data);
      ws.onerror = () => { connected.value = false; };
      ws.onclose = () => { connected.value = false; };
    }
    function reconnect() { connect(); }

    function clearMessages() {
      allMessages.value = [];
      messagesByAgent.value = {};
      Object.keys(agents.value).forEach(id => { messagesByAgent.value[id] = []; });
    }

    async function fetchAgents() {
      try {
        const res = await fetch('/api/agents');
        const data = await res.json();
        if (Array.isArray(data.agents)) data.agents.forEach(a => ensureAgent(a.id, true));
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
      connected, showConnectPanel, tokenInput, subscribeInput, searchQuery,
      selectedAgent, agents, allMessages, agentCount, forceExpandThinking, forceExpandArgs,
      allThinkingExpanded, allArgsExpanded, messageArea, bottomAnchor,
      filteredMessages, msgCountFor, shortId, selectAgent, reconnect, clearMessages,
      toggleAllThinking, toggleAllArgs,
      recentAgents, agentColor, relativeTime
    };
  }
}).mount('#app');
<\/script>
</body>
</html>`
	return part1 + string(vue3JS) + part2 + string(markedJS) + part3
}
