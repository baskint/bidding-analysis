// src/contexts/AuthContext.tsx
"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { User, onAuthStateChanged, signOut as firebaseSignOut } from "firebase/auth";
import { auth } from "@/lib/firebase";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  signOut: async () => {},
});

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    console.log("AuthProvider: Setting up auth listener");

    const unsubscribe = onAuthStateChanged(
      auth,
      (user) => {
        console.log(
          "AuthProvider: Auth state changed",
          user ? "User logged in" : "User logged out"
        );
        setUser(user);
        setLoading(false);
      },
      (error) => {
        console.error("AuthProvider: Auth state change error:", error);
        setLoading(false);
      }
    );

    return () => {
      console.log("AuthProvider: Cleaning up auth listener");
      unsubscribe();
    };
  }, []);

  const signOut = async () => {
    try {
      console.log("AuthProvider: Signing out user");
      await firebaseSignOut(auth);
    } catch (error) {
      console.error("AuthProvider: Error signing out:", error);
      throw error;
    }
  };

  const value = {
    user,
    loading,
    signOut,
  };

  console.log("AuthProvider: Rendering with state:", { user: !!user, loading });

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
