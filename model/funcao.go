package model

type Funcao struct {
	ID     int`json:"-"`
	Funcao string`json:"funcao"`
}

type FuncaoDto struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
}
