# Blytz.Auto Documentation Index

> Last updated: 2026-04-08

This repo started as `blytz.booking`, a booking management prototype. The current direction is to evolve it into **Blytz.Auto**, a simpler SaaS for automotive workshops.

## Read This First In A Fresh Session

1. `docs/product.md`
2. `docs/mvp.md`
3. `docs/architecture.md`
4. `docs/implementation-plan.md`
5. `README.md`
6. `AGENTS.md`

## Current Product Direction

- product name: `Blytz.Auto`
- shape: multi-tenant SaaS
- deployment: single VPS first
- stack: Vite/React frontend + Go backend + Postgres
- billing: workshop pays platform subscription
- no payout/disbursement flow in MVP
- no per-tenant isolated infra in MVP

## Current Build Status

- Slice 1 booking-core hardening is complete in code
- Slice 2 workshop tenancy and membership enforcement is complete in core paths
- Slice 3 customers, vehicles, and jobs are implemented in protected workshop flows
- Slice 4 SaaS controls and final auth hardening are the active next slice before shipping

## Why This Exists

The repo contains useful booking-flow groundwork, but the old prototype docs are stale and mix prototype assumptions with partially implemented auth/service work. These docs define the new direction clearly so future sessions can continue with less confusion.
