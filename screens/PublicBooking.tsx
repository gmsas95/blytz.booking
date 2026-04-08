import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ChevronLeft, Check, Lock, Clock, AlertCircle } from 'lucide-react';
import { api, ApiError, Service, Slot, CustomerDetails, Business } from '../api';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Input } from '../components/Input';
import { formatMoneyFromMinor, subtractMinorAmounts } from '../utils/money';

type Step = 'SERVICE' | 'SLOT' | 'DETAILS' | 'PAYMENT';

export const PublicBooking: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const [business, setBusiness] = useState<Business | null>(null);
  const [loadingBusiness, setLoadingBusiness] = useState(true);
  const [businessError, setBusinessError] = useState<string | null>(null);
  const [step, setStep] = useState<Step>('SERVICE');
  const [selectedService, setSelectedService] = useState<Service | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [customer, setCustomer] = useState<CustomerDetails>({ name: '', email: '', phone: '' });
  const [isProcessing, setIsProcessing] = useState(false);
  const [services, setServices] = useState<Service[]>([]);
  const [slots, setSlots] = useState<Slot[]>([]);
  const [loadingServices, setLoadingServices] = useState(true);
  const [loadingSlots, setLoadingSlots] = useState(true);
  const [dataError, setDataError] = useState<string | null>(null);
  const [bookingError, setBookingError] = useState<string | null>(null);

  useEffect(() => {
    const fetchBusiness = async () => {
      if (!slug) {
        setBusiness(null);
        setBusinessError('Missing workshop slug.');
        setLoadingBusiness(false);
        return;
      }

      try {
        setBusiness(null);
        setServices([]);
        setSlots([]);
        setSelectedService(null);
        setSelectedSlot(null);
        setDataError(null);
        setBookingError(null);
        setBusinessError(null);
        const businesses = await api.getBusinesses();
        const found = businesses.find((candidate) => candidate.slug === slug);
        if (!found) {
          setBusiness(null);
          setBusinessError('Workshop not found.');
          return;
        }
        setBusiness(found);
      } catch (error) {
        console.error('Failed to fetch business:', error);
        setBusiness(null);
        setBusinessError('Failed to load workshop details. Please try again.');
      } finally {
        setLoadingBusiness(false);
      }
    };

    fetchBusiness();
  }, [slug]);

  useEffect(() => {
    const fetchData = async () => {
      if (!business) {
        setLoadingServices(false);
        setLoadingSlots(false);
        return;
      }

      try {
        setLoadingServices(true);
        setLoadingSlots(true);
        setDataError(null);

        const [servicesData, slotsData] = await Promise.all([
          api.getServicesByBusiness(business.id),
          api.getSlotsByBusiness(business.id),
        ]);

        setServices(servicesData);
        setSlots(slotsData);
      } catch (error) {
        console.error('Failed to fetch booking data:', error);
        setDataError('Failed to load services or slots. Please refresh and try again.');
      } finally {
        setLoadingServices(false);
        setLoadingSlots(false);
      }
    };

    fetchData();
  }, [business]);

  const slotDayOptions = useMemo(() => {
    return [0, 1, 2, 3].map((offset) => {
      const date = new Date();
      date.setDate(date.getDate() + offset);
      return {
        key: offset,
        date,
      };
    });
  }, []);

  const fmtMoney = (amountMinor: number, currencyCode: string = 'USD') => formatMoneyFromMinor(amountMinor, currencyCode);
  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { weekday: 'short', month: 'short', day: 'numeric' }).format(new Date(iso));
  const fmtTime = (iso: string) => new Intl.DateTimeFormat('en-US', { hour: 'numeric', minute: '2-digit' }).format(new Date(iso));

  const handleServiceSelect = (service: Service) => {
    setSelectedService(service);
    setSelectedSlot(null);
    setBookingError(null);
    setStep('SLOT');
  };

  const handleSlotSelect = (slot: Slot) => {
    if (slot.is_booked) return;
    setSelectedSlot(slot);
    setBookingError(null);
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

  const refreshSlots = async () => {
    if (!business) return;
    setLoadingSlots(true);
    try {
      const nextSlots = await api.getSlotsByBusiness(business.id);
      setSlots(nextSlots);
      setSelectedSlot(null);
    } catch (error) {
      console.error('Failed to refresh slots:', error);
    } finally {
      setLoadingSlots(false);
    }
  };

  const handlePayment = async () => {
    if (!business || !selectedService || !selectedSlot) return;

    setIsProcessing(true);
    setBookingError(null);
    try {
      await api.createBooking({
        business_id: business.id,
        service_id: selectedService.id,
        slot_id: selectedSlot.id,
        customer,
      });

      navigate('/confirmation');
    } catch (error) {
      console.error('Failed to create booking:', error);
      if (error instanceof ApiError && error.status === 409) {
        setBookingError('That slot was just taken. Please choose another available time.');
        setStep('SLOT');
        await refreshSlots();
      } else {
        setBookingError('Failed to create booking. Please try again.');
      }
    } finally {
      setIsProcessing(false);
    }
  };

  const goBack = () => {
    setBookingError(null);
    if (step === 'SERVICE') navigate('/');
    if (step === 'SLOT') setStep('SERVICE');
    if (step === 'DETAILS') setStep('SLOT');
    if (step === 'PAYMENT') setStep('DETAILS');
  };

  if (loadingBusiness) {
    return <div className="min-h-screen bg-white flex items-center justify-center text-gray-500">Loading workshop...</div>;
  }

  if (!business) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center px-4">
        <Card className="max-w-md w-full text-center space-y-4">
          <div className="flex justify-center text-red-500"><AlertCircle className="h-6 w-6" /></div>
          <h1 className="text-xl font-semibold text-gray-900">Workshop unavailable</h1>
          <p className="text-sm text-gray-500">{businessError || 'Unable to find that workshop.'}</p>
          <Button onClick={() => navigate('/')}>Back to home</Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="w-full min-h-screen bg-white pb-24 font-sans">
      <div className="bg-white px-4 py-6 border-b border-gray-100 sticky top-0 z-10 flex items-center justify-between">
        <div>
          <button onClick={goBack} className="flex items-center text-sm text-gray-500 hover:text-gray-900 mb-1">
            <ChevronLeft className="h-4 w-4 mr-1" />
            {step === 'SERVICE' ? 'Back to Directory' : 'Back'}
          </button>
          <h2 className="text-lg font-bold text-gray-900">{business.name}</h2>
        </div>
        <div className="h-8 w-8 rounded-full bg-gray-900 text-white flex items-center justify-center text-xs font-bold">
          {business.name.substring(0, 1)}
        </div>
      </div>

      <div className="max-w-lg mx-auto pt-6 px-4">
        <div className="mb-6">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            {step === 'SERVICE' && 'Select a Service'}
            {step === 'SLOT' && 'Pick a Time'}
            {step === 'DETAILS' && 'Your Details'}
            {step === 'PAYMENT' && 'Secure Deposit'}
          </h1>
          <p className="text-gray-500 text-sm mt-1">
            {step === 'SERVICE' && 'Choose the package that fits your needs.'}
            {step === 'SLOT' && 'Select an available opening for your session.'}
            {step === 'DETAILS' && 'Where should we send the confirmation?'}
            {step === 'PAYMENT' && selectedService && `A ${fmtMoney(selectedService.deposit_amount_minor, selectedService.currency_code)} deposit is required to confirm.`}
          </p>
        </div>

        {(dataError || bookingError) && (
          <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {bookingError || dataError}
          </div>
        )}

        {step === 'SERVICE' && (
          <div className="space-y-4">
            {loadingServices ? (
              <p className="text-gray-500">Loading services...</p>
            ) : services.length === 0 ? (
              <p className="text-gray-500">No services available for this workshop.</p>
            ) : (
              services.map((service) => (
                <Card
                  key={service.id}
                  onClick={() => handleServiceSelect(service)}
                  className="hover:scale-[1.01] transition-transform"
                >
                  <div className="flex justify-between items-start gap-4">
                    <div>
                      <h3 className="font-semibold text-gray-900">{service.name}</h3>
                      <p className="text-sm text-gray-500 mt-1 line-clamp-2">{service.description}</p>
                      <p className="text-xs text-gray-400 mt-2">{service.duration_min} minutes</p>
                    </div>
                    <div className="text-right">
                      <span className="block font-medium text-gray-900">{fmtMoney(service.total_price_minor, service.currency_code)}</span>
                      {service.deposit_amount_minor < service.total_price_minor && (
                        <span className="block text-xs font-medium text-primary-600">
                          Deposit: {fmtMoney(service.deposit_amount_minor, service.currency_code)}
                        </span>
                      )}
                    </div>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}

        {step === 'SLOT' && (
          <div className="h-full flex flex-col">
            <div className="flex space-x-2 overflow-x-auto pb-4 no-scrollbar">
              {slotDayOptions.map(({ key, date }) => {
                const isSelected = selectedSlot
                  ? new Date(selectedSlot.start_time).toDateString() === date.toDateString()
                  : false;

                return (
                  <button
                    key={key}
                    className={`flex flex-col items-center justify-center min-w-[4rem] h-16 rounded-lg border text-sm transition-colors ${isSelected ? 'border-primary-500 bg-primary-50 text-primary-700' : 'border-gray-200 bg-white text-gray-600'}`}
                  >
                    <span className="text-xs uppercase font-medium">{date.toLocaleDateString('en-US', { weekday: 'short' })}</span>
                    <span className="text-lg font-bold">{date.getDate()}</span>
                  </button>
                );
              })}
            </div>

            <div className="space-y-3 mt-2">
              {loadingSlots ? (
                <p className="text-gray-500 text-sm">Loading slots...</p>
              ) : slots.length === 0 ? (
                <p className="text-gray-500 text-sm">No slots configured.</p>
              ) : (
                slots.map((slot) => (
                  <button
                    key={slot.id}
                    disabled={slot.is_booked}
                    onClick={() => handleSlotSelect(slot)}
                    className={`w-full flex items-center justify-between p-4 rounded-lg border text-left transition-all ${
                      slot.is_booked
                        ? 'bg-gray-100 border-gray-100 text-gray-400 cursor-not-allowed'
                        : selectedSlot?.id === slot.id
                          ? 'border-primary-500 bg-primary-50 ring-1 ring-primary-500'
                          : 'bg-white border-gray-200 hover:border-primary-300'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <Clock className={`h-4 w-4 ${selectedSlot?.id === slot.id ? 'text-primary-600' : 'text-gray-400'}`} />
                      <span className={`font-medium ${selectedSlot?.id === slot.id ? 'text-primary-900' : 'text-gray-900'}`}>
                        {fmtTime(slot.start_time)} - {fmtTime(slot.end_time)}
                      </span>
                    </div>
                    {slot.is_booked && <span className="text-xs font-medium">Booked</span>}
                    {selectedSlot?.id === slot.id && <Check className="h-4 w-4 text-primary-600" />}
                  </button>
                ))
              )}
            </div>

            <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 safe-area-pb z-20">
              <div className="max-w-lg mx-auto">
                <Button fullWidth disabled={!selectedSlot} onClick={handleSlotConfirm}>Continue</Button>
              </div>
            </div>
          </div>
        )}

        {step === 'DETAILS' && (
          <form onSubmit={handleDetailsSubmit} className="space-y-4">
            <Input
              label="Full Name"
              placeholder="Jane Doe"
              required
              value={customer.name}
              onChange={(e) => setCustomer({ ...customer, name: e.target.value })}
            />
            <Input
              label="Email Address"
              type="email"
              placeholder="jane@example.com"
              required
              value={customer.email}
              onChange={(e) => setCustomer({ ...customer, email: e.target.value })}
            />
            <Input
              label="Phone Number"
              type="tel"
              placeholder="(555) 123-4567"
              required
              value={customer.phone}
              onChange={(e) => setCustomer({ ...customer, phone: e.target.value })}
            />

            <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 safe-area-pb z-20">
              <div className="max-w-lg mx-auto">
                <Button fullWidth type="submit">Review & Pay</Button>
              </div>
            </div>
          </form>
        )}

        {step === 'PAYMENT' && selectedService && selectedSlot && (
          <div className="space-y-6">
            <div className="bg-gray-50 rounded-lg p-4 border border-gray-200 space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Service</span>
                <span className="font-medium text-gray-900">{selectedService.name}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Date</span>
                <span className="font-medium text-gray-900">{fmtDate(selectedSlot.start_time)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Time</span>
                <span className="font-medium text-gray-900">{fmtTime(selectedSlot.start_time)}</span>
              </div>
              <div className="h-px bg-gray-200 my-2"></div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Total</span>
                <span className="text-gray-900">{fmtMoney(selectedService.total_price_minor, selectedService.currency_code)}</span>
              </div>
              <div className="flex justify-between text-base font-semibold text-primary-700">
                <span>Due Now (Deposit)</span>
                <span>{fmtMoney(selectedService.deposit_amount_minor, selectedService.currency_code)}</span>
              </div>
              {selectedService.deposit_amount_minor < selectedService.total_price_minor && (
                <p className="text-xs text-gray-500 text-right mt-1">
                  Remaining {fmtMoney(subtractMinorAmounts(selectedService.total_price_minor, selectedService.deposit_amount_minor), selectedService.currency_code)} due at appointment.
                </p>
              )}
            </div>

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
                Pay {fmtMoney(selectedService.deposit_amount_minor, selectedService.currency_code)} & Confirm
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
