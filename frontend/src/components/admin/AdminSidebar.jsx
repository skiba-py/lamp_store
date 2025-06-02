import React from 'react';
import { NavLink } from 'react-router-dom';

const AdminSidebar = () => (
  <aside className="admin-sidebar">
    <NavLink to="/admin/products" className="admin-menu-btn">Товары</NavLink>
    <NavLink to="/admin/orders" className="admin-menu-btn">Заказы</NavLink>
  </aside>
);

export default AdminSidebar; 