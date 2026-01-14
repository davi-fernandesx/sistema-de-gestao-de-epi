package model

import (

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/shopspring/decimal"
)

type ItemParaInserir struct {
	ID_epi           int64             `json:"id_epi" binding:"required"`
	ID_tamanho       int64             `json:"id_tamanho" binding:"required"`
	Quantidade       int             `json:"quantidade" binding:"required,numeric,gt=0"`
	Valor_unitario   decimal.Decimal `json:"valor_unitario" binding:"required,gt=0"`
}

type EntregaParaInserir struct {
	ID_funcionario     int64               `json:"id_funcionario" binding:"required"`
	Data_entrega       configs.DataBr    `json:"data_entrega" binding:"required"`
	Assinatura_Digital string            `json:"assinatura_digital" binding:"required"`
	Itens              []ItemParaInserir `json:"itens" binding:"required,min=1,dive"`
 
}

type ItemEntregueDto struct {
	Id            int64             `json:"id"`
	Epi           EpiDto          `json:"epi"`
	Tamanho       TamanhoDto      `json:"tamanho"`
	Quantidade    int             `json:"quantidade"`
	ValorUnitario decimal.Decimal `json:"valor_unitario"`
}

type EntregaDto struct {
	Id                 int64              `json:"id"`
	Funcionario        Funcionario_Dto   `json:"funcionario"`
	Data_entrega       configs.DataBr    `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemEntregueDto `json:"itens"`
}
