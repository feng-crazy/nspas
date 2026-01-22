import React, { useState, useEffect } from 'react';
import type { Tool } from '../types';
import { getUserTools, deleteTool } from '../services/api';
import './ToolsPage.css';
import '../components/Card.css';
import '../components/Modal.css';

const Tools: React.FC = () => {
  const [tools, setTools] = useState<Tool[]>([]);
  const [showToolModal, setShowToolModal] = useState(false);
  const [currentTool, setCurrentTool] = useState<Tool | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 获取用户工具
  const fetchTools = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getUserTools();
      // 转换日期字符串为Date对象
      const formattedTools = data.map(tool => ({
        ...tool,
        createdAt: new Date(tool.createdAt)
      }));
      setTools(formattedTools);
    } catch (err) {
      console.error('Failed to fetch tools:', err);
      setError('获取工具失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载工具
  useEffect(() => {
    fetchTools();
  }, []);

  const handleOpenTool = (tool: Tool) => {
    setCurrentTool(tool);
    setShowToolModal(true);
  };

  const handleDeleteTool = async (toolId: string) => {
    try {
      await deleteTool(toolId);
      // 更新工具列表
      setTools(prevTools => prevTools.filter(tool => tool.id !== toolId));
    } catch (err) {
      console.error('Failed to delete tool:', err);
      setError('删除工具失败，请稍后重试');
    }
  };

  return (
    <div className="tools-page">
      <div className="tools-header">
        <h1>🔧 我的修行工具</h1>
      </div>
      
      {error && <div className="tools-error">{error}</div>}
      
      <div className="tools-list">
        {loading ? (
          <div className="tools-loading">
            <p>正在加载工具...</p>
          </div>
        ) : tools.length === 0 ? (
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
