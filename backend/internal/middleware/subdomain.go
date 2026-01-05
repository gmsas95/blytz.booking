package middleware

import (
	"net/http"
	"strings"

	"blytz.cloud/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubdomainConfig struct {
	BaseDomain string
}

type SubdomainMiddleware struct {
	businessService *services.BusinessService
	config          *SubdomainConfig
}

func NewSubdomainMiddleware(businessService *services.BusinessService, config *SubdomainConfig) *SubdomainMiddleware {
	return &SubdomainMiddleware{
		businessService: businessService,
		config:          config,
	}
}

func (m *SubdomainMiddleware) ExtractAndValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host

		subdomain := m.extractSubdomain(host)

		if subdomain == "" {
			c.Set("subdomain", "")
			c.Set("business_id", uuid.Nil)
			c.Next()
			return
		}

		business, err := m.businessService.GetBySlug(subdomain)
		if err != nil {
			if err == services.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate business"})
			c.Abort()
			return
		}

		c.Set("subdomain", subdomain)
		c.Set("business_id", business.ID)
		c.Set("business_slug", business.Slug)
		c.Next()
	}
}

func (m *SubdomainMiddleware) extractSubdomain(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	host = strings.TrimSpace(host)

	baseDomain := m.config.BaseDomain

	if !strings.HasSuffix(host, baseDomain) {
		return ""
	}

	host = strings.TrimSuffix(host, baseDomain)
	host = strings.Trim(host, ".")

	if host == "" || strings.Contains(host, ".") {
		return ""
	}

	ignoredSubdomains := map[string]bool{
		"www":          true,
		"api":          true,
		"admin":        true,
		"app":          true,
		"localhost":    true,
		"127-0-0-1":    true,
		"127-0-0-1-ip": true,
	}

	if ignoredSubdomains[host] {
		return ""
	}

	if strings.HasPrefix(host, "192-168-") || strings.HasPrefix(host, "10-0-") {
		return ""
	}

	return host
}
