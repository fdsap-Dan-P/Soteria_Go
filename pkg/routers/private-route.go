package routers

import (
	"soteria_go/pkg/controllers/healthchecks"
	middleware "soteria_go/pkg/utils"
	"time"

	fiberUtils "soteria_go/pkg/utils/go-utils/fiber"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App) {

	// JWT Setup
	app.Use(fiberUtils.AuthenticationMiddleware(fiberUtils.JWTConfig{
		Duration:     15 * time.Minute,
		CookieMaxAge: 15 * 60,
		SetCookies:   true,
		SecretKey:    []byte(middleware.GetEnv("SECRET_KEY")),
	}))

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/private")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

}
