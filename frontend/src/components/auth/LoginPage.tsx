'use client';;

// Example: Login component that supports both auth methods

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";

// Default export for the main LoginPage component
export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isRegisterMode, setIsRegisterMode] = useState(false);

  const { signInWithGoogle, login, register, loading, error, clearError, isAuthenticated } =
    useAuth();
  const router = useRouter();

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/dashboard");
    }
  }, [isAuthenticated, router]);

  const handleGoogleSignIn = async () => {
    clearError();
    const success = await signInWithGoogle();
    if (success) {
      router.push("/dashboard");
    }
  };

  const handleCredentialsSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();

    const success = isRegisterMode
      ? await register(username, password)
      : await login(username, password);

    if (success) {
      router.push("/dashboard");
    }
  };

  const toggleMode = () => {
    setIsRegisterMode(!isRegisterMode);
    clearError();
    setUsername("");
    setPassword("");
  };

  return (
    <div className='min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4'>
      <div className='max-w-md w-full bg-white rounded-xl shadow-xl p-8'>
        <div className='text-center mb-8'>
          <h1 className='text-3xl font-bold text-gray-900 mb-2'>
            {isRegisterMode ? "Create Account" : "Welcome Back"}
          </h1>
          <p className='text-gray-600'>
            {isRegisterMode
              ? "Join us to start analyzing your bidding data"
              : "Sign in to access your bidding analysis dashboard"}
          </p>
        </div>

        {/* Error Message */}
        {error && (
          <div className='bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6'>
            {error}
          </div>
        )}

        {/* Google Sign In */}
        <button
          onClick={handleGoogleSignIn}
          disabled={loading}
          className='w-full bg-white border border-gray-300 text-gray-700 px-4 py-3 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200 flex items-center justify-center gap-3 mb-6 disabled:opacity-50 disabled:cursor-not-allowed'
        >
          <svg className='w-5 h-5' viewBox='0 0 24 24'>
            <path
              fill='#4285F4'
              d='M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z'
            />
            <path
              fill='#34A853'
              d='M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z'
            />
            <path
              fill='#FBBC05'
              d='M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z'
            />
            <path
              fill='#EA4335'
              d='M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z'
            />
          </svg>
          {loading ? "Please wait..." : "Continue with Google"}
        </button>

        <div className='relative mb-6'>
          <div className='absolute inset-0 flex items-center'>
            <div className='w-full border-t border-gray-300'></div>
          </div>
          <div className='relative flex justify-center text-sm'>
            <span className='bg-white px-4 text-gray-500'>Or continue with username</span>
          </div>
        </div>

        {/* Username/Password Form */}
        <form onSubmit={handleCredentialsSubmit} className='space-y-4'>
          <div>
            <label htmlFor='username' className='block text-sm font-medium text-gray-700 mb-2'>
              Username
            </label>
            <input
              id='username'
              type='text'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200 text-gray-700'
              placeholder='Enter your username'
              required
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor='password' className='block text-sm font-medium text-gray-700 mb-2'>
              Password
            </label>
            <input
              id='password'
              type='password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200 text-gray-700'
              placeholder='Enter your password'
              required
              disabled={loading}
              minLength={6}
            />
          </div>

          <button
            type='submit'
            disabled={loading}
            className='w-full bg-blue-600 text-white px-4 py-3 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition duration-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed'
          >
            {loading ? "Please wait..." : isRegisterMode ? "Create Account" : "Sign In"}
          </button>
        </form>

        {/* Toggle between login/register */}
        <div className='mt-6 text-center'>
          <button
            onClick={toggleMode}
            disabled={loading}
            className='text-blue-600 hover:text-blue-700 font-medium disabled:opacity-50'
          >
            {isRegisterMode ? "Already have an account? Sign in" : "Don't have an account? Sign up"}
          </button>
        </div>

        {/* Requirements for registration */}
        {isRegisterMode && (
          <div className='mt-4 p-3 bg-blue-50 rounded-lg'>
            <p className='text-sm text-blue-700 font-medium mb-1'>Requirements:</p>
            <ul className='text-xs text-blue-600 space-y-1'>
              <li>• Username: at least 3 characters</li>
              <li>• Password: at least 6 characters</li>
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

// Named export for UserProfile component
export function UserProfile() {
  const { user, signOut } = useAuth();
  const router = useRouter();

  const handleSignOut = async () => {
    await signOut();
    router.push("/login");
  };

  if (!user) return null;

  return (
    <div className='user-profile'>
      <div className='flex items-center gap-4'>
        <div>
          <h3>{user.username}</h3>
          {user.email && <p className='text-sm text-gray-500'>{user.email}</p>}
          <span className='inline-block bg-gray-100 text-gray-800 text-xs px-2 py-1 rounded'>
            {user.authProvider === "google" ? "🔗 Google" : "🔐 Local Account"}
          </span>
        </div>
        <button
          onClick={handleSignOut}
          className='bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700'
        >
          Sign Out
        </button>
      </div>
    </div>
  );
}
