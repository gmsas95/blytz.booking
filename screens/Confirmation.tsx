import React from 'react';
import { CheckCircle, Download, Home } from 'lucide-react';
import { Button } from '../components/Button';
import { Card } from '../components/Card';

interface ConfirmationProps {
  bookingDetails: any; // Using any for simplicity in demo
  onDone: () => void;
}

export const Confirmation: React.FC<ConfirmationProps> = ({ bookingDetails, onDone }) => {
  if (!bookingDetails) return null;

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-6 text-center">
      <div className="w-full max-w-md space-y-8 animate-fade-in-up">
        <div className="flex justify-center">
          <CheckCircle className="h-20 w-20 text-green-500" />
        </div>
        
        <div>
            <h1 className="text-2xl font-bold text-gray-900">Booking Confirmed!</h1>
            <p className="text-gray-500 mt-2">
                A receipt has been sent to <span className="font-medium text-gray-900">{bookingDetails.customer.email}</span>.
            </p>
        </div>

        <Card className="text-left space-y-4 bg-gray-50 border-gray-200">
            <div>
                <p className="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-1">Service</p>
                <p className="font-medium text-gray-900">{bookingDetails.service.name}</p>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
                <div>
                    <p className="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-1">Date</p>
                    <p className="font-medium text-gray-900">{new Date(bookingDetails.slot.startTime).toLocaleDateString()}</p>
                </div>
                <div>
                    <p className="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-1">Time</p>
                    <p className="font-medium text-gray-900">
                        {new Date(bookingDetails.slot.startTime).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
                    </p>
                </div>
            </div>

            <div className="pt-4 border-t border-gray-200">
                <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-600">Deposit Paid</span>
                    <span className="font-bold text-gray-900 text-lg">
                        ${bookingDetails.service.depositAmount}
                    </span>
                </div>
            </div>
        </Card>

        <div className="space-y-3">
            <Button variant="outline" fullWidth className="gap-2" onClick={() => alert("Downloading receipt PDF...")}>
                <Download className="h-4 w-4" /> Download Receipt
            </Button>
            <Button variant="ghost" fullWidth onClick={onDone}>
                Book Another Service
            </Button>
        </div>
      </div>
    </div>
  );
};
