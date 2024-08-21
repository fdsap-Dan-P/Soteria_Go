package response

type (
	ApplicationDetails struct {
		App_id                  int    `json:"app_id"`
		Application_code        string `json:"app_code"`
		Application_name        string `json:"application_name"`
		Application_description string `json:"application_description"`
		Api_key                 string `json:"api_key"`
		Created_at              string `json:"created_at"`
		Updated_at              string `json:"updated_at"`
	}
)
