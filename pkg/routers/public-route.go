package routers

import (
	"soteria_go/pkg/controllers/healthchecks"
	securitymanagement "soteria_go/pkg/controllers/security-management"
	setparameters "soteria_go/pkg/controllers/security-management/set-parameters"
	setuserpassword "soteria_go/pkg/controllers/security-management/set-user-password.go"
	userlogs "soteria_go/pkg/controllers/user-logs"
	usermanagement "soteria_go/pkg/controllers/user-management"
	"soteria_go/pkg/controllers/user-management/memberVerification"
	registernewuser "soteria_go/pkg/controllers/user-management/register-new-user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	appName := app.Group("/soteria-go")
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
	userLogs.Post("/logout", userlogs.LogOut)

	//--- U S E R    M A N A G E M E N T ---//
	userManagement := auth.Group("/user-management")
	userManagement.Post("/hcis-inquiry", registernewuser.HCISUserDetailsProvider)
	userManagement.Post("/register-new-user/staff", registernewuser.StaffRegistration)
	userManagement.Post("/register-new-user/non-staff", registernewuser.NonStaffRegistraion)
	userManagement.Post("/update-user/:user_category/:user_identity", registernewuser.UpdateUserDetails)
	userManagement.Post("/delete-user", usermanagement.DeleteUser)

	//--- S E C U R I T Y    M A N A G E M E N T ---//
	secManagement := auth.Group("/security-management")
	secManagement.Get("/validate-header", securitymanagement.ThirdPartyHeaderValidation)
	secManagement.Post("/register-application", securitymanagement.AppRegistration)           // admin only
	secManagement.Post("/retrieve-api-key/:app-code", securitymanagement.RetrievePlainApiKey) // admin only
	secManagement.Get("/encrypt-api-key", securitymanagement.EncryptApiKey)
	secManagement.Post("/change-password", setuserpassword.UserInitiatedPasswordChange)
	secManagement.Post("/expire-password", setuserpassword.UserChangePasswordAfterExpired)
	secManagement.Post("/reset-password", setuserpassword.ResetUserPasswordToTemporary)

	// Set Parameters
	setParams := secManagement.Group("/parameters")
	setParams.Post("/update", setparameters.SetParams)
	setParams.Get("/list", setparameters.ParameterList)

	// OttoKonek Rose
	member := auth.Group("/member")
	member.Post("/verify", memberVerification.MemberVerification)
}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health checkss
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)
}
