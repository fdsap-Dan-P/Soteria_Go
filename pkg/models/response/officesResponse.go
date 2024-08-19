package response

type InstitutionDetails struct {
	Institution_id          int    `json:"institution_id"`
	Institution_code        string `json:"institution_code"`
	Institution_name        string `json:"institution_name"`
	Institution_description string `json:"institution_description"`
	Created_at              string `json:"created_at"`
	Updated_at              string `json:"updated_at"`
}

type AreaDetails struct {
	Area_id          int    `json:"area_id"`
	Area_code        string `json:"area_code"`
	Area_name        string `json:"area_name"`
	Area_description string `json:"area_description"`
	Institution_id   int    `json:"institution_id"`
	Created_at       string `json:"created_at"`
	Updated_at       string `json:"updated_at"`
}

type UnitDetails struct {
	Unit_id          int    `json:"unit_id"`
	Unit_code        string `json:"unit_code"`
	Unit_name        string `json:"unit_name"`
	Unit_description string `json:"unit_description"`
	Area_id          int    `json:"area_id"`
	Created_at       string `json:"created_at"`
	Updated_at       string `json:"updated_at"`
}

type CenterDetails struct {
	Center_id          int    `json:"center_id"`
	Center_code        string `json:"center_code"`
	Center_name        string `json:"center_name"`
	Center_description string `json:"center_description"`
	Unit_id            int    `json:"unit_id"`
	Created_at         string `json:"created_at"`
	Updated_at         string `json:"updated_at"`
}
