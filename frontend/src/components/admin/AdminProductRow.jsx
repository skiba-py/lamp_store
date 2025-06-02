import React, { useState } from 'react';
import { PRODUCTS_API_URL } from '../../services/api';

const AdminProductRow = ({ product, isNew, onChange, onSave, onDelete }) => {
  const [edit, setEdit] = useState(isNew || false);
  const [form, setForm] = useState(product);

  React.useEffect(() => {
    setForm(product);
  }, [product]);

  const handleChange = e => {
    const { name, value } = e.target;
    const processedValue = name === 'stock' || name === 'price' ? Number(value) || 0 : value;
    setForm(f => ({ ...f, [name]: processedValue }));
    if (onChange) onChange({ ...form, [name]: processedValue });
  };

  const handleImageChange = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await fetch(`${PRODUCTS_API_URL}/products/${product.id}/image`, {
        method: 'POST',
        body: formData,
      });
      
      if (!response.ok) throw new Error('Failed to upload image');
      
      const data = await response.json();
      setForm(f => ({ ...f, image: data.image }));
      if (onChange) onChange({ ...form, image: data.image });
    } catch (error) {
      console.error('Error uploading image:', error);
    }
  };

  const handleSave = () => {
    if (onSave) onSave(form);
    setEdit(false);
  };

  if (isNew) {
    return (
      <tr className="admin-product-row new-row">
        <td><input name="id" value={form.id || ''} disabled /></td>
        <td><input name="name" value={form.name || ''} onChange={handleChange} placeholder="Наименование" /></td>
        <td><input name="description" value={form.description || ''} onChange={handleChange} placeholder="Описание" /></td>
        <td><input name="price" value={form.price || ''} onChange={handleChange} placeholder="Цена" type="number" /></td>
        <td><input name="stock" value={form.stock || ''} onChange={handleChange} placeholder="Количество" type="number" /></td>
        <td>
          <input 
            type="file" 
            accept="image/*" 
            onChange={handleImageChange}
            style={{ display: 'none' }}
            id="image-upload"
          />
          <label htmlFor="image-upload" className="upload-button">
            Загрузить картинку
          </label>
        </td>
        <td><button onClick={handleSave}>Добавить</button></td>
      </tr>
    );
  }

  return (
    <tr className="admin-product-row">
      <td>{product.id}</td>
      <td>{edit ? <input name="name" value={form.name || ''} onChange={handleChange} /> : product.name}</td>
      <td>{edit ? <input name="description" value={form.description || ''} onChange={handleChange} /> : product.description}</td>
      <td>{edit ? <input name="price" value={form.price || ''} onChange={handleChange} type="number" /> : product.price}</td>
      <td>{edit ? <input name="stock" value={form.stock || ''} onChange={handleChange} type="number" /> : product.stock}</td>
      <td>
        {edit ? (
          <>
            <input 
              type="file" 
              accept="image/*" 
              onChange={handleImageChange}
              style={{ display: 'none' }}
              id={`image-upload-${product.id}`}
            />
            <label htmlFor={`image-upload-${product.id}`} className="upload-button">
              Загрузить картинку
            </label>
          </>
        ) : (
          product.image && (
            <img 
              src={product.image} 
              alt={product.name} 
              style={{ width: '50px', height: '50px', objectFit: 'cover' }} 
            />
          )
        )}
      </td>
      <td>
        {edit ? (
          <>
            <button onClick={handleSave}>Сохранить</button>
            <button onClick={() => setEdit(false)}>Отмена</button>
          </>
        ) : (
          <>
            <button onClick={() => setEdit(true)}>Редактировать</button>
            <button onClick={onDelete}>Удалить</button>
          </>
        )}
      </td>
    </tr>
  );
};

export default AdminProductRow; 