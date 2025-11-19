import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/contexts/AuthContext";
// 1. Import the ThemeProvider from the new .tsx file

export const metadata: Metadata = {
  title: "Bidding Analysis Dashboard",
  description: "Advanced bidding analytics and campaign management",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      {/* 2. Apply theme-aware background to the body using slate for dark mode */}
      <body className="bg-gray-50">
        {/* 3. Wrap your providers with ThemeProvider at the top level */}
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}