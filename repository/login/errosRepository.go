package login

import (

	"errors"

)


//conjunto de erros
var  ErrusuarioJaExistente = errors.New("usuario já cadastrado")
var  ErrAoApagarUmLogin = errors.New("erro ao apagar login")
var ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
