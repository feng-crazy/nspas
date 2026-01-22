import React from 'react';
import type { Conversation, ConversationType } from '../types';
import './ConversationLayout.css';

interface ConversationHistoryProps {
  conversations: Conversation[];
  selectedConversationId: string | null;
  conversationType: ConversationType;
  onSelectConversation: (conversation: Conversation) => void;
  onNewConversation: () => void;
  loading?: boolean;
  isCollapsed?: boolean;
}

const ConversationHistory: React.FC<ConversationHistoryProps> = ({
  conversations,
  selectedConversationId,
  conversationType,
  onSelectConversation,
  onNewConversation,
  loading = false,
  isCollapsed = false
}) => {
  return (
    <div className={`conversation-history-panel ${isCollapsed ? 'collapsed' : ''}`}>
      <div className="conversation-history-header">
        {!isCollapsed && (
          <h2>
            {conversationType === 'analysis' && '🧠 神经科学分析'}
            {conversationType === 'mapping' && '✨ 修行映射'}
            {conversationType === 'assistant' && '🔧 修行小助手'}
          </h2>
        )}
        <button 
          className="new-conversation-button"
          onClick={onNewConversation}
          aria-label="新建对话"
        >
          {isCollapsed ? '+' : '新建对话'}
        </button>
      </div>
      
      <div className="conversation-list">
        {loading ? (
          <div className="loading-conversations">
            {!isCollapsed && <p>加载对话列表中...</p>}
          </div>
        ) : conversations.length === 0 ? (
          <div className="no-conversations">
            {!isCollapsed && <p>暂无对话历史</p>}
            {!isCollapsed && (
              <button 
                className="create-first-conversation-button"
                onClick={onNewConversation}
              >
                开始第一次对话
              </button>
            )}
          </div>
        ) : (
          conversations.map(conversation => (
            <div 
              key={conversation.id} 
              className={`conversation-item ${selectedConversationId === conversation.id ? 'active' : ''}`}
              onClick={() => onSelectConversation(conversation)}
              title={conversation.title}
            >
              {!isCollapsed && (
                <div className="conversation-item-title">
                  {conversation.title}
                </div>
              )}
              {!isCollapsed && (
                <div className="conversation-item-meta">
                  <span className="conversation-item-date">
                    {conversation.updatedAt.toLocaleDateString()}
                  </span>
                  <span className="conversation-item-message-count">
                    {conversation.messages.length} 条消息
                  </span>
                </div>
              )}
              {isCollapsed && (
                <div className="conversation-item-collapsed">
                  <div className="conversation-item-collapsed-icon">💬</div>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ConversationHistory;