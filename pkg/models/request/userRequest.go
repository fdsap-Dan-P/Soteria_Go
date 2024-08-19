package request

type LogInRequest struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	Institution_code string `json:"institution_code"`
	Application_code string `json:"application_code"`
}

type ValidateInHCIS struct {
	StaffID string `json:"StaffID"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Staff_id string `json:"staff_id"`
}
