import React, { useEffect, useState } from 'react';
import type { ConversationType, Conversation, Message } from '../types';
import ConversationHistory from './ConversationHistory';
import ChatInterface from './ChatInterface';
import { getConversations, getConversationById } from '../services/api';
import './ConversationLayout.css';

interface ConversationLayoutProps {
  conversationType: ConversationType;
  onSaveTool?: (htmlContent: string) => void;
}

const ConversationLayout: React.FC<ConversationLayoutProps> = ({ conversationType, onSaveTool }) => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null);
  const [currentMessages, setCurrentMessages] = useState<Message[]>([]);

  // 获取历史会话列表
  const fetchConversations = async () => {
    try {
      setLoading(true);
      const data = await getConversations(conversationType);
      // 转换日期字符串为Date对象
      const formattedConversations = data.map(conv => ({
        ...conv,
        createdAt: new Date(conv.createdAt),
        updatedAt: new Date(conv.updatedAt),
        messages: conv.messages.map(msg => ({
          ...msg,
          createdAt: new Date(msg.createdAt)
        }))
      }));
      setConversations(formattedConversations);
    } catch (error) {
      console.error('Failed to fetch conversations:', error);
      // 出错时显示空列表，让用户知道发生了错误
      setConversations([]);
    } finally {
      setLoading(false);
    }
  };

  // 初始加载会话列表
  useEffect(() => {
    fetchConversations();
  }, [conversationType]);

  // 新建对话
  const handleNewConversation = () => {
    setCurrentConversationId(null);
    setCurrentMessages([]);
  };

  // 选择历史对话
  const handleSelectConversation = async (conversation: Conversation) => {
    try {
      // 从API获取完整的会话详情
      const fullConversation = await getConversationById(conversation.id);
      setCurrentConversationId(fullConversation.id);
      // 转换日期字符串为Date对象
      const formattedMessages = fullConversation.messages.map(msg => ({
        ...msg,
        createdAt: new Date(msg.createdAt)
      }));
      setCurrentMessages(formattedMessages);
    } catch (error) {
      console.error('Failed to fetch conversation details:', error);
      // 出错时使用列表中的数据
      setCurrentConversationId(conversation.id);
      setCurrentMessages(conversation.messages);
    }
  };

  return (
    <div className="conversation-layout">
      {/* 左侧历史对话列表 */}
      <ConversationHistory 
        conversations={conversations}
        selectedConversationId={currentConversationId}
        conversationType={conversationType}
        onSelectConversation={handleSelectConversation}
        onNewConversation={handleNewConversation}
        loading={loading}
      />
      
      {/* 右侧聊天界面 */}
      <div className="chat-interface-panel">
        <ChatInterface 
          conversationType={conversationType} 
          conversationId={currentConversationId}
          messages={currentMessages}
          onSaveTool={onSaveTool}
          onConversationUpdate={(convId, updatedMessages) => {
            setCurrentConversationId(convId);
            setCurrentMessages(updatedMessages);
            // 刷新会话列表，确保最新会话显示在列表中
            fetchConversations();
          }}
        />
      </div>
    </div>
  );
};

export default ConversationLayout;