package model

import "time"

type EntregaEpiInserir struct {
	ID                 int
	ID_funcionario     int
	ID_epi             int
	Data_entrega       time.Time
	Assinatura_Digital string
	Quantidade         int
}

type EntregaEpi struct {
	Id                 int
	IDfuncionario      int
	NomeFuncionario    string
	Funcao             string
	Departamento       string
	IDEpi              int
	Nome               string
	Fabricante         string
	CA                 string
	Descricao          string
	DataFabricacao     time.Time
	DataValidade       time.Time
	DataValidadeCa     time.Time
	IDprotecao         int
	NomeProtecao       string
	Data_entrega       time.Time
	Assinatura_Digital string
	Quantidade         int
}

type Entrega_epi_dto struct {
	Funcionario        Funcionario_Dto `json:"funcionario"`
	Epi                EpiDto          `json:"epi"`
	Data_entrega       time.Time       `json:"data_entrega"`
	Assinatura_Digital string          `json:"assinatura_digital"`
	Quantidade         int             `json:"quantidade"`
}
