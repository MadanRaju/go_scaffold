package handlers

import (
	"log"
	"net/http"

	"inventory-optimisation-server/internal/mid"
	"inventory-optimisation-server/internal/platform/auth"
	"inventory-optimisation-server/internal/platform/db"
	"inventory-optimisation-server/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, masterDB *db.DB, authenticator *auth.Authenticator) http.Handler {

	// authmw is used for authentication/authorization middleware.
	authmw := mid.Auth{
		Authenticator: authenticator,
	}

	app := web.New(log, mid.RequestLogger, mid.Metrics, mid.ErrorHandler)

	// Register health check endpoint. This route is not authenticated.
	h := Health{
		MasterDB: masterDB,
	}
	app.Handle("GET", "/v1/health", h.Check)

	// Register user management and authentication endpoints.
	u := User{
		MasterDB:       masterDB,
		TokenGenerator: authenticator,
	}

	app.Handle("GET", "/v1/users", u.List, authmw.Authenticate)
	app.Handle("POST", "/v1/users", u.Create, authmw.Authenticate)
	app.Handle("GET", "/v1/users/:id", u.Retrieve, authmw.Authenticate)
	app.Handle("PUT", "/v1/users/:id", u.Update, authmw.Authenticate)
	app.Handle("DELETE", "/v1/users/:id", u.Delete, authmw.Authenticate)

	// This route is not authenticated
	app.Handle("GET", "/v1/users/token", u.Token)

	o := OptimisationRequest{
		MasterDB: masterDB,
	}
	app.Handle("GET", "/v1/validate", o.Validate)
	app.Handle("POST", "/v1/optimisation-requests", o.Create)

	app.Handle("GET", "/v1/optimisation-requests", o.List)
	app.Handle("GET", "/v1/optimisation-requests/:id", o.Retrieve)

	return app
}
