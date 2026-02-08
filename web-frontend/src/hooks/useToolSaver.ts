import { useState } from 'react';
import { saveTool } from '../services/api';
import type { ToolSaveData } from '../types';

/**
 * 工具保存 Hook
 * 提供工具保存的状态管理和逻辑处理
 */
export const useToolSaver = () => {
  const [showSaveModal, setShowSaveModal] = useState(false);
  const [toolName, setToolName] = useState('');
  const [toolDescription, setToolDescription] = useState('');
  const [currentHtmlContent, setCurrentHtmlContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null);
  
  /**
   * 打开保存工具模态框
   * @param htmlContent HTML内容
   * @param conversationId 会话ID
   */
  const openSaveModal = (htmlContent: string, conversationId?: string) => {
    setCurrentHtmlContent(htmlContent);
    if (conversationId) {
      setCurrentConversationId(conversationId);
    }
    setShowSaveModal(true);
    setError(null);
  };
  
  /**
   * 保存工具
   */
  const handleSaveTool = async () => {
    if (!currentConversationId) {
      setError('缺少会话ID，无法保存工具');
      return false;
    }
    
    if (!toolName.trim() || !toolDescription.trim()) {
      setError('工具名称和描述不能为空');
      return false;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      const saveData: ToolSaveData = {
        name: toolName,
        description: toolDescription,
        html_content: currentHtmlContent,
        conversation_id: currentConversationId
      };
      
      await saveTool(saveData);
      
      // 关闭模态框并重置表单
      resetForm();
      return true;
    } catch (err) {
      console.error('Failed to save tool:', err);
      setError('保存工具失败，请稍后重试');
      return false;
    } finally {
      setLoading(false);
    }
  };
  
  /**
   * 重置表单
   */
  const resetForm = () => {
    setShowSaveModal(false);
    setToolName('');
    setToolDescription('');
    setCurrentHtmlContent('');
    setCurrentConversationId(null);
    setError(null);
  };
  
  /**
   * 关闭模态框
   */
  const closeModal = () => {
    setShowSaveModal(false);
  };
  
  return {
    showSaveModal,
    toolName,
    setToolName,
    toolDescription,
    setToolDescription,
    loading,
    error,
    openSaveModal,
    handleSaveTool,
    resetForm,
    closeModal
  };
};