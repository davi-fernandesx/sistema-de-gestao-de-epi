package model

import "time"

type EntradaEpiInserir struct {
	ID           int
	ID_epi       int
	Id_tamanho   int
	Data_entrada time.Time
	Quantidade   int
	Lote         string
	Fornecedor   string
}

type EntradaEpi struct {
	ID               int
	ID_epi           int
	Nome             string
	Fabricante       string
	CA               string
	Descricao        string
	DataFabricacao   time.Time
	DataValidade     time.Time
	DataValidadeCa   time.Time
	IDprotecao       int
	NomeProtecao     string
	Id_Tamanho       int
	TamanhoDescricao string
	Quantidade       int
	Lote             string
	Fornecedor       string
}

type EntradaEpiDto struct {
	Epi          EpiDto     `json:"epi"`
	Tamanho      TamanhoDto `json:"tamanho"`
	Data_entrada time.Time  `json:"data_entrada"`
	Quantidade   int        `json:"quantidade"`
	Lote         string     `json:"lote"`
	Fornecedor   string     `json:"fornecedor"`
}
