package model

import (
	"time"
)


//modelo para ser usado ao inserir/atualizar no banco de dados

type EpiInserir struct {
	Nome           string`json:"nome"`
	Fabricante     string`json:"fabricante"`
	CA             string`json:"ca"`
	Descricao      string`json:"descricao"`
	DataFabricacao time.Time`json:"data_fabricacao"`
	DataValidade   time.Time`json:"data_validade"`
	DataValidadeCa time.Time`json:"data_validade_ca"`
	Idtamanho      []int`json:"id_tamanho"`
	IDprotecao     int`json:"id_protecao"`
	AlertaMinimo   int`json:"alerta_minimo"`
}
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



// modelo para ser usado no controller e services
type EpiDto struct {
	Id             int             `json:"id"`
	Nome           string          `json:"nome"`
	Fabricante     string          `json:"fabricante"`
	CA             string          `json:"ca"`
	Tamanho        []TamanhoDto    `json:"tamanhos"`
	Descricao      string          `json:"descricao"`
	DataFabricacao time.Time       `json:"data_fabricante"`
	DataValidade   time.Time       `json:"data_validade"`
	DataValidadeCa time.Time       `json:"data_validadeCa"`
	Protecao       TipoProtecaoDto `json:"protecao"`
}
