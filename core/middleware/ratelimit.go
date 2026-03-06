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
