package model

type Departamento struct {
	ID           int    `json:"-"`
	Departamento string `json:"departamento" binding:"required,min=2,max=50"`
}

type DepartamentoDto struct {
	ID           int    `json:"id"`
	Departamento string `json:"departamento"`
}
