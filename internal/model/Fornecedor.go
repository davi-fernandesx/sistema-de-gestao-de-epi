package model



type FornecedorInserir struct {
	RazaoSocial       string         `json:"razao_social" binding:"required,max=100"`
	NomeFantasia      string         `json:"nome_fantasia" binding:"required,max=100"`
	CNPJ              string         `json:"cnpj" binding:"required,cnpj"` // Valide 14 digitos no validator
	InscricaoEstadual string         `json:"inscricao_estadual" binding:"required"`
}

type Fornecedor struct {
    ID                int       `json:"id"`
    TenantID          int       `json:"-"` // Não precisa retornar pro front
    RazaoSocial       string    `json:"razao_social"`
    NomeFantasia      string    `json:"nome_fantasia"`
    CNPJ              string    `json:"cnpj"` // Valide 14 digitos no validator
    InscricaoEstadual string    `json:"inscricao_estadual"`
    Ativo             bool      `json:"ativo"`
    
}

type FornecedorDto struct {
    ID                int       `json:"id"`// Não precisa retornar pro front
    RazaoSocial       string    `json:"razao_social"`
    NomeFantasia      string    `json:"nome_fantasia"`
    CNPJ              string    `json:"cnpj"` // Valide 14 digitos no validator
    InscricaoEstadual string    `json:"inscricao_estadual"` 
}

type FornecedorUpdate struct {
	RazaoSocial       *string         `json:"razao_social"`
	NomeFantasia      *string         `json:"nome_fantasia"`
	CNPJ              *string         `json:"cnpj" binding:"cnpj"` // Valide 14 digitos no validator
	InscricaoEstadual *string         `json:"inscricao_estadual"`
}