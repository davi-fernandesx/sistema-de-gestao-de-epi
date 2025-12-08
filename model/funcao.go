package model

type Funcao struct {
	ID     int`json:"-"`
	Funcao string`json:"funcao"`
	IdDepartamento int `json:"id_departamento"`
}

type FuncaoDto struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
	Departamento DepartamentoDto `json:"departamento"`
}
