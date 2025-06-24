// frontend/src/app/dashboard/layout.tsx
"use client";
import { Sidebar } from "@/components/dashboard/Sidebar";
import { Header } from "@/components/dashboard/Header";
import ProtectedRoute from "@/components/auth/ProtectedRoute";
import DebugAuth from "@/components/auth/DebugAuth";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <DebugAuth />
      <ProtectedRoute>
        <div className='min-h-screen bg-gradient-to-br from-slate-50 to-slate-100'>
          <div className='flex'>
            {/* Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <div className='flex-1 flex flex-col min-h-screen ml-64'>
              <Header />

              {/* Page Content */}
              <main className='flex-1 p-6'>
                <div className='max-w-7xl mx-auto'>{children}</div>
              </main>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    </>
  );
}
