import React from 'react';
import './TopPage.css'; 

const TopPage = () => {

  const redirectLogin = () => {
    window.location.href = "http://localhost:8080/api/spotify/login"
  }
  
  return (
    <div className="top-page">
      <div className="top-page-content">
        <h1 className="main-title">
           <span className="highlight">Welcome To SpoTube!</span> 
        </h1>
          <button onClick={redirectLogin}>
            Go To SpoTube
          </button>
      </div>
    </div>
  );
}

export default TopPage;
