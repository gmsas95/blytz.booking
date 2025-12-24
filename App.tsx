import React, { useState } from 'react';
import { ViewState, Business } from './types';
import { SaaSLanding } from './screens/SaaSLanding';
import { PublicBooking } from './screens/PublicBooking';
import { Confirmation } from './screens/Confirmation';
import { Login } from './screens/Login';
import { OperatorDashboard } from './screens/OperatorDashboard';

const App: React.FC = () => {
  const [currentView, setCurrentView] = useState<ViewState>(ViewState.SAAS_LANDING);
  const [selectedBusiness, setSelectedBusiness] = useState<Business | null>(null);
  const [lastBooking, setLastBooking] = useState<any>(null);

  // Router Actions
  const handleBusinessSelect = (biz: Business) => {
    setSelectedBusiness(biz);
    setCurrentView(ViewState.PUBLIC_BOOKING);
  };

  const handleBookingComplete = (bookingData: any) => {
    setLastBooking(bookingData);
    setCurrentView(ViewState.CONFIRMATION);
  };

  const handleOperatorLogin = () => {
    setCurrentView(ViewState.LOGIN);
  };

  const handleLoginSuccess = () => {
    setCurrentView(ViewState.DASHBOARD);
  };

  const handleLogout = () => {
    setCurrentView(ViewState.SAAS_LANDING);
    setSelectedBusiness(null);
  };

  const renderView = () => {
    switch (currentView) {
      case ViewState.SAAS_LANDING:
        return (
          <SaaSLanding 
            onSelectBusiness={handleBusinessSelect} 
            onOperatorLogin={handleOperatorLogin}
          />
        );
      case ViewState.PUBLIC_BOOKING:
        if (!selectedBusiness) return <SaaSLanding onSelectBusiness={handleBusinessSelect} onOperatorLogin={handleOperatorLogin} />;
        return (
          <PublicBooking 
            business={selectedBusiness}
            onComplete={handleBookingComplete} 
            onExit={() => setCurrentView(ViewState.SAAS_LANDING)}
          />
        );
      case ViewState.CONFIRMATION:
        return (
          <Confirmation 
            bookingDetails={lastBooking} 
            onDone={() => setCurrentView(ViewState.SAAS_LANDING)} 
          />
        );
      case ViewState.LOGIN:
        return (
          <Login onLogin={handleLoginSuccess} />
        );
      case ViewState.DASHBOARD:
        return (
          <OperatorDashboard onLogout={handleLogout} />
        );
      default:
        return <SaaSLanding onSelectBusiness={handleBusinessSelect} onOperatorLogin={handleOperatorLogin} />;
    }
  };

  return (
    <div className="antialiased text-gray-900 bg-white">
      {renderView()}
    </div>
  );
};

export default App;
