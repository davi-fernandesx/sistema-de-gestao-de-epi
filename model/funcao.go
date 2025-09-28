package model

type Funcao struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
}

type FuncaoDto struct {
	Funcao string `json:"cargo"`
}
