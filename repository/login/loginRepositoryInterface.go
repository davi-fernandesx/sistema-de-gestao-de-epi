package auth

import (
	"context"
	"errors"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type loginRepository interface{


	AddLogin(ctx context.Context, model *model.Login)(*model.Login, error) 
	DeletarLogin(ctx context.Context, id int) error
	Login(ctx context.Context, Nome string) (*model.Login, error)
}
//conjunto de erros
var  ErrusuarioJaExistente = errors.New("usuario já cadastrado")
var  ErrAoApagarUmLogin = errors.New("erro ao apagar login")
var ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
