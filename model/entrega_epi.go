package model

import "time"

type Entrega_epi struct {
	ID                 int       `json:"id"`
	ID_funcionario     int       `json:"id_funcionario"`
	ID_epi             int       `json:"id_epi"`
	Data_entrega       time.Time `json:"data_entrega"`
	Assinatura_Digital string    `json:"assinatura_digital"`
	Quantidade         int       `json:"quantidade"`
}

type Entrega_epi_dto struct {
	Funcionario        Funcionario_Dto `json:"funcionario"`
	Epi                EpiDto        `json:"epi"`
	Data_entrega       time.Time       `json:"data_entrega"`
	Assinatura_Digital string          `json:"assinatura_digital"`
	Quantidade         int             `json:"quantidade"`
}
