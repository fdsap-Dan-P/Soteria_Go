package request

type (
	UserRegistrationRequest struct {
		Username string `json:"username"`
		Staff_id string `json:"staff_id"`
	}

	LoginCredentialsRequest struct {
		User_identity string `json:"user_identity"`
		Password      string `json:"password"`
	}

	ChangePasswordRequest struct {
		Old_password string `json:"old_password"`
		New_password string `json:"new_password"`
	}
)
