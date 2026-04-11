package model

type MembersDB struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	JoiningDate string `json:"joining_date"`
	Duration    int64  `json:"duration"`
	ExpiryDate  string `json:"expiry_date"`
}

type Members struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	JoiningDate string `json:"joining_date"`
	Duration    int64  `json:"duration"`
}
