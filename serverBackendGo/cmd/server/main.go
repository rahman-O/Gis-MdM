// @title Headwind MDM API (Go)
// @version 1.0
// @description Gradual Go migration of Headwind MDM REST API.
// @BasePath /rest
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
