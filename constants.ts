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
    name: 'TintLab Studio',
    slug: 'tint-lab',
    vertical: 'Automotive',
    description: 'Window tinting and heat rejection packages for daily drivers.',
    themeColor: 'emerald',
  },
  {
    id: 'b3',
    name: 'ShineBay Detailing',
    slug: 'shine-bay',
    vertical: 'Automotive',
    description: 'Wash, polish, and protection packages for busy workshop teams.',
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
    totalPriceMinor: 20000,
    depositAmountMinor: 5000,
    currencyCode: 'USD',
  },
  {
    id: 's2',
    businessId: 'b1',
    name: 'Ceramic Coating Gold',
    description: '5-year protection package with paint correction.',
    durationMin: 360,
    totalPriceMinor: 120000,
    depositAmountMinor: 30000,
    currencyCode: 'USD',
  },
  // b2: Tint shop
  {
    id: 's3',
    businessId: 'b2',
    name: 'Nano Ceramic Tint',
    description: 'Premium tint package with high heat rejection film.',
    durationMin: 60,
    totalPriceMinor: 12000,
    depositAmountMinor: 4000,
    currencyCode: 'USD',
  },
  {
    id: 's4',
    businessId: 'b2',
    name: 'Front Two Windows Tint',
    description: 'Quick tint installation for front window pairs.',
    durationMin: 45,
    totalPriceMinor: 18000,
    depositAmountMinor: 6000,
    currencyCode: 'USD',
  },
  // b3: Detailing bay
  {
    id: 's5',
    businessId: 'b3',
    name: 'Paint Decontamination Detail',
    description: '2-hour decon wash, clay treatment, and gloss enhancement.',
    durationMin: 120,
    totalPriceMinor: 45000,
    depositAmountMinor: 15000,
    currencyCode: 'USD',
  },
  {
    id: 's6',
    businessId: 'b3',
    name: 'Express Exterior Wash',
    description: '30-minute wash and dry for quick turnaround jobs.',
    durationMin: 30,
    totalPriceMinor: 15000,
    depositAmountMinor: 15000, // Full payment
    currencyCode: 'USD',
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
    depositPaidMinor: 5000,
    totalPriceMinor: 20000,
    currencyCode: 'USD',
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
    depositPaidMinor: 4000,
    totalPriceMinor: 12000,
    currencyCode: 'USD',
    createdAt: new Date().toISOString()
  }
];
