package model

import (
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/shopspring/decimal"
)

type EntradaEpiInserir struct {
	ID_epi             int             `json:"id_epi" binding:"required,numeric"`
	Id_tamanho         int             `json:"id_tamanho" binding:"required,numeric"`
	Data_entrada       configs.DataBr  `json:"data_entrada" binding:"required"`
	Quantidade_Atual   int             `json:"quantidade_Atual" binding:"required,numeric,gt=0"`
	Quantidade         int             `json:"quantidade" binding:"required,numeric,gt=0"`
	DataFabricacao     configs.DataBr  `json:"data_fabricacao" binding:"required"`
	DataValidade       configs.DataBr  `json:"data_validade" binding:"required,gtfield=DataFabricacao"`
	Lote               string          `json:"lote" binding:"required,numeric,max=6"`
	Fornecedor         string          `json:"fornecedor" binding:"required,max=50"`
	Nota_fiscal_serie  string          `json:"notaFicalSerie" binding:"required,max=20,numeric"`
	Nota_fiscal_numero string          `json:"notaFiscalNumero" binding:"required,max=10,numeric"`
	ValorUnitario      decimal.Decimal `json:"valorUnitario" binding:"required,gte=0"`
}

type EntradaEpiDto struct {
	ID                 int             `json:"id"`
	Epi                EpiDto          `json:"epi"`
	Data_entrada       configs.DataBr  `json:"data_entrada"`
	Quantidade         int             `json:"quantidade"`
	Quantidade_Atual   int             `json:"quantidade_Atual"`
	Lote               string          `json:"lote"`
	Fornecedor         string          `json:"fornecedor"`
	Nota_fiscal_serie  string          `json:"notaFicalSerie"`
	Nota_fiscal_numero string          `json:"notaFiscalNumero"`
	ValorUnitario      decimal.Decimal `json:"valor_unitario"`
}
