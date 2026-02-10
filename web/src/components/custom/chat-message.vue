<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import MarkdownIt from 'markdown-it';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';

defineOptions({ name: 'ChatMessage' });

interface Props {
  role?: 'user' | 'assistant';
  content: string;
  userAvatar?: string;
  botAvatar?: string;
}

const props = withDefaults(defineProps<Props>(), {
  role: 'assistant',
  userAvatar: 'https://avatars.githubusercontent.com/u/1?v=4',
  botAvatar: 'https://avatars.githubusercontent.com/u/2?v=4'
});

const fullText = ref('');
const typingIndex = ref(0);
const isTyping = ref(false);

const typeNextChar = () => {
  if (!props.content || typingIndex.value >= props.content.length) {
    isTyping.value = false;
    return;
  }

  fullText.value += props.content[typingIndex.value];
  typingIndex.value += 1;

  setTimeout(() => {
    typeNextChar(); // 不使用 requestAnimationFrame + setTimeout 混用，避免过快重复触发
  }, 15);
};

watch(
  () => props.content,
  (newVal, _oldVal) => {
    if (props.role !== 'assistant') {
      fullText.value = newVal;
      return;
    }

    // 如果是 assistant 且内容增加了，就开始或继续打字
    if (newVal.length > typingIndex.value) {
      if (!isTyping.value) {
        isTyping.value = true;
        typeNextChar();
      }
    }
  }
);

onMounted(() => {
  if (props.role === 'assistant') {
    fullText.value = '';
    typingIndex.value = 0;
    isTyping.value = true;
    typeNextChar();
  } else {
    fullText.value = props.content;
  }
});

const md = new MarkdownIt({
  html: true,
  linkify: true,
  highlight: (str: string, lang: string): string => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `
          <div class="code-block">
            <button class="copy-btn" onclick="navigator.clipboard.writeText(${JSON.stringify(str)})">📋</button>
            <pre><code class="hljs language-${lang}">${hljs.highlight(str, { language: lang }).value}</code></pre>
          </div>
        `;
      } catch {}
    }
    return `<pre><code class="hljs">${md.utils.escapeHtml(str)}</code></pre>`;
  }
});

const renderedMarkdown = computed(() => md.render(fullText.value));
</script>

<template>
  <div
    class="chat-message"
    :class="{
      'chat-ai': role === 'assistant',
      'chat-user': role === 'user'
    }"
  >
    <NAvatar v-if="role === 'assistant'" round size="small" :src="botAvatar" />
    <div class="chat-content" v-html="renderedMarkdown" />
    <NAvatar v-if="role === 'user'" round size="small" :src="userAvatar" />
  </div>
</template>

<style scoped>
.chat-message {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  margin: 1rem 0;
}

.chat-ai {
  flex-direction: row;
}

.chat-user {
  flex-direction: row-reverse;
}

.chat-content {
  max-width: 80%;
  /* background: #f4f4f4; */
  border-radius: 8px;
  padding: 0.75rem 1rem;
  line-height: 1.6;
  font-size: 16px;
  white-space: pre-wrap;
  word-break: break-word;
}

.code-block {
  position: relative;
}

.copy-btn {
  position: absolute;
  top: 6px;
  right: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 14px;
}
</style>
