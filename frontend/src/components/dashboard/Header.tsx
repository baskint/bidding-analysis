// frontend/src/components/dashboard/Header.tsx
'use client';

import { Bell, Search, LogOut } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/contexts/AuthContext";
import { useState } from "react";

export function Header() {
  const { user } = useAuth();
  console.log("Header user:", user);
  const [showUserMenu, setShowUserMenu] = useState(false);

  return (
    <header className='bg-white border-b border-slate-200 px-6 py-4'>
      <div className='flex items-center justify-between'>
        {/* Search */}
        <div className='flex-1 max-w-md'>
          <div className='relative'>
            <Search className='absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400' />
            <input
              type='text'
              placeholder='Search campaigns, bids, or analytics...'
              className='w-full pl-10 pr-4 py-2 border border-slate-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent'
            />
          </div>
        </div>

        {/* Actions */}
        <div className='flex items-center space-x-4'>
          {/* Notifications */}
          <button className='relative p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors'>
            <Bell className='w-5 h-5' />
            <span className='absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full'></span>
          </button>

          {/* User Menu */}
          <div className='relative'>
            <div className='flex items-center space-x-3'>
              <div className='text-right'>
                <p className='text-sm font-medium text-slate-900'>{user?.username || "User"}</p>
              </div>
              <button
                onClick={() => setShowUserMenu(!showUserMenu)}
                className='relative w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-full flex items-center justify-center hover:from-blue-700 hover:to-purple-700 transition-all'
              >
                {user?.photoURL ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={user.photoURL}
                    alt='Profile'
                    className='w-8 h-8 rounded-full object-cover'
                  />
                ) : (
                  <span className='text-white text-sm font-medium'>
                    {(user?.username || user?.email || "U").charAt(0).toUpperCase()}
                  </span>
                )}
              </button>
            </div>

            {/* Dropdown Menu */}
            {showUserMenu && (
              <div className='absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-slate-200 py-1 z-50'>
                <div className='px-4 py-2 border-b border-slate-100'>
                  <p className='text-sm font-medium text-slate-900 truncate'>
                    {user?.username || "User"}
                  </p>
                  <p className='text-xs text-slate-500 truncate'>{user?.email}</p>
                </div>

                <Link
                  href='/logout'
                  className='flex items-center px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors'
                  onClick={() => setShowUserMenu(false)}
                >
                  <LogOut className='w-4 h-4 mr-2' />
                  Sign Out
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Overlay to close menu when clicking outside */}
      {showUserMenu && (
        <div className='fixed inset-0 z-40' onClick={() => setShowUserMenu(false)} />
      )}
    </header>
  );
}
