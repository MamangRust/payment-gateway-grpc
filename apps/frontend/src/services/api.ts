import axios from 'axios';
import type { AuthResponse, LoginRequest, RegisterRequest } from '../types/api';

const API_BASE_URL = 'http://localhost:5000/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper function to create API instance with auth header
const createApiWithAuth = (token: string | null) => {
  return axios.create({
    baseURL: API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` })
    },
  });
};

export const authApi = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', data);
    return response.data.data;
  },
  
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },
  
  getMe: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/auth/me');
    return response.data.data;
  },

  refreshToken: async (refreshToken: string) => {
    if (!refreshToken) {
      throw new Error('No refresh token');
    }
    const response = await api.post('/auth/refresh-token', {
      refresh_token: refreshToken
    });
    return response.data;
  }
};

export const cardApi = {
  getCards: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/card');
    return response.data;
  },
  
  getMyCards: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/card/user');

    console.log('Card response:', response);

    const cardData = response.data.data;
    
    // Always return an array - if it's an object, convert to array
    if (Array.isArray(cardData)) {
      return cardData;
    } else if (cardData && typeof cardData === 'object') {
      return [cardData]; // Convert single object to array
    } else {
      return []; // Return empty array if null/undefined
    }
  },

  createCard: async (data: {
    card_number: string;
    card_type: string;
    expire_date: string;
    cvv: string;
    card_provider: string;
    user_id: number;
  }, token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.post('/card/create', data);
    return response.data;
  }
};

export const saldoApi = {
  getMySaldo: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/saldo/user');
    return response.data;
  }
};

export const topupApi = {
  createTopup: async (data: any, token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.post('/topups/create', data);
    return response.data;
  },
  
  getTopupHistory: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/topups');
    return response.data;
  }
};

export const transferApi = {
  createTransfer: async (data: any, token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.post('/transfer/create', data);
    return response.data;
  },
  
  getTransferHistory: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/transfer');
    return response.data;
  }
};

export const withdrawApi = {
  createWithdraw: async (data: any, token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.post('/withdraw/create', data);
    return response.data;
  },
  
  getWithdrawHistory: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/withdraw');
    return response.data;
  }
};

export const transactionApi = {
  getTransactions: async (token: string) => {
    const authApi = createApiWithAuth(token);
    const response = await authApi.get('/transaction');
    return response.data;
  }
};
