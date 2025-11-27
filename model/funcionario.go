package model

type FuncionarioINserir struct {
	Nome            string `json:"nome"`
	Matricula       string `json:"matricula"`
	ID_departamento *int   `json:"id_departamento"`
	ID_funcao       *int   `json:"id_funcao"`
}

type Funcionario struct {
	Id              int
	Nome            string
	Matricula       int
	ID_departamento int
	Departamento    string
	ID_funcao       int
	Funcao          string
}

type Funcionario_Dto struct {
	ID           int             `json:"id"`
	Nome         string          `json:"nome"`
	Matricula    int             `json:"matricula"`
	Departamento DepartamentoDto `json:"departamento"`
	Funcao       FuncaoDto       `json:"funcao"`
}
