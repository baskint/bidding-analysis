// src/contexts/AuthContext.tsx
"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import {
  User as FirebaseUser,
  onAuthStateChanged,
  signOut as firebaseSignOut,
  signInWithPopup,
  GoogleAuthProvider,
} from "firebase/auth";
import { auth } from "@/lib/firebase";
import { getCurrentUser } from "@/lib/api";

// Get API base URL from environment
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Combined user type
interface CombinedUser {
  id: string;
  username: string;
  email?: string;
  authProvider: "google" | "local";
  created_at?: string;
  photoURL?: string;
}

interface AuthContextType {
  user: CombinedUser | null;
  loading: boolean;
  isAuthenticated: boolean;
  signInWithGoogle: () => Promise<boolean>;
  login: (username: string, password: string) => Promise<boolean>;
  register: (username: string, password: string) => Promise<boolean>;
  signOut: () => Promise<void>;
  error: string | null;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  isAuthenticated: false,
  signInWithGoogle: async () => false,
  login: async () => false,
  register: async () => false,
  signOut: async () => {},
  error: null,
  clearError: () => {},
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
  const [user, setUser] = useState<CombinedUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Initialize auth state
  useEffect(() => {
    console.log("AuthProvider: Initializing auth state");
    console.log("AuthProvider: Using API URL:", API_BASE_URL);

    // Check for stored JWT token (Go backend auth)
    const token = localStorage.getItem("auth_token");
    if (token) {
      console.log("AuthProvider: Found JWT token, verifying with Go backend");

      getCurrentUser()
        .then((response) => {
          setUser({
            id: response.user.id,
            username: response.user.username,
            authProvider: "local",
            created_at: response.user.created_at,
          });
          console.log("AuthProvider: JWT token verified, user authenticated");
          setLoading(false);
        })
        .catch(() => {
          console.log("AuthProvider: JWT token invalid, clearing auth");
          localStorage.removeItem("auth_token");
          localStorage.removeItem("user");
          setLoading(false);
        });
    } else {
      setLoading(false);
    }

    // Listen for Firebase auth state changes (Google auth)
    const unsubscribe = onAuthStateChanged(
      auth,
      (firebaseUser: FirebaseUser | null) => {
        if (firebaseUser && !token) {
          // User signed in with Google (and no existing JWT)
          console.log("AuthProvider: Google user authenticated");
          setUser({
            id: firebaseUser.uid,
            username: firebaseUser.displayName || firebaseUser.email || "Google User",
            email: firebaseUser.email || undefined,
            authProvider: "google",
          });
        } else if (!firebaseUser && !token) {
          // No auth from either source
          console.log("AuthProvider: No authentication found");
          setUser(null);
        }
        setLoading(false);
      },
      (authError) => {
        console.error("AuthProvider: Firebase auth error:", authError);
        setLoading(false);
      }
    );

    return () => {
      console.log("AuthProvider: Cleaning up auth listeners");
      unsubscribe();
    };
  }, []);

  const signInWithGoogle = async (): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      console.log("AuthProvider: Attempting Google sign in");

      const provider = new GoogleAuthProvider();
      await signInWithPopup(auth, provider);

      // User will be set via onAuthStateChanged listener
      console.log("AuthProvider: Google sign in successful");
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Google sign in failed";
      console.error("AuthProvider: Google sign in error:", errorMessage);
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const login = async (username: string, password: string): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      console.log("AuthProvider: Attempting Go backend login");

      const response = await fetch(`${API_BASE_URL}/trpc/auth.login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Login failed");
      }

      const loginResult = await response.json();

      // Store auth data
      localStorage.setItem("auth_token", loginResult.token);
      localStorage.setItem("user", JSON.stringify(loginResult.user));

      setUser({
        id: loginResult.user.id,
        username: loginResult.user.username,
        authProvider: "local",
        created_at: loginResult.user.created_at,
      });

      console.log("AuthProvider: Go backend login successful");
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Login failed";
      console.error("AuthProvider: Login error:", errorMessage);
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const register = async (username: string, password: string): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      console.log("AuthProvider: Attempting Go backend registration");

      const response = await fetch(`${API_BASE_URL}/trpc/auth.register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Registration failed");
      }

      const registerResult = await response.json();

      // Store auth data
      localStorage.setItem("auth_token", registerResult.token);
      localStorage.setItem("user", JSON.stringify(registerResult.user));

      setUser({
        id: registerResult.user.id,
        username: registerResult.user.username,
        authProvider: "local",
        created_at: registerResult.user.created_at,
      });

      console.log("AuthProvider: Go backend registration successful");
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Registration failed";
      console.error("AuthProvider: Registration error:", errorMessage);
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const signOut = async (): Promise<void> => {
    console.log("AuthProvider: Signing out user");
    setLoading(true);

    try {
      // Sign out from Google if applicable
      if (user?.authProvider === "google") {
        await firebaseSignOut(auth);
        console.log("AuthProvider: Signed out from Google");
      }

      // Clear JWT token and localStorage
      localStorage.removeItem("auth_token");
      localStorage.removeItem("user");

      // Clear state
      setUser(null);
      setError(null);

      console.log("AuthProvider: User signed out successfully");
    } catch (err) {
      console.error("AuthProvider: Error during sign out:", err);
      // Still clear local state even if Firebase signout fails
      localStorage.removeItem("auth_token");
      localStorage.removeItem("user");
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const clearError = () => {
    setError(null);
  };

  const isAuthenticated = !!user;

  const value = {
    user,
    loading,
    isAuthenticated,
    signInWithGoogle,
    login,
    register,
    signOut,
    error,
    clearError,
  };

  console.log("AuthProvider: Rendering with state:", {
    user: !!user,
    authProvider: user?.authProvider,
    loading,
    isAuthenticated,
    hasError: !!error,
    apiUrl: API_BASE_URL,
  });

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
