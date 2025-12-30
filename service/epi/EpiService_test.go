package epi

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) AddEpi(ctx context.Context, epi *model.EpiInserir) error {

	args := m.Called(ctx, epi)
	return args.Error(0)
}
func (m *MockRepo) DeletarEpi(ctx context.Context, id int) error {

	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepo) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {

	args := m.Called(ctx, id)
	var d *model.Epi

	if args.Get(0) != nil {

		d = args.Get(0).(*model.Epi)
	}

	return d, args.Error(1)
}
func (m *MockRepo) BuscarTodosEpi(ctx context.Context) ([]model.Epi, error) {

	args := m.Called(ctx)
	var d []model.Epi

	if args.Get(0) != nil {

		d = args.Get(0).([]model.Epi)
	}

	return d, args.Error(1)

}
func (m *MockRepo) UpdateEpiNome(ctx context.Context, id int, nome string) error {

	args := m.Called(ctx, id, nome)
	return args.Error(0)
}
func (m *MockRepo) UpdateEpiCa(ctx context.Context, id int, ca string) error {
	args := m.Called(ctx, id, ca)
	return args.Error(0)
}
func (m *MockRepo) UpdateEpiFabricante(ctx context.Context, id int, fabricante string) error {

	args := m.Called(ctx, id, fabricante)
	return args.Error(0)
}
func (m *MockRepo) UpdateEpiDescricao(ctx context.Context, id int, descricao string) error {

	args := m.Called(ctx, id, descricao)
	return args.Error(0)
}
func (m *MockRepo) UpdateEpiDataValidadeCa(ctx context.Context, id int, dataValidadeCa time.Time) error {

	args := m.Called(ctx, id, dataValidadeCa)
	return args.Error(0)
}

func Mock() (*MockRepo, *EpiService) {

	mock := new(MockRepo)
	serv := &EpiService{EpiRepo: mock}

	return mock, serv
}

func TestSalvarEpi(t *testing.T) {

	ctx := context.Background()
	epi := &model.EpiInserir{

		Nome:           "luva",
		Fabricante:     "teste",
		CA:             "3241",
		Descricao:      "luva de borracha",
		DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(4)),
		Idtamanho:      []int{2, 4},
		IDprotecao:     1,
		AlertaMinimo:   10,
	}

	t.Run("sucesso ao salvar epi", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("AddEpi", ctx, epi).Return(nil)

		err := serv.SalvarEpi(ctx, epi)
		require.NoError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("erro - idtamanho nulo", func(t *testing.T) {

		test := &model.EpiInserir{Idtamanho: []int{}}

		mock, serv := Mock()

		err := serv.SalvarEpi(ctx, test)
		fmt.Println(err)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNulo))
		mock.AssertExpectations(t)

	})

	t.Run("Ca ja cadastrado", func(t *testing.T) {

		test := &model.EpiInserir{CA: "12345", Idtamanho: []int{1}}

		mock, serv := Mock()

		mock.On("AddEpi", ctx, test).Return(Errors.ErrSalvar)

		err := serv.SalvarEpi(ctx, test)
		fmt.Println(err)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrCaCadastrado))
	})

	t.Run("id protecao nulo", func(t *testing.T) {
		test := &model.EpiInserir{CA: "12345", Idtamanho: []int{1}, IDprotecao: 11}

		mock, serv := Mock()

		mock.On("AddEpi", ctx, test).Return(Errors.ErrDadoIncompativel)

		err := serv.SalvarEpi(ctx, test)
		fmt.Println(err)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))

	})

	t.Run("erro de conexao com o banco de dados", func(t *testing.T) {
		test := &model.EpiInserir{CA: "12345", Idtamanho: []int{1}, IDprotecao: 11}

		mock, serv := Mock()

		mock.On("AddEpi", ctx, test).Return(Errors.ErrConexaoDb)

		err := serv.SalvarEpi(ctx, test)
		require.Error(t, err)

	})

	t.Run(" teste para ver se o trimSpace esta realmemnte tirando as espaçoes em branco", func(t *testing.T) {

		input := &model.EpiInserir{
			Nome:       "   Luva de Nitrila   ",
			Fabricante: "  3M Brasil  ",
			CA:         " 12345 ",
			Descricao:  "  Equipamento de proteção  ",
			Idtamanho:  []int{1}, // Para passar na validação de nulo
		}

		mockTest, serv := Mock()

		mockTest.On("AddEpi", ctx, mock.MatchedBy(func(e *model.EpiInserir) bool {

			return e.Nome == "Luva de Nitrila" &&
				e.Fabricante == "3M Brasil" &&
				e.CA == "12345" &&
				e.Descricao == "Equipamento de proteção"
		})).Return(nil)

		err := serv.SalvarEpi(ctx, input)

		// 4. Asserts
		require.NoError(t, err)
		mockTest.AssertExpectations(t)

		// Verificação extra no objeto original (ponteiro foi modificado)
		require.Equal(t, "Luva de Nitrila", input.Nome)
		require.Equal(t, "12345", input.CA)
		require.Equal(t, "3M Brasil", input.Fabricante)
		require.Equal(t, "Equipamento de proteção", input.Descricao)

	})
}

func TestListarTodosEpis(t *testing.T) {

	ctx := context.Background()

	epis := []model.Epi{

		{ID: 1, Nome: "luva", Fabricante: "test1", CA: "12345", Descricao: "luva 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(3)),
			AlertaMinimo: 10, IDprotecao: 1, NomeProtecao: "protecao maos", Tamanhos: []model.Tamanhos{{ID: 1, Tamanho: "G"}}},

		{ID: 1, Nome: "bota", Fabricante: "test2", CA: "67890", Descricao: "bota 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(6)),
			AlertaMinimo: 7, IDprotecao: 1, NomeProtecao: "protecao pes", Tamanhos: []model.Tamanhos{{ID: 2, Tamanho: "P"}}},

		{ID: 1, Nome: "capacete", Fabricante: "tes3", CA: "2468", Descricao: "capacete 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(9)),
			AlertaMinimo: 4, IDprotecao: 1, NomeProtecao: "protecao cabeca", Tamanhos: []model.Tamanhos{{ID: 3, Tamanho: "M"}}},
	}

	t.Run("sucesso ao buscar todos os epis", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarTodosEpi", ctx).Return(epis, nil)

		test, err := serv.ListasTodosEpis(ctx)
		require.NoError(t, err)
		require.Equal(t, test[0].Id, epis[0].ID)
		require.Equal(t, test[1].Id, epis[1].ID)
		require.Equal(t, test[2].Id, epis[2].ID)

		mock.AssertExpectations(t)
	})

	t.Run("falha no banco de dados", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarTodosEpi", ctx).Return(nil, Errors.ErrBuscarTodos)

		test, err := serv.ListasTodosEpis(ctx)
		require.Error(t, err)
		require.Nil(t, test)
	})

	t.Run("epis vazios", func(t *testing.T) {

		var episVazio []model.Epi

		mock, serv := Mock()

		mock.On("BuscarTodosEpi", ctx).Return(episVazio, nil)

		test, err := serv.ListasTodosEpis(ctx)
		require.NoError(t, err)
		require.Nil(t, test)

	})
}

func TestListasEpi(t *testing.T) {

	ctx := context.Background()

	epis := []model.Epi{

		{ID: 1, Nome: "luva", Fabricante: "test1", CA: "12345", Descricao: "luva 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(3)),
			AlertaMinimo: 10, IDprotecao: 1, NomeProtecao: "protecao maos", Tamanhos: []model.Tamanhos{{ID: 1, Tamanho: "G"}}},

		{ID: 1, Nome: "bota", Fabricante: "test2", CA: "67890", Descricao: "bota 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().AddDate(0, 3, 0)),
			AlertaMinimo: 7, IDprotecao: 1, NomeProtecao: "protecao pes", Tamanhos: []model.Tamanhos{{ID: 2, Tamanho: "P"}}},

		{ID: 1, Nome: "capacete", Fabricante: "tes3", CA: "2468", Descricao: "capacete 1", DataValidadeCa: *configs.NewDataBrPtr(time.Now().Add(9)),
			AlertaMinimo: 4, IDprotecao: 1, NomeProtecao: "protecao cabeca", Tamanhos: []model.Tamanhos{{ID: 3, Tamanho: "M"}}},
	}

	t.Run("sucesso ao buscar um epi", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("BuscarEpi", ctx, epis[0].ID).Return(&epis[0], nil)

		test, err := serv.ListarEpi(ctx, epis[0].ID)

		require.NoError(t, err)
		require.Equal(t, epis[0].ID, test.Id)
		require.Equal(t, epis[0].DataValidadeCa.Time(), test.DataValidadeCa)

		mock.AssertExpectations(t)

	})

	t.Run("epi nao encontrado", func(t *testing.T) {
		mock, serv := Mock()
		var episVazio *model.Epi

		mock.On("BuscarEpi", ctx, epis[0].ID).Return(episVazio, Errors.ErrNaoEncontrado)

		test, err := serv.ListarEpi(ctx, epis[0].ID)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEpiNaoEncontrado))
		require.Empty(t, test)

		mock.AssertExpectations(t)

	})

	t.Run("falha do banco de dados - escaneaer os dados", func(t *testing.T) {
		mock, serv := Mock()
		var episVazio *model.Epi

		mock.On("BuscarEpi", ctx, epis[0].ID).Return(episVazio, Errors.ErrFalhaAoEscanearDados)

		test, err := serv.ListarEpi(ctx, epis[0].ID)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))
		require.Empty(t, test)

		mock.AssertExpectations(t)

	})

	t.Run("falha do banco de dados - erro generico", func(t *testing.T) {
		mock, serv := Mock()
		var episVazio *model.Epi

		mock.On("BuscarEpi", ctx, epis[0].ID).Return(episVazio, errors.New("erro generico bd"))

		test, err := serv.ListarEpi(ctx, epis[0].ID)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))
		require.Empty(t, test)

		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {
		mock, serv := Mock()

		test, err := serv.ListarEpi(ctx, -11)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))
		require.Empty(t, test)

		mock.AssertExpectations(t)

	})
}

func TestUpdateCa(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao atualizar o CA", func(t *testing.T) {

		mock, serv := Mock()
		mock.On("UpdateEpiCa", ctx, 1, "12653").Return(nil)

		err := serv.AtualizarEpiCa(ctx, 1, "12653")
		require.NoError(t, err)
		require.Nil(t, err)
		mock.AssertExpectations(t)

	})

	t.Run("id passado negativo", func(t *testing.T) {

		mock, serv := Mock()
		err := serv.AtualizarEpiCa(ctx, -11, "43423")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))

		mock.AssertExpectations(t)
	})

	t.Run("Ca ja esta cadastrado no sistema", func(t *testing.T) {

		mock, serv := Mock()
		mock.On("UpdateEpiCa", ctx, 1, "12653").Return(Errors.ErrSalvar)

		err := serv.AtualizarEpiCa(ctx, 1, " 12653 ")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrCaCadastrado))

		mock.AssertExpectations(t)

	})

	t.Run("falha no banco de dados", func(t *testing.T) {
		mock, serv := Mock()
		mock.On("UpdateEpiCa", ctx, 1, "12653").Return(Errors.ErrInternal)

		err := serv.AtualizarEpiCa(ctx, 1, "12653")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

	t.Run("Ca nao esta seguindo as ordens (apenas numeros de 0-9, ate 6 numeros)", func(t *testing.T) {

		mock, serv := Mock()
		err := serv.AtualizarEpiCa(ctx, 1, "12T53") // contendo letras
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrCa))
		mock.AssertExpectations(t)

		err = serv.AtualizarEpiCa(ctx, 1, "12537575") // passando dos 6 numeros
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrCa))
		mock.AssertExpectations(t)

	})
}

func TestAtualizaFabricante(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao atualizar o fabricante", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiFabricante", ctx, 1, "master").Return(nil)

		err := serv.AtualizarEpiFabricante(ctx, 1, "master")
		require.NoError(t, err)
		require.Nil(t, err)

		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizarEpiFabricante(ctx, -4, "master")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrId))

		mock.AssertExpectations(t)
	})

	t.Run("Erro interno (provavelmente o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiFabricante", ctx, 1, "master").Return(Errors.ErrInternal)

		err := serv.AtualizarEpiFabricante(ctx, 1, "master")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

	t.Run("Erro interno (qualquer erro que nao envolva o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiFabricante", ctx, 1, "master").Return(errors.New("erro generico"))

		err := serv.AtualizarEpiFabricante(ctx, 1, "master")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})
}

func TestUpdateNome(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao atualizar o nome", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiNome", ctx, 1, "luva acrilico").Return(nil)

		err := serv.AtualizarEpiNome(ctx, 1, "luva acrilico")
		require.NoError(t, err)
		require.Nil(t, err)

		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizarEpiNome(ctx, -4, "luva acrilico")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrId))

		mock.AssertExpectations(t)
	})

	t.Run("Erro interno (provavelmente o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiNome", ctx, 1, "luva acrilico").Return(Errors.ErrInternal)

		err := serv.AtualizarEpiNome(ctx, 1, "luva acrilico")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

	t.Run("Erro interno (qualquer erro que nao envolva o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiNome", ctx, 1, "luva acrilico").Return(errors.New("erro generico"))

		err := serv.AtualizarEpiNome(ctx, 1, "luva acrilico")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

}

func TestUpdateDescricao(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao atualizar a descricao", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiDescricao", ctx, 1, "descricao x").Return(nil)

		err := serv.AtualizaDescricao(ctx, 1, "descricao x")
		require.NoError(t, err)
		require.Nil(t, err)

		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizaDescricao(ctx, -4, "descricao x")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrId))

		mock.AssertExpectations(t)
	})

	t.Run("Erro interno (provavelmente o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiDescricao", ctx, 1, "descricao x").Return(Errors.ErrInternal)

		err := serv.AtualizaDescricao(ctx, 1, "descricao x")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

	t.Run("Erro interno (qualquer erro que nao envolva o banco de dados)", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiDescricao", ctx, 1, "descricao x").Return(errors.New("erro generico"))

		err := serv.AtualizaDescricao(ctx, 1, "descricao x")
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mock.AssertExpectations(t)

	})

}

func TestUpadateValidadeCa(t *testing.T) {

	ctx := context.Background()
	// Datas para os testes
	hoje := time.Now().Truncate(24 * time.Hour)
	amanha := hoje.AddDate(0, 0, 1)
	ontem := hoje.AddDate(0, 0, -1)
	anos := hoje.AddDate(2, 4, 3)

	t.Run("erro - caso a data zerada(January 1, year 1, 00:00:00 UTC.)", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizaDataValidadeCa(ctx, 1, configs.DataBr{})
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrDataZero))

		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizaDataValidadeCa(ctx, -11, configs.DataBr(amanha))
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrId))

		mock.AssertExpectations(t)
	})

	t.Run("erro - data no passado", func(t *testing.T) {

		mock, serv := Mock()

		err := serv.AtualizaDataValidadeCa(ctx, 11, configs.DataBr(ontem))
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, errDataMenor))

		mock.AssertExpectations(t)

	})

	t.Run("sucesso - data de hoje", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiDataValidadeCa", ctx, 1, hoje).Return(nil)
		err := serv.AtualizaDataValidadeCa(ctx, 1, configs.DataBr(hoje))
		require.NoError(t, err)
		require.Nil(t, err)

		mock.AssertExpectations(t)
	})

	t.Run("sucesso - data futura", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("UpdateEpiDataValidadeCa", ctx, 1, anos).Return(nil)
		err := serv.AtualizaDataValidadeCa(ctx, 1, configs.DataBr(anos))
		require.NoError(t, err)
		require.Nil(t, err)

		mock.AssertExpectations(t)
	})

}

func TestDeleteEpi(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao deletar um epi", func(t *testing.T) {

		mock, serv := Mock()

		mock.On("DeletarEpi", ctx, 1).Return(nil)

		err := serv.DeletarEpi(ctx, 1)
		require.NoError(t, err)
		require.Nil(t, err)
		mock.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mock, serv := Mock()
		err := serv.DeletarEpi(ctx, -1)
		require.Error(t, err)
		require.NotNil(t, err)

		require.True(t, errors.Is(err, ErrId))
		mock.AssertExpectations(t)
	})

	t.Run("algum erro do id passado", func(t *testing.T) {


		mock, serv := Mock()

		mock.On("DeletarEpi", ctx, 1).Return(Errors.ErrAoapagar)

		err := serv.DeletarEpi(ctx, 1)
		require.Error(t, err)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrId))
		mock.AssertExpectations(t)

	})
}
