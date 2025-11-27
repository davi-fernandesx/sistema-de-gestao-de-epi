package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"time"
)

type ItemParaInserir struct {
	ID_epi         int             `json:"id_epi"`
	ID_tamanho     int             `json:"id_tamanho"`
	Quantidade     int             `json:"quantidade"`
	IdEntrada      int             `json:"-"`
	IdEntrega      int             `json:"-"`
	Valor_unitario decimal.Decimal `json:"valor_unitario"`
}

type EntregaParaInserir struct {
	ID_funcionario     int               `json:"id_funcionario"`
	Data_entrega       time.Time         `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemParaInserir `json:"itens"`
	Id_troca           int               `json:""`
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
	Valor_unitario  decimal.Decimal
	Id_troca        sql.NullInt64
}

type ItemEntregueDto struct {
	Id            int             `json:"id"`
	Epi           EpiDto          `json:"epi"`
	Tamanho       TamanhoDto      `json:"tamanho"`
	Quantidade    int             `json:"quantidade"`
	ValorUnitario decimal.Decimal `json:"valor_unitario"`
}

type EntregaDto struct {
	Id                 int               `json:"id"`
	Funcionario        Funcionario_Dto   `json:"funcionario"`
	Data_entrega       time.Time         `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemEntregueDto `json:"itens"`
}
