package model

import "time"

type DevolucaoInserir struct {
	IdFuncionario     int       `json:"id_funcionario"`
	IdEpi             int       `json:"id_epi"`
	IdMotivo          int       `json:"id_motivo"`
	DataDevolucao     time.Time `json:"data_devolucao"`
	Quantidade        int       `json:"quantidade"`
	IdEpiNovo         int       `json:"id_novo_epi"`
	IdTamanhoNovo     int       `json:"tamanhoEpi_novo"`
	AssinaturaDigital string    `json:"assinatura_digital"`
}

type Devolucao struct {
	Id                int
	DataEntrega       time.Time
	Id_funcionario    int
	NomeFuncionario   string
	Id_departamento   int
	Departamento      string
	Id_funcao         int
	Funcao            string
	ID_epiTroca       int
	NomeEpiTroca      string
	FabricanteTroca   string
	CAtroca           string
	IdMotivo          int
	Motivo            string
	AssinaturaDigital string

	ID_epiNovo int
	NomeEpi    string
	Fabricante string
	CA         string
	Id_tamanho int
	Tamanho    string
}

type DevolucaoDto struct {
	Id                int             `json:"id"`
	IdFuncionario     Funcionario_Dto `json:"id_funcionario"`
	IdEpi             EpiDto          `json:"id_epi"`
	MotivoDevolucao   DevolucaoEpiDto `json:"motivoDaDevolucao"`
	DataDevolucao     time.Time       `json:"dataDevolucao"`
	AssinaturaDigital string          `json:"assinatura_digital"`

	IdEpiNovo  EpiDto     `json:"id_novo_epi"`
	Tamanho    TamanhoDto `json:"tamanho"`
	Quantidade int        `json:"quantidade"`
}
