import React, { useEffect, useRef, useState } from 'react';
import './TopPage.css'; 
import sflogo from './imges/Spotify LOGO.png';
import ytlogo from './imges/youtube_logo.png';
import { useNavigate } from 'react-router-dom';

const TopPage = () => {
  const navigate = useNavigate();
  const [isDescriptionVisible, setIsDescriptionVisible] = useState(false);
  const descriptionRef = useRef(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsDescriptionVisible(true);
          observer.unobserve(entry.target);
        }
      },
      { threshold: 0.75 }
    );

    const currentRef = descriptionRef.current; // refをローカル変数にコピー

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef); // クリーンアップでcurrentRefを使用
      }
    };
  }, []); // 依存配列を空にすることで、一度だけ実行

  const handleLogin = () => {
    navigate('/main');
  };
  
  return (
    <div className="top-page">
      {/* トップページのセクション */}
      <section className="top-page-content">
        <h1 className="main-title">
           <span className="highlight">Welcome To SpoTube!</span> 
        </h1>
        <button className="top-page-button" onClick={handleLogin}>
          Go To SpoTube
        </button>
      </section>

      {/* 説明ページのセクション */}
      <section
        ref={descriptionRef}
        className={`description-section ${isDescriptionVisible ? 'fade-in' : ''}`}
      >
        <div className="logo-container">
          <img src={sflogo} className="spotify_logo" alt="Spotify logo" />
          <img src={ytlogo} className="youtube_logo" alt="YouTube logo" />
        </div>
        <div className="description-text">
          <h2 className='subtitle'>What is SpoTube?</h2>
          <p className='explanation'>SpoTubeは、Spotifyで聞いた曲の履歴を元にMVを表示させるサービスです。</p>
        </div>
      </section>
    </div>
  );
};

export default TopPage;
