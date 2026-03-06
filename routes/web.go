package routes

import (
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/http/controllers"
)

// RegisterWeb defines web (HTML) routes.
func RegisterWeb(r *router.Router) {
	r.Get("/", controllers.Home)
}
