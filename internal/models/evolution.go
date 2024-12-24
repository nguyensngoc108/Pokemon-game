package models

type Up struct {
	NationalId int    `json:"nationalId"`
	Name       string `json:"name"`
	Method     string `json:"method"`
	Level      int    `json:"level"`
}

type Evolution struct {
	From []Up   `json:"from"`
	To   []Up   `json:"to"`
	ID   string `json:"_id"`
	Rev  string `json:"_rev"`
}
