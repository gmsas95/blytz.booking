import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Mail, ArrowRight } from 'lucide-react';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { api } from '../api';
import { useAuth } from '../context/AuthContext';

export const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [isRegister, setIsRegister] = useState(false);
  const [name, setName] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      let response;
      if (isRegister) {
        response = await api.register({ email, name, password });
      } else {
        response = await api.login({ email, password });
      }

      api.setToken(response.token);
      login();
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-white">
      <div className="max-w-sm w-full space-y-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900">
            {isRegister ? 'Create Account' : 'Operator Login'}
          </h1>
          <p className="mt-2 text-sm text-gray-500">
            {isRegister ? 'Sign up to access your dashboard.' : 'Access your dashboard and bookings.'}
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm">
            {error}
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          {isRegister && (
            <Input 
              label="Full Name" 
              type="text" 
              required 
              placeholder="John Doe"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
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
            placeholder="••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <Button type="submit" fullWidth isLoading={loading} className="gap-2">
            {isRegister ? 'Create Account' : 'Login'} <ArrowRight className="h-4 w-4" />
          </Button>
        </form>
        
        <div className="text-center">
          <button
            type="button"
            onClick={() => {
              setIsRegister(!isRegister);
              setError('');
            }}
            className="text-sm text-gray-600 hover:text-gray-900 underline"
          >
            {isRegister ? 'Already have an account? Login' : "Don't have an account? Register"}
          </button>
        </div>
      </div>
    </div>
  );
};
