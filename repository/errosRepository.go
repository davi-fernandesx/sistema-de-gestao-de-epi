package repository

import (

	"errors"
)

//model do repository

// erros comuns
var ErrConexaoDb = errors.New("erro ao se conectar com o banco de dados")
var ErrFalhaAoEscanearDados = errors.New("erro ao escanear os dados")
var ErrLinhasAfetadas = errors.New( "erro ao verificar as linhas afetadas")
var ErrDadoIncompativel = errors.New("tipo de dado invalido")

//conjunto de erros usuario
var  ErrusuarioJaExistente = errors.New("usuario já cadastrado")
var  ErrAoApagarUmLogin = errors.New("erro ao apagar login")
var ErrUsuarioNaoEncontrado = errors.New("usuario não encontrado")


// erros departamentos

var  ErrDepartamentoJaExistente = errors.New("departamento já cadastrado")
var  ErrAoApagarUmDepartamento = errors.New("erro ao apagar departamento")
var ErrDepartamentoNaoEncontrado = errors.New("departamento não encontrado")
var ErrBuscarTodosDepartamentos = errors.New("erro ao buscar todos os departamentos")
var ErrIterarSobreDepartamentos = errors.New("erro ao iterar sobre os departamentos")


//epi

var ErrEpiAoAdicionarEpi = errors.New("erro ao executar o comendo db exec ao adicionar um epi no banco de dados")
var ErrAoProcurarEpi = errors.New("epi não encontrado")
var ErrAoBuscarTodosOsEpis = errors.New("erro ao buscar todos os epis")
var ErrAoInterarSobreEpis = errors.New("erro ao iterar sobre os epis")
var ErrEpiNaoEncontrado = errors.New("epi nao encontrado na base de dados")

// erros funcao

var ErrFuncaoJaExistente = errors.New("funcao ja cadastrada") 
