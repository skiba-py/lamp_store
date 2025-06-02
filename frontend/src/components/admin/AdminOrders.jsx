import React, { useEffect, useState } from 'react';
import { adminApi, productsApi } from '../../services/api';

const AdminOrders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const token = localStorage.getItem('adminToken');
  const [productNames, setProductNames] = useState({});

  const fetchOrders = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await adminApi.getOrders(token);
      console.log('full response:', response);
      setOrders(response);
    } catch (err) {
      setError('Ошибка загрузки заказов');
    } finally {
      setLoading(false);
    }
  };

  const fetchProductNames = async (orders) => {
    const names = {};
    for (const order of orders) {
      if (Array.isArray(order.items)) {
        for (const item of order.items) {
          if (item.product_id && !names[item.product_id]) {
            try {
              const product = await productsApi.getProductById(item.product_id);
              names[item.product_id] = product.name;
            } catch (e) {
              names[item.product_id] = '';
            }
          }
        }
      }
    }
    setProductNames(names);
  };

  useEffect(() => {
    fetchOrders();
    // eslint-disable-next-line
  }, []);

  useEffect(() => {
    if (orders && orders.length > 0) {
      fetchProductNames(orders);
    }
  }, [orders]);

  useEffect(() => {
    console.log('orders:', orders);
  }, [orders]);

  const handleStatusChange = async (id, status) => {
    try {
      const order = orders.find(o => o.id === id);
      if (!order) return;
      const updatedOrder = { ...order, status };
      console.log('updateOrder payload:', updatedOrder);
      await adminApi.updateOrder(id, updatedOrder, token);
      fetchOrders();
    } catch (e) {
      setError('Ошибка обновления статуса заказа');
      console.error(e);
    }
  };

  return (
    <div className="admin-orders">
      <h2>Заказы</h2>
      {error && <div className="admin-error">{error}</div>}
      <div className="admin-orders-table-wrapper">
        <table className="admin-orders-table">
          <thead>
            <tr>
              <th>id заказа</th>
              <th>Пользователь</th>
              <th>Статус</th>
              <th>Сумма</th>
              <th>Товары</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {Array.isArray(orders) && orders.map(order => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>{order.user_id}</td>
                <td>{order.status}</td>
                <td>{order.total}</td>
                <td>
                  <ul>
                    {Array.isArray(order.items) && order.items.map(item => (
                      <li key={item.id}>{productNames[item.product_id] ? `${productNames[item.product_id]} (${item.product_id})` : item.product_id} × {item.quantity}</li>
                    ))}
                  </ul>
                </td>
                <td>
                  <select value={order.status} onChange={e => handleStatusChange(order.id, e.target.value)}>
                    <option value="pending">pending</option>
                    <option value="completed">completed</option>
                    <option value="cancelled">cancelled</option>
                  </select>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AdminOrders; 