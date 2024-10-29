import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import TopPage from './pages/TopPage/TopPage'; 
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import MainPage from './pages/MainPage/MainPage';

const AppContent = () => {
  return (
    <>
      {/* 全ページでヘッダーとフッターを表示 */}
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<TopPage />} />
          <Route path="/main" element={<MainPage />} />
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
