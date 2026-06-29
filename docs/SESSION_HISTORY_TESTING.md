# Testing Guide for Session History Feature

## Test Checklist

### 1. Build and Deploy
```bash
cd ~/.openclaw/ClawWatch/hub
go build -o clawwatch-hub .
# Deploy to your hub server
```

### 2. Test Session List Updates

**Expected Behavior:**
- Agent sends `session_list` message with sessions array
- Hub updates `agent.Sessions` field
- `/api/agents` response includes sessions

**Test Request:**
```bash
curl http://hub:4848/api/agents | jq '.agents[0].sessions'
```

**Expected Response:**
```json
[
  {
    "agentId": "agent:main",
    "sessionId": "session-123",
    "mtime": 1234567890,
    "size": 1024
  }
]
```

### 3. Test Session History Fetch

**Test Request:**
```bash
curl "http://hub:4848/api/session-history?agentId=my-mac&openclawAgentId=agent:main&sessionId=session-123" | jq
```

**Expected Flow:**
1. Hub sends `fetch_session` to agent
2. Agent responds with `session_history` or `session_history_error`
3. HTTP response contains session data or error

**Success Response (200):**
```json
{
  "history": [...],
  "sessionId": "session-123"
}
```

**Error Response (500):**
```json
{
  "error": "session not found"
}
```

**Timeout Response (504):**
```
timeout waiting for agent response
```

### 4. Verify Internal Communication

**Check that session_history messages are NOT broadcast to consoles:**
1. Connect console WebSocket
2. Trigger session history fetch
3. Verify console does NOT receive `session_history` or `session_history_error` messages
4. Verify console DOES receive other agent messages normally

### 5. Test Edge Cases

#### Missing Agent
```bash
curl "http://hub:4848/api/session-history?agentId=nonexistent&openclawAgentId=X&sessionId=Y"
# Expected: 404 agent not found
```

#### Missing Parameters
```bash
curl "http://hub:4848/api/session-history?agentId=my-mac"
# Expected: 400 bad request
```

#### Agent Disconnected During Request
1. Start request
2. Disconnect agent before response
3. Expected: 504 timeout after 10s

#### Send Queue Full
1. Fill agent send buffer
2. Try session history request
3. Expected: 504 timeout with "failed to send request to agent"

## Protocol Messages

### Hub → Agent (fetch_session)
```json
{
  "type": "fetch_session",
  "requestId": "req-1719123456000000",
  "openclawAgentId": "agent:main",
  "sessionId": "session-abc123"
}
```

### Agent → Hub (success)
```json
{
  "type": "session_history",
  "requestId": "req-1719123456000000",
  "data": {
    "history": [...],
    "sessionId": "session-abc123",
    ...
  }
}
```

### Agent → Hub (error)
```json
{
  "type": "session_history_error",
  "requestId": "req-1719123456000000",
  "error": "session not found"
}
```

### Agent → Hub (session list update)
```json
{
  "type": "session_list",
  "sessions": [
    {
      "agentId": "agent:main",
      "sessionId": "session-123",
      "mtime": 1234567890,
      "size": 1024
    }
  ]
}
```

## Debugging Tips

### Check Pending Requests
Add logging in Hub to monitor pending requests:
```go
log.Printf("[debug] pending requests: %d", len(hub.pendingRequests))
```

### Monitor Agent Send Channel
```go
log.Printf("[debug] agent %s send buffer: %d/%d", agent.ID, len(agent.Send), cap(agent.Send))
```

### Verify Message Routing
Enable verbose logging in `handleSessionHistoryResponse`:
```go
log.Printf("[agent] received %s for requestId: %s", msg.Type, msg.RequestID)
```

## Performance Notes

- Max concurrent session history requests: Limited by agent send buffer size (256)
- Request timeout: 10 seconds total (2s send + 10s response)
- No queueing: If agent send buffer is full, request fails immediately
- Memory: Each pending request holds one channel (~96 bytes overhead)

## Security Considerations

- No authentication on `/api/session-history` (add if needed)
- Session data may contain sensitive information
- Consider rate limiting for session history requests
- Agent can DoS hub by sending malicious requestIds (mitigation: bounded map size)
