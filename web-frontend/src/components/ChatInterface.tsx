import React, { useState, useRef, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Message, ConversationType } from '../types';
import { chatWithAI } from '../services/api';
import './ChatInterface.css';
import './ChatMessage.css';

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

    // 创建AI消息占位符
    const aiMessage: Message = {
      id: `ai-${Date.now()}`,
      content: '',
      isUser: false,
      createdAt: new Date()
    };
    
    // 添加AI消息占位符到消息列表
    const messagesWithAiPlaceholder = [...updatedMessages, aiMessage];
    setMessages(messagesWithAiPlaceholder);

    try {
      // 调用AI服务获取响应，使用流式更新
      const aiResponse = await chatWithAI(
        updatedMessages,
        conversationType,
        conversationId || undefined,
        (updateData) => {
          // 实时更新AI消息内容
          setMessages(prevMessages => {
            const newMessages = [...prevMessages];
            // 更新最后一条AI消息的内容
            if (newMessages.length > 0) {
              newMessages[newMessages.length - 1] = {
                ...newMessages[newMessages.length - 1],
                content: updateData.full_content
              };
            }
            return newMessages;
          });
          
          // 更新conversationId（如果是新建对话）
          if (!conversationId) {
            setConversationId(updateData.conversation_id);
          }
        }
      );
      
      // 更新消息列表，包含完整AI响应
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
      // 更新错误消息
      setMessages(prevMessages => {
        const newMessages = [...prevMessages];
        // 更新最后一条AI消息为错误消息
        if (newMessages.length > 0) {
          newMessages[newMessages.length - 1] = {
            ...newMessages[newMessages.length - 1],
            content: '抱歉，发送消息失败，请稍后重试。'
          };
        }
        return newMessages;
      });
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
