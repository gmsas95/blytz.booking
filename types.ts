export enum BookingStatus {
  PENDING = 'PENDING',
  CONFIRMED = 'CONFIRMED',
  COMPLETED = 'COMPLETED',
  CANCELLED = 'CANCELLED'
}

export interface Business {
  id: string;
  name: string;
  slug: string;
  vertical: string; // e.g., 'Automotive', 'Wellness', 'Professional'
  description: string;
  themeColor: string; // Hex code for branding
}

export interface Service {
  id: string;
  businessId: string;
  name: string;
  description: string;
  durationMin: number;
  totalPrice: number;
  depositAmount: number;
}

export interface Slot {
  id: string;
  businessId: string;
  startTime: string; // ISO string
  endTime: string; // ISO string
  isBooked: boolean;
}

export interface CustomerDetails {
  name: string;
  email: string;
  phone: string;
}

export interface Booking {
  id: string;
  businessId: string;
  serviceId: string;
  serviceName: string;
  slotId: string;
  slotTime: string;
  customer: CustomerDetails;
  status: BookingStatus;
  depositPaid: number;
  totalPrice: number;
  createdAt: string;
}

export enum ViewState {
  SAAS_LANDING = 'SAAS_LANDING',
  PUBLIC_BOOKING = 'PUBLIC_BOOKING',
  CONFIRMATION = 'CONFIRMATION',
  LOGIN = 'LOGIN',
  DASHBOARD = 'DASHBOARD'
}
