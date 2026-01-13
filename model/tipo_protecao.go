package model

type TipoProtecao struct {
	ID   int      `json:"-"`
	Nome Protecao `json:"nome" binding:"required, min=6, max=50"`
}

type TipoProtecaoDto struct {
	ID   int64      `json:"id"`
	Nome Protecao `json:"nome"`
}

type Protecao string

const (
	Proteção_para_os_Pés_e_Pernas Protecao = " Proteção para os Pés e Pernas "
	Proteção_das_Mãos_e_Braços    Protecao = " Proteção das Mãos e Braços"
	Proteção_do_Corpo             Protecao = " Proteção do Corpo "
	Proteção_da_Cabeça_e_Face     Protecao = "Proteção da Cabeça e Face"
)
