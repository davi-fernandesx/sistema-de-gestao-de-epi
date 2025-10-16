package model

type FuncionarioINserir struct {
	ID              int   
	Nome            string 
	ID_departamento int    
	ID_funcao       int    
}

type Funcionario struct {

	Id int
	Nome string
	ID_departamento int
	Departamento string
	ID_funcao int
	Funcao string
}


type Funcionario_Dto struct {
	ID           int    `json:"id"`
	Nome         string          `json:"nome"`
	Departamento DepartamentoDto `json:"departamento"`
	Funcao       FuncaoDto       `json:"funcao"`
}