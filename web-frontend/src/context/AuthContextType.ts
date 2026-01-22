import { createContext } from 'react';
import type { User } from '../types';

// 定义 AuthContext 的类型
export interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  register: (email: string, password: string, phone?: string) => Promise<void>;
}

// 创建并导出 AuthContext
export const AuthContext = createContext<AuthContextType | undefined>(undefined);
