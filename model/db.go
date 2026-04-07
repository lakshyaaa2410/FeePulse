package model

type MemebersDB struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	JoiningDate string `json:"joining_date"`
	Duration    int64  `json:"duration"`
	ExpiryDate  int64  `json:"expiry_date"`
}
