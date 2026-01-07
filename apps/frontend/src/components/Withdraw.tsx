import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { cardApi, withdrawApi, saldoApi } from '../services/api';
import { useAuthStore } from '../store/authStore';
import Layout from './Layout';
import type { Card, Saldo } from '../types/api';

export default function Withdraw() {
  const [cards, setCards] = useState<Card[]>([]);
  const [saldo, setSaldo] = useState<Saldo[]>([]);
  const [selectedCard, setSelectedCard] = useState('');
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState('');
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
        console.log('Withdraw cards:', cardsData);
        setCards(cardsData);
        setSaldo(saldoData);
        
        if (cardsData.length > 0) {
          setSelectedCard(cardsData[0].card_number);
        }
      } catch (err: any) {
        console.error('Withdraw fetch error:', err);
        setError('Failed to fetch cards');
      }
    };

    fetchData();
  }, [token]);

  const getCardBalance = (cardNumber: string) => {
    const cardSaldo = saldo.find(s => s.card_number === cardNumber);
    return cardSaldo ? cardSaldo.total_balance : 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    if (!selectedCard || !amount || parseFloat(amount) <= 0) {
      setError('Please fill all fields correctly');
      setLoading(false);
      return;
    }

    const balance = getCardBalance(selectedCard);
    if (parseFloat(amount) > balance) {
      setError('Insufficient balance');
      setLoading(false);
      return;
    }

    try {
      const withdrawData = {
        card_number: selectedCard,
        withdraw_amount: parseFloat(amount),
        withdraw_time: new Date().toISOString()
      };

      if (!token) {
        throw new Error('No authentication token');
      }

      await withdrawApi.createWithdraw(withdrawData, token);
      setSuccess('Withdrawal successful!');
      setAmount('');
    } catch (err: any) {
      console.error('Withdraw error:', err);
      setError(err.response?.data?.message || 'Withdrawal failed');
    } finally {
      setLoading(false);
    }
  };

  const quickAmounts = [100, 250, 500, 1000, 2500];
  const maxWithdrawal = selectedCard ? getCardBalance(selectedCard) : 0;

  return (
    <Layout>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center py-6">
            <Link to="/dashboard" className="text-indigo-600 hover:text-indigo-800 mr-4">
              ← Back
            </Link>
            <h1 className="text-2xl font-bold text-gray-900">Withdraw Money</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Withdraw from your card</h2>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
              {error}
            </div>
          )}

          {success && (
            <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-6">
              {success}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Card Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Select Card
              </label>
              {cards.length === 0 ? (
                <p className="text-gray-500">No cards available</p>
              ) : (
                <div>
                  <select
                    value={selectedCard}
                    onChange={(e) => setSelectedCard(e.target.value)}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 mb-2"
                  >
                    {cards.map((card) => (
                      <option key={card.card_id} value={card.card_number}>
                        {card.card_provider} - {card.card_type} (•••• {card.card_number.slice(-4)})
                      </option>
                    ))}
                  </select>
                  <p className="text-sm text-gray-500">
                    Available balance: ${getCardBalance(selectedCard).toLocaleString()}
                  </p>
                </div>
              )}
            </div>

            {/* Amount Input */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Withdrawal Amount
              </label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">
                  $
                </span>
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="block w-full pl-8 pr-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  placeholder="0.00"
                  min="0"
                  step="0.01"
                  max={maxWithdrawal > 0 ? maxWithdrawal : undefined}
                />
              </div>
              {selectedCard && (
                <p className="text-sm text-gray-500 mt-1">
                  Maximum withdrawal: ${maxWithdrawal.toLocaleString()}
                </p>
              )}
            </div>

            {/* Quick Amount Buttons */}
            <div>
              <p className="text-sm font-medium text-gray-700 mb-2">Quick amounts:</p>
              <div className="grid grid-cols-5 gap-3">
                {quickAmounts.map((quickAmount) => (
                  <button
                    key={quickAmount}
                    type="button"
                    onClick={() => setAmount(quickAmount.toString())}
                    disabled={quickAmount > maxWithdrawal}
                    className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    ${quickAmount.toLocaleString()}
                  </button>
                ))}
              </div>
            </div>

            {/* Custom Amount Buttons */}
            <div>
              <p className="text-sm font-medium text-gray-700 mb-2">Quick percentages:</p>
              <div className="grid grid-cols-4 gap-3">
                {[25, 50, 75, 100].map((percentage) => {
                  const calcAmount = (maxWithdrawal * percentage) / 100;
                  return (
                    <button
                      key={percentage}
                      type="button"
                      onClick={() => setAmount(calcAmount.toString())}
                      disabled={calcAmount <= 0}
                      className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                    >
                      {percentage}% (${calcAmount.toLocaleString()})
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Warning */}
            <div className="bg-yellow-50 border border-yellow-200 p-4 rounded-lg">
              <div className="flex">
                <div className="flex-shrink-0">
                  <span className="text-yellow-400">⚠️</span>
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-yellow-800">Withdrawal Notice</h3>
                  <div className="mt-2 text-sm text-yellow-700">
                    <ul className="list-disc list-inside space-y-1">
                      <li>Withdrawals may take 1-3 business days to process</li>
                      <li>Bank fees may apply depending on your withdrawal method</li>
                      <li>Daily withdrawal limits may apply</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            {/* Submit Button */}
            <div>
              <button
                type="submit"
                disabled={loading || cards.length === 0 || maxWithdrawal <= 0}
                className="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500 disabled:opacity-50"
              >
                {loading ? 'Processing...' : `Withdraw ${amount ? '$' + parseFloat(amount).toLocaleString() : ''}`}
              </button>
            </div>
          </form>

          {/* Withdrawal Limits Info */}
          <div className="mt-8 p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <h3 className="text-sm font-semibold text-gray-900 mb-2">Withdrawal Limits:</h3>
            <div className="text-sm text-gray-600 space-y-1">
              <p>• Maximum per transaction: $10,000</p>
              <p>• Maximum daily: $25,000</p>
              <p>• Maximum weekly: $75,000</p>
              <p>• Card available balance: ${maxWithdrawal.toLocaleString()}</p>
            </div>
          </div>

          {/* Recent Withdrawals */}
          <div className="mt-8 pt-6 border-t border-gray-200">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Withdrawals</h3>
            <div className="text-gray-500">
              <p>No recent withdrawals found</p>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}
