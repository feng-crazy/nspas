import React, { useState } from 'react';
import ConversationLayout from '../components/ConversationLayout';
import { saveTool } from '../services/api';

const Assistant: React.FC = () => {
  const [showSaveModal, setShowSaveModal] = useState(false);
  const [toolName, setToolName] = useState('');
  const [toolDescription, setToolDescription] = useState('');
  const [currentHtmlContent, setCurrentHtmlContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null);

  const handleSaveTool = (htmlContent: string, conversationId?: string) => {
    setCurrentHtmlContent(htmlContent);
    if (conversationId) {
      setCurrentConversationId(conversationId);
    }
    setShowSaveModal(true);
  };

  const handleModalSave = async () => {
    if (!currentConversationId) {
      setError('缺少会话ID，无法保存工具');
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      await saveTool({
        name: toolName,
        description: toolDescription,
        html_content: currentHtmlContent,
        conversation_id: currentConversationId
      });
      
      // 关闭模态框并重置表单
      setShowSaveModal(false);
      setToolName('');
      setToolDescription('');
      setCurrentHtmlContent('');
      setCurrentConversationId(null);
    } catch (err) {
      console.error('Failed to save tool:', err);
      setError('保存工具失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="assistant-page">
      <ConversationLayout 
        conversationType="assistant" 
        onSaveTool={handleSaveTool} 
      />
      
      {/* 保存工具模态框 */}
      {showSaveModal && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>保存修行工具</h3>
            
            {error && <div className="modal-error">{error}</div>}
            
            <div className="modal-form">
              <div className="form-group">
                <label htmlFor="tool-name">工具名称：</label>
                <input
                  type="text"
                  id="tool-name"
                  value={toolName}
                  onChange={(e) => setToolName(e.target.value)}
                  placeholder="输入工具名称"
                  required
                  disabled={loading}
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="tool-description">工具描述：</label>
                <textarea
                  id="tool-description"
                  value={toolDescription}
                  onChange={(e) => setToolDescription(e.target.value)}
                  placeholder="描述这个工具的用途和使用方法"
                  rows={4}
                  required
                  disabled={loading}
                />
              </div>
            </div>
            
            <div className="modal-actions">
              <button 
                className="modal-button cancel"
                onClick={() => setShowSaveModal(false)}
                disabled={loading}
              >
                取消
              </button>
              <button 
                className="modal-button save"
                onClick={handleModalSave}
                disabled={loading || !toolName.trim() || !toolDescription.trim()}
              >
                {loading ? '保存中...' : '保存'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Assistant;
