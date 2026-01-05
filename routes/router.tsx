import React from 'react';
import { createBrowserRouter, RouterProvider, Outlet } from 'react-router-dom';
import { SaaSLanding } from '../screens/SaaSLanding';
import { PublicBooking } from '../screens/PublicBooking';
import { Confirmation } from '../screens/Confirmation';
import { Login } from '../screens/Login';
import { ForgotPassword } from '../screens/ForgotPassword';
import { ResetPassword } from '../screens/ResetPassword';
import { Availability } from '../screens/Availability';
import { OperatorDashboard } from '../screens/OperatorDashboard';
import { ProtectedRoute } from './ProtectedRoute';
import { AuthProvider } from '../context/AuthContext';
import { getSubdomain } from '../utils/subdomain';

const Layout = () => (
  <AuthProvider>
    <Outlet />
  </AuthProvider>
);

const isSubdomain = getSubdomain() !== null;

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      // Subdomain routes - Public booking only
      ...(isSubdomain ? [
        {
          index: true,
          element: <PublicBooking />,
        },
      ] : []),

      // Main domain routes - SaaS landing + Operator routes
      ...(!isSubdomain ? [
        {
          index: true,
          element: <SaaSLanding />,
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
        {
          path: 'availability',
          element: (
            <ProtectedRoute>
              <Availability />
            </ProtectedRoute>
          ),
        },
      ] : []),

      // Universal routes - accessible on both main domain and subdomains
      {
        path: 'confirmation',
        element: <Confirmation />,
      },
    ],
  },
]);
