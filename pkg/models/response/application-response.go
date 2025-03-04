package response

type (
	ApplicationDetails struct {
		Application_id          int    `json:"application_id"`
		Application_code        string `json:"application_code"`
		Application_name        string `json:"application_name"`
		Application_description string `json:"application_description"`
		Api_key                 string `json:"api_key"`
		Api_key_plain           string `json:"api_key_plain"`
		Created_at              string `json:"created_at"`
		Updated_at              string `json:"updated_at"`
	}

	UserAppResponse struct {
		User_id        int    `json:"user_id"`
		Application_id int    `json:"application_id"`
		Created_at     string `json:"created_at"`
		Updated_at     string `json:"updated_at"`
	}
)
