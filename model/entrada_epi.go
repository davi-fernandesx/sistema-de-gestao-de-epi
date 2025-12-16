package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type EntradaEpiInserir struct {
	ID_epi         int             `json:"id_epi"`
	Id_tamanho     int             `json:"id_tamanho"`
	Data_entrada   time.Time       `json:"data_entrada"`
	Quantidade     int             `json:"quantidade"`
	DataFabricacao time.Time       `json:"data_fabricacao"`
	DataValidade   time.Time       `json:"data_validade"`
	Lote           string          `json:"lote"`
	Fornecedor     string          `json:"fornecedor"`
	ValorUnitario  decimal.Decimal `json:"valorUnitario"`
}

type EntradaEpi struct {
	ID               int
	ID_epi           int
	Nome             string
	Fabricante       string
	CA               string
	Descricao        string
	DataFabricacao   time.Time
	DataValidade     time.Time
	DataValidadeCa   time.Time
	IDprotecao       int
	NomeProtecao     string
	Id_Tamanho       int
	TamanhoDescricao string
	Quantidade       int
	Data_entrada 	 time.Time
	Lote             string
	Fornecedor       string
	ValorUnitario    decimal.Decimal
}

type EntradaEpiDto struct {
	ID            int             `json:"id"`
	Epi           EpiDto          `json:"epi"`
	Tamanho       TamanhoDto      `json:"tamanho"`
	Data_entrada  time.Time       `json:"data_entrada"`
	Quantidade    int             `json:"quantidade"`
	Lote          string          `json:"lote"`
	Fornecedor    string          `json:"fornecedor"`
	ValorUnitario decimal.Decimal `json:"valor_unitario"`
}
