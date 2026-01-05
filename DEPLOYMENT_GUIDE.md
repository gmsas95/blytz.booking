# Quick Deployment Guide for blytz.cloud

## Current Environment Status

✅ Backend compiles successfully
✅ Frontend builds successfully
✅ All subdomain middleware files created
✅ Concurrency fix implemented

---

## Required Environment Variables (Add These)

### Backend (Dokploy Backend Service)

Add these to your backend environment in Dokploy:

```bash
BASE_DOMAIN=blytz.cloud
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=5
```

### Frontend (Dokploy Frontend Service)

Add this to your frontend environment in Dokploy:

```bash
VITE_BASE_DOMAIN=blytz.cloud
```

---

## Deployment Steps

### Option 1: Via Dokploy UI (Recommended)

1. **Backend Service**
   - Go to your backend app in Dokploy
   - Settings → Environment Variables
   - Add `BASE_DOMAIN=blytz.cloud`
   - Add `DB_MAX_OPEN_CONNS=30`
   - Add `DB_MAX_IDLE_CONNS=5`
   - Click Save
   - Click Redeploy

2. **Frontend Service**
   - Go to your frontend app in Dokploy
   - Settings → Environment Variables
   - Add `VITE_BASE_DOMAIN=blytz.cloud`
   - Click Save
   - Click Redeploy

### Option 2: Via Docker Compose

1. **Create `.env` file in project root:**

```bash
cd /home/sas/blytz.booking
cat > .env << 'EOF'
BASE_DOMAIN=blytz.cloud
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=5
VITE_BASE_DOMAIN=blytz.cloud
EOF
```

2. **Redeploy:**

```bash
docker-compose down
docker-compose up -d --build
```

### Option 3: Via Docker CLI (Dokploy)

```bash
# Stop services
docker stop blytz-booking-backend blytz-booking-frontend

# Add env vars to backend
docker update \
  --env-add BASE_DOMAIN=blytz.cloud \
  --env-add DB_MAX_OPEN_CONNS=30 \
  --env-add DB_MAX_IDLE_CONNS=5 \
  blytz-booking-backend

# Add env var to frontend
docker update \
  --env-add VITE_BASE_DOMAIN=blytz.cloud \
  blytz-booking-frontend

# Restart services
docker start blytz-booking-backend blytz-booking-frontend
```

---

## Verify Deployment

### 1. Check Backend Health

```bash
curl https://api.blytz.cloud/health
```

Expected response:
```json
{"status":"healthy"}
```

### 2. Test Subdomain Lookup

```bash
curl https://api.blytz.cloud/api/v1/business/by-subdomain?slug=detail-pro
```

Expected response:
```json
{
  "id": "...",
  "name": "DetailPro Automotive",
  "slug": "detail-pro",
  ...
}
```

### 3. Test Main Domain (SaaS Landing)

```bash
curl -I https://blytz.cloud
```

Expected: `200 OK` and HTML content

### 4. Test Subdomain (Public Booking)

```bash
curl -I https://detail-pro.blytz.cloud
```

Expected: `200 OK` and HTML content (booking page for that business)

### 5. Test Invalid Subdomain

```bash
curl -I https://invalid-subdomain-test.blytz.cloud
```

Expected: Redirects to `https://blytz.cloud`

---

## Test Concurrency Fix

### Load Test Script

Create a test script `test-booking-concurrency.sh`:

```bash
#!/bin/bash

API_URL="https://api.blytz.cloud"
BUSINESS_ID="your-business-uuid"
SERVICE_ID="your-service-uuid"
SLOT_ID="your-slot-uuid"

echo "Testing concurrent bookings on same slot..."
echo "Expected: Exactly MaxBookings bookings succeed, rest fail with 409"
echo ""

for i in {1..10}; do
  curl -X POST "$API_URL/api/v1/bookings" \
    -H "Content-Type: application/json" \
    -d "{
      \"businessId\": \"$BUSINESS_ID\",
      \"serviceId\": \"$SERVICE_ID\",
      \"slotId\": \"$SLOT_ID\",
      \"customer\": {
        \"name\": \"Test User $i\",
        \"email\": \"test$i@example.com\",
        \"phone\": \"555-000$i\"
      }
    }" \
    -w "\nHTTP Status: %{http_code}\n" \
    -o /dev/null \
    2>&1 | grep -E "(HTTP Status|error)" &

  # Small delay to simulate real users
  sleep 0.1
done

wait

echo ""
echo "Test complete. Check your dashboard for bookings."
```

Run it:
```bash
chmod +x test-booking-concurrency.sh
./test-booking-concurrency.sh
```

**Expected Output:**
- Exactly `MaxBookings` bookings with `HTTP Status: 200`
- Rest with `HTTP Status: 409` (Conflict - slot full)

---

## Troubleshooting

### Issue: Subdomain routing not working

**Check:**
```bash
# Check if BASE_DOMAIN is set
docker exec blytz-booking-backend env | grep BASE_DOMAIN

# Check if VITE_BASE_DOMAIN is set
docker exec blytz-booking-frontend env | grep VITE_BASE_DOMAIN
```

**Fix:** Add the missing environment variables (see "Deployment Steps" above)

### Issue: "Business not found" on valid subdomain

**Check:**
```bash
# Test API directly
curl "https://api.blytz.cloud/api/v1/business/by-subdomain?slug=detail-pro"
```

**Fix:**
1. Verify business slug exists in database
2. Check `BASE_DOMAIN` matches actual domain
3. Check Cloudflare wildcard DNS: `*.blytz.cloud` → VPS IP

### Issue: Operator routes accessible on subdomain

**Expected behavior:** Should redirect to main domain

**Check:**
```bash
# Test redirect
curl -I https://detail-pro.blytz.cloud/dashboard
```

**Expected:** `302 Found` with `Location: https://blytz.cloud/dashboard`

### Issue: Concurrent bookings still failing

**Check:**
```sql
-- Connect to PostgreSQL
docker exec -it blytz-booking-postgres psql -U postgres -d blytz

-- Check if row locking works
SELECT slot_id, booking_count, is_booked
FROM slots
WHERE id = 'your-slot-id';

-- Check transaction isolation level
SHOW default_transaction_isolation_level;
```

**Expected:** `read committed` or `serializable`

**Fix:** Ensure PostgreSQL version >= 9.5 (supports row locking)

---

## Rollback Plan

If something breaks, rollback by:

### 1. Remove Subdomain Middleware

**Backend:**
```bash
# Edit cmd/server/main.go
# Comment out these lines:
# subdomainMiddleware := middleware.NewSubdomainMiddleware(...)
# r.Use(subdomainMiddleware.ExtractAndValidate())
# v1.GET("/business/by-subdomain", ...)
```

**Frontend:**
```bash
# Revert to path-based routing
# Edit routes/router.tsx
# Add back: { path: 'business/:slug', element: <PublicBooking /> }
```

### 2. Remove Environment Variables

```bash
# Backend
docker update --env-rm BASE_DOMAIN blytz-booking-backend
docker update --env-rm DB_MAX_OPEN_CONNS blytz-booking-backend
docker update --env-rm DB_MAX_IDLE_CONNS blytz-booking-backend

# Frontend
docker update --env-rm VITE_BASE_DOMAIN blytz-booking-frontend
```

### 3. Redeploy

```bash
docker restart blytz-booking-backend blytz-booking-frontend
```

---

## Success Criteria

You know it's working when:

- [ ] `https://blytz.cloud` → Shows SaaS landing page
- [ ] `https://detail-pro.blytz.cloud` → Shows booking page for "DetailPro"
- [ ] `https://lumina-spa.blytz.cloud` → Shows booking page for "Lumina Spa"
- [ ] `https://invalid.blytz.cloud` → Redirects to `https://blytz.cloud`
- [ ] `https://blytz.cloud/dashboard` → Operator dashboard works
- [ ] `https://detail-pro.blytz.cloud/dashboard` → Redirects to `blytz.cloud/dashboard`
- [ ] Concurrent bookings on same slot → Only `MaxBookings` succeed
- [ ] API health check → Returns `{"status":"healthy"}`

---

## What's Different Now?

### Before (Path-Based Routing)

```
https://blytz.cloud/business/detail-pro  ❌
https://blytz.cloud/business/lumina-spa  ❌
```

### After (Subdomain Routing)

```
https://detail-pro.blytz.cloud/  ✅
https://lumina-spa.blytz.cloud/  ✅
```

### Before (Race Conditions)

```
User A books slot X at T=0.001
User B books slot X at T=0.002
→ Both succeed even if MaxBookings=1
→ OVERBOOKING!
```

### After (Row-Level Locking)

```
User A books slot X at T=0.001 → Locks slot, creates booking, unlocks
User B books slot X at T=0.002 → Waits for lock, checks capacity, fails
→ NO OVERBOOKING!
```

---

## Next Steps

1. **Add environment variables** (BASE_DOMAIN, VITE_BASE_DOMAIN, DB_MAX_OPEN_CONNS, DB_MAX_IDLE_CONNS)
2. **Redeploy** backend and frontend
3. **Test** subdomain routing with actual businesses
4. **Run** concurrency test script
5. **Monitor** logs for errors

---

## Support

If you encounter issues:

1. Check logs:
   ```bash
   docker logs blytz-booking-backend
   docker logs blytz-booking-frontend
   ```

2. Check environment:
   ```bash
   docker exec blytz-booking-backend env | sort
   docker exec blytz-booking-frontend env | sort
   ```

3. Test API directly:
   ```bash
   curl -v https://api.blytz.cloud/health
   ```

---

**🎉 Ready to deploy! Add the environment variables and redeploy!**
