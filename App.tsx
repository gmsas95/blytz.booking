import React from 'react';
import { RouterProvider } from 'react-router-dom';
import { router } from './routes/router';

const App: React.FC = () => {
  return (
    <div className="antialiased text-gray-900 bg-white">
      <RouterProvider router={router} />
    </div>
  );
};

export default App;
