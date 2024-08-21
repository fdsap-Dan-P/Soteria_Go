package response

type (
	InstitutionDetails struct {
		Institution_id   int    `json:"institution_id"`
		Institution_code string `json:"institution_code"`
		Institution_name string `json:"institution_name"`
		Created_at       string `json:"created_at"`
		Updated_at       string `json:"updated_at"`
	}
)
