# Session History Implementation Summary

## Overview
Added session history replay functionality and session list synchronization to ClawWatch Hub.

## Changes Made

### 1. hub.go
- **Added `SessionInfo` struct**: Represents OpenClaw session metadata
  - Fields: `AgentID`, `SessionID`, `MTime`, `Size`
  
- **Extended `AgentConn` struct**: Added `Sessions []SessionInfo` field

- **Extended `Hub` struct**: Added pending request tracking
  - `pendingMu sync.Mutex`: Protects pendingRequests map
  - `pendingRequests map[string]chan []byte`: Tracks requests awaiting responses
  
- **Updated `NewHub()`**: Initializes `pendingRequests` map

### 2. agent.go
- **Updated `agentReadPump()`**: 
  - Added call to `handleSessionHistoryResponse()` before broadcasting
  - Messages identified as `session_history` or `session_history_error` are NOT broadcast to consoles
  
- **Updated `updateAgentMetadata()`**:
  - Now handles both `agent_hello` and `session_list` message types
  - When `type=session_list`, updates `agent.Sessions` field
  
- **Added `handleSessionHistoryResponse()` function**:
  - Detects `session_history` and `session_history_error` messages
  - Extracts `requestId` from response
  - Finds corresponding pending request channel
  - Sends response data to waiting HTTP handler
  - Cleans up pending request entry
  - Returns `true` to prevent broadcast to console

### 3. main.go
- **Added `/api/session-history` endpoint**: New HTTP handler for session history requests

- **Updated `handleAgentsList()`**: Now includes `sessions` field in agent list responses

- **Added `handleSessionHistory()` function**:
  - Parameters: `agentId` (host-level agent), `openclawAgentId`, `sessionId`
  - Generates unique `requestId`
  - Creates response channel and registers in `hub.pendingRequests`
  - Constructs `fetch_session` message: `{"type":"fetch_session","requestId":"xxx","openclawAgentId":"xxx","sessionId":"xxx"}`
  - Sends message to agent via `agent.Send` channel (2s send timeout)
  - Waits for response on pending channel (10s response timeout)
  - Returns session history data or error to HTTP client
  - Cleans up pending request on timeout

## Protocol Flow

### Session History Request Flow
```
1. HTTP Client → Hub: GET /api/session-history?agentId=X&openclawAgentId=Y&sessionId=Z
2. Hub → Agent: {"type":"fetch_session","requestId":"req-123","openclawAgentId":"Y","sessionId":"Z"}
3. Agent → Hub: {"type":"session_history","requestId":"req-123","data":{...}}
4. Hub → HTTP Client: {session history data}
```

### Session List Update Flow
```
1. Agent → Hub: {"type":"session_list","sessions":[...]}
2. Hub: Updates agent.Sessions field
3. Hub → Consoles: (broadcasts normally to subscribed consoles)
```

## Key Design Decisions

1. **Internal Communication**: `session_history` and `session_history_error` messages are NOT broadcast to consoles (internal communication only)

2. **Timeout Management**: 
   - 2s timeout for sending request to agent
   - 10s timeout for receiving response from agent
   - Pending requests are cleaned up on timeout

3. **Thread Safety**: `pendingRequests` map is protected by `pendingMu` mutex

4. **Error Handling**: Both send failures and agent-reported errors are properly handled and returned to HTTP client

## API Specification

### GET /api/session-history
**Query Parameters:**
- `agentId` (required): Host-level agent ID
- `openclawAgentId` (required): OpenClaw agent ID
- `sessionId` (required): Session ID to fetch

**Success Response (200):**
```json
{
  "history": [...],
  "sessionId": "...",
  ...
}
```

**Error Responses:**
- 400: Missing required parameters
- 404: Agent not found
- 504: Timeout (send or response timeout)
- 500: Agent reported error

### GET /api/agents
**Enhanced Response:**
Now includes `sessions` field for each agent:
```json
{
  "agents": [
    {
      "id": "...",
      "hostname": "...",
      "sessions": [
        {
          "agentId": "...",
          "sessionId": "...",
          "mtime": 1234567890,
          "size": 1024
        }
      ],
      ...
    }
  ]
}
```

## Verification
- ✅ Code passes `go vet ./...`
- ✅ Code compiles successfully
- ✅ No race conditions in pending request handling
- ✅ Proper cleanup on timeout/disconnect

## Next Steps (Agent-side implementation needed)
1. Agent must handle `fetch_session` messages
2. Agent must read session history from disk
3. Agent must respond with `session_history` or `session_history_error` message
4. Agent should periodically send `session_list` updates
