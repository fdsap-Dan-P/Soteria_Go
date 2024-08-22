package response

type (
	UserDetails struct {
		User_id                 int    `json:"user_id"`
		Username                string `json:"username"`
		Staff_id                string `json:"staff_id"`
		First_name              string `json:"first_name"`
		Middle_name             string `json:"middle_name"`
		Last_name               string `json:"last_name"`
		Email                   string `json:"email"`
		Phone_no                string `json:"phone_no"`
		Last_login              string `json:"last_login"`
		Institution_id          int    `json:"institution_id,omitempty"`
		Institution_code        string `json:"institution_code,omitempty"`
		Institution_name        string `json:"institution_name,omitempty"`
		Birthdate               string `json:"birthdate,omitempty"`
		Requires_password_reset bool   `json:"requires_password_reset,omitempty"`
		Last_password_reset     string `json:"last_password_reset,omitempty"`
		Token                   string `json:"token,omitempty"`
	}

	UserPasswordDetails struct {
		User_id                 int    `json:"user_id"`
		Password_hash           string `json:"password_hash"`
		Requires_password_reset bool   `json:"requires_password_reset"`
		Last_password_reset     string `json:"last_password_reset"`
		Created_at              string `json:"created_at"`
		Updated_at              string `json:"updated_at"`
	}

	UserTokenDetails struct {
		Token_id   int    `json:"token_id"`
		Username   string `json:"username"`
		Staff_id   string `json:"staff_id"`
		Token      string `json:"token"`
		Insti_code string `json:"insti_code"`
		App_code   string `json:"app_code"`
	}
)
