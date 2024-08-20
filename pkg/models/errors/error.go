package errors

type ErrorModel struct {
	Message   string `json:"message,omitempty"`
	IsSuccess bool   `json:"isSuccess,omitempty"`
	Error     error  `json:"error,omitempty"`
	// User_id    int    `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Staff_id string `json:"staff_id,omitempty"`
	// Session_id string `json:"session_id,omitempty"`
	// Is_active  bool   `json:"is_active,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}
