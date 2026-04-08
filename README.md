# Blytz.Cloud / Blytz.Auto - Booking SaaS Iteration Base

A cloud-based booking management prototype that is now being repurposed into **Blytz.Auto**, a simpler SaaS for automotive workshops. **This repo is not production ready yet.**

## Current Direction

Read these first:

- `docs/INDEX.md`
- `docs/product.md`
- `docs/mvp.md`
- `docs/architecture.md`
- `docs/implementation-plan.md`

## ⚠️ Important Notice

**This application lacks critical production features including authentication, proper security, and complete business logic. It should be considered a functional prototype/demo only.**

## Quick Start

### Prerequisites
- Node.js 18+
- PostgreSQL 16+
- Go 1.22+ (for backend development)
- Docker (optional, for containerized deployment)

### Local Development

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd blytz.booking
   ```

2. **Database setup:**
   ```bash
   # Create PostgreSQL database
   createdb blytz
   # Or use Docker: docker run -d --name postgres -e POSTGRES_DB=blytz -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16
   ```

3. **Backend setup:**
   ```bash
   cd backend
   go mod download
   go run cmd/server/main.go
   # Backend runs on http://localhost:8080
   ```

4. **Frontend setup:**
   ```bash
   # In root directory
   npm install
   npm run dev
   # Frontend runs on http://localhost:3000
   ```

### Docker Deployment (Development)

```bash
# Start all services
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f

# Stop everything
docker-compose down
```

## Project Status

### ✅ Implemented Features
- Basic multi-step booking flow (Service → Slot → Details → Payment)
- Business and service catalog management
- Time slot availability system
- Docker containerization
- Basic CRUD API endpoints
- Responsive UI with Tailwind CSS

### ❌ Critical Missing Features
- **Auth Hardening**: Registration enumeration hardening is still partial
- **Security**: Input validation and error handling still need more production hardening
- **Payment Processing**: Only simulation, no real payment integration
- **Business Logic**: Notifications and subscription enforcement are still missing

## Architecture Overview

### Frontend (React + TypeScript)
```
src/
├── components/     # Reusable UI components (Button, Card, Input)
├── screens/        # Page components (Login, Dashboard, Booking, etc.)
├── App.tsx         # Main app with fake routing via state switching
├── api.ts          # API client (mixed with types - needs cleanup)
├── types.ts        # TypeScript interfaces
└── constants.ts    # Mock data (deprecated)
```

### Backend (Go + Gin)
```
backend/
├── cmd/server/     # Application entry point
├── internal/
│   ├── handlers/   # HTTP handlers (direct DB access - needs service layer)
│   ├── models/     # GORM database models
│   └── repository/ # Database connection layer
└── config/         # Configuration management
```

## API Documentation

### Business Endpoints
```
GET  /api/v1/businesses              # List all businesses
GET  /api/v1/businesses/:id          # Get business details
GET  /api/v1/businesses/:id/services # Get services for business
GET  /api/v1/businesses/:id/slots    # Get available time slots
```

### Booking Endpoints
```
POST /api/v1/bookings                # Create new public booking
GET  /api/v1/businesses/:id/bookings # List bookings for business (auth + membership required)
```

### Health Check
```
GET  /health                         # Service health status
```

**Note:** Authentication uses httpOnly cookie sessions, and operator workshop access is membership scoped.

## Environment Configuration

### Frontend (.env.local)
```
VITE_API_URL=http://localhost:8080
```

### Backend Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=blytz
DB_SSLMODE=disable

# Server
SERVER_PORT=8080
ENV=development
AUTO_MIGRATE=true
SEED_DATA=false
BACKFILL_MONEY_FIELDS=true
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# JWT
JWT_SECRET=your-secret-key
JWT_COOKIE_NAME=blytz_session
JWT_COOKIE_SECURE=false
TRUSTED_PROXIES=127.0.0.1
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
VITE_API_URL=http://localhost:8080
APP_DOMAIN=localhost
API_DOMAIN=api.localhost
```

## Development Issues & Limitations

### Frontend Issues
- **Mixed Concerns**: API client and types in same file
- **No Error Boundaries**: App crashes on errors

### Backend Issues
- **No Validation**: Only basic GORM model validation
- **No Error Handling**: Generic error messages
- **Security Issues**: Registration enumeration hardening is still partial

### Security Problems
- No input validation or sanitization
- Registration still returns a distinct conflict status for existing emails
- Client-side active workshop preference is still stored in localStorage
- JWT logout/session revocation now invalidates prior session cookies

## Roadmap to Production

### Phase 1: Authentication (2-3 weeks)
- [ ] Implement JWT authentication system
- [ ] Add user registration and login endpoints
- [ ] Add password hashing with bcrypt
- [ ] Create auth middleware for protected routes
- [ ] Add user management (profile, password reset)

### Phase 2: Architecture (2-3 weeks)
- [ ] Add service layer to backend (separate business logic from handlers)
- [ ] Implement React Router for proper frontend routing
- [ ] Add comprehensive input validation
- [ ] Implement proper error handling and logging
- [ ] Add DTOs to separate API contracts from database models

### Phase 3: Business Logic (3-4 weeks)
- [ ] Integrate real payment processing (Stripe/PayPal)
- [ ] Add booking conflict detection and resolution
- [ ] Implement email/SMS notifications
- [ ] Add booking status management workflow
- [ ] Create admin dashboard for business management

### Phase 4: Security & Performance (1-2 weeks)
- [x] Implement auth rate limiting
- [ ] Add input sanitization and SQL injection prevention
- [x] Configure proper CORS policies
- [ ] Add request/response logging
- [ ] Implement database indexing optimization

## Release Checklist

- [ ] `JWT_SECRET` explicitly set
- [ ] `DB_PASSWORD` explicitly set
- [ ] `JWT_COOKIE_SECURE=true` in staging/production
- [ ] `TRUSTED_PROXIES` matches Dokploy/reverse-proxy network
- [ ] `CORS_ALLOWED_ORIGINS` matches deployed frontend origins
- [ ] `SEED_DATA=false`
- [ ] `AUTO_MIGRATE` reviewed for deploy target
- [ ] `go test ./...` passes
- [ ] `npm run build` passes

## Common Development Commands

```bash
# Frontend development
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build

# Backend development
cd backend
go run cmd/server/main.go              # Run server
go test ./...                          # Run tests (if any exist)
go mod tidy                            # Clean up dependencies

# Database operations
docker exec -it blytzbooking-postgres psql -U postgres -d blytz  # Access DB

# Docker operations
docker-compose up -d                 # Start all services
docker-compose logs -f backend       # View backend logs
docker-compose restart backend       # Restart backend service
```

## Contributing Guidelines

### When Adding Features
1. **Start with authentication** - Don't build on the mock login system
2. **Use proper routing** - Implement React Router for navigation
3. **Add service layers** - Separate business logic from HTTP handlers
4. **Include validation** - Both frontend and backend validation
5. **Add error handling** - Proper error messages and boundaries
6. **Write tests** - No test framework exists yet

### Code Style
- Frontend: Functional React components with TypeScript
- Backend: Standard Go conventions with clean architecture
- Database: GORM with PostgreSQL best practices
- Follow existing patterns but improve architecture

## Troubleshooting

### Common Issues

**Frontend won't start:**
```bash
# Check Node version
node --version  # Should be 18+
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

**Backend won't connect to database:**
```bash
# Check PostgreSQL is running
docker ps | grep postgres
# Verify connection string
cd backend && go run cmd/server/main.go
```

**API calls failing:**
- Check backend is running on port 8080
- Verify CORS settings in backend
- Check browser console for specific errors
- Frontend falls back to mock data silently

**Docker containers not starting:**
```bash
# Check port conflicts
netstat -tulpn | grep -E '(3000|8080|5432|6379)'
# Check Docker logs
docker-compose logs
```

---

## ⚠️ Final Warning

**This is a prototype demonstrating booking flow concepts. It should NOT be deployed to production without implementing the missing critical features, especially authentication, security, and proper business logic.**

For questions or to contribute to making this production-ready, please refer to the roadmap above and start with Phase 1: Authentication.
