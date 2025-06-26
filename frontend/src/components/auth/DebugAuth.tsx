// src/components/auth/DebugAuth.tsx (temporary component for debugging)
'use client';

import { useAuth } from '@/contexts/AuthContext';

export default function DebugAuth() {
  const { user, loading } = useAuth();

  return (
    <div className="fixed top-0 right-0 bg-black text-white p-4 text-xs z-50 max-w-xs">
      <h3 className="font-bold">Auth Debug</h3>
      <p>Loading: {loading ? 'true' : 'false'}</p>
      <p>User: {user ? 'logged in' : 'null'}</p>
      {user && (
        <>
          <p>Email: {user.email}</p>
          <p>Name: {user.username || 'No name'}</p>
        </>
      )}
    </div>
  );
}
