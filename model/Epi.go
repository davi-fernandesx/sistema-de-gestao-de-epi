package model

import (
	
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
)


//modelo para ser usado ao inserir/atualizar no banco de dados

type EpiInserir struct {
	Nome           string`json:"nome" binding:"required"`
	Fabricante     string`json:"fabricante" binding:"required,max=50"`
	CA             string`json:"ca" binding:"required,numeric,min=1,max=6"`
	Descricao      string`json:"descricao" binding:"lte=250"`
	DataValidadeCa configs.DataBr `json:"data_validade_ca" binding:"required"`
	Idtamanho      []int`json:"id_tamanho" binding:"required,min=1"`
	IDprotecao     int`json:"id_protecao" binding:"required,numeric"`
	AlertaMinimo   int`json:"alerta_minimo" binding:"required,gte=0"`
}   
// model banco de dados (com campos trazidos do inner join)
type Epi struct {
	ID             int
	Nome           string
	Fabricante     string
	CA             string
	Descricao      string
	DataValidadeCa configs.DataBr
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
	DataValidadeCa configs.DataBr       `json:"data_validadeCa"`
	Protecao       TipoProtecaoDto `json:"protecao"`
}
