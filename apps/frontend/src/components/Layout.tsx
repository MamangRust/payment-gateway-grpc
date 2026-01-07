import { ReactNode } from 'react';
import BottomNav from './BottomNav';

interface LayoutProps {
  children: ReactNode;
  showBottomNav?: boolean;
}

export default function Layout({ children, showBottomNav = true }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 pb-16 md:pb-0">
      {/* Mobile padding for bottom navigation */}
      <div className="md:hidden h-16"></div>
      
      {/* Main content */}
      <main>
        {children}
      </main>

      {/* Bottom navigation for mobile */}
      {showBottomNav && <BottomNav />}
    </div>
  );
}
