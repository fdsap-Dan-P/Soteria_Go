package request

type (
	TokenValidationRequest struct {
		Token string `json:"token"`
	}

	ApplicationRequest struct {
		App_name string `json:"app_name"`
		App_desc string `json:"app_desc"`
	}
)