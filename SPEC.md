---
type: spec
title: Booking System Template - Specification
resource: blytz-booking
description: "**Purpose:** Reusable white-label booking system template for service businesses (workshops, clinics, salons, etc.)"
tags: [go, react, tailwind]
updated: 2026-06-18
---

# Booking System Template - Specification

## Project Summary

**Purpose:** Reusable white-label booking system template for service businesses (workshops, clinics, salons, etc.)
**Stack:** Laravel 12 + React (Inertia.js) + MySQL + CHIP Collect
**Deployment:** DigitalOcean Droplet per client
**Storage:** GitHub

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Laravel 12 (PHP 8.3+) |
| Frontend | React + Inertia.js + Tailwind CSS |
| Database | MySQL 8.0 |
| Payment | CHIP Collect (aliff/chip-in) |
| Email | Laravel Mail (SMTP/SendGrid) |
| WhatsApp | Twilio API |
| Auth | Laravel Breeze (Sanctum) |
| Queue | Laravel Queue (database driver) |
| Storage | Local filesystem (logo uploads) |

---

## Part 1: Customer Booking Page

**URL:** `/booking`

### Step 1: Select Service
- Display service cards with: name, description, price, duration
- Filter by category (optional)
- Click to select → proceed to slot selection

### Step 2: Pick Time Slot
- Calendar view (month/week/day)
- Available slots shown in green
- Booked slots shown in red
- Click slot → proceed to details

### Step 3: Enter Details
- Name (required)
- Phone (required)
- Email (required)
- Vehicle/Notes (optional)
- Validate inputs

### Step 4: Pay via CHIP
- Display booking summary
- "Pay Now" button → redirect to CHIP checkout
- CHIP handles payment (FPX, credit card, etc.)
- Webhook callback updates booking status

### Step 5: Confirmation
- Success page with booking reference
- Email confirmation sent
- WhatsApp confirmation sent (if configured)

---

## Part 2: Back Office

### Dashboard
- Today's bookings count
- Revenue today (paid bookings)
- Upcoming bookings list
- Quick stats (week/month)

### Calendar View
- Day/week/month view
- Color-coded: pending, paid, completed, cancelled
- Click booking to view details
- Drag to reschedule (optional)

### Bookings Management
- List all bookings
- Filter by: date, status, service
- Actions: view, edit, cancel, mark complete
- Export to CSV (optional)

### Services Management
- CRUD services
- Set: name, description, price, duration, category
- Active/inactive toggle

### Time Slots
- Set business hours (Mon-Sun)
- Add blocked dates (holidays)
- Set slot duration (15/30/60 min)
- Buffer time between bookings

### Customers
- View customer list
- View booking history per customer
- Search by name/phone/email

### Reports
- Daily/weekly/monthly bookings
- Revenue report
- Service popularity
- Export PDF/CSV

### Staff/Users
- Add/manage staff accounts
- Roles: admin, staff (view only)
- Activity log (optional)

---

## Part 3: Client Admin (Settings)

### Branding
- Upload logo (PNG/JPG, max 2MB)
- Primary color picker
- Business name
- Business address
- Business phone

### Business Hours
- Set open/close time per day
- Closed days toggle
- Holiday dates

### Services
- Add/edit/delete services
- Set pricing (MYR)
- Set duration
- Set category

### Payment
- CHIP API Key
- CHIP Brand ID
- Sandbox/Live mode toggle

### Notifications
- Email: SMTP settings (host, port, username, password)
- WhatsApp: Twilio SID, Auth Token, From number
- Enable/disable each notification type

### Users
- Add/manage admin accounts
- Reset passwords

---

## Part 4: Notifications

| Trigger | Email | WhatsApp |
|---------|-------|----------|
| Booking confirmed (paid) | ✅ | ✅ (optional) |
| Booking reminder (24h before) | ✅ | ✅ (optional) |
| Booking cancelled | ✅ | ✅ (optional) |
| Booking completed | ✅ | ❌ |

---

## Database Schema

### bookings
- id, service_id, customer_name, customer_phone, customer_email, vehicle_notes, slot_datetime, status (pending/paid/completed/cancelled), chip_payment_id, amount, created_at, updated_at

### services
- id, name, description, price, duration_minutes, category, is_active, created_at, updated_at

### time_slots
- id, date, start_time, end_time, is_booked, is_blocked, created_at, updated_at

### business_settings
- id, key, value, created_at, updated_at

### users
- id, name, email, password, role (admin/staff), created_at, updated_at

### payments
- id, chip_id, booking_id, status, amount, currency, reference, payload, paid_at, created_at, updated_at

---

## API Endpoints

### Public
```
GET  /api/services          # List active services
GET  /api/slots             # Get available slots (date range)
POST /api/bookings          # Create booking
GET  /api/bookings/{id}     # Get booking status
```

### Webhook
```
POST /chipin/callback       # CHIP payment callback
```

### Protected (Back Office)
```
GET    /api/admin/dashboard      # Dashboard stats
GET    /api/admin/bookings       # List bookings
PUT    /api/admin/bookings/{id}  # Update booking
GET    /api/admin/services       # List services
POST   /api/admin/services       # Create service
PUT    /api/admin/services/{id}  # Update service
DELETE /api/admin/services/{id}  # Delete service
GET    /api/admin/settings       # Get settings
PUT    /api/admin/settings       # Update settings
```

---

## Payment Flow (CHIP)

1. Customer clicks "Pay Now"
2. Backend creates CHIP purchase:
   ```php
   $purchase = new Purchase(app(Client::class));
   return $purchase->createAndRedirect([
       'client' => ['email' => $booking->customer_email],
       'purchase' => [
           'products' => [[
               'name' => $service->name,
               'price' => $service->price * 100, // cents
               'quantity' => 1,
           ]],
       ],
       'reference' => $booking->id,
       'success_redirect' => route('booking.success', ['id' => $booking->id]),
       'failed_redirect' => route('booking.failed', ['id' => $booking->id]),
   ]);
   ```
3. Customer pays on CHIP checkout
4. CHIP sends webhook to `/chipin/callback`
5. Backend updates booking status to "paid"
6. Sends email + WhatsApp confirmation

---

## Deployment

### Per Client Setup
1. Clone template from GitHub
2. Create MySQL database
3. Copy `.env.example` to `.env`
4. Configure database credentials
5. Configure CHIP API keys
6. Configure email/SMTP
7. Configure WhatsApp/Twilio (optional)
8. Run migrations: `php artisan migrate`
9. Create admin user: `php artisan make:admin`
10. Upload logo in admin panel
11. Add services in admin panel
12. Set business hours

### Environment Variables
```env
APP_NAME="Booking System"
APP_ENV=production
APP_URL=https://client-domain.com

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=booking_db
DB_USERNAME=root
DB_PASSWORD=secret

CHIPIN_API_KEY=chip_api_key_here
CHIPIN_BRAND_ID=chip_brand_id_here
CHIPIN_MODE=live

MAIL_MAILER=smtp
MAIL_HOST=smtp.sendgrid.net
MAIL_PORT=587
MAIL_USERNAME=apikey
MAIL_PASSWORD=sendgrid_api_key
MAIL_FROM_ADDRESS=noreply@client-domain.com
MAIL_FROM_NAME="${APP_NAME}"

TWILIO_SID=twilio_sid_here
TWILIO_AUTH_TOKEN=twilio_token_here
TWILIO_WHATSAPP_FROM=whatsapp:+1234567890
```

---

## Estimated Build Timeline

| Phase | Duration | Details |
|-------|----------|---------|
| Setup + Auth | 1 week | Laravel install, Breeze auth, database |
| Customer Booking | 2 weeks | Service selection, slot picker, booking form |
| CHIP Integration | 1 week | Payment flow, webhooks, callbacks |
| Back Office | 3 weeks | Dashboard, calendar, bookings, services, reports |
| Admin Settings | 1 week | Branding, business hours, notifications |
| Email/WhatsApp | 1 week | Notification templates, queue jobs |
| Testing + Polish | 1 week | Bug fixes, UI polish, documentation |

**Total: ~10 weeks**

---

## File Structure

```
booking-template/
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── BookingController.php
│   │   │   ├── Admin/
│   │   │   │   ├── DashboardController.php
│   │   │   │   ├── BookingController.php
│   │   │   │   ├── ServiceController.php
│   │   │   │   ├── SettingController.php
│   │   │   │   └── ReportController.php
│   │   │   └── ChipInController.php
│   │   └── Requests/
│   ├── Models/
│   │   ├── Booking.php
│   │   ├── Service.php
│   │   ├── TimeSlot.php
│   │   ├── Payment.php
│   │   └── Setting.php
│   └── Services/
│       ├── BookingService.php
│       ├── SlotService.php
│       └── NotificationService.php
├── config/
│   └── chipin.php
├── database/
│   └── migrations/
├── resources/
│   ├── js/
│   │   ├── Pages/
│   │   │   ├── Booking/
│   │   │   │   ├── SelectService.jsx
│   │   │   │   ├── SelectSlot.jsx
│   │   │   │   ├── EnterDetails.jsx
│   │   │   │   ├── Payment.jsx
│   │   │   │   └── Success.jsx
│   │   │   └── Admin/
│   │   │       ├── Dashboard.jsx
│   │   │       ├── Bookings/
│   │   │       ├── Services/
│   │   │       ├── Calendar.jsx
│   │   │       ├── Customers.jsx
│   │   │       ├── Reports.jsx
│   │   │       └── Settings/
│   │   └── Components/
│   └── views/
│       ├── success.blade.php
│       └── failed.blade.php
├── routes/
│   ├── web.php
│   └── api.php
└── .env.example
```

---

## Next Steps

1. ✅ Specification complete
2. Create Laravel project
3. Install dependencies (Inertia, Breeze, CHIP)
4. Build database schema
5. Implement customer booking flow
6. Implement back office
7. Add CHIP payment
8. Add notifications
9. Test and deploy

---

**Status:** Specification ready for development
