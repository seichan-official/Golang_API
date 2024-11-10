import React, { useEffect, useState } from "react";
import './MainPage.css';

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

        // YouTubeデータの取得
        for (const track of spotifyData.recently_played_tracks) {
          const postData = { "query": track.youtube_search_query };
          const postResponse = await fetch(youtubeURL, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(postData),
          });

          if (postResponse.ok) {
            const youtubeData = await postResponse.json();
            track.youtubeVideo = {
              title: youtubeData[0].title,
              url: youtubeData[0].video_url,
              thumbnail: youtubeData[0].thumbnail_url,
            };
          } else {
            track.youtubeVideo = null;
          }
        }
        setData(spotifyData);
      } catch (error) {
        console.error("Fetch error:", error);
      }
    };

    fetchData();
  }, []);

  if (!post) return <div>Loading...</div>;

  const redirectToYoutube = (url) => {
    window.open(url, "_blank");
  };

  return (
    <div className="main-page">
      <h1>{post.profile.display_name}さんの最近再生した曲</h1>
      <img src={post.profile.profile_image_url} alt="User Profile" className="profile-image" />
      <p>Country: {post.profile.country}</p>
      <div className="track-list">
        {post.recently_played_tracks.map((track, index) => (
          <div key={index} className="track-item">
            <h2>{track.track_name}</h2>
            <img src={track.album.smallest_image_url} alt={track.album.name} className="album-image" />
            <p>Album: {track.album.name}</p>
            
            {track.artists?.map((artist, idx) => (
              <div key={idx} className="artist-info">
                <p>Artist: <a href={artist.spotify_url} target="_blank" rel="noopener noreferrer">{artist.name}</a></p>
                <img src={artist.smallest_image_url} alt={artist.name} className="artist-image" />
              </div>
            ))}
            
            {track.youtubeVideo && (
              <div className="youtube-info">
                <img 
                  src={track.youtubeVideo.thumbnail} 
                  alt={track.youtubeVideo.title} 
                  className="youtube-thumbnail" 
                  onClick={() => redirectToYoutube(track.youtubeVideo.url)}
                />
                <p className="youtube-title" onClick={() => redirectToYoutube(track.youtubeVideo.url)}>
                  {track.youtubeVideo.title}
                </p>
                <button onClick={() => redirectToYoutube(track.youtubeVideo.url)}>Watch on YouTube</button>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default MainPage;
