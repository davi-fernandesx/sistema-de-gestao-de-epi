package model

import (
	"time"
)

// model banco de dados (com campos trazidos do inner join)
type Epi struct {
	ID             int
	Nome           string
	Fabricante     string
	CA             string
	Descricao      string
	DataFabricacao time.Time
	DataValidade   time.Time
	DataValidadeCa time.Time
	AlertaMinimo   int
	IDprotecao     int
	NomeProtecao   string
	Tamanhos       []Tamanhos
}

//modelo para ser usado ao inserir/atualizar no banco de dados

type EpiInserir struct {
	ID             int
	Nome           string
	Fabricante     string
	CA             string
	Descricao      string
	DataFabricacao time.Time
	DataValidade   time.Time
	DataValidadeCa time.Time
	Idtamanho      []int
	IDprotecao     int
	AlertaMinimo   int
}

// modelo para ser usado no controller e services
type EpiDto struct {
	Id             int             `json:"id"`
	Nome           string          `json:"nome"`
	Fabricante     string          `json:"fabricante"`
	CA             string          `json:"ca"`
	Tamanho        []TamanhoDto    `json:"tamanhos"`
	Descricao      string          `json:"descricao"`
	DataFabricacao time.Time       `json:"dataFabricante"`
	DataValidade   time.Time       `json:"dataValidade"`
	DataValidadeCa time.Time       `json:"DataValidadeCa"`
	Protecao       TipoProtecaoDto `json:"Protecao"`
}
