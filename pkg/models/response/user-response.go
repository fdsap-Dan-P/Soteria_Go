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
		Password                string `json:"password,omitempty"`
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

	HeaderValidationResponse struct {
		Username   string `json:"username"`
		Insti_code string `json:"insti_code"`
		Insti_name string `json:"insti_name"`
		App_code   string `json:"app_code"`
		App_name   string `json:"app_name"`
	}

	LastPasswordUsed struct {
		Password_hash string `json:"password_hash"`
	}

	MemberVerificationResponse struct {
		Phone_no             string      `json:"phone_no"`
		First_name           string      `json:"first_name"`
		Last_name            string      `json:"last_name"`
		Birthdate            string      `json:"birthdate"`
		Is_member            bool        `json:"is_member"`
		Institution_code     string      `json:"institution_code"`
		Institution_name     string      `json:"institution_name"`
		No_phone_number_user int         `json:"no_phone_number_user"`
		Member_details       interface{} `json:"member_details"`
	}
)

type UserApplicationDetails struct {
	Username                string `json:"username"`
	Staff_id                string `json:"staff_id"`
	Email                   string `json:"email"`
	Phone_no                string `json:"phone_no"`
	Application_code        string `json:"application_code"`
	Application_name        string `json:"application_name"`
	Application_description string `json:"application_description"`
	Institution_code        string `json:"institution_code"`
	Institution_name        string `json:"institution_name"`
}
