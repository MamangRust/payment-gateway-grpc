import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi } from '../services/api';

interface User {
  id: string;
  username: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

interface AuthActions {
  login: (email: string, password: string) => Promise<void>;
  register: (firstname: string, lastname: string, email: string, password: string, confirm_password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState & AuthActions>()(
  persist(
    (set, get) => ({
      // State
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      loading: true,
      error: null,

      login: async (email: string, password: string) => {
        try {
          set({ loading: true, error: null });
          const response = await authApi.login({ email, password });

          const {access_token, refresh_token} = response;

          console.log("access", access_token)

          
          set({
            token: access_token,
            refreshToken: refresh_token,
            isAuthenticated: true,
            loading: false,
            error: null
          });
          
          // Verify the token and get user data
          await get().checkAuth();
        } catch (error: any) {
          set({
            loading: false,
            error: error.response?.data?.message || 'Login failed'
          });
          throw error;
        }
      },

      register: async (firstname: string, lastname: string, email: string, password: string, confirm_password: string) => {
        try {
          set({ loading: true, error: null });
          const response = await authApi.register({ firstname, lastname, email, password, confirm_password });
          
          set({
            token: response.access_token,
            refreshToken: response.refresh_token,
            isAuthenticated: true,
            loading: false,
            error: null
          });
          
          // Fetch user data after successful registration
          // await get().checkAuth();
        } catch (error: any) {
          set({
            loading: false,
            error: error.response?.data?.message || 'Registration failed'
          });
          throw error;
        }
      },

      logout: () => {
        set({
          user: null,
          token: null,
          refreshToken: null,
          isAuthenticated: false,
          error: null
        });
      },

      checkAuth: async () => {
        const { token } = get();

        console.log("token", token)

        if (!token) {
          set({ loading: false, isAuthenticated: false });
          return;
        }

        try {
          const userData = await authApi.getMe(token);
          console.log("user ", userData)

          set({
            user: userData,
            isAuthenticated: true,
            loading: false,
            error: null
          });
        } catch (error) {
          // Token expired or invalid
          set({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
            loading: false
          });
        }
      },

      clearError: () => set({ error: null })
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
        user: state.user
      })
    }
  )
);
