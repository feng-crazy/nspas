import type { Conversation, Message, User, Tool } from '../types';

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

// 调用AI服务进行聊天（流式）
export const chatWithAI = async (messages: Message[], conversationType: string, conversationId?: string, onUpdate?: (data: {
  content: string;
  full_content: string;
  conversation_id: string;
  messages: Message[];
}) => void): Promise<{
  content: string;
  conversation_id: string;
  messages: Message[];
}> => {
  return new Promise((resolve, reject) => {
    // 创建AbortController用于取消请求
    const controller = new AbortController();
    const { signal } = controller;
    
    // 发送POST请求，使用流式响应
    fetch(`${API_BASE_URL}/ai/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        messages,
        conversation_type: conversationType,
        conversation_id: conversationId || '',
      }),
      signal,
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to call AI service');
      }
      
      // 检查是否支持流式响应
      if (!response.body) {
        throw new Error('Response body is not readable');
      }
      
      // 创建ReadableStream读取器
      const reader = response.body.getReader();
      const decoder = new TextDecoder('utf-8');
      let buffer = '';
      let fullResponse: {
        content: string;
        conversation_id: string;
        messages: Message[];
      } | null = null;
      
      // 定义读取函数
      const readStream = async () => {
        try {
          const { done, value } = await reader.read();
          
          if (done) {
            // 流结束，返回完整响应
            if (fullResponse) {
              resolve(fullResponse);
            } else {
              reject(new Error('No complete response received'));
            }
            return;
          }
          
          // 解码新接收到的数据
          buffer += decoder.decode(value, { stream: true });
          
          // 处理缓冲区中的SSE事件
          let eventEndIndex;
          while ((eventEndIndex = buffer.indexOf('\n\n')) !== -1) {
            const eventData = buffer.substring(0, eventEndIndex);
            buffer = buffer.substring(eventEndIndex + 2);
            
            // 解析SSE事件
            if (eventData.startsWith('data: ')) {
              const jsonData = eventData.substring(6);
              try {
                const data = JSON.parse(jsonData);
                
                // 根据事件类型处理
                if (data) {
                  // 调用更新回调
                  if (onUpdate) {
                    onUpdate({
                      content: data.content,
                      full_content: data.full_content,
                      conversation_id: data.conversation_id,
                      messages: data.messages,
                    });
                  }
                  
                  // 保存完整响应
                  if (data.completed) {
                    fullResponse = {
                      content: data.content,
                      conversation_id: data.conversation_id,
                      messages: data.messages,
                    };
                  }
                }
              } catch (e) {
                console.error('Failed to parse SSE data:', e);
              }
            }
          }
          
          // 继续读取
          readStream();
        } catch (e) {
          reject(e);
        }
      };
      
      // 开始读取流
      readStream();
    })
    .catch(error => {
      reject(error);
    });
  });
};

// 生成会话标题
export const generateConversationTitle = async (messages: Message[]): Promise<string> => {
	// 这里应该调用python-ai-service的接口生成标题
	// 暂时返回一个默认标题
	return messages[0]?.content.substring(0, 30) || '新会话';
};

// 认证相关API
// 注册新用户
export const register = async (email: string, password: string, phone?: string): Promise<{ user: User; token: string }> => {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password, phone }),
  });
  if (!response.ok) {
    throw new Error('Registration failed');
  }
  return response.json();
};

// 用户登录
export const login = async (email: string, password: string): Promise<{ user: User; token: string }> => {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });
  if (!response.ok) {
    throw new Error('Login failed');
  }
  return response.json();
};

// 获取当前用户信息
export const getCurrentUser = async (): Promise<User> => {
  const response = await fetch(`${API_BASE_URL}/user`);
  if (!response.ok) {
    throw new Error('Failed to get current user');
  }
  return response.json();
};

// 工具相关API
// 获取用户所有工具
export const getUserTools = async (): Promise<Tool[]> => {
  const response = await fetch(`${API_BASE_URL}/tools`);
  if (!response.ok) {
    throw new Error('Failed to fetch tools');
  }
  return response.json();
};

// 获取工具详情
export const getToolById = async (id: string): Promise<Tool> => {
  const response = await fetch(`${API_BASE_URL}/tools/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch tool');
  }
  return response.json();
};

// 删除工具
export const deleteTool = async (id: string): Promise<void> => {
  const response = await fetch(`${API_BASE_URL}/tools/${id}`, {
    method: 'DELETE'
  });
  if (!response.ok) {
    throw new Error('Failed to delete tool');
  }
};

// 保存工具
export const saveTool = async (tool: {
  name: string;
  description: string;
  html_content: string;
  conversation_id: string;
}): Promise<Tool> => {
  const response = await fetch(`${API_BASE_URL}/tools`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(tool)
  });
  if (!response.ok) {
    throw new Error('Failed to save tool');
  }
  return response.json();
};