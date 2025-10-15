package model

type Login struct {
	ID    int
	Nome  string
	Senha string
}

type LoginDto struct {
	ID    int    `json:"id"`
	Nome  string `json:""`
	Senha string `json:""`
}
