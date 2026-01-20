import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar: React.FC = () => {
  const { logout, user } = useAuth();
  const location = useLocation();

  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <div className="navbar-logo">
          <Link to="/" className="logo-link">
            🧠 神经科学AI修行助手
          </Link>
        </div>
        
        <div className="navbar-links">
          <Link 
            to="/" 
            className={`navbar-link ${location.pathname === '/' ? 'active' : ''}`}
          >
            首页
          </Link>
          <Link 
            to="/analysis" 
            className={`navbar-link ${location.pathname === '/analysis' ? 'active' : ''}`}
          >
            神经科学分析
          </Link>
          <Link 
            to="/mapping" 
            className={`navbar-link ${location.pathname === '/mapping' ? 'active' : ''}`}
          >
            修行映射
          </Link>
          <Link 
            to="/assistant" 
            className={`navbar-link ${location.pathname === '/assistant' ? 'active' : ''}`}
          >
            修行小助手
          </Link>
          <Link 
            to="/tools" 
            className={`navbar-link ${location.pathname === '/tools' ? 'active' : ''}`}
          >
            我的工具
          </Link>
        </div>
        
        <div className="navbar-user">
          <span className="user-email">{user?.email}</span>
          <button 
            className="logout-button"
            onClick={handleLogout}
          >
            退出登录
          </button>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
