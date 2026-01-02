import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { LogOut, Calendar, LayoutDashboard, Settings, Plus, Trash2, Download, TrendingUp, DollarSign, Users, ChevronDown, Briefcase, Edit2, CreditCard, Bell, ShieldCheck, Loader2, Clock } from 'lucide-react';
import { BookingStatus, Business, Service, Booking, Slot } from '../types';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { useAuth } from '../context/AuthContext';
import { api } from '../api';

export const OperatorDashboard: React.FC = () => {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [activeTab, setActiveTab] = useState<'DASHBOARD' | 'BOOKINGS' | 'SERVICES' | 'SLOTS' | 'SETTINGS'>('DASHBOARD');
  const [currentBusiness, setCurrentBusiness] = useState<Business | null>(null);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [slots, setSlots] = useState<Slot[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  
  const [showServiceForm, setShowServiceForm] = useState(false);
  const [editingService, setEditingService] = useState<Service | null>(null);
  const [showSlotForm, setShowSlotForm] = useState(false);
  const [serviceForm, setServiceForm] = useState({
    name: '',
    description: '',
    durationMin: 60,
    totalPrice: 0,
    depositAmount: 0
  });
  const [slotForm, setSlotForm] = useState({
    startTime: '',
    endTime: ''
  });

  const [businessForm, setBusinessForm] = useState({
    name: '',
    slug: '',
    vertical: '',
    description: '',
    theme_color: 'blue'
  });

  useEffect(() => {
    fetchData();
  }, [currentBusiness]);

  useEffect(() => {
    if (currentBusiness) {
      setBusinessForm({
        name: currentBusiness.name,
        slug: currentBusiness.slug,
        vertical: currentBusiness.vertical,
        description: currentBusiness.description,
        theme_color: currentBusiness.themeColor
      });
    }
  }, [currentBusiness]);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);

      const businessesData = await api.getBusinesses();
      setBusinesses(businessesData);

      if (businessesData.length > 0) {
        const selectedBusiness = currentBusiness || businessesData[0];
        setCurrentBusiness(selectedBusiness);

        const [bookingsData, servicesData, slotsData] = await Promise.all([
          api.getBookingsByBusiness(selectedBusiness.id),
          api.getServicesByBusiness(selectedBusiness.id),
          api.getSlotsByBusiness(selectedBusiness.id)
        ]);

        setBookings(bookingsData);
        setServices(servicesData);
        setSlots(slotsData);
      }
    } catch (err) {
      console.error('Failed to fetch data:', err);
      setError('Failed to load data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleBusinessChange = async (business: Business) => {
    setCurrentBusiness(business);
  };

  const handleSaveBusiness = async () => {
    if (!currentBusiness) return;
    
    try {
      setSaving(true);
      await api.updateBusiness(currentBusiness.id, businessForm);
      await fetchData();
      alert('Business updated successfully!');
    } catch (err) {
      console.error('Failed to update business:', err);
      alert('Failed to update business. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleCreateBusiness = async () => {
    if (!businessForm.name || !businessForm.slug) {
      alert('Please fill in all required fields');
      return;
    }

    try {
      setSaving(true);
      await api.createBusiness({
        name: businessForm.name,
        slug: businessForm.slug,
        vertical: businessForm.vertical || 'General',
        description: businessForm.description,
        theme_color: businessForm.theme_color
      });
      await fetchData();
      alert('Business created successfully!');
    } catch (err) {
      console.error('Failed to create business:', err);
      alert('Failed to create business. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleCreateService = async () => {
    if (!currentBusiness) return;

    try {
      setSaving(true);
      await api.createService(currentBusiness.id, serviceForm);
      setShowServiceForm(false);
      setServiceForm({
        name: '',
        description: '',
        durationMin: 60,
        totalPrice: 0,
        depositAmount: 0
      });
      await fetchData();
      alert('Service created successfully!');
    } catch (err) {
      console.error('Failed to create service:', err);
      alert('Failed to create service. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteService = async (serviceId: string) => {
    if (!currentBusiness) return;
    
    if (!confirm('Are you sure you want to delete this service?')) return;

    try {
      setSaving(true);
      await api.deleteService(currentBusiness.id, serviceId);
      await fetchData();
      alert('Service deleted successfully!');
    } catch (err) {
      console.error('Failed to delete service:', err);
      alert('Failed to delete service. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleEditService = (service: Service) => {
    setEditingService(service);
    setServiceForm({
      name: service.name,
      description: service.description,
      durationMin: service.durationMin,
      totalPrice: service.totalPrice,
      depositAmount: service.depositAmount
    });
    setShowServiceForm(true);
  };

  const handleUpdateService = async () => {
    if (!currentBusiness || !editingService) return;

    try {
      setSaving(true);
      await api.updateService(currentBusiness.id, editingService.id, serviceForm);
      setShowServiceForm(false);
      setEditingService(null);
      setServiceForm({
        name: '',
        description: '',
        durationMin: 60,
        totalPrice: 0,
        depositAmount: 0
      });
      await fetchData();
      alert('Service updated successfully!');
    } catch (err) {
      console.error('Failed to update service:', err);
      alert('Failed to update service. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }).format(new Date(iso));
  const fmtMoney = (n: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);
  const fmtTime = (iso: string) => new Intl.DateTimeFormat('en-US', { hour: 'numeric', minute: '2-digit' }).format(new Date(iso));

  const handleCreateSlot = async () => {
    if (!currentBusiness) return;

    try {
      setSaving(true);
      await api.createSlot(currentBusiness.id, slotForm);
      setShowSlotForm(false);
      setSlotForm({
        startTime: '',
        endTime: ''
      });
      await fetchData();
      alert('Slot created successfully!');
    } catch (err) {
      console.error('Failed to create slot:', err);
      alert('Failed to create slot. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteSlot = async (slotId: string) => {
    if (!currentBusiness) return;

    if (!confirm('Are you sure you want to delete this slot?')) return;

    try {
      setSaving(true);
      await api.deleteSlot(currentBusiness.id, slotId);
      await fetchData();
      alert('Slot deleted successfully!');
    } catch (err) {
      console.error('Failed to delete slot:', err);
      alert('Failed to delete slot. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleCancelBooking = async (bookingId: string) => {
    if (!currentBusiness) return;

    if (!confirm('Are you sure you want to cancel this booking? This will refund the deposit and make the slot available again.')) return;

    try {
      setSaving(true);
      await api.cancelBooking(currentBusiness.id, bookingId);
      await fetchData();
      alert('Booking cancelled successfully!');
    } catch (err) {
      console.error('Failed to cancel booking:', err);
      alert('Failed to cancel booking. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const myBookings = bookings;
  const myServices = services;
  const mySlots = slots;
  const totalRev = myBookings.reduce((acc, curr) => acc + curr.depositPaid, 0);

  const NavItem = ({ id, label, icon: Icon }: any) => (
    <button
      onClick={() => setActiveTab(id)}
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

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin text-gray-400 mx-auto mb-4" />
          <p className="text-gray-500">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">{error}</p>
          <Button onClick={fetchData}>Retry</Button>
        </div>
      </div>
    );
  }

  if (businesses.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center max-w-md w-full mx-4">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">No Businesses Found</h2>
          <p className="text-gray-500 mb-6">You don't have any businesses set up yet. Create one to get started.</p>
          <Card className="p-6">
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Business Name</label>
                <input
                  type="text"
                  value={businessForm.name}
                  onChange={(e) => setBusinessForm({...businessForm, name: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                  placeholder="e.g., My Business"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Vertical</label>
                <select
                  value={businessForm.vertical}
                  onChange={(e) => setBusinessForm({...businessForm, vertical: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select vertical...</option>
                  <option value="Automotive">Automotive</option>
                  <option value="Wellness">Wellness</option>
                  <option value="Professional">Professional Services</option>
                  <option value="Creative">Creative</option>
                  <option value="Education">Education</option>
                  <option value="General">General</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={businessForm.description}
                  onChange={(e) => setBusinessForm({...businessForm, description: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                  rows={3}
                  placeholder="Describe your business..."
                />
              </div>
              <Button onClick={handleCreateBusiness} disabled={saving} className="w-full">
                {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                Create Your First Business
              </Button>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white flex flex-col md:flex-row font-sans text-gray-900">
      {/* Sidebar */}
      <aside className="w-full md:w-64 bg-white border-r border-gray-200 flex-shrink-0 flex flex-col h-auto md:h-screen sticky top-0">
        <div className="p-6 flex items-center gap-2 border-b border-gray-100">
          <div className="h-8 w-8 bg-primary-500 rounded-lg flex items-center justify-center text-white font-bold">
            B
          </div>
          <span className="font-bold text-lg tracking-tight">Blytz.Cloud</span>
        </div>

        <div className="px-4 py-6 space-y-1 flex-1">
          <div className="mb-6 px-4">
             <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">Organization</p>
             <select
               value={currentBusiness?.id || ''}
               onChange={(e) => {
                 const business = businesses.find(b => b.id === e.target.value);
                 if (business) handleBusinessChange(business);
               }}
               className="w-full flex items-center gap-2 p-2 rounded-lg bg-gray-50 border border-gray-200 text-sm"
             >
                {businesses.map(b => (
                  <option key={b.id} value={b.id}>{b.name}</option>
                ))}
             </select>
          </div>

          <p className="px-4 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2 mt-6">Main Menu</p>
          <NavItem id="DASHBOARD" label="Overview" icon={LayoutDashboard} />
          <NavItem id="BOOKINGS" label="Bookings" icon={Calendar} />
          <NavItem id="SERVICES" label="Services" icon={Briefcase} />
          <NavItem id="SLOTS" label="Slots" icon={Clock} />
          <NavItem id="SETTINGS" label="Settings" icon={Settings} />
        </div>

        <div className="p-4 border-t border-gray-100">
          <button onClick={logout} className="flex items-center gap-2 text-sm text-gray-500 hover:text-red-600 transition-colors w-full px-4 py-2">
            <LogOut className="h-4 w-4" />
            Sign Out
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 overflow-y-auto bg-gray-50 p-6 md:p-10">
        <header className="mb-8 flex justify-between items-start">
           <div>
              <h1 className="text-2xl font-bold text-gray-900">
                  {activeTab === 'DASHBOARD' && 'Dashboard Overview'}
                  {activeTab === 'BOOKINGS' && 'Bookings Management'}
                  {activeTab === 'SERVICES' && 'Service Packages'}
                  {activeTab === 'SLOTS' && 'Availability Slots'}
                  {activeTab === 'SETTINGS' && 'Business Settings'}
              </h1>
              <p className="text-gray-500 mt-1">
                {activeTab === 'DASHBOARD' && 'Welcome back, Operator.'}
                {activeTab === 'BOOKINGS' && 'Track and manage your customer appointments.'}
                {activeTab === 'SERVICES' && 'Configure what your customers can book.'}
                {activeTab === 'SLOTS' && 'Manage available booking times.'}
                {activeTab === 'SETTINGS' && 'Manage your profile and preferences.'}
              </p>
           </div>
             {activeTab === 'BOOKINGS' && <Button>Export CSV</Button>}
             {activeTab === 'SERVICES' && (
               <Button className="gap-2" onClick={() => setShowServiceForm(true)}>
                 <Plus className="h-4 w-4" /> Add Service
               </Button>
             )}
             {activeTab === 'SLOTS' && (
               <Button className="gap-2" onClick={() => setShowSlotForm(true)}>
                 <Plus className="h-4 w-4" /> Add Slot
               </Button>
             )}
        </header>

        {activeTab === 'DASHBOARD' && (
           <div className="space-y-6">
              {/* Stats Grid */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                 <Card className="p-6">
                    <div className="flex items-center justify-between mb-4">
                       <span className="text-gray-500 text-sm font-medium">Total Revenue</span>
                       <div className="bg-green-100 p-2 rounded-full text-green-600"><DollarSign className="h-4 w-4" /></div>
                    </div>
                    <p className="text-3xl font-bold text-gray-900">{fmtMoney(totalRev)}</p>
                    <p className="text-xs text-green-600 mt-1 flex items-center font-medium"><TrendingUp className="h-3 w-3 mr-1" /> +12% from last week</p>
                 </Card>
                 <Card className="p-6">
                    <div className="flex items-center justify-between mb-4">
                       <span className="text-gray-500 text-sm font-medium">Active Bookings</span>
                       <div className="bg-blue-100 p-2 rounded-full text-blue-600"><Calendar className="h-4 w-4" /></div>
                    </div>
                    <p className="text-3xl font-bold text-gray-900">{myBookings.length}</p>
                    <p className="text-xs text-gray-400 mt-1">Next: Today, 2:00 PM</p>
                 </Card>
                 <Card className="p-6">
                    <div className="flex items-center justify-between mb-4">
                       <span className="text-gray-500 text-sm font-medium">Conversion Rate</span>
                       <div className="bg-purple-100 p-2 rounded-full text-purple-600"><Users className="h-4 w-4" /></div>
                    </div>
                    <p className="text-3xl font-bold text-gray-900">24%</p>
                    <p className="text-xs text-gray-400 mt-1">Page views to deposit</p>
                 </Card>
              </div>

              <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
                 <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-center">
                    <h3 className="font-bold text-gray-900">Recent Activity</h3>
                    <button className="text-sm text-primary-600 font-medium hover:underline">View All</button>
                 </div>
                 {myBookings.length === 0 ? (
                    <div className="p-12 text-center text-gray-400">No activity yet.</div>
                 ) : (
                    <div className="divide-y divide-gray-100">
                       {myBookings.map(b => (
                          <div key={b.id} className="p-4 px-6 flex items-center justify-between hover:bg-gray-50 transition-colors">
                             <div className="flex items-center gap-4">
                                <div className={`h-10 w-10 rounded-full flex items-center justify-center text-sm font-bold ${
                                   b.status === BookingStatus.CONFIRMED ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                                }`}>
                                   {b.customer.name[0]}
                                </div>
                                <div>
                                   <p className="font-medium text-gray-900">{b.customer.name}</p>
                                   <p className="text-sm text-gray-500">{b.serviceName} • {fmtMoney(b.depositPaid)} paid</p>
                                </div>
                             </div>
                             <div className="text-right">
                                <p className="text-sm text-gray-900">{fmtDate(b.createdAt)}</p>
                                <span className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 font-medium">{b.status}</span>
                             </div>
                          </div>
                       ))}
                    </div>
                 )}
              </div>
            </div>
         )}

         {/* BOOKINGS TAB */}
        {activeTab === 'BOOKINGS' && (
          <div className="space-y-4">
            {myBookings.length === 0 ? (
                <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                    <p className="text-gray-500">No bookings found for this business.</p>
                </div>
            ) : (
                myBookings.map(booking => (
                  <Card key={booking.id} className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                          booking.status === BookingStatus.CONFIRMED ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                        }`}>
                          {booking.status}
                        </span>
                        <span className="text-xs text-gray-400 font-mono">#{booking.id.slice(-4)}</span>
                      </div>
                      <h3 className="font-semibold text-gray-900">{booking.customer.name}</h3>
                      <p className="text-sm text-gray-500">{booking.serviceName}</p>
                      <div className="flex items-center gap-2 mt-2 text-sm text-gray-700">
                        <Calendar className="h-4 w-4 text-gray-400" />
                        {fmtDate(booking.slotTime)}
                      </div>
                    </div>
                    
                     <div className="flex items-center gap-3">
                       <div className="text-right mr-4 hidden sm:block">
                         <p className="text-sm font-medium text-gray-900">Paid: {fmtMoney(booking.depositPaid)}</p>
                         <p className="text-xs text-gray-500">Total: {fmtMoney(booking.totalPrice)}</p>
                       </div>
                       {booking.status !== BookingStatus.CANCELLED && (
                         <Button variant="outline" className="text-xs py-2 h-9" onClick={() => handleCancelBooking(booking.id)} disabled={saving}>
                           {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                           Cancel Booking
                         </Button>
                       )}
                       <Button variant="ghost" className="text-gray-400 hover:text-gray-900 p-2">
                         <Download className="h-4 w-4" />
                       </Button>
                     </div>
                  </Card>
                ))
            )}
          </div>
        )}

        {/* SERVICES MANAGEMENT TAB */}
        {activeTab === 'SERVICES' && (
          <div className="space-y-4">
             {showServiceForm && (
                <Card className="p-6 bg-blue-50 border-blue-200">
                  <h3 className="text-lg font-bold text-gray-900 mb-4">{editingService ? 'Edit Service' : 'Add New Service'}</h3>
                 <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                   <div className="col-span-1 md:col-span-2">
                     <label className="block text-sm font-medium text-gray-700 mb-1">Service Name</label>
                     <input
                       type="text"
                       value={serviceForm.name}
                       onChange={(e) => setServiceForm({...serviceForm, name: e.target.value})}
                       className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                       placeholder="e.g., Haircut, Massage"
                     />
                   </div>
                   <div className="col-span-1 md:col-span-2">
                     <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                     <textarea
                       value={serviceForm.description}
                       onChange={(e) => setServiceForm({...serviceForm, description: e.target.value})}
                       className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                       rows={3}
                       placeholder="Describe your service..."
                     />
                   </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Duration (minutes)</label>
                      <input
                        type="number"
                        value={serviceForm.durationMin}
                        onChange={(e) => setServiceForm({...serviceForm, durationMin: parseInt(e.target.value)})}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                        min="1"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Total Price ($)</label>
                      <input
                        type="number"
                        value={serviceForm.totalPrice}
                        onChange={(e) => setServiceForm({...serviceForm, totalPrice: parseFloat(e.target.value)})}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                        min="0"
                        step="0.01"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Deposit Amount ($)</label>
                      <input
                        type="number"
                        value={serviceForm.depositAmount}
                        onChange={(e) => setServiceForm({...serviceForm, depositAmount: parseFloat(e.target.value)})}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                        min="0"
                        step="0.01"
                      />
                    </div>
                 </div>
                  <div className="flex justify-end gap-3 mt-4">
                    <Button variant="outline" onClick={() => {
                      setShowServiceForm(false);
                      setEditingService(null);
                      setServiceForm({
                        name: '',
                        description: '',
                        durationMin: 60,
                        totalPrice: 0,
                        depositAmount: 0
                      });
                    }}>
                      Cancel
                    </Button>
                    <Button onClick={editingService ? handleUpdateService : handleCreateService} disabled={saving}>
                      {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                      {editingService ? 'Update Service' : 'Create Service'}
                    </Button>
                  </div>
               </Card>
             )}

             {myServices.length === 0 ? (
               <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                   <p className="text-gray-500">No services configured.</p>
                   <Button className="mt-4" variant="outline" onClick={() => setShowServiceForm(true)}>Create your first Service</Button>
               </div>
             ) : (
                 myServices.map((service) => (
                    <Card key={service.id} className="group flex flex-col sm:flex-row sm:items-center justify-between gap-6 hover:shadow-md transition-shadow">
                       <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                             <h3 className="text-lg font-bold text-gray-900">{service.name}</h3>
                             <div className="px-2 py-0.5 bg-gray-100 rounded text-xs text-gray-500 font-medium">
                                {service.durationMin} min
                             </div>
                          </div>
                          <p className="text-sm text-gray-500 max-w-xl">{service.description}</p>
                       </div>

                       <div className="flex items-center gap-6 border-t sm:border-t-0 sm:border-l border-gray-100 pt-4 sm:pt-0 sm:pl-6">
                          <div className="text-right min-w-[100px]">
                             <p className="text-sm text-gray-500">Total Price</p>
                             <p className="font-semibold text-gray-900">{fmtMoney(service.totalPrice)}</p>
                          </div>
                          <div className="text-right min-w-[100px]">
                             <p className="text-sm text-gray-500">Required Deposit</p>
                             <p className="font-bold text-primary-600">{fmtMoney(service.depositAmount)}</p>
                          </div>
                            <div className="flex items-center gap-1">
                               <Button variant="ghost" className="p-2 text-gray-400 hover:text-blue-600" onClick={() => handleEditService(service)}>
                                  <Edit2 className="h-4 w-4" />
                               </Button>
                               <Button variant="ghost" className="p-2 text-gray-400 hover:text-red-600" onClick={() => handleDeleteService(service.id)} disabled={saving}>
                                  {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                               </Button>
                           </div>
                      </div>
                   </Card>
                ))
              )}
           </div>
         )}

         {/* SLOTS TAB */}
         {activeTab === 'SLOTS' && (
           <div className="space-y-4">
             {showSlotForm && (
               <Card className="p-6 bg-blue-50 border-blue-200">
                 <h3 className="text-lg font-bold text-gray-900 mb-4">Add New Slot</h3>
                 <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                   <div>
                     <label className="block text-sm font-medium text-gray-700 mb-1">Start Time</label>
                     <input
                       type="datetime-local"
                       value={slotForm.startTime}
                       onChange={(e) => setSlotForm({...slotForm, startTime: e.target.value})}
                       className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                     />
                   </div>
                   <div>
                     <label className="block text-sm font-medium text-gray-700 mb-1">End Time</label>
                     <input
                       type="datetime-local"
                       value={slotForm.endTime}
                       onChange={(e) => setSlotForm({...slotForm, endTime: e.target.value})}
                       className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                     />
                   </div>
                 </div>
                 <div className="flex justify-end gap-3 mt-4">
                   <Button
                     variant="outline"
                     onClick={() => {
                       setShowSlotForm(false);
                       setSlotForm({ startTime: '', endTime: '' });
                     }}
                   >
                     Cancel
                   </Button>
                   <Button onClick={handleCreateSlot} disabled={saving}>
                     {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                     Create Slot
                   </Button>
                 </div>
               </Card>
             )}

             {mySlots.length === 0 ? (
               <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                 <p className="text-gray-500">No slots configured.</p>
                 <Button className="mt-4" variant="outline" onClick={() => setShowSlotForm(true)}>Create your first Slot</Button>
               </div>
             ) : (
               mySlots.map((slot) => (
                 <Card key={slot.id} className="flex items-center justify-between">
                   <div className="flex-1">
                     <div className="flex items-center gap-3">
                       <div className={`h-10 w-10 rounded-full flex items-center justify-center text-sm font-bold ${slot.isBooked ? 'bg-gray-100 text-gray-500' : 'bg-green-100 text-green-700'}`}>
                         <Clock className="h-4 w-4" />
                       </div>
                       <div>
                         <p className="font-semibold text-gray-900">{fmtDate(slot.startTime)}</p>
                         <p className="text-sm text-gray-500">{fmtTime(slot.startTime)} - {fmtTime(slot.endTime)}</p>
                       </div>
                     </div>
                   </div>
                   <div className="flex items-center gap-3 border-t border-gray-100 pt-4">
                     <Button variant="ghost" className="p-2 text-gray-400 hover:text-red-600" onClick={() => handleDeleteSlot(slot.id)} disabled={saving}>
                       {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                     </Button>
                     {slot.isBooked && <span className="text-xs font-medium text-gray-500">Booked</span>}
                   </div>
                 </Card>
               ))
             )}
           </div>
         )}

         {/* SETTINGS TAB */}
        {activeTab === 'SETTINGS' && (
            <div className="max-w-4xl space-y-6">
                {/* General Profile */}
                <Card className="p-6">
                   <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <Briefcase className="h-5 w-5 text-gray-400" />
                      Business Profile
                   </h3>
                   <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Business Name</label>
                        <input
                          type="text"
                          value={businessForm.name}
                          onChange={(e) => setBusinessForm({...businessForm, name: e.target.value})}
                          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">URL Slug</label>
                        <input
                          type="text"
                          value={businessForm.slug || ''}
                          onChange={(e) => {
                            const slug = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-');
                            setBusinessForm({...businessForm, slug: slug});
                          }}
                          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                          placeholder="my-business"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                          Your public URL: blytz.cloud/business/{businessForm.slug || 'your-slug'}
                        </p>
                      </div>
                      <div className="col-span-1 md:col-span-2">
                         <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                         <textarea
                           value={businessForm.description}
                           onChange={(e) => setBusinessForm({...businessForm, description: e.target.value})}
                           className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                           rows={3}
                         />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Vertical</label>
                        <select
                          value={businessForm.vertical}
                          onChange={(e) => setBusinessForm({...businessForm, vertical: e.target.value})}
                          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
                        >
                          <option value="">Select vertical...</option>
                          <option value="Automotive">Automotive</option>
                          <option value="Wellness">Wellness</option>
                          <option value="Professional">Professional Services</option>
                          <option value="Creative">Creative</option>
                          <option value="Education">Education</option>
                          <option value="General">General</option>
                        </select>
                      </div>
                   </div>
                </Card>

                {/* Branding */}
                <Card className="p-6">
                   <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <Settings className="h-5 w-5 text-gray-400" />
                      Branding & Theme
                   </h3>
                    <div className="flex items-center gap-4">
                       <div
                         className={`w-12 h-12 rounded-full cursor-pointer ring-2 ring-offset-2 ${businessForm?.theme_color === 'blue' ? 'ring-blue-500' : 'opacity-50 hover:opacity-100'} bg-blue-500`}
                         onClick={() => businessForm && setBusinessForm({...businessForm, theme_color: 'blue'})}
                       ></div>
                       <div
                         className={`w-12 h-12 rounded-full cursor-pointer ring-2 ring-offset-2 ${businessForm?.theme_color === 'emerald' ? 'ring-emerald-500' : 'opacity-50 hover:opacity-100'} bg-emerald-500`}
                         onClick={() => businessForm && setBusinessForm({...businessForm, theme_color: 'emerald'})}
                       ></div>
                       <div
                         className={`w-12 h-12 rounded-full cursor-pointer ring-2 ring-offset-2 ${businessForm?.theme_color === 'zinc' ? 'ring-zinc-800' : 'opacity-50 hover:opacity-100'} bg-zinc-800`}
                         onClick={() => businessForm && setBusinessForm({...businessForm, theme_color: 'zinc'})}
                       ></div>
                       <div
                         className={`w-12 h-12 rounded-full cursor-pointer ring-2 ring-offset-2 ${businessForm?.theme_color === 'purple' ? 'ring-purple-500' : 'opacity-50 hover:opacity-100'} bg-purple-500`}
                         onClick={() => businessForm && setBusinessForm({...businessForm, theme_color: 'purple'})}
                       ></div>
                   </div>
                   <p className="text-sm text-gray-500 mt-4">Selected color will be applied to your public booking page, buttons, and customer emails.</p>
                </Card>

                {/* Banking / Stripe */}
                <Card className="p-6">
                   <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <CreditCard className="h-5 w-5 text-gray-400" />
                      Payment Connection
                   </h3>
                   <div className="bg-green-50 border border-green-200 rounded-lg p-4 flex items-center justify-between">
                      <div className="flex items-center gap-3">
                         <div className="bg-green-100 p-2 rounded-full">
                            <ShieldCheck className="h-5 w-5 text-green-600" />
                         </div>
                         <div>
                            <p className="font-bold text-green-900 text-sm">Stripe Connected</p>
                            <p className="text-xs text-green-700">Payouts active • Blytz fees: 1.5%</p>
                         </div>
                      </div>
                      <Button variant="outline" className="text-xs h-8 bg-white">Manage Payouts</Button>
                   </div>
                </Card>

                {/* Action Bar */}
                <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
                    <Button
                      variant="ghost"
                      onClick={() => {
                        if (currentBusiness) {
                          setBusinessForm({
                            name: currentBusiness.name,
                            slug: currentBusiness.slug,
                            vertical: currentBusiness.vertical,
                            description: currentBusiness.description,
                            theme_color: currentBusiness.themeColor
                          });
                        }
                      }}
                    >
                      Discard Changes
                    </Button>
                    <Button onClick={handleSaveBusiness} disabled={saving}>
                      {saving ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                      Save Configuration
                    </Button>
                </div>
            </div>
        )}
      </main>
    </div>
  );
};