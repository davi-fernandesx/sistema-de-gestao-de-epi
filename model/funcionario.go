package model

type Funcionario struct {
	ID              int   
	Nome            string 
	ID_departamento int    
	ID_funcao       int    
}



type Funcionario_Dto struct {
	Nome         string          `json:"nome"`
	Departamento DepartamentoDto `json:"departamento"`
	Funcao       FuncaoDto       `json:"funcao"`
}
