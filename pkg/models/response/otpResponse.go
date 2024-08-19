package response

type OTPResponse struct {
	User_id    int    `json:"user_id,omitempty"`
	Otp        string `json:"otp,omitempty"`
	Method     string `json:"method,omitempty"`
	Via        string `json:"via,omitempty"`
	Status     string `json:"status,omitempty"`
	Created_at string `json:"created_at,omitempty"`
}

type ValidateOtpResponse struct {
	Date string `json:"date"`
	Msg  string `json:"msg"`
}
