// @title Headwind MDM API (Go)
// @version 1.0
// @description Gradual Go migration of Headwind MDM REST API.
//
// **Private endpoints** require authentication. In Swagger UI:
// 1. Call `POST /public/jwt/login` with `{"login":"admin","password":"admin"}` (raw password is accepted).
// 2. Copy the `Authorization` value from the response **headers** (format: `Bearer <jwt>`).
// 3. Click **Authorize**, paste that full value, then try `/private/*` requests.
//
// @BasePath /rest
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT from POST /public/jwt/login response header (e.g. Bearer eyJhbGciOiJIUzI1NiIs...)
package main

import (
	"log"
	"os"

	"github.com/gis-mdm/server-backend-go/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("server exited: %v", err)
		os.Exit(1)
	}
}
