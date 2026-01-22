// 添加Authorization头到所有API请求
export const setupAuthInterceptor = () => {
  const originalFetch = window.fetch;
  
  window.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
    const token = localStorage.getItem('token');
    const headers = {
      ...init?.headers,
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    };
    
    return originalFetch(input, { ...init, headers });
  };
};

// 初始化拦截器
setupAuthInterceptor();
