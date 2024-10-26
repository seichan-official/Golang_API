import React from 'react';
import {BrowserRouter, Routes, Route} from 'react-router-dom';
import TopPage from './pages/TopPage'; 
import './pages/TopPage.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import Header from './components/Header';
import Footer from './components/Footer';

const App = () => {
  return (
    <BrowserRouter>
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<TopPage />} />
        </Routes>
      </main>
      <Footer />
    </BrowserRouter>
  );
};

export default App;
