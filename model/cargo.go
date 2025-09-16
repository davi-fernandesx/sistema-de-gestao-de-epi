package model


type Cargo struct {

	ID 	int `json:" id"`
	Cargo string `json:"cargo"`
}

type CargoDto struct {

	Cargo string `json:"cargo"`
}


