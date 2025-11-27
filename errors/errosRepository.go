package Errors

import (

	"errors"
)

var (

	ErrConexaoDb = errors.New("erro ao se conectar com o banco de dados")
	ErrFalhaAoEscanearDados = errors.New("erro ao escanear os dados")
	ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
	ErrDadoIncompativel = errors.New("tipo de dado invalido")
	ErrAoapagar = errors.New("erro ao apagar")
	ErrAoIterar = errors.New("erro ao iterar")
	ErrSalvar = errors.New("erro ao salvar")
	ErrBuscarTodos = errors.New("erro ao buscar todo os itens")
    ErrNaoEncontrado = errors.New(" não encontrado")
	ErrInternal      = errors.New("erro interno do repositório")
	ErrEstoqueInsuficiente = errors.New("estoque do epi zerado")

)
