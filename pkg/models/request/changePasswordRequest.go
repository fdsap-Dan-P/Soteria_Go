package request

type ForgotPassword struct {
	Username string `json:"username"`
}

type ResetPasswordIfExpiredRequest struct {
	Otp               string `json:"otp"`
	Previous_password string `json:"previous_password"`
	New_password      string `json:"new_password"`
	Confirm_password  string `json:"confirm_password"`
	Institution_code  string `json:"institution_code"`
	Application_code  string `json:"application_code"`
}

type ResetPasswordRequest struct {
	Otp              string `json:"otp"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
	Institution_code string `json:"institution_code"`
	Application_code string `json:"application_code"`
}

type OTPRequest struct {
	Otp              string `json:"otp"`
	Institution_code string `json:"institution_code"`
	Application_code string `json:"application_code"`
}

type SendingOTPRequest struct {
	Institution_code string `json:"institution_code"`
	Application_code string `json:"application_code"`
}
