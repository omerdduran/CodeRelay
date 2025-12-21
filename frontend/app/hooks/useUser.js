'use client';

import { useState, useEffect, useCallback } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export function useUser() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Load user from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem('coderelay_user');
    if (stored) {
      try {
        setUser(JSON.parse(stored));
      } catch {
        localStorage.removeItem('coderelay_user');
      }
    }
    setLoading(false);
  }, []);

  // Create or get user by nickname
  const login = useCallback(async (nickname) => {
    try {
      // First try to create user via API
      const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname }),
      });

      let userData;
      if (res.ok) {
        userData = await res.json();
      } else {
        // User might already exist, create local user object
        userData = { id: Date.now(), nickname };
      }

      // Store in localStorage
      localStorage.setItem('coderelay_user', JSON.stringify(userData));
      setUser(userData);
      return userData;
    } catch (error) {
      // Offline mode - create local user
      const userData = { id: Date.now(), nickname };
      localStorage.setItem('coderelay_user', JSON.stringify(userData));
      setUser(userData);
      return userData;
    }
  }, []);

  // Logout
  const logout = useCallback(() => {
    localStorage.removeItem('coderelay_user');
    setUser(null);
  }, []);

  return { user, loading, login, logout };
}
