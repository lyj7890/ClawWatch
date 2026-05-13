'use strict';

/**
 * AgentPusher
 * 监听本地 OpenClaw session 日志，增量推送到远端 ClawWatch Hub。
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const WebSocket = require('ws');

const RECONNECT_DELAY_MS = 3000;
const MAX_RECONNECT_DELAY_MS = 30000;
const PING_INTERVAL_MS = 25000;
const POLL_INTERVAL_MS = 1000; // fallback poll when fs.watch unreliable

class AgentPusher {
  /**
   * @param {object} opts
   * @param {string} opts.hubUrl        - Hub WebSocket URL，如 ws://hub.example.com:4848
   * @param {string} [opts.agentId]     - Agent 标识，默认 hostname
   * @param {string} [opts.token]       - Hub agent token
   * @param {string} [opts.watchDir]    - 监听目录，默认 ~/.openclaw/agents
   */
  constructor(opts = {}) {
    this.hubUrl = opts.hubUrl;
    this.agentId = opts.agentId || os.hostname();
    this.token = opts.token || '';
    this.watchDir = opts.watchDir || path.join(os.homedir(), '.openclaw', 'agents');

    this.ws = null;
    this.connected = false;
    this.reconnectDelay = RECONNECT_DELAY_MS;
    this.pingTimer = null;
    this.reconnectTimer = null;

    // filePath -> { size, watcher }
    this.fileStates = new Map();
    // dir watcher for detecting new session files
    this.dirWatchers = new Map();
  }

  // ─── Public ──────────────────────────────────────────────

  start() {
    if (!this.hubUrl) {
      console.error('[agent] HUB_URL is required. Use --hub <url> or set HUB_URL env var.');
      process.exit(1);
    }

    console.log(`[agent] Starting ClawWatch Agent`);
    console.log(`[agent]   agentId  : ${this.agentId}`);
    console.log(`[agent]   hub      : ${this.hubUrl}`);
    console.log(`[agent]   watchDir : ${this.watchDir}`);

    this._connect();
    this._watchAgentsDir();

    process.on('SIGINT', () => this._shutdown());
    process.on('SIGTERM', () => this._shutdown());
  }

  // ─── WebSocket ────────────────────────────────────────────

  _buildWsUrl() {
    const u = new URL(`${this.hubUrl}/ws/agent`);
    u.searchParams.set('agentId', this.agentId);
    if (this.token) u.searchParams.set('token', this.token);
    return u.toString();
  }

  _connect() {
    const url = this._buildWsUrl();
    console.log(`[agent] Connecting to ${url.replace(/token=[^&]+/, 'token=***')} ...`);

    this.ws = new WebSocket(url, {
      handshakeTimeout: 10000,
    });

    this.ws.on('open', () => {
      console.log('[agent] ✅ Connected to Hub');
      this.connected = true;
      this.reconnectDelay = RECONNECT_DELAY_MS;
      this._startPing();
      // 发送当前 session 快照
      this._sendSessionList();
    });

    this.ws.on('message', (data) => {
      // Hub 发来的控制消息（预留，如 reload 指令）
      try {
        const msg = JSON.parse(data.toString());
        console.log(`[agent] ← hub: ${msg.type || 'unknown'}`);
      } catch {}
    });

    this.ws.on('close', (code, reason) => {
      console.log(`[agent] ❌ Disconnected (${code}). Reconnecting in ${this.reconnectDelay / 1000}s...`);
      this.connected = false;
      this._stopPing();
      this._scheduleReconnect();
    });

    this.ws.on('error', (err) => {
      console.error(`[agent] WS error: ${err.message}`);
    });
  }

  _scheduleReconnect() {
    clearTimeout(this.reconnectTimer);
    this.reconnectTimer = setTimeout(() => {
      this._connect();
    }, this.reconnectDelay);
    // 指数退避，最大 30s
    this.reconnectDelay = Math.min(this.reconnectDelay * 2, MAX_RECONNECT_DELAY_MS);
  }

  _startPing() {
    this.pingTimer = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.ping();
      }
    }, PING_INTERVAL_MS);
  }

  _stopPing() {
    clearInterval(this.pingTimer);
  }

  _send(obj) {
    if (!this.connected || !this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    try {
      this.ws.send(JSON.stringify(obj));
    } catch (e) {
      console.error('[agent] send error:', e.message);
    }
  }

  // ─── Session List ─────────────────────────────────────────

  _sendSessionList() {
    const sessions = this._getAllSessions();
    this._send({
      type: 'session_list',
      sessions: sessions.map(s => ({
        agent: s.agent,
        id: s.id,
        path: s.path,
        mtime: s.mtime,
        size: s.size,
      })),
    });
    console.log(`[agent] Sent session_list: ${sessions.length} sessions`);

    // 开始监听最新的几个 session 文件
    sessions.slice(0, 10).forEach(s => this._watchFile(s.path, s.agent));
  }

  _getAllSessions() {
    if (!fs.existsSync(this.watchDir)) return [];

    const results = [];
    let agentNames;
    try {
      agentNames = fs.readdirSync(this.watchDir);
    } catch { return []; }

    for (const agentName of agentNames) {
      const sessionsDir = path.join(this.watchDir, agentName, 'sessions');
      if (!fs.existsSync(sessionsDir)) continue;

      let files;
      try {
        files = fs.readdirSync(sessionsDir)
          .filter(f => f.endsWith('.jsonl') && !f.includes('.trajectory') && !f.includes('.checkpoint') && !f.endsWith('.lock'));
      } catch { continue; }

      for (const f of files) {
        const filePath = path.join(sessionsDir, f);
        try {
          const stat = fs.statSync(filePath);
          results.push({
            agent: agentName,
            id: f.replace('.jsonl', ''),
            path: filePath,
            mtime: stat.mtimeMs,
            size: stat.size,
          });
        } catch {}
      }
    }

    return results.sort((a, b) => b.mtime - a.mtime);
  }

  // ─── File Watching ────────────────────────────────────────

  _watchAgentsDir() {
    if (!fs.existsSync(this.watchDir)) {
      console.log(`[agent] watchDir not found: ${this.watchDir}, will retry...`);
      setTimeout(() => this._watchAgentsDir(), 5000);
      return;
    }

    // 监听 agents 目录，检测新 agent / 新 session 目录出现
    try {
      const watcher = fs.watch(this.watchDir, { recursive: true }, (eventType, filename) => {
        if (!filename) return;
        if (filename.endsWith('.jsonl') && !filename.includes('.trajectory') && !filename.includes('.checkpoint') && !filename.endsWith('.lock')) {
          const fullPath = path.join(this.watchDir, filename);
          if (fs.existsSync(fullPath) && !this.fileStates.has(fullPath)) {
            // 新 session 文件出现
            const parts = filename.split(path.sep); // [agentName, 'sessions', 'xxx.jsonl']
            const agentName = parts[0];
            console.log(`[agent] New session detected: ${filename}`);
            this._watchFile(fullPath, agentName);
            // 通知 Hub 有新 session
            this._send({
              type: 'new_session',
              agent: agentName,
              session: path.basename(fullPath, '.jsonl'),
              path: fullPath,
              timestamp: Date.now(),
            });
          }
        }
      });
      this.dirWatchers.set(this.watchDir, watcher);
      console.log(`[agent] Watching agents dir: ${this.watchDir}`);
    } catch (e) {
      console.error(`[agent] Failed to watch dir: ${e.message}`);
      // 降级到定期扫描
      this._startPolling();
    }
  }

  _watchFile(filePath, agentName) {
    if (this.fileStates.has(filePath)) return;

    let size = 0;
    try {
      size = fs.existsSync(filePath) ? fs.statSync(filePath).size : 0;
    } catch {}

    const state = { size, agentName, watcher: null };
    this.fileStates.set(filePath, state);

    try {
      const watcher = fs.watch(filePath, (eventType) => {
        if (eventType === 'change') {
          this._onFileChange(filePath);
        }
      });
      watcher.on('error', () => {}); // 文件被删除时静默处理
      state.watcher = watcher;
    } catch (e) {
      // fs.watch 失败时依赖轮询
    }
  }

  _onFileChange(filePath) {
    const state = this.fileStates.get(filePath);
    if (!state) return;

    try {
      const stat = fs.statSync(filePath);
      if (stat.size <= state.size) return;

      const fd = fs.openSync(filePath, 'r');
      const buf = Buffer.alloc(stat.size - state.size);
      fs.readSync(fd, buf, 0, buf.length, state.size);
      fs.closeSync(fd);

      const newContent = buf.toString('utf-8');
      state.size = stat.size;

      const lines = newContent.split('\n').filter(l => l.trim());
      const parsed = [];
      for (const line of lines) {
        try { parsed.push(JSON.parse(line)); } catch {}
      }

      if (parsed.length > 0) {
        this._send({
          type: 'log',
          agent: state.agentName,
          session: path.basename(filePath, '.jsonl'),
          lines: parsed,
          timestamp: Date.now(),
        });
      }
    } catch (e) {
      // 文件可能被删除，忽略
    }
  }

  // ─── Polling Fallback ─────────────────────────────────────

  _startPolling() {
    console.log('[agent] Using polling mode (fallback)');
    setInterval(() => {
      const sessions = this._getAllSessions();
      sessions.slice(0, 20).forEach(s => {
        const state = this.fileStates.get(s.path);
        if (!state) {
          this._watchFile(s.path, s.agent);
        } else if (s.size > state.size) {
          this._onFileChange(s.path);
        }
      });
    }, POLL_INTERVAL_MS);
  }

  // ─── Shutdown ─────────────────────────────────────────────

  _shutdown() {
    console.log('\n[agent] Shutting down...');
    this._stopPing();
    clearTimeout(this.reconnectTimer);
    for (const [, state] of this.fileStates) {
      if (state.watcher) try { state.watcher.close(); } catch {}
    }
    for (const [, w] of this.dirWatchers) {
      try { w.close(); } catch {}
    }
    if (this.ws) this.ws.close();
    process.exit(0);
  }
}

module.exports = { AgentPusher };
