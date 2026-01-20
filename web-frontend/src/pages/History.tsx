import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import type { Conversation } from '../types';

const History: React.FC = () => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const type = searchParams.get('type') || 'all';

  useEffect(() => {
    // 模拟获取对话历史数据
    const fetchConversations = async () => {
      try {
        setLoading(true);
        // 这里应该调用API获取真实数据
        // const response = await fetch('/api/conversations');
        // const data = await response.json();
        
        // 模拟数据
        const mockData: Conversation[] = [
          {
            id: '1',
            type: 'analysis',
            title: '思维过程分析',
            messages: [
              { id: '1-1', content: '我最近总是感到焦虑', isUser: true, createdAt: new Date('2024-01-18T10:00:00') },
              { id: '1-2', content: '焦虑是一种常见的情绪反应...', isUser: false, createdAt: new Date('2024-01-18T10:01:00') }
            ],
            createdAt: new Date('2024-01-18T10:00:00'),
            updatedAt: new Date('2024-01-18T10:01:00')
          },
          {
            id: '2',
            type: 'mapping',
            title: '修行语录映射',
            messages: [
              { id: '2-1', content: '"活在当下"的神经科学解释是什么？', isUser: true, createdAt: new Date('2024-01-17T15:30:00') },
              { id: '2-2', content: '"活在当下"涉及大脑的前额叶皮层...', isUser: false, createdAt: new Date('2024-01-17T15:31:00') }
            ],
            createdAt: new Date('2024-01-17T15:30:00'),
            updatedAt: new Date('2024-01-17T15:31:00')
          }
        ];
        
        // 过滤对话类型
        const filteredData = type === 'all' 
          ? mockData 
          : mockData.filter(conv => conv.type === type);
        
        setConversations(filteredData);
      } catch (error) {
        console.error('Failed to fetch conversations:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchConversations();
  }, [type]);

  const handleDelete = async (id: string) => {
    try {
      // 这里应该调用API删除对话
      // await fetch(`/api/conversations/${id}`, { method: 'DELETE' });
      setConversations(conversations.filter(conv => conv.id !== id));
    } catch (error) {
      console.error('Failed to delete conversation:', error);
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'analysis':
        return '🧠 神经科学分析';
      case 'mapping':
        return '✨ 修行映射';
      case 'assistant':
        return '🔧 修行小助手';
      default:
        return '未知类型';
    }
  };

  return (
    <div className="history-page">
      <div className="history-header">
        <h1>📚 对话历史</h1>
        <div className="filter-tabs">
          <button 
            className={`filter-tab ${type === 'all' ? 'active' : ''}`}
            onClick={() => window.location.href = '/history?type=all'}
          >
            全部
          </button>
          <button 
            className={`filter-tab ${type === 'analysis' ? 'active' : ''}`}
            onClick={() => window.location.href = '/history?type=analysis'}
          >
            神经科学分析
          </button>
          <button 
            className={`filter-tab ${type === 'mapping' ? 'active' : ''}`}
            onClick={() => window.location.href = '/history?type=mapping'}
          >
            修行映射
          </button>
          <button 
            className={`filter-tab ${type === 'assistant' ? 'active' : ''}`}
            onClick={() => window.location.href = '/history?type=assistant'}
          >
            修行小助手
          </button>
        </div>
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : conversations.length === 0 ? (
        <div className="no-conversations">
          <p>暂无对话历史</p>
          <Link to="/" className="create-button">开始新对话</Link>
        </div>
      ) : (
        <div className="conversations-list">
          {conversations.map(conversation => (
            <div key={conversation.id} className="conversation-item">
              <div className="conversation-info">
                <div className="conversation-type">{getTypeLabel(conversation.type)}</div>
                <h3 className="conversation-title">{conversation.title}</h3>
                <div className="conversation-meta">
                  <span className="conversation-date">
                    {conversation.updatedAt.toLocaleString()}
                  </span>
                  <span className="message-count">
                    {conversation.messages.length} 条消息
                  </span>
                </div>
              </div>
              <div className="conversation-actions">
                <Link 
                  to={`/${conversation.type}?convId=${conversation.id}`} 
                  className="action-button view-button"
                >
                  继续对话
                </Link>
                <button 
                  className="action-button delete-button"
                  onClick={() => handleDelete(conversation.id)}
                >
                  删除
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default History;