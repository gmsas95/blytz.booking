# Blytz.Auto Implementation Plan

> Last updated: 2026-04-08

## Current Reality

This repo is useful as a base, but still contains prototype assumptions:

- stale docs
- remaining auth hardening work before shipping
- localStorage token storage still increases XSS blast radius
- rate limiting and auth enumeration defenses are still missing
- subscriptions/billing controls are not implemented yet

### Slice 1 status

Slice 1 is now complete in code:

1. Redis is no longer a required runtime dependency
2. service and booking money fields use integer minor units plus `currency_code`
3. booking creation returns conflict on slot race loss
4. CORS is config-driven instead of `*`
5. public booking no longer falls back to mock data silently
6. startup behavior is controlled by explicit flags:
   - `AUTO_MIGRATE`
   - `SEED_DATA`
   - `BACKFILL_MONEY_FIELDS`
   - `CORS_ALLOWED_ORIGINS`

## Recommended Build Order

### Slice 1 — harden the booking core ✅ complete

Goal: make the current prototype safer and more real before expanding features.

1. remove Redis as required runtime dependency
2. replace float money fields with integer minor units
3. make booking creation concurrency-safe
4. tighten CORS and config behavior
5. remove mock-data fallback from public booking flow

### Slice 2 — introduce workshop tenancy ✅ complete

1. add tenant/workshop ownership model
2. add membership model for operators
3. scope protected API access by tenant
4. reshape business language toward workshop domain
5. expose active workshop membership in auth context

### Slice 3 — add workshop entities ✅ complete in core paths

1. customers
2. vehicles
3. jobs/status
4. dashboard modules for those entities

### Slice 4 — add SaaS controls ← active

1. subscription model/state
2. subscription enforcement
3. billing UX for workshop owners

## Priority Rules

When choosing what to do next, prefer in this order:

1. correctness and safety
2. tenant isolation
3. production-path realism
4. workshop-specific domain support
5. visual polish

## File-Level Guidance

Likely next files to change:

- `backend/internal/auth/jwt.go`
- `backend/internal/handlers/handlers.go`
- `backend/internal/middleware/access.go`
- `backend/internal/services/auth_service.go`
- `backend/internal/services/job_service.go`
- `backend/internal/services/vehicle_service.go`
- `api.ts`
- `context/AuthContext.tsx`
- `screens/OperatorDashboard.tsx`
- `docs/architecture.md`
- `docs/mvp.md`
- `README.md`

## Verification Checklist For Future Sessions

Before declaring meaningful implementation work done:

1. run backend tests
2. run frontend build
3. re-read modified files to confirm changes landed
4. verify protected workshop routes reject non-members with `403`
5. verify `/auth/me` returns membership/workshop context
6. manually exercise booking flow
7. verify booking conflict path under repeated booking attempts
8. verify old float money fields are gone where intended
9. verify stale Slice 1 doc claims are removed
10. verify tenant-scoped customers, vehicles, and jobs cannot cross workshops
11. verify JWT secret is explicitly configured in the target environment

## Fresh-Session Prompt Starter

Use this in a new conversation if needed:

> Read `docs/INDEX.md`, `docs/product.md`, `docs/mvp.md`, `docs/architecture.md`, and `docs/implementation-plan.md`. Then inspect the current codebase and implement Slice 4 only: add subscription state and enforcement, improve auth hardening (rate limiting, safer session handling), and verify with backend tests plus frontend build.
