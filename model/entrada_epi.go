package model

import (
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/shopspring/decimal"
)

type EntradaEpiInserir struct {
	ID_epi           int             `json:"id_epi" binding:"required,numeric"`
	Id_tamanho       int             `json:"id_tamanho" binding:"required,numeric"`
	Data_entrada     configs.DataBr  `json:"data_entrada" binding:"required"`
	Quantidade       int             `json:"quantidade" binding:"required,numeric,gt=0"`
	Quantidade_Atual int             `json:"quantidade_Atual" binding:"required,numeric,gt=0"`
	DataFabricacao   configs.DataBr  `json:"data_fabricacao" binding:"required"`
	DataValidade     configs.DataBr  `json:"data_validade" binding:"required,gtfield=DataFabricacao"`
	Lote             string          `json:"lote" binding:"required,numeric,max=6"`
	Fornecedor       string          `json:"fornecedor" binding:"required,max=50"`
	ValorUnitario    decimal.Decimal `json:"valorUnitario" binding:"required,gte=0"`
}

type EntradaEpi struct {
	ID               int
	ID_epi           int
	Nome             string
	Fabricante       string
	CA               string
	Descricao        string
	DataFabricacao   configs.DataBr
	DataValidade     configs.DataBr
	DataValidadeCa   configs.DataBr
	IDprotecao       int
	NomeProtecao     string
	Id_Tamanho       int
	TamanhoDescricao string
	Quantidade       int
	Quantidade_Atual int
	Data_entrada     configs.DataBr
	Lote             string
	Fornecedor       string
	ValorUnitario    decimal.Decimal
}

type EntradaEpiDto struct {
	ID            int             `json:"id"`
	Epi           EpiDto          `json:"epi"`
	Data_entrada  configs.DataBr  `json:"data_entrada"`
	Quantidade    int             `json:"quantidade"`
	Lote          string          `json:"lote"`
	Fornecedor    string          `json:"fornecedor"`
	ValorUnitario decimal.Decimal `json:"valor_unitario"`
}
