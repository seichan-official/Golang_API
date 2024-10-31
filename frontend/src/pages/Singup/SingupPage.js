import React, { useState } from 'react';
import './SingupPage.css'; // CSSファイルをインポート

const SignupPage = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  return (
    <div className="signup-page">
      <h2>サインアップは現在開発中となります。 </h2>
      <h2>
          開発完了するまではSpotifyアカウントでお楽しみください。
      </h2>
    </div>
  );
};

export default SignupPage;

