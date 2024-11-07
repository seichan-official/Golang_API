import React from 'react';
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import TopPage from './pages/TopPage/TopPage'; 
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import MainPage_Header from './components/Header/MainPage_Header';
import MainPage from './pages/MainPage/MainPage';
import SignupPage from './pages/Singup/SingupPage'; 

const AppContent = () => {
  const location = useLocation();

  // 現在のURLパスに基づいてヘッダーを切り替え
  const renderHeader = () => {
    if (location.pathname === '/') {
      return <Header />;
    } else if (location.pathname === '/main') {
      return <MainPage_Header />;
    } else {
      return <Header />;
    }
  };

  return (
    <>
      {renderHeader()} {/* 条件に基づいたヘッダーをレンダリング */}
      <main>
        <Routes>
          <Route path="/" element={<TopPage />} />
          <Route path="/main" element={<MainPage />} />
          <Route path="/signup" element={<SignupPage />} />
        </Routes>
      </main>
      <Footer />
    </>
  );
};

const App = () => {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
};

export default App;
