import type { Conversation, Message } from '../types';

// API基础URL
const API_BASE_URL = '/api';

// 获取历史会话列表
export const getConversations = async (type?: string): Promise<Conversation[]> => {
  const url = new URL(`${API_BASE_URL}/conversations`);
  if (type) {
    url.searchParams.append('type', type);
  }
  
  const response = await fetch(url.toString());
  if (!response.ok) {
    throw new Error('Failed to fetch conversations');
  }
  
  return response.json();
};

// 获取会话详情
export const getConversationById = async (id: string): Promise<Conversation> => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch conversation');
  }
  
  return response.json();
};

// 创建新会话
export const createConversation = async (type: string, title: string): Promise<Conversation> => {
  const response = await fetch(`${API_BASE_URL}/conversations`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      type,
      title,
    }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create conversation');
  }
  
  return response.json();
};

// 更新会话消息
export const updateConversationMessages = async (id: string, messages: Message[]): Promise<Conversation> => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      messages,
    }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to update conversation');
  }
  
  return response.json();
};

// 删除会话
export const deleteConversation = async (id: string): Promise<void> => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`, {
    method: 'DELETE',
  });
  
  if (!response.ok) {
    throw new Error('Failed to delete conversation');
  }
};

// 调用AI服务进行聊天
export const chatWithAI = async (messages: Message[], conversationType: string, conversationId?: string): Promise<{
  content: string;
  conversation_id: string;
  messages: Message[];
}> => {
  const response = await fetch(`${API_BASE_URL}/ai/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      messages,
      conversation_type: conversationType,
      conversation_id: conversationId || '',
    }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to call AI service');
  }
  
  return response.json();
};

// 生成会话标题
export const generateConversationTitle = async (messages: Message[]): Promise<string> => {
  // 这里应该调用python-ai-service的接口生成标题
  // 暂时返回一个默认标题
  return messages[0]?.content.substring(0, 30) || '新会话';
};