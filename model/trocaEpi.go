package model

import (
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
)

type DevolucaoInserir struct {
	IdFuncionario       int            `json:"id_funcionario" binding:"required"`
	IdEpi               int            `json:"id_epi" binding:"required"`
	IdMotivo            int            `json:"id_motivo" binding:"required"`
	IdTamanho           int            `json:"id_tamanho" binding:"required"`
	DataDevolucao       configs.DataBr `json:"data_devolucao" binding:"required"`
	QuantidadeADevolver int            `json:"quantidade_a_devolver" binding:"required,numeric,gt=0"`
	NovaQuantidade      *int          `json:"nova_quantidade"`
	IdEpiNovo           *int          `json:"id_novo_epi" `
	IdTamanhoNovo       *int          `json:"tamanhoEpi_novo"`
	AssinaturaDigital   string         `json:"assinatura_digital" binding:"required"`
}

type Devolucao struct {
	Id                  int
	DataEntrega         time.Time
	Id_funcionario      int
	NomeFuncionario     string
	Id_departamento     int
	Departamento        string
	Id_funcao           int
	Funcao              string
	ID_epiTroca         int
	NomeEpiTroca        string
	FabricanteTroca     string
	CAtroca             string
	IdTamanho           int
	Tamanho             string
	IdMotivo            int
	Motivo              string
	QuantidadeADevolver int

	AssinaturaDigital string

	ID_epiNovo     *int
	NomeEpiNovo    *string
	FabricanteNovo *string
	CANovo         *string
	Id_tamanhoNovo *int
	TamanhoNovo    *string
	NovaQuantidade *int
}

type DevolucaoDto struct {
	Id                  int             `json:"id"`
	IdFuncionario       Funcionario_Dto `json:"id_funcionario"`
	IdEpi               EpiDto          `json:"id_epi"`
	MotivoDevolucao     DevolucaoEpiDto `json:"motivoDaDevolucao"`
	DataDevolucao       time.Time       `json:"dataDevolucao"`
	QuantidadeADevolver int             `json:"quantidade_a_devolver"`
	AssinaturaDigital   string          `json:"assinatura_digital"`

	IdEpiNovo      *EpiDto     `json:"id_novo_epi"`
	Tamanho        *TamanhoDto `json:"tamanho"`
	NovaQuantidade *int        `json:"quantidade_nova"`
}
