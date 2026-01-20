import React, { useEffect, useState } from 'react';
import type { ConversationType, Conversation, Message } from '../types';
import ConversationHistory from './ConversationHistory';
import ChatInterface from './ChatInterface';
import { getConversations, createConversation, getConversationById } from '../services/api';
import './ConversationLayout.css';

interface ConversationLayoutProps {
  conversationType: ConversationType;
  onSaveTool?: (htmlContent: string) => void;
}

const ConversationLayout: React.FC<ConversationLayoutProps> = ({ conversationType, onSaveTool }) => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null);
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null);
  const [currentMessages, setCurrentMessages] = useState<Message[]>([]);

  // 获取历史会话列表
  const fetchConversations = async () => {
    try {
      setLoading(true);
      const data = await getConversations(conversationType);
      setConversations(data);
    } catch (error) {
      console.error('Failed to fetch conversations:', error);
      // 出错时使用模拟数据
      setConversations([
        {
          id: '1',
          userId: 'user1',
          type: conversationType,
          title: '思维过程分析',
          messages: [
            { id: '1-1', content: '我最近总是感到焦虑', isUser: true, createdAt: new Date('2024-01-18T10:00:00') },
            { id: '1-2', content: '焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应...', isUser: false, createdAt: new Date('2024-01-18T10:01:00') },
            { id: '2-1', content: '我最近总是感到焦虑', isUser: true, createdAt: new Date('2024-01-18T10:00:00') },
            { id: '2-2', content: '焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应...', isUser: false, createdAt: new Date('2024-01-18T10:01:00') },
            { id: '3-1', content: '我最近总是感到焦虑', isUser: true, createdAt: new Date('2024-01-18T10:00:00') },
            { id: '3-2', content: '焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应...', isUser: false, createdAt: new Date('2024-01-18T10:01:00') },
            { id: '4-1', content: '我最近总是感到焦虑', isUser: true, createdAt: new Date('2024-01-18T10:00:00') },
            { id: '4-2', content: '焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应,焦虑是一种常见的情绪反应...', isUser: false, createdAt: new Date('2024-01-18T10:01:00') }
          ],
          createdAt: new Date('2024-01-18T10:00:00'),
          updatedAt: new Date('2024-01-18T10:01:00')
        },
        {
          id: '2',
          userId: 'user1',
          type: conversationType,
          title: '正念冥想的神经科学',
          messages: [
            { id: '2-1', content: '正念冥想对大脑有什么影响？', isUser: true, createdAt: new Date('2024-01-17T15:30:00') },
            { id: '2-2', content: '正念冥想可以改变大脑的结构和功能...', isUser: false, createdAt: new Date('2024-01-17T15:31:00') }
          ],
          createdAt: new Date('2024-01-17T15:30:00'),
          updatedAt: new Date('2024-01-17T15:31:00')
        }
      ]);
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
    setSelectedConversation(null);
    setCurrentConversationId(null);
    setCurrentMessages([]);
  };

  // 选择历史对话
  const handleSelectConversation = async (conversation: Conversation) => {
    try {
      // 从API获取完整的会话详情
      const fullConversation = await getConversationById(conversation.id);
      setSelectedConversation(fullConversation);
      setCurrentConversationId(fullConversation.id);
      setCurrentMessages(fullConversation.messages);
    } catch (error) {
      console.error('Failed to fetch conversation details:', error);
      // 出错时使用列表中的数据
      setSelectedConversation(conversation);
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