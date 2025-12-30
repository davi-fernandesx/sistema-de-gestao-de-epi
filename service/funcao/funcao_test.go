package funcao

import (
	"context"
	"errors"
	"testing"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) AddFuncao(ctx context.Context, funcao *model.FuncaoInserir) error {

	args := m.Called(ctx, funcao)
	return args.Error(0)
}
func (m *MockRepo) DeletarFuncao(ctx context.Context, id int) error {

	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepo) BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error) {

	args := m.Called(ctx, id)
	var d *model.Funcao

	if args.Get(0) != nil {

		d = args.Get(0).(*model.Funcao)
	}

	return d, args.Error(1)
}

func (m *MockRepo) UpdateFuncao(ctx context.Context, id int, funcao string) (int64, error) {

	args := m.Called(ctx, id, funcao)

	return args.Get(0).(int64), args.Error(1)
}
func (m *MockRepo) BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error) {

	args := m.Called(ctx)
	var d []model.Funcao

	if args.Get(0) != nil {

		d = args.Get(0).([]model.Funcao)
	}

	return d, args.Error(1)
}

func (m *MockRepo) PossuiFuncionariosVinculados(ctx context.Context, id int) (bool, error) {

	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func Mock() (*MockRepo, *FuncaoService) {

	mock := new(MockRepo)
	serv := &FuncaoService{FuncaoRepo: mock}

	return mock, serv
}

func TestSalvarFuncao(t *testing.T) {

	ctx := context.Background()
	Func := &model.FuncaoInserir{ID: 1, Funcao: "dev", IdDepartamento: 1}

	t.Run("sucesso ao salvar", func(t *testing.T) {

		mock, serv := Mock()
		mock.On("AddFuncao", ctx, Func).Return(nil)

		err := serv.SalvarFuncao(ctx, Func)
		assert.NoError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("Erro ao tentar adicionar uma funcao com menos de 2 caracteres", func(t *testing.T) {

		Func := &model.FuncaoInserir{ID: 1, Funcao: "d", IdDepartamento: 1}
		mock, serv := Mock()

		err := serv.SalvarFuncao(ctx, Func)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrFuncaoMinCaracteres))
		mock.AssertExpectations(t)

	})

	t.Run("erro ao tentar adicionar uma funcao que ja existe", func(t *testing.T) {
		Func := &model.FuncaoInserir{ID: 1, Funcao: "dev", IdDepartamento: 1}

		mock, serv := Mock()

		mock.On("AddFuncao", ctx, Func).Return(Errors.ErrDadoIncompativel)

		err := serv.SalvarFuncao(ctx, Func)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrId))
		mock.AssertExpectations(t)
	})

	t.Run("erro do banco de dados", func(t *testing.T) {
		Func := &model.FuncaoInserir{ID: 1, Funcao: "dev", IdDepartamento: 1}

		mock, serv := Mock()

		mock.On("AddFuncao", ctx, Func).Return(Errors.ErrInternal)

		err := serv.SalvarFuncao(ctx, Func)
		assert.Error(t, err)
		mock.AssertExpectations(t)
	})

}

func TestBuscaFuncao(t *testing.T) {

	ctx := context.Background()
	funcs := []*model.Funcao{

		{ID: 1, Funcao: "dev", IdDepartamento: 1, NomeDepartamento: "t1"},
		{ID: 2, Funcao: "rh", IdDepartamento: 2, NomeDepartamento: "adm"},
		{ID: 3, Funcao: "analista", IdDepartamento: 3, NomeDepartamento: "adm"},
		{ID: 4, Funcao: "qa", IdDepartamento: 4, NomeDepartamento: "ti"},
	}
	t.Run("sucesso ao buscar uma funcao", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarFuncao", ctx, 3).Return(funcs[2], nil)

		test, err := serv.ListarFuncao(ctx, 3)
		assert.NoError(t, err)
		assert.Equal(t, test.ID, funcs[2].ID)
		assert.Equal(t, test.Funcao, funcs[2].Funcao)

		mock.AssertExpectations(t)

	})
	t.Run("funcao nao encontrada", func(t *testing.T) {
		mock, serv := Mock()

		mock.On("BuscarFuncao", ctx, 5).Return(nil, Errors.ErrNaoEncontrado)

		_, err := serv.ListarFuncao(ctx, 5)
		assert.Error(t, err)

		assert.True(t, errors.Is(err, ErrRegistroNaoEncontrado))
		mock.AssertExpectations(t)
	})

	t.Run("erro do banco de dados", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarFuncao", ctx, 1).Return(nil, Errors.ErrFalhaAoEscanearDados)

		test, err := serv.ListarFuncao(ctx, 1)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrRegistroNaoEncontrado))
		require.Empty(t, test)
		mock.AssertExpectations(t)
	})



}

func TestBuscaTodasFuncao(t *testing.T) {

	ctx := context.Background()
	funcs := []model.Funcao{

		{ID: 1, Funcao: "dev", IdDepartamento: 1, NomeDepartamento: "t1"},
		{ID: 2, Funcao: "rh", IdDepartamento: 2, NomeDepartamento: "adm"},
		{ID: 3, Funcao: "analista", IdDepartamento: 3, NomeDepartamento: "adm"},
		{ID: 4, Funcao: "qa", IdDepartamento: 4, NomeDepartamento: "ti"},
	}

	t.Run("sucesso ao buscar todas as funcoes", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarTodasFuncao", ctx).Return(funcs, nil)

		test, err := serv.ListasTodasFuncao(ctx)
		assert.NoError(t, err)
		assert.Equal(t, test[0].ID, funcs[0].ID)
		assert.Equal(t, test[1].ID, funcs[1].ID)
		assert.Equal(t, test[2].ID, funcs[2].ID)
		assert.Equal(t, test[3].ID, funcs[3].ID)

		mock.AssertExpectations(t)
	})

	t.Run("erro generico do banco de dados", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarTodasFuncao", ctx).Return(nil, errors.New("erro generico db"))

		test, err := serv.ListasTodasFuncao(ctx)

		assert.Error(t, err)
		assert.Nil(t, test)

		mock.AssertExpectations(t)

	})
}

func TestFuncaoUpdate(t *testing.T) {

	ctx := context.Background()
	funcs := []model.Funcao{

		{ID: 1, Funcao: "dev", IdDepartamento: 1, NomeDepartamento: "t1"},
		{ID: 2, Funcao: "rh", IdDepartamento: 2, NomeDepartamento: "adm"},
	}

	t.Run("sucesso ao atualizar uma funcao", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateFuncao", ctx, funcs[0].ID, funcs[0].Funcao).Return(int64(1), nil)

		err := serv.AtualizarFuncao(ctx, funcs[0].ID, funcs[0].Funcao)

		assert.NoError(t, err)
		assert.Nil(t, err)

		mock.AssertExpectations(t)

	})

	t.Run("erro de funcao ja cadastrada no sistema", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateFuncao", ctx, funcs[0].ID, funcs[0].Funcao).Return(int64(0), Errors.ErrSalvar)

		err := serv.AtualizarFuncao(ctx, funcs[0].ID, funcs[0].Funcao)

		assert.Error(t, err)

		assert.True(t, errors.Is(err, ErrFuncaoCadastrada))
		mock.AssertExpectations(t)
	})

	t.Run("outros erros do banco de dados", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateFuncao", ctx, funcs[0].ID, funcs[0].Funcao).Return(int64(0), errors.New("erro db"))

		err := serv.AtualizarFuncao(ctx, funcs[0].ID, funcs[0].Funcao)

		assert.Error(t, err)
		assert.NotNil(t, err)

		mock.AssertExpectations(t)
	})

}

func TestDeletarFuncao(t *testing.T) {

	ctx := context.Background()
	funcs := []model.Funcao{

		{ID: 1, Funcao: "dev", IdDepartamento: 1, NomeDepartamento: "t1"},
		{ID: 2, Funcao: "rh", IdDepartamento: 2, NomeDepartamento: "adm"},
	}

	t.Run("sucesso ao deletar funcao", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("PossuiFuncionariosVinculados", ctx, 1).Return(false, nil)
		mock.On("DeletarFuncao", ctx, funcs[0].ID).Return(nil)

		err := serv.DeletarFuncao(ctx, funcs[0].ID)
		assert.NoError(t, err)
		assert.Nil(t, err)

		mock.AssertExpectations(t)
	})

	t.Run("id nao encontrado para o delete", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("PossuiFuncionariosVinculados", ctx, 11).Return(false, nil)
		mock.On("DeletarFuncao", ctx, 11).Return(Errors.ErrNaoEncontrado)

		err := serv.DeletarFuncao(ctx, 11)
		assert.Error(t, err)
		assert.NotNil(t, err)

		mock.AssertExpectations(t)
	})

}
