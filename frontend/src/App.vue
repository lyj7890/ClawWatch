<template>
  <div class="h-screen flex flex-col bg-gray-50">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 px-6 py-3 shadow-sm">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-3">
          <span class="text-2xl">🦞</span>
          <h1 class="text-xl font-bold">Claw Watch</h1>
          <div class="flex items-center gap-2 ml-4">
            <div :class="connected ? 'bg-green-500' : 'bg-red-500'" class="w-2 h-2 rounded-full"></div>
            <span class="text-sm text-gray-600">{{ connected ? 'Connected' : 'Disconnected' }}</span>
          </div>
        </div>
        
        <!-- 批量控制按钮 -->
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
      </div>

      <!-- Session Tabs -->
      <div v-if="recentSessions.length > 0" class="flex items-center gap-2 mt-3">
        <button
          v-for="session in recentSessions"
          :key="session.id"
          @click="switchSession(session.id)"
          class="px-3 py-1.5 rounded-t text-xs font-medium transition-colors border-b-2 relative"
          :class="{
            'bg-blue-50 text-blue-700 border-blue-500': currentSessionId === session.id,
            'bg-gray-50 text-gray-600 border-transparent hover:bg-gray-100': currentSessionId !== session.id
          }"
        >
          <div class="flex items-center gap-1.5">
            <span v-if="session.isLatest">🔴</span>
            <code class="font-mono">{{ session.id.substring(0, 8) }}</code>
            <span class="text-gray-400">·</span>
            <span>{{ formatTime(session.mtime) }}</span>
            <!-- 未读标识 -->
            <span 
              v-if="session.hasUnread" 
              class="ml-1 w-2 h-2 bg-red-500 rounded-full animate-pulse"
              title="有新消息"
            ></span>
          </div>
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar -->
      <aside class="w-64 bg-white border-r border-gray-200 flex flex-col shadow-sm">
        <!-- Agent Selector -->
        <div class="p-4 border-b border-gray-200">
          <label class="text-xs text-gray-500 uppercase tracking-wide mb-2 block font-medium">Agent</label>
          <select 
            v-model="currentAgent" 
            @change="onAgentChange"
            class="w-full bg-gray-50 text-gray-900 rounded px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 outline-none border border-gray-200"
          >
            <option v-for="agent in agents" :key="agent" :value="agent">
              {{ agent }}
            </option>
          </select>
        </div>

        <!-- Filters -->
        <div class="p-4 border-b border-gray-200">
          <label class="text-xs text-gray-500 uppercase tracking-wide mb-2 block font-medium">Filters</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-gray-900 text-gray-700">
              <input type="checkbox" v-model="filters.user" class="rounded">
              <span>👤 User</span>
            </label>
            <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-gray-200">
              <input type="checkbox" v-model="filters.assistant" class="rounded">
              <span>🤖 Assistant</span>
            </label>
            <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-gray-200">
              <input type="checkbox" v-model="filters.tool" class="rounded">
              <span>🔧 Tools</span>
            </label>
            <label class="flex items-center gap-2 text-sm cursor-pointer hover:text-gray-200">
              <input type="checkbox" v-model="showThinking" class="rounded">
              <span>💭 Thinking</span>
            </label>
          </div>
        </div>

        <!-- Search -->
        <div class="p-4 border-b border-gray-200">
          <label class="text-xs text-gray-500 uppercase tracking-wide mb-2 block font-medium">Search</label>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Search messages..."
            class="w-full bg-gray-50 text-gray-900 rounded px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 outline-none border border-gray-200"
          >
        </div>

        <!-- Actions -->
        <div class="p-4 mt-auto">
          <button 
            @click="loadSession"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm font-medium transition shadow-sm"
          >
            🔄 Refresh
          </button>
        </div>
      </aside>

      <!-- Message List -->
      <main ref="messageListEl" class="flex-1 overflow-y-auto relative" @scroll="onScroll">
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

      <!-- 暂停自动滚动提示（fixed 居中，不受滚动容器影响） -->
      <div 
        v-if="autoScrollPaused" 
        @click="resumeAutoScroll"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 cursor-pointer bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-full shadow-lg text-sm font-medium transition-all flex items-center gap-2 z-50"
      >
        <span>⏸</span>
        <span>自动滚动已暂停</span>
        <span class="text-blue-200">· 点击恢复</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import MessageCard from './components/MessageCard.vue'

// API 在同端口
const API_BASE = ''

const connected = ref(false)
const loading = ref(true)
const currentAgent = ref('main')
const currentSessionId = ref('')
const agents = ref([])
const sessions = ref([])
const sessionMessageCounts = ref({}) // 记录每个 session 的消息数量
const messages = ref([])
const searchQuery = ref('')
const showThinking = ref(true)
const forceExpandThinking = ref(0)
const forceExpandArgs = ref(0)
const allThinkingExpanded = ref(false)
const allArgsExpanded = ref(false)
const sortOrder = ref('desc') // 'desc' = 最新在上, 'asc' = 最旧在上
const autoScrollPaused = ref(false)
const messageListEl = ref(null)

// 最近活跃的 sessions (最多5个)
const recentSessions = computed(() => {
  return sessions.value.slice(0, 5).map((session, index) => {
    const unreadInfo = sessionMessageCounts.value[session.id]
    return {
      ...session,
      isLatest: index === 0,
      hasUnread: unreadInfo?.hasUnread === true && session.id !== currentSessionId.value
    }
  })
})

const filters = ref({
  user: true,
  assistant: true,
  tool: true
})

// 轮询
let pollInterval = null

const filteredMessages = computed(() => {
  // 先把 toolResult 合并到对应的 assistant 消息中
  const messagesWithResults = []
  const toolResultMap = new Map()
  
  // 第一遍：收集所有 toolResult
  messages.value.forEach(msg => {
    if (msg.type === 'message' && msg.message?.role === 'toolResult') {
      const callId = msg.message.toolCallId
      if (callId) {
        toolResultMap.set(callId, msg.message)
      }
    }
  })
  
  // 第二遍：处理消息并附加 toolResult
  messages.value.forEach(msg => {
    if (msg.type !== 'message') return
    const role = msg.message?.role || 'unknown'
    
    // 跳过单独的 toolResult（已经合并到 assistant 中）
    if (role === 'toolResult') return
    
    // 克隆消息避免修改原始数据
    const clonedMsg = JSON.parse(JSON.stringify(msg))
    
    // 如果是 assistant 消息，附加对应的 toolResult
    if (role === 'assistant' && clonedMsg.message?.content) {
      clonedMsg.message.content = clonedMsg.message.content.map(item => {
        if (item.type === 'toolCall' && item.id) {
          const result = toolResultMap.get(item.id)
          if (result) {
            // 检查是否是错误：isError 字段 或 details.status === 'error' 或 content 里包含 error
            let hasError = result.isError
            if (!hasError && result.details?.status === 'error') {
              hasError = true
            }
            if (!hasError && result.content) {
              // 检查 content 里是否有 error 标记
              const contentStr = JSON.stringify(result.content)
              if (contentStr.includes('"status":"error"') || contentStr.includes('"status": "error"')) {
                hasError = true
              }
            }
            
            return {
              ...item,
              result: result.content,
              resultDetails: result.details,
              isError: hasError
            }
          }
        }
        return item
      })
    }
    
    messagesWithResults.push(clonedMsg)
  })
  
  // 角色过滤
  let filtered = messagesWithResults.filter(msg => {
    const role = msg.message?.role || 'unknown'
    
    if (role === 'user' && !filters.value.user) return false
    if (role === 'assistant' && !filters.value.assistant) return false
    
    // 搜索过滤
    if (searchQuery.value) {
      const content = JSON.stringify(msg.message?.content || []).toLowerCase()
      if (!content.includes(searchQuery.value.toLowerCase())) return false
    }
    
    return true
  })
  
  // 排序
  return filtered.sort((a, b) => {
    const timeA = new Date(a.timestamp).getTime()
    const timeB = new Date(b.timestamp).getTime()
    return sortOrder.value === 'desc' ? timeB - timeA : timeA - timeB
  })
})

// 自动滚动控制
function onScroll() {
  const el = messageListEl.value
  if (!el) return
  // 判断是否在底部（距离底部 80px 以内算"在底部"）
  const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 80
  if (!atBottom && !autoScrollPaused.value) {
    autoScrollPaused.value = true
  }
}

function resumeAutoScroll() {
  autoScrollPaused.value = false
  // 如果是 desc 排序（最新在上），滚到顶部；asc（最新在下），滚到底部
  const el = messageListEl.value
  if (el) {
    if (sortOrder.value === 'desc') {
      el.scrollTop = 0
    } else {
      el.scrollTop = el.scrollHeight
    }
  }
}

// 排序切换
function toggleSortOrder() {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
}

// 批量展开/折叠
function toggleAllThinking() {
  if (allThinkingExpanded.value) {
    forceExpandThinking.value--
    allThinkingExpanded.value = false
  } else {
    forceExpandThinking.value++
    allThinkingExpanded.value = true
  }
}

function toggleAllArgs() {
  if (allArgsExpanded.value) {
    forceExpandArgs.value--
    allArgsExpanded.value = false
  } else {
    forceExpandArgs.value++
    allArgsExpanded.value = true
  }
}

async function loadAgents() {
  try {
    const res = await fetch(`${API_BASE}/api/agents`)
    const data = await res.json()
    agents.value = data.agents
  } catch (error) {
    console.error('Failed to load agents:', error)
  }
}

async function loadSessions() {
  try {
    const res = await fetch(`${API_BASE}/api/sessions?agent=${currentAgent.value}`)
    const data = await res.json()
    sessions.value = data.sessions
  } catch (error) {
    console.error('Failed to load sessions:', error)
  }
}

async function loadSession(sessionId = null) {
  try {
    loading.value = true
    
    // 找到对应的 session
    let sessionPath
    if (sessionId && sessions.value.length > 0) {
      const session = sessions.value.find(s => s.id === sessionId)
      if (session) {
        sessionPath = session.path
        currentSessionId.value = sessionId
      }
    }
    
    // 如果没找到，使用最新会话
    if (!sessionPath) {
      const sessionRes = await fetch(`${API_BASE}/api/latest-session?agent=${currentAgent.value}`)
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
    
    // 读取会话内容
    const dataRes = await fetch(`${API_BASE}/api/session?path=${encodeURIComponent(sessionPath)}`)
    const { data } = await dataRes.json()
    
    messages.value = data
    connected.value = true
    
    // 标记为已读（清除未读标识）
    if (currentSessionId.value) {
      sessionMessageCounts.value[currentSessionId.value] = {
        count: data.length,
        hasUnread: false
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

function onAgentChange() {
  sessions.value = []
  loadSessions()
  loadSession()
}

async function switchSession(sessionId) {
  await loadSession(sessionId)
}

function formatTime(mtime) {
  const date = new Date(mtime)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function startPolling() {
  if (pollInterval) clearInterval(pollInterval)
  
  pollInterval = setInterval(async () => {
    try {
      // 刷新 sessions 列表
      await loadSessions()
      
      // 检查所有最近 sessions 的消息数量
      for (const session of recentSessions.value) {
        try {
          const dataRes = await fetch(`${API_BASE}/api/session?path=${encodeURIComponent(session.path)}`)
          const { data } = await dataRes.json()
          
          const currentCount = data.length
          const lastInfo = sessionMessageCounts.value[session.id]
          
          // 如果是当前 session，直接更新消息
          if (session.id === currentSessionId.value) {
            if (data.length > messages.value.length) {
              const wasAtBottom = (() => {
                const el = messageListEl.value
                if (!el) return true
                return el.scrollHeight - el.scrollTop - el.clientHeight < 80
              })()
              messages.value = data
              // 只有没暂停且原来在底部时才自动滚动
              if (!autoScrollPaused.value && wasAtBottom) {
                // Vue nextTick 后滚动
                setTimeout(() => {
                  const el = messageListEl.value
                  if (el) {
                    if (sortOrder.value === 'desc') {
                      el.scrollTop = 0
                    } else {
                      el.scrollTop = el.scrollHeight
                    }
                  }
                }, 50)
              }
            }
            sessionMessageCounts.value[session.id] = {
              count: currentCount,
              hasUnread: false
            }
          } else {
            // 其他 session
            if (!lastInfo) {
              // 第一次记录，不算未读
              sessionMessageCounts.value[session.id] = {
                count: currentCount,
                hasUnread: false
              }
            } else {
              // 检查是否有新消息
              const hasNewMessages = currentCount > lastInfo.count
              sessionMessageCounts.value[session.id] = {
                count: currentCount,
                hasUnread: hasNewMessages || lastInfo.hasUnread // 保持未读状态，直到切换过去
              }
            }
          }
        } catch (err) {
          console.error(`Failed to poll session ${session.id}:`, err)
        }
      }
    } catch (error) {
      console.error('Poll error:', error)
    }
  }, 3000)
}

onMounted(async () => {
  await loadAgents()
  await loadSessions()
  await loadSession()
  startPolling()
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>
