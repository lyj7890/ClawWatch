#!/usr/bin/env node
/**
 * Claw Watch - Backend Server
 * OpenClaw 实时会话监控器 - 提供 WebSocket 实时推送和 HTTP API
 */
const http = require('http');
const fs = require('fs');
const path = require('path');
const { WebSocketServer } = require('ws');

const PORT = process.env.PORT || 3939;
const OPENCLAW_DIR = path.join(process.env.HOME, '.openclaw');

// 存储所有活跃的 WebSocket 连接
const clients = new Map();

// 文件监控器
const watchers = new Map();

/**
 * 查找指定 agent 的最新会话文件
 */
function findLatestSession(agentName = 'main') {
  const agentDir = path.join(OPENCLAW_DIR, 'agents', agentName, 'sessions');

  if (!fs.existsSync(agentDir)) {
    return null;
  }

  const files = fs.readdirSync(agentDir)
    .filter(f => f.endsWith('.jsonl') && !f.includes('.trajectory.') && !f.endsWith('.lock'))
    .map(f => ({
      name: f,
      path: path.join(agentDir, f),
      mtime: fs.statSync(path.join(agentDir, f)).mtime
    }))
    .sort((a, b) => b.mtime - a.mtime);

  return files.length > 0 ? files[0] : null;
}

/**
 * 获取所有可用的 agent 列表
 */
function getAgentList() {
  const agentsDir = path.join(OPENCLAW_DIR, 'agents');

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
 * 读取会话文件内容
 */
function readSessionFile(filePath) {
  if (!fs.existsSync(filePath)) {
    return [];
  }

  const content = fs.readFileSync(filePath, 'utf-8');
  const lines = content.split('\n').filter(line => line.trim());

  return lines.map((line, index) => {
    try {
      return JSON.parse(line);
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

            // 发送给客户端
            const client = clients.get(clientId);
            if (client && client.readyState === 1) {
              client.send(JSON.stringify({
                type: 'update',
                data: data
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
    const agentName = url.searchParams.get('agent') || 'main';
    const agentDir = path.join(OPENCLAW_DIR, 'agents', agentName, 'sessions');

    if (!fs.existsSync(agentDir)) {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ sessions: [] }));
      return;
    }

    const sessions = fs.readdirSync(agentDir)
      .filter(f => f.endsWith('.jsonl') && !f.includes('.trajectory.') && !f.endsWith('.lock'))
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
    const agentName = url.searchParams.get('agent') || 'main';
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
  console.log(`
╭────────────────────────────────────────────────────╮
│  🦞 Claw Watch - Server Running                    │
├────────────────────────────────────────────────────┤
│  🌐 URL: http://localhost:${PORT}                   │
│  📁 OpenClaw: ${OPENCLAW_DIR}
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
