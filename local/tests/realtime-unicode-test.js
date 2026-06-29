#!/usr/bin/env node

const assert = require('assert');
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawn } = require('child_process');
const WebSocket = require('ws');

const root = path.resolve(__dirname, '..');
const home = fs.mkdtempSync(path.join(os.tmpdir(), 'clawwatch-unicode-'));
const liveURL = process.env.CLAWWATCH_TEST_URL;
const sessionDir = liveURL
  ? home
  : path.join(home, '.openclaw', 'agents', 'main', 'sessions');
const sessionPath = path.join(sessionDir, 'unicode-test.jsonl');
const port = 43939;

fs.mkdirSync(sessionDir, { recursive: true });
fs.writeFileSync(
  sessionPath,
  JSON.stringify({
    type: 'message',
    message: { role: 'user', content: '已有中文历史内容' },
  }) + '\n'
);

const server = liveURL
  ? null
  : spawn(process.execPath, [path.join(root, 'web-viewer-server.js')], {
      cwd: root,
      env: { ...process.env, HOME: home, PORT: String(port) },
      stdio: ['ignore', 'ignore', 'inherit'],
    });

const timeout = setTimeout(() => finish(new Error('Timed out waiting for realtime update')), 8000);
let socket;

function finish(error) {
  clearTimeout(timeout);
  socket?.close();
  server?.kill();
  fs.rmSync(home, { recursive: true, force: true });
  if (error) {
    console.error(error);
    process.exitCode = 1;
  }
}

server?.once('error', finish);

setTimeout(() => {
  socket = new WebSocket(liveURL || `ws://localhost:${port}`);
  socket.once('open', () => {
    socket.send(JSON.stringify({ type: 'watch', path: sessionPath }));
  });

  socket.on('message', raw => {
    const envelope = JSON.parse(raw);
    if (envelope.type === 'watching') {
      fs.appendFileSync(
        sessionPath,
        JSON.stringify({
          type: 'message',
          message: {
            role: 'assistant',
            content: [{ type: 'text', text: '最终回复：中文实时推送正常' }],
          },
        }) + '\n'
      );
      return;
    }

    if (envelope.type === 'update') {
      try {
        assert.strictEqual(
          envelope.data.message.content[0].text,
          '最终回复：中文实时推送正常'
        );
        console.log('unicode realtime update passed');
        finish();
      } catch (error) {
        finish(error);
      }
    }
  });
}, 500);
