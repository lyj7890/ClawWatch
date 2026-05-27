<template>
  <div 
    class="rounded border transition shadow-sm"
    :class="{
      'bg-blue-50 border-blue-200 hover:border-blue-300': role === 'user',
      'bg-white border-gray-200 hover:border-gray-300': role === 'assistant' && !errorMessage,
      'bg-red-50/30 border-red-200 hover:border-red-300': role === 'assistant' && errorMessage,
      'bg-gray-50 border-gray-200 hover:border-gray-300': role === 'tool'
    }"
  >
    <!-- Header - 单行紧凑 -->
    <div 
      class="px-4 py-2 border-b flex items-center justify-between text-sm"
      :class="{
        'border-blue-200 bg-blue-100/50': role === 'user',
        'border-gray-200 bg-gray-50': role === 'assistant' && !errorMessage,
        'border-red-200 bg-red-50': role === 'assistant' && errorMessage,
        'border-gray-200': role === 'tool'
      }"
    >
      <div class="flex items-center gap-2">
        <span>{{ roleIcon }}</span>
        <span 
          class="font-medium"
          :class="{
            'text-blue-700': role === 'user',
            'text-green-700': role === 'assistant',
            'text-gray-700': role === 'tool'
          }"
        >{{ roleText }}</span>
        <span class="text-gray-400">·</span>
        <span class="text-gray-500 text-xs">{{ timestamp }}</span>
      </div>
      <div class="flex items-center gap-3 text-xs text-gray-500">
        <span v-if="message.message?.model">{{ message.message.model }}</span>
        <span v-if="tokenInfo">📊 {{ tokenInfo }}</span>
      </div>
    </div>

    <!-- Error Banner -->
    <div 
      v-if="errorMessage" 
      class="px-4 py-2 bg-red-50 border-b border-red-200 flex items-center gap-2"
    >
      <span class="text-red-500 text-sm">⚠️</span>
      <span class="text-red-700 text-xs font-medium">Model Error</span>
      <span class="text-red-600 text-xs">·</span>
      <span class="text-red-600 text-xs">{{ errorMessage }}</span>
      <span v-if="errorProvider" class="text-red-400 text-xs ml-1">({{ errorProvider }}/{{ errorModel }})</span>
    </div>

    <!-- Content -->
    <div class="px-4 py-3">
      <div v-for="(item, index) in content" :key="index">
        <!-- Text -->
        <div 
          v-if="item.type === 'text' && role === 'assistant'" 
          class="text-sm leading-relaxed markdown-body"
          v-html="renderMarkdown(item.text)"
        ></div>
        <div 
          v-else-if="item.type === 'text'" 
          class="text-sm leading-relaxed whitespace-pre-wrap"
        >
          {{ item.text }}
        </div>

        <!-- Thinking -->
        <div v-if="item.type === 'thinking' && showThinking" class="mt-2">
          <button 
            @click="toggleThinking"
            class="flex items-center gap-2 text-sm text-blue-600 hover:text-blue-700"
          >
            <span :class="{ 'rotate-90': isThinkingExpanded }" class="transition-transform">▶</span>
            <span>💭 Thinking</span>
          </button>
          <div 
            v-if="isThinkingExpanded" 
            class="mt-2 p-3 bg-blue-50 border border-blue-200 rounded text-xs text-gray-700 max-h-96 overflow-y-auto leading-relaxed whitespace-pre-wrap"
          >
            {{ item.thinking || item.text }}
          </div>
        </div>

        <!-- Tool Call -->
        <div 
          v-if="item.type === 'toolCall'" 
          class="mt-2 rounded p-3 transition-colors"
          :class="{
            'bg-red-50 border border-red-300': item.result && item.isError,
            'bg-green-50 border border-green-300': item.result && !item.isError,
            'bg-amber-50 border border-amber-200': !item.result
          }"
        >
          <div class="flex items-center justify-between">
            <div 
              class="text-sm font-medium flex items-center gap-1.5"
              :class="{
                'text-red-700': item.result && item.isError,
                'text-green-700': item.result && !item.isError,
                'text-amber-700': !item.result
              }"
            >
              <span v-if="item.result && item.isError">❌</span>
              <span v-else-if="item.result && !item.isError">✅</span>
              <span v-else>🔧</span>
              {{ item.name || 'Tool' }}
            </div>
            <div class="flex items-center gap-2">
              <button 
                v-if="hasArgs(item)"
                @click="toggleArgs(index)"
                class="text-xs text-gray-600 hover:text-gray-700"
              >
                {{ isArgExpanded(index) ? '▼' : '▶' }} args
              </button>
              <button 
                v-if="item.result"
                @click="toggleResult(index)"
                class="text-xs text-gray-600 hover:text-gray-700"
              >
                {{ isResultExpanded(index) ? '▼' : '▶' }} result
              </button>
            </div>
          </div>
          
          <!-- Arguments -->
          <div 
            v-if="isArgExpanded(index) && hasArgs(item)" 
            class="mt-2 p-2 bg-white rounded text-xs font-mono text-gray-700 max-h-48 overflow-y-auto border"
            :class="{
              'border-red-300': item.result && item.isError,
              'border-green-300': item.result && !item.isError,
              'border-amber-200': !item.result
            }"
          >
            <pre>{{ JSON.stringify(item.arguments, null, 2) }}</pre>
          </div>
          <div v-else-if="!hasArgs(item) && !item.result" class="mt-1 text-xs text-gray-500">(no parameters)</div>
          
          <!-- Result -->
          <div v-if="isResultExpanded(index) && item.result" class="mt-2">
            <div 
              v-for="(resultItem, rIdx) in item.result" 
              :key="rIdx"
              :class="{
                'p-2 bg-white border rounded text-xs max-h-48 overflow-y-auto': true,
                'border-red-300 text-red-700': item.isError,
                'border-green-300 text-gray-700': !item.isError
              }"
            >
              <div v-if="resultItem.type === 'text'" class="whitespace-pre-wrap">{{ resultItem.text }}</div>
            </div>
            <!-- Details (diff, etc) -->
            <div 
              v-if="item.resultDetails && item.resultDetails.diff" 
              class="mt-1 p-2 bg-gray-800 text-gray-100 rounded text-xs font-mono max-h-32 overflow-y-auto"
            >
              <pre>{{ item.resultDetails.diff }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { marked } from 'marked'

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true
})

const props = defineProps({
  message: Object,
  showThinking: Boolean,
  forceExpandThinking: Number,
  forceExpandArgs: Number
})

const thinkingExpanded = ref(false)
const argsExpanded = ref({})
const resultExpanded = ref({})

const role = computed(() => props.message.message?.role || 'unknown')
const content = computed(() => {
  const raw = props.message.message?.content || []
  // Filter out the generic failure placeholder when we have an actual errorMessage
  if (props.message.message?.errorMessage) {
    return raw.filter(item => !(item.type === 'text' && item.text === '[assistant turn failed before producing content]'))
  }
  return raw
})
const errorMessage = computed(() => props.message.message?.errorMessage || null)
const errorProvider = computed(() => props.message.message?.provider || null)
const errorModel = computed(() => props.message.message?.model || null)

const roleIcon = computed(() => {
  const icons = {
    user: '👤',
    assistant: '🤖',
    tool: '⚙️'
  }
  return icons[role.value] || '❓'
})

const roleText = computed(() => {
  return role.value.charAt(0).toUpperCase() + role.value.slice(1)
})

const timestamp = computed(() => {
  return new Date(props.message.timestamp).toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const tokenInfo = computed(() => {
  const usage = props.message.message?.usage
  if (!usage || !usage.totalTokens) return null
  return `${usage.totalTokens} tokens`
})

const isThinkingExpanded = computed(() => {
  return thinkingExpanded.value || props.forceExpandThinking > 0
})

function isArgExpanded(index) {
  return argsExpanded.value[index] || props.forceExpandArgs > 0
}

function isResultExpanded(index) {
  return resultExpanded.value[index] || props.forceExpandArgs > 0
}

function toggleThinking() {
  thinkingExpanded.value = !thinkingExpanded.value
}

function hasArgs(item) {
  return item.arguments && Object.keys(item.arguments).length > 0
}

function toggleArgs(index) {
  argsExpanded.value[index] = !argsExpanded.value[index]
}

function toggleResult(index) {
  resultExpanded.value[index] = !resultExpanded.value[index]
}

function renderMarkdown(text) {
  if (!text) return ''
  return marked.parse(text)
}

// 监听批量控制，当折叠所有时重置本地状态
watch(() => props.forceExpandThinking, (newVal, oldVal) => {
  if (newVal < oldVal) {
    thinkingExpanded.value = false
  }
})

watch(() => props.forceExpandArgs, (newVal, oldVal) => {
  if (newVal < oldVal) {
    argsExpanded.value = {}
    resultExpanded.value = {}
  }
})
</script>
