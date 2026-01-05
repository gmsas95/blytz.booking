# 🛠️ CRITICAL FIX: Business Ownership Security

## Problem Identified

**Current Issue:** Any logged-in user can see/edit/delete ALL businesses, not just their own.

This is a **severe security vulnerability** for multi-tenant SaaS.

---

## ✅ SOLUTION: Implement Business Ownership (One-to-One)

### 1. Database Model Updates

**File:** `backend/internal/models/models.go`

**Changes:**
```go
// Add OwnerID to Business struct
type Business struct {
    ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    OwnerID         uuid.UUID `json:"-" gorm:"type:uuid;not null;index;unique"`  // NEW
    Name            string    `json:"name" gorm:"not null"`
    Slug            string    `json:"slug" gorm:"uniqueIndex;not null"`
    // ... other fields
    Owner           User      `json:"-" gorm:"foreignKey:OwnerID"`  // NEW
}

// Add Employee model (for future staff access)
type Employee struct {
    ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
    UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    Email      string    `json:"email" gorm:"not null;uniqueIndex"`
    Role       string    `json:"role" gorm:"not null;default:'staff'"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    Business   Business  `json:"-" gorm:"foreignKey:BusinessID"`
    User       User      `json:"-" gorm:"foreignKey:UserID"`
}
```

### 2. Auto-Create Business on User Registration

**File:** `backend/internal/services/auth_service.go`

**Add import:**
```go
import "strings"
```

**Update Register() function:**
```go
func (s *AuthService) Register(email, name, password string) (*models.User, string, error) {
    // Check if user already exists
    var existingUser models.User
    if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
        return nil, "", ErrConflict
    }

    // Hash password
    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return nil, "", err
    }

    // Create user
    user := models.User{
        Email:        email,
        Name:         name,
        PasswordHash: hashedPassword,
    }

    if err := s.DB.Create(&user).Error; err != nil {
        return nil, "", err
    }

    // CREATE DEFAULT BUSINESS for new user
    businessSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
    business := &models.Business{
        OwnerID:        user.ID,  // NEW
        Name:            name + "'s Business",
        Slug:            businessSlug,
        Vertical:        "General",
        Description:     "Your business for scheduling and bookings",
        ThemeColor:      "blue",
        SlotDurationMin: 30,
        MaxBookings:     1,
    }

    if err := s.DB.Create(business).Error; err != nil {
        // Log error but don't fail registration
        // Business creation is optional
    }

    // Generate token
    token, err := auth.GenerateToken(user.ID.String(), user.Email)
    if err != nil {
        return nil, "", err
    }

    return &user, token, nil
}
```

### 3. Add Ownership Checks to Business Handlers

**File:** `backend/internal/handlers/handlers.go`

**Update ListBusinesses:**
```go
func (h *Handler) ListBusinesses(c *gin.Context) {
    userID := c.GetString("user_id")

    var businesses []models.Business
    if err := h.Repo.DB.Where("owner_id = ?", userID).Find(&businesses).Error; err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch businesses"})
        return
    }

    response := make([]dto.BusinessResponse, len(businesses))
    for i, b := range businesses {
        response[i] = dto.BusinessResponse{
            ID:              b.ID.String(),
            Name:            b.Name,
            Slug:            b.Slug,
            Vertical:        b.Vertical,
            Description:     b.Description,
            ThemeColor:      b.ThemeColor,
            SlotDurationMin: b.SlotDurationMin,
            MaxBookings:     b.MaxBookings,
            CreatedAt:       b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt:       b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
        }
    }

    c.JSON(http.StatusOK, response)
}
```

**Update UpdateBusiness:**
```go
func (h *Handler) UpdateBusiness(c *gin.Context) {
    id := c.Param("businessId")
    businessID, err := uuid.Parse(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    userID := c.GetString("user_id")

    var existingBusiness models.Business
    if err := h.Repo.DB.Where("id = ? AND owner_id = ?", businessID, userID).First(&existingBusiness).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }

    var req dto.UpdateBusinessRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
        return
    }

    updates := &models.Business{}
    if req.Name != nil {
        updates.Name = *req.Name
    }
    // ... other updates (keep existing logic)

    if err := h.BusinessService.Update(businessID, updates); err != nil {
        if err == services.ErrNotFound {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update business"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Business updated successfully"})
}
```

**Update GetServicesByBusiness:**
```go
func (h *Handler) GetServicesByBusiness(c *gin.Context) {
    businessIDStr := c.Param("businessId")
    businessUUID, err := uuid.Parse(businessIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    userID := c.GetString("user_id")

    var existingBusiness models.Business
    if err := h.Repo.DB.Where("id = ? AND owner_id = ?", businessUUID, userID).First(&existingBusiness).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }

    services, err := h.ServiceService.GetByBusiness(businessUUID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch services"})
        return
    }

    response := make([]dto.ServiceResponse, len(services))
    for i, s := range services {
        response[i] = dto.ServiceResponse{
            ID:            s.ID.String(),
            BusinessID:    s.BusinessID.String(),
            Name:          s.Name,
            Description:   s.Description,
            DurationMin:   s.DurationMin,
            TotalPrice:    s.TotalPrice,
            DepositAmount: s.DepositAmount,
            CreatedAt:     s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt:     s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
        }
    }

    c.JSON(http.StatusOK, response)
}
```

**Update CreateService:** (Add ownership check)
```go
func (h *Handler) CreateService(c *gin.Context) {
    businessIDStr := c.Param("businessId")
    businessID, err := uuid.Parse(businessIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    userID := c.GetString("user_id")

    var existingBusiness models.Business
    if err := h.Repo.DB.Where("id = ? AND owner_id = ?", businessID, userID).First(&existingBusiness).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }

    var req dto.CreateServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
        return
    }

    service := &models.Service{
        ID:            uuid.New(),
        BusinessID:    businessID,
        Name:          req.Name,
        Description:   req.Description,
        DurationMin:   req.DurationMin,
        TotalPrice:    req.TotalPrice,
        DepositAmount: req.DepositAmount,
    }

    if err := h.ServiceService.Create(service); err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create service"})
        return
    }

    c.JSON(http.StatusCreated, dto.ServiceResponse{
        ID:            service.ID.String(),
        BusinessID:    service.BusinessID.String(),
        Name:          service.Name,
        Description:   service.Description,
        DurationMin:   service.DurationMin,
        TotalPrice:    service.TotalPrice,
        DepositAmount: service.DepositAmount,
        CreatedAt:     service.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:     service.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    })
}
```

**Update UpdateService:** (Add ownership check)
```go
func (h *Handler) UpdateService(c *gin.Context) {
    businessIDStr := c.Param("businessId")
    serviceID := c.Param("serviceId")

    businessUUID, err := uuid.Parse(businessIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    serviceUUID, err := uuid.Parse(serviceID)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid service ID"})
        return
    }

    userID := c.GetString("user_id")

    var business models.Business
    if err := h.Repo.DB.Where("id = ?", businessUUID).First(&business).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }

    if business.OwnerID != userID {
        c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "You don't have permission to update this service"})
        return
    }

    var req dto.UpdateServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
        return
    }

    service, err := h.ServiceService.GetByID(serviceUUID)
    if err != nil {
        if err == services.ErrNotFound {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Service not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch service"})
        return
    }

    updates := &models.Service{}
    if req.Name != nil {
        updates.Name = *req.Name
    }
    if req.Description != nil {
        updates.Description = *req.Description
    }
    if req.DurationMin != nil {
        updates.DurationMin = *req.DurationMin
    }
    if req.TotalPrice != nil {
        updates.TotalPrice = *req.TotalPrice
    }
    if req.DepositAmount != nil {
        updates.DepositAmount = *req.DepositAmount
    }

    if err := h.ServiceService.Update(serviceUUID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update service"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}
```

**Update DeleteService:** (Add ownership check)
```go
func (h *Handler) DeleteService(c *gin.Context) {
    businessIDStr := c.Param("businessId")
    serviceID := c.Param("serviceId")

    businessUUID, err := uuid.Parse(businessIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    serviceUUID, err := uuid.Parse(serviceID)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid service ID"})
        return
    }

    userID := c.GetString("user_id")

    var business models.Business
    if err := h.Repo.DB.Where("id = ?", businessUUID).First(&business).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }

    if business.OwnerID != userID {
        c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "You don't have permission to delete this service"})
        return
    }

    service, err := h.ServiceService.GetByID(serviceUUID)
    if err != nil {
        if err == services.ErrNotFound {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Service not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch service"})
        return
    }

    if service.BusinessID != businessUUID {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Service does not belong to this business"})
        return
    }

    if err := h.ServiceService.Delete(serviceUUID); err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete service"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
```

### 4. Update Repository Migration

**File:** `backend/internal/repository/repository.go`

**Update AutoMigrate:**
```go
func (r *Repository) AutoMigrate() error {
    return r.DB.AutoMigrate(
        &models.Business{},
        &models.BusinessAvailability{},
        &models.Service{},
        &models.Slot{},
        &models.Booking{},
        &models.User{},
        &models.Employee{},  // NEW: For future staff access
    )
}
```

---

### 5. Frontend Changes

**File:** `screens/OperatorDashboard.tsx`

**Remove business dropdown** - Users only have 1 business, so no dropdown needed:

```typescript
// Remove or comment out business selector UI
// Show single business name only

{currentBusiness && (
  <div className="business-header">
    <h1>{currentBusiness.name}</h1>
  </div>
)}

// If no business, show setup message
{!currentBusiness && (
  <div className="no-business">
    <h2>No Business Found</h2>
    <p>Creating your business now...</p>
  </div>
)}
```

**Remove handleBusinessChange function** - No longer needed

**Update handleSaveBusiness** - Already fixed (previous commits)

---

## 🎯 TRUE SUCCESS METRICS After Implementation

### Security ✅
- Users can ONLY access their own businesses
- Cross-business data access prevented
- API endpoints verify ownership on ALL business operations

### User Experience ✅
- Simple: One business per user (no dropdown confusion)
- Clean: No accidental editing of wrong business
- Transparent: Users only see what they own

### Data Isolation ✅
- Businesses isolated by owner_id
- Services isolated per business
- Slots isolated per business
- Bookings isolated per business

---

## 📋 Implementation Checklist

- [ ] **Database Models:** Add OwnerID to Business, create Employee model
- [ ] **Repository:** Update AutoMigrate to include Employee table
- [ ] **Auth Service:** Update Register() to auto-create business
- [ ] **Handlers:** Update ListBusinesses (filter by owner_id)
- [ ] **Handlers:** Update UpdateBusiness (check ownership)
- [ ] **Handlers:** Update GetServicesByBusiness (check ownership)
- [ ] **Handlers:** Update CreateService (check ownership)
- [ ] **Handlers:** Update UpdateService (check ownership)
- [ ] **Handlers:** Update DeleteService (check ownership)
- [ ] **Frontend:** Remove business dropdown
- [ ] **Frontend:** Remove handleBusinessChange
- [ ] **Frontend:** Update to handle single-business case
- [ ] **Test:** Verify users can't access other businesses
- [ ] **Test:** Verify business creation on registration

---

## 🚨 IMMEDIATE ACTION REQUIRED

**This is CRITICAL for production.** Do NOT deploy until business ownership is implemented.

**Without these changes, ANY user can:**
- Edit ANY business settings
- Delete ANY business services
- See ANY business bookings
- Access ANY business data

---

## ❓ Need Confirmation

Before I proceed with implementation, confirm:

1. **One-to-One Relationship?**
   - Each user has exactly ONE business
   - Simple ownership model
   - No team members initially

2. **Business Auto-Creation on Registration?**
   - Yes: Create default business when user signs up
   - User sees their business immediately after login

3. **Frontend Simplicity?**
   - Remove business dropdown
   - Show single business name
   - No business switching UI

**This is much simpler than many-to-many and more secure than the current open access model.**
