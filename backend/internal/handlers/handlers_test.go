package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"blytz.cloud/backend/internal/auth"
	"blytz.cloud/backend/internal/middleware"
	"blytz.cloud/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testOrigin = "http://localhost:3000"

func setupHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	statements := []string{
		`CREATE TABLE users (id text PRIMARY KEY, email text NOT NULL, name text, password_hash text NOT NULL, token_version integer NOT NULL DEFAULT 1, created_at datetime, updated_at datetime)`,
		`CREATE TABLE businesses (id text PRIMARY KEY, name text NOT NULL, slug text NOT NULL, vertical text NOT NULL, description text, theme_color text, created_at datetime, updated_at datetime)`,
		`CREATE TABLE memberships (id text PRIMARY KEY, user_id text NOT NULL, business_id text NOT NULL, role text NOT NULL, created_at datetime, updated_at datetime)`,
		`CREATE TABLE bookings (id text PRIMARY KEY, business_id text NOT NULL, service_id text NOT NULL, slot_id text NOT NULL, service_name text NOT NULL, slot_time datetime NOT NULL, name text NOT NULL, email text NOT NULL, phone text NOT NULL, status text NOT NULL, deposit_paid_minor integer NOT NULL, total_price_minor integer NOT NULL, currency_code text NOT NULL, created_at datetime, updated_at datetime)`,
		`CREATE TABLE customers (id text PRIMARY KEY, business_id text NOT NULL, name text NOT NULL, email text NOT NULL, phone text NOT NULL, notes text, created_at datetime, updated_at datetime)`,
		`CREATE TABLE vehicles (id text PRIMARY KEY, business_id text NOT NULL, customer_id text NOT NULL, year integer, make text NOT NULL, model text NOT NULL, color text, license_plate text, created_at datetime, updated_at datetime)`,
		`CREATE TABLE jobs (id text PRIMARY KEY, business_id text NOT NULL, customer_id text NOT NULL, vehicle_id text NOT NULL, booking_id text, title text NOT NULL, status text NOT NULL, scheduled_at datetime NOT NULL, notes text, created_at datetime, updated_at datetime)`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("create test schema: %v", err)
		}
	}

	return db
}

func seedHandlerTestData(t *testing.T, db *gorm.DB) (string, string, string) {
	t.Helper()

	userID := uuid.New().String()
	businessID := uuid.New().String()
	otherBusinessID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	queries := []string{
		fmt.Sprintf(`INSERT INTO users (id, email, name, password_hash, token_version, created_at, updated_at) VALUES ('%s', 'owner@example.com', 'Owner', '%s', 1, '%s', '%s')`, userID, hashedPassword, now, now),
		fmt.Sprintf(`INSERT INTO businesses (id, name, slug, vertical, description, theme_color, created_at, updated_at) VALUES ('%s', 'DetailPro Automotive', 'detail-pro', 'Automotive', 'Premium detailing workshop', 'blue', '%s', '%s')`, businessID, now, now),
		fmt.Sprintf(`INSERT INTO businesses (id, name, slug, vertical, description, theme_color, created_at, updated_at) VALUES ('%s', 'Other Workshop', 'other-workshop', 'Automotive', 'Second workshop', 'zinc', '%s', '%s')`, otherBusinessID, now, now),
		fmt.Sprintf(`INSERT INTO memberships (id, user_id, business_id, role, created_at, updated_at) VALUES ('%s', '%s', '%s', 'OWNER', '%s', '%s')`, uuid.New().String(), userID, businessID, now, now),
		fmt.Sprintf(`INSERT INTO bookings (id, business_id, service_id, slot_id, service_name, slot_time, name, email, phone, status, deposit_paid_minor, total_price_minor, currency_code, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', 'Full Interior Detail', '%s', 'Alice Smith', 'alice@example.com', '555-0101', 'CONFIRMED', 5000, 20000, 'USD', '%s', '%s')`, uuid.New().String(), businessID, uuid.New().String(), uuid.New().String(), now, now, now),
		fmt.Sprintf(`INSERT INTO customers (id, business_id, name, email, phone, notes, created_at, updated_at) VALUES ('%s', '%s', 'Alice Smith', 'alice@example.com', '555-0101', 'VIP detail client', '%s', '%s')`, uuid.New().String(), businessID, now, now),
	}

	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			t.Fatalf("seed test data: %v", err)
		}
	}

	return userID, businessID, otherBusinessID
}

func setupHandlerRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	auth.SetJWTSecret("test-secret")
	auth.SetCookieName("blytz_session")
	SetForceSecureCookies(false)
	repo := &repository.Repository{DB: db}
	handler := NewHandler(repo)
	router := gin.New()
	v1 := router.Group("/api/v1")
	authRoutes := v1.Group("/auth")
	authRoutes.Use(middleware.RequireAllowedOrigin([]string{testOrigin}), middleware.RateLimitByIP(30, time.Minute), middleware.RateLimitByIPAndEmail(10, time.Minute))
	authRoutes.POST("/login", handler.Login)
	authRoutes.POST("/register", handler.Register)
	v1.GET("/auth/me", auth.AuthMiddleware(handler.AuthService), handler.GetCurrentUser)
	v1.POST("/auth/logout", middleware.RequireAllowedOrigin([]string{testOrigin}), auth.AuthMiddleware(handler.AuthService), handler.Logout)
	operator := v1.Group("/businesses/:businessId")
	operator.Use(auth.AuthMiddleware(handler.AuthService), middleware.RequireBusinessMembership(handler.AuthService))
	operator.GET("/bookings", handler.ListBookings)
	operator.GET("/customers", handler.ListCustomers)
	operator.POST("/vehicles", middleware.RequireAllowedOrigin([]string{testOrigin}), handler.CreateVehicle)
	return router
}

func authHeaderForTest(t *testing.T, userID string) string {
	t.Helper()
	auth.SetJWTSecret("test-secret")
	token, err := auth.GenerateToken(userID, "owner@example.com", 1)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return "Bearer " + token
}

func TestGetCurrentUserIncludesMembershipContext(t *testing.T) {
	db := setupHandlerTestDB(t)
	userID, businessID, _ := seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", authHeaderForTest(t, userID))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Memberships []struct {
			BusinessID string `json:"business_id"`
			Role       string `json:"role"`
		} `json:"memberships"`
		ActiveBusinessID string `json:"active_business_id"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.User.ID != userID {
		t.Fatalf("expected user id %s, got %s", userID, payload.User.ID)
	}
	if len(payload.Memberships) != 1 {
		t.Fatalf("expected 1 membership, got %d", len(payload.Memberships))
	}
	if payload.Memberships[0].BusinessID != businessID {
		t.Fatalf("expected business id %s, got %s", businessID, payload.Memberships[0].BusinessID)
	}
	if payload.ActiveBusinessID != businessID {
		t.Fatalf("expected active business id %s, got %s", businessID, payload.ActiveBusinessID)
	}
}

func TestListBookingsRequiresMatchingMembership(t *testing.T) {
	db := setupHandlerTestDB(t)
	userID, businessID, otherBusinessID := seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	allowedRequest := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+businessID+"/bookings", nil)
	allowedRequest.Header.Set("Authorization", authHeaderForTest(t, userID))
	allowedRecorder := httptest.NewRecorder()
	router.ServeHTTP(allowedRecorder, allowedRequest)

	if allowedRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for member access, got %d", allowedRecorder.Code)
	}

	forbiddenRequest := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+otherBusinessID+"/bookings", nil)
	forbiddenRequest.Header.Set("Authorization", authHeaderForTest(t, userID))
	forbiddenRecorder := httptest.NewRecorder()
	router.ServeHTTP(forbiddenRecorder, forbiddenRequest)

	if forbiddenRecorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member access, got %d", forbiddenRecorder.Code)
	}
}

func TestListCustomersRequiresMatchingMembership(t *testing.T) {
	db := setupHandlerTestDB(t)
	userID, businessID, otherBusinessID := seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	allowedRequest := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+businessID+"/customers", nil)
	allowedRequest.Header.Set("Authorization", authHeaderForTest(t, userID))
	allowedRecorder := httptest.NewRecorder()
	router.ServeHTTP(allowedRecorder, allowedRequest)

	if allowedRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for member customer access, got %d", allowedRecorder.Code)
	}

	forbiddenRequest := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+otherBusinessID+"/customers", nil)
	forbiddenRequest.Header.Set("Authorization", authHeaderForTest(t, userID))
	forbiddenRecorder := httptest.NewRecorder()
	router.ServeHTTP(forbiddenRecorder, forbiddenRequest)

	if forbiddenRecorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member customer access, got %d", forbiddenRecorder.Code)
	}
}

func TestCreateVehicleRejectsForeignCustomerReference(t *testing.T) {
	db := setupHandlerTestDB(t)
	userID, businessID, _ := seedHandlerTestData(t, db)
	foreignCustomerID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	if err := db.Exec(fmt.Sprintf(`INSERT INTO customers (id, business_id, name, email, phone, notes, created_at, updated_at) VALUES ('%s', '%s', 'Other Customer', 'other@example.com', '555-0202', '', '%s', '%s')`, foreignCustomerID, uuid.New().String(), now, now)).Error; err != nil {
		t.Fatalf("seed foreign customer: %v", err)
	}
	router := setupHandlerRouter(db)

	body := `{"customer_id":"` + foreignCustomerID + `","year":2022,"make":"Tesla","model":"Model Y","color":"White","license_plate":"TEST123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+businessID+"/vehicles", strings.NewReader(body))
	req.Header.Set("Authorization", authHeaderForTest(t, userID))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", testOrigin)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for foreign customer reference, got %d", recorder.Code)
	}
}

func TestLoginIsRateLimitedByIP(t *testing.T) {
	db := setupHandlerTestDB(t)
	_, _, _ = seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	body := `{"email":"owner@example.com","password":"wrong-password"}`
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", testOrigin)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 before limit, got %d on attempt %d", recorder.Code, i+1)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", testOrigin)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after rate limit, got %d", recorder.Code)
	}
}

func TestLoginSetsSessionCookieAndAuthMeAcceptsIt(t *testing.T) {
	db := setupHandlerTestDB(t)
	userID, businessID, _ := seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	loginBody := `{"email":"owner@example.com","password":"password123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Origin", testOrigin)
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginReq)

	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 login, got %d", loginRecorder.Code)
	}
	cookies := loginRecorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected auth cookie to be set")
	}

	authMeReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	authMeReq.AddCookie(cookies[0])
	authMeRecorder := httptest.NewRecorder()
	router.ServeHTTP(authMeRecorder, authMeReq)

	if authMeRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 auth/me with cookie, got %d", authMeRecorder.Code)
	}

	var payload struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		ActiveBusinessID string `json:"active_business_id"`
	}
	if err := json.Unmarshal(authMeRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode auth/me: %v", err)
	}
	if payload.User.ID != userID {
		t.Fatalf("expected user id %s, got %s", userID, payload.User.ID)
	}
	if payload.ActiveBusinessID != businessID {
		t.Fatalf("expected business id %s, got %s", businessID, payload.ActiveBusinessID)
	}
}

func TestLogoutRevokesCurrentSession(t *testing.T) {
	db := setupHandlerTestDB(t)
	_, _, _ = seedHandlerTestData(t, db)
	router := setupHandlerRouter(db)

	loginBody := `{"email":"owner@example.com","password":"password123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Origin", testOrigin)
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginReq)
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 login, got %d", loginRecorder.Code)
	}
	cookies := loginRecorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected auth cookie to be set")
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutReq.AddCookie(cookies[0])
	logoutReq.Header.Set("Origin", testOrigin)
	logoutRecorder := httptest.NewRecorder()
	router.ServeHTTP(logoutRecorder, logoutReq)
	if logoutRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 logout, got %d", logoutRecorder.Code)
	}

	authMeReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	authMeReq.AddCookie(cookies[0])
	authMeRecorder := httptest.NewRecorder()
	router.ServeHTTP(authMeRecorder, authMeReq)
	if authMeRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout revocation, got %d", authMeRecorder.Code)
	}
}
