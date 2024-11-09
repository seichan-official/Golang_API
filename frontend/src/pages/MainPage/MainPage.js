// MainPage.js
import React, { useEffect, useState } from "react";
import './MainPage.css'

const baseURL = "http://localhost:8080/api/spotify/history";
const youtubeURL = "http://localhost:8080/api/youtube/search";

const MainPage = () => {
  const [post, setData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(baseURL, {
          method: "GET",
          credentials: "include",
        });
        if (!response.ok) {
          throw new Error("GET request failed");
        }
        const spotifyData = await response.json();
        let count = 0;
        
        for (const value of spotifyData.recently_played_tracks) {
          if (count >= 3) break;

          const postData = { "query": value.youtube_search_query };

          const postResponse = await fetch(youtubeURL, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(postData),
          });

          if (postResponse.ok) {
            const youtubeData = await postResponse.json();
            value.youtubeURL = youtubeData[0]?.video_url;
            value.youtubeThumbnail = youtubeData[0]?.thumbnail;
          } else {
            value.youtubeURL = "";
            value.youtubeThumbnail = "";
          }
          count++;
        }
        
        setData(spotifyData);
      } catch (error) {
        console.error("Fetch error:", error);
      }
    };

    fetchData();
  }, []);

  if (!post) return <div>Loading...</div>;

  return (
    <div className="main-page">
      <h1>{post.profile.display_name}さんの最近再生した曲</h1>
      <img src={post.profile.profile_image_url} alt="User Profile" width="100" />
      <p>Country: {post.profile.country}</p>
      <div className="track-list">
        {post.recently_played_tracks.slice(0, 3).map((track, index) => (
          <div key={index} className="track-item">
            <a href={track.youtubeURL} target="_blank" rel="noopener noreferrer">
              <img src={track.youtubeThumbnail || track.album.smallest_image_url} alt={track.track_name} className="thumbnail" />
            </a>
            <div className="track-info">
              <h2 className="track-title">{track.track_name}</h2>
              <p>Album: {track.album.name}</p>
              {track.artists?.map((artist, idx) => (
                <div key={idx}>
                  <p>Artist: <a href={artist.spotify_url}>{artist.name}</a></p>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MainPage;
