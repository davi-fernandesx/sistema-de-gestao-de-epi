package auth

import (
	"golang.org/x/crypto/bcrypt"
	
)

//funcao criada para codificar a senha de login do usuario

func HashPassword(senha string) ([]byte, error) {

	hashPass, err:=bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return  nil, err
	} 

	return  hashPass, nil

}

//funcao que ira servir para comparar a senha que o usuario digitar, com a que esta salva no banco de dados
func HashCompare(hash []byte, senha []byte)(bool, error){

	err:= bcrypt.CompareHashAndPassword(hash, senha)
	if err != nil {
		return false , err
	}

	return  true, nil
}
