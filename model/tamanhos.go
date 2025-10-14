package model

type Tamanhos struct {
	ID      int   
	Tamanho string 
}

type TamanhoDto struct {
	Tamanho int16 `json:"tamanho"`
}
