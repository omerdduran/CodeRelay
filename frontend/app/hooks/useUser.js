'use client';

import { useState, useEffect, useCallback } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const TOKEN_KEY = 'coderelay_token';
const USER_KEY = 'coderelay_user';

export function useUser() {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);

  // Load user and token from localStorage on mount
  useEffect(() => {
    const storedToken = localStorage.getItem(TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedToken && storedUser) {
      try {
        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch {
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }
    setLoading(false);
  }, []);

  // Register new user
  const register = useCallback(async (nickname, password) => {
    const res = await fetch(`${API_URL}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nickname, password }),
    });

    if (!res.ok) {
      const data = await res.json();
      if (res.status === 409) {
        throw new Error('nickname_taken');
      }
      throw new Error(data.error || 'registration_failed');
    }

    const data = await res.json();

    // Store token and user
    localStorage.setItem(TOKEN_KEY, data.token);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setToken(data.token);
    setUser(data.user);

    return data.user;
  }, []);

  // Login existing user
  const login = useCallback(async (nickname, password) => {
    const res = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nickname, password }),
    });

    if (!res.ok) {
      if (res.status === 401) {
        throw new Error('invalid_credentials');
      }
      throw new Error('login_failed');
    }

    const data = await res.json();

    // Store token and user
    localStorage.setItem(TOKEN_KEY, data.token);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setToken(data.token);
    setUser(data.user);

    return data.user;
  }, []);

  // Logout
  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    setToken(null);
    setUser(null);
  }, []);

  // Get auth headers for API calls
  const getAuthHeaders = useCallback(() => {
    if (!token) return {};
    return { 'Authorization': `Bearer ${token}` };
  }, [token]);

  return { user, token, loading, register, login, logout, getAuthHeaders };
}
