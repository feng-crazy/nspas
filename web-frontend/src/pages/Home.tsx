import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="home-page">
      <div className="home-header">
        <h1>🧠 神经科学AI修行助手</h1>
        <p>探索大脑奥秘，提升认知能力</p>
      </div>
      
      <div className="home-cards">
        <Link to="/analysis" className="home-card">
          <div className="card-content">
            <h2>🧠 神经科学分析</h2>
            <p>分析你的思维过程、认知模式和行为倾向，从神经科学角度解释大脑运作机制。</p>
          </div>
        </Link>
        
        <Link to="/mapping" className="home-card">
          <div className="card-content">
            <h2>✨ 修行映射</h2>
            <p>将修行语录映射到脑科学机制与神经通路，提供科学解释和神经科学依据。</p>
          </div>
        </Link>
        
        <Link to="/assistant" className="home-card">
          <div className="card-content">
            <h2>🔧 修行小助手</h2>
            <p>根据你的需求，个性化生成脑科学修行工具，帮助你提升认知能力和情绪管理。</p>
          </div>
        </Link>
      </div>
    </div>
  );
};

export default Home;
