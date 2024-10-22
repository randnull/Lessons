import React from 'react';
import { useLocation } from 'react-router-dom';
import './Profile.css';

const Profile = () => {
  const location = useLocation();
  const { grade, subject, goal, price } = location.state || {}; // Данные из формы

  return (
    <div className="profile">
      <h2>Личный кабинет</h2>
      <div className="profile-info">
        {grade && subject && goal && price ? (
          <>
            <p><strong>Класс/Курс:</strong> {grade}</p>
            <p><strong>Предмет:</strong> {subject}</p>
            <p><strong>Цель:</strong> {goal}</p>
            <p><strong>Цена занятия:</strong> {price} ₽</p>
          </>
        ) : (
          <p>Данные отсутствуют. Вернитесь к регистрации.</p>
        )}
      </div>
    </div>
  );
};

export default Profile;
