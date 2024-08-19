package response

type SystemConfigurationResponse struct {
	Config_id          int    `json:"config_id,omitempty"`
	Config_code        string `json:"config_code,omitempty"`
	Config_name        string `json:"config_name,omitempty"`
	Config_description string `json:"config_description,omitempty"`
	Config_type        string `json:"config_type,omitempty"`
	Config_value       string `json:"config_value,omitempty"`
	Updated_by         int    `json:"updated_by,omitempty"`
	Created_at         string `json:"created_at,omitempty"`
	Updated_at         string `json:"updated_at,omitempty"`
}

type IdleTimeResponse struct {
	Config_id          string `json:"config_id"`
	Config_code        string `json:"config_code"`
	Config_name        string `json:"config_name"`
	Config_description string `json:"config_description"`
	Config_type        string `json:"config_type"`
	Config_value       int    `json:"config_value"`
	Updated_by         int    `json:"updated_by"`
	Created_at         string `json:"created_at"`
	Updated_at         string `json:"updated_at"`
}
