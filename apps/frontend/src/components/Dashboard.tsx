import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { cardApi, saldoApi, authApi } from '../services/api';
import { useAuthStore } from '../store/authStore';
import Layout from './Layout';
import type { Card, Saldo, User } from '../types/api';

export default function Dashboard() {
  const [user, setUser] = useState<User | null>(null);
  const [cards, setCards] = useState<Card[]>([]);
  const [saldo, setSaldo] = useState<Saldo[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const token = useAuthStore(state => state.token);
  const logout = useAuthStore(state => state.logout);

  useEffect(() => {
    const fetchData = async () => {
      if (!token) return;
      
      try {
        const [userData, cardsData, saldoData] = await Promise.all([
          authApi.getMe(token),
          cardApi.getMyCards(token),
          saldoApi.getMySaldo(token)
        ]);
        
        setUser(userData);
        setCards(cardsData);
        setSaldo(saldoData);
      } catch (error) {
        console.error('Failed to fetch data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [token]);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  const totalBalance = saldo.reduce((sum, s) => sum + s.total_balance, 0);

  return (
    <Layout>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <h1 className="text-3xl font-bold text-gray-900">
              Welcome, {user?.firstname} {user?.lastname}
            </h1>
            <button
              onClick={handleLogout}
              className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Balance Card */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">Total Balance</h2>
          <p className="text-4xl font-bold text-indigo-600">
            ${totalBalance.toLocaleString()}
          </p>
          <p className="text-sm text-gray-500 mt-2">Across {cards.length} card(s)</p>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          <Link
            to="/topup"
            className="bg-green-500 hover:bg-green-600 text-white rounded-lg p-6 text-center transition-colors"
          >
            <div className="text-3xl mb-2">💰</div>
            <h3 className="font-semibold">Top Up</h3>
            <p className="text-sm opacity-90">Add funds to card</p>
          </Link>

          <Link
            to="/transfer"
            className="bg-blue-500 hover:bg-blue-600 text-white rounded-lg p-6 text-center transition-colors"
          >
            <div className="text-3xl mb-2">💸</div>
            <h3 className="font-semibold">Transfer</h3>
            <p className="text-sm opacity-90">Send money to others</p>
          </Link>

          <Link
            to="/withdraw"
            className="bg-orange-500 hover:bg-orange-600 text-white rounded-lg p-6 text-center transition-colors"
          >
            <div className="text-3xl mb-2">🏧</div>
            <h3 className="font-semibold">Withdraw</h3>
            <p className="text-sm opacity-90">Withdraw from card</p>
          </Link>

          <Link
            to="/cards"
            className="bg-purple-500 hover:bg-purple-600 text-white rounded-lg p-6 text-center transition-colors"
          >
            <div className="text-3xl mb-2">💳</div>
            <h3 className="font-semibold">Cards</h3>
            <p className="text-sm opacity-90">Manage cards</p>
          </Link>

          <Link
            to="/transactions"
            className="bg-pink-500 hover:bg-pink-600 text-white rounded-lg p-6 text-center transition-colors"
          >
            <div className="text-3xl mb-2">📊</div>
            <h3 className="font-semibold">History</h3>
            <p className="text-sm opacity-90">View transactions</p>
          </Link>
        </div>

        {/* Cards List */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Your Cards</h2>
          {cards.length === 0 ? (
            <p className="text-gray-500">No cards found</p>
          ) : (
            <div className="space-y-4">
              {cards.map((card) => (
                <div key={card.card_id} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <div>
                      <p className="font-medium text-gray-900">•••• •••• •••• {card.card_number.slice(-4)}</p>
                      <p className="text-sm text-gray-500">{card.card_provider} - {card.card_type}</p>
                    </div>
                    <Link
                      to="/saldo"
                      className="text-indigo-600 hover:text-indigo-800 font-medium"
                    >
                      View Balance
                    </Link>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Recent Transactions */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold text-gray-900">Quick Access</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Link
              to="/saldo"
              className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
            >
              <div className="text-2xl mb-2">💳</div>
              <h3 className="font-medium">View Saldo</h3>
              <p className="text-sm text-gray-500">Check card balances</p>
            </Link>
            
            <Link
              to="/transactions"
              className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors">
              <div className="text-2xl mb-2">📋</div>
              <h3 className="font-medium">Transactions</h3>
              <p className="text-sm text-gray-500">View all transactions</p>
            </Link>
            
            <Link
              to="/transfer"
              className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors">
              <div className="text-2xl mb-2">🔄</div>
              <h3 className="font-medium">Quick Transfer</h3>
              <p className="text-sm text-gray-500">Send money instantly</p>
            </Link>
          </div>
        </div>
      </main>
    </Layout>
  );
}
