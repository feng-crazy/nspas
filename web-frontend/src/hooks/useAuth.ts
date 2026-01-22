import { useContext } from 'react';
import { AuthContext } from '../context/AuthContextType';

/**
 * 自定义 Hook，用于获取 AuthContext 中的认证状态和方法
 * @returns AuthContext 中的认证状态和方法
 */
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
