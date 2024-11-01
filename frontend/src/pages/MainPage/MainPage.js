import React, {useEffect, useState} from "react";
import './MainPage.css'

const baseURL = "http://localhost:8080/api/spotify/history"
const youtubeURL = "http://localhost:8080/api/youtube/search"



const MainPage = () => {
  const [post, setData] = useState(null);

  useEffect(() => {

    const fetchData = async () => {
      try {
        const response = await fetch(baseURL, {
          method: "GET",
          //credentials: "include",
        });
        if (!response.ok) {
          throw new Error("GET request failed")
        };
        const spotifyData = await response.json();
        console.log("spotiy", spotifyData);
        console.log(spotifyData)
        for(const value of spotifyData.recently_played_tracks){
          console.log(value.youtube_search_query)
          const postData = {"query": value.youtube_search_query}

          const postResponse = await fetch(youtubeURL, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(postData),
          });

          if (!postResponse.ok) {
            value.youtubeURL =""
          }
          else{
            const youtubeData = await postResponse.json();
            value.youtubeURL = youtubeData[0].video_url       
          }

  
          console.log("value", value);
          console.log("youtubeURL", value.youtubeURL)
        }
        setData(spotifyData);
      } catch (error) {
        console.error("Fetch error:", error);
      }
    };

    fetchData();
  }, []);

  if (!post) return <div>Loading...</div>;

  const redirectYoutube = (url) => {
    window.location.href = url;
  }

  return (
    <div className="main-page">
      <h1>{post.user_profile.display_name}さんの最近再生した曲</h1>
      <img src={post.user_profile.profile_image_url} alt="User Profile" width="100" />
      <p>Country: {post.user_profile.country}</p>
      <div className="track-list">
        {post.recently_played_tracks.map((track, index) => (
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
              <button onClick={() => redirectYoutube(track.youtubeURL)}>YouubeLink</button>
            </p>
          </div>
        ))}
      </div>
    </div>

  );
};

export default MainPage;