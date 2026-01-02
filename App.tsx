import React from 'react';
import { RouterProvider } from 'react-router-dom';
import { router } from './routes/router';
import { ErrorBoundary } from './components/ErrorBoundary';

const App: React.FC = () => {
  return (
    <ErrorBoundary>
      <div className="antialiased text-gray-900 bg-white">
        <RouterProvider router={router} />
      </div>
    </ErrorBoundary>
  );
};

export default App;

