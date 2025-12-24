# AGENTS.md

## Overview

Blytz.Cloud is a cloud-based booking management solution for freelancers that enforces upfront payment before appointment confirmation. It's a React + TypeScript + Vite application with a multi-tenant SaaS architecture.

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
├── components/          # Reusable UI components
│   ├── Button.tsx
│   ├── Card.tsx
│   └── Input.tsx
├── screens/             # Page-level / view components
│   ├── Confirmation.tsx
│   ├── Login.tsx
│   ├── OperatorDashboard.tsx
│   ├── PublicBooking.tsx
│   └── SaaSLanding.tsx
├── App.tsx              # Main app router (view state management)
├── index.tsx            # Entry point
├── types.ts             # TypeScript interfaces and enums
├── constants.ts         # Mock data (businesses, services, slots, bookings)
├── package.json
├── tsconfig.json
├── vite.config.ts
└── index.html
```

## Code Patterns & Conventions

### Component Structure
- **Functional Components**: All components use React functional components with TypeScript
- **Props Interfaces**: Define props interfaces explicitly, either inline or in `types.ts`
- **Export Pattern**: Named exports for all components (`export const ComponentName: React.FC<Props> = ...`)

### State Management
- **Local State**: Use `useState` for component-local state
- **Derived State**: Use `useMemo` for filtering and computed values
- **View Routing**: The app uses a custom view state pattern in `App.tsx` with `ViewState` enum

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

### Mock Data Pattern
All mock data is centralized in `constants.ts`:
- `MOCK_BUSINESSES`: Business entities
- `MOCK_SERVICES`: Service offerings per business
- `MOCK_SLOTS`: Available time slots
- `MOCK_BOOKINGS`: Booking records

Mock data uses ISO date strings for timestamps and is filtered by `businessId` when needed.

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

### Files Created
- `Dockerfile`: Multi-stage build (Node build → Node serve with `serve` package)
- `docker-compose.yml`: Service definition with Traefik labels
- `.dockerignore`: Optimizes build context

### Build Process
1. **Builder Stage**: Uses Node.js 22 Alpine to install deps and run `npm run build`
2. **Production Stage**: Serves static files with `serve` package (lightweight Node.js static server)
3. **Output**: Static files in `/dist` served on port 80

### Dokploy Setup with Traefik
1. Connect your Git repository to Dokploy
2. Use `docker-compose.yml` as the deployment configuration
3. Update `your-domain.com` in Traefik labels to your actual domain
4. Set environment variable `GEMINI_API_KEY` in Dokploy's env vars section
5. Deploy - Dokploy will build and run the container
6. Traefik will automatically route traffic to the container

### Environment Variables
Only required env var:
```
GEMINI_API_KEY=your_key_here
```
Set this in Dokploy's environment variables section (not in `.env.local`).

## Development Notes

### Vite Configuration
- Dev server: `localhost:3000` with `host: 0.0.0.0`
- Environment variables: `GEMINI_API_KEY` is exposed via `process.env`
- Path resolution: `@` alias maps to project root

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

### Multi-Step Form Pattern
The `PublicBooking` component demonstrates a multi-step wizard pattern:
1. **SERVICE**: Select a service
2. **SLOT**: Select a time slot
3. **DETAILS**: Enter customer information
4. **PAYMENT**: Simulated payment flow

Each step has its own state and validation logic.
