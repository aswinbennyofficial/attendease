// Update your routes.go file
package routes

import (

	"github.com/aswinbennyofficial/attendease/internal/controllers"
	"github.com/aswinbennyofficial/attendease/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux) {
	r.Get("/health", controllers.HandleHealth)

	r.With(middleware.AdminLoginRequired).Get("/private", controllers.HandlePrivate)

	// API for organisation
	r.Post("/api/admin/login", controllers.HandleAdminSignin)
	r.Post("/api/admin/signup", controllers.HandleAdminSignup)
	r.Post("/api/admin/refresh", controllers.HandleRefresh)
	r.Post("/api/logout", controllers.HandleLogout)

	// API for events
	r.With(middleware.AdminLoginRequired).Post("/api/events", controllers.HandleCreateEvent) // Create event
	r.With(middleware.AdminLoginRequired).Get("/api/events", controllers.HandleGetEvents)   // Get all events
	r.With(middleware.AdminLoginRequired).Get("/api/events/{eventid}", controllers.HandleGetAnEvent) // Get event by eventid

	

}
