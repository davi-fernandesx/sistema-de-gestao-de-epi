package model

type Tamanhos struct {
	ID      int   `json:"id"`
	Tamanho int16 `json:"tamanho"`
}

type TamanhoDto struct {
	Tamanho int16 `json:"tamanho"`
}
