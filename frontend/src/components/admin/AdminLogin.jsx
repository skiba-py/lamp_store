import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminApi } from '../../services/api';

const AdminLogin = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const { token } = await adminApi.login({ username, password });
      localStorage.setItem('adminToken', token);
      navigate('/admin/products');
    } catch (err) {
      setError('Неверный логин или пароль');
    }
  };

  return (
    <div className="admin-login-wrapper">
      <form className="admin-login-form" onSubmit={handleSubmit}>
        <h2>Вход для администратора</h2>
        <input type="text" placeholder="Логин" value={username} onChange={e => setUsername(e.target.value)} required />
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} required />
        <button type="submit">Войти</button>
        {error && <div className="admin-login-error">{error}</div>}
      </form>
    </div>
  );
};

export default AdminLogin; 