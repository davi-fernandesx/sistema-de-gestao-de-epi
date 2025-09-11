package model


type Login struct {

	ID int
	Nome string
	Senha string
}

type LoginDto struct {
	Nome string
	Senha string
}