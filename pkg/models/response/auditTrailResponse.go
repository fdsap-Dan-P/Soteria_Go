package response

type UserDetailsAuditTrail struct {
	User_id          int    `json:"user_id,omitempty"`
	Staff_id         string `json:"staff_id,omitempty"`
	Username         string `json:"username,omitempty"`
	User_status      string `jaon:"user_status,omitempty"`
	Institution_name string `json:"institution_name,omitempty"`
	User_action      string `json:"user_action,omitempty"`
	Old_value        string `json:"old_value,omitempty"`
	New_value        string `json:"new_value,omitempty"`
	Log_date         string `json:"log_date,omitempty"`
	Log_time         string `json:"log_time,omitempty"`
}
