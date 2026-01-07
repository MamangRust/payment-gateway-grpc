import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { cardApi, topupApi } from '../services/api';
import { useAuthStore } from '../store/authStore';
import Layout from './Layout';
import type { Card } from '../types/api';

export default function Topup() {
  const [cards, setCards] = useState<Card[]>([]);
  const [selectedCard, setSelectedCard] = useState('');
  const [amount, setAmount] = useState('');
  const [method, setMethod] = useState('Alfamart');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const token = useAuthStore(state => state.token);

  useEffect(() => {
    const fetchCards = async () => {
      if (!token) return;

      try {
        const cardsData = await cardApi.getMyCards(token);
        console.log('Topup cards:', cardsData);
        setCards(cardsData);
        if (cardsData.length > 0) {
          setSelectedCard(cardsData[0].card_number);
        }
      } catch (err: any) {
        console.error('Topup fetch cards error:', err);
        setError('Failed to fetch cards');
      }
    };

    fetchCards();
  }, [token]);

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

    try {
      const topupData = {
        card_number: selectedCard,
        topup_amount: parseInt(amount),
        topup_method: method
      };

      if (!token) {
        throw new Error('No authentication token');
      }

      await topupApi.createTopup(topupData, token);
      setSuccess('Topup successful!');
      setAmount('');
    } catch (err: any) {
      console.error('Topup error:', err);
      setError(err.response?.data?.message || 'Topup failed');
    } finally {
      setLoading(false);
    }
  };

  const quickAmounts = [10000, 25000, 50000, 100000];

  return (
    <Layout>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center py-6">
            <Link to="/dashboard" className="text-indigo-600 hover:text-indigo-800 mr-4">
              ← Back
            </Link>
            <h1 className="text-2xl font-bold text-gray-900">Top Up</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Add funds to your card</h2>

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
                <select
                  value={selectedCard}
                  onChange={(e) => setSelectedCard(e.target.value)}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  {cards.map((card) => (
                    <option key={card.card_id} value={card.card_number}>
                      {card.card_provider} - {card.card_type} (•••• {card.card_number.slice(-4)})
                    </option>
                  ))}
                </select>
              )}
            </div>

            {/* Amount Input */}
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
                />
              </div>
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

            {/* Payment Method */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Payment Method
              </label>
              <select
                value={method}
                onChange={(e) => setMethod(e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="">-- Pilih Metode Pembayaran --</option>
                <option value="Alfamart">Alfamart</option>
                <option value="Indomart">Indomart</option>
                <option value="Lawson">Lawson</option>
                <option value="Dana">Dana</option>
                <option value="OVO">OVO</option>
                <option value="Gopay">Gopay</option>
                <option value="LinkAja">LinkAja</option>
                <option value="Jenius">Jenius</option>
                <option value="Fastpay">Fastpay</option>
                <option value="Kudo">Kudo</option>
                <option value="BRI">BRI</option>
                <option value="Mandiri">Mandiri</option>
                <option value="BCA">BCA</option>
                <option value="BNI">BNI</option>
                <option value="Bukopin">Bukopin</option>
                <option value="E-Banking">E-Banking</option>
                <option value="Visa">Visa</option>
                <option value="MasterCard">MasterCard</option>
                <option value="Discover">Discover</option>
                <option value="American Express">American Express</option>
                <option value="PayPal">PayPal</option>
              </select>

            </div>

            {/* Submit Button */}
            <div>
              <button
                type="submit"
                disabled={loading || cards.length === 0}
                className="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
              >
                {loading ? 'Processing...' : `Top Up ${amount ? '$' + parseFloat(amount).toLocaleString() : ''}`}
              </button>
            </div>
          </form>

          {/* Topup History */}
          <div className="mt-8 pt-6 border-t border-gray-200">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Topups</h3>
            <div className="text-gray-500">
              <p>No recent topups found</p>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}
