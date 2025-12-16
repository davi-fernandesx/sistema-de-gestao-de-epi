package model

type FuncaoInserir struct {
	ID     int`json:"id"`
	Funcao string`json:"funcao"`
	IdDepartamento int `json:"id_departamento"`
}

type Funcao struct {
	ID     int`json:"id"`
	Funcao string`json:"funcao"`
	IdDepartamento int `json:"id_departamento"`
	NomeDepartamento string `json:"departamento"`
}

type FuncaoDto struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
	Departamento DepartamentoDto `json:"departamento"`
}
