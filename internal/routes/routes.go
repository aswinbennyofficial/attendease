// Update your routes.go file
package routes

import (
	"net/http"

	"github.com/aswinbennyofficial/attendease/internal/controllers"
	"github.com/aswinbennyofficial/attendease/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux) {
	r.Get("/health", controllers.HandleHealth)

	r.With(middleware.LoginRequired).Get("/welcome", controllers.HandleWelcome)

	r.Post("/api/users/login", controllers.HandleSignin)
	r.Post("/api/users/signup", controllers.HandleSignup)
	r.Post("/api/users/refresh", controllers.HandleRefresh)
	r.Post("/api/users/logout", controllers.HandleLogout)

	r.Get("/cookie/check-cookie", func(w http.ResponseWriter, r *http.Request) {
		// Check if the "JWtoken" cookie is present
		cookie, err := r.Cookie("JWtoken")
		if err != nil {
			// Cookie not present
			w.Write([]byte("Cookie 'JWtoken' not found",))
			return
		}

		// Cookie present
		w.Write([]byte("Cookie 'JWtoken' found with value: " + cookie.Value))
	})

}
