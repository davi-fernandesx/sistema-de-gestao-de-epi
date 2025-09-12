package login

import (
	"context"
	"errors"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type loginRepository interface{


	AddLogin( model *model.Login)(error) 
	DeletarLogin(ctx context.Context, id int) error
	Login(model *model.Login) (*model.Login, error)
}
//conjunto de erros
var  ErrusuarioJaExistente = errors.New("usuario já cadastrado")
var  ErrAoApagarUmLogin = errors.New("erro ao apagar login")
var ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
