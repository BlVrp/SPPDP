import React, { createContext, useContext, useEffect, useState } from 'react';
import { useRouter } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface User {
    id: string;
    city: string;
    email: string;
    fileName: string;
    firstName: string;
    lastName: string;
    phoneNumber: string;
    post: string;
    postDepartment: string;
    website: string;
  }
  

  interface AuthContextType {
    token: string | null;
    setToken: (token: string | null) => void;
    isLoggedIn: boolean;
    logout: () => void;
    loading: boolean;
    user: User | null;
    setUser: (user: User | null) => void;
  }
  

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState<User | null>(null);
    const router = useRouter();
  
    useEffect(() => {
        const loadToken = async () => {
          const savedToken = await AsyncStorage.getItem('token');
          const savedUser = await AsyncStorage.getItem('user');
          setToken(savedToken);
          if (savedUser) {
            setUser(JSON.parse(savedUser));
          }
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
        await AsyncStorage.multiRemove(['token', 'user']);
        setToken(null);
        setUser(null);
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
          user,
          setUser,
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
