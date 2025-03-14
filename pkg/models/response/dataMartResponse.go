package response

type (
	MemberResponse struct {
		RetCode string     `json:"retCode"`
		Message string     `json:"message"`
		Data    MemberData `json:"data"`
	}

	MemberData struct {
		Message   string `json:"message"`
		IsSuccess bool   `json:"isSuccess"`
		// Error     error         `json:"error,omitempty"`
		No_phone_number_user int           `json:"no_phone_number_user"`
		Details              MemberDetails `json:"details,omitempty"`
	}

	MemberDetails struct {
		Cid           string `json:"cid"`
		First_name    string `json:"first_name"`
		Middle_name   string `json:"middle_name"`
		Last_name     string `json:"last_name"`
		Date_of_birth string `json:"date_of_birth"`
		Phone_1       string `json:"phone_1"`
		Phone_2       string `json:"phone_2"`
		Email         string `json:"email"`
		Center_code   string `json:"center_code"`
		Center_name   string `json:"center_name"`
		Unit_code     string `json:"unit_code"`
		Unit_name     string `json:"unit_name"`
		Branch_code   string `json:"branch_code"`
		Branch_name   string `json:"branch_name"`
		Ao_name       string `json:"ao_name"`
		Ao_staff_id   string `json:"ao_staff_id"`
		Insti_code    string `json:"insti_code"`
	}
)

type (
	MemberSavingResponse struct {
		RetCode string           `json:"retCode"`
		Message string           `json:"message"`
		Data    MemberSavingData `json:"data"`
	}

	MemberSavingData struct {
		Message   string              `json:"message"`
		IsSuccess bool                `json:"isSuccess"`
		Error     error               `json:"error,omitempty"`
		Details   MemberSavingDetails `json:"details,omitempty"`
	}

	MemberSavingDetails struct {
		Account_number   string `json:"account_number"`
		Account_title    string `json:"account_title"`
		Account_status   string `json:"account_status"`
		Account_category string `json:"account_category"`
		Cid              string `json:"cid"`
		Insti_code       string `json:"insti_code"`
	}
)
