import React, { useState, useRef, useEffect, useMemo } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Message, ConversationType } from '../types';
import { chatWithAI } from '../services/api';
import { useAuth } from '../hooks/useAuth';
import './ChatInterface.css';
import './ChatMessage.css';

interface ChatInterfaceProps {
  conversationType: ConversationType;
  conversationId?: string | null;
  messages?: Message[];
  onSaveTool?: (htmlContent: string, conversationId?: string) => void;
  onConversationUpdate?: (conversationId: string, messages: Message[]) => void;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({ 
  conversationType, 
  conversationId: propConversationId, 
  messages: propMessages, 
  onSaveTool,
}) => {
  const { user } = useAuth();
  const [messages, setMessages] = useState<Message[]>(propMessages || []);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [conversationId, setConversationId] = useState<string | null>(propConversationId || null);
  const [toastMessage, setToastMessage] = useState<string | null>(null);
  const [aiResponseText, setAiResponseText] = useState<string>('');
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const aiMessageIdRef = useRef<string>('');
  
  // 使用 useMemo 缓存最新的 AI 消息
  const latestAiMessage = useMemo(() => {
    return messages.filter(msg => !msg.isUser).pop();
  }, [messages]);
  
  // 显示提示信息
  const showToast = (message: string) => {
    setToastMessage(message);
    setTimeout(() => {
      setToastMessage(null);
    }, 3000);
  };

  // 滚动到底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, aiResponseText]);

  // 从props更新conversationId和messages
  useEffect(() => {
    if (propConversationId !== undefined) {
      setConversationId(propConversationId);
    }
    if (propMessages !== undefined) {
      const processedMessages = propMessages.map((msg, index) => ({
        ...msg,
        id: msg.id && msg.id !== '000000000000000000000000' 
          ? msg.id 
          : generateUniqueId(`${msg.isUser ? 'user-' : 'ai-'}-${index}`)
      }));
      setMessages(processedMessages);
    }
  }, [propConversationId, propMessages]);

  // 生成唯一ID
  const generateUniqueId = (prefix: string = ''): string => {
    return `${prefix}${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  };

  // 发送消息
  const handleSend = async () => {
    if (!input.trim() || isTyping) return;

    setIsTyping(true);
    setAiResponseText('');

    const userMessage: Message = {
      id: generateUniqueId('user-'),
      content: input.trim(),
      isUser: true,
      createdAt: new Date()
    };

    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    setInput('');

    const aiMessageId = generateUniqueId('ai-');
    aiMessageIdRef.current = aiMessageId;
    
    const aiMessage: Message = {
      id: aiMessageId,
      content: '正在思考...',
      isUser: false,
      createdAt: new Date()
    };
    
    const messagesWithAiPlaceholder = [...updatedMessages, aiMessage];
    setMessages(messagesWithAiPlaceholder);

    try {
      if (!user?.id) {
        throw new Error('用户未登录');
      }

      const currentConversationId = conversationId;

      await chatWithAI(
        user.id,
        updatedMessages,
        conversationType,
        currentConversationId || undefined,
        (updateData) => {
          console.log('[ChatInterface] 收到数据:', updateData.full_content);
          setAiResponseText(updateData.full_content);
          
          setMessages(prevMessages => {
            const newMessages = [...prevMessages];
            const aiIndex = newMessages.findIndex(m => m.id === aiMessageId);
            if (aiIndex !== -1) {
              newMessages[aiIndex] = {
                ...newMessages[aiIndex],
                content: updateData.full_content
              };
            }
            return newMessages;
          });
          
          if (!currentConversationId && updateData.conversation_id) {
            setConversationId(updateData.conversation_id);
          }
        }
      );
      
    } catch (error) {
      console.error('Failed to send message:', error);
      setMessages(prevMessages => {
        const newMessages = [...prevMessages];
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
            onClick={() => {
              if (latestAiMessage) {
                const content = latestAiMessage.content;
                
                const isHtml = /<[^>]+>/.test(content);
                
                if (isHtml) {
                  onSaveTool(content, conversationId || undefined);
                } else {
                  console.warn('The latest AI message is not in HTML format');
                  showToast('最新的AI消息不是HTML格式，无法保存为工具');
                }
              } else {
                console.warn('No AI messages found to save as tool');
                showToast('没有找到AI消息来保存为工具');
              }
            }}
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
              ) : message.id === aiMessageIdRef.current && aiResponseText ? (
                <div dangerouslySetInnerHTML={{ __html: aiResponseText }} />
              ) : (
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {message.content}
                </ReactMarkdown>
              )}
            </div>
          </div>
        ))}
        {isTyping && !aiResponseText && (
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
      
      {/* 提示信息 */}
      {toastMessage && (
        <div className="toast-notification">
          {toastMessage}
        </div>
      )}
      
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
