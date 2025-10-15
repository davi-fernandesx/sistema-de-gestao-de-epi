package model

type Funcao struct {
	ID     int
	Funcao string
}

type FuncaoDto struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
}
