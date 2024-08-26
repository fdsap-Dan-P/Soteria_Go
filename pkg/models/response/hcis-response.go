package response

type StaffInfo struct {
	Cluster            string `json:"cluster"`
	LastName           string `json:"lastName"`
	ZipCode            string `json:"zipCode"`
	FirstName          string `json:"firstname"`
	ProvAddress        string `json:"provAddress"`
	DateHired          string `json:"dateHired"`
	Gender             string `json:"gender"`
	EmploymentStatus   string `json:"employmentStatus"`
	DateRegularization string `json:"dateRegularization"`
	CivilStatus        string `json:"civilStatus"`
	Institution        string `json:"institution"`
	EmailAddress       string `json:"emailAddress"`
	BrgyAddress        string `json:"brgyAddress"`
	CityAddress        string `json:"cityAddress"`
	Department         string `json:"department"`
	StaffID            string `json:"staffID"`
	Area               string `json:"area"`
	EmploymentType     string `json:"employmentType"`
	LocationAssignment string `json:"locationAssignment"`
	NickName           string `json:"nickName"`
	BirthDate          string `json:"birthDate"`
	JobLevel           string `json:"jobLevel"`
	JobPosition        string `json:"jobPosition"`
	Unit               string `json:"unit"`
	JobGrade           string `json:"jobGrade"`
	MobilePhone        string `json:"mobilePhone"`
	MiddleName         string `json:"middleName"`
	Region             string `json:"region"`
	Age                string `json:"age"`
	Status             string `json:"status"`
}

type StaffInfoResponse struct {
	StaffInfo []StaffInfo `json:"StaffInfo"`
}
