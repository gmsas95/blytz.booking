import React from 'react';
import { createBrowserRouter, RouterProvider, Outlet } from 'react-router-dom';
import { SaaSLanding } from '../screens/SaaSLanding';
import { PublicBooking } from '../screens/PublicBooking';
import { Confirmation } from '../screens/Confirmation';
import { Login } from '../screens/Login';
import { ForgotPassword } from '../screens/ForgotPassword';
import { ResetPassword } from '../screens/ResetPassword';
import { OperatorDashboard } from '../screens/OperatorDashboard';
import { ProtectedRoute } from './ProtectedRoute';
import { AuthProvider } from '../context/AuthContext';

const Layout = () => (
  <AuthProvider>
    <Outlet />
  </AuthProvider>
);

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <SaaSLanding />,
      },
      {
        path: 'business/:slug',
        element: <PublicBooking />,
      },
      {
        path: 'confirmation',
        element: <Confirmation />,
      },
       {
         path: 'login',
         element: <Login />,
       },
       {
         path: 'forgot-password',
         element: <ForgotPassword />,
       },
       {
         path: 'reset-password',
         element: <ResetPassword />,
       },
      {
        path: 'dashboard',
        element: (
          <ProtectedRoute>
            <OperatorDashboard />
          </ProtectedRoute>
        ),
      },
    ],
  },
]);
