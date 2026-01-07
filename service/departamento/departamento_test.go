package departamento

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) AddDepartamento(ctx context.Context, departamento *model.Departamento) error {

	args := m.Called(ctx, departamento)

	return args.Error(0)
}
func (m *MockRepo) DeletarDepartamento(ctx context.Context, id int) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
func (m *MockRepo) BuscarDepartamento(ctx context.Context, id int) (*model.DepartamentoDto, error) {

	args := m.Called(ctx, id)
	var d *model.DepartamentoDto

	if args.Get(0) != nil {

		d = args.Get(0).(*model.DepartamentoDto)
	}

	return d, args.Error(1)
}
func (m *MockRepo) BuscarTodosDepartamentos(ctx context.Context) ([]model.DepartamentoDto, error) {

	args := m.Called(ctx)
	var d []model.DepartamentoDto

	if args.Get(0) != nil {

		d = args.Get(0).([]model.DepartamentoDto)
	}

	return d, args.Error(1)
}
func (m *MockRepo) UpdateDepartamento(ctx context.Context, id int, departamento string) (int64, error) {

	args := m.Called(ctx, id, departamento)

	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) PossuiFuncoesVinculadas(ctx context.Context, id int) (bool, error) {

	args := m.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

// --- TESTES DO SERVICE ---

func TestSalvarDepartamento(t *testing.T) {
	ctx := context.Background()
	dep := &model.Departamento{Departamento: "TI"}

	t.Run("Sucesso ao salvar", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("AddDepartamento", ctx, dep).Return(nil)

		err := svc.SalvarDepartamento(ctx, dep)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("Erro no banco ao salvar", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("AddDepartamento", ctx, dep).Return(errors.New("db error"))

		err := svc.SalvarDepartamento(ctx, dep)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao salvar")
	})
}

func TestListarDepartamento(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao listar", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}
		depFake := &model.DepartamentoDto{ID: 1, Departamento: "RH"}

		m.On("BuscarDepartamento", ctx, 1).Return(depFake, nil)

		res, err := svc.ListarDepartamento(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "RH", res.Departamento)
		assert.Equal(t, 1, res.ID)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		svc := &DepartamentoServices{}
		_, err := svc.ListarDepartamento(ctx, 0)
		assert.Equal(t, ErrId, err)
	})

	t.Run("Departamento não encontrado (sql.ErrNoRows)", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("BuscarDepartamento", ctx, 99).Return(nil, sql.ErrNoRows)

		_, err := svc.ListarDepartamento(ctx, 99)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não encontrado")
	})
}

func TestListarTodosDepartamentos(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao listar todos", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}
		listaFake := []model.DepartamentoDto{
			{Departamento: "TI"},
			{Departamento: "Financeiro"},
		}

		m.On("BuscarTodosDepartamentos", ctx).Return(listaFake, nil)

		res, err := svc.ListarTodosDepartamentos(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "TI", res[0].Departamento)
	})

	t.Run("Retorno vazio quando deps é nil", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("BuscarTodosDepartamentos", ctx).Return(nil, nil)

		res, err := svc.ListarTodosDepartamentos(ctx)
		assert.NoError(t, err)
		assert.Empty(t, res)
		assert.NotNil(t, res) // Garante que retornou [] e não null
	})
}

func TestDeletarDepartamento(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao deletar", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("PossuiFuncoesVinculadas", ctx, 1).Return(false, nil)
		m.On("DeletarDepartamento", ctx, 1).Return(nil)

		err := svc.DeletarDepartamento(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Erro de integridade (FK Constraint 547)", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		
		m.On("PossuiFuncoesVinculadas", ctx, 1).Return(false, nil)
		m.On("DeletarDepartamento", ctx, 1).Return(errors.New("error 547 check constraint"))

		err := svc.DeletarDepartamento(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "departamento ja pode estar inativo")
	})

	t.Run("erro - departamento possui vinculo", func(t *testing.T) {

		m:= new(MockRepo)

		svc:= &DepartamentoServices{DepartamentoRepo: m}

		m.On("PossuiFuncoesVinculadas", ctx, 1).Return(true, nil)
		m.On("DeletarDepartamento", ctx, 1).Return(ErrFuncaoComVinculo)

		err := svc.DeletarDepartamento(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "departamento com vinculo ")

	})
}

func TestAtualizarDepartamento(t *testing.T) {
	ctx := context.Background()

	t.Run("Sucesso ao atualizar", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("UpdateDepartamento", ctx, 1, "Vendas").Return(int64(1), nil)

		err := svc.AtualizarDepartamento(ctx, 1, "Vendas")
		assert.NoError(t, err)
	})

	t.Run("Erro de Duplicidade (Unique 2627)", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("UpdateDepartamento", ctx, 1, "Vendas").Return(int64(0), errors.New("error 2627 unique key"))

		err := svc.AtualizarDepartamento(ctx, 1, "Vendas")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro tecnico ao")
	})

	t.Run("Erro ID não existe (0 linhas)", func(t *testing.T) {
		m := new(MockRepo)
		svc := &DepartamentoServices{DepartamentoRepo: m}

		m.On("UpdateDepartamento", ctx, 999, "Vendas").Return(int64(0), nil)

		err := svc.AtualizarDepartamento(ctx, 999, "Vendas")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não encontrado")
	})
}
