package model

import (
	"time"
)

type ItemParaEntrega struct {
	Id         int
	ID_epi     int
	ID_tamanho int
	Quantidade int
	Id_entrega int
}

type EntregaParaInserir struct {
	Id                 int
	ID_funcionario     int               `json:"id_funcionario"`
	Data_entrega       time.Time         `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemParaEntrega `json:"itens"`
}

type Entrega struct {
	Id              int
	DataEntrega     time.Time
	Id_funcionario  int
	NomeFuncionario string
	Id_departamento int
	Departamento    string
	Id_funcao       int
	Funcao          string
	ID_epi          int
	NomeEpi         string
	Fabricante      string
	CA              string
	Descricao       string
	DataFabricacao  time.Time
	DataValidade    time.Time
	DataValidadeCa  time.Time
	IDprotecao      int
	NomeProtecao    string
	Tamanhos        string
	Quantidade      int
}

type ItemEntregueDto struct {
	Id         int        `json:"id"`
	Epi        EpiDto     `json:"epi"`
	Tamanho    TamanhoDto `json:"tamanho"`
	Quantidade int        `json:"quantidade"`
}

type EntregaDto struct {
	Id                 int               `json:"id"`
	Funcionario        Funcionario_Dto   `json:"funcionario"`
	Data_entrega       time.Time         `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemEntregueDto `json:"itens"`
}
