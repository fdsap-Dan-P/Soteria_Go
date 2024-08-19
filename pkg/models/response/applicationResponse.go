package response

type ApplicationResponse struct {
	App_id          int    `json:"app_id"`
	App_code        string `json:"app_code"`
	App_name        string `json:"app_name"`
	App_description string `json:"app_description"`
	Created_at      string `json:"created_at"`
	Updated_at      string `json:"updated_at"`
}
