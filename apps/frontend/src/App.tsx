import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { useAuthStore } from './store/authStore';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Dashboard from './components/Dashboard';
import Topup from './components/Topup';
import Transfer from './components/Transfer';
import Withdraw from './components/Withdraw';
import Saldo from './components/Saldo';
import Transactions from './components/Transactions';
import Card from './components/Card';

function App() {
  const { isAuthenticated, loading, checkAuth } = useAuthStore();

  useEffect(() => {
    checkAuth();

    console.log(isAuthenticated)
  }, [checkAuth]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route 
            path="/login" 
            element={!isAuthenticated ? <Login /> : <Navigate to="/dashboard" />} 
          />
          <Route 
            path="/register" 
            element={!isAuthenticated ? <Register /> : <Navigate to="/dashboard" />} 
          />
          <Route 
            path="/dashboard" 
            element={isAuthenticated ? <Dashboard /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/topup" 
            element={isAuthenticated ? <Topup /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/transfer" 
            element={isAuthenticated ? <Transfer /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/withdraw" 
            element={isAuthenticated ? <Withdraw /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/saldo" 
            element={isAuthenticated ? <Saldo /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/transactions" 
            element={isAuthenticated ? <Transactions /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/cards" 
            element={isAuthenticated ? <Card /> : <Navigate to="/login" />} 
          />
          <Route path="/" element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
