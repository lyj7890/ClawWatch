#!/usr/bin/env node
/**
 * Claw Watch - Backend Server
 * OpenClaw / Hermes 实时会话监控器 - 提供 WebSocket 实时推送和 HTTP API
 */
const http = require('http');
const fs = require('fs');
const path = require('path');
const { WebSocketServer } = require('ws');

const PORT = process.env.PORT || 3939;

// ── 模式检测 ──
const MODE = process.argv.includes('--hermes') ? 'hermes' : 'openclaw';
const SESSION_DIR = MODE === 'hermes'
  ? path.join(process.env.HOME, '.hermes', 'sessions')
  : path.join(process.env.HOME, '.openclaw', 'agents');

// 存储所有活跃的 WebSocket 连接
const clients = new Map();

// 文件监控器
const watchers = new Map();

/**
 * Hermes → OpenClaw 格式转换
 * Hermes:  { role, content (string/JSON), timestamp, model, ... }
 * OpenClaw: { type: "message", timestamp, message: { role, content: [{type, text}], model, usage } }
 */
function hermesToOpenClawFormat(msg) {
  const role = msg.role || 'unknown';
  const rawContent = msg.content || '';

  // 构建 content 数组 (OpenClaw 格式)
  let contentArray;
  if (role === 'tool') {
    // Hermes tool 消息: content 是 JSON 字符串 (tool result/error)
    contentArray = [{ type: 'text', text: typeof rawContent === 'string' ? rawContent : JSON.stringify(rawContent) }];
  } else if (role === 'assistant' || role === 'user') {
    // Hermes assistant/user: content 是纯文本或 markdown
    contentArray = [{ type: 'text', text: typeof rawContent === 'string' ? rawContent : JSON.stringify(rawContent) }];
  } else if (role === 'session_meta') {
    // Hermes session_meta: 第一行元数据，显示为 system 消息
    contentArray = [{ type: 'text', text: typeof rawContent === 'string' ? rawContent : JSON.stringify(rawContent) }];
  } else {
    contentArray = [{ type: 'text', text: typeof rawContent === 'string' ? rawContent : JSON.stringify(rawContent) }];
  }

  return {
    type: 'message',
    timestamp: msg.timestamp || msg.created_at || new Date().toISOString(),
    message: {
      role: role === 'session_meta' ? 'system' : role,
      content: contentArray,
      model: msg.model || null,
      usage: msg.usage || null,
      session_id: msg.session_id || null
    }
  };
}

/**
 * 列出所有会话文件 (Hermes: 扁平目录, OpenClaw: agent 子目录)
 */
function listSessionFiles(agentName) {
  if (MODE === 'hermes') {
    if (!fs.existsSync(SESSION_DIR)) return [];
    return fs.readdirSync(SESSION_DIR)
      .filter(f => f.endsWith('.jsonl'))
      .map(f => {
        const filePath = path.join(SESSION_DIR, f);
        const stats = fs.statSync(filePath);
        return {
          name: f,
          path: filePath,
          mtime: stats.mtime,
          size: stats.size
        };
      })
      .sort((a, b) => b.mtime - a.mtime);
  }

  // OpenClaw mode
  const dir = agentName
    ? path.join(SESSION_DIR, agentName, 'sessions')
    : path.join(SESSION_DIR, 'main', 'sessions');

  if (!fs.existsSync(dir)) return [];

  return fs.readdirSync(dir)
    .filter(f => f.endsWith('.jsonl') && !f.includes('.trajectory') && !f.includes('.checkpoint') && !f.endsWith('.lock'))
    .map(f => {
      const filePath = path.join(dir, f);
      const stats = fs.statSync(filePath);
      return {
        name: f,
        path: filePath,
        mtime: stats.mtime,
        size: stats.size
      };
    })
    .sort((a, b) => b.mtime - a.mtime);
}

/**
 * 查找最新的会话文件
 */
function findLatestSession(agentName = 'main') {
  const files = listSessionFiles(agentName);
  return files.length > 0 ? files[0] : null;
}

/**
 * 获取所有可用的 agent 列表
 */
function getAgentList() {
  if (MODE === 'hermes') {
    // Hermes 没有 agent 概念，返回单个 'hermes' agent
    return ['hermes'];
  }

  const agentsDir = SESSION_DIR;

  if (!fs.existsSync(agentsDir)) {
    return [];
  }

  return fs.readdirSync(agentsDir)
    .filter(name => {
      const sessionsDir = path.join(agentsDir, name, 'sessions');
      return fs.existsSync(sessionsDir);
    });
}

/**
 * 归一化 content 字段：确保始终为 [{type, text}] 数组格式
 * OpenClaw 新版 user message 的 content 可能是纯字符串，需要兼容
 */
function normalizeContent(msg) {
  if (!msg || !msg.message) return msg;
  const c = msg.message.content;
  if (typeof c === 'string') {
    msg.message.content = [{ type: 'text', text: c }];
  } else if (Array.isArray(c)) {
    // 已经是数组，保持原样
  } else if (c == null) {
    msg.message.content = [];
  }
  return msg;
}

/**
 * 读取会话文件内容并转为 OpenClaw 兼容格式
 */
function readSessionFile(filePath) {
  if (!fs.existsSync(filePath)) {
    return [];
  }

  const content = fs.readFileSync(filePath, 'utf-8');
  const lines = content.split('\n').filter(line => line.trim());

  return lines.map((line, index) => {
    try {
      const parsed = JSON.parse(line);
      // Hermes 模式: 转换为 OpenClaw 格式; OpenClaw 模式: 归一化 content
      if (MODE === 'hermes') return hermesToOpenClawFormat(parsed);
      return normalizeContent(parsed);
    } catch (e) {
      return null;
    }
  }).filter(Boolean);
}

/**
 * 监控会话文件变化
 */
function watchSessionFile(filePath, clientId) {
  const watchKey = `${clientId}:${filePath}`;

  // 如果已经在监控，先停止
  if (watchers.has(watchKey)) {
    watchers.get(watchKey).close();
  }

  let lastSize = 0;

  if (fs.existsSync(filePath)) {
    lastSize = fs.statSync(filePath).size;
  }

  const watcher = fs.watch(filePath, (eventType) => {
    if (eventType === 'change') {
      const stats = fs.statSync(filePath);

      if (stats.size > lastSize) {
        const content = fs.readFileSync(filePath, 'utf-8');
        const newContent = content.slice(lastSize);
        lastSize = stats.size;

        const newLines = newContent.split('\n').filter(l => l.trim());

        newLines.forEach(line => {
          try {
            const data = JSON.parse(line);

            // 归一化 content 后发送给客户端
            const normalized = MODE === 'hermes' ? hermesToOpenClawFormat(data) : normalizeContent(data);
            const client = clients.get(clientId);
            if (client && client.readyState === 1) {
              client.send(JSON.stringify({
                type: 'update',
                data: normalized
              }));
            }
          } catch (e) {
            // 忽略解析错误
          }
        });
      }
    }
  });

  watchers.set(watchKey, watcher);
}

/**
 * 停止监控会话文件
 */
function unwatchSessionFile(clientId, filePath = null) {
  if (filePath) {
    const watchKey = `${clientId}:${filePath}`;
    if (watchers.has(watchKey)) {
      watchers.get(watchKey).close();
      watchers.delete(watchKey);
    }
  } else {
    // 停止该客户端的所有监控
    for (const [key, watcher] of watchers.entries()) {
      if (key.startsWith(`${clientId}:`)) {
        watcher.close();
        watchers.delete(key);
      }
    }
  }
}

/**
 * HTTP 服务器处理
 */
const server = http.createServer((req, res) => {
  // 设置 CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }

  const url = new URL(req.url, `http://${req.headers.host}`);

  // API: 获取 agent 列表
  if (url.pathname === '/api/agents') {
    const agents = getAgentList();
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ agents }));
    return;
  }

  // API: 获取指定 agent 的所有 sessions
  if (url.pathname === '/api/sessions') {
    const agentName = url.searchParams.get('agent') || (MODE === 'hermes' ? 'hermes' : 'main');

    if (MODE === 'hermes') {
      const files = listSessionFiles();
      const sessions = files.map(f => ({
        id: f.name.replace('.jsonl', ''),
        path: f.path,
        mtime: f.mtime,
        size: f.size
      }));

      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ sessions }));
      return;
    }

    // OpenClaw mode
    const agentDir = path.join(SESSION_DIR, agentName, 'sessions');

    if (!fs.existsSync(agentDir)) {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ sessions: [] }));
      return;
    }

    const sessions = fs.readdirSync(agentDir)
      .filter(f => f.endsWith('.jsonl') && !f.includes('.trajectory') && !f.includes('.checkpoint') && !f.endsWith('.lock'))
      .map(f => {
        const filePath = path.join(agentDir, f);
        const stats = fs.statSync(filePath);
        return {
          id: f.replace('.jsonl', ''),
          path: filePath,
          mtime: stats.mtime,
          size: stats.size
        };
      })
      .sort((a, b) => b.mtime - a.mtime); // 按修改时间倒序

    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ sessions }));
    return;
  }

  // API: 获取指定 agent 的最新会话
  if (url.pathname === '/api/latest-session') {
    const agentName = url.searchParams.get('agent') || (MODE === 'hermes' ? 'hermes' : 'main');
    const session = findLatestSession(agentName);

    if (!session) {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Session not found' }));
      return;
    }

    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      name: session.name,
      path: session.path,
      mtime: session.mtime
    }));
    return;
  }

  // API: 读取会话内容
  if (url.pathname === '/api/session') {
    const filePath = url.searchParams.get('path');

    if (!filePath || !fs.existsSync(filePath)) {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'File not found' }));
      return;
    }

    const data = readSessionFile(filePath);
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ data }));
    return;
  }

  // 提供前端静态文件
  const distPath = path.join(__dirname, 'frontend', 'dist');
  let filePath = path.join(distPath, url.pathname === '/' ? 'index.html' : url.pathname);
  
  // 检查文件是否存在
  if (!fs.existsSync(filePath)) {
    // 如果是 SPA 路由，返回 index.html
    if (!url.pathname.startsWith('/api')) {
      filePath = path.join(distPath, 'index.html');
    } else {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Not Found' }));
      return;
    }
  }
  
  // 确定 MIME 类型
  const ext = path.extname(filePath);
  const mimeTypes = {
    '.html': 'text/html',
    '.js': 'application/javascript',
    '.css': 'text/css',
    '.json': 'application/json',
    '.png': 'image/png',
    '.jpg': 'image/jpeg',
    '.svg': 'image/svg+xml',
  };
  const contentType = mimeTypes[ext] || 'application/octet-stream';
  
  // 读取并返回文件
  try {
    const content = fs.readFileSync(filePath);
    res.writeHead(200, { 'Content-Type': contentType });
    res.end(content);
  } catch (error) {
    res.writeHead(500, { 'Content-Type': 'text/plain' });
    res.end('Internal Server Error');
  }
});

// WebSocket 服务器
const wss = new WebSocketServer({ server });

wss.on('connection', (ws, req) => {
  const clientId = Math.random().toString(36).substr(2, 9);
  clients.set(clientId, ws);

  console.log(`[WebSocket] Client connected: ${clientId}`);

  ws.on('message', (message) => {
    try {
      const msg = JSON.parse(message);

      if (msg.type === 'watch') {
        // 开始监控指定文件
        const { path: filePath } = msg;
        console.log(`[WebSocket] Client ${clientId} watching: ${filePath}`);
        watchSessionFile(filePath, clientId);

        ws.send(JSON.stringify({
          type: 'watching',
          path: filePath
        }));
      } else if (msg.type === 'unwatch') {
        // 停止监控
        const { path: filePath } = msg;
        console.log(`[WebSocket] Client ${clientId} unwatching: ${filePath}`);
        unwatchSessionFile(clientId, filePath);
      }
    } catch (e) {
      console.error('[WebSocket] Error parsing message:', e);
    }
  });

  ws.on('close', () => {
    console.log(`[WebSocket] Client disconnected: ${clientId}`);
    unwatchSessionFile(clientId);
    clients.delete(clientId);
  });

  ws.on('error', (error) => {
    console.error(`[WebSocket] Error for client ${clientId}:`, error);
  });
});

// 启动服务器
server.listen(PORT, () => {
  const modeLabel = MODE === 'hermes' ? 'Hermes' : 'OpenClaw';
  const modeDir = MODE === 'hermes' ? path.join(process.env.HOME, '.hermes', 'sessions') : path.join(process.env.HOME, '.openclaw');
  console.log(`
╭────────────────────────────────────────────────────╮
│  🦞 Claw Watch - Server Running (${modeLabel} mode)        │
├────────────────────────────────────────────────────┤
│  🌐 URL: http://localhost:${PORT}                   │
│  📁 ${modeLabel}: ${modeDir}
│  🔌 WebSocket: ws://localhost:${PORT}               │
╰────────────────────────────────────────────────────╯

Press Ctrl+C to stop
  `);
});

// 优雅退出
process.on('SIGINT', () => {
  console.log('\n\n👋 Shutting down...');

  // 关闭所有监控
  for (const watcher of watchers.values()) {
    watcher.close();
  }

  // 关闭所有 WebSocket 连接
  for (const client of clients.values()) {
    client.close();
  }

  server.close(() => {
    console.log('✅ Server stopped');
    process.exit(0);
  });
});
