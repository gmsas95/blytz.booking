import React from 'react';
import { useNavigate } from 'react-router-dom';
import { CheckCircle, Download, Home } from 'lucide-react';
import { Button } from '../components/Button';
import { Card } from '../components/Card';

export const Confirmation: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-6 text-center">
      <div className="w-full max-w-md space-y-8 animate-fade-in-up">
        <div className="flex justify-center">
          <CheckCircle className="h-20 w-20 text-green-500" />
        </div>

        <div>
            <h1 className="text-2xl font-bold text-gray-900">Booking Confirmed!</h1>
            <p className="text-gray-500 mt-2">
                Your booking has been confirmed.
            </p>
        </div>

        <Card className="text-left space-y-4 bg-gray-50 border-gray-200">
            <div className="pt-4 border-t border-gray-200">
                <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-600">Status</span>
                    <span className="font-bold text-green-600 text-lg">
                        Confirmed
                    </span>
                </div>
            </div>
        </Card>

        <div className="space-y-3">
            <Button variant="outline" fullWidth className="gap-2" onClick={() => alert("Downloading receipt PDF...")}>
                <Download className="h-4 w-4" /> Download Receipt
            </Button>
            <Button variant="ghost" fullWidth onClick={() => navigate('/')}>
                Book Another Service
            </Button>
        </div>
      </div>
    </div>
  );
};
