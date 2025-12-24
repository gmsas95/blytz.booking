# AGENTS.md

## Overview

Blytz.Cloud is a cloud-based booking management solution for freelancers that enforces upfront payment before appointment confirmation. Full-stack application with:

- **Frontend**: React + TypeScript + Vite
- **Backend**: Go with Gin framework + GORM
- **Database**: PostgreSQL
- **Cache**: Redis

## Commands

### Development
```bash
npm install           # Install dependencies
npm run dev          # Start dev server on port 3000, host 0.0.0.0
```

### Build & Preview
```bash
npm run build        # Build for production
npm run preview      # Preview production build locally
```

### Docker (Dokploy)
```bash
docker-compose up -d --build    # Build and start container
docker-compose down             # Stop and remove container
docker-compose logs -f          # View logs
```

### Environment Setup
The app requires a `.env.local` file with:
```
GEMINI_API_KEY=your_gemini_api_key_here
```

## Project Structure

```
/home/gmsas95/blytz.booking/
├── frontend/            # React + TypeScript frontend
│   ├── components/      # Reusable UI components
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   └── Input.tsx
│   ├── screens/         # Page-level / view components
│   │   ├── Confirmation.tsx
│   │   ├── Login.tsx
│   │   ├── OperatorDashboard.tsx
│   │   ├── PublicBooking.tsx
│   │   └── SaaSLanding.tsx
│   ├── App.tsx         # Main app router (view state management)
│   ├── index.tsx       # Entry point
│   ├── api.ts          # API client for backend communication
│   ├── types.ts        # TypeScript interfaces and enums
│   ├── constants.ts    # Mock data (deprecated, use API)
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── index.html
├── backend/             # Go backend API
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Application entry point
│   ├── internal/
│   │   ├── handlers/              # HTTP request handlers
│   │   │   └── handlers.go
│   │   ├── models/                # Data models
│   │   │   └── models.go
│   │   └── repository/            # Data access layer
│   │       └── repository.go
│   ├── config/
│   │   └── config.go             # Configuration management
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml   # Multi-service orchestration
├── .env.example       # Environment variables template
└── AGENTS.md         # This file
```

## Code Patterns & Conventions

### Go Backend Patterns
- **Architecture**: Clean architecture with handlers, models, repository separation
- **Framework**: Gin HTTP router
- **ORM**: GORM for database operations
- **Configuration**: Environment-based config with struct-based loading
- **API Design**: RESTful with `/api/v1/` prefix
- **Naming**: Go conventions - CamelCase exported, camelCase unexported

### Frontend Component Structure
- **Functional Components**: All components use React functional components with TypeScript
- **Props Interfaces**: Define props interfaces explicitly, either inline or in `types.ts`
- **Export Pattern**: Named exports for all components (`export const ComponentName: React.FC<Props> = ...`)

### State Management
- **Frontend**: Local state with `useState`, derived state with `useEffect` for API calls
- **Backend**: Database-driven state with GORM, optional Redis caching
- **API Communication**: Centralized API client in `api.ts`
- **View Routing**: Custom view state pattern in `App.tsx` with `ViewState` enum

### Naming Conventions
- **Components**: PascalCase (`PublicBooking`, `OperatorDashboard`)
- **Interfaces**: PascalCase, no `I` prefix (`Business`, `Service`, `Booking`)
- **Enums**: PascalCase with UPPER_CASE values (`BookingStatus.PENDING`)
- **Constants**: UPPERCASE with `MOCK_` prefix for mock data (`MOCK_BUSINESSES`)
- **Event Handlers**: `handle` prefix (`handleServiceSelect`, `handleBookingComplete`)
- **Callback Props**: `on` prefix (`onSelectBusiness`, `onComplete`, `onLogin`)

### Styling
- **Tailwind CSS**: Utility-first CSS framework
- **Color Scheme**: Uses `primary-*`, `gray-*`, `zinc-*` color scale
- **Responsive**: Mobile-first design with `sm:`, `md:` breakpoints
- **Icons**: `lucide-react` for all icons

### TypeScript Configuration
- **Target**: ES2022
- **Module Resolution**: bundler
- **JSX**: react-jsx
- **Path Alias**: `@/` maps to project root (e.g., `@/components/Button`)
- **Import Extensions**: `.ts` and `.tsx` extensions allowed in imports

### Mock Data Pattern (Frontend)
Mock data is centralized in `constants.ts` but **deprecated** - use API instead:
- `MOCK_BUSINESSES`: Business entities (deprecated)
- `MOCK_SERVICES`: Service offerings per business (deprecated)
- `MOCK_SLOTS`: Available time slots (deprecated)
- `MOCK_BOOKINGS`: Booking records (deprecated)

New components use the API client in `api.ts` for real data.

### Database Seeding (Backend)
Backend automatically seeds initial data on startup in `repository.SeedData()`:
- 3 businesses (Automotive, Wellness, Creative)
- 4 services (2 per business)
- Slots are generated dynamically or seeded via SQL

### View Routing
The app uses a custom view state machine rather than a traditional router:

```typescript
enum ViewState {
  SAAS_LANDING = 'SAAS_LANDING',
  PUBLIC_BOOKING = 'PUBLIC_BOOKING',
  CONFIRMATION = 'CONFIRMATION',
  LOGIN = 'LOGIN',
  DASHBOARD = 'DASHBOARD'
}
```

Navigation happens by changing the `currentView` state and passing data through props/callbacks.

### Form Handling
- Use `<form>` with `onSubmit` handlers
- Prevent default with `e.preventDefault()`
- Controlled inputs with `value` and `onChange`
- Validate before state transitions

### Type Definitions (types.ts)
Key interfaces:
- `Business`: Business entity with branding (id, name, slug, vertical, description, themeColor)
- `Service`: Service offering (id, businessId, name, description, durationMin, totalPrice, depositAmount)
- `Slot`: Time slot (id, businessId, startTime, endTime, isBooked)
- `CustomerDetails`: Customer info (name, email, phone)
- `Booking`: Booking record with status tracking
- `BookingStatus`: PENDING | CONFIRMED | COMPLETED | CANCELLED

## Important Gotchas

### Dynamic Theme Colors
The app supports dynamic brand colors per business. Since Tailwind doesn't support JIT color interpolation without a safelist, inline styles are used for business-specific colors:

```typescript
const brandStyle = { color: business.themeColor };
```

### Fixed Positioning
Mobile booking flow uses fixed positioning for action buttons at the bottom:
```typescript
className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 safe-area-pb z-20"
```

The `safe-area-pb` class is used for iOS safe area handling.

### Date/Time Formatting
Helpers are defined inline per component for date/time formatting:
```typescript
const fmtMoney = (n: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);
const fmtDate = (iso: string) => new Intl.DateTimeFormat('en-US', { ... }).format(new Date(iso));
const fmtTime = (iso: string) => new Intl.DateTimeFormat('en-US', { ... }).format(new Date(iso));
```

### No Traditional Router
This app uses a view state pattern rather than React Router. All state lives in `App.tsx` and flows down through props.

### Icon Library
Icons are imported from `lucide-react`. Always check if the icon exists before using.

### No Testing Framework
Currently no tests, no test files, no test runner configured. This is a demo/prototype codebase.

## Docker & Dokploy Deployment

### Full Stack Services
The `docker-compose.yml` orchestrates 4 services:

1. **PostgreSQL**: Database (port 5432)
2. **Redis**: Cache (port 6379)
3. **Backend**: Go API (port 8080)
4. **Frontend**: React app served with `serve` (port 80)

### Docker Files Created
- **Frontend Dockerfile**: Multi-stage build (Node build → Node serve)
- **Backend Dockerfile**: Multi-stage build (Go build → Alpine binary)
- **docker-compose.yml**: Full-stack orchestration with Traefik labels
- **.dockerignore**: Optimizes build context (both frontend and backend)

### Build Process
**Frontend**:
1. **Builder Stage**: Node.js 22 Alpine, installs deps, runs `npm run build`
2. **Production Stage**: Serves static files with `serve` package

**Backend**:
1. **Builder Stage**: Golang 1.22, downloads deps, compiles binary
2. **Production Stage**: Alpine with compiled binary

### Dokploy Setup with Traefik
1. Connect your Git repository to Dokploy
2. Use `docker-compose.yml` as the deployment configuration
3. Update `your-domain.com` in Traefik labels to your actual domain
4. Set environment variables:
   - `DB_USER`, `DB_PASSWORD`, `DB_NAME` for PostgreSQL
   - `JWT_SECRET` for authentication
   - `GEMINI_API_KEY` for frontend (if using AI features)
5. Deploy - Dokploy will build and run all containers
6. Traefik will automatically route:
   - `your-domain.com` → Frontend
   - `your-domain.com/api` → Backend

### Environment Variables
Copy `.env.example` to `.env` and configure:

```bash
# Database
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=blytz

# JWT
JWT_SECRET=your-super-secret-key

# Frontend (development)
VITE_API_URL=http://localhost:8080
```

In production/Dokploy, set these via the Dokploy UI.

## Development Notes

### Vite Configuration (Frontend)
- Dev server: `localhost:3000` with `host: 0.0.0.0`
- Environment variables: `VITE_API_URL` for backend URL
- Path resolution: `@` alias maps to project root
- Proxy: In production, Traefik handles `/api` routing

### Go Configuration (Backend)
- Server: Gin router on port 8080
- Database: PostgreSQL with GORM (auto-migration on startup)
- Redis: Optional caching layer (currently unused)
- CORS: Enabled for all origins (configure for production)
- Seeding: Automatically seeds initial business/service data

### TypeScript Strictness
- `skipLibCheck: true` for faster builds
- `noEmit: true` - Vite handles compilation
- `allowImportingTsExtensions: true` - Can import `.ts/.tsx` files directly

### Component Props Patterns
Components use `React.HTMLAttributes` extensions for native element props:

```typescript
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  fullWidth?: boolean;
  isLoading?: boolean;
}
```

This allows passing standard HTML attributes (`disabled`, `className`, `onClick`, etc.) to components.

## API Architecture

### Backend API Endpoints
Base URL: `http://localhost:8080` (or via Traefik)

**Health**
- `GET /health` - Health check endpoint

**Businesses**
- `GET /api/v1/businesses` - List all businesses
- `GET /api/v1/businesses/:id` - Get business details

**Services**
- `GET /api/v1/businesses/:businessId/services` - List services for a business

**Slots**
- `GET /api/v1/businesses/:businessId/slots` - Get available slots (not booked)

**Bookings**
- `POST /api/v1/bookings` - Create a new booking
- `GET /api/v1/businesses/:businessId/bookings` - List bookings for a business

### API Client (Frontend)
Frontend uses centralized API client in `api.ts`:
```typescript
import { api } from './api';

// Get all businesses
const businesses = await api.getBusinesses();

// Get services for a business
const services = await api.getServicesByBusiness(businessId);

// Create booking
const booking = await api.createBooking({
  business_id,
  service_id,
  slot_id,
  // ...
});
```

### Data Models
**Business**: id, name, slug, vertical, description, theme_color
**Service**: id, business_id, name, description, duration_min, total_price, deposit_amount
**Slot**: id, business_id, start_time, end_time, is_booked
**Booking**: id, business_id, service_id, slot_id, service_name, slot_time, customer, status, deposit_paid, total_price

## Multi-Step Form Pattern

```typescript
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  fullWidth?: boolean;
  isLoading?: boolean;
}
```

This allows passing standard HTML attributes (`disabled`, `className`, `onClick`, etc.) to components.

### Multi-Step Form Pattern
The `PublicBooking` component demonstrates a multi-step wizard pattern:
1. **SERVICE**: Select a service
2. **SLOT**: Select a time slot
3. **DETAILS**: Enter customer information
4. **PAYMENT**: Simulated payment flow

Each step has its own state and validation logic.
