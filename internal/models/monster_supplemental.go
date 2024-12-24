package models

// Move represents a Pokémon move.
type MonsterSupplemental struct {
	SpecialAttackEV  int         `json:"specialAttackEV"`
	HpEV             int         `json:"hpEV"`
	HatchSteps       int         `json:"hatchSteps"`
	DefenseEV        int         `json:"defenseEV"`
	AttackEV         int         `json:"attackEV"`
	SpecialDefenseEV int         `json:"specialDefenseEV"`
	SpeedEV          int         `json:"speedEV"`
	GenderRatio      interface{} `json:"genderRatio"`
	Species          string      `json:"species"`
	JapaneseName     string      `json:"japaneseName"`
	HepburnName      string      `json:"Guranburu"`
	EggGroups        string      `json:"eggGroups"`

	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}
