package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for the CORS middleware.
type CORSConfig struct {
	AllowOrigins     []string // Default: from CORS_ALLOWED_ORIGINS env, or ["*"]
	AllowMethods     []string // Default: ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
	AllowHeaders     []string // Default: ["Origin","Content-Type","Accept","Authorization","X-Request-ID","X-CSRF-Token"]
	ExposeHeaders    []string // Default: ["Content-Length","X-Request-ID"]
	AllowCredentials bool     // Default: true
	MaxAge           int      // Default: 43200 (12 hours), in seconds
}

// defaultCORSConfig returns the default CORS configuration.
func defaultCORSConfig() CORSConfig {
	origins := []string{"*"}
	if env := os.Getenv("CORS_ALLOWED_ORIGINS"); env != "" {
		origins = strings.Split(env, ",")
	}

	return CORSConfig{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           43200,
	}
}

// CORS returns middleware that handles cross-origin requests.
// If no config is provided, sensible defaults are used.
func CORS(configs ...CORSConfig) gin.HandlerFunc {
	cfg := defaultCORSConfig()
	if len(configs) > 0 {
		cfg = configs[0]
	}

	origins := strings.Join(cfg.AllowOrigins, ", ")
	methods := strings.Join(cfg.AllowMethods, ", ")
	headers := strings.Join(cfg.AllowHeaders, ", ")
	maxAge := fmt.Sprintf("%d", cfg.MaxAge)

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", maxAge)

		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if len(cfg.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
