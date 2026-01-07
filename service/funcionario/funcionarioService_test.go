package funcionario

import (
	"context"
	"database/sql"
	"testing"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFuncionarioRepo struct {
	mock.Mock
}

func (m *MockFuncionarioRepo) AddFuncionario(ctx context.Context, Funcionario *model.FuncionarioINserir) error {

	args := m.Called(ctx, Funcionario)

	return args.Error(0)
}

func (m *MockFuncionarioRepo) BuscaFuncionario(ctx context.Context, matricula string) (*model.Funcionario, error) {

	args := m.Called(ctx, matricula) //passando os argumentos

	var f *model.Funcionario
	if args.Get(0) != nil {
		f = args.Get(0).(*model.Funcionario)
	}

	return f, args.Error(1)
}

func (m *MockFuncionarioRepo) BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error) {

	args := m.Called(ctx)

	var funcs []model.Funcionario

	if args.Get(0) != nil {

		funcs = args.Get(0).([]model.Funcionario)
	}
	return funcs, args.Error(1)
}

func (m *MockFuncionarioRepo) DeletarFuncionario(ctx context.Context, matricula string) error {
	args := m.Called(ctx, matricula)

	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioNome(ctx context.Context, id int, funcionario string) error {

	args := m.Called(ctx, id, funcionario)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioDepartamento(ctx context.Context, id int, idDepartamento string) error {

	args := m.Called(ctx, id, idDepartamento)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioFuncao(ctx context.Context, id int, idFuncao string) error {

	args := m.Called(ctx, id, idFuncao)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}

func TestSalvarFuncionario(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Erro Matricula ja cadastrada", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}
		
		input := model.FuncionarioINserir{Matricula: "123", Nome: "Jose"}
		
		// Simula que a busca encontrou um funcionário (Matrícula já existe)
		m.On("BuscaFuncionario", ctx, "123").Return(&model.Funcionario{Id: 1}, nil)

		err := svc.SalvarFuncionario(ctx, input)
		
		assert.Error(t, err)
		assert.Equal(t, ErrMatricula, err)
	})

	t.Run("Sucesso ao salvar novo funcionario", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}
		
		input := model.FuncionarioINserir{Matricula: "124", Nome: "Maria"}
		
		// 1. Busca não encontra ninguém (sql.ErrNoRows)
		m.On("BuscaFuncionario", ctx, "124").Return(nil, sql.ErrNoRows)
		// 2. Chama o Add
		m.On("AddFuncionario", ctx, mock.AnythingOfType("*model.FuncionarioINserir")).Return(nil)

		err := svc.SalvarFuncionario(ctx, input)
		
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})
}

func TestListarFuncionarioPorMatricula(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao buscar", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}
		
		funcFake := &model.Funcionario{
			Id: 1, Nome: "Alan", Matricula: "500", 
			ID_departamento: 1, Departamento: "TI",
			ID_funcao: 1, Funcao: "Dev",
		}

		m.On("BuscaFuncionario", ctx, "500").Return(funcFake, nil)

		res, err := svc.ListarFuncionarioPorMatricula(ctx, "500")
		
		assert.NoError(t, err)
		assert.Equal(t, "Alan", res.Nome)
		assert.Equal(t, "TI", res.Funcao.Departamento.Departamento)
	})

	t.Run("Funcionario nao encontrado no banco", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}

		m.On("BuscaFuncionario", ctx, "999").Return(nil, sql.ErrNoRows)

		res, err := svc.ListarFuncionarioPorMatricula(ctx, "999")
		
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "funcionario nao encontrado")
	})
}

func TestListaTodosFuncionarios(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso retornar lista", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}
		
		lista := []model.Funcionario{
			{Id: 1, Nome: "Func 1", Matricula: "10"},
			{Id: 2, Nome: "Func 2", Matricula: "20"},
		}

		m.On("BuscarTodosFuncionarios", ctx).Return(lista, nil)

		res, err := svc.ListaTodosFuncionarios(ctx)
		
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "Func 1", res[0].Nome)
	})

	t.Run("Erro de scan no repositorio", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}

		m.On("BuscarTodosFuncionarios", ctx).Return(nil, Errors.ErrFalhaAoEscanearDados)

		res, err := svc.ListaTodosFuncionarios(ctx)
		
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro interno ao processar dados")
	})
}

func TestDeletarFuncionario(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao deletar", func(t *testing.T) {
		m := new(MockFuncionarioRepo)
		svc := &FuncionarioService{FuncionarioRepo: m}

		m.On("DeletarFuncionario", ctx, "100").Return(nil)

		err := svc.DeletarFuncionario(ctx, "100")
		assert.NoError(t, err)
	})

	t.Run("Erro matricula invalida", func(t *testing.T) {
		svc := &FuncionarioService{}
		// Aqui não precisa de mock pois a função VerificaMatricula falha antes
		err := svc.DeletarFuncionario(ctx, "abc")
		assert.Error(t, err)
	})
}