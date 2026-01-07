package model

type FuncaoInserir struct {
	Funcao string`json:"funcao"  binding:"required,min=2,max=50"`
	IdDepartamento int `json:"id_departamento" binding:"required,min=1"`
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
