const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface Business {
  id: string;
  name: string;
  slug: string;
  vertical: string;
  description: string;
  theme_color: string;
}

export interface Service {
  id: string;
  business_id: string;
  name: string;
  description: string;
  duration_min: number;
  total_price: number;
  deposit_amount: number;
}

export interface Slot {
  id: string;
  business_id: string;
  start_time: string;
  end_time: string;
  is_booked: boolean;
}

export interface CustomerDetails {
  name: string;
  email: string;
  phone: string;
}

export interface Booking {
  id: string;
  business_id: string;
  service_id: string;
  slot_id: string;
  service_name: string;
  slot_time: string;
  customer: CustomerDetails;
  status: string;
  deposit_paid: number;
  total_price: number;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Businesses
  async getBusinesses(): Promise<Business[]> {
    return this.request<Business[]>('/api/v1/businesses');
  }

  async getBusiness(id: string): Promise<Business> {
    return this.request<Business>(`/api/v1/businesses/${id}`);
  }

  // Services
  async getServicesByBusiness(businessId: string): Promise<Service[]> {
    return this.request<Service[]>(`/api/v1/businesses/${businessId}/services`);
  }

  // Slots
  async getSlotsByBusiness(businessId: string): Promise<Slot[]> {
    return this.request<Slot[]>(`/api/v1/businesses/${businessId}/slots`);
  }

  // Bookings
  async createBooking(booking: Omit<Booking, 'id' | 'status' | 'created_at' | 'updated_at'>): Promise<Booking> {
    return this.request<Booking>('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify(booking),
    });
  }

  async getBookingsByBusiness(businessId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/api/v1/businesses/${businessId}/bookings`);
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health');
  }
}

export const api = new ApiClient(API_BASE_URL);
