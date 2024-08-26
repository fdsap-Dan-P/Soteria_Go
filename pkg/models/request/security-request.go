package request

type (
	TokenValidationRequest struct {
		Token string `json:"token"`
	}

	ApplicationRequest struct {
		App_name string `json:"app_name"`
		App_desc string `json:"app_desc"`
	}

	ParameterRequest struct {
		No_of_minutes    string `json:"no_of_minutes"`
		Institution_code string `json:"institution_code"`
		Application_code string `json:"application_code"`
	}
)
