const ORDERS_API_URL = 'http://localhost:8001/api';
export const PRODUCTS_API_URL = 'http://localhost:8000/api';
const ADMIN_API_URL = 'http://localhost:8003/api/admin';

const handleResponse = async (response) => {
  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || 'Произошла ошибка при выполнении запроса');
  }
  return response.json();
};

export const ordersApi = {
  createOrder: async (orderData) => {
    const response = await fetch(`${ORDERS_API_URL}/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(orderData),
    });
    return handleResponse(response);
  },

  getOrders: async (userId) => {
    const response = await fetch(`${ORDERS_API_URL}/orders/user/${userId}`);
    return handleResponse(response);
  },

  getOrderById: async (orderId) => {
    const response = await fetch(`${ORDERS_API_URL}/orders/${orderId}`);
    return handleResponse(response);
  },

  updateOrder: async (orderId, orderData) => {
    const response = await fetch(`${ORDERS_API_URL}/orders/${orderId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(orderData),
    });
    return handleResponse(response);
  },

  deleteOrder: async (orderId) => {
    const response = await fetch(`${ORDERS_API_URL}/orders/${orderId}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    return true;
  }
};

export const productsApi = {
  getProducts: async () => {
    const response = await fetch(`${PRODUCTS_API_URL}/products`);
    return handleResponse(response);
  },

  getProductById: async (productId) => {
    const response = await fetch(`${PRODUCTS_API_URL}/products/${productId}`);
    return handleResponse(response);
  },

  checkProductAvailability: async (productId, quantity) => {
    const response = await fetch(`${PRODUCTS_API_URL}/products/${productId}/availability`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ quantity }),
    });
    return handleResponse(response);
  }
};

export const confirmReservation = async (reservationId) => {
  try {
    const response = await fetch(`${PRODUCTS_API_URL}/products/reserve/${reservationId}/confirm`, {
      method: 'POST',
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
  } catch (error) {
    console.error('Ошибка при подтверждении резервирования:', error);
    throw error;
  }
};

export const cancelReservation = async (reservationId) => {
  try {
    const response = await fetch(`${PRODUCTS_API_URL}/products/reserve/${reservationId}/cancel`, {
      method: 'POST',
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
  } catch (error) {
    console.error('Ошибка при отмене резервирования:', error);
    throw error;
  }
};

export const adminApi = {
  login: async (credentials) => {
    const response = await fetch(`${ADMIN_API_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });
    return handleResponse(response);
  },

  getProducts: async (token) => {
    const response = await fetch(`${ADMIN_API_URL}/products`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return handleResponse(response);
  },

  createProduct: async (product, token) => {
    const response = await fetch(`${ADMIN_API_URL}/products`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(product),
    });
    return handleResponse(response);
  },

  updateProduct: async (id, product, token) => {
    const response = await fetch(`${ADMIN_API_URL}/products/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(product),
    });
    return handleResponse(response);
  },

  deleteProduct: async (id, token) => {
    const response = await fetch(`${ADMIN_API_URL}/products/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    return true;
  },

  getOrders: async (token) => {
    const response = await fetch(`${ADMIN_API_URL}/orders`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return handleResponse(response);
  },

  updateOrder: async (id, order, token) => {
    const response = await fetch(`${ADMIN_API_URL}/orders/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(order),
    });
    return handleResponse(response);
  },
}; 