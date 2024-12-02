package routers

import (
	"soteria_go/pkg/controllers/healthchecks"
	securitymanagement "soteria_go/pkg/controllers/security-management"
	setparameters "soteria_go/pkg/controllers/security-management/set-parameters"
	setuserpassword "soteria_go/pkg/controllers/security-management/set-user-password.go"
	userlogs "soteria_go/pkg/controllers/user-logs"
	usermanagement "soteria_go/pkg/controllers/user-management"
	registernewuser "soteria_go/pkg/controllers/user-management/register-new-user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	appName := app.Group("/cagabay-ua")
	apiEndpoint := appName.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	//////////////////////////////////////////////////////////////////////////////////////////////
	// Monitoring
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)
	v1Endpoint.Get("/monitor", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	// Initialize default config (Assign the middleware to /metrics)
	v1Endpoint.Get("/monitor", monitor.New())
	auth := v1Endpoint.Group("/auth")

	//--- U S E R    L O G S ---//
	userLogs := auth.Group("/user-logs")
	userLogs.Post("/login", userlogs.Login)
	userLogs.Get("/:username/logout", userlogs.LogOut)

	//--- U S E R    M A N A G E M E N T ---//
	userManagement := auth.Group("/user-management")
	userManagement.Post("/hcis-inquiry", registernewuser.HCISUserDetailsProvider)
	userManagement.Post("/register-new-user", registernewuser.RegisterUser)
	userManagement.Post("/update-user/:user_identity", registernewuser.UpdateUserDetails)
	userManagement.Get("/delete-user/:user_identity", usermanagement.DeleteUser)

	//--- S E C U R I T Y    M A N A G E M E N T ---//
	secManagement := auth.Group("/security-management")
	secManagement.Get("/validate-header", securitymanagement.ThirdPartyHeaderValidation)
	v1Endpoint.Post("/register-application", securitymanagement.AppRegistration) // no validation
	secManagement.Post("/change-password/:username", setuserpassword.UserInitiatedPasswordChange)
	secManagement.Post("/expire-password/:username", setuserpassword.UserChangePasswordAfterExpired)
	secManagement.Get("/reset-password/:username", setuserpassword.ResetUserPasswordToTemporary)

	// Set Parameters
	setParams := secManagement.Group("/parameters")
	setParams.Post("/update", setparameters.SetParams)
	setParams.Get("/list", setparameters.ParameterList)

	// OttoKonek Rose
	member := auth.Group("/member")
	member.Post("/verify", usermanagement.MemberVerification)
}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health checkss
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)
}
