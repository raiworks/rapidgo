package routes

import (
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/http/controllers"
)

// RegisterAPI defines API routes under the /api prefix.
func RegisterAPI(r *router.Router) {
	api := r.Group("/api")
	api.APIResource("/posts", &controllers.PostController{})
}
