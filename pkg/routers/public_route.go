package routers

import (
	"soteria_go/pkg/controllers/healthchecks"
	loginmanagement "soteria_go/pkg/controllers/userLog/logInManagement"
	logoutmanagement "soteria_go/pkg/controllers/userLog/logOutManagement"
	usermanagement "soteria_go/pkg/controllers/userManagement"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	//////////////////////////////////////////////////////////////////////////////////////////////
	// Monitoring
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

	// Initialize default config (Assign the middleware to /metrics)
	v1Endpoint.Get("/monitor", monitor.New())
	auth := v1Endpoint.Group("/auth")

	// Or extend your config for customization
	// Assign the middleware to /metrics
	// and change the Title to `MyService Metrics Page`
	v1Endpoint.Get("/monitor", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	//--- U S E R    L O G S ---//
	userLogs := v1Endpoint.Group("/user-logs")
	userLogs.Post("/login", loginmanagement.MainUserLogIn)
	userLogs.Get("/logout/:username/:session_id", logoutmanagement.MainLogOut)

	//--- U S E R    M A N A G E M E N T ---//
	userManagement := auth.Group("/user-management")
	userManagement.Post("/register", usermanagement.RegisterUserFromHCIS)
}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health checkss
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)
}
