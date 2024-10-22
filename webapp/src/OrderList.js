import React, { useEffect, useState } from 'react';
import './OrderList.css'; // Подключение стилей

const OrderList = () => {
  const [orders, setOrders] = useState([]); // Здесь сохраняются заказы после получения данных
  const [loading, setLoading] = useState(true); // Состояние для показа "загрузка"
  const [error, setError] = useState(null); // Состояние для ошибок

  useEffect(() => {
    // Выполняем GET запрос к эндпоинту /orders для получения списка заказов
    fetch('ХОСТ ДЛЯ БЭКА', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => {
        if (!response.ok) { // Проверка успешного ответа
          throw new Error('Network response was not ok');
        }
        return response.json(); // Преобразуем ответ в JSON
      })
      .then(data => {
        setOrders(data); // Сохраняем заказы в state (orders)
        setLoading(false); // Загрузка завершена
      })
      .catch(error => {
        console.error('Error fetching orders:', error); // Обрабатываем ошибку
        setError(error); // Сохраняем ошибку в state
        setLoading(false); // Загрузка завершена (даже при ошибке)
      });
  }, []); // Пустой массив зависимостей, чтобы запрос выполнялся один раз при монтировании компонента

  // Если данные еще загружаются, показываем "Loading..."
  if (loading) {
    return <div>Loading...</div>;
  }

  // Если возникла ошибка, показываем сообщение об ошибке
  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="order-list">
      <h2>Список заказов</h2>
      <div className="order-cards">
        {orders.length > 0 ? (
          orders.map((order) => (
            <div key={order.id} className="order-card">
              <h3>Предмет: {order.subject}</h3>
              <p>Описание: {order.description}</p>
              <div className="order-footer">
                <span className="price">Мин. цена: {order.min_price}</span>
                <span className="price">Макс. цена: {order.max_price}</span>
                <span className="created-at">Создан: {new Date(order.created_at).toLocaleString()}</span>
              </div>
            </div>
          ))
        ) : (
          <p>Нет доступных заказов.</p>
        )}
      </div>
    </div>
  );
};

export default OrderList;
