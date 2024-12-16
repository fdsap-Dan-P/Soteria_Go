package request

type (
	UserRegistrationRequest struct {
		Username         string `json:"username"`
		First_name       string `json:"first_name"`
		Middle_name      string `json:"middle_name"`
		Last_name        string `json:"last_name"`
		Email            string `json:"email"`
		Phone_no         string `json:"phone_no"`
		Staff_id         string `json:"staff_id"`
		Institution_code string `json:"institution_code"`
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
		First_name       string `json:"first_ame,omitempty"`
		Last_name        string `json:"last_name,omitempty"`
		Birthdate        string `json:"birthdate,omitempty"`
		Phone_no         string `json:"phone_no,omitempty"`
		Cid              string `json:"cid,omitempty"`
		Institution_code string `json:"institution_code,omitempty"`
	}
)
