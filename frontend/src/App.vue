<template>
  <div class="h-screen flex flex-col bg-gray-50">
    <!-- Header（去掉 session tab 行，只保留工具按钮） -->
    <header class="bg-white border-b border-gray-200 px-6 py-3 shadow-sm flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="text-2xl">🦞</span>
        <h1 class="text-xl font-bold">Claw Watch</h1>
        <div class="flex items-center gap-2 ml-4">
          <div :class="connected ? 'bg-green-500' : 'bg-red-500'" class="w-2 h-2 rounded-full"></div>
          <span class="text-sm text-gray-600">{{ connected ? 'Connected' : 'Disconnected' }}</span>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="toggleSortOrder"
          class="px-3 py-1.5 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded text-xs font-medium transition flex items-center gap-1.5"
        >
          <span>{{ sortOrder === 'desc' ? '🔽' : '🔼' }}</span>
          <span>{{ sortOrder === 'desc' ? 'Newest First' : 'Oldest First' }}</span>
        </button>
        <div class="w-px h-6 bg-gray-300"></div>
        <button
          @click="toggleAutoRefresh"
          class="px-3 py-1.5 rounded text-xs font-medium transition flex items-center gap-1.5 border"
          :class="autoRefresh ? 'bg-green-50 text-green-600 border-green-200 hover:bg-green-100' : 'bg-gray-100 text-gray-400 border-gray-200 hover:bg-gray-200'"
        >
          <span>{{ autoRefresh ? '🟢' : '⏸️' }}</span>
          <span>{{ autoRefresh ? 'Auto Refresh' : 'Paused' }}</span>
        </button>
        <div class="w-px h-6 bg-gray-300"></div>
        <button
          @click="toggleAllThinking"
          class="px-3 py-1.5 bg-blue-50 hover:bg-blue-100 text-blue-600 rounded text-xs font-medium transition border border-blue-200"
        >
          💭 {{ allThinkingExpanded ? 'Collapse' : 'Expand' }} All Thinking
        </button>
        <button
          @click="toggleAllArgs"
          class="px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-600 rounded text-xs font-medium transition border border-amber-200"
        >
          🔧 {{ allArgsExpanded ? 'Hide' : 'Show' }} All Args
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <div class="flex flex-1 overflow-hidden">

      <!-- Sidebar：两级 agent → sessions -->
      <aside class="w-64 bg-white border-r border-gray-200 flex flex-col shadow-sm">

        <!-- Sessions 列表（可滚动，占满剩余空间） -->
        <div class="flex-1 overflow-y-auto">
          <div
            v-for="agentGroup in groupedSessions"
            :key="agentGroup.agent"
          >
            <!-- Agent 分组 header -->
            <button
              @click="toggleAgentGroup(agentGroup.agent)"
              class="w-full flex items-center gap-2 px-3 py-2 text-left hover:bg-gray-50 transition sticky top-0 z-10 bg-gray-50 border-b border-gray-100"
            >
              <span
                class="w-2 h-2 rounded-full flex-shrink-0"
                :class="agentGroup.online ? 'bg-green-500' : 'bg-gray-300'"
              ></span>
              <span
                class="text-xs font-semibold flex-1 truncate"
                :class="agentGroup.online ? 'text-gray-700' : 'text-gray-400'"
              >{{ agentGroup.agent }}</span>
              <!-- agent 级未读汇总 -->
              <span
                v-if="agentGroup.totalUnread > 0"
                class="text-xs bg-red-500 text-white rounded-full px-1.5 font-semibold flex-shrink-0"
              >{{ agentGroup.totalUnread }}</span>
              <span v-else class="text-xs text-gray-400 bg-gray-100 rounded-full px-1.5">{{ agentGroup.sessions.length }}</span>
              <span class="text-gray-400 text-xs transition-transform duration-200" :class="collapsedAgents[agentGroup.agent] ? '' : 'rotate-90'">▶</span>
            </button>

            <!-- Sessions 列表 -->
            <div v-show="!collapsedAgents[agentGroup.agent]">
              <button
                v-for="(session, idx) in agentGroup.sessions"
                :key="session.id"
                @click="switchSession(session.id, agentGroup.agent)"
                class="w-full flex items-center gap-2 px-3 py-2 text-left border-b border-gray-50 transition pl-6 border-l-2"
                :class="{
                  'bg-blue-50 border-l-blue-500': currentSessionId === session.id,
                  'border-l-transparent hover:bg-gray-50': currentSessionId !== session.id
                }"
              >
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1.5">
                    <code class="text-xs font-mono" :class="currentSessionId === session.id ? 'text-blue-700 font-semibold' : 'text-gray-700'">
                      {{ session.id.substring(0, 8) }}
                    </code>
                    <span v-if="idx === 0 && agentGroup.online" class="text-xs bg-green-100 text-green-700 border border-green-200 rounded px-1 leading-4 font-semibold">LATEST</span>
                  </div>
                  <div class="text-xs text-gray-400 mt-0.5">{{ formatTime(session.mtime) }}</div>
                </div>
                <!-- 未读标识 -->
                <span
                  v-if="sessionMessageCounts[session.id]?.hasUnread"
                  class="w-2 h-2 bg-red-500 rounded-full animate-pulse flex-shrink-0"
                  title="有新消息"
                ></span>
                <span
                  v-if="sessionMessageCounts[session.id]?.hasUnread"
                  class="text-xs bg-red-100 text-red-600 border border-red-200 rounded-full px-1.5 font-semibold flex-shrink-0"
                >
                  {{ sessionMessageCounts[session.id]?.newCount || '' }}
                </span>
              </button>
            </div>
          </div>

          <!-- 没有 agent 时 -->
          <div v-if="groupedSessions.length === 0" class="flex flex-col items-center justify-center h-32 text-gray-400 text-sm">
            <span class="text-2xl mb-2">📭</span>
            <span>No sessions found</span>
          </div>
        </div>

        <!-- 底部固定：Filters + Search + Refresh -->
        <div class="border-t border-gray-200 flex-shrink-0">
          <!-- Filters -->
          <div class="p-3 border-b border-gray-100">
            <label class="text-xs text-gray-400 uppercase tracking-wide mb-2 block font-semibold">Filters</label>
            <div class="space-y-1.5">
              <label class="flex items-center gap-2 text-xs cursor-pointer text-gray-700 hover:text-gray-900">
                <input type="checkbox" v-model="filters.user" class="rounded">
                <span>👤 User</span>
              </label>
              <label class="flex items-center gap-2 text-xs cursor-pointer text-gray-700 hover:text-gray-900">
                <input type="checkbox" v-model="filters.assistant" class="rounded">
                <span>🤖 Assistant</span>
              </label>
              <label class="flex items-center gap-2 text-xs cursor-pointer text-gray-700 hover:text-gray-900">
                <input type="checkbox" v-model="filters.system" class="rounded">
                <span>⚙️ System / Context</span>
              </label>
            </div>
          </div>

          <!-- Search -->
          <div class="p-3 border-b border-gray-100">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="🔍 Search messages..."
              class="w-full bg-gray-50 text-gray-900 rounded px-3 py-1.5 text-xs focus:ring-2 focus:ring-blue-500 outline-none border border-gray-200"
            >
          </div>

          <!-- Refresh -->
          <div class="p-3">
            <button
              @click="loadSession()"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-1.5 rounded text-xs font-medium transition shadow-sm"
            >
              🔄 Refresh
            </button>
          </div>
        </div>
      </aside>

      <!-- Message List -->
      <main ref="messageListEl" class="flex-1 overflow-y-auto relative">
        <div v-if="loading" class="flex items-center justify-center h-full">
          <div class="text-gray-600">Loading...</div>
        </div>

        <div v-else-if="filteredMessages.length === 0" class="flex items-center justify-center h-full">
          <div class="text-gray-500">No messages to display</div>
        </div>

        <div v-else class="max-w-5xl mx-auto p-4 space-y-1">
          <MessageCard
            v-for="(msg, index) in filteredMessages"
            :key="index"
            :message="msg"
            :show-thinking="showThinking"
            :force-expand-thinking="forceExpandThinking"
            :force-expand-args="forceExpandArgs"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, reactive } from 'vue'
import MessageCard from './components/MessageCard.vue'

const API_BASE = ''

const connected = ref(false)
const loading = ref(true)
const currentAgent = ref('main')
const currentSessionId = ref('')
const agents = ref([])
const sessions = ref([])                    // 当前 agent 的 sessions（原有逻辑保留）
const allSessions = ref({})                 // { agentName: [session, ...] }
const sessionMessageCounts = reactive({})   // { sessionId: { count, hasUnread, newCount } }
const messages = ref([])
const searchQuery = ref('')
const showThinking = ref(true)
const forceExpandThinking = ref(0)
const forceExpandArgs = ref(0)
const allThinkingExpanded = ref(false)
const allArgsExpanded = ref(false)
const sortOrder = ref('desc')
const autoRefresh = ref(true)
const messageListEl = ref(null)
const collapsedAgents = reactive({})        // { agentName: true/false }

const filters = ref({
  user: true,
  assistant: true,
  tool: true,
  system: true
})

let pollInterval = null

// ── 按 agent 分组的 sessions（sidebar 数据源） ──
const groupedSessions = computed(() => {
  return agents.value.map(agent => {
    const agentSessions = (allSessions.value[agent] || []).slice().sort((a, b) => b.mtime - a.mtime)
    // 简单判断在线：最新 session mtime 在 10 分钟内认为 online
    const latestMtime = agentSessions[0]?.mtime || 0
    const online = Date.now() - latestMtime < 10 * 60 * 1000
    // 汇总该 agent 下所有 session 的未读数
    const totalUnread = agentSessions.reduce((sum, s) => {
      const info = sessionMessageCounts[s.id]
      return sum + (info?.hasUnread ? (info.newCount || 1) : 0)
    }, 0)
    return { agent, sessions: agentSessions, online, totalUnread }
  }).sort((a, b) => {
    // 在线 agent 排前面
    if (a.online !== b.online) return a.online ? -1 : 1
    return 0
  })
})

// ── 消息过滤 ──
const filteredMessages = computed(() => {
  const messagesWithResults = []
  const toolResultMap = new Map()

  messages.value.forEach(msg => {
    if (msg.type === 'message' && msg.message?.role === 'toolResult') {
      const callId = msg.message.toolCallId
      if (callId) toolResultMap.set(callId, msg.message)
    }
  })

  messages.value.forEach(msg => {
    if (msg.type !== 'message') {
      messagesWithResults.push(msg)
      return
    }

    const role = msg.message?.role || 'unknown'
    if (role === 'toolResult') return

    const clonedMsg = JSON.parse(JSON.stringify(msg))

    if (role === 'assistant' && clonedMsg.message?.content) {
      clonedMsg.message.content = clonedMsg.message.content.map(item => {
        if (item.type === 'toolCall' && item.id) {
          const result = toolResultMap.get(item.id)
          if (result) {
            let hasError = result.isError
            if (!hasError && result.details?.status === 'error') hasError = true
            if (!hasError && result.content) {
              const contentStr = JSON.stringify(result.content)
              if (contentStr.includes('"status":"error"') || contentStr.includes('"status": "error"')) hasError = true
            }
            return { ...item, result: result.content, resultDetails: result.details, isError: hasError }
          }
        }
        return item
      })
    }

    messagesWithResults.push(clonedMsg)
  })

  let filtered = messagesWithResults.filter(msg => {
    if (msg.type !== 'message') return filters.value.system

    const role = msg.message?.role || 'unknown'
    if (role === 'user' && !filters.value.user) return false
    if (role === 'assistant' && !filters.value.assistant) return false

    if (searchQuery.value) {
      const content = JSON.stringify(msg.message?.content || []).toLowerCase()
      const metaStr = JSON.stringify(msg).toLowerCase()
      if (!content.includes(searchQuery.value.toLowerCase()) && !metaStr.includes(searchQuery.value.toLowerCase())) return false
    }

    return true
  })

  return filtered.sort((a, b) => {
    const timeA = new Date(a.timestamp).getTime()
    const timeB = new Date(b.timestamp).getTime()
    return sortOrder.value === 'desc' ? timeB - timeA : timeA - timeB
  })
})

// ── Agent 分组折叠 ──
function toggleAgentGroup(agent) {
  collapsedAgents[agent] = !collapsedAgents[agent]
}

function toggleAutoRefresh() { autoRefresh.value = !autoRefresh.value }
function toggleSortOrder() { sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc' }

function toggleAllThinking() {
  if (allThinkingExpanded.value) { forceExpandThinking.value--; allThinkingExpanded.value = false }
  else { forceExpandThinking.value++; allThinkingExpanded.value = true }
}

function toggleAllArgs() {
  if (allArgsExpanded.value) { forceExpandArgs.value--; allArgsExpanded.value = false }
  else { forceExpandArgs.value++; allArgsExpanded.value = true }
}

// ── 加载 agents ──
async function loadAgents() {
  try {
    const res = await fetch(`${API_BASE}/api/agents`)
    const data = await res.json()
    agents.value = data.agents

    // 初始化 collapsedAgents：offline agent 默认折叠（首次加载时先全部展开，轮询后再判断）
    data.agents.forEach(agent => {
      if (collapsedAgents[agent] === undefined) collapsedAgents[agent] = true  // 默认折叠
    })
  } catch (error) {
    console.error('Failed to load agents:', error)
  }
}

// ── 加载所有 agent 的 sessions ──
async function loadAllSessions() {
  for (const agent of agents.value) {
    try {
      const res = await fetch(`${API_BASE}/api/sessions?agent=${agent}`)
      const data = await res.json()
      allSessions.value = { ...allSessions.value, [agent]: data.sessions || [] }
    } catch (e) {
      console.error(`Failed to load sessions for ${agent}:`, e)
    }
  }
}

// ── 加载某个 agent 的 sessions（兼容旧逻辑） ──
async function loadSessions() {
  try {
    const res = await fetch(`${API_BASE}/api/sessions?agent=${currentAgent.value}`)
    const data = await res.json()
    sessions.value = data.sessions
    allSessions.value = { ...allSessions.value, [currentAgent.value]: data.sessions || [] }
  } catch (error) {
    console.error('Failed to load sessions:', error)
  }
}

// ── 加载消息 ──
async function loadSession(sessionId = null) {
  try {
    loading.value = true

    let sessionPath
    let targetAgent = currentAgent.value

    if (sessionId) {
      // 在所有 agent 的 sessions 里找
      for (const [agent, agentSessions] of Object.entries(allSessions.value)) {
        const found = agentSessions.find(s => s.id === sessionId)
        if (found) {
          sessionPath = found.path
          currentSessionId.value = sessionId
          currentAgent.value = agent
          targetAgent = agent
          break
        }
      }
    }

    if (!sessionPath) {
      const sessionRes = await fetch(`${API_BASE}/api/latest-session?agent=${targetAgent}`)
      const session = await sessionRes.json()

      if (session.error) {
        messages.value = []
        currentSessionId.value = ''
        return
      }

      sessionPath = session.path
      const sessionFile = session.path.split('/').pop()
      currentSessionId.value = sessionFile.replace('.jsonl', '')
    }

    const dataRes = await fetch(`${API_BASE}/api/session?path=${encodeURIComponent(sessionPath)}`)
    const { data } = await dataRes.json()

    messages.value = data
    connected.value = true

    if (currentSessionId.value) {
      sessionMessageCounts[currentSessionId.value] = {
        count: data.length,
        hasUnread: false,
        newCount: 0
      }
    }
  } catch (error) {
    console.error('Failed to load session:', error)
    connected.value = false
    currentSessionId.value = ''
  } finally {
    loading.value = false
  }
}

async function switchSession(sessionId, agent) {
  currentAgent.value = agent
  await loadSession(sessionId)
}

function formatTime(mtime) {
  if (!mtime) return ''
  const now = Date.now()
  const diff = now - mtime
  if (diff < 60000) return 'just now'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
  if (diff < 172800000) return 'yesterday'
  const date = new Date(mtime)
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// ── 轮询：刷新所有 sessions + 当前 session 消息 ──
function startPolling() {
  if (pollInterval) clearInterval(pollInterval)

  pollInterval = setInterval(async () => {
    try {
      await loadAllSessions()

      // 更新 offline agent 的折叠状态
      groupedSessions.value.forEach(g => {
        if (!g.online && collapsedAgents[g.agent] === false && g.sessions.length > 0) {
          // 如果 agent 变 offline 且用户没主动操作过，默认折叠
          // 只在首次检测到 offline 时折叠
        }
      })

      // 检查所有 sessions 的未读状态
      for (const [agent, agentSessions] of Object.entries(allSessions.value)) {
        for (const session of agentSessions) {
          try {
            const dataRes = await fetch(`${API_BASE}/api/session?path=${encodeURIComponent(session.path)}`)
            const { data } = await dataRes.json()
            const currentCount = data.length
            const lastInfo = sessionMessageCounts[session.id]

            if (session.id === currentSessionId.value) {
              if (autoRefresh.value && data.length > messages.value.length) {
                messages.value = data
              }
              sessionMessageCounts[session.id] = { count: currentCount, hasUnread: false, newCount: 0 }
            } else {
              if (!lastInfo) {
                sessionMessageCounts[session.id] = { count: currentCount, hasUnread: false, newCount: 0 }
              } else {
                const newCount = Math.max(0, currentCount - lastInfo.count)
                const hasNewMessages = newCount > 0
                sessionMessageCounts[session.id] = {
                  count: currentCount,
                  hasUnread: hasNewMessages || lastInfo.hasUnread,
                  newCount: hasNewMessages ? (lastInfo.newCount || 0) + newCount : lastInfo.newCount || 0
                }
              }
            }
          } catch (e) {
            // 忽略单个 session 的错误
          }
        }
      }
    } catch (error) {
      console.error('Poll error:', error)
    }
  }, 3000)
}

onMounted(async () => {
  await loadAgents()
  await loadAllSessions()
  await loadSession()
  startPolling()
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>
