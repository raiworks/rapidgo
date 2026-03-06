package router

import (
	"html/template"
	"net/http"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/gin-gonic/gin"
)

// Router wraps the Gin engine and provides framework-level route registration.
type Router struct {
	engine *gin.Engine
}

// New creates a new Router with Gin mode set based on APP_ENV.
func New() *Router {
	setGinMode()
	engine := gin.New()
	return &Router{engine: engine}
}

// setGinMode configures Gin's mode based on the APP_ENV environment variable.
func setGinMode() {
	switch config.AppEnv() {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// Engine returns the underlying Gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// ServeHTTP implements http.Handler, delegating to the Gin engine.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// Get registers a GET route.
func (r *Router) Get(path string, handlers ...gin.HandlerFunc) {
	r.engine.GET(path, handlers...)
}

// Post registers a POST route.
func (r *Router) Post(path string, handlers ...gin.HandlerFunc) {
	r.engine.POST(path, handlers...)
}

// Put registers a PUT route.
func (r *Router) Put(path string, handlers ...gin.HandlerFunc) {
	r.engine.PUT(path, handlers...)
}

// Delete registers a DELETE route.
func (r *Router) Delete(path string, handlers ...gin.HandlerFunc) {
	r.engine.DELETE(path, handlers...)
}

// Patch registers a PATCH route.
func (r *Router) Patch(path string, handlers ...gin.HandlerFunc) {
	r.engine.PATCH(path, handlers...)
}

// Options registers an OPTIONS route.
func (r *Router) Options(path string, handlers ...gin.HandlerFunc) {
	r.engine.OPTIONS(path, handlers...)
}

// Group creates a new route group with a shared prefix and optional middleware.
func (r *Router) Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{group: r.engine.Group(prefix, handlers...)}
}

// Use adds global middleware to the router.
func (r *Router) Use(middleware ...gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

// SetFuncMap sets the template function map on the Gin engine.
// Must be called before LoadTemplates.
func (r *Router) SetFuncMap(funcMap template.FuncMap) {
	r.engine.SetFuncMap(funcMap)
}

// LoadTemplates loads HTML templates matching the glob pattern.
func (r *Router) LoadTemplates(pattern string) {
	r.engine.LoadHTMLGlob(pattern)
}

// Static serves files from a local directory under the given URL path.
func (r *Router) Static(urlPath, dirPath string) {
	r.engine.Static(urlPath, dirPath)
}

// StaticFile serves a single file at the given URL path.
func (r *Router) StaticFile(urlPath, filePath string) {
	r.engine.StaticFile(urlPath, filePath)
}

// Run starts the HTTP server on the given address.
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
