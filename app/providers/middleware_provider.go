package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/middleware"
)

// MiddlewareProvider registers built-in middleware aliases and groups.
type MiddlewareProvider struct{}

// Register is a no-op — middleware has no singleton to register.
func (p *MiddlewareProvider) Register(c *container.Container) {}

// Boot registers built-in middleware aliases and default groups.
func (p *MiddlewareProvider) Boot(c *container.Container) {
	middleware.RegisterAlias("recovery", middleware.Recovery())
	middleware.RegisterAlias("requestid", middleware.RequestID())
	middleware.RegisterAlias("cors", middleware.CORS())
	middleware.RegisterAlias("error_handler", middleware.ErrorHandler())
	middleware.RegisterAlias("auth", middleware.AuthMiddleware())

	middleware.RegisterGroup("global",
		middleware.Recovery(),
		middleware.RequestID(),
	)
}
