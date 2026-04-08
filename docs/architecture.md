# Blytz.Auto Architecture

> Last updated: 2026-04-08

## Architecture Summary

Use a **simple monolith with clear internal boundaries**.

### Stack

- frontend: React + Vite
- routing/state: React Router now; TanStack can be introduced incrementally later if justified
- backend: Go + Gin
- database: Postgres
- cache/queue: Redis optional later, not required for MVP
- deployment: single VPS

## Why This Architecture

This repo already contains a usable Vite + Go + Postgres prototype. The right move is to **iterate and harden**, not rebuild into a more complex platform.

This architecture is chosen for:

- low solo-maintenance burden
- easier deployment and debugging
- simpler security review
- faster path to first paying workshop

## System Shape

### Frontend

One frontend app serves:

- landing page
- login/register
- public workshop booking pages
- workshop dashboard

Suggested route shape:

- `/`
- `/pricing`
- `/login`
- `/register`
- `/w/:slug`
- `/dashboard`
- `/dashboard/bookings`
- `/dashboard/customers`
- `/dashboard/vehicles`
- `/dashboard/services`
- `/dashboard/jobs`
- `/dashboard/settings`

### Backend

One Go API serves:

- auth
- tenant/workshop management
- services
- bookings
- customers
- vehicles
- jobs
- subscriptions

Keep clear internal boundaries:

- handlers: HTTP transport only
- services: business logic
- models: persistence structures
- auth: token/password utilities
- repository: DB wiring

## Tenancy Model

Use **shared-table multitenancy** for MVP.

Current implementation note:

- the current codebase still uses `business_id` as the effective tenant scope key
- Slice 2 adds memberships and enforces workshop-scoped access using that key first
- a future slice can rename `business` terminology more aggressively if needed
- Slice 3 adds tenant-scoped customers, vehicles, and jobs on the same key

### Rules

- each workshop is a tenant
- each operator belongs to one or more tenants via membership
- tenant-owned entities must carry the effective tenant scope key (`business_id` in the current codebase)
- all protected reads/writes must verify tenant access

## Data Model Direction

Core entities now in the codebase:

- tenants/workshops
- users
- memberships
- workshop_services
- availability_slots
- customers
- vehicles
- bookings
- jobs

Still planned:

- notes
- subscriptions
- payment_events

## Booking Concurrency

Do not rely on a simple read-then-write check.

Current status:

- booking creation now uses a transaction plus atomic slot reservation path
- API returns conflict when the slot is already taken

Recommended MVP approach:

1. customer selects a discrete slot
2. backend begins a transaction
3. backend locks or atomically updates the slot row
4. backend inserts the booking
5. backend commits
6. if the slot was already taken, return a conflict response

## Security Baseline

- restrict CORS to configured frontend origins
- validate request DTOs explicitly
- do not trust only frontend route guards
- protect every tenant-owned endpoint with auth + tenant checks
- store money in integer minor units
- remove mock fallback behavior from production paths
- avoid startup auto-seeding in production mode

Current status:

- CORS is now driven by configured allowed origins
- service and booking money fields now use minor units plus `currency_code`
- public booking no longer silently falls back to mock data
- startup behavior is now controlled by flags instead of always migrating/seeding
- operator bookings/customers/vehicles/jobs now require auth + workshop membership
- JWT secret must now be explicitly configured at startup
- auth sessions use httpOnly cookies with server-side token-version revocation
- login/register are rate limited with both per-IP and IP+email limits for the current single-instance VPS model

Remaining ship blockers:

- auth now uses httpOnly cookie sessions, but registration still returns a distinct conflict status
- subscription enforcement is not implemented yet

## Deployment Baseline

Start with:

- VPS
- reverse proxy
- frontend build served statically or via separate container
- Go API process/container
- Postgres

Redis is optional and should only be reintroduced if a real need appears.
