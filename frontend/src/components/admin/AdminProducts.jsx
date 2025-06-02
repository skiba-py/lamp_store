import React, { useEffect, useState } from 'react';
import { adminApi } from '../../services/api';
import AdminProductRow from './AdminProductRow';

const PAGE_SIZE = 10;

const AdminProducts = () => {
  const [products, setProducts] = useState([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [newProduct, setNewProduct] = useState({ name: '', description: '', price: '', stock: '', image: '' });
  const token = localStorage.getItem('adminToken');

  const fetchProducts = async () => {
    setLoading(true);
    setError('');
    try {
      const products = await adminApi.getProducts(token);
      setProducts(products);
    } catch (err) {
      setError('Ошибка загрузки товаров');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
    // eslint-disable-next-line
  }, []);

  const handleCreate = async () => {
    try {
      await adminApi.createProduct({ ...newProduct, price: Number(newProduct.price), stock: Number(newProduct.stock) }, token);
      setNewProduct({ name: '', description: '', price: '', stock: '', image: '' });
      fetchProducts();
    } catch {
      setError('Ошибка создания товара');
    }
  };

  const handleUpdate = async (id, updated) => {
    try {
      await adminApi.updateProduct(id, updated, token);
      fetchProducts();
    } catch {
      setError('Ошибка обновления товара');
    }
  };

  const handleDelete = async (id) => {
    try {
      await adminApi.deleteProduct(id, token);
      fetchProducts();
    } catch {
      setError('Ошибка удаления товара');
    }
  };

  // Пагинация
  const pageCount = Math.ceil(products.length / PAGE_SIZE);
  const pagedProducts = products.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  return (
    <div className="admin-products">
      <h2>Товары</h2>
      {error && <div className="admin-error">{error}</div>}
      <div className="admin-products-table-wrapper">
        <table className="admin-products-table">
          <thead>
            <tr>
              <th>id товара</th>
              <th>Наименование</th>
              <th>Описание</th>
              <th>Цена</th>
              <th>Количество</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <AdminProductRow
              isNew
              product={newProduct}
              onChange={setNewProduct}
              onSave={handleCreate}
            />
            {pagedProducts.map(product => (
              <AdminProductRow
                key={product.id}
                product={product}
                onSave={updated => handleUpdate(product.id, updated)}
                onDelete={() => handleDelete(product.id)}
              />
            ))}
          </tbody>
        </table>
      </div>
      <div className="admin-pagination">
        <button onClick={() => setPage(1)} disabled={page === 1}>{'«'}</button>
        <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>{'<'}</button>
        {Array.from({ length: pageCount }, (_, i) => (
          <button
            key={i + 1}
            className={page === i + 1 ? 'active' : ''}
            onClick={() => setPage(i + 1)}
          >{i + 1}</button>
        ))}
        <button onClick={() => setPage(p => Math.min(pageCount, p + 1))} disabled={page === pageCount}>{'>'}</button>
        <button onClick={() => setPage(pageCount)} disabled={page === pageCount}>{'»'}</button>
      </div>
    </div>
  );
};

export default AdminProducts; 