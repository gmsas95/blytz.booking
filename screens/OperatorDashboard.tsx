import React, { useState, useEffect } from 'react';
import { LogOut, Calendar, LayoutDashboard, Settings, Plus, Trash2, Download, TrendingUp, DollarSign, Users, Briefcase, Edit2, CreditCard, Bell, ShieldCheck } from 'lucide-react';
import { api, Business, Service, Booking } from '../api';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Input } from '../components/Input';

interface OperatorDashboardProps {
  onLogout: () => void;
}

export const OperatorDashboard: React.FC<OperatorDashboardProps> = ({ onLogout }) => {
  const [activeTab, setActiveTab] = useState<'DASHBOARD' | 'BOOKINGS' | 'SERVICES' | 'SETTINGS'>('DASHBOARD');
  const [loading, setLoading] = useState(false);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [currentBusiness, setCurrentBusiness] = useState<Business | null>(null);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [services, setServices] = useState<Service[]>([]);

  useEffect(() => {
    const fetchBusinesses = async () => {
      try {
        setLoading(true);
        const data = await api.getBusinesses();
        setBusinesses(data);
        
        if (data.length > 0) {
          setCurrentBusiness(data[0]);
          await fetchBusinessData(data[0].id);
        }
      } catch (err) {
        console.error('Failed to fetch businesses:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchBusinesses();
  }, []);

  const fetchBusinessData = async (businessId: string) => {
    try {
      const [bookingsData, servicesData] = await Promise.all([
        api.getBookingsByBusiness(businessId),
        api.getServicesByBusiness(businessId),
      ]);
      setBookings(bookingsData);
      setServices(servicesData);
    } catch (err) {
      console.error('Failed to fetch business data:', err);
    }
  };

  const handleBusinessChange = (businessId: string) => {
    const business = businesses.find(b => b.id === businessId);
    if (business) {
      setCurrentBusiness(business);
      fetchBusinessData(business.id);
      setActiveTab('DASHBOARD');
    }
  };

  const NavItem = ({ id, label, icon: Icon }: { id: string; label: string; icon: any }) => (
    <button
      onClick={() => setActiveTab(id as any)}
      className={`w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg transition-colors ${
        activeTab === id
          ? 'bg-gray-100 text-gray-900' 
          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
      }`}
    >
      <Icon className="h-5 w-5" />
      {label}
    </button>
  );

  const fmtMoney = (n: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);
  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', year: 'numeric' }).format(new Date(iso));
  const fmtTime = (iso: string) => new Intl.DateTimeFormat('en-US', { hour: 'numeric', minute: '2-digit' }).format(new Date(iso));

  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center p-6">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white flex flex-col md:flex-row font-sans text-gray-900">
      <aside className="w-full md:w-64 bg-white border-r border-gray-200 flex-shrink-0 flex flex-col h-auto sticky top-0">
        <div className="p-6 flex items-center gap-2 border-b border-gray-100">
          <div className="h-8 w-8 bg-primary-500 rounded-lg flex items-center justify-center text-white font-bold">
            B
          </div>
          <span className="font-bold text-lg tracking-tight">Blytz.Cloud</span>
        </div>

        <div className="px-4 py-2 space-y-1 flex-1">
          <div>
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">Organization</p>
            <select 
              value={currentBusiness?.id || ''} 
              onChange={(e) => handleBusinessChange(e.target.value)}
              className="w-full p-2 rounded-lg bg-gray-50 border border-gray-200 text-sm"
            >
              {businesses.map(b => (
                <option key={b.id} value={b.id}>{b.name}</option>
              ))}
            </select>
          </div>

          <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2 mt-4">Main Menu</p>
          <NavItem id="DASHBOARD" label="Overview" icon={LayoutDashboard} />
          <NavItem id="BOOKINGS" label="Bookings" icon={Calendar} />
          <NavItem id="SERVICES" label="Services" icon={Briefcase} />
          <NavItem id="SETTINGS" label="Settings" icon={Settings} />
        </div>

        <div className="p-4 border-t border-gray-100">
          <button onClick={onLogout} className="flex items-center gap-2 text-sm text-gray-500 hover:text-red-600 transition-colors w-full px-4 py-2.5">
            <LogOut className="h-4 w-4" />
            Sign Out
          </button>
        </div>
      </aside>

      <main className="flex-1 overflow-y-auto bg-gray-50 p-6 md:p-10">
        <header className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900">
            {activeTab === 'DASHBOARD' && 'Dashboard Overview'}
            {activeTab === 'BOOKINGS' && 'Bookings Management'}
            {activeTab === 'SERVICES' && 'Service Packages'}
            {activeTab === 'SETTINGS' && 'Business Settings'}
          </h1>
          <p className="text-gray-500 mt-1">
            {activeTab === 'DASHBOARD' && 'Welcome back, Operator.'}
            {activeTab === 'BOOKINGS' && 'Track and manage your customer appointments.'}
            {activeTab === 'SERVICES' && 'Configure what your customers can book.'}
            {activeTab === 'SETTINGS' && 'Manage your profile and preferences.'}
          </p>
        </header>

        {activeTab === 'DASHBOARD' && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Total Revenue</span>
                  <div className="bg-green-100 p-2 rounded-full text-green-600"><DollarSign className="h-4 w-4" /></div>
                </div>
                <div className="text-3xl font-bold text-gray-900">
                  {fmtMoney(bookings.reduce((sum, b) => sum + b.deposit_paid, 0))}
                </div>
                <div className="text-xs text-green-600 mt-1"><TrendingUp className="h-3 w-3 mr-1 inline" /> Total deposits collected</div>
              </Card>

              <Card>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Active Bookings</span>
                  <div className="bg-blue-100 p-2 rounded-full text-blue-600"><Calendar className="h-4 w-4" /></div>
                </div>
                <div className="text-3xl font-bold text-gray-900">{bookings.filter(b => b.status !== 'CANCELLED').length}</div>
                <div className="text-xs text-gray-400 mt-1">Total bookings</div>
              </Card>

              <Card>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Services</span>
                  <div className="bg-purple-100 p-2 rounded-full text-purple-600"><Briefcase className="h-4 w-4" /></div>
                </div>
                <div className="text-3xl font-bold text-gray-900">{services.length}</div>
                <div className="text-xs text-gray-400 mt-1">Configured packages</div>
              </Card>
            </div>

            <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
              <div className="px-6 py-4 flex justify-between items-center">
                <h3 className="font-bold text-gray-900">Recent Activity</h3>
                <button className="text-sm text-primary-600 font-medium hover:underline" onClick={() => setActiveTab('BOOKINGS')}>View All</button>
              </div>
              <div>
                {bookings.length === 0 ? (
                  <div className="p-12 text-center text-gray-400">No bookings yet.</div>
                ) : (
                  <div className="divide-y divide-gray-100">
                    {bookings.slice(0, 5).map(b => (
                      <div key={b.id} className="p-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
                        <div className="flex items-center gap-3">
                          <div className={`p-2 px-2 py-0.5 rounded-full text-xs font-medium ${
                            b.status === 'CONFIRMED' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                          }`}>
                            {b.status}
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium">{b.customer.name}</p>
                            <p className="text-sm text-gray-500">{b.service_name}</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-sm font-medium">{fmtMoney(b.deposit_paid)}</div>
                          <div className="text-xs text-gray-400">{fmtDate(b.slot_time)}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'BOOKINGS' && (
          <div className="space-y-6">
            {bookings.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                <p className="text-gray-500">No bookings found for this business.</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {bookings.map(booking => (
                  <Card key={booking.id} className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div className="flex-1">
                      <div className={`p-2 px-2 py-0.5 rounded-full text-xs font-medium inline-block mb-2 ${
                        booking.status === 'CONFIRMED' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                      }`}>
                        {booking.status}
                      </div>
                      <p className="text-sm font-medium">{booking.customer.name}</p>
                      <p className="text-sm text-gray-500">{booking.service_name}</p>
                      <p className="text-sm text-gray-900 mt-1">{fmtDate(booking.slot_time)} at {fmtTime(booking.slot_time)}</p>
                    </div>

                    <div className="flex items-center gap-2">
                      <div className="text-right min-w-[100px]">
                        <p className="text-sm font-medium text-gray-900">Paid: {fmtMoney(booking.deposit_paid)}</p>
                        <div className="text-sm text-gray-400 mt-1">Total: {fmtMoney(booking.total_price)}</div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Button variant="outline" className="text-xs py-2 h-9">
                          Details
                        </Button>
                        <Button variant="ghost" className="p-2 h-9">
                          <Download className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </Card>
                ))}
              </div>
            )}
          </div>
        )}

        {activeTab === 'SERVICES' && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {services.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                <p className="text-gray-500">No services configured.</p>
              </div>
            ) : (
              services.map((service) => (
                <Card key={service.id} className="group flex flex-col sm:flex-row sm:items-center justify-between gap-6 hover:shadow-md transition-shadow">
                  <div className="flex-1">
                    <div className="p-4">
                      <h3 className="text-lg font-bold text-gray-900">{service.name}</h3>
                      <p className="text-sm text-gray-500 mb-2 line-clamp-2">{service.description}</p>
                      <p className="text-xs text-gray-400 font-mono">
                        {service.duration_min} min
                      </p>
                    </div>

                    <div className="flex items-center gap-2 border-t sm:border-t-0 sm:border-l sm:border-t border-gray-200 pt-4">
                      <span className="text-sm text-gray-500">Total Price</span>
                      <div className="text-2xl font-bold text-gray-900">{fmtMoney(service.total_price)}</div>
                      <div className="text-sm text-gray-500 mt-1">Required Deposit</div>
                      <div className="text-xl font-bold text-primary-600">{fmtMoney(service.deposit_amount)}</div>
                    </div>

                    <div className="flex items-center gap-2 mt-6 sm:border-t sm:border-t border-gray-200 pt-4">
                      <Button variant="outline" className="text-xs py-2 h-9">
                        Edit
                      </Button>
                      <Button variant="ghost" className="p-2 h-9">
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}

        {activeTab === 'SETTINGS' && (
          <div className="max-w-4xl space-y-6">
            <Card className="p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                <Briefcase className="h-5 w-5 text-gray-400" />
                Business Profile
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input label="Business Name" defaultValue={currentBusiness?.name || ''} />
                <Input label="URL Slug" defaultValue={currentBusiness?.slug || ''} />
                <Input label="Description" defaultValue={currentBusiness?.description || ''} />
                <Input label="Vertical" defaultValue={currentBusiness?.vertical || ''} disabled />
              </div>
            </Card>

            <Card className="p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                <Settings className="h-5 w-5 text-gray-400" />
                Branding & Theme
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Primary Brand Color</label>
                  <div className="flex items-center gap-2">
                    <input type="color" defaultValue={currentBusiness?.theme_color || '#3b82f6'} className="h-10 w-20 rounded border" />
                    <Input placeholder="#3b82f6" defaultValue={currentBusiness?.theme_color || ''} />
                  </div>
                </div>
              </div>
              <div className="col-span-2 text-xs text-gray-400 mt-2">Selected color will be applied to your public booking page.</div>
            </Card>

            <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
              <Button variant="ghost">Discard Changes</Button>
              <Button>Save Configuration</Button>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};
