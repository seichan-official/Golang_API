import React from 'react';
import './TopPage.css'; 
import { useNavigate } from 'react-router-dom';

const TopPage = () => {
  const navigate = useNavigate(); // コメントアウトを解除

  const handleLogin = () => {
    navigate('/main'); // navigateを使用してページ遷移
  };
  
  return (
    <div className="top-page">
      <div className="top-page-content">
        <h1 className="main-title">
           <span className="highlight">Welcome To SpoTube!</span> 
        </h1>
          <button onClick={handleLogin}>
            Go To SpoTube
          </button>
      </div>
    </div>
  );
}

export default TopPage;
