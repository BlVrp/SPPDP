import React, { createContext, useContext, useEffect, useState } from 'react';
import { useRouter } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface AuthContextType {
  token: string | null;
  setToken: (token: string | null) => void;
  isLoggedIn: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const router = useRouter();
  
    useEffect(() => {
      const loadToken = async () => {
        const savedToken = await AsyncStorage.getItem('token');
        setToken(savedToken);
        setLoading(false);
      };
      loadToken();
    }, []);
  
    useEffect(() => {
      console.log("🔐 Auth State Changed:");
      console.log("Token:", token);
      console.log("Is Logged In:", !!token && !loading);
    }, [token, loading]);
  
    const logout = async () => {
      await AsyncStorage.removeItem('token');
      setToken(null);
      router.replace('/auth/login');
    };
  
    return (
      <AuthContext.Provider
        value={{
          token,
          setToken,
          isLoggedIn: !!token && !loading,
          logout,
          loading,
        }}
      >
        {children}
      </AuthContext.Provider>
    );
  };

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};
