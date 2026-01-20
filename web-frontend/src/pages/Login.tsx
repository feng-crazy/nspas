import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Login: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const { login, register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      if (isLogin) {
        await login(email, password);
      } else {
        await register(email, password);
      }
      navigate('/');
    } catch (err) {
      setError('操作失败，请检查您的输入或网络连接');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  // 处理微信登录
  const handleWeChatLogin = () => {
    // 生成state参数，用于防止CSRF攻击
    const state = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
    
    // 保存state到localStorage，用于后续验证
    localStorage.setItem('wechat_state', state);
    
    // 调用后端API获取微信授权URL
    fetch('http://localhost:8080/api/auth/wechat?state=' + state)
      .then(response => response.json())
      .then(data => {
        // 重定向到微信授权页面
        window.location.href = data.url;
      })
      .catch(error => {
        console.error('Failed to get WeChat auth URL:', error);
        setError('获取微信授权URL失败');
      });
  };

  // 检查是否是微信回调
  React.useEffect(() => {
    // 获取URL参数
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const state = urlParams.get('state');
    
    // 如果有code参数，说明是微信回调
    if (code) {
      // 验证state
      const savedState = localStorage.getItem('wechat_state');
      if (state !== savedState) {
        setError('Invalid state parameter');
        return;
      }
      
      // 清除保存的state
      localStorage.removeItem('wechat_state');
      
      // 处理微信登录回调
      handleWeChatCallback(code, state);
    }
  }, []);

  // 处理微信登录回调
  const handleWeChatCallback = (code: string, state: string) => {
    setIsLoading(true);
    
    // 调用后端API处理微信登录
    fetch('http://localhost:8080/api/auth/wechat/callback?code=' + code + '&state=' + state)
      .then(response => response.json())
      .then(data => {
        // 保存token和用户信息
        // 这里需要根据后端返回的数据结构进行调整
        if (data.token) {
          // 登录成功，跳转到首页
          navigate('/');
        } else {
          setError('微信登录失败');
        }
      })
      .catch(error => {
        console.error('Failed to login with WeChat:', error);
        setError('微信登录失败');
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-header">
          <h1>🧠 神经科学AI修行助手</h1>
          <h2>{isLogin ? '登录' : '注册'}</h2>
        </div>
        
        {error && <div className="login-error">{error}</div>}
        
        <form onSubmit={handleSubmit} className="login-form">
          <div className="form-group">
            <label htmlFor="email">邮箱：</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="请输入您的邮箱"
              required
              disabled={isLoading}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">密码：</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="请输入您的密码"
              required
              disabled={isLoading}
              minLength={6}
            />
          </div>
          
          <button 
            type="submit" 
            className="login-button"
            disabled={isLoading}
          >
            {isLoading ? '处理中...' : isLogin ? '登录' : '注册'}
          </button>
        </form>
        
        {/* 微信登录按钮 */}
        <div className="login-divider">
          <span>或</span>
        </div>
        
        <button 
          className="wechat-login-button"
          onClick={handleWeChatLogin}
          disabled={isLoading}
        >
          <span className="wechat-icon">💬</span>
          使用微信登录
        </button>
        
        <div className="login-toggle">
          <p>
            {isLogin ? '还没有账号？' : '已有账号？'}
            <button 
              className="toggle-button"
              onClick={() => setIsLogin(!isLogin)}
              disabled={isLoading}
            >
              {isLogin ? '立即注册' : '立即登录'}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
