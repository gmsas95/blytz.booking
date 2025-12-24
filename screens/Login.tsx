import React, { useState } from 'react';
import { Mail, ArrowRight } from 'lucide-react';
import { Button } from '../components/Button';
import { Input } from '../components/Input';

interface LoginProps {
  onLogin: () => void;
}

export const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const [email, setEmail] = useState('');
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    // Simulate API call
    setTimeout(() => {
      setLoading(false);
      setSent(true);
      // Auto login for demo after a short delay
      setTimeout(onLogin, 1500);
    }, 1000);
  };

  if (sent) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6 bg-gray-50">
        <div className="max-w-md w-full text-center space-y-4">
           <div className="bg-green-100 rounded-full h-16 w-16 flex items-center justify-center mx-auto text-green-600">
             <Mail className="h-8 w-8" />
           </div>
           <h2 className="text-xl font-bold text-gray-900">Magic Link Sent!</h2>
           <p className="text-gray-500">Check your inbox at <strong>{email}</strong> to log in.</p>
           <p className="text-xs text-gray-400 mt-4">(Redirecting to dashboard...)</p>
        </div>
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
          <Input 
            label="Email Address" 
            type="email" 
            required 
            placeholder="you@service.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />

          <Button type="submit" fullWidth isLoading={loading} className="gap-2">
            Send Magic Link <ArrowRight className="h-4 w-4" />
          </Button>
        </form>
        
        <div className="text-center">
             <button type="button" onClick={onLogin} className="text-xs text-gray-400 hover:text-gray-600 underline">
                 Skip (Demo Mode)
             </button>
        </div>
      </div>
    </div>
  );
};
