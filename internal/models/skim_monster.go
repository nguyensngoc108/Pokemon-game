package models

// Move represents a Pokémon move.
type ListMapObject struct {
	Name        string `json:"name"`
	ResourceURI string `json:"resource_uri"`
}

type SkimMonster struct {
	Description []ListMapObject `json:"descriptions"`
	Type        []ListMapObject `json:"types"`
	Abilities   []ListMapObject `json:"abilities"`
	/*
		"attack":49,"defense":49,"speed":45,"sp_atk":65,"sp_def":65,"hp":45,"weight":"69","height":"7","national_id":1,"name":"Bulbasaur","male_female_ratio":"87.5/12.5","abilities":[{"name":"chlorophyll","resource_uri":"/api/v1/ability/34/"
	*/
	Attack          int    `json:"attack"`
	Defense         int    `json:"defense"`
	Speed           int    `json:"speed"`
	SpAtk           int    `json:"sp_atk"`
	SpDef           int    `json:"sp_def"`
	HP              int    `json:"hp"`
	Weight          string `json:"weight"`
	Height          string `json:"height"`
	NationalID      int    `json:"national_id"`
	MaleFemaleRatio string `json:"male_female_ratio"`
	CatchRate       int    `json:"catch_rate"`
	ID              string `json:"_id"`
	Name            string `json:"name"`
}

func (m SkimMonster) ToMonster() *Monster {
	var abilities []string
	for _, ability := range m.Abilities {
		abilities = append(abilities, ability.Name)
	}

	var types []string
	for _, t := range m.Type {
		types = append(types, t.Name)
	}
	return &Monster{
		Attack:          m.Attack,
		Defense:         m.Defense,
		Speed:           m.Speed,
		SpAtk:           m.SpAtk,
		SpDef:           m.SpDef,
		HP:              m.HP,
		Weight:          m.Weight,
		Height:          m.Height,
		NationalID:      m.NationalID,
		MaleFemaleRatio: m.MaleFemaleRatio,
		CatchRate:       m.CatchRate,
		ID:              m.ID,
		Name:            m.Name,
		Abilities:       abilities,
		Types:           types,
	}
}

type Monster struct {
	Name            string   `json:"name"`
	Attack          int      `json:"attack"`
	Defense         int      `json:"defense"`
	Speed           int      `json:"speed"`
	SpAtk           int      `json:"sp_atk"`
	SpDef           int      `json:"sp_def"`
	HP              int      `json:"hp"`
	Weight          string   `json:"weight"`
	Height          string   `json:"height"`
	NationalID      int      `json:"national_id"`
	MaleFemaleRatio string   `json:"male_female_ratio"`
	CatchRate       int      `json:"catch_rate"`
	ID              string   `json:"_id"`
	Abilities       []string `json:"abilities"`
	Types           []string `json:"types"`
}
