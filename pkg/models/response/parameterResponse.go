package response

type (
	ConfigDetails struct {
		Config_id          int    `json:"config_id"`
		Config_code        string `json:"config_code"`
		Config_name        string `json:"config_name"`
		Config_description string `json:"config_description"`
		Config_type        string `json:"config_type"`
		Config_value       string `json:"config_value"`
		Config_insti_code  string `json:"config_insti_code"`
		Config_app_code    string `json:"config_app_code"`
		Created_at         string `json:"created_at"`
		Updated_at         string `json:"updated_at"`
	}
)
