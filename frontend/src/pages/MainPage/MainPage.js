import React, {useEffect, useState} from "react";
import './MainPage.css'

const baseURL = "http://localhost:8080/api/spotify/history"

const MainPage = () => {
  const [data, setData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(baseURL, {
          method: "GET",
          //credentials: "include",
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const jsonData = await response.json();
        console.log(jsonData)

        setData(jsonData);
      } catch (error) {
        console.error("Fetch error:", error);
      }
    };

    fetchData();
  }, []);

  if (!data) return <div>Loading...</div>;

  return (
    <div className="main-page">
      <h1>{data.user_profile.display_name}さんの最近再生した曲</h1>
      <img src={data.user_profile.profile_image_url} alt="User Profile" width="100" />
      <p>Country: {data.user_profile.country}</p>
      <div className="track-list">
        {data.recently_played_tracks.map((track, index) => (
          <div key={index} className="track-item">
            <h2>{track.track_name}</h2>
            <img src={track.album.smallest_image_url} alt={track.album.name} width="50" />
            <p>Album: {track.album.name}</p>
            {track.artists?.map((artist, idx) => (
              <div key={idx}>
                <p>Artist: <a href={artist.spotify_url}>{artist.name}</a></p>
                <img src={artist.smallest_image_url} alt={artist.name} width="30" />
              </div>
            ))}
            <p>
              YouTube:{" "}
              <a
                href={`https://www.youtube.com/results?search_query=${encodeURIComponent(
                  track.youtube_search_query
                )}`}
                target="_blank"
                rel="noopener noreferrer"
              >
                {track.youtube_search_query}
              </a>
            </p>
          </div>
        ))}
      </div>
    </div>

  );
};

export default MainPage;