package response

type ResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DataModel struct {
	Message   string      `json:"message"`
	IsSuccess bool        `json:"isSuccess"`
	Error     error       `json:"error,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

type ReturnModel struct {
	RetCode string    `json:"retCode"`
	Message string    `json:"message"`
	Data    DataModel `json:"data,omitempty"`
}

type DBFuncResponse struct {
	Remark string `json:"remark"`
}

type RespFromDB struct {
	RetCode       string `json:"retCode"`
	Category      string `json:"category"`
	Error_message string `json:"error_message"`
	Is_success    bool   `json:"is_success"`
}

type Total struct {
	Count int `json:"count"`
}
