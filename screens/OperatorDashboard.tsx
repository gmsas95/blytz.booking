import React, { useEffect, useState } from 'react';
import { Calendar, LayoutDashboard, Settings, Plus, Trash2, Download, TrendingUp, DollarSign, Users, ChevronDown, Briefcase, Edit2, CreditCard, ShieldCheck, LogOut, CarFront, ClipboardList } from 'lucide-react';
import { api, Booking, Service, CustomerRecord, VehicleRecord, JobRecord } from '../api';
import { BookingStatus } from '../types';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { useAuth } from '../context/AuthContext';
import { formatMoneyFromMinor } from '../utils/money';

export const OperatorDashboard: React.FC = () => {
  const { logout, currentUser, memberships, activeMembership, activeBusinessId, setActiveBusinessId } = useAuth();
  const [activeTab, setActiveTab] = useState<'DASHBOARD' | 'BOOKINGS' | 'CUSTOMERS' | 'VEHICLES' | 'JOBS' | 'SERVICES' | 'SETTINGS'>('DASHBOARD');
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [customers, setCustomers] = useState<CustomerRecord[]>([]);
  const [vehicles, setVehicles] = useState<VehicleRecord[]>([]);
  const [jobs, setJobs] = useState<JobRecord[]>([]);
  const [loadingData, setLoadingData] = useState(true);
  const [dataError, setDataError] = useState<string | null>(null);

  useEffect(() => {
    const loadWorkshopData = async () => {
      if (!activeBusinessId) {
        setBookings([]);
        setServices([]);
        setLoadingData(false);
        return;
      }

      setLoadingData(true);
      setDataError(null);

      try {
        const [bookingsResponse, servicesResponse, customersResponse, vehiclesResponse, jobsResponse] = await Promise.all([
          api.getBookingsByBusiness(activeBusinessId),
          api.getServicesByBusiness(activeBusinessId),
          api.getCustomersByBusiness(activeBusinessId),
          api.getVehiclesByBusiness(activeBusinessId),
          api.getJobsByBusiness(activeBusinessId),
        ]);
        setBookings(bookingsResponse);
        setServices(servicesResponse);
        setCustomers(customersResponse);
        setVehicles(vehiclesResponse);
        setJobs(jobsResponse);
      } catch (error: any) {
        setBookings([]);
        setServices([]);
        setCustomers([]);
        setVehicles([]);
        setJobs([]);
        setDataError(error?.message || 'Failed to load workshop data.');
      } finally {
        setLoadingData(false);
      }
    };

    loadWorkshopData();
  }, [activeBusinessId]);

  const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }).format(new Date(iso));
  const fmtMoney = (amountMinor: number, currencyCode: string = 'USD') => formatMoneyFromMinor(amountMinor, currencyCode);
  const totalRev = bookings.reduce((acc, curr) => acc + curr.deposit_paid_minor, 0);

  const NavItem = ({ id, label, icon: Icon }: any) => (
    <button
      onClick={() => setActiveTab(id)}
      className={`w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg transition-colors ${
        activeTab === id ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
      }`}
    >
      <Icon className="h-5 w-5" />
      {label}
    </button>
  );

  if (!memberships.length) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-6">
        <Card className="max-w-lg w-full text-center space-y-4">
          <h1 className="text-2xl font-bold text-gray-900">No workshop membership yet</h1>
          <p className="text-sm text-gray-500">
            Your account is authenticated, but it is not assigned to a workshop yet. For demo access, use the seeded owner account or add a membership in the backend.
          </p>
          <div className="text-xs text-gray-400">Current user: {currentUser?.email || 'unknown'}</div>
          <Button variant="outline" onClick={logout}>Sign Out</Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white flex flex-col md:flex-row font-sans text-gray-900">
      <aside className="w-full md:w-64 bg-white border-r border-gray-200 flex-shrink-0 flex flex-col h-auto md:h-screen sticky top-0">
        <div className="p-6 flex items-center gap-2 border-b border-gray-100">
          <div className="h-8 w-8 bg-primary-500 rounded-lg flex items-center justify-center text-white font-bold">B</div>
          <span className="font-bold text-lg tracking-tight">Blytz.Auto</span>
        </div>

        <div className="px-4 py-6 space-y-1 flex-1">
          <div className="mb-6 px-4">
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">Workshop</p>
            <div className="space-y-2">
              {memberships.map((membership) => {
                const isActive = membership.business_id === activeBusinessId;
                return (
                  <button
                    key={membership.id}
                    onClick={() => setActiveBusinessId(membership.business_id)}
                    className={`w-full flex items-center gap-2 p-2 rounded-lg border text-left ${isActive ? 'bg-gray-100 border-gray-300' : 'bg-gray-50 border-gray-200'}`}
                  >
                    <div className="h-6 w-6 rounded bg-gray-800 text-white flex items-center justify-center text-xs">
                      {membership.business.name[0]}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">{membership.business.name}</p>
                      <p className="text-xs text-gray-500">{membership.role}</p>
                    </div>
                    <ChevronDown className="h-4 w-4 text-gray-400" />
                  </button>
                );
              })}
            </div>
          </div>

          <p className="px-4 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2 mt-6">Main Menu</p>
          <NavItem id="DASHBOARD" label="Overview" icon={LayoutDashboard} />
          <NavItem id="BOOKINGS" label="Bookings" icon={Calendar} />
          <NavItem id="CUSTOMERS" label="Customers" icon={Users} />
          <NavItem id="VEHICLES" label="Vehicles" icon={CarFront} />
          <NavItem id="JOBS" label="Jobs" icon={ClipboardList} />
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

      <main className="flex-1 overflow-y-auto bg-gray-50 p-6 md:p-10">
        <header className="mb-8 flex justify-between items-start gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {activeTab === 'DASHBOARD' && 'Dashboard Overview'}
              {activeTab === 'BOOKINGS' && 'Bookings Management'}
              {activeTab === 'CUSTOMERS' && 'Customer Records'}
              {activeTab === 'VEHICLES' && 'Vehicle Records'}
              {activeTab === 'JOBS' && 'Job Board'}
              {activeTab === 'SERVICES' && 'Service Packages'}
              {activeTab === 'SETTINGS' && 'Workshop Settings'}
            </h1>
            <p className="text-gray-500 mt-1">
              {activeTab === 'DASHBOARD' && `Welcome back, ${currentUser?.name || 'Operator'}.`}
              {activeTab === 'BOOKINGS' && 'Track and manage your workshop appointments.'}
              {activeTab === 'CUSTOMERS' && 'View the customers assigned to this workshop.'}
              {activeTab === 'VEHICLES' && 'Track vehicles connected to workshop customers.'}
              {activeTab === 'JOBS' && 'Monitor scheduled and active workshop jobs.'}
              {activeTab === 'SERVICES' && 'Configure what your customers can book.'}
              {activeTab === 'SETTINGS' && 'Manage your workshop profile and preferences.'}
            </p>
          </div>
          {activeTab === 'BOOKINGS' && <Button>Export CSV</Button>}
          {activeTab === 'CUSTOMERS' && <Button className="gap-2"><Plus className="h-4 w-4" /> Add Customer</Button>}
          {activeTab === 'VEHICLES' && <Button className="gap-2"><Plus className="h-4 w-4" /> Add Vehicle</Button>}
          {activeTab === 'JOBS' && <Button className="gap-2"><Plus className="h-4 w-4" /> New Job</Button>}
          {activeTab === 'SERVICES' && <Button className="gap-2"><Plus className="h-4 w-4" /> Add Service</Button>}
        </header>

        {dataError && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{dataError}</div>
        )}

        {activeTab === 'DASHBOARD' && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Deposits Collected</span>
                  <div className="bg-green-100 p-2 rounded-full text-green-600"><DollarSign className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{fmtMoney(totalRev)}</p>
                <p className="text-xs text-green-600 mt-1 flex items-center font-medium"><TrendingUp className="h-3 w-3 mr-1" /> Live workshop-scoped total</p>
              </Card>
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Active Bookings</span>
                  <div className="bg-blue-100 p-2 rounded-full text-blue-600"><Calendar className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{loadingData ? '—' : bookings.length}</p>
                <p className="text-xs text-gray-400 mt-1">Protected by workshop membership</p>
              </Card>
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Customers</span>
                  <div className="bg-purple-100 p-2 rounded-full text-purple-600"><Users className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{loadingData ? '—' : customers.length}</p>
                <p className="text-xs text-gray-400 mt-1">Visible for the selected workshop only</p>
              </Card>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Vehicles</span>
                  <div className="bg-gray-100 p-2 rounded-full text-gray-700"><CarFront className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{loadingData ? '—' : vehicles.length}</p>
                <p className="text-xs text-gray-400 mt-1">Linked to workshop customers</p>
              </Card>
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Jobs</span>
                  <div className="bg-orange-100 p-2 rounded-full text-orange-600"><ClipboardList className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{loadingData ? '—' : jobs.length}</p>
                <p className="text-xs text-gray-400 mt-1">Scheduled and active board items</p>
              </Card>
              <Card className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-gray-500 text-sm font-medium">Active Services</span>
                  <div className="bg-indigo-100 p-2 rounded-full text-indigo-600"><Briefcase className="h-4 w-4" /></div>
                </div>
                <p className="text-3xl font-bold text-gray-900">{loadingData ? '—' : services.length}</p>
                <p className="text-xs text-gray-400 mt-1">Available in the public booking flow</p>
              </Card>
            </div>

            <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-center">
                <h3 className="font-bold text-gray-900">Recent Activity</h3>
                <button className="text-sm text-primary-600 font-medium hover:underline">View All</button>
              </div>
              {loadingData ? (
                <div className="p-12 text-center text-gray-400">Loading activity...</div>
              ) : bookings.length === 0 ? (
                <div className="p-12 text-center text-gray-400">No activity yet for this workshop.</div>
              ) : (
                <div className="divide-y divide-gray-100">
                  {bookings.map((booking) => (
                    <div key={booking.id} className="p-4 px-6 flex items-center justify-between hover:bg-gray-50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className={`h-10 w-10 rounded-full flex items-center justify-center text-sm font-bold ${
                          booking.status === BookingStatus.CONFIRMED ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                        }`}>
                          {booking.customer.name[0]}
                        </div>
                        <div>
                          <p className="font-medium text-gray-900">{booking.customer.name}</p>
                          <p className="text-sm text-gray-500">{booking.service_name} • {fmtMoney(booking.deposit_paid_minor, booking.currency_code)} paid</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-gray-900">{fmtDate(booking.slot_time)}</p>
                        <span className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 font-medium">{booking.status}</span>
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
            {loadingData ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">Loading bookings...</div>
            ) : bookings.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                <p className="text-gray-500">No bookings found for this workshop.</p>
              </div>
            ) : (
              bookings.map((booking) => (
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
                    <p className="text-sm text-gray-500">{booking.service_name}</p>
                    <div className="flex items-center gap-2 mt-2 text-sm text-gray-700">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      {fmtDate(booking.slot_time)}
                    </div>
                  </div>

                  <div className="flex items-center gap-3">
                    <div className="text-right mr-4 hidden sm:block">
                      <p className="text-sm font-medium text-gray-900">Paid: {fmtMoney(booking.deposit_paid_minor, booking.currency_code)}</p>
                      <p className="text-xs text-gray-500">Total: {fmtMoney(booking.total_price_minor, booking.currency_code)}</p>
                    </div>
                    <Button variant="outline" className="text-xs py-2 h-9">Details</Button>
                    <Button variant="ghost" className="text-gray-400 hover:text-gray-900 p-2">
                      <Download className="h-4 w-4" />
                    </Button>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}

        {activeTab === 'SERVICES' && (
          <div className="grid grid-cols-1 gap-4">
            {loadingData ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">Loading services...</div>
            ) : services.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300">
                <p className="text-gray-500">No services configured.</p>
                <Button className="mt-4" variant="outline">Create your first Service</Button>
              </div>
            ) : (
              services.map((service) => (
                <Card key={service.id} className="group flex flex-col sm:flex-row sm:items-center justify-between gap-6 hover:shadow-md transition-shadow">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h3 className="text-lg font-bold text-gray-900">{service.name}</h3>
                      <div className="px-2 py-0.5 bg-gray-100 rounded text-xs text-gray-500 font-medium">{service.duration_min} min</div>
                    </div>
                    <p className="text-sm text-gray-500 max-w-xl">{service.description}</p>
                  </div>

                  <div className="flex items-center gap-6 border-t sm:border-t-0 sm:border-l border-gray-100 pt-4 sm:pt-0 sm:pl-6">
                    <div className="text-right min-w-[100px]">
                      <p className="text-sm text-gray-500">Total Price</p>
                      <p className="font-semibold text-gray-900">{fmtMoney(service.total_price_minor, service.currency_code)}</p>
                    </div>
                    <div className="text-right min-w-[100px]">
                      <p className="text-sm text-gray-500">Required Deposit</p>
                      <p className="font-bold text-primary-600">{fmtMoney(service.deposit_amount_minor, service.currency_code)}</p>
                    </div>
                    <div className="flex items-center gap-1">
                      <Button variant="ghost" className="p-2 text-gray-400 hover:text-blue-600"><Edit2 className="h-4 w-4" /></Button>
                      <Button variant="ghost" className="p-2 text-gray-400 hover:text-red-600"><Trash2 className="h-4 w-4" /></Button>
                    </div>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}

        {activeTab === 'CUSTOMERS' && (
          <div className="space-y-4">
            {loadingData ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">Loading customers...</div>
            ) : customers.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">No customer records for this workshop.</div>
            ) : (
              customers.map((customer) => (
                <Card key={customer.id} className="flex items-center justify-between gap-4">
                  <div>
                    <h3 className="font-semibold text-gray-900">{customer.name}</h3>
                    <p className="text-sm text-gray-500">{customer.email} • {customer.phone}</p>
                    {customer.notes && <p className="text-xs text-gray-400 mt-1">{customer.notes}</p>}
                  </div>
                  <div className="text-right text-xs text-gray-400">Added {fmtDate(customer.created_at)}</div>
                </Card>
              ))
            )}
          </div>
        )}

        {activeTab === 'VEHICLES' && (
          <div className="space-y-4">
            {loadingData ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">Loading vehicles...</div>
            ) : vehicles.length === 0 ? (
              <div className="text-center py-20 bg-white rounded-lg border border-dashed border-gray-300 text-gray-500">No vehicle records for this workshop.</div>
            ) : (
              vehicles.map((vehicle) => (
                <Card key={vehicle.id} className="flex items-center justify-between gap-4">
                  <div>
                    <h3 className="font-semibold text-gray-900">{vehicle.year} {vehicle.make} {vehicle.model}</h3>
                    <p className="text-sm text-gray-500">Owner: {vehicle.customer.name}</p>
                    <p className="text-xs text-gray-400 mt-1">{vehicle.color || 'Color not set'} • Plate: {vehicle.license_plate || 'N/A'}</p>
                  </div>
                  <div className="text-right text-xs text-gray-400">Added {fmtDate(vehicle.created_at)}</div>
                </Card>
              ))
            )}
          </div>
        )}

        {activeTab === 'JOBS' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            {['SCHEDULED', 'IN_PROGRESS', 'READY', 'DELIVERED'].map((status) => (
              <div key={status} className="bg-white border border-gray-200 rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-4">{status.replace('_', ' ')}</h3>
                <div className="space-y-3">
                  {jobs.filter((job) => job.status === status).length === 0 ? (
                    <div className="text-sm text-gray-400">No jobs in this column.</div>
                  ) : jobs.filter((job) => job.status === status).map((job) => (
                    <Card key={job.id} className="p-4">
                      <h4 className="font-medium text-gray-900">{job.title}</h4>
                      <p className="text-sm text-gray-500 mt-1">{job.customer.name} • {job.vehicle.year} {job.vehicle.make} {job.vehicle.model}</p>
                      <p className="text-xs text-gray-400 mt-2">Scheduled {fmtDate(job.scheduled_at)}</p>
                      {job.notes && <p className="text-xs text-gray-500 mt-2">{job.notes}</p>}
                    </Card>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'SETTINGS' && activeMembership && (
          <div className="max-w-4xl space-y-6">
            <Card className="p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                <Briefcase className="h-5 w-5 text-gray-400" />
                Workshop Profile
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input label="Workshop Name" defaultValue={activeMembership.business.name} />
                <Input label="URL Slug" defaultValue={activeMembership.business.slug} />
                <div className="col-span-1 md:col-span-2">
                  <Input label="Description" defaultValue={activeMembership.business.description} />
                </div>
                <Input label="Vertical" defaultValue={activeMembership.business.vertical} disabled className="bg-gray-50 text-gray-500" />
              </div>
            </Card>

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
                    <p className="font-bold text-green-900 text-sm">Workshop billing placeholder</p>
                    <p className="text-xs text-green-700">Subscription controls arrive in a later slice.</p>
                  </div>
                </div>
                <Button variant="outline" className="text-xs h-8 bg-white">Manage Billing</Button>
              </div>
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
