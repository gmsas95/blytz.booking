# Critical Fixes Deployment Guide

## Issues Fixed

### Issue 1: Database Migration Failure ❌
**Problem:** `owner_id` column added as NOT NULL to existing businesses with NULL values

**Solution:** Made `owner_id` nullable in model
- Removed `not null` constraint from `OwnerID` field
- Migration will now succeed

### Issue 2: Auth Middleware Blocking Everything ❌
**Problem:** `/api/v1/businesses` requires auth, breaking login flow

**Solution:** Made `/api/v1/businesses` publicly accessible with conditional filtering
- Unauthenticated users → receive empty array `[]`
- Authenticated users → receive their businesses only
- Removes chicken-and-egg problem

---

## Deployment Steps

### Step 1: Redeploy Backend in Dokploy (CRITICAL - Do First!)

1. Go to Dokploy → Backend Service
2. Trigger redeploy from `staging` branch
3. Wait for container to start (~1-2 minutes)
4. Check logs: Should see "Starting server on port 8080"
5. Verify migration success: Should NOT see `owner_id` errors

**Expected Logs After Success:**
```
Running migrations...
Migrated models: User, Business, Employee, Service, Slot, Booking, BusinessAvailability
Starting server on port 8080...
```

**No More Errors Like:**
```
ERROR: column "owner_id" of relation "businesses" contains null values
```

---

### Step 2: Run SQL Migration for Existing Businesses

**IMPORTANT:** After backend starts successfully, run this SQL to assign `owner_id` to existing businesses.

**Option A: Via Docker CLI**
```bash
# Find backend container name
docker ps | grep backend

# Run migration script
docker exec -it <backend-container-name> psql -U <db-user> -d <db-name> < migrate_owner_id.sql

# Example:
docker exec -it blytz-booking-backend psql -U postgres -d blytz < migrate_owner_id.sql
```

**Option B: Via Database Client**
1. Connect to PostgreSQL database (via pgAdmin, DBeaver, or CLI)
2. Copy and paste SQL from `migrate_owner_id.sql`
3. Execute script
4. Verify output shows updated business count

**Option C: Via Direct SQL**
```sql
-- Assign all existing businesses to first user
UPDATE businesses
SET owner_id = (SELECT id FROM users LIMIT 1)
WHERE owner_id IS NULL;

-- Verify
SELECT COUNT(*) as total_businesses,
       COUNT(owner_id) as businesses_with_owner
FROM businesses;
```

---

### Step 3: Redeploy Frontend in Dokploy

1. Go to Dokploy → Frontend Service
2. Verify repository: `git@github.com:gmsas95/blytz.booking.git`
3. Verify branch: `staging`
4. Trigger redeploy
5. Wait for build to complete (~2-3 minutes)

**Expected Frontend Build:**
```
vite v6.x.x building for production...
✓ built in X.XXs
```

---

### Step 4: Verify Fixes

**Test 1: Backend Health**
```bash
curl https://api.blytz.cloud/health
# Expected: {"status":"healthy"}
```

**Test 2: Unauthenticated Businesses API**
```bash
curl https://api.blytz.cloud/api/v1/businesses
# Expected: [] (empty array - user not logged in)
```

**Test 3: Login Flow**
1. Visit `https://blytz.cloud/login`
2. Should NOT see CORS errors
3. Should NOT see 404 errors
4. Login form should load correctly

**Test 4: Registration with Auto-Create Business**
```bash
curl -X POST https://api.blytz.cloud/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","password":"test123"}'
# Expected: User created + business auto-created
```

**Test 5: Business Ownership Security**
```bash
# Get token from registration
TOKEN="your-jwt-token-here"

# Get businesses (should return only user's business)
curl -H "Authorization: Bearer $TOKEN" \
  https://api.blytz.cloud/api/v1/businesses
# Expected: Array with 1 business (Test User's business)
```

---

## Verification Checklist

- [ ] Backend redeployed successfully from staging branch
- [ ] No `owner_id` migration errors in backend logs
- [ ] SQL migration script executed (existing businesses have owner_id)
- [ ] Frontend redeployed successfully
- [ ] Login page loads without CORS errors
- [ ] Register page loads without CORS errors
- [ ] New users get auto-created business
- [ ] Logged-in users see only their business
- [ ] Unauthenticated users see empty businesses array

---

## Troubleshooting

**Problem: Backend still showing `owner_id` errors**

**Solution:**
1. Verify backend pulled latest staging code
2. Check `backend/internal/models/models.go` - OwnerID should NOT have `not null`
3. Hard reset backend container in Dokploy
4. Check PostgreSQL: `\d businesses` - owner_id should be nullable

**Problem: Frontend still showing CORS errors**

**Solution:**
1. Verify frontend pulled latest staging code
2. Check backend CORS middleware in `main.go`
3. Hard reset frontend container in Dokploy
4. Clear browser cache and cookies

**Problem: New users not getting auto-created business**

**Solution:**
1. Check backend logs for registration errors
2. Verify `backend/internal/services/auth_service.go` has business creation code
3. Check if slug generation is causing uniqueness errors
4. Verify OwnerID is being set correctly in registration

**Problem: Existing businesses still showing in dashboard**

**Solution:**
1. Run SQL migration script (see Step 2)
2. Verify owner_id is set in database
3. Check that `ListBusinesses` handler is using `GetByUser()`
4. Verify user_id is set correctly in auth context

---

## After Successful Deployment

1. **Test Complete User Flow:**
   - Register new user → should auto-create business
   - Login as user → should see only their business
   - Edit business → should work
   - Create service → should work
   - Create slot → should work

2. **Verify Security:**
   - User A cannot access User B's business → should get 403 Forbidden
   - Unauthenticated users cannot create services/slots → should get 401

3. **Monitor Logs:**
   ```bash
   # Watch backend logs for any errors
   docker logs blytz-booking-backend -f
   ```

---

## Production Deployment (After Staging Success)

1. Merge staging → production
2. Redeploy production backend
3. Run SQL migration on production database
4. Redeploy production frontend
5. Run final verification tests
6. Monitor for errors for 1-2 hours

---

## Notes

- **owner_id is nullable** - This is intentional for backwards compatibility
- **Existing businesses need migration** - SQL script assigns them to first user
- **Auth is optional for GET /businesses** - Returns empty array for unauthenticated
- **All mutations require auth** - Create/update/delete still protected
- **One-to-one model enforced** - New users get auto-created business

---

## Support

If issues persist after following this guide:
1. Check backend logs: `docker logs blytz-booking-backend -f`
2. Check frontend logs: `docker logs blytz-booking-frontend -f`
3. Verify database state: Connect via database client and check `SELECT * FROM businesses;`
4. Run verification script: `./verify-deployment.sh`
