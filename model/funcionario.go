package model

type Funcionario struct {
	ID              int    `json:"id"`
	Nome            string `json:"nome"`
	ID_departamento int    `json:"id_departamento"`
	ID_funcao       int    `json:"id_funcao"`
}



type Funcionario_Dto struct {
	Nome         string          `json:"nome"`
	Departamento DepartamentoDto `json:"departamento"`
	Funcao       FuncaoDto       `json:"funcao"`
}
