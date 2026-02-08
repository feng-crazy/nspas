import React from 'react';
import ConversationLayout from '../components/ConversationLayout';
import { useToolSaver } from '../hooks/useToolSaver';
import '../components/Modal.css';

const Assistant: React.FC = () => {
  const {
    showSaveModal,
    toolName,
    setToolName,
    toolDescription,
    setToolDescription,
    loading,
    error,
    openSaveModal,
    handleSaveTool,
    closeModal
  } = useToolSaver();

  // 处理保存工具的回调
  const handleSaveToolCallback = (htmlContent: string, conversationId?: string) => {
    openSaveModal(htmlContent, conversationId);
  };

  return (
    <div className="assistant-page">
      <ConversationLayout 
        conversationType="assistant" 
        onSaveTool={handleSaveToolCallback} 
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
                onClick={closeModal}
                disabled={loading}
              >
                取消
              </button>
              <button 
                className="modal-button save"
                onClick={handleSaveTool}
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
