package response

type (
	MemberResponse struct {
		RetCode string     `json:"retCode"`
		Message string     `json:"message"`
		Data    MemberData `json:"data"`
	}

	MemberData struct {
		Message   string        `json:"message"`
		IsSuccess bool          `json:"isSuccess"`
		Error     error         `json:"error,omitempty"`
		Details   MemberDetails `json:"details,omitempty"`
	}

	MemberDetails struct {
		Cid           string `json:"cid"`
		First_name    string `json:"first_name"`
		Last_name     string `json:"last_name"`
		Date_of_birth string `json:"date_of_birth"`
		Phone_1       string `json:"phone_1"`
		Phone_2       string `json:"phone_2"`
		Center_code   string `json:"center_code"`
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
