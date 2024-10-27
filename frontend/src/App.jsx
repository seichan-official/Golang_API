import React from 'react';
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import TopPage from './pages/TopPage/TopPage'; 
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import MainPage from './pages/MainPage/MainPage';

const AppContent = () => {
  const location = useLocation();
  
  // 現在のパスが"/main"かどうかを判定
  const isMainPage = location.pathname === '/main';

  return (
    <>
      {/* /mainページ以外でのみヘッダーとフッターを表示 */}
      {!isMainPage && <Header />}
      <main>
        <Routes>
          <Route path="/" element={<TopPage />} />
          <Route path="/main" element={<MainPage />} />
        </Routes>
      </main>
      {!isMainPage && <Footer />}
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
