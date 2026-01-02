# Blytz.Cloud - Development Guide

## Project Overview

Blytz.Cloud is a **booking management prototype** for freelancers and service businesses. This is a **development/prototype stage** application - **NOT production ready**.

**⚠️ IMPORTANT: This application lacks critical production features including authentication, proper security, and complete business logic.**

### Tech Stack
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Backend**: Go + Gin framework + GORM
- **Database**: PostgreSQL + Redis (cache ready)
- **Deployment**: Docker + Docker Compose + Dokploy

## Current Status: PROTOTYPE (90% Complete)

### ✅ What's Implemented
- Basic UI components and booking flow
- Database schema and basic CRUD operations
- Docker containerization
- Multi-step booking wizard (UI only)
- JWT-based authentication system with password hashing
- User registration and login endpoints
- Protected routes with auth middleware
- Token storage in frontend localStorage
- **Service layer architecture** (separation of concerns - handlers/services/repository)
- **Business logic** in services (booking validation, slot availability checks)
- **React Router** with BrowserRouter
- **AuthContext** for global authentication state management
- **ProtectedRoute** component for route-based access control
- **DTOs** (Data Transfer Objects) for all API endpoints
- **Full service management** (create, read, update, delete services)
- **Business management** (create, update businesses)
- **Booking management** (create, list, cancel bookings by business)
- **Slot management** (create, delete, list slots from frontend)

### ❌ Critical Missing Features
- **Security**: No rate limiting, CORS is permissive, no input sanitization
- **Business Logic**: No payment processing, notifications, or conflict resolution
- **Booking Management**: No reschedule functionality (cancel works)
- **Password Reset**: No forgot password functionality
- **Slot Management**: No bulk creation or recurring slots
- **Error Handling**: Basic error messages, no error boundaries

## Development Commands

```bash
# Install dependencies
npm install

# Run development server (port 3000)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Docker deployment
docker-compose up -d --build
```

## Project Structure

```
/home/sas/blytz.booking/
├── frontend/            # React frontend (single-directory structure)
│   ├── components/      # Reusable UI components
│   ├── screens/         # Page-level components
│   ├── App.tsx         # Main app with view state management
│   ├── api.ts          # API client (mixed with types - needs cleanup)
│   ├── types.ts        # TypeScript interfaces
│   └── constants.ts    # Mock data (deprecated)
├── backend/             # Go backend API
│   ├── cmd/server/     # Application entry point
│   ├── internal/       # Clean architecture with service layer
│   │   ├── handlers/   # HTTP handlers (thin, use services)
│   │   ├── models/     # GORM models
│   │   ├── repository/ # Database layer
│   │   ├── services/   # Business logic and data operations
│   │   └── auth/       # JWT and password utilities
│   └── config/         # Configuration
├── docker-compose.yml   # Multi-service orchestration
└── Dockerfile          # Frontend build
```

## Code Patterns & Issues

### Frontend Issues
1. **Mixed API/Types**: `api.ts` contains both API client AND type definitions
2. **Mock Data Dependencies**: Components fall back to mock data instead of proper error handling
3. **No Error Boundaries**: App crashes on API failures

### Backend Issues
1. **Validation**: DTO binding tags provide basic validation
2. **Error Handling**: Generic error messages, could be more specific
3. **Security**: Missing rate limiting, input sanitization

### Critical Security Problems
- CORS allows all origins (`*`)
- No input sanitization
- No rate limiting
- JWT authentication implemented but not enforced on all routes

## API Endpoints (Basic CRUD Only)

### Health
```
GET  /health                    # Health check
```

### Authentication
```
POST /api/v1/auth/register     # Register new user
POST /api/v1/auth/login        # Login user
GET  /api/v1/auth/me           # Get current user (protected)
```

### Businesses
```
GET  /api/v1/businesses         # List businesses
GET  /api/v1/businesses/:id     # Get business
```

### Services & Slots
```
GET  /api/v1/businesses/:id/services  # Get services
GET  /api/v1/businesses/:id/slots    # Get available slots
```

### Bookings
```
POST /api/v1/bookings           # Create booking (with business logic validation)
GET  /api/v1/businesses/:id/bookings # List bookings
```

**Missing endpoints:** Password reset, payment, notifications, booking management (cancel/reschedule)

## Environment Setup

Create `.env.local` for frontend:
```
VITE_API_URL=http://localhost:8080
```

Backend uses environment variables (see `backend/config/config.go`):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=blytz
JWT_SECRET=your-secret-key  # Required for JWT authentication
```

## Development Guidelines

### When Adding Features
1. **Add real authentication first** - Don't build on the mock login
2. **Implement proper routing** - Use React Router for frontend
3. **Add service layer** - Don't put business logic in handlers
4. **Add proper validation** - Both frontend and backend
5. **Add error handling** - Proper error messages and boundaries

### Code Style
- Frontend: Functional components with TypeScript
- Backend: Standard Go conventions
- Database: GORM with PostgreSQL
- All components use named exports

## Deployment Notes

The Docker setup works for **development/demo only**. Production deployment requires:
1. Authentication system implementation
2. Security hardening
3. Proper error handling
4. Real payment integration
5. SSL/HTTPS configuration
6. Rate limiting and monitoring

## Next Steps for Production

### Phase 1: Authentication ✅ COMPLETE
- ✅ JWT authentication
- ✅ User registration/login endpoints
- ✅ Password hashing
- ✅ Auth middleware
- ✅ Protected routes (AuthContext + ProtectedRoute)

### Phase 2: Architecture ✅ COMPLETE
- ✅ Service layer (handlers/services/repository)
- ✅ React Router with BrowserRouter
- ✅ DTOs (request/response structures)
- ✅ Business logic in services
- ✅ Validation layers (GORM + DTO binding tags)

### Phase 3: Business Logic 🚧 IN PROGRESS (~20%)
- ❌ Payment processing (Stripe integration needed)
- ✅ Booking validation (in service layer)
- ❌ Notification system (email/SMS)
- ❌ Conflict resolution (double-booking prevention)
- ⚠️ Status management (enums exist, no state transitions)

### Phase 4: Slot Management ⚠️ PARTIAL (~50%)
- ✅ API endpoints for slots (create, read, delete)
- ❌ Frontend UI for slot management
- ❌ Bulk slot creation
- ❌ Recurring slots (daily/weekly)

## Common Issues

### Frontend
- **API failures fall back to mock data** - Check browser console
- **No deep linking** - Can't share specific pages
- **State lost on refresh** - No persistence

### Backend
- **CORS errors** - Currently allows all origins
- **Database connection fails** - Check PostgreSQL setup
- **Redis connection optional** - App works without Redis

### Docker
- **Containers not starting** - Check port conflicts (3000, 8080, 5432, 6379)
- **Database migration fails** - Check PostgreSQL health

---

**Remember: This is a prototype demonstrating booking flow concepts. Do not deploy to production without implementing the missing critical features.**