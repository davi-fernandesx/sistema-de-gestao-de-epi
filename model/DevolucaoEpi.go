package model

type DevolucaoEpi struct {
	Id     int             `jsion:"-"`
	Motivo MotivoDevolucao `json:"motivo"`
}

type DevolucaoEpiDto struct {
	Id int `json:"id"`
	Motivo MotivoDevolucao `json:"motivo"`
}

type MotivoDevolucao string

const (
	Numeracao_ou_tamanho_errado                  MotivoDevolucao = "Numeração ou tamanho errado"
	Substituição_por_Desgaste_ou_Dano MotivoDevolucao = "Substituição por Desgaste ou Dano"
	Data_de_validade_Vencida          MotivoDevolucao = "Vencimento da validade ou do CA"
	Mudança_de_Função_ou_Setor        MotivoDevolucao = "Mudança de Função ou Setor"
	Demissao                          MotivoDevolucao = "Demissão"
)
