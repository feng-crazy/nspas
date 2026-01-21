import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import type { Conversation } from '../types';
import { getConversations } from '../services/api';

const History: React.FC = () => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const type = searchParams.get('type') || 'all';

  useEffect(() => {
    // 获取对话历史数据
    const fetchConversations = async () => {
      try {
        setLoading(true);
        let data: Conversation[] = [];
        
        // 如果是获取全部对话，需要分别获取每种类型的对话
        if (type === 'all') {
          // 获取所有类型的对话
          const analysisConvs = await getConversations('analysis');
          const mappingConvs = await getConversations('mapping');
          const assistantConvs = await getConversations('assistant');
          data = [...analysisConvs, ...mappingConvs, ...assistantConvs];
        } else {
          // 获取特定类型的对话
          data = await getConversations(type);
        }
        
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
        setConversations([]);
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