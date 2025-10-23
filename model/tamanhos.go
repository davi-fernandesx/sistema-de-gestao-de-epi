package model

type Tamanhos struct {
	ID      int   
	Tamanho string 
}

type TamanhoDto struct {
	ID      int `json:"id"`
	Tamanho string `json:"tamanho"`
}
