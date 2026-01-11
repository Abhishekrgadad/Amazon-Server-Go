import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authAPI } from '../services/api';
import { AuthContextType, User, RegisterUserRequest } from '../types';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));

  useEffect(() => {
    if (token) {
      // Decode token to get user info (simplified - in production, verify token)
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        // Store minimal user info from token
        setUser({
          email: payload.email,
          role: payload.role,
          full_name: '',
          phone_number: '',
        });
      } catch (error) {
        console.error('Error decoding token:', error);
        logout();
      }
    }
  }, [token]);

  const login = async (email: string, password: string, role: string) => {
    try {
      const response = await authAPI.login(email, password, role);
      const newToken = response.token;
      setToken(newToken);
      localStorage.setItem('token', newToken);
      
      // Set user from token payload
      const payload = JSON.parse(atob(newToken.split('.')[1]));
      setUser({
        email: payload.email,
        role: payload.role,
        full_name: '',
        phone_number: '',
      });
      // Store user ID from token (backend uses 'user_id' field)
      if (payload.user_id) {
        localStorage.setItem('userId', payload.user_id);
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed');
    }
  };

  const register = async (data: RegisterUserRequest) => {
    try {
      await authAPI.registerUser(data);
      // Auto login after registration
      await login(data.email, data.password, 'user');
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Registration failed');
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        register,
        logout,
        isAuthenticated: !!token,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
