import React, { createContext, useContext, useState, useEffect } from 'react';
import { setAuthToken } from '../services/authService';
import { authService } from '../services/authService';



const AuthContext = createContext(undefined);

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [userInfo, setUserInfo] = useState(() => {

    const stored = localStorage.getItem('userInfo');
    return stored ? JSON.parse(stored) : null;
  });

  useEffect(() => {
    if (token) {
      setAuthToken(token);
    }
  }, [token]);

  const login = async (email, password) => {
    const response = await authService.login(email, password);
    if (response.success && response.token) {
      setToken(response.token);
      setUserInfo(response);
      localStorage.setItem('token', response.token);
      localStorage.setItem('userInfo', JSON.stringify(response));
      setAuthToken(response.token);
    } else {
      throw new Error(response.message || 'Login failed');
    }
  };

  const logout = () => {
    if (token) {
      authService.revokeToken(token).catch(console.error);
    }
    setToken(null);
    setUserInfo(null);
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    setAuthToken(null);
  };

  return (
    <AuthContext.Provider
      value={{
        token,
        userInfo,
        login,
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
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
