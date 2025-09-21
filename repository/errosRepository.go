package repository

import (

	"errors"
)

//model do repository


//conjunto de erros
var  ErrusuarioJaExistente = errors.New("usuario já cadastrado")
var  ErrAoApagarUmLogin = errors.New("erro ao apagar login")
var ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
var ErrUsuarioNaoEncontrado = errors.New("usuario não encontrado")
var ErrConexaoDb = errors.New("erro ao se conectar com o banco de dados")
var ErrFalhaAoEscanearDados = errors.New("erro ao escanear os dados")


