export interface User {
  _id?: string;
  full_name: string;
  email: string;
  phone_number: string;
  shipping_address?: string;
  role: string;
}

export interface Product {
  _id: string;
  name: string;
  description: string;
  price: number;
  stock_quantity: number;
  category: string;
  brand: string;
  visibility: boolean;
  average_rating?: number;
  total_reviews?: number;
  review_comments?: string[];
  created_at?: string;
  updated_at?: string;
}

export interface CartItem {
  product_id: string;
  product_name: string;
  price: number;
  quantity: number;
  subtotal: number;
}

export interface Cart {
  user_id: string;
  cart_id: string;
  items: CartItem[];
  total_price: number;
  message?: string;
}

export interface Order {
  _id: string;
  user_id: string;
  address: string;
  items: OrderItem[];
  total_price: number;
  payment_type: string;
  coupon_code?: string;
  discount: number;
  status: string;
  user_details: User;
  created_at: string;
}

export interface OrderItem {
  product_id: string;
  product_name: string;
  quantity: number;
  price: number;
}

export interface Review {
  _id?: string;
  user_id: string;
  product_id: string;
  order_id: string;
  rating: number;
  comment: string;
  created_at?: string;
}

export interface Coupon {
  _id: string;
  code: string;
  discount: number;
  active: boolean;
  is_percent: boolean;
  expiry_date: string;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
  role: 'customer' | 'seller' | 'admin';
}

export interface RegisterUserRequest {
  full_name: string;
  email: string;
  phone_number: string;
  password: string;
  shipping_address?: string;
  role: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string, role: string) => Promise<void>;
  register: (data: RegisterUserRequest) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}
