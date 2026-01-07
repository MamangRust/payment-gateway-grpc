import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { cardApi } from '../services/api';
import Layout from './Layout';

export default function Card() {
  const [cards, setCards] = useState<any[]>([]);
  const [showAddCardForm, setShowAddCardForm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const token = useAuthStore(state => state.token);
  const user = useAuthStore(state => state.user);

  const [formData, setFormData] = useState({
    card_number: '',
    card_type: 'credit',
    expire_date: '',
    cvv: '',
    card_provider: ''
  });

  useEffect(() => {
    fetchCards();
  }, [token, user]);

  const fetchCards = async () => {
    if (!token) return;

    try {
      const cardsData = await cardApi.getMyCards(token);
      console.log('Cards data from API:', cardsData);
      setCards(cardsData); // API now always returns array
    } catch (err: any) {
      console.error('Fetch cards error:', err);
      setError(err.response?.data?.message || 'Failed to fetch cards');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  // Helper function to convert MM/YY to ISO datetime format
  const convertExpireDateToISO = (expiryDate: string): string => {
    // Parse MM/YY format
    const parts = expiryDate.split('/');
    if (parts.length !== 2) {
      throw new Error('Invalid date format. Use MM/YY');
    }

    const month = parts[0];
    const year = parts[1];

    // Add 2000 to year to get full year (e.g., 25 -> 2025)
    const fullYear = `20${year}`;

    // Create date for the last day of the month at midnight UTC
    const lastDayOfMonth = new Date(parseInt(fullYear), parseInt(month), 0);

    // Format as ISO datetime (end of the month at 00:00:00Z)
    const isoString = `${fullYear}-${month.padStart(2, '0')}-${lastDayOfMonth.getDate().toString().padStart(2, '0')}T00:00:00Z`;

    return isoString;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !user) return;

    setLoading(true);
    setError('');

    try {
      const cardData = {
        ...formData,
        expire_date: convertExpireDateToISO(formData.expire_date), // Convert to ISO format
        user_id: parseInt(user.id) // Convert to number if needed
      };

      console.log('Card data being sent:', cardData);

      const createResponse = await cardApi.createCard(cardData, token);
      console.log('Card created:', createResponse);

      // Refresh cards list after creation
      await fetchCards();
      setShowAddCardForm(false);
      setFormData({
        card_number: '',
        card_type: 'credit',
        expire_date: '',
        cvv: '',
        card_provider: ''
      });
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create card');
    } finally {
      setLoading(false);
    }
  };

  const formatCardNumber = (value: string) => {
    const v = value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
    const matches = v.match(/\d{4,16}/g);
    const match = matches && matches[0] || '';
    const parts = [];
    for (let i = 0, len = match.length; i < len; i += 4) {
      parts.push(match.substring(i, i + 4));
    }
    if (parts.length) {
      return parts.join(' ');
    } else {
      return v;
    }
  };

  return (
    <Layout>
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">My Cards</h1>
          <p className="text-gray-600">Manage your payment cards</p>
        </div>

        <div className="mb-6">
          <button
            onClick={() => setShowAddCardForm(!showAddCardForm)}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition duration-200"
          >
            {showAddCardForm ? 'Cancel' : 'Add New Card'}
          </button>
        </div>

        {error && (
          <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded-md mb-6">
            <p className="text-red-700">{error}</p>
          </div>
        )}

        {showAddCardForm && (
          <div className="bg-white shadow-lg rounded-xl p-6 mb-8">
            <h2 className="text-xl font-semibold mb-6">Add New Card</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="card_number" className="block text-sm font-medium text-gray-700 mb-2">
                    Card Number
                  </label>
                  <input
                    id="card_number"
                    name="card_number"
                    type="text"
                    required
                    maxLength={19}
                    className="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="1234 5678 9012 3456"
                    value={formatCardNumber(formData.card_number)}
                    onChange={(e) => setFormData({
                      ...formData,
                      card_number: e.target.value.replace(/\s/g, '')
                    })}
                  />
                </div>

                <div>
                  <label htmlFor="cvv" className="block text-sm font-medium text-gray-700 mb-2">
                    CVV
                  </label>
                  <input
                    id="cvv"
                    name="cvv"
                    type="text"
                    required
                    maxLength={4}
                    className="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="123"
                    value={formData.cvv}
                    onChange={(e) => setFormData({
                      ...formData,
                      cvv: e.target.value.replace(/[^0-9]/g, '')
                    })}
                  />
                </div>

                <div>
                  <label htmlFor="expire_date" className="block text-sm font-medium text-gray-700 mb-2">
                    Expiry Date (MM/YY)
                  </label>
                  <input
                    id="expire_date"
                    name="expire_date"
                    type="text"
                    required
                    placeholder="12/25"
                    className="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    value={formData.expire_date}
                    onChange={(e) => {
                      const value = e.target.value.replace(/[^0-9/]/g, '');
                      if (value.length <= 5) {
                        setFormData({
                          ...formData,
                          expire_date: value
                        });
                      }
                    }}
                  />
                  {formData.expire_date && formData.expire_date.length === 5 && (
                    <p className="mt-1 text-xs text-gray-500">
                      Will be sent as: {(() => {
                        try {
                          return convertExpireDateToISO(formData.expire_date);
                        } catch {
                          return 'Invalid format';
                        }
                      })()}
                    </p>
                  )}
                </div>

                <div>
                  <label htmlFor="card_type" className="block text-sm font-medium text-gray-700 mb-2">
                    Card Type
                  </label>
                  <select
                    id="card_type"
                    name="card_type"
                    required
                    className="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    value={formData.card_type}
                    onChange={handleChange}
                  >
                    <option value="credit">Credit Card</option>
                    <option value="debit">Debit Card</option>
                  </select>
                </div>

                <div className="md:col-span-2">
                  <label htmlFor="card_provider" className="block text-sm font-medium text-gray-700 mb-2">
                    Card Provider
                  </label>
                  <input
                    id="card_provider"
                    name="card_provider"
                    type="text"
                    required
                    className="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Visa, Mastercard, etc."
                    value={formData.card_provider}
                    onChange={handleChange}
                  />
                </div>
              </div>

              <div className="flex justify-end space-x-4">
                <button
                  type="button"
                  onClick={() => setShowAddCardForm(false)}
                  className="px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition duration-200"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition duration-200"
                >
                  {loading ? 'Adding...' : 'Add Card'}
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {cards.length > 0 ? (
            cards.map((card: any) => (
              <div key={card.card_id} className="bg-gradient-to-br from-blue-600 to-blue-700 text-white rounded-xl p-6 shadow-lg">
                <div className="flex justify-between items-start mb-8">
                  <div className="text-sm opacity-90">
                    {card.card_type === 'credit' ? 'Credit Card' : 'Debit Card'}
                  </div>
                  <div className="text-sm">
                    {card.card_provider}
                  </div>
                </div>

                <div className="mb-6">
                  <div className="text-xs opacity-75 mb-2">Card Number</div>
                  <div className="text-lg font-mono">
                    {card.card_number.replace(/(\d{4})(?=\d)/g, '$1 ')}
                  </div>
                </div>

                <div className="flex justify-between items-end">
                  <div>
                    <div className="text-xs opacity-75 mb-1">Expires</div>
                    <div className="text-sm">{card.expire_date}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs opacity-75 mb-1">CVV</div>
                    <div className="text-sm">***</div>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="col-span-full text-center text-gray-500">
              No cards found
            </div>
          )}
        </div>

        {cards.length === 0 && !showAddCardForm && (
          <div className="text-center py-12 bg-white rounded-xl shadow-lg">
            <div className="text-gray-400 mb-4">
              <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No cards yet</h3>
            <p className="text-gray-600 mb-6">Add your first payment card to get started</p>
            <button
              onClick={() => setShowAddCardForm(true)}
              className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition duration-200"
            >
              Add Your First Card
            </button>
          </div>
        )}

        <div className="mt-8">
          <Link
            to="/dashboard"
            className="text-blue-600 hover:text-blue-500 font-medium"
          >
            ← Back to Dashboard
          </Link>
        </div>
      </div>
    </Layout>
  );
}
