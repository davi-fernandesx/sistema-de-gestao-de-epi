package model

import "time"


//model banco de dados 
type Epi struct {

	ID int `json:"id"`
	Nome string `json:"nome"`
	Fabricante string `json:"fabricante"`
	CA string `json:"ca"`
	Descricao string `json:"descricao"`
	DataFabricante time.Timer `json:"dataFabricante"`
	DataValidade time.Timer `json:"dataValidade"`
	DataValidadeCa time.Timer `json:"DataValidadadeCa"`
	IDprotecao int `json:"idProtecao"`
	AlertaMinimo int `json:"alertaMinimo"`
	
}

