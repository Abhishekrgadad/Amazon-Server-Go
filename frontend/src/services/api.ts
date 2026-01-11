import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth APIs
export const authAPI = {
  login: async (email: string, password: string, role: string) => {
    const response = await api.post('/auth/login', { email, password, role });
    return response.data;
  },
  registerUser: async (data: any) => {
    const response = await api.post('/auth/register/user', data);
    return response.data;
  },
  registerSeller: async (data: any) => {
    const response = await api.post('/auth/register/seller', data);
    return response.data;
  },
  registerAdmin: async (data: any) => {
    const response = await api.post('/auth/register/admin', data);
    return response.data;
  },
  resetPassword: async (email: string, role: string) => {
    const response = await api.post('/auth/reset-password', { email, role });
    return response.data;
  },
  updatePassword: async (data: any) => {
    const response = await api.post('/auth/update-password', data);
    return response.data;
  },
};

// Product APIs
export const productAPI = {
  getActiveProducts: async (page: number = 1) => {
    const response = await api.get(`/auth/products/true/page${page}`);
    return response.data;
  },
  getProductById: async (id: string) => {
    const response = await api.get(`/auth/products/get/${id}`);
    return response.data;
  },
  filterProducts: async (filters: {
    name?: string;
    category?: string;
    brand?: string;
    min_price?: number;
    max_price?: number;
  }) => {
    const params = new URLSearchParams();
    if (filters.name) params.append('name', filters.name);
    if (filters.category) params.append('category', filters.category);
    if (filters.brand) params.append('brand', filters.brand);
    if (filters.min_price) params.append('min_price', filters.min_price.toString());
    if (filters.max_price) params.append('max_price', filters.max_price.toString());
    
    const response = await api.get(`/auth/products/filter?${params.toString()}`);
    return response.data;
  },
};

// Cart APIs
export const cartAPI = {
  addToCart: async (userId: string, items: { product_id: string; quantity: number }[]) => {
    const response = await api.post('/auth/cart/add', {
      user_id: userId,
      items,
    });
    return response.data;
  },
  getCart: async (userId: string) => {
    const response = await api.get(`/auth/cart/view?user_id=${userId}`);
    return response.data;
  },
  updateCart: async (userId: string, productId: string, quantity: number) => {
    const response = await api.put('/auth/cart/update', {
      user_id: userId,
      product_id: productId,
      quantity,
    });
    return response.data;
  },
  removeFromCart: async (userId: string, productId: string) => {
    const response = await api.delete(`/auth/cart/delete?user_id=${userId}&product_id=${productId}`);
    return response.data;
  },
  clearCart: async (userId: string) => {
    const response = await api.delete(`/auth/cart/clear?user_id=${userId}`);
    return response.data;
  },
};

// Order APIs
export const orderAPI = {
  checkout: async (data: {
    user_id: string;
    cart_id: string;
    payment_type: string;
    address: string;
    coupon_code?: string;
  }) => {
    const response = await api.post('/auth/order/checkout', data);
    return response.data;
  },
  viewOrders: async (userId: string, page: number = 1) => {
    const response = await api.get(`/auth/order/view?user_id=${userId}&page=${page}`);
    return response.data;
  },
  cancelOrder: async (orderId: string, userId: string) => {
    const response = await api.post(`/auth/order/cancel?order_id=${orderId}&user_id=${userId}`);
    return response.data;
  },
  returnOrder: async (orderId: string) => {
    const response = await api.post(`/auth/order/return?order_id=${orderId}`);
    return response.data;
  },
  checkStatus: async (orderId: string) => {
    const response = await api.post(`/auth/order/status?order_id=${orderId}`);
    return response.data;
  },
};

// Review APIs
export const reviewAPI = {
  addReview: async (data: {
    product_id: string;
    order_id: string;
    rating: number;
    comment: string;
  }) => {
    const response = await api.post('/auth/review/add', data);
    return response.data;
  },
  getProductReviews: async (productId: string) => {
    const response = await api.get(`/auth/review/view/${productId}`);
    return response.data;
  },
  updateReview: async (data: {
    product_id: string;
    rating: number;
    comment: string;
  }) => {
    const response = await api.put('/auth/review/update', data);
    return response.data;
  },
  deleteReview: async (productId: string) => {
    const response = await api.delete(`/auth/review/delete/${productId}`);
    return response.data;
  },
};

// Coupon APIs
export const couponAPI = {
  getAllCoupons: async () => {
    const response = await api.get('/auth/coupons/view');
    return response.data;
  },
  getCouponById: async (id: string) => {
    const response = await api.get(`/auth/coupons/view/${id}`);
    return response.data;
  },
};

export default api;
