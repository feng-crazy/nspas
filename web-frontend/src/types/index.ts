// 对话类型
export type ConversationType = 'analysis' | 'mapping' | 'assistant';

// 消息类型
export interface Message {
  id: string;
  content: string;
  isUser: boolean;
  createdAt: Date;
}

// 对话类型
export interface Conversation {
  id: string;
  userId: string;
  type: ConversationType;
  title: string;
  messages: Message[];
  createdAt: Date;
  updatedAt: Date;
}

// 工具类型
export interface Tool {
  id: string;
  userId: string;
  name: string;
  description: string;
  htmlContent: string;
  conversationId: string;
  createdAt: Date;
}

// 用户类型
export interface User {
  id: string;
  email: string;
  phone?: string;
  role: 'user' | 'admin';
  createdAt: Date;
}
