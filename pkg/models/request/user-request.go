package request

type (
	UserRegistrationRequest struct {
		Username         string `json:"username"`
		First_name       string `json:"first_name,omitempty"`
		Middle_name      string `json:"middle_name,omitempty"`
		Last_name        string `json:"last_name,omitempty"`
		Email            string `json:"email,omitempty"`
		Phone_no         string `json:"phone_no,omitempty"`
		Staff_id         string `json:"staff_id,omitempty"`
		Institution_code string `json:"institution_code,omitempty"`
		Birthdate        string `json:"birthdate,omitempty"`
	}

	LoginCredentialsRequest struct {
		User_identity string `json:"user_identity"`
		Password      string `json:"password"`
	}

	ChangePasswordRequest struct {
		Old_password string `json:"old_password"`
		New_password string `json:"new_password"`
	}

	MemberVerificationRequest struct {
		First_name       string `json:"first_name,omitempty"`
		Last_name        string `json:"last_name,omitempty"`
		Birthdate        string `json:"birthdate,omitempty"`
		Phone_no         string `json:"phone_no,omitempty"`
		Cid              string `json:"cid,omitempty"`
		Institution_code string `json:"institution_code,omitempty"`
	}
)
