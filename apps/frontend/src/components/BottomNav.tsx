import { useLocation, Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';

export default function BottomNav() {
  const location = useLocation();
  const isAuthenticated = useAuthStore(state => state.isAuthenticated);

  // Only show bottom nav on mobile when authenticated
  if (!isAuthenticated) return null;

  const navigationItems = [
    { path: '/dashboard', label: 'Home', icon: '🏠' },
    { path: '/cards', label: 'Cards', icon: '💳' },
    { path: '/topup', label: 'Topup', icon: '💰' },
    { path: '/transfer', label: 'Transfer', icon: '💸' },
    { path: '/saldo', label: 'Balance', icon: '💵' },
  ];

  return (
    <div className="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50">
      <nav className="bg-white shadow-lg">
        <div className="flex justify-around items-center h-16">
          {navigationItems.map((item) => {
            const isActive = location.pathname === item.path;
            
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex flex-col items-center justify-center py-1 px-2 min-w-[60px] transition-colors ${
                  isActive
                    ? 'text-blue-600'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <span className="text-xl mb-1">{item.icon}</span>
                <span className={`text-xs font-medium ${isActive ? 'text-blue-600' : ''}`}>
                  {item.label}
                </span>
                {isActive && (
                  <div className="w-1 h-1 bg-blue-600 rounded-full mt-1"></div>
                )}
              </Link>
            );
          })}
        </div>
      </nav>
    </div>
  );
}
