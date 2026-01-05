const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080';

export interface Business {
  id: string;
  name: string;
  slug: string;
  vertical: string;
  description: string;
  theme_color: string;
  slotDurationMin: number;
  maxBookings: number;
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
  startTime: string;
  endTime: string;
  isBooked: boolean;
  bookingCount: number;
}

export interface BusinessAvailability {
  id: string;
  businessId: string;
  dayOfWeek: number; // 0=Monday, 6=Sunday
  startTime: string;
  endTime: string;
  isClosed: boolean;
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
  slotId: string;
  serviceName: string;
  slotTime: string;
  customer: CustomerDetails;
  status: string;
  depositPaid: number;
  totalPrice: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const token = localStorage.getItem('token');
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options?.headers as Record<string, string>,
    };

    if (token && !options?.headers?.['Authorization']) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        headers,
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

  // Auth
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getCurrentUser(): Promise<User> {
    return this.request<User>('/api/v1/auth/me');
  }

  setToken(token: string): void {
    localStorage.setItem('token', token);
  }

  getToken(): string | null {
    return localStorage.getItem('token');
  }

  logout(): void {
    localStorage.removeItem('token');
  }

  // Password Reset
  async forgotPassword(email: string): Promise<{ message: string }> {
    return this.request<{ message: string }>('/api/v1/auth/forgot-password', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  async resetPassword(token: string, password: string): Promise<{ message: string }> {
    return this.request<{ message: string }>('/api/v1/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify({ token, password }),
    });
  }

  // Businesses
  async getBusinesses(): Promise<Business[]> {
    return this.request<Business[]>('/api/v1/businesses');
  }

  async getBusiness(id: string): Promise<Business> {
    return this.request<Business>(`/api/v1/businesses/${id}`);
  }

  async getBusinessBySubdomain(): Promise<Business | null> {
    const slug = (window as any).location.hostname.split('.')[0];
    try {
      return await this.request<Business>(`/api/v1/business/by-subdomain?slug=${slug}`);
    } catch {
      return null;
    }
  }

  async createBusiness(data: {
    name: string;
    slug: string;
    vertical: string;
    description?: string;
    theme_color?: string;
  }): Promise<Business> {
    return this.request<Business>('/api/v1/businesses', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateBusiness(id: string, data: {
    name?: string;
    slug?: string;
    vertical?: string;
    description?: string;
    themeColor?: string;
  }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Services
  async getServicesByBusiness(businessId: string): Promise<Service[]> {
    return this.request<Service[]>(`/api/v1/businesses/${businessId}/services`);
  }

  async createService(businessId: string, data: {
    name: string;
    description?: string;
    durationMin: number;
    totalPrice: number;
    depositAmount: number;
  }): Promise<Service> {
    return this.request<Service>(`/api/v1/businesses/${businessId}/services`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateService(businessId: string, serviceId: string, data: {
    name?: string;
    description?: string;
    durationMin?: number;
    totalPrice?: number;
    depositAmount?: number;
  }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${businessId}/services/${serviceId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteService(businessId: string, serviceId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${businessId}/services/${serviceId}`, {
      method: 'DELETE',
    });
  }

  // Slots
  async getSlotsByBusiness(businessId: string): Promise<Slot[]> {
    return this.request<Slot[]>(`/api/v1/businesses/${businessId}/slots`);
  }

  async createSlot(businessId: string, data: {
    startTime: string;
    endTime: string;
  }): Promise<Slot> {
    return this.request<Slot>(`/api/v1/businesses/${businessId}/slots`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteSlot(businessId: string, slotId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${businessId}/slots/${slotId}`, {
      method: 'DELETE',
    });
  }

  // Bookings
  async createBooking(booking: {
    businessId: string;
    serviceId: string;
    slotId: string;
    customer: CustomerDetails;
  }): Promise<Booking> {
    return this.request<Booking>('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify(booking),
    });
  }

  async getBookingsByBusiness(businessId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/api/v1/businesses/${businessId}/bookings`);
  }

  async cancelBooking(businessId: string, bookingId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${businessId}/bookings/${bookingId}`, {
      method: 'DELETE',
    });
  }

  // Availability
  async getAvailability(businessId: string): Promise<BusinessAvailability[]> {
    return this.request<BusinessAvailability[]>(`/api/v1/businesses/${businessId}/availability`);
  }

  async setAvailability(businessId: string, data: {
    dayOfWeek?: number;
    startTime?: string;
    endTime?: string;
    isClosed?: boolean;
  }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/businesses/${businessId}/availability`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async generateSlots(businessId: string, data: {
    startDate: string;
    endDate: string;
    durationMin: number;
  }): Promise<Slot[]> {
    return this.request<Slot[]>(`/api/v1/businesses/${businessId}/slots/generate`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health');
  }
}

export const api = new ApiClient(API_BASE_URL);
