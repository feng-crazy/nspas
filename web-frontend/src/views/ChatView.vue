<template>
  <div class="chat-container">
    <div class="chat-header">
      <h2>神经科学修行助手</h2>
      <p>探索大脑奥秘，理解思维本质</p>
    </div>
    
    <div class="chat-messages" ref="messagesContainer">
      <div 
        v-for="message in messages" 
        :key="message.id"
        :class="['message', message.role]"
      >
        <div class="message-content">
          <div class="sender">{{ message.role === 'user' ? '我' : '助手' }}</div>
          <div class="text">{{ message.content }}</div>
        </div>
      </div>
    </div>
    
    <div class="chat-input">
      <textarea 
        v-model="inputMessage" 
        placeholder="请输入您的问题或想法..."
        @keydown.enter.exact.prevent="sendMessage"
        :disabled="isLoading"
      ></textarea>
      <div class="input-actions">
        <button @click="clearHistory" class="clear-btn">清除历史</button>
        <button @click="sendMessage" :disabled="!inputMessage.trim() || isLoading">
          {{ isLoading ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, nextTick, onMounted } from 'vue'
import { chatAPI } from '../services/api'

export default {
  name: 'ChatView',
  setup() {
    const inputMessage = ref('')
    const messagesContainer = ref(null)
    const isLoading = ref(false)
    
    const messages = reactive([
      {
        id: 1,
        role: 'assistant',
        content: '您好！我是您的神经科学修行助手。请告诉我您当前的心理状态或想要探讨的哲学问题，我会从神经科学的角度为您解析。'
      }
    ])
    
    // Load chat history on mount
    onMounted(async () => {
      try {
        const response = await chatAPI.getHistory()
        if (response.data.messages && response.data.messages.length > 0) {
          messages.splice(0, messages.length) // Clear initial message
          response.data.messages.forEach(msg => {
            messages.push({
              id: msg.id || Date.now() + Math.random(),
              role: msg.role,
              content: msg.message
            })
          })
          await scrollToBottom()
        }
      } catch (error) {
        console.error('Failed to load chat history:', error)
      }
    })
    
    const sendMessage = async () => {
      if (!inputMessage.value.trim() || isLoading.value) return
      
      // 添加用户消息
      const userMessage = {
        id: Date.now(),
        role: 'user',
        content: inputMessage.value
      }
      
      messages.push(userMessage)
      
      // 清空输入框
      const userContent = inputMessage.value
      inputMessage.value = ''
      
      // 滚动到底部
      await scrollToBottom()
      
      // 设置加载状态
      isLoading.value = true
      
      try {
        // 准备上下文
        const context = messages.slice(0, -1).map(msg => ({
          role: msg.role,
          message: msg.content
        }))
        
        // 调用后端API
        const response = await chatAPI.sendMessage(userContent, context)
        
        const aiResponse = {
          id: Date.now() + 1,
          role: 'assistant',
          content: response.data.response
        }
        messages.push(aiResponse)
        await scrollToBottom()
      } catch (error) {
        console.error('Failed to send message:', error)
        const errorResponse = {
          id: Date.now() + 1,
          role: 'assistant',
          content: '抱歉，处理您的请求时出现错误。请稍后重试。'
        }
        messages.push(errorResponse)
      } finally {
        isLoading.value = false
      }
    }
    
    const scrollToBottom = async () => {
      await nextTick()
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    }
    
    const clearHistory = async () => {
      try {
        await chatAPI.clearHistory()
        messages.splice(0, messages.length)
        messages.push({
          id: Date.now(),
          role: 'assistant',
          content: '对话历史已清除。有什么可以帮您的吗？'
        })
      } catch (error) {
        console.error('Failed to clear history:', error)
      }
    }
    
    return {
      inputMessage,
      messages,
      messagesContainer,
      isLoading,
      sendMessage,
      clearHistory
    }
  }
}
</script>

<style scoped>
.chat-container {
  max-width: 800px;
  margin: 0 auto;
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.chat-header {
  background-color: #42b983;
  color: white;
  padding: 20px;
  text-align: center;
}

.chat-header h2 {
  margin: 0 0 10px;
  font-size: 24px;
}

.chat-header p {
  margin: 0;
  opacity: 0.9;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f9f9f9;
}

.message {
  margin-bottom: 20px;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 18px;
}

.message.user .message-content {
  background-color: #42b983;
  color: white;
  border-bottom-right-radius: 4px;
}

.message.assistant .message-content {
  background-color: #ebebeb;
  color: #333;
  border-bottom-left-radius: 4px;
}

.sender {
  font-size: 12px;
  margin-bottom: 4px;
  opacity: 0.8;
}

.text {
  line-height: 1.5;
}

.chat-input {
  padding: 20px;
  background-color: white;
  border-top: 1px solid #eee;
}

.chat-input textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  resize: none;
  height: 60px;
  font-family: inherit;
  margin-bottom: 10px;
}

.chat-input textarea:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.chat-input button {
  padding: 8px 24px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.chat-input .clear-btn {
  background-color: #999;
}

.chat-input .clear-btn:hover {
  background-color: #777;
}

.chat-input button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.chat-input button:hover:not(:disabled) {
  background-color: #359c6d;
}
</style>