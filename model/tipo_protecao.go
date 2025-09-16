package model

type TipoProtecao struct {

	ID int  `json:"id"`
	Nome Protecao `json:"nome"`

}


type TipoProtecaoDto struct {

	Nome Protecao `json:"nome"`
}

type Protecao string

const (

	Proteção_para_os_Pés_e_Pernas Protecao = " Proteção para os Pés e Pernas "
	Proteção_das_Mãos_e_Braços Protecao = " Proteção das Mãos e Braços"
	Proteção_do_Corpo  Protecao = " Proteção do Corpo "
	Proteção_da_Cabeça_e_Face Protecao = "Proteção da Cabeça e Face"

)


