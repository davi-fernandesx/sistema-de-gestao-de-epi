package model


type Departamento struct {
	Departamento string `json:"departamento" binding:"required,max=50"`
}

type DepartamentoDto struct {
	ID           int    `json:"id"`
	Departamento string `json:"departamento" example:"Recursos Humanos"`
}