package model

type TipoProtecao struct {
	Nome string `json:"nome" binding:"required,max=50"`
}

type TipoProtecaoDto struct {
	ID   int64  `json:"id"`
	Nome string `json:"nome"`
}
