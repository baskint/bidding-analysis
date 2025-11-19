// frontend/src/app/dashboard/layout.tsx
'use client';

import { Sidebar } from "@/components/dashboard/Sidebar";
import { Header } from "@/components/dashboard/Header";
import ProtectedRoute from "@/components/auth/ProtectedRoute";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <ProtectedRoute>
        {/* FIX: The outer container must be updated to use theme-aware background classes.
          1. Removed static gradient.
          2. Added light mode default background (bg-gray-50)
          4. Added transition for a smooth visual switch.
        */}
        <div className='min-h-screen bg-gray-50 dark:bg-slate-900 transition-colors duration-500'>
          <div className='flex'>
            {/* Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <div className='flex-1 flex flex-col min-h-screen ml-64 bg-white dark:bg-slate-900 text-gray-800 dark:text-slate-100'>
              <Header />

              {/* Page Content */}
              <main className='flex-1 p-6'>
                <div className='max-w-7xl mx-auto bg-white dark:bg-slate-900'>{children}</div>
              </main>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    </>
  );
}