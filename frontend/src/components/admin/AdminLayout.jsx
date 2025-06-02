import React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import AdminSidebar from './AdminSidebar';
import './admin.css';

const AdminLayout = () => {
  const navigate = useNavigate();
  const handleLogout = () => {
    localStorage.removeItem('adminToken');
    navigate('/admin/login');
  };
  return (
    <div className="admin-root">
      <header className="admin-header">
        <div className="admin-logo">LAMPS <span role="img" aria-label="lamp">💡</span></div>
        <div className="admin-title">ADMIN PANEL</div>
        <button className="admin-logout-btn" onClick={handleLogout} style={{marginLeft: 'auto'}}>Выйти</button>
      </header>
      <div className="admin-main">
        <AdminSidebar />
        <div className="admin-content"><Outlet /></div>
      </div>
      <footer className="admin-footer">
        <div className="admin-footer-logo">LAMPS <span role="img" aria-label="lamp">💡</span></div>
        <div className="admin-footer-copy">© 2025</div>
      </footer>
    </div>
  );
};

export default AdminLayout; 