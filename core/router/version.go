package router

import "github.com/gin-gonic/gin"

// Version creates a route group prefixed with /api/{version}.
// The returned RouteGroup supports all existing route registration methods.
//
// Example:
//
//	v1 := r.Version("v1") // prefix: /api/v1
//	v1.Get("/users", listUsers)
//	v1.APIResource("/posts", &PostController{})
func (r *Router) Version(version string) *RouteGroup {
	return r.Group("/api/" + version)
}

// DeprecatedVersion creates a versioned route group (like Version) but injects
// middleware that adds deprecation headers to every response:
//   - Sunset: {sunsetDate} — RFC 8594 sunset date in HTTP-date format
//   - X-API-Deprecated: true — simple boolean signal for clients/monitoring
//
// Example:
//
//	v1 := r.DeprecatedVersion("v1", "Sat, 01 Jun 2026 00:00:00 GMT")
//	v1.Get("/users", listUsersV1)
func (r *Router) DeprecatedVersion(version, sunsetDate string) *RouteGroup {
	g := r.Version(version)
	g.Use(deprecationHeaders(sunsetDate))
	return g
}

// deprecationHeaders returns middleware that sets Sunset and X-API-Deprecated
// headers on every response.
func deprecationHeaders(sunsetDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Sunset", sunsetDate)
		c.Header("X-API-Deprecated", "true")
		c.Next()
	}
}
