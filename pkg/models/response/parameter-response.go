package response

type (
	ConfigDetails struct {
		Config_id          int    `json:"config_id"`
		Config_code        string `json:"config_code,omitempty"`
		Config_name        string `json:"config_name,omitempty"`
		Config_description string `json:"config_description,omitempty"`
		Config_type        string `json:"config_type,omitempty"`
		Config_value       string `json:"config_value,omitempty"`
		Config_insti_code  string `json:"config_insti_code,omitempty"`
		Config_app_code    string `json:"config_app_code,omitempty"`
		Created_at         string `json:"created_at,omitempty"`
		Updated_at         string `json:"updated_at,omitempty"`
	}
)
