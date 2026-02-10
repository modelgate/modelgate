<script setup lang="ts">
import { ref } from 'vue';
const emit = defineEmits(['send']);
const text = ref('');

defineOptions({
  name: 'ChatInput'
});

function send() {
  if (!text.value.trim()) return;
  emit('send', text.value);
  text.value = '';
}
</script>

<template>
  <div class="input-wrapper">
    <NInput v-model:value="text" placeholder="输入你的问题..." type="textarea" autosize @keydown.enter.prevent="send" />
    <NButton type="primary" :disabled="!text" @click="send">发送</NButton>
  </div>
</template>

<style scoped>
.input-wrapper {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
  border-top: 1px solid #eee;
}
</style>
