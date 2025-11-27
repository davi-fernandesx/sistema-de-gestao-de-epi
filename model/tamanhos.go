package model

type Tamanhos struct {
	ID      int    `json:"-"`
	Tamanho string `json:"tamanho"`
}

type TamanhoDto struct {
	ID      int    `json:"id"`
	Tamanho string `json:"tamanho"`
}
