package models

type MonsterMove struct {
	Move []struct {
		LearnType string `json:"learn_type"`
		Level     int    `json:"level"`
		Id        int    `json:"id"`
	} `json:"moves"`
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}
