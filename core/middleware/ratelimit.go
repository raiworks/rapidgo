package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rate := os.Getenv("RATE_LIMIT")
	if rate == "" {
		rate = "60-M" // 60 requests per minute
	}
	r, _ := limiter.NewRateFromFormatted(rate)
	store := memory.NewStore()
	instance := limiter.New(store, r)
	return mgin.NewMiddleware(instance)
}

// RateLimitConfig allows customising the rate limiter.
type RateLimitConfig struct {
	// Rate in ulule format, e.g. "100-M" (100 per minute). Default: "60-M".
	Rate string
	// KeyFunc extracts the rate-limit key from the request.
	// If nil, the client IP is used (default behaviour).
	KeyFunc func(c *gin.Context) string
}

// RateLimitWithConfig creates a rate-limit middleware with the given config.
func RateLimitWithConfig(cfg RateLimitConfig) gin.HandlerFunc {
	rate := cfg.Rate
	if rate == "" {
		rate = "60-M"
	}
	r, _ := limiter.NewRateFromFormatted(rate)
	store := memory.NewStore()

	opts := []limiter.Option{}
	if cfg.KeyFunc != nil {
		opts = append(opts, limiter.WithClientIPHeader(""))
	}
	instance := limiter.New(store, r, opts...)

	if cfg.KeyFunc == nil {
		return mgin.NewMiddleware(instance)
	}

	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)
		ctx, err := instance.Get(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "rate limiter error"})
			return
		}
		if ctx.Reached {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
