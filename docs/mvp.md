# Blytz.Auto MVP Scope

> Last updated: 2026-04-08

## MVP Goal

Ship a production-lean but simple SaaS that one workshop can realistically pay for.

## In Scope

### Public Side

- landing page
- pricing page
- workshop booking page by slug

### Workshop App

- auth
- workshop profile/settings
- service catalog
- customer records
- vehicle records
- bookings list
- public booking intake
- simple job status board
- subscription status visibility

### Platform Rules

- single VPS deployment first
- single Postgres database
- shared-table multitenancy with strict tenant scoping
- subscription-only billing

## Explicitly Out Of Scope

- customer online payments
- payout/disbursement to workshop bank accounts
- Stripe Connect / marketplace money movement
- per-tenant infrastructure isolation
- custom domains in MVP
- WhatsApp automation
- inventory, accounting, parts, procurement
- advanced analytics/reporting
- multi-branch support

## Required Engineering Behaviors

### Tenant isolation

Every tenant-owned record must be scoped by the tenant key and every authenticated query must enforce tenant ownership.

Current implementation note:

- the current codebase uses `business_id` as the effective workshop/tenant scope key
- Slice 2 adds memberships and route-level tenant authorization on top of that
- Slice 3 adds workshop-scoped customers, vehicles, and jobs in both API and dashboard views

### Booking safety

Slot booking must be concurrency-safe.

Preferred MVP approach:

- discrete slots for predictable-duration services
- Postgres transaction
- row lock / safe update path
- unique/consistency constraints where possible

### Money safety

Do not use floating point for money.

Use:

- `amount_minor`
- `price_minor`
- `deposit_minor`
- separate `currency_code`

## Success Criteria

The MVP is successful when:

1. a workshop can sign up and manage its own workspace
2. a customer can book a service on the public page
3. the booking cannot double-book the same slot under normal concurrent requests
4. workshop staff can view workshop-scoped customers, vehicles, bookings, and jobs
5. the app can be deployed and maintained by one person on a VPS

## Remaining MVP Hardening Before Shipping

- explicit JWT secret in every environment
- auth enumeration hardening on registration
- subscription state and enforcement for paid workshops
