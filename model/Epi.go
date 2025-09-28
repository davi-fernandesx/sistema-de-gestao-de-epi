package model

import "time"

// model banco de dados
type Epi struct {
	ID             int       `json:"id"`
	Nome           string    `json:"nome"`
	Fabricante     string    `json:"fabricante"`
	CA             string    `json:"ca"`
	Descricao      string    `json:"descricao"`
	DataFabricacao time.Time `json:"dataFabricante"`
	DataValidade   time.Time `json:"dataValidade"`
	DataValidadeCa time.Time `json:"DataValidadadeCa"`
	IDprotecao     int       `json:"idProtecao"`
	AlertaMinimo   int       `json:"alertaMinimo"`
}
