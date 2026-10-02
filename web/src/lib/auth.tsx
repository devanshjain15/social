import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { api } from './api';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  register: (username: string, email: string, password: string) => Promise<{
    message: string;
    activation_url: string;
    activation_token: string;
  }>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    // Check if user has a valid token/cookie
    const checkAuth = async () => {
      const token = api.getToken();
      if (token) {
        try {
          // Try to fetch user data to verify token is valid
          await api.getAllUsers();
          setIsAuthenticated(true);
        } catch (error) {
          // Token is invalid, clear it
          api.clearToken();
          setIsAuthenticated(false);
        }
      } else {
        setIsAuthenticated(false);
      }
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const response = await api.login(email, password);
    // Store user info for filtering in explore page
    if (response && response.id) {
      localStorage.setItem('current_user', JSON.stringify(response));
    }
    setIsAuthenticated(true);
  };

  const logout = async () => {
    await api.logout();
    localStorage.removeItem('current_user');
    setIsAuthenticated(false);
  };

  const register = async (username: string, email: string, password: string) => {
    return await api.register(username, email, password);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, login, logout, register }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
