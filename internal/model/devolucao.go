package model

import "github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"

type DevolucaoInserir struct {
	IdFuncionario       int            `json:"id_funcionario" binding:"required"`
	IdEpi               int            `json:"id_epi" binding:"required"`
	IdMotivo            int            `json:"id_motivo" binding:"required"`
	IdTamanho           int            `json:"id_tamanho" binding:"required"`
	DataDevolucao       configs.DataBr `json:"data_devolucao" binding:"required"`
	QuantidadeADevolver int            `json:"quantidade_a_devolver" binding:"required,numeric,gt=0"`
	NovaQuantidade      *int           `json:"nova_quantidade"`
	IdEpiNovo           *int           `json:"id_novo_epi" `
	IdTamanhoNovo       *int           `json:"tamanhoEpi_novo"`
	Troca               bool           `json:"É_troca" binding:"required"`
	AssinaturaDigital   string         `json:"assinatura_digital" binding:"required"`
	IdUser              int            `json:"usuario" binding:"required"`
}

type DevolucaoDto struct {
	Id                  int             `json:"id"`
	IdFuncionario       Funcionario_Dto `json:"id_funcionario"`
	IdEpi               EpiDto          `json:"id_epi"`
	MotivoDevolucao     string          `json:"motivoDaDevolucao"`
	DataDevolucao       configs.DataBr  `json:"dataDevolucao"`
	QuantidadeADevolver int             `json:"quantidade_a_devolver"`
	AssinaturaDigital   string          `json:"assinatura_digital"`

	IdEpiNovo      *EpiDto     `json:"id_novo_epi"`
	Tamanho        *TamanhoDto `json:"tamanho"`
	NovaQuantidade *int        `json:"quantidade_nova"`
}
