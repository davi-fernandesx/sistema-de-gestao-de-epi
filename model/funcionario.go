package model

type FuncionarioINserir struct {
	Nome            string
	Matricula       string
	ID_departamento *int
	ID_funcao       *int
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
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	Matricula    int
	Departamento DepartamentoDto `json:"departamento"`
	Funcao       FuncaoDto       `json:"funcao"`
}
