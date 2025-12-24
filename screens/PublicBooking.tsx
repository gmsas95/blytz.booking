import React, { useState, useMemo } from 'react';
import { ChevronLeft, Check, Lock, Clock } from 'lucide-react';
import { MOCK_SERVICES, MOCK_SLOTS } from '../constants';
import { Service, Slot, CustomerDetails, Business } from '../types';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Input } from '../components/Input';

interface PublicBookingProps {
  business: Business;
  onComplete: (bookingDetails: any) => void;
  onExit: () => void;
}

type Step = 'SERVICE' | 'SLOT' | 'DETAILS' | 'PAYMENT';

export const PublicBooking: React.FC<PublicBookingProps> = ({ business, onComplete, onExit }) => {
  const [step, setStep] = useState<Step>('SERVICE');
  const [selectedService, setSelectedService] = useState<Service | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [customer, setCustomer] = useState<CustomerDetails>({ name: '', email: '', phone: '' });
  const [isProcessing, setIsProcessing] = useState(false);

  // Filter data for this business
  const businessServices = useMemo(() => MOCK_SERVICES.filter(s => s.businessId === business.id), [business.id]);
  const businessSlots = useMemo(() => MOCK_SLOTS.filter(s => s.businessId === business.id), [business.id]);

  // Helper to format currency
  const fmtMoney = (n: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);
  
  // Helper to format date
  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { weekday: 'short', month: 'short', day: 'numeric' }).format(new Date(iso));
  const fmtTime = (iso: string) => new Intl.DateTimeFormat('en-US', { hour: 'numeric', minute: '2-digit' }).format(new Date(iso));

  const handleServiceSelect = (service: Service) => {
    setSelectedService(service);
    setStep('SLOT');
  };

  const handleSlotSelect = (slot: Slot) => {
    if (slot.isBooked) return;
    setSelectedSlot(slot);
  };

  const handleSlotConfirm = () => {
    if (selectedSlot) setStep('DETAILS');
  };

  const handleDetailsSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (customer.name && customer.email && customer.phone) {
      setStep('PAYMENT');
    }
  };

  const handlePayment = () => {
    setIsProcessing(true);
    // Simulate network request
    setTimeout(() => {
      setIsProcessing(false);
      onComplete({
        business,
        service: selectedService,
        slot: selectedSlot,
        customer,
      });
    }, 2000);
  };

  const goBack = () => {
    if (step === 'SERVICE') onExit();
    if (step === 'SLOT') setStep('SERVICE');
    if (step === 'DETAILS') setStep('SLOT');
    if (step === 'PAYMENT') setStep('DETAILS');
  };

  // Dynamic Theme Color (Tailwind doesn't support dynamic JIT class interpolation easily without safelist, 
  // so we'll use inline styles for the specific brand color)
  const brandStyle = { color: business.themeColor };

  return (
    <div className="w-full min-h-screen bg-white pb-24 font-sans">
      {/* Brand Header */}
      <div className="bg-white px-4 py-6 border-b border-gray-100 sticky top-0 z-10 flex items-center justify-between">
         <div>
            <button onClick={goBack} className="flex items-center text-sm text-gray-500 hover:text-gray-900 mb-1">
              <ChevronLeft className="h-4 w-4 mr-1" />
              {step === 'SERVICE' ? 'Back to Directory' : 'Back'}
            </button>
            <h2 className="text-lg font-bold text-gray-900">{business.name}</h2>
         </div>
         <div className="h-8 w-8 rounded-full bg-gray-900 text-white flex items-center justify-center text-xs font-bold">
            {business.name.substring(0,1)}
         </div>
      </div>

      <div className="max-w-lg mx-auto pt-6 px-4">
        <div className="mb-6">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            {step === 'SERVICE' && "Select a Service"}
            {step === 'SLOT' && "Pick a Time"}
            {step === 'DETAILS' && "Your Details"}
            {step === 'PAYMENT' && "Secure Deposit"}
          </h1>
          <p className="text-gray-500 text-sm mt-1">
            {step === 'SERVICE' && "Choose the package that fits your needs."}
            {step === 'SLOT' && "Select an available opening for your session."}
            {step === 'DETAILS' && "Where should we send the confirmation?"}
            {step === 'PAYMENT' && `A ${selectedService ? fmtMoney(selectedService.depositAmount) : ''} deposit is required to confirm.`}
          </p>
        </div>

        {/* STEP 1: SERVICE SELECTION */}
        {step === 'SERVICE' && (
          <div className="space-y-4">
            {businessServices.length === 0 ? (
               <p className="text-gray-500">No services available for this business.</p>
            ) : (
              businessServices.map((service) => (
                <Card 
                  key={service.id} 
                  onClick={() => handleServiceSelect(service)}
                  className="hover:scale-[1.01] transition-transform"
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-semibold text-gray-900">{service.name}</h3>
                      <p className="text-sm text-gray-500 mt-1 line-clamp-2">{service.description}</p>
                      <p className="text-xs text-gray-400 mt-2">{service.durationMin} minutes</p>
                    </div>
                    <div className="text-right">
                      <span className="block font-medium text-gray-900">{fmtMoney(service.totalPrice)}</span>
                      {service.depositAmount < service.totalPrice && (
                        <span className="block text-xs font-medium text-primary-600">
                          Deposit: {fmtMoney(service.depositAmount)}
                        </span>
                      )}
                    </div>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}

        {/* STEP 2: SLOT SELECTION */}
        {step === 'SLOT' && (
          <div className="h-full flex flex-col">
            {/* Horizontal Date Scroller */}
            <div className="flex space-x-2 overflow-x-auto pb-4 no-scrollbar">
              {[0, 1, 2, 3].map(offset => {
                const d = new Date();
                d.setDate(d.getDate() + offset);
                const isSelected = offset === (selectedSlot ? (new Date(selectedSlot.startTime).getDate() - new Date().getDate()) : -1);
                return (
                  <button 
                    key={offset}
                    className={`flex flex-col items-center justify-center min-w-[4rem] h-16 rounded-lg border text-sm transition-colors ${isSelected ? 'border-primary-500 bg-primary-50 text-primary-700' : 'border-gray-200 bg-white text-gray-600'}`}
                  >
                    <span className="text-xs uppercase font-medium">{d.toLocaleDateString('en-US', { weekday: 'short' })}</span>
                    <span className="text-lg font-bold">{d.getDate()}</span>
                  </button>
                )
              })}
            </div>

            <div className="space-y-3 mt-2">
              {businessSlots.length === 0 && <p className="text-gray-500 text-sm">No slots configured.</p>}
              {businessSlots.map((slot) => {
                return (
                  <button
                    key={slot.id}
                    disabled={slot.isBooked}
                    onClick={() => handleSlotSelect(slot)}
                    className={`w-full flex items-center justify-between p-4 rounded-lg border text-left transition-all ${
                      slot.isBooked 
                        ? 'bg-gray-100 border-gray-100 text-gray-400 cursor-not-allowed' 
                        : selectedSlot?.id === slot.id
                          ? 'border-primary-500 bg-primary-50 ring-1 ring-primary-500'
                          : 'bg-white border-gray-200 hover:border-primary-300'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <Clock className={`h-4 w-4 ${selectedSlot?.id === slot.id ? 'text-primary-600' : 'text-gray-400'}`} />
                      <span className={`font-medium ${selectedSlot?.id === slot.id ? 'text-primary-900' : 'text-gray-900'}`}>
                        {fmtTime(slot.startTime)} - {fmtTime(slot.endTime)}
                      </span>
                    </div>
                    {slot.isBooked && <span className="text-xs font-medium">Booked</span>}
                    {selectedSlot?.id === slot.id && <Check className="h-4 w-4 text-primary-600" />}
                  </button>
                );
              })}
            </div>

            <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 safe-area-pb z-20">
              <div className="max-w-lg mx-auto">
                <Button 
                  fullWidth 
                  disabled={!selectedSlot} 
                  onClick={handleSlotConfirm}
                >
                  Continue
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* STEP 3: CUSTOMER DETAILS */}
        {step === 'DETAILS' && (
          <form onSubmit={handleDetailsSubmit} className="space-y-4">
            <Input 
              label="Full Name" 
              placeholder="Jane Doe"
              required
              value={customer.name}
              onChange={e => setCustomer({...customer, name: e.target.value})}
            />
            <Input 
              label="Email Address" 
              type="email" 
              placeholder="jane@example.com"
              required
              value={customer.email}
              onChange={e => setCustomer({...customer, email: e.target.value})}
            />
            <Input 
              label="Phone Number" 
              type="tel" 
              placeholder="(555) 123-4567"
              required
              value={customer.phone}
              onChange={e => setCustomer({...customer, phone: e.target.value})}
            />
            
            <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 safe-area-pb z-20">
              <div className="max-w-lg mx-auto">
                <Button fullWidth type="submit">
                  Review & Pay
                </Button>
              </div>
            </div>
          </form>
        )}

        {/* STEP 4: PAYMENT */}
        {step === 'PAYMENT' && selectedService && selectedSlot && (
          <div className="space-y-6">
            {/* Summary Card */}
            <div className="bg-gray-50 rounded-lg p-4 border border-gray-200 space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Service</span>
                <span className="font-medium text-gray-900">{selectedService.name}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Date</span>
                <span className="font-medium text-gray-900">{fmtDate(selectedSlot.startTime)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Time</span>
                <span className="font-medium text-gray-900">{fmtTime(selectedSlot.startTime)}</span>
              </div>
              <div className="h-px bg-gray-200 my-2"></div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Total</span>
                <span className="text-gray-900">{fmtMoney(selectedService.totalPrice)}</span>
              </div>
              <div className="flex justify-between text-base font-semibold text-primary-700">
                <span>Due Now (Deposit)</span>
                <span>{fmtMoney(selectedService.depositAmount)}</span>
              </div>
              {selectedService.depositAmount < selectedService.totalPrice && (
                 <p className="text-xs text-gray-500 text-right mt-1">Remaining {fmtMoney(selectedService.totalPrice - selectedService.depositAmount)} due at appointment.</p>
              )}
            </div>

            {/* Fake Payment Form */}
            <div className="space-y-4">
               <div className="flex items-center gap-2 text-sm text-gray-600 mb-2">
                  <Lock className="h-3 w-3" />
                  <span>Secure SSL Encryption</span>
               </div>
               
               <Input label="Card Number" placeholder="4242 4242 4242 4242" />
               <div className="grid grid-cols-2 gap-4">
                  <Input label="Expiry" placeholder="MM/YY" />
                  <Input label="CVC" placeholder="123" />
               </div>
               <Input label="Cardholder Name" placeholder={customer.name} />
            </div>

            <div className="pt-4">
              <Button 
                fullWidth 
                onClick={handlePayment} 
                isLoading={isProcessing}
                variant="primary"
              >
                Pay {fmtMoney(selectedService.depositAmount)} & Confirm
              </Button>
              <div className="mt-6 flex justify-center">
                 <div className="text-[10px] text-gray-400 uppercase tracking-widest font-semibold flex items-center gap-1">
                    Powered by Blytz.Cloud
                 </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};