package model

type FuncionarioINserir struct {
	Nome            string `json:"nome" binding:"required,min=3,max=150"`
	Matricula       string `json:"matricula" binding:"required,min=7,max=7"`
	ID_departamento *int   `json:"id_departamento" binding:"required,min=1"`
	ID_funcao       *int   `json:"id_funcao"  binding:"required,min=1"`
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
	Funcao       FuncaoDto       `json:"funcao"`
}
