# Dokploy Domain Setup Guide

This guide explains how to configure your domains for Blytz.Cloud using Dokploy.

## Overview

We've updated the `docker-compose.yml` to use Dokploy's recommended approach for domain management (Method 1 from [Dokploy Domains Documentation](https://docs.dokploy.com/docs/core/docker-compose/domains)).

## What Changed

1. **Removed manual Traefik labels** - Dokploy will handle these automatically
2. **Added `dokploy-network`** - Required for Dokploy routing
3. **Updated frontend API URL** - Now uses HTTPS: `https://api.blytz.cloud`

## Step-by-Step Setup

### 1. Deploy to Dokploy

Push your code to Git repository and connect it to Dokploy.

### 2. Create Docker Compose Application

In Dokploy:
- Create a new application
- Select "Docker Compose"
- Connect your Git repository
- Use the updated `docker-compose.yml`

### 3. Configure Domains

**For Frontend (blytz.cloud):**
1. Navigate to your Docker Compose application
2. Go to the **Domains** tab
3. Click **Add Domain**
4. Enter: `blytz.cloud`
5. Select the `blytz-cloud` service
6. Port: `80`
7. Enable HTTPS (if available)
8. Save

**For Backend (api.blytz.cloud):**
1. Click **Add Domain** again
2. Enter: `api.blytz.cloud`
3. Select the `backend` service
4. Port: `8080`
5. Enable HTTPS (if available)
6. Save

### 4. DNS Configuration

Create A records for your domains pointing to your server's IP:

```
Type: A
Name: @
Value: YOUR_SERVER_IP

Type: A
Name: api
Value: YOUR_SERVER_IP
```

### 5. Deploy

Click **Deploy** in Dokploy. Dokploy will automatically:
- Add the services to `dokploy-network`
- Generate Traefik labels for routing
- Configure SSL certificates (if enabled)

### 6. Verify

Test that everything is working:

```bash
# Frontend
curl https://blytz.cloud

# Backend health check
curl https://api.blytz.cloud/health
```

## Network Architecture

```
Internet
    ↓
Traefik (Dokploy)
    ↓
dokploy-network (external)
    ↓
├── blytz-frontend (port 80)
└── backend (port 8080)
        ↓
    blytz-network (internal)
        ↓
├── postgres
└── redis
```

## Troubleshooting

### Domain Not Working

1. **Check DNS:** Ensure A records are pointing to correct IP
   ```bash
   dig blytz.cloud
   dig api.blytz.cloud
   ```

2. **Check Traefik Dashboard:** View Traefik dashboard in Dokploy to see if routers are registered

3. **Preview Compose:** Use "Preview Compose" button in Dokploy to verify labels were added correctly

4. **Check Service Health:** Ensure backend and frontend containers are running
   ```bash
   docker ps | grep blytz
   ```

### HTTPS Not Working

1. Ensure DNS records are publicly accessible
2. Wait a few minutes for Let's Encrypt propagation
3. Check Traefik logs in Dokploy

### Backend Not Accessible

1. Verify `VITE_API_URL` is set correctly in frontend environment
2. Check CORS settings in backend (currently allows all origins - restrict in production)
3. Ensure backend is healthy: `curl http://localhost:8080/health`

### Preview Compose

Click the "Preview Compose" button in Dokploy to see the final generated configuration. It should include Traefik labels like:

```yaml
labels:
  - traefik.enable=true
  - traefik.http.routers.{service-name}.rule=Host(`blytz.cloud`)
  - traefik.http.routers.{service-name}.entrypoints=websecure
  - traefik.http.services.{service-name}.loadbalancer.server.port=80
```

## Security Notes

### Current Configuration (Insecure)

The backend currently allows all CORS origins. For production, update `backend/cmd/server/main.go`:

```go
// Replace lines 44-56 with:
r.Use(func(c *gin.Context) {
    allowedOrigins := []string{"https://blytz.cloud"}
    origin := c.Request.Header.Get("Origin")

    // Allow only specific origins
    for _, allowedOrigin := range allowedOrigins {
        if origin == allowedOrigin {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            break
        }
    }

    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
```

### Environment Variables

Make sure to set strong passwords in Dokploy:

- `DB_PASSWORD` - Strong database password
- `JWT_SECRET` - Random 32+ character secret

## Alternative: Manual Traefik Configuration

If you prefer manual configuration instead of using Dokploy's domain management, modify the `docker-compose.yml` to add labels:

```yaml
backend:
  # ... existing config ...
  labels:
    - traefik.enable=true
    - traefik.http.routers.blytz-backend.rule=Host(`api.blytz.cloud`)
    - traefik.http.routers.blytz-backend.entrypoints=websecure
    - traefik.http.routers.blytz-backend.tls=true
    - traefik.http.routers.blytz-backend.tls.certresolver=letsencrypt
    - traefik.http.services.blytz-backend.loadbalancer.server.port=8080

blytz-cloud:
  # ... existing config ...
  labels:
    - traefik.enable=true
    - traefik.http.routers.blytz-frontend.rule=Host(`blytz.cloud`)
    - traefik.http.routers.blytz-frontend.entrypoints=websecure
    - traefik.http.routers.blytz-frontend.tls=true
    - traefik.http.routers.blytz-frontend.tls.certresolver=letsencrypt
    - traefik.http.services.blytz-frontend.loadbalancer.server.port=80
```

Note: This requires your Dokploy Traefik to be configured with a `letsencrypt` certificate resolver.
