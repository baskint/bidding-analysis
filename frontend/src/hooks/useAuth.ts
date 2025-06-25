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
      localStorage.setItem('authToken', result.token);
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
      localStorage.setItem('authToken', result.token);
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

  const logout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
  };

  return {
    login,
    register,
    logout,
    loading,
    error,
  };
}
