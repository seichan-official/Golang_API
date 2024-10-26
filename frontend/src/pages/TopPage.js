import React from 'react';
import './TopPage.css'; 
import {newUserSessionPath}  
import App from '../App';
const TopPage = () => {

  return (
    <div className="top-page">
      <div className="top-page-content">
        <h1 className="main-title">
           <span className="highlight">Welcome To SpoTube!</span> 
        </h1>
          <button className="consultation-btn">
            <a href={newUserSessionPath} 
             className="btn btn-info btn-sm btn-block mb-3 sign_in">
              Go to SpoTube World
            </a>
          </button>
      </div>
    </div>
  );
}

export default TopPage;
