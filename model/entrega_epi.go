package model

import (
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/shopspring/decimal"
)

type ItemParaInserir struct {
	ID_epi         int             `json:"id_epi" binding:"requered,numeric"`
	ID_tamanho     int             `json:"id_tamanho" binding:"required,numeric"`
	Quantidade     int             `json:"quantidade" binding:"required,numeric,gt=0"`
	Valor_unitario decimal.Decimal `json:"valor_unitario" binding:"required,gt=0"`
}

type EntregaParaInserir struct {
	ID_funcionario     int               `json:"id_funcionario" binding:"requered,numeric"`
	Data_entrega       configs.DataBr    `json:"data_entrega" binding:"requered"`
	Assinatura_Digital string            `json:"assinatura_digital" binding:"requered"`
	Itens              []ItemParaInserir `json:"itens" binding:"requered,min=1,dive"`
	Id_troca           sql.NullInt64     `json:"-"`
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
	Data_entrega       configs.DataBr    `json:"data_entrega"`
	Assinatura_Digital string            `json:"assinatura_digital"`
	Itens              []ItemEntregueDto `json:"itens"`
}
