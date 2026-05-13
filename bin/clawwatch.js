#!/usr/bin/env node
/**
 * ClawWatch CLI
 * OpenClaw 实时会话监控器 - 命令行管理工具
 */

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const PORT = process.env.PORT || 3939;
const DATA_DIR = path.join(os.homedir(), '.clawwatch');
const PID_FILE = path.join(DATA_DIR, 'clawwatch.pid');
const LOG_FILE = path.join(DATA_DIR, 'clawwatch.log');
const SERVER_FILE = path.join(__dirname, '..', 'web-viewer-server.js');

// 确保数据目录存在
if (!fs.existsSync(DATA_DIR)) {
  fs.mkdirSync(DATA_DIR, { recursive: true });
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function isRunning() {
  if (!fs.existsSync(PID_FILE)) return false;
  const pid = parseInt(fs.readFileSync(PID_FILE, 'utf-8').trim(), 10);
  if (isNaN(pid)) return false;
  try {
    process.kill(pid, 0);
    return pid;
  } catch (e) {
    return false;
  }
}

async function start() {
  const pid = isRunning();
  if (pid) {
    console.log(`⚠️  ClawWatch is already running (PID: ${pid})`);
    console.log(`   Use 'clawwatch restart' to restart`);
    return;
  }

  // 检查前端是否已构建
  const distPath = path.join(__dirname, '..', 'frontend', 'dist', 'index.html');
  if (!fs.existsSync(distPath)) {
    console.log('🔨 Building frontend...');
    try {
      await buildFrontend();
    } catch (e) {
      console.error('❌ Frontend build failed:', e.message);
      process.exit(1);
    }
  }

  console.log('🚀 Starting ClawWatch...');

  const logStream = fs.openSync(LOG_FILE, 'a');
  const child = spawn(process.execPath, [SERVER_FILE], {
    detached: true,
    stdio: ['ignore', logStream, logStream],
    env: { ...process.env, PORT: String(PORT) }
  });
  child.unref();
  fs.closeSync(logStream);
  fs.writeFileSync(PID_FILE, String(child.pid));

  await sleep(2000);

  if (isRunning()) {
    console.log('✅ ClawWatch started successfully!');
    console.log('');
    console.log(`   PID:  ${child.pid}`);
    console.log(`   Port: ${PORT}`);
    console.log(`   URL:  http://localhost:${PORT}`);
    console.log(`   Logs: clawwatch logs`);
  } else {
    console.log('❌ Failed to start ClawWatch');
    console.log('   Check logs: clawwatch logs');
    if (fs.existsSync(PID_FILE)) fs.unlinkSync(PID_FILE);
    process.exit(1);
  }
}

async function stop() {
  const pid = isRunning();
  if (!pid) {
    console.log('⚠️  ClawWatch is not running');
    if (fs.existsSync(PID_FILE)) fs.unlinkSync(PID_FILE);
    return;
  }

  console.log(`🛑 Stopping ClawWatch (PID: ${pid})...`);
  try {
    process.kill(pid, 'SIGTERM');
  } catch (e) { /* already gone */ }

  // 等待进程退出（最多 5 秒）
  for (let i = 0; i < 20; i++) {
    await sleep(250);
    if (!isRunning()) break;
  }

  // 强制杀死
  if (isRunning()) {
    try { process.kill(pid, 'SIGKILL'); } catch (e) { /* ignore */ }
    await sleep(500);
  }

  if (fs.existsSync(PID_FILE)) fs.unlinkSync(PID_FILE);
  console.log('✅ ClawWatch stopped');
}

async function restart() {
  console.log('🔄 Restarting ClawWatch...\n');
  await stop();
  console.log('');
  await start();
}

function status() {
  const pid = isRunning();
  if (pid) {
    console.log('✅ ClawWatch is running');
    console.log('');
    console.log(`   PID:  ${pid}`);
    console.log(`   Port: ${PORT}`);
    console.log(`   URL:  http://localhost:${PORT}`);
    console.log(`   Log:  ${LOG_FILE}`);
  } else {
    console.log('⚠️  ClawWatch is not running');
    console.log('');
    console.log("   Start with: clawwatch start");
  }
}

function logs(follow = false) {
  if (!fs.existsSync(LOG_FILE)) {
    console.log(`⚠️  Log file not found: ${LOG_FILE}`);
    return;
  }

  if (follow) {
    const tail = spawn('tail', ['-f', LOG_FILE], { stdio: 'inherit' });
    tail.on('error', () => {
      // fallback: 用 Node.js 实现
      let size = fs.statSync(LOG_FILE).size;
      fs.watch(LOG_FILE, () => {
        const newSize = fs.statSync(LOG_FILE).size;
        if (newSize > size) {
          const buf = Buffer.alloc(newSize - size);
          const fd = fs.openSync(LOG_FILE, 'r');
          fs.readSync(fd, buf, 0, buf.length, size);
          fs.closeSync(fd);
          process.stdout.write(buf);
          size = newSize;
        }
      });
    });
  } else {
    const content = fs.readFileSync(LOG_FILE, 'utf-8');
    const lines = content.split('\n');
    const last50 = lines.slice(Math.max(0, lines.length - 51)).join('\n');
    console.log('📋 ClawWatch Logs (last 50 lines):');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log(last50);
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('');
    console.log('Follow logs: clawwatch logs --follow');
  }
}

function buildFrontend() {
  return new Promise((resolve, reject) => {
    const frontendDir = path.join(__dirname, '..', 'frontend');
    if (!fs.existsSync(frontendDir)) {
      return reject(new Error('frontend directory not found'));
    }

    // 先安装前端依赖
    const install = spawn('npm', ['install'], {
      cwd: frontendDir,
      stdio: 'inherit',
      shell: process.platform === 'win32'
    });
    install.on('close', code => {
      if (code !== 0) return reject(new Error('npm install failed'));
      // 再构建
      const build = spawn('npm', ['run', 'build'], {
        cwd: frontendDir,
        stdio: 'inherit',
        shell: process.platform === 'win32'
      });
      build.on('close', code2 => {
        if (code2 !== 0) return reject(new Error('npm run build failed'));
        resolve();
      });
    });
  });
}

function showHelp() {
  console.log(`🦞 ClawWatch - OpenClaw Session Monitor

Usage: clawwatch [command]

Commands:
  start           Start the local web monitor server (default)
  stop            Stop the server
  restart         Restart the server
  status          Show server status
  logs            Show recent logs
  logs --follow   Follow logs in real-time
  agent           Push local session logs to a remote Hub
  help            Show this help

Agent Options (clawwatch agent):
  --hub <url>     Hub WebSocket URL (or HUB_URL env)
  --id  <name>    Agent ID, default: hostname (or HUB_AGENT_ID env)
  --token <tok>   Hub agent token (or HUB_AGENT_TOKEN env)
  --dir <path>    Watch directory, default: ~/.openclaw/agents

Examples:
  clawwatch                         # Start local web monitor
  clawwatch agent --hub ws://hub:4848   # Push logs to Hub
  clawwatch stop                    # Stop local server
  clawwatch logs -f                 # Follow logs

URL: http://localhost:${PORT}
`);
}

// 主入口
(async () => {
  const cmd = process.argv[2] || 'start';
  const args = process.argv.slice(3);

  switch (cmd) {
    case 'start':
      await start();
      break;
    case 'agent': {
      const { AgentPusher } = require('../lib/agent-pusher');
      // 解析 --key value 参数
      const argMap = {};
      for (let i = 0; i < args.length - 1; i++) {
        if (args[i].startsWith('--')) argMap[args[i].slice(2)] = args[i + 1];
      }
      new AgentPusher({
        hubUrl:  argMap.hub   || process.env.HUB_URL,
        agentId: argMap.id    || process.env.HUB_AGENT_ID,
        token:   argMap.token || process.env.HUB_AGENT_TOKEN,
        watchDir: argMap.dir  || process.env.HUB_WATCH_DIR,
      }).start();
      break;
    }
    case 'stop':
      await stop();
      break;
    case 'restart':
      await restart();
      break;
    case 'status':
      status();
      break;
    case 'logs':
      logs(args.includes('--follow') || args.includes('-f'));
      break;
    case 'help':
    case '--help':
    case '-h':
      showHelp();
      break;
    default:
      console.log(`❌ Unknown command: ${cmd}\n`);
      showHelp();
      process.exit(1);
  }
})();
