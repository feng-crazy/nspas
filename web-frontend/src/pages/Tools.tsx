import React, { useState } from 'react';
import type { Tool } from '../types';

// 模拟工具数据
const mockTools: Tool[] = [
  {
    id: '1',
    userId: '1',
    name: '注意力N-back训练工具',
    description: '通过N-back任务训练工作记忆和注意力，提升前额叶功能。',
    htmlContent: '<div>N-back tool HTML content</div>',
    conversationId: 'conv-1',
    createdAt: new Date('2025-12-27T14:30:00')
  },
  {
    id: '2',
    userId: '1',
    name: '正念呼吸引导工具',
    description: '引导式呼吸练习，提升专注力和情绪调节能力。',
    htmlContent: '<div>Breathing tool HTML content</div>',
    conversationId: 'conv-2',
    createdAt: new Date('2025-12-26T10:20:00')
  },
  {
    id: '3',
    userId: '1',
    name: '情绪调节训练器',
    description: '认知重构练习，改善情绪反应和思维模式。',
    htmlContent: '<div>Emotion regulation tool HTML content</div>',
    conversationId: 'conv-3',
    createdAt: new Date('2025-12-25T16:15:00')
  }
];

const Tools: React.FC = () => {
  const [tools, setTools] = useState<Tool[]>(mockTools);
  const [showToolModal, setShowToolModal] = useState(false);
  const [currentTool, setCurrentTool] = useState<Tool | null>(null);

  const handleOpenTool = (tool: Tool) => {
    setCurrentTool(tool);
    setShowToolModal(true);
  };

  const handleDeleteTool = (toolId: string) => {
    // 这里应该调用API删除工具
    setTools(prevTools => prevTools.filter(tool => tool.id !== toolId));
  };

  return (
    <div className="tools-page">
      <div className="tools-header">
        <h1>🔧 我的修行工具</h1>
      </div>
      
      <div className="tools-list">
        {tools.length === 0 ? (
          <div className="no-tools">
            <p>您还没有保存任何修行工具。</p>
            <p>在修行小助手中创建并保存工具后，它们会显示在这里。</p>
          </div>
        ) : (
          tools.map(tool => (
            <div key={tool.id} className="tool-card">
              <div className="tool-card-content">
                <h3>{tool.name}</h3>
                <p className="tool-description">{tool.description}</p>
                <p className="tool-date">
                  创建于：{tool.createdAt.toLocaleDateString()} {tool.createdAt.toLocaleTimeString()}
                </p>
              </div>
              
              <div className="tool-card-actions">
                <button 
                  className="tool-button open"
                  onClick={() => handleOpenTool(tool)}
                >
                  打开
                </button>
                <button 
                  className="tool-button delete"
                  onClick={() => handleDeleteTool(tool.id)}
                >
                  删除
                </button>
              </div>
            </div>
          ))
        )}
      </div>
      
      {/* 工具预览模态框 */}
      {showToolModal && currentTool && (
        <div className="modal-overlay">
          <div className="modal-content tool-modal">
            <div className="modal-header">
              <h3>{currentTool.name}</h3>
              <button 
                className="modal-close"
                onClick={() => setShowToolModal(false)}
              >
                ×
              </button>
            </div>
            
            <div className="tool-preview">
              <h4>工具描述</h4>
              <p>{currentTool.description}</p>
              
              <h4>工具内容</h4>
              <div 
                className="tool-html"
                dangerouslySetInnerHTML={{ __html: currentTool.htmlContent }}
              />
            </div>
            
            <div className="modal-actions">
              <button 
                className="modal-button close"
                onClick={() => setShowToolModal(false)}
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Tools;
