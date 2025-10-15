package model

type Departamento struct {
	ID           int    `json:"id"`
	Departamento string `json:"departamento"`
}

type DepartamentoDto struct {
	ID           int    `json:"id"`
	Departamento string `json:"departamento"`
}
