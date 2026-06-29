<template>
  <template v-if="info?.hasUnread">
    <span class="w-2 h-2 bg-red-500 rounded-full animate-pulse flex-shrink-0" title="有新消息"></span>
    <span class="text-xs bg-red-100 text-red-600 border border-red-200 rounded-full px-1.5 font-semibold flex-shrink-0">
      {{ info.newCount || '' }}
    </span>
  </template>
</template>

<script setup>
defineProps({
  info: Object   // { hasUnread, newCount }
})
</script>
