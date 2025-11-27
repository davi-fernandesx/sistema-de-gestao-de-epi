package model

type Login struct {
	ID    int`json:"-"`
	Nome  string`json:"nome"`
	Senha string`json:"senha"`
}

type LoginDto struct {
	ID    int    `json:"-"`
	Nome  string `json:"nome"`
	Senha string `json:"senha"`
}
