import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { LogOut, Calendar, LayoutDashboard, Settings, Plus, Trash2, Download, TrendingUp, DollarSign, Users, ChevronDown, Briefcase, Edit2, CreditCard, Bell, ShieldCheck } from 'lucide-react';
import { MOCK_BOOKINGS, MOCK_SERVICES, MOCK_BUSINESSES } from '../constants';
import { BookingStatus, Business } from '../types';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { useAuth } from '../context/AuthContext';

export const OperatorDashboard: React.FC = () => {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [activeTab, setActiveTab] = useState<'DASHBOARD' | 'BOOKINGS' | 'SERVICES' | 'SETTINGS'>('DASHBOARD');
  // Mock business selection (Simulating DetailPro logged in)
  const [currentBusiness, setCurrentBusiness] = useState<Business>(MOCK_BUSINESSES[0]);

  // Format helpers
  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }).format(new Date(iso));
  const fmtMoney = (n: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);

  // Filter Data
  const myBookings = MOCK_BOOKINGS.filter(b => b.businessId === currentBusiness.id);
  const myServices = MOCK_SERVICES.filter(s => s.businessId === currentBusiness.id);
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
             <div className="flex items-center gap-2 p-2 rounded-lg bg-gray-50 border border-gray-200">
                <div className="h-6 w-6 rounded bg-gray-800 text-white flex items-center justify-center text-xs">
                    {currentBusiness.name[0]}
                </div>
                <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{currentBusiness.name}</p>
                </div>
                <ChevronDown className="h-4 w-4 text-gray-400" />
             </div>
          </div>

          <p className="px-4 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2 mt-6">Main Menu</p>
          <NavItem id="DASHBOARD" label="Overview" icon={LayoutDashboard} />
          <NavItem id="BOOKINGS" label="Bookings" icon={Calendar} />
          <NavItem id="SERVICES" label="Services" icon={Briefcase} />
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
                  {activeTab === 'SETTINGS' && 'Business Settings'}
              </h1>
              <p className="text-gray-500 mt-1">
                {activeTab === 'DASHBOARD' && 'Welcome back, Operator.'}
                {activeTab === 'BOOKINGS' && 'Track and manage your customer appointments.'}
                {activeTab === 'SERVICES' && 'Configure what your customers can book.'}
                {activeTab === 'SETTINGS' && 'Manage your profile and preferences.'}
              </p>
           </div>
           {activeTab === 'BOOKINGS' && <Button>Export CSV</Button>}
           {activeTab === 'SERVICES' && <Button className="gap-2"><Plus className="h-4 w-4" /> Add Service</Button>}
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
                      <Button variant="outline" className="text-xs py-2 h-9">
                        Details
                      </Button>
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
          <div className="grid grid-cols-1 gap-4">
             {myServices.length === 0 ? (
               <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                  <p className="text-gray-500">No services configured.</p>
                  <Button className="mt-4" variant="outline">Create your first Service</Button>
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
                            <Button variant="ghost" className="p-2 text-gray-400 hover:text-blue-600">
                               <Edit2 className="h-4 w-4" />
                            </Button>
                            <Button variant="ghost" className="p-2 text-gray-400 hover:text-red-600">
                               <Trash2 className="h-4 w-4" />
                            </Button>
                         </div>
                      </div>
                   </Card>
                ))
             )}
          </div>
        )}
        
        {/* COMPLETED SETTINGS TAB */}
        {activeTab === 'SETTINGS' && (
            <div className="max-w-4xl space-y-6">
                {/* General Profile */}
                <Card className="p-6">
                   <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <Briefcase className="h-5 w-5 text-gray-400" /> 
                      Business Profile
                   </h3>
                   <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <Input label="Business Name" defaultValue={currentBusiness.name} />
                      <Input label="URL Slug" defaultValue={currentBusiness.slug} />
                      <div className="col-span-1 md:col-span-2">
                         <Input label="Description" defaultValue={currentBusiness.description} />
                      </div>
                      <Input label="Vertical" defaultValue={currentBusiness.vertical} disabled className="bg-gray-50 text-gray-500" />
                   </div>
                </Card>

                {/* Branding */}
                <Card className="p-6">
                   <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <Settings className="h-5 w-5 text-gray-400" /> 
                      Branding & Theme
                   </h3>
                   <div className="flex items-center gap-4">
                      <div className="w-12 h-12 rounded-full bg-blue-500 ring-2 ring-offset-2 ring-blue-500 cursor-pointer"></div>
                      <div className="w-12 h-12 rounded-full bg-emerald-500 cursor-pointer opacity-50 hover:opacity-100"></div>
                      <div className="w-12 h-12 rounded-full bg-zinc-800 cursor-pointer opacity-50 hover:opacity-100"></div>
                      <div className="w-12 h-12 rounded-full bg-primary-500 cursor-pointer opacity-50 hover:opacity-100"></div>
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
                    <Button variant="ghost">Discard Changes</Button>
                    <Button>Save Configuration</Button>
                </div>
            </div>
        )}
      </main>
    </div>
  );
};