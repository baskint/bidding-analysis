// src/hooks/useAuth.ts
import { useState } from 'react';
import { loginUser, registerUser, type AuthResponse } from '@/lib/api';

export function useAuthApi() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (username: string, password: string): Promise<AuthResponse | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await loginUser(username, password);
      // Store token and user info
      localStorage.setItem('auth_token', result.token);
      localStorage.setItem('user', JSON.stringify(result.user));
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Login failed';
      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const register = async (username: string, password: string): Promise<AuthResponse | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await registerUser(username, password);
      // Store token and user info
      localStorage.setItem('auth_token', result.token);
      localStorage.setItem('user', JSON.stringify(result.user));
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Registration failed';
      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const logout = (): void => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
  };

  // Helper to check if user is authenticated
  const isAuthenticated = (): boolean => {
    return !!localStorage.getItem('auth_token');
  };

  // Helper to get current user from localStorage
  const getCurrentUser = () => {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  };

  // Clear error
  const clearError = (): void => {
    setError(null);
  };

  return {
    login,
    register,
    logout,
    loading,
    error,
    isAuthenticated,
    getCurrentUser,
    clearError,
  };
}
