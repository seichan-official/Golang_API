import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css'; // CSSファイルのインポートを修正

const Header = () => {
  return (
    <header>
      <div className="bg-dark text-white py-2">
        <div className="container">
          <Link to="/" className="main-title">
            <h1>
              SpoTube
            </h1>
          </Link>
        </div>
      </div>
    </header>
  );
};

export default Header;
