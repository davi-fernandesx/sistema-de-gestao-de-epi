package model

type Funcao struct {
	ID     int    
	Funcao string 
}

type FuncaoDto struct {
	Funcao string `json:"cargo"`
}
