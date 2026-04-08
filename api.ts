const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
  status: number;
  body?: unknown;

  constructor(message: string, status: number, body?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.body = body;
  }
}

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
  total_price_minor: number;
  deposit_amount_minor: number;
  currency_code: string;
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
  deposit_paid_minor: number;
  total_price_minor: number;
  currency_code: string;
}

export interface CreateBookingRequest {
  business_id: string;
  service_id: string;
  slot_id: string;
  customer: CustomerDetails;
}

export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

export interface MembershipBusiness {
  id: string;
  name: string;
  slug: string;
  vertical: string;
  description: string;
  theme_color: string;
}

export interface Membership {
  id: string;
  business_id: string;
  role: 'OWNER' | 'STAFF';
  business: MembershipBusiness;
}

export interface CurrentUserResponse {
  user: User;
  memberships: Membership[];
  active_business_id?: string;
}

export interface CustomerRecord {
  id: string;
  business_id: string;
  name: string;
  email: string;
  phone: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface VehicleRecord {
  id: string;
  business_id: string;
  customer_id: string;
  year: number;
  make: string;
  model: string;
  color: string;
  license_plate: string;
  customer: CustomerRecord;
  created_at: string;
  updated_at: string;
}

export interface JobRecord {
  id: string;
  business_id: string;
  customer_id: string;
  vehicle_id: string;
  booking_id?: string;
  title: string;
  status: 'SCHEDULED' | 'IN_PROGRESS' | 'READY' | 'DELIVERED';
  scheduled_at: string;
  notes: string;
  customer: CustomerRecord;
  vehicle: VehicleRecord;
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
    const headers = new Headers(options?.headers);
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }
    if (token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${token}`);
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const contentType = response.headers.get('content-type') || '';
    const isJSON = contentType.includes('application/json');
    const payload = isJSON ? await response.json() : await response.text();

    if (!response.ok) {
      const message = typeof payload === 'object' && payload && 'error' in payload
        ? String((payload as { error: string }).error)
        : `API error: ${response.status} ${response.statusText}`;
      throw new ApiError(message, response.status, payload);
    }

    return payload as T;
  }

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

  async getCurrentUser(): Promise<CurrentUserResponse> {
    return this.request<CurrentUserResponse>('/api/v1/auth/me');
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

  async getBusinesses(): Promise<Business[]> {
    return this.request<Business[]>('/api/v1/businesses');
  }

  async getBusiness(id: string): Promise<Business> {
    return this.request<Business>(`/api/v1/businesses/${id}`);
  }

  async getServicesByBusiness(businessId: string): Promise<Service[]> {
    return this.request<Service[]>(`/api/v1/businesses/${businessId}/services`);
  }

  async getSlotsByBusiness(businessId: string): Promise<Slot[]> {
    return this.request<Slot[]>(`/api/v1/businesses/${businessId}/slots`);
  }

  async createBooking(booking: CreateBookingRequest): Promise<Booking> {
    return this.request<Booking>('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify(booking),
    });
  }

  async getBookingsByBusiness(businessId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/api/v1/businesses/${businessId}/bookings`);
  }

  async getCustomersByBusiness(businessId: string): Promise<CustomerRecord[]> {
    return this.request<CustomerRecord[]>(`/api/v1/businesses/${businessId}/customers`);
  }

  async getVehiclesByBusiness(businessId: string): Promise<VehicleRecord[]> {
    return this.request<VehicleRecord[]>(`/api/v1/businesses/${businessId}/vehicles`);
  }

  async getJobsByBusiness(businessId: string): Promise<JobRecord[]> {
    return this.request<JobRecord[]>(`/api/v1/businesses/${businessId}/jobs`);
  }

  async healthCheck(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health');
  }
}

export const api = new ApiClient(API_BASE_URL);
