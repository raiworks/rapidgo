package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raiworks/rapidgo/v2/core/router"
	"gorm.io/gorm"
)

// Routes registers liveness and readiness health-check endpoints.
// The dbFn callback defers database resolution until the first request.
// An optional version string adds version info to the liveness response.
func Routes(r *router.Router, dbFn func() *gorm.DB, version ...string) {
	r.Get("/health", func(c *gin.Context) {
		resp := gin.H{"status": "ok"}
		if len(version) > 0 && version[0] != "" {
			resp["version"] = version[0]
		}
		c.JSON(http.StatusOK, resp)
	})

	r.Get("/health/ready", func(c *gin.Context) {
		sqlDB, err := dbFn().DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "db": "connected"})
	})
}
