package model

type Login struct {
	ID    int    `json:"-"`
	Nome  string `json:""`
	Senha string `json:""`
}

type LoginDto struct {
	Nome  string `json:""`
	Senha string `json:""`
}
