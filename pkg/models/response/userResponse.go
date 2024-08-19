package response

type UserInfo struct {
	User_id             int    `json:"user_id"`
	Username            string `json:"username,omitempty"`
	Staff_id            string `json:"staff_id"`
	Email               string `json:"email,omitempty"`
	Phone_no            string `json:"phone_no,omitempty"`
	Status              string `json:"status,omitempty"`
	Last_login          string `json:"last_login,omitempty"`
	Last_login_date     string `json:"last_login_date,omitempty"`
	Last_login_time     string `json:"last_login_time,omitempty"`
	Nick_name           string `json:"nick_name,omitempty"`
	First_name          string `json:"first_name,omitempty"`
	Middle_name         string `json:"middle_name,omitempty"`
	Last_name           string `json:"last_name,omitempty"`
	Gender              string `json:"gender"`
	Birthdate           string `json:"birthdate,omitempty"`
	Age                 int    `json:"age,omitempty"`
	User_address        string `json:"user_address,omitempty"`
	Date_hired          string `json:"date_hired,omitempty"`
	Date_regularization string `json:"date_regularization,omitempty"`
	Employment_status   string `json:"employment_status,omitempty"`
	Employment_type     string `json:"employment_type,omitempty"`
	Job_level           string `json:"job_level,omitempty"`
	Job_position        string `json:"job_position,omitempty"`
	Job_grade           string `json:"job_grade,omitempty"`
	Institution_code    string `json:"institution_code,omitempty"`
	Institution_name    string `json:"institution_name,omitempty"`
	Area_code           string `json:"area_code,omitempty"`
	Area_name           string `json:"area_name,omitempty"`
	Unit_code           string `json:"unit_code,omitempty"`
	Unit_name           string `json:"unit_name,omitempty"`
	Center_code         string `json:"center_code,omitempty"`
	Center_name         string `json:"center_name,omitempty"`
	Role_code           string `json:"role_code,omitempty"`
	Role_name           string `json:"role_name,omitempty"`
	Activate            bool   `json:"activate,omitempty"`
	Deactivate          bool   `json:"deactivate,omitempty"`
	Unblock             bool   `json:"unblock,omitempty"`
	Unlock              bool   `json:"unlock,omitempty"`
}

type UserAccountResponse struct {
	User_id    int    `json:"user_id"`
	Username   string `json:"username,omitempty"`
	First_name string `json:"first_name,omitempty"`
	Last_name  string `json:"last_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone_no   string `json:"phone_no,omitempty"`
	Staff_id   string `json:"staff_id"`
	Is_active  bool   `json:"is_active,omitempty"`
	Status_id  int    `json:"status_id,omitempty"`
	Last_login string `json:"last_login,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type UserPasswordsResponse struct {
	User_id                 int    `json:"user_id"`
	Password_hash           string `json:"password_hash,omitempty"`
	Requires_password_reset bool   `json:"requires_password_reset,omitempty"`
	Last_password_reset     string `json:"last_password_reset,omitempty"`
	Created_at              string `json:"created_at,omitempty"`
	Updated_at              string `json:"updated_at,omitempty"`
}

type UserEmploymentResponse struct {
	User_id             int    `json:"user_id"`
	Region              string `json:"region"`
	Cluster             string `json:"cluster"`
	Department          string `json:"department"`
	Location_assignment string `json:"location_assignment"`
	Date_hired          string `json:"date_hired"`
	Date_regularization string `json:"date_regularization"`
	Employment_status   string `json:"employment_status"`
	Employment_type     string `json:"employment_type"`
	Job_level           string `json:"job_level"`
	Job_position        string `json:"job_position"`
	Job_grade           string `json:"job_grade"`
	Institution_code    string `json:"institution_code"`
	Area_code           string `json:"area_code"`
	Unit_code           string `json:"unit_code"`
	Center_code         string `json:"center_code"`
}

type UserProfilesResponse struct {
	User_id    int    `json:"user_id,omitempty"`
	First_name string `json:"first_name,omitempty"`
	Last_name  string `json:"last_name,omitempty"`
	Birthdate  string `json:"birthdate,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type LastUsed struct {
	Password_hash string `json:"password_hash"`
}

type LogInResponse struct {
	User_details   UserDetailUponLoggedIn `json:"user_details"`
	Administration Administration         `json:"administration"`
}

type SessionLogs struct {
	Session_id  string `json:"session_id"`
	User_id     int    `json:"user_id"`
	Login_time  string `json:"login_time"`
	Logout_time string `json:"logout_time"`
	Ip          string `json:"ip"`
	Client_info string `json:"client_info"`
	Is_active   bool   `json:"is_active"`
}

type UserDetailUponLoggedIn struct {
	User_id     int    `json:"user_id"`
	Full_name   string `json:"full_name,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone_no    string `json:"phone_no,omitempty"`
	Birthdate   string `json:"birthdate,omitempty"`
	Staff_id    string `json:"staff_id"`
	Role_id     int    `json:"role_id"`
	Role_code   string `json:"role_code"`
	Role_name   string `json:"role_name"`
	Branch_name string `json:"branch_name,omitempty"`
	Branch_code string `json:"branch_codes,omitempty"`
	Created_at  string `json:"created_at"`
	Last_login  string `json:"last_login"`
	Token       string `json:"token,omitempty"`
	Session_id  string `json:"session_id"`
}

type Administration struct {
	// User        Permission_struct `json:"user"`
	// Role        Permission_struct `json:"role"`
	// Security    Permission_struct `json:"security"`
	// Report      Permission_struct `json:"report"`
	// Dashboard   Permission_struct `json:"dashboard"`
	// Access      Permission_struct `json:"access"`
	User_access bool `json:"user_access"`
}

type UserStatusResponse struct {
	Status_id          int    `json:"status_id"`
	Status_code        string `json:"status_code"`
	Status_name        string `json:"status_name"`
	Status_description string `json:"status_description"`
	Created_at         string `json:"created_at,omitempty"`
	Updated_at         string `json:"updated_at,omitempty"`
}

type RestrictedResponse struct {
	User_id    int    `json:"user_id"`
	Ip         string `json:"ip"`
	Activity   string `json:"activity"`
	Created_at string `json:"created_at"`
}

type LogInSuccess struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Staff_id string `json:"staff_id"`
	Token    string `json:"token,omitempty"`
}
