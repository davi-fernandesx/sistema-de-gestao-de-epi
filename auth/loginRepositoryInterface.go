package auth

import (
	"context"
	"errors"
)

type loginRepository interface{


	AddLogin(ctx context.Context, model *Login)(*Login, error) 
	DeletarLogin(ctx context.Context, model Login) error
	Login(ctx context.Context, model Login) bool
}

//conjunto de erros
var  UsuarioJaExistente = errors.New("usuario já cadastrado")
