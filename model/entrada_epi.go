package model

import "time"

type Entrada_epi struct {
	ID           int       `json:"id"`
	ID_epi       int       `json:"id_epi"`
	Data_entrada time.Time `json:"data_entrada"`
	Quantidade   int       `json:"quantidade"`
	Lote         string    `json:"lote"`
	Fornecedor   string    `json:"fornecedor"`
}

type Entrada_epi_dto struct {
	Epi          Epi_dto   `json:"epi"`
	Data_entrada time.Time `json:"data_entrada"`
	Quantidade   int       `json:"quantidade"`
	Lote         string    `json:"lote"`
	Fornecedor   string    `json:"fornecedor"`
}
