import React from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import StudentRegistrationForm from './StudentRegistrationForm';
import OrderList from './OrderList';
import Profile from './Profile';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<OrderList />} />
        <Route path="/order" element={<StudentRegistrationForm />} />
      </Routes>
    </Router>
  );
}

export default App;
