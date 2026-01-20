import React, { useState, useRef, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Message, ConversationType } from '../types';
import { chatWithAI, createConversation, generateConversationTitle } from '../services/api';
import './ChatInterface.css';

interface ChatInterfaceProps {
  conversationType: ConversationType;
  conversationId?: string | null;
  messages?: Message[];
  onSaveTool?: (htmlContent: string) => void;
  onConversationUpdate?: (conversationId: string, messages: Message[]) => void;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({ 
  conversationType, 
  conversationId: propConversationId, 
  messages: propMessages, 
  onSaveTool,
  onConversationUpdate
}) => {
  const [messages, setMessages] = useState<Message[]>(propMessages || []);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [conversationId, setConversationId] = useState<string | null>(propConversationId || null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 滚动到底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // 从props更新conversationId和messages
  useEffect(() => {
    if (propConversationId !== undefined) {
      setConversationId(propConversationId);
    }
    if (propMessages !== undefined) {
      setMessages(propMessages);
    }
  }, [propConversationId, propMessages]);

  // 模拟AI回复
  const simulateAIResponse = (_userMessage: string) => {
    setIsTyping(true);

    // 模拟不同类型对话的回复
    const getAIResponse = () => {
      switch (conversationType) {
        case 'analysis':
          return `## 神经科学分析

您的消息触发了以下大脑区域：

### 1. 前额叶皮层
- **功能**：执行控制、决策制定
- **激活程度**：中等

### 2. 杏仁核
- **功能**：情绪处理、恐惧反应
- **激活程度**：低

> 建议：尝试正念冥想，有助于调节前额叶与杏仁核的连接。`;
        case 'mapping':
          return `## 修行映射

### "观呼吸"的神经科学原理

| 脑区 | 功能 | 作用 |
|------|------|------|
| 前扣带回 | 注意力控制 | 维持专注 |
| 岛叶 | 内感受 | 觉察呼吸 |
| 前额叶 | 执行控制 | 抑制分心 |

### 神经可塑性效应
1. 增强注意力网络
2. 提升情绪调节能力
3. 改善自我觉察`;
        case 'assistant':
          return `## 注意力训练工具

我为您设计了一个**数字N-back训练**工具，可以有效提升工作记忆和注意力。

### 训练原理
- 激活前额叶皮层
- 增强工作记忆容量
- 提升注意力持续时间

### 使用方法
1. 选择难度级别（1-back到3-back）
2. 观察屏幕上出现的数字
3. 判断当前数字是否与N步前相同

<div style="border: 1px solid #ccc; padding: 20px; border-radius: 8px; margin: 20px 0;">
  <h3>N-back训练工具</h3>
  <div style="display: flex; flex-direction: column; gap: 10px;">
    <div>
      <label>难度级别：</label>
      <select>
        <option value="1">1-back</option>
        <option value="2">2-back</option>
        <option value="3">3-back</option>
      </select>
    </div>
    <div style="font-size: 48px; text-align: center; margin: 20px 0;">
      5
    </div>
    <div style="display: flex; gap: 10px;">
      <button style="flex: 1; padding: 10px;">相同</button>
      <button style="flex: 1; padding: 10px;">不同</button>
    </div>
  </div>
</div>`;
        default:
          return '感谢您的消息！';
      }
    };

    // 模拟打字机效果
    setTimeout(() => {
      const response = getAIResponse();
      let index = 0;
      const aiMessage: Message = {
        id: `ai-${Date.now()}`,
        content: '',
        isUser: false,
        createdAt: new Date()
      };

      setMessages(prev => [...prev, aiMessage]);

      const typingInterval = setInterval(() => {
        if (index < response.length) {
          setMessages(prev => {
            const updatedMessages = [...prev];
            const lastMessage = updatedMessages[updatedMessages.length - 1];
            if (lastMessage.id === aiMessage.id) {
              lastMessage.content = response.substring(0, index + 1);
            }
            return updatedMessages;
          });
          index++;
        } else {
          clearInterval(typingInterval);
          setIsTyping(false);
        }
      }, 20);
    }, 1000);
  };

  // 发送消息
  const handleSend = async () => {
    if (!input.trim() || isTyping) return;

    setIsTyping(true);

    const userMessage: Message = {
      id: `user-${Date.now()}`,
      content: input.trim(),
      isUser: true,
      createdAt: new Date()
    };

    // 更新本地消息列表，显示用户消息
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    setInput('');

    try {
      // 调用AI服务获取响应
      const aiResponse = await chatWithAI(updatedMessages, conversationType, conversationId || undefined);
      
      // 更新消息列表，包含AI响应
      const finalMessages = aiResponse.messages;
      setMessages(finalMessages);
      
      // 更新conversationId（如果是新建对话）
      if (!conversationId) {
        setConversationId(aiResponse.conversation_id);
      }
      
      // 通知父组件会话已更新
      if (onConversationUpdate) {
        onConversationUpdate(aiResponse.conversation_id, finalMessages);
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      // 显示错误消息
      const errorMessage: Message = {
        id: `ai-${Date.now()}`,
        content: '抱歉，发送消息失败，请稍后重试。',
        isUser: false,
        createdAt: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsTyping(false);
    }
  };

  // 处理键盘事件
  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="chat-interface">
      <div className="chat-header">
        <h2>
          {conversationType === 'analysis' && '🧠 神经科学分析'}
          {conversationType === 'mapping' && '✨ 修行映射'}
          {conversationType === 'assistant' && '🔧 修行小助手'}
        </h2>
        {onSaveTool && (
          <button 
            className="save-tool-button"
            onClick={() => onSaveTool('<div>Sample tool HTML</div>')}
          >
            保存工具
          </button>
        )}
      </div>
      
      <div className="chat-messages">
        {messages.map((message) => (
          <div 
            key={message.id} 
            className={`message ${message.isUser ? 'user-message' : 'ai-message'}`}
          >
            <div className="message-content">
              {message.isUser ? (
                <p>{message.content}</p>
              ) : (
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {message.content}
                </ReactMarkdown>
              )}
            </div>
          </div>
        ))}
        {isTyping && (
          <div className="message ai-message">
            <div className="message-content">
              <div className="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>
      
      <div className="chat-input">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={`请输入您的${conversationType === 'analysis' ? '思维过程' : conversationType === 'mapping' ? '修行语录' : '工具需求'}...`}
          rows={3}
        />
        <button 
          className="send-button"
          onClick={handleSend}
          disabled={!input.trim() || isTyping}
        >
          发送
        </button>
      </div>
    </div>
  );
};

export default ChatInterface;
