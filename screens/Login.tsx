import React, { useState } from 'react';
import { Mail, ArrowRight } from 'lucide-react';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { api } from '../api';

interface LoginProps {
  onLogin: () => void;
}

export const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const data = await api.login(email, password);
      localStorage.setItem('token', data.token);

      setLoading(false);
      onLogin();
    } catch (err: {
      console.error('Login failed:', err);
      setError('Invalid email or password. Please try again.');
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6 bg-gray-50">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-white">
      <div className="max-w-sm w-full space-y-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900">Operator Login</h1>
          <p className="mt-2 text-sm text-gray-500">Access your dashboard and bookings.</p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}

          <Input 
            label="Email Address" 
            type="email" 
            required 
            placeholder="you@service.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />

          <Input 
            label="Password" 
            type="password" 
            required 
            placeholder="•••••••••••••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <Button type="submit" fullWidth isLoading={loading} className="gap-2">
            Log In <ArrowRight className="h-4 w-4" />
          </Button>

          <div className="text-center">
            <button type="button" onClick={onLogin} className="text-xs text-gray-400 hover:text-gray-600 underline">
              Skip (Demo Mode)
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
