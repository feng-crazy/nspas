import React, { useState } from 'react';
import ConversationLayout from '../components/ConversationLayout';

const Assistant: React.FC = () => {
  const [showSaveModal, setShowSaveModal] = useState(false);
  const [toolName, setToolName] = useState('');
  const [toolDescription, setToolDescription] = useState('');
  const [currentHtmlContent, setCurrentHtmlContent] = useState('');

  const handleSaveTool = (htmlContent: string) => {
    setCurrentHtmlContent(htmlContent);
    setShowSaveModal(true);
  };

  const handleModalSave = () => {
    // 这里应该调用API保存工具
    console.log('Saving tool:', {
      name: toolName,
      description: toolDescription,
      htmlContent: currentHtmlContent
    });
    
    // 关闭模态框并重置表单
    setShowSaveModal(false);
    setToolName('');
    setToolDescription('');
    setCurrentHtmlContent('');
  };

  return (
    <div className="assistant-page">
      <ConversationLayout conversationType="assistant" onSaveTool={handleSaveTool} />
      
      {/* 保存工具模态框 */}
      {showSaveModal && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>保存修行工具</h3>
            
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
                />
              </div>
            </div>
            
            <div className="modal-actions">
              <button 
                className="modal-button cancel"
                onClick={() => setShowSaveModal(false)}
              >
                取消
              </button>
              <button 
                className="modal-button save"
                onClick={handleModalSave}
                disabled={!toolName.trim() || !toolDescription.trim()}
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Assistant;
