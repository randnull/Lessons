import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
//import './StudentRegistrationForm.css';

const StudentRegistrationForm = () => {
  const [formData, setFormData] = useState({
    grade: '',
    subject: '',
    goal: '',
    price: ''
  });

  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    // Передача данных на страницу профиля через navigate с использованием state
    navigate('/profile', { state: formData });
  };

  return (
    <div className="registration-form">
      <h2>Регистрация Ученика</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="grade">Класс/Курс</label>
          <input
            type="text"
            id="grade"
            name="grade"
            value={formData.grade}
            onChange={handleChange}
            required
          />
        </div>

        <div>
          <label htmlFor="subject">Предмет</label>
          <select
            id="subject"
            name="subject"
            value={formData.subject}
            onChange={handleChange}
            required
          >
            <option value="">Выберите предмет</option>
            <option value="math">Математика</option>
            <option value="science">Наука</option>
            <option value="history">История</option>
            <option value="english">Английский</option>
          </select>
        </div>

        <div>
          <label htmlFor="goal">Цель</label>
          <textarea
            id="goal"
            name="goal"
            value={formData.goal}
            onChange={handleChange}
            required
          />
        </div>

        <div>
          <label htmlFor="price">Цена занятия</label>
          <input
            type="number"
            id="price"
            name="price"
            value={formData.price}
            onChange={handleChange}
            required
          />
        </div>

        <button type="submit">Зарегистрироваться</button>
      </form>
    </div>
  );
};

export default StudentRegistrationForm;
