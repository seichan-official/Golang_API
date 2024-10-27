import React, { useEffect, useState } from "react";
import './MainPage.css';

const MainPage = () => {
  const [videos, setVideos] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/youtube")
      .then(response => response.json())
      .then(data => setVideos(data))
      .catch(error => console.error("Error fetching YouTube data:", error));
  }, []);

  return (
    <div className="main-page">
      <h1>Main Page</h1>
      <div className="video-list">
        {videos.map((video, index) => (
          <div key={index} className="video-item">
            <h2>{video.title}</h2>
            <p>{video.description}</p>
            <a href={video.url} target="_blank" rel="noopener noreferrer">
              Watch on YouTube
            </a>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MainPage;
