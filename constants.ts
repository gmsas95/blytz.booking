import { Business, Service, Slot, Booking, BookingStatus } from './types';

export const MOCK_BUSINESSES: Business[] = [
  {
    id: 'b1',
    name: 'DetailPro Automotive',
    slug: 'detail-pro',
    vertical: 'Automotive',
    description: 'Premium mobile detailing and ceramic coating.',
    themeColor: 'blue',
  },
  {
    id: 'b2',
    name: 'Lumina Wellness Spa',
    slug: 'lumina-spa',
    vertical: 'Wellness',
    description: 'Massage therapy, facials, and holistic healing.',
    themeColor: 'emerald',
  },
  {
    id: 'b3',
    name: 'FlashFrame Studios',
    slug: 'flash-frame',
    vertical: 'Creative',
    description: 'Editorial portraiture and high-end fashion photography.',
    themeColor: 'zinc',
  }
];

export const MOCK_SERVICES: Service[] = [
  // b1: Automotive
  {
    id: 's1',
    businessId: 'b1',
    name: 'Full Interior Detail',
    description: 'Deep clean, steam, shampoo, and leather conditioning.',
    durationMin: 120,
    totalPrice: 200,
    depositAmount: 50,
  },
  {
    id: 's2',
    businessId: 'b1',
    name: 'Ceramic Coating Gold',
    description: '5-year protection package with paint correction.',
    durationMin: 360,
    totalPrice: 1200,
    depositAmount: 300,
  },
  // b2: Wellness
  {
    id: 's3',
    businessId: 'b2',
    name: 'Deep Tissue Massage',
    description: '60-minute therapeutic massage for stress relief.',
    durationMin: 60,
    totalPrice: 120,
    depositAmount: 40,
  },
  {
    id: 's4',
    businessId: 'b2',
    name: 'Hydrafacial Signature',
    description: 'Cleanse, extract, and hydrate skin.',
    durationMin: 45,
    totalPrice: 180,
    depositAmount: 60,
  },
  // b3: Creative (Photography)
  {
    id: 's5',
    businessId: 'b3',
    name: 'Editorial Portrait Session',
    description: '2-hour studio session, 3 outfit changes, 10 retouched edits.',
    durationMin: 120,
    totalPrice: 450,
    depositAmount: 150,
  },
  {
    id: 's6',
    businessId: 'b3',
    name: 'Mini Headshot Blitzer',
    description: '30-minute express session. Perfect for LinkedIn.',
    durationMin: 30,
    totalPrice: 150,
    depositAmount: 150, // Full payment
  },
];

// Helper to generate slots
const today = new Date();
const tomorrow = new Date(today);
tomorrow.setDate(tomorrow.getDate() + 1);

const formatSlot = (date: Date, hour: number, minute: number = 0) => {
  const d = new Date(date);
  d.setHours(hour, minute, 0, 0);
  return d.toISOString();
};

// Generate slots for each business
export const MOCK_SLOTS: Slot[] = [
  // b1 slots
  { id: 'sl1', businessId: 'b1', startTime: formatSlot(today, 9), endTime: formatSlot(today, 11), isBooked: true },
  { id: 'sl2', businessId: 'b1', startTime: formatSlot(today, 13), endTime: formatSlot(today, 15), isBooked: false },
  { id: 'sl3', businessId: 'b1', startTime: formatSlot(tomorrow, 9), endTime: formatSlot(tomorrow, 11), isBooked: false },
  
  // b2 slots
  { id: 'sl4', businessId: 'b2', startTime: formatSlot(today, 10), endTime: formatSlot(today, 11), isBooked: false },
  { id: 'sl5', businessId: 'b2', startTime: formatSlot(today, 14), endTime: formatSlot(today, 15), isBooked: false },
  
  // b3 slots
  { id: 'sl6', businessId: 'b3', startTime: formatSlot(today, 10), endTime: formatSlot(today, 10, 30), isBooked: false },
  { id: 'sl7', businessId: 'b3', startTime: formatSlot(tomorrow, 11), endTime: formatSlot(tomorrow, 11, 30), isBooked: false },
];

export const MOCK_BOOKINGS: Booking[] = [
  {
    id: 'bk_1',
    businessId: 'b1',
    serviceId: 's1',
    serviceName: 'Full Interior Detail',
    slotId: 'sl1',
    slotTime: formatSlot(today, 9),
    customer: { name: 'Alice Smith', email: 'alice@example.com', phone: '555-0101' },
    status: BookingStatus.CONFIRMED,
    depositPaid: 50,
    totalPrice: 200,
    createdAt: new Date(Date.now() - 86400000).toISOString()
  },
  {
    id: 'bk_2',
    businessId: 'b2',
    serviceId: 's3',
    serviceName: 'Deep Tissue Massage',
    slotId: 'sl_old',
    slotTime: formatSlot(today, 14),
    customer: { name: 'Bob Jones', email: 'bob@example.com', phone: '555-0202' },
    status: BookingStatus.PENDING,
    depositPaid: 40,
    totalPrice: 120,
    createdAt: new Date().toISOString()
  }
];