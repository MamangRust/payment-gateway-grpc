import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { cardApi, transferApi, saldoApi } from '../services/api';
import { useAuthStore } from '../store/authStore';
import Layout from './Layout';
import type { Card, Saldo } from '../types/api';

export default function Transfer() {
  const [cards, setCards] = useState<Card[]>([]);
  const [saldo, setSaldo] = useState<Saldo[]>([]);
  const [fromCard, setFromCard] = useState('');
  const [toCard, setToCard] = useState('');
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
        console.log('Transfer cards:', cardsData);
        setCards(cardsData);
        setSaldo(saldoData);
        
        if (cardsData.length > 0) {
          setFromCard(cardsData[0].card_number);
        }
      } catch (err: any) {
        console.error('Transfer fetch error:', err);
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

    if (!fromCard || !toCard || !amount || parseFloat(amount) <= 0) {
      setError('Please fill all fields correctly');
      setLoading(false);
      return;
    }

    if (fromCard === toCard) {
      setError('Cannot transfer to the same card');
      setLoading(false);
      return;
    }

    const fromBalance = getCardBalance(fromCard);
    if (parseFloat(amount) > fromBalance) {
      setError('Insufficient balance');
      setLoading(false);
      return;
    }

    try {
      const transferData = {
        transfer_from: fromCard,
        transfer_to: toCard,
        transfer_amount: parseFloat(amount)
      };

      if (!token) {
        throw new Error('No authentication token');
      }

      await transferApi.createTransfer(transferData, token);
      setSuccess('Transfer successful!');
      setAmount('');
    } catch (err: any) {
      console.error('Transfer error:', err);
      setError(err.response?.data?.message || 'Transfer failed');
    } finally {
      setLoading(false);
    }
  };

  const quickAmounts = [100, 500, 1000, 5000];

  return (
    <Layout>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center py-6">
            <Link to="/dashboard" className="text-indigo-600 hover:text-indigo-800 mr-4">
              ← Back
            </Link>
            <h1 className="text-2xl font-bold text-gray-900">Transfer Money</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Send money to another card</h2>

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
            {/* From Card */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                From Card
              </label>
              {cards.length === 0 ? (
                <p className="text-gray-500">No cards available</p>
              ) : (
                <div>
                  <select
                    value={fromCard}
                    onChange={(e) => setFromCard(e.target.value)}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 mb-2"
                  >
                    {cards.map((card) => (
                      <option key={card.card_id} value={card.card_number}>
                        {card.card_provider} - {card.card_type} (•••• {card.card_number.slice(-4)})
                      </option>
                    ))}
                  </select>
                  <p className="text-sm text-gray-500">
                    Available balance: ${getCardBalance(fromCard).toLocaleString()}
                  </p>
                </div>
              )}
            </div>

            {/* To Card */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                To Card Number
              </label>
              <input
                type="text"
                value={toCard}
                onChange={(e) => setToCard(e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Enter destination card number"
              />
            </div>

            {/* Amount */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Amount
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
                  max={fromCard ? getCardBalance(fromCard) : undefined}
                />
              </div>
              {fromCard && (
                <p className="text-sm text-gray-500 mt-1">
                  Max: ${getCardBalance(fromCard).toLocaleString()}
                </p>
              )}
            </div>

            {/* Quick Amount Buttons */}
            <div>
              <p className="text-sm font-medium text-gray-700 mb-2">Quick amounts:</p>
              <div className="grid grid-cols-4 gap-3">
                {quickAmounts.map((quickAmount) => (
                  <button
                    key={quickAmount}
                    type="button"
                    onClick={() => setAmount(quickAmount.toString())}
                    className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  >
                    ${quickAmount.toLocaleString()}
                  </button>
                ))}
              </div>
            </div>

            {/* Submit Button */}
            <div>
              <button
                type="submit"
                disabled={loading || cards.length === 0}
                className="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
              >
                {loading ? 'Processing...' : `Send ${amount ? '$' + parseFloat(amount).toLocaleString() : ''}`}
              </button>
            </div>
          </form>

          {/* Transfer Tips */}
          <div className="mt-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <h3 className="text-sm font-semibold text-blue-900 mb-2">Transfer Tips:</h3>
            <ul className="text-sm text-blue-700 space-y-1">
              <li>• Double-check the destination card number</li>
              <li>• Make sure you have sufficient balance</li>
              <li>• Transfers are usually instant</li>
              <li>• Keep your receipt for reference</li>
            </ul>
          </div>

          {/* Recent Transfers */}
          <div className="mt-8 pt-6 border-t border-gray-200">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Transfers</h3>
            <div className="text-gray-500">
              <p>No recent transfers found</p>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}
