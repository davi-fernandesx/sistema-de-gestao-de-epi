package auth

import (
	"golang.org/x/crypto/bcrypt"
	
)

func Hash(senha string) ([]byte, error) {

	hashPass, err:=bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return  nil, err
	} 

	return  hashPass, nil

}
