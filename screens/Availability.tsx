import React, { useState, useEffect } from 'react';
import { Clock, Plus, Save, Trash2, Calendar as CalendarIcon, X } from 'lucide-react';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { useAuth } from '../context/AuthContext';
import { api, Business, BusinessAvailability } from '../api';

const DAYS = [
  { value: 0, name: 'Sunday' },
  { value: 1, name: 'Monday' },
  { value: 2, name: 'Tuesday' },
  { value: 3, name: 'Wednesday' },
  { value: 4, name: 'Thursday' },
  { value: 5, name: 'Friday' },
  { value: 6, name: 'Saturday' },
];

interface DayAvailability {
  dayOfWeek: number;
  startTime: string;
  endTime: string;
  isClosed: boolean;
}

export const Availability: React.FC = () => {
  const { logout } = useAuth();
  const [currentBusiness, setCurrentBusiness] = useState<Business | null>(null);
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [availability, setAvailability] = useState<BusinessAvailability[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  
  const [durationMin, setDurationMin] = useState(30);
  const [maxBookings, setMaxBookings] = useState(1);
  const [generating, setGenerating] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const businessesData = await api.getBusinesses();
      setBusinesses(businessesData);
      
      if (businessesData.length > 0) {
        const selectedBusiness = currentBusiness || businessesData[0];
        setCurrentBusiness(selectedBusiness);
        
        const [availabilityData] = await Promise.all([
          api.getAvailability(selectedBusiness.id),
        ]);
        
        setAvailability(availabilityData[0]);
        setDurationMin(selectedBusiness.slotDurationMin);
        setMaxBookings(selectedBusiness.maxBookings);
      }
    } catch (err) {
      console.error('Failed to fetch data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBusinessChange = (business: Business) => {
    setCurrentBusiness(business);
  };

  const handleSaveDay = async (dayIndex: number, updatedData: Partial<DayAvailability>) => {
    if (!currentBusiness) return;

    try {
      setSaving(true);
      await api.setAvailability(currentBusiness.id, {
        dayOfWeek: updatedData.dayOfWeek,
        startTime: updatedData.startTime,
        endTime: updatedData.endTime,
        isClosed: updatedData.isClosed,
      });
      await fetchData();
    } catch (err) {
      console.error('Failed to save availability:', err);
      alert('Failed to save availability. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleGenerateSlots = async () => {
    if (!currentBusiness) return;

    const startDate = new Date();
    const endDate = new Date();
    endDate.setMonth(endDate.getMonth() + 2); // Generate for next 2 months

    try {
      setGenerating(true);
      const slots = await api.generateSlots(currentBusiness.id, {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        durationMin,
      });
      alert(`Generated ${slots.length} slots successfully!`);
      await fetchData();
    } catch (err) {
      console.error('Failed to generate slots:', err);
      alert('Failed to generate slots. Please try again.');
    } finally {
      setGenerating(false);
    }
  };

  const getDayAvailability = (dayOfWeek: number): DayAvailability | undefined => {
    return availability.find(a => a.dayOfWeek === dayOfWeek);
  };

  const updateDayAvailability = (dayOfWeek: number, updates: Partial<DayAvailability>) => {
    setAvailability(prev => prev.map(a => 
      a.dayOfWeek === dayOfWeek ? { ...a, ...updates } : a
    ));
  };

  const getAvailabilityForDay = (dayOfWeek: number): string => {
    const avail = getDayAvailability(dayOfWeek);
    if (avail?.isClosed) {
      return <span className="text-red-600">Closed</span>;
    }
    if (!avail?.startTime || !avail?.endTime) {
      return <span className="text-gray-500">Not set</span>;
    }
    return <span className="text-green-600">
      {avail.startTime} - {avail.endTime}
    </span>;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Clock className="h-8 w-8 animate-spin text-gray-400 mx-auto mb-4" />
          <p className="text-gray-500">Loading availability...</p>
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
          <span className="font-bold text-lg tracking-tight">Availability</span>
        </div>

        <div className="p-4 flex-1">
          {businesses.length > 0 && (
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
          )}
        </div>

        <div className="p-4 border-t border-gray-100">
          <button onClick={logout} className="flex items-center gap-2 text-sm text-gray-500 hover:text-red-600 transition-colors w-full px-4 py-2">
            Sign Out
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto bg-gray-50 p-6">
        <div className="max-w-5xl mx-auto space-y-6">
          {/* Settings Header */}
          <Card className="p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Slot Settings</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Slot Duration (minutes)</label>
                <select
                  value={durationMin}
                  onChange={(e) => setDurationMin(parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value={15}>15 min</option>
                  <option value={30}>30 min</option>
                  <option value={45}>45 min</option>
                  <option value={60}>60 min</option>
                  <option value={90}>90 min</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Max Bookings per Slot</label>
                <input
                  type="number"
                  min={1}
                  value={maxBookings}
                  onChange={(e) => setMaxBookings(parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  placeholder="1"
                />
              </div>
              <div className="flex items-end">
                <Button 
                  onClick={handleGenerateSlots}
                  disabled={generating || !currentBusiness}
                  className="w-full"
                >
                  {generating ? <Clock className="h-4 w-4 animate-spin mr-2" /> : null}
                  Generate Slots
                </Button>
              </div>
            </div>
          </Card>

          {/* Weekly Schedule */}
          <Card className="p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-6 flex items-center justify-between">
              <span>Weekly Schedule</span>
              <span className="text-sm text-gray-500">Set working hours for each day</span>
            </h3>

            <div className="space-y-4">
              {DAYS.map((day) => {
                const avail = getDayAvailability(day.value);
                const isClosed = avail?.isClosed || false;

                return (
                  <div key={day.value} className="border border-gray-200 rounded-lg p-4">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-2">
                        <CalendarIcon className="h-5 w-5 text-gray-400" />
                        <span className="font-semibold text-gray-900">{day.name}</span>
                      </div>
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={isClosed}
                          onChange={(e) => {
                            updateDayAvailability(day.value, { isClosed: e.target.checked });
                            handleSaveDay(day.value, {
                              dayOfWeek: day.value,
                              isClosed: e.target.checked,
                            });
                          }}
                          className="w-4 h-4 text-gray-600 rounded"
                        />
                        <span className="text-sm text-gray-600">Closed</span>
                      </label>
                    </div>

                    {!isClosed && (
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <label className="block text-sm text-gray-600 mb-1">Start Time</label>
                          <input
                            type="time"
                            value={avail?.startTime || ''}
                            onChange={(e) => {
                              updateDayAvailability(day.value, { startTime: e.target.value });
                            }}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                            step="900"
                          />
                        </div>
                        <div>
                          <label className="block text-sm text-gray-600 mb-1">End Time</label>
                          <input
                            type="time"
                            value={avail?.endTime || ''}
                            onChange={(e) => {
                              updateDayAvailability(day.value, { endTime: e.target.value });
                            }}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                            step="900"
                          />
                        </div>
                      </div>
                    )}

                    {!isClosed && avail?.startTime && avail?.endTime && (
                      <div className="flex items-center justify-end gap-2 mt-3 pt-3 border-t border-gray-100">
                        <Button
                          onClick={() => handleSaveDay(day.value, getDayAvailability(day.value))}
                          disabled={saving}
                          variant="outline"
                          className="gap-2"
                        >
                          <Save className="h-4 w-4" />
                          Save
                        </Button>
                        <Button
                          onClick={() => {
                            updateDayAvailability(day.value, { startTime: '', endTime: '' });
                            handleSaveDay(day.value, getDayAvailability(day.value));
                          }}
                          disabled={saving}
                          variant="ghost"
                        >
                          <X className="h-4 w-4 text-gray-400" />
                        </Button>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </Card>
        </div>
      </main>
    </div>
  );
};
