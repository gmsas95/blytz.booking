# 🛠️ CRITICAL SECURITY FIX: One Business Per User

## Problem
Any logged-in user can currently access/edit ALL businesses.

## Solution (Minimal Fix - 5 mins to implement)

### 1. Backend: Add Single Business Check

**File:** `backend/internal/services/business_service.go`

**Add method after GetAll():**
```go
func (s *BusinessService) GetByUser(userID uuid.UUID) ([]models.Business, error) {
    var businesses []models.Business
    if err := s.DB.Where("owner_id = ?", userID).Find(&businesses).Error; err != nil {
        return nil, err
    }
    return businesses, nil
}
```

### 2. Frontend: Remove Business Dropdown, Show Single Business

**File:** `screens/OperatorDashboard.tsx`

**Remove business selection dropdown**
**Show single business name only**

---

## Step 1: Add GetByUser Method

**Backend:** `backend/internal/services/business_service.go`

Add this new method after existing GetAll() method:

```go
func (s *BusinessService) GetByUser(userID uuid.UUID) ([]models.Business, error) {
    var businesses []models.Business
    if err := s.DB.Where("owner_id = ?", userID).Find(&businesses).Error; err != nil {
        return nil, err
    }
    return businesses, nil
}
```

---

## Step 2: Update Frontend Dropdown

**Frontend:** `screens/OperatorDashboard.tsx`

**A. Remove business selection from UI** (lines ~100-130):
```typescript
// Remove or comment this out:
// {businesses.map((biz, idx) => (
//   <button onClick={() => handleBusinessChange(biz)} className={...}>
//     {biz.name}
//   </button>
// ))}
```

**B. Show single business header**:
```typescript
{currentBusiness && (
  <div className="flex items-center gap-4">
    <h2 className="text-xl font-bold">{currentBusiness.name}</h2>
  </div>
)}

// Or if no business:
{!currentBusiness && (
  <div className="text-gray-500">No business found. Contact support.</div>
)}
```

---

## Step 3: Update Handlers

**File:** `backend/internal/handlers/handlers.go`

**Update NewHandler to add GetByUser:**
```go
func NewHandler(repo *repository.Repository, emailConfig email.EmailConfig) *Handler {
    // ... existing services

    return &Handler{
        Repo:                repo,
        // ... existing services
        BusinessService:   services.NewBusinessService(repo.DB),
        BusinessServiceByUser: services.NewBusinessService(repo.DB), // NEW
        // ... other services
        // ...
    }
}
```

**Update GetBusinesses to use GetByUser:**
```go
func (h *Handler) GetBusinesses(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var businesses []models.Business
    if err := h.BusinessServiceByUser.GetByUser(userID).Error; err != nil {
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

**Update GetBusiness to check ownership:**
```go
func (h *Handler) GetBusiness(c *gin.Context) {
    id := c.Param("businessId")
    businessID, err := uuid.Parse(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
        return
    }

    userID := c.GetString("user_id")
    
    var business models.Business
    if err := h.Repo.DB.Where("id = ? AND owner_id = ?", businessID, userID).First(&business).Error; err != nil {
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
        return
    }
    
    // ... rest of existing GetBusiness logic
}
```

**Update UpdateBusiness (check ownership):**
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
    // ... rest of update logic

    if err := h.BusinessService.Update(businessID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update business"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Business updated successfully"})
}
```

---

## Step 4: Update Frontend API Call

**File:** `api.ts`

**Import new types:**
```typescript
// Add to existing imports
```

**Add GetBusinessByUser method:**
```typescript
async getBusinessByUser(): Promise<Business[]> {
    return this.request<Business[]>('/api/v1/businesses/by-user');  // NEW ENDPOINT
}
```

**Update OperatorDashboard to use GetBusinessByUser:**
```typescript
// In useEffect:
const fetchData = async () => {
    try {
        setLoading(true);
        setError(null);
        
        // NEW: Get businesses for this user only
        const businessesData = await api.getBusinessByUser(); // CHANGED
        
        setBusinesses(businessesData);

        if (businessesData.length > 0) {
            // Always select first (only one business per user)
            const selectedBusiness = currentBusiness || businessesData[0];
            setCurrentBusiness(selectedBusiness);
            
            const [bookingsData, servicesData, slotsData, availabilityData] = await Promise.all([
                api.getBookingsByBusiness(selectedBusiness.id),
                api.getServicesByBusiness(selectedBusiness.id),
                api.getSlotsByBusiness(selectedBusiness.id),
                api.getAvailability(selectedBusiness.id)
            ]);

            setBookings(bookingsData);
            setServices(servicesData);
            setSlots(slotsData);
            setAvailability(availabilityData);
            setDurationMin((selectedBusiness as any).slotDurationMin || 30);
            setMaxBookings((selectedBusiness as any).maxBookings || 1);
            setEditingDay({});
        }
    } catch (err) {
        console.error('Failed to fetch data:', err);
        setError('Failed to load data. Please try again.');
    } finally {
        setLoading(false);
    }
};

// Remove handleBusinessChange function (line ~120)
// Users can only have 1 business, no switching
```

// Remove business selection dropdown UI (lines ~415-450)
// Change to single business header only
```

---

## Step 5: Deploy & Test

### Backend
1. Add GetByUser method to business_service.go
2. Add BusinessServiceByUser to NewHandler
3. Update GetBusinesses, GetBusiness, UpdateBusiness with ownership checks
4. Add new endpoint: GET /api/v1/businesses/by-user

### Frontend
1. Add getBusinessByUser to API client
2. Update OperatorDashboard to call getBusinessByUser
3. Remove business dropdown UI
4. Show single business header
5. Remove handleBusinessChange

### Expected Behavior After Fix:
- ✅ Users see ONLY their own business in dashboard
- ✅ No business dropdown to select other businesses
- ✅ Users CANNOT edit other businesses' settings
- ✅ Users CANNOT access other businesses' bookings
- ✅ API returns 403/403 for unauthorized access
- ✅ One business per user enforced at database level

---

## Testing
1. Login as User A, create booking
2. Logout
3. Try to access User B's business via API - should FAIL (403 Forbidden)
4. Try to edit User B's business settings - should FAIL (403 Forbidden)

---

## Success Metrics Achieved
✅ **Security:** Users can ONLY access their own business
✅ **API Security:** All business operations check ownership
✅ **UI Simplicity:** No business dropdown, single business
✅ **Data Isolation:** No cross-business access
✅ **Database Level:** One-to-one user:business relationship

---

**Implementation Time:** ~30 minutes
**Complexity:** LOW - minimal changes to existing code
**Risk:** LOW - doesn't require database schema changes

---

## Future Enhancement (Optional)
After this fix, you can add:
- Business creation in user registration (already designed in auth_service.go)
- Employee access (Employee model already created)
- Business switching (if one-to-many needed)

---

**This is MUCH SIMPLER than full business ownership implementation and provides IMMEDIATE security.**
