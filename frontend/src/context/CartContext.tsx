import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { cartAPI } from '../services/api';
import { Cart, CartItem } from '../types';
import { useAuth } from './AuthContext';

interface CartContextType {
  cart: Cart | null;
  loading: boolean;
  addToCart: (productId: string, quantity: number) => Promise<void>;
  updateCartItem: (productId: string, quantity: number) => Promise<void>;
  removeFromCart: (productId: string) => Promise<void>;
  clearCart: () => Promise<void>;
  refreshCart: () => Promise<void>;
  getCartItemCount: () => number;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const { user, isAuthenticated } = useAuth();
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(false);

  const refreshCart = async () => {
    if (!isAuthenticated || !user) return;
    
    try {
      setLoading(true);
      // Extract user ID from token or user object
      const userId = localStorage.getItem('userId') || '';
      if (userId) {
        const cartData = await cartAPI.getCart(userId);
        setCart(cartData);
      }
    } catch (error) {
      console.error('Error fetching cart:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isAuthenticated && user) {
      refreshCart();
    }
  }, [isAuthenticated, user]);

  const addToCart = async (productId: string, quantity: number) => {
    if (!isAuthenticated || !user) {
      throw new Error('Please login to add items to cart');
    }
    
    try {
      setLoading(true);
      const userId = localStorage.getItem('userId') || '';
      await cartAPI.addToCart(userId, [{ product_id: productId, quantity }]);
      await refreshCart();
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to add to cart');
    } finally {
      setLoading(false);
    }
  };

  const updateCartItem = async (productId: string, quantity: number) => {
    if (!isAuthenticated || !user) return;
    
    try {
      setLoading(true);
      const userId = localStorage.getItem('userId') || '';
      await cartAPI.updateCart(userId, productId, quantity);
      await refreshCart();
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to update cart');
    } finally {
      setLoading(false);
    }
  };

  const removeFromCart = async (productId: string) => {
    if (!isAuthenticated || !user) return;
    
    try {
      setLoading(true);
      const userId = localStorage.getItem('userId') || '';
      await cartAPI.removeFromCart(userId, productId);
      await refreshCart();
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to remove from cart');
    } finally {
      setLoading(false);
    }
  };

  const clearCart = async () => {
    if (!isAuthenticated || !user) return;
    
    try {
      setLoading(true);
      const userId = localStorage.getItem('userId') || '';
      await cartAPI.clearCart(userId);
      setCart(null);
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to clear cart');
    } finally {
      setLoading(false);
    }
  };

  const getCartItemCount = () => {
    if (!cart || !cart.items) return 0;
    return cart.items.reduce((sum, item) => sum + item.quantity, 0);
  };

  return (
    <CartContext.Provider
      value={{
        cart,
        loading,
        addToCart,
        updateCartItem,
        removeFromCart,
        clearCart,
        refreshCart,
        getCartItemCount,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
