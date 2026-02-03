package model

type Funcao struct {
	Funcao string`json:"funcao"  binding:"required,max=50"`
	IdDepartamento int `json:"id_departamento" binding:"required,min=1"`
}

type FuncaoDto struct {
	ID     int    `json:"id"`
	Funcao string `json:"cargo"`
	Departamento DepartamentoDto `json:"departamento"`
}