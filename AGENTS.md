# Blytz.Cloud - Development Guide

## Project Overview

Blytz.Cloud is a **booking management prototype** for freelancers and service businesses. This is a **development/prototype stage** application - **NOT production ready**.

**⚠️ IMPORTANT: This application lacks critical production features including authentication, proper security, and complete business logic.**

### Tech Stack
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Backend**: Go + Gin framework + GORM
- **Database**: PostgreSQL + Redis (cache ready)
- **Deployment**: Docker + Docker Compose + Dokploy

## Current Status: PROTOTYPE (40% Complete)

### ✅ What's Implemented
- Basic UI components and booking flow
- Database schema and basic CRUD operations
- Docker containerization
- Multi-step booking wizard (UI only)
- Mock data system with API fallback

### ❌ Critical Missing Features
- **Authentication System**: No user registration, login, or session management
- **Security**: No input validation, rate limiting, or proper error handling
- **Business Logic**: No booking conflicts, payment processing, or notifications
- **Frontend Routing**: Not real routing - just state switching
- **Backend Architecture**: No service layer, direct DB access in handlers

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
│   ├── internal/       # Clean architecture attempt
│   │   ├── handlers/   # HTTP handlers (thin, but direct DB access)
│   │   ├── models/     # GORM models
│   │   └── repository/ # Database layer
│   └── config/         # Configuration
├── docker-compose.yml   # Multi-service orchestration
└── Dockerfile          # Frontend build
```

## Code Patterns & Issues

### Frontend Issues
1. **Fake Routing**: `App.tsx` uses view state switching, not real React Router
2. **Mixed API/Types**: `api.ts` contains both API client AND type definitions
3. **Mock Data Dependencies**: Components fall back to mock data instead of proper error handling
4. **No State Management**: All state in top-level component
5. **No Error Boundaries**: App crashes on API failures

### Backend Issues
1. **No Service Layer**: Handlers directly access database
2. **No Authentication**: Zero auth endpoints despite User model existing
3. **No Validation**: Only basic GORM model validation
4. **No Business Logic**: Just CRUD operations
5. **Poor Error Handling**: Generic error messages

### Critical Security Problems
- CORS allows all origins (`*`)
- No input sanitization
- No rate limiting
- No JWT implementation (secret exists but unused)
- No password hashing
- No protected routes

## API Endpoints (Basic CRUD Only)

```
GET  /health                    # Health check
GET  /api/v1/businesses         # List businesses
GET  /api/v1/businesses/:id     # Get business
GET  /api/v1/businesses/:id/services  # Get services
GET  /api/v1/businesses/:id/slots    # Get available slots
POST /api/v1/bookings           # Create booking
GET  /api/v1/businesses/:id/bookings # List bookings
```

**Missing endpoints:** Authentication, user management, payment, notifications, booking management

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
JWT_SECRET=your-secret-key  # Currently unused
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

### Phase 1: Authentication (2-3 weeks)
- [ ] Implement JWT authentication
- [ ] Add user registration/login endpoints
- [ ] Add password hashing
- [ ] Add auth middleware
- [ ] Add protected routes

### Phase 2: Architecture (2-3 weeks)
- [ ] Add service layer to backend
- [ ] Implement React Router for frontend
- [ ] Add proper error handling
- [ ] Add validation layers
- [ ] Add DTOs and business logic

### Phase 3: Business Logic (2-3 weeks)
- [ ] Implement real payment processing
- [ ] Add booking validation
- [ ] Add notification system
- [ ] Add conflict resolution
- [ ] Add proper status management

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