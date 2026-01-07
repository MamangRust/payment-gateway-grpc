import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { cardApi, saldoApi } from '../services/api';
import { useAuthStore } from '../store/authStore';
import Layout from './Layout';
import type { Card, Saldo } from '../types/api';

export default function SaldoPage() {
  const [cards, setCards] = useState<Card[]>([]);
  const [saldo, setSaldo] = useState<Saldo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const token = useAuthStore(state => state.token);

  useEffect(() => {
    const fetchData = async () => {
      if (!token) return;
      
      try {
        const [cardsData, saldoData] = await Promise.all([
          cardApi.getMyCards(token),
          saldoApi.getMySaldo(token)
        ]);

        console.log('Saldo cards:', cardsData);
        setCards(cardsData);
        setSaldo(saldoData);
      } catch (err: any) {
        console.error('Saldo fetch error:', err);
        setError('Failed to fetch card balances');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [token]);

  const totalBalance = saldo.reduce((sum, s) => sum + s.total_balance, 0);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <Layout>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center py-6">
            <Link to="/dashboard" className="text-indigo-600 hover:text-indigo-800 mr-4">
              ← Back
            </Link>
            <h1 className="text-2xl font-bold text-gray-900">Card Balances</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Total Balance Card */}
        <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-lg font-medium text-white mb-2">Total Balance</h2>
          <p className="text-4xl font-bold text-white">
            ${totalBalance.toLocaleString()}
          </p>
          <p className="text-indigo-100 mt-2">Across {cards.length} card(s)</p>
        </div>

        {/* Cards with Balances */}
        <div className="space-y-6">
          {cards.length === 0 ? (
            <div className="bg-white rounded-lg shadow-md p-8 text-center">
              <div className="text-gray-400 text-5xl mb-4">💳</div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">No Cards Found</h3>
              <p className="text-gray-500 mb-4">You don't have any cards yet</p>
              <Link
                to="/topup"
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
              >
                Add a Card
              </Link>
            </div>
          ) : (
            cards.map((card) => {
              const cardSaldo = saldo.find(s => s.card_number === card.card_number);
              const balance = cardSaldo ? cardSaldo.total_balance : 0;

              return (
                <div key={card.card_id} className="bg-white rounded-lg shadow-md overflow-hidden">
                  <div className="md:flex">
                    {/* Card Representation */}
                    <div className="md:w-1/3 bg-gradient-to-br from-gray-800 to-gray-900 p-6 text-white">
                      <div className="mb-4">
                        <div className="text-xs uppercase tracking-wider opacity-70 mb-1">Card Number</div>
                        <div className="text-lg font-mono">•••• •••• •••• {card.card_number.slice(-4)}</div>
                      </div>
                      <div className="mb-4">
                        <div className="text-xs uppercase tracking-wider opacity-70 mb-1">Card Type</div>
                        <div className="text-sm font-medium">{card.card_type} Card</div>
                      </div>
                      <div>
                        <div className="text-xs uppercase tracking-wider opacity-70 mb-1">Provider</div>
                        <div className="text-sm font-medium">{card.card_provider}</div>
                      </div>
                    </div>

                    {/* Balance Info and Actions */}
                    <div className="md:w-2/3 p-6">
                      <div className="flex justify-between items-start mb-4">
                        <div>
                          <h3 className="text-xl font-semibold text-gray-900">Current Balance</h3>
                          <p className="text-3xl font-bold text-indigo-600 mt-1">
                            ${balance.toLocaleString()}
                          </p>
                        </div>

                        <div className="flex space-x-2">
                          <Link
                            to="/topup"
                            className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white text-sm font-medium rounded-md transition-colors"
                          >
                            Top Up
                          </Link>
                          <Link
                            to="/transfer"
                            className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium rounded-md transition-colors"
                          >
                            Transfer
                          </Link>
                          <Link
                            to="/withdraw"
                            className="px-4 py-2 bg-orange-500 hover:bg-orange-600 text-white text-sm font-medium rounded-md transition-colors"
                          >
                            Withdraw
                          </Link>
                        </div>
                      </div>

                      {/* Quick Stats */}
                      <div className="grid grid-cols-3 gap-4 pt-4 border-t border-gray-200">
                        <div className="text-center">
                          <div className="text-2xl font-bold text-gray-900">3</div>
                          <div className="text-xs text-gray-500">Transactions</div>
                        </div>
                        <div className="text-center">
                          <div className="text-2xl font-bold text-green-600">+2</div>
                          <div className="text-xs text-gray-500">Income</div>
                        </div>
                        <div className="text-center">
                          <div className="text-2xl font-bold text-red-600">01</div>
                          <div className="text-xs text-gray-500">Expense</div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Balance Summary */}
        {cards.length > 0 && (
          <div className="mt-8 bg-white rounded-lg shadow-md p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Balance Summary</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {cards.map((card) => {
                const cardSaldo = saldo.find(s => s.card_number === card.card_number);
                const balance = cardSaldo ? cardSaldo.total_balance : 0;

                return (
                  <div key={card.card_id} className="border border-gray-200 rounded-md p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-900">
                        {card.card_provider}
                      </span>
                      <span className="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded">
                        •••• {card.card_number.slice(-4)}
                      </span>
                    </div>
                    <div className="text-xl font-bold text-indigo-600">
                      ${balance.toLocaleString()}
                    </div>
                    <div className="text-xs text-gray-500 mt-1">
                      {card.card_type} Card
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </main>
    </Layout>
  );
}
