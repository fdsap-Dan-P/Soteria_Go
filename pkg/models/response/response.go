package response

import "soteria_go/pkg/models/errors"

type ResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ReturnModel struct {
	RetCode string            `json:"retCode"`
	Message string            `json:"message"`
	Data    errors.ErrorModel `json:"data,omitempty"`
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
