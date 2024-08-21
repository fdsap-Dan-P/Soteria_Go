package request

type (
	UserRegistrationRequest struct {
		Username string `json:"username"`
		Staff_id string `json:"staff_id"`
	}
)
