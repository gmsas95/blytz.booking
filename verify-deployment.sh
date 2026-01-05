#!/bin/bash

# Post-Deployment Verification Script for blytz.cloud
# Run this after adding environment variables and redeploying

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="https://api.blytz.cloud"
BASE_DOMAIN="blytz.cloud"

echo -e "${YELLOW}=== blytz.cloud Deployment Verification ===${NC}"
echo ""

# Test 1: Backend Health Check
echo -e "${YELLOW}Test 1: Backend Health Check${NC}"
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/health" || echo "000")
if [ "$HEALTH" = "200" ]; then
    echo -e "${GREEN}✓${NC} Backend is healthy"
else
    echo -e "${RED}✗${NC} Backend health check failed (HTTP $HEALTH)"
    echo "  Check: docker logs blytz-booking-backend"
fi
echo ""

# Test 2: Main Domain (SaaS Landing)
echo -e "${YELLOW}Test 2: Main Domain (blytz.cloud)${NC}"
MAIN_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "https://$BASE_DOMAIN" || echo "000")
if [ "$MAIN_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✓${NC} Main domain loads (SaaS Landing)"
else
    echo -e "${RED}✗${NC} Main domain failed (HTTP $MAIN_RESPONSE)"
    echo "  Check: docker logs blytz-booking-frontend"
    echo "  Check: VITE_BASE_DOMAIN environment variable"
fi
echo ""

# Test 3: Valid Subdomain (DetailPro)
echo -e "${YELLOW}Test 3: Valid Subdomain (detail-pro.blytz.cloud)${NC}"
SUBDOMAIN_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "https://detail-pro.$BASE_DOMAIN" || echo "000")
if [ "$SUBDOMAIN_RESPONSE" = "200" ]; then
    echo -e "${GREEN}✓${NC} Valid subdomain loads (Public Booking)"
else
    echo -e "${RED}✗${NC} Valid subdomain failed (HTTP $SUBDOMAIN_RESPONSE)"
    echo "  Check: Cloudflare wildcard DNS (*.blytz.cloud → VPS IP)"
    echo "  Check: BASE_DOMAIN environment variable in backend"
fi
echo ""

# Test 4: API - Get All Businesses
echo -e "${YELLOW}Test 4: API - List All Businesses${NC}"
BUSINESSES=$(curl -s "$API_URL/api/v1/businesses" || echo "ERROR")
if echo "$BUSINESSES" | grep -q "detail-pro\|lumina-spa\|flash-frame"; then
    echo -e "${GREEN}✓${NC} Businesses API returns data"
    BUSINESS_COUNT=$(echo "$BUSINESSES" | grep -o '"id"' | wc -l)
    echo "  Found $BUSINESS_COUNT business(es)"
else
    echo -e "${RED}✗${NC} Businesses API failed or no data"
    echo "  Response: $BUSINESSES"
fi
echo ""

# Test 5: API - Get Business by Subdomain
echo -e "${YELLOW}Test 5: API - Get Business by Subdomain (detail-pro)${NC}"
BUSINESS_BY_SLUG=$(curl -s "$API_URL/api/v1/business/by-subdomain?slug=detail-pro" || echo "ERROR")
if echo "$BUSINESS_BY_SLUG" | grep -q "DetailPro\|detail-pro"; then
    echo -e "${GREEN}✓${NC} Business by subdomain API works"
else
    echo -e "${RED}✗${NC} Business by subdomain API failed"
    echo "  Response: $BUSINESS_BY_SLUG"
    echo "  Check: BASE_DOMAIN environment variable in backend"
fi
echo ""

# Test 6: Invalid Subdomain (Should 404 or redirect)
echo -e "${YELLOW}Test 6: Invalid Subdomain (invalid-subdomain-test.blytz.cloud)${NC}"
INVALID_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}\n%{redirect_url}" "https://invalid-subdomain-test.$BASE_DOMAIN" || echo "000")
if echo "$INVALID_RESPONSE" | grep -q "404\|302"; then
    echo -e "${GREEN}✓${NC} Invalid subdomain handled correctly (404 or redirect)"
else
    echo -e "${RED}✗${NC} Invalid subdomain not handled (HTTP ${INVALID_RESPONSE%%$'\n'*})"
fi
echo ""

# Test 7: Operator Route on Main Domain
echo -e "${YELLOW}Test 7: Operator Route on Main Domain (/dashboard)${NC}"
DASHBOARD_MAIN=$(curl -s -o /dev/null -w "%{http_code}" "https://$BASE_DOMAIN/dashboard" || echo "000")
if [ "$DASHBOARD_MAIN" = "200" ] || [ "$DASHBOARD_MAIN" = "401" ]; then
    echo -e "${GREEN}✓${NC} Dashboard accessible on main domain (200 or 401 if not logged in)"
else
    echo -e "${RED}✗${NC} Dashboard failed (HTTP $DASHBOARD_MAIN)"
fi
echo ""

# Test 8: Operator Route Redirect from Subdomain
echo -e "${YELLOW}Test 8: Operator Route Redirect from Subdomain (/dashboard on subdomain)${NC}"
DASHBOARD_REDIRECT=$(curl -s -I "https://detail-pro.$BASE_DOMAIN/dashboard" | grep -i "Location" | head -1)
if echo "$DASHBOARD_REDIRECT" | grep -q "$BASE_DOMAIN/dashboard"; then
    echo -e "${GREEN}✓${NC} Dashboard redirects to main domain from subdomain"
    echo "  Redirect: $DASHBOARD_REDIRECT"
else
    echo -e "${RED}✗${NC} Dashboard doesn't redirect from subdomain"
    echo "  Expected: redirect to https://$BASE_DOMAIN/dashboard"
    echo "  Got: $DASHBOARD_REDIRECT"
fi
echo ""

# Test 9: Check Environment Variables (Backend)
echo -e "${YELLOW}Test 9: Check Backend Environment Variables${NC}"
echo "Checking for required environment variables..."

BACKEND_CONTAINER=$(docker ps --filter "name=backend" --format "{{.Names}}" | head -1)
if [ -n "$BACKEND_CONTAINER" ]; then
    BASE_DOMAIN_SET=$(docker exec "$BACKEND_CONTAINER" env | grep "^BASE_DOMAIN=" || echo "")
    DB_POOL_SET=$(docker exec "$BACKEND_CONTAINER" env | grep "^DB_MAX_OPEN_CONNS=" || echo "")

    if [ -n "$BASE_DOMAIN_SET" ]; then
        echo -e "${GREEN}✓${NC} BASE_DOMAIN is set: $BASE_DOMAIN_SET"
    else
        echo -e "${RED}✗${NC} BASE_DOMAIN is NOT set"
    fi

    if [ -n "$DB_POOL_SET" ]; then
        echo -e "${GREEN}✓${NC} DB_MAX_OPEN_CONNS is set: $DB_POOL_SET"
    else
        echo -e "${RED}✗${NC} DB_MAX_OPEN_CONNS is NOT set"
    fi
else
    echo -e "${RED}✗${NC} Backend container not found"
fi
echo ""

# Test 10: Check Environment Variables (Frontend)
echo -e "${YELLOW}Test 10: Check Frontend Environment Variables${NC}"
echo "Checking for required environment variables..."

FRONTEND_CONTAINER=$(docker ps --filter "name=frontend\|blytz-cloud" --format "{{.Names}}" | head -1)
if [ -n "$FRONTEND_CONTAINER" ]; then
    VITE_BASE_SET=$(docker exec "$FRONTEND_CONTAINER" env | grep "^VITE_BASE_DOMAIN=" || echo "")

    if [ -n "$VITE_BASE_SET" ]; then
        echo -e "${GREEN}✓${NC} VITE_BASE_DOMAIN is set: $VITE_BASE_SET"
    else
        echo -e "${RED}✗${NC} VITE_BASE_DOMAIN is NOT set"
    fi
else
    echo -e "${RED}✗${NC} Frontend container not found"
fi
echo ""

# Summary
echo -e "${YELLOW}=== Verification Summary ===${NC}"
echo ""
echo "If all tests passed (✓), your deployment is successful!"
echo ""
echo "Next steps:"
echo "1. Test actual booking flow on https://detail-pro.blytz.cloud"
echo "2. Run concurrency test: ./test-booking-concurrency.sh"
echo "3. Check logs: docker logs blytz-booking-backend -f"
echo "4. Monitor metrics in Dokploy dashboard"
echo ""
