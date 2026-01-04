package entrada

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service/entrada/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSalvar(t *testing.T) {

	ctx := context.Background()
	mockRepo := mocks.NewEntradaRepository(t)
	service := NewEntradaService(mockRepo)

	entrada := &model.EntradaEpiInserir{
		ID_epi:         1,
		Id_tamanho:     1,
		Data_entrada:   *configs.NewDataBrPtr(time.Now().AddDate(0, 0, 0)),
		Quantidade:     1,
		DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
		DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
		Lote:           "trf-45",
		Fornecedor:     "xtc",
		ValorUnitario:  decimal.NewFromFloat(12.99),
	}

	t.Run("sucesso ao salvar entrada", func(t *testing.T) {

		mockRepo.On("AddEntradaEpi", ctx, mock.Anything).Return(nil)

		err := service.SalvarEntrada(ctx, entrada)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})

	t.Run("data de fabricacao igual ou menor que a data de validade", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)
		dataBase := time.Now().AddDate(-1, 0, 0)
		test := &model.EntradaEpiInserir{
			ID_epi:         1,
			Id_tamanho:     1,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     1,
			DataFabricacao: *configs.NewDataBrPtr(time.Now()),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(-1, 0, 0)),
			Lote:           "trf-45",
			Fornecedor:     "xtc",
			ValorUnitario:  decimal.NewFromFloat(12.99),
		}

		test2 := &model.EntradaEpiInserir{
			ID_epi:         1,
			Id_tamanho:     1,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     1,
			DataFabricacao: *configs.NewDataBrPtr(dataBase),
			DataValidade:   *configs.NewDataBrPtr(dataBase),
			Lote:           "trf-45",
			Fornecedor:     "xtc",
			ValorUnitario:  decimal.NewFromFloat(12.99),
		}

		err := service.SalvarEntrada(ctx, test)

		require.Error(t, err)
		require.True(t, errors.Is(err, errDataMenorValidade))

		err = service.SalvarEntrada(ctx, test2)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrDataIgual))
		mockRepo.AssertExpectations(t)

	})

	t.Run("verificando se a data de entrada é menor que hoje", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)
		entrada := &model.EntradaEpiInserir{
			ID_epi:         1,
			Id_tamanho:     1,
			Data_entrada:   *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
			Quantidade:     1,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			Lote:           "trf-45",
			Fornecedor:     "xtc",
			ValorUnitario:  decimal.NewFromFloat(12.99),
		}

		err := service.SalvarEntrada(ctx, entrada)

		require.Error(t, err)
		require.True(t, errors.Is(err, errDataMenor))
		mockRepo.AssertExpectations(t)

	})

	t.Run("epi ou tamanho nao esta cadastrado no sistema", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)
		entrada := &model.EntradaEpiInserir{
			ID_epi:         21,
			Id_tamanho:     13,
			Data_entrada:   *configs.NewDataBrPtr(time.Now().AddDate(0, 0, 1)),
			Quantidade:     1,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			Lote:           "trf-45",
			Fornecedor:     "xtc",
			ValorUnitario:  decimal.NewFromFloat(12.99),
		}
		mockRepo.On("AddEntradaEpi", mock.Anything, mock.Anything).Return(Errors.ErrDadoIncompativel)

		err := service.SalvarEntrada(ctx, entrada)
		fmt.Println(err)
		require.Error(t, err)

		require.True(t, errors.Is(err, ErrNaoCadastrado))
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do generico do banco de dados", func(t *testing.T) {

		entrada := &model.EntradaEpiInserir{
			ID_epi:         21,
			Id_tamanho:     13,
			Data_entrada:   *configs.NewDataBrPtr(time.Now().AddDate(0, 0, 1)),
			Quantidade:     1,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			Lote:           "trf-45",
			Fornecedor:     "xtc",
		}

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)
		mockRepo.On("AddEntradaEpi", mock.Anything, mock.Anything).Return(Errors.ErrSalvar)
		err := service.SalvarEntrada(ctx, entrada)
		fmt.Println(err)
		require.Error(t, err)

		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))
		mockRepo.AssertExpectations(t)

	})
}

func TestListarEntrada(t *testing.T) {

	ctx := context.Background()

	entrada := model.EntradaEpi{
		ID:               1,
		ID_epi:           2,
		Nome:             "luva",
		Fabricante:       "test",
		CA:               "32321",
		Descricao:        "hjfhjf",
		DataFabricacao:   *configs.NewDataBrPtr(time.Now()),
		DataValidade:     *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
		DataValidadeCa:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 1)),
		IDprotecao:       1,
		NomeProtecao:     "maos",
		Id_Tamanho:       2,
		TamanhoDescricao: "g",
		Quantidade:       12,
		Data_entrada:     *configs.NewDataBrPtr(time.Now()),
		Lote:             "323-tg",
		Fornecedor:       "tydfyrf",
		ValorUnitario:    decimal.NewFromFloat(34.99),
	}

	t.Run("sucesso ao buscar uma entrada", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarEntrada", ctx, 1).Return(entrada, nil)

		test, err := service.ListarEntrada(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, test)
		require.Equal(t, entrada.ID, test.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha no banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarEntrada", ctx, 1).Return(model.EntradaEpi{}, errors.New("erros genericos"))

		test, err := service.ListarEntrada(ctx, 1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		test, err := service.ListarEntrada(ctx, -1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)
	})

	t.Run("entrada não encontrada", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarEntrada", ctx, 1).Return(model.EntradaEpi{}, Errors.ErrBuscarTodos)

		test, err := service.ListarEntrada(ctx, 1)
		require.Nil(t, err)
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

}

func TestListasTodasEntradas(t *testing.T) {

	ctx := context.Background()

	entradas := []model.EntradaEpi{

		{
			ID:               1,
			ID_epi:           2,
			Nome:             "luva",
			Fabricante:       "test",
			CA:               "32321",
			Descricao:        "hjfhjf",
			DataFabricacao:   *configs.NewDataBrPtr(time.Now()),
			DataValidade:     *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			DataValidadeCa:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 1)),
			IDprotecao:       1,
			NomeProtecao:     "maos",
			Id_Tamanho:       2,
			TamanhoDescricao: "g",
			Quantidade:       12,
			Data_entrada:     *configs.NewDataBrPtr(time.Now()),
			Lote:             "323-tg",
			Fornecedor:       "tydfyrf",
			ValorUnitario:    decimal.NewFromFloat(34.99),
		},
		{
			ID:               2,
			ID_epi:           3,
			Nome:             "bota",
			Fabricante:       "test",
			CA:               "32321",
			Descricao:        "hjfhjf",
			DataFabricacao:   *configs.NewDataBrPtr(time.Now()),
			DataValidade:     *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			DataValidadeCa:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 1)),
			IDprotecao:       1,
			NomeProtecao:     "maos",
			Id_Tamanho:       2,
			TamanhoDescricao: "g",
			Quantidade:       12,
			Data_entrada:     *configs.NewDataBrPtr(time.Now()),
			Lote:             "323-tg",
			Fornecedor:       "tydfyrf",
			ValorUnitario:    decimal.NewFromFloat(34.99),
		},
		{
			ID:               5,
			ID_epi:           4,
			Nome:             "luva",
			Fabricante:       "test",
			CA:               "32321",
			Descricao:        "hjfhjf",
			DataFabricacao:   *configs.NewDataBrPtr(time.Now()),
			DataValidade:     *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
			DataValidadeCa:   *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 1)),
			IDprotecao:       1,
			NomeProtecao:     "maos",
			Id_Tamanho:       2,
			TamanhoDescricao: "g",
			Quantidade:       12,
			Data_entrada:     *configs.NewDataBrPtr(time.Now()),
			Lote:             "323-tg",
			Fornecedor:       "tydfyrf",
			ValorUnitario:    decimal.NewFromFloat(34.99),
		},
	}

	t.Run("sucesso ao buscar todas as entradas", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarTodasEntradas", ctx).Return(entradas, nil)

		tests, err := service.ListasTodasEntradas(ctx)
		require.NoError(t, err)
		require.NotNil(t, tests)
		require.NotEmpty(t, tests)

		require.Equal(t, entradas[0].ID, tests[0].ID)
		require.Equal(t, entradas[1].ID, tests[1].ID)
		require.Equal(t, entradas[2].ID, tests[2].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("falha no banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarTodasEntradas", ctx).Return([]model.EntradaEpi{}, errors.New("erro generico bd"))

		test, err := service.ListasTodasEntradas(ctx)
		require.Error(t, err)
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

	t.Run("banco nao tras nenhuma entrada", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("BuscarTodasEntradas", ctx).Return(nil, nil)

		test, err := service.ListasTodasEntradas(ctx)
		require.NoError(t, err)
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)
	})

}

func TestCancelarEntrada(t *testing.T) {

	ctx := context.Background()
	t.Run("sucesso ao cancelar uma entrada", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("CancelarEntrada", ctx, 1).Return(nil)

		err := service.DeletarEntradas(ctx, 1)
		require.NoError(t, err)
		require.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("id negativo passado como parametro", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		err := service.DeletarEntradas(ctx, -1)
		require.Error(t, err)
		require.NotNil(t, err)

		require.True(t, errors.Is(err, ErrId))

		mockRepo.AssertExpectations(t)
	})

	t.Run("entrada nao encontrada para cancelamento", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("CancelarEntrada", ctx, 1).Return(Errors.ErrNaoEncontrado)

		err := service.DeletarEntradas(ctx, 1)
		require.Error(t, err)
		require.NotNil(t, err)
 
		require.True(t, errors.Is(err, ErrId))

		mockRepo.AssertExpectations(t)

	})

	t.Run("erros do banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntradaRepository(t)
		service := NewEntradaService(mockRepo)

		mockRepo.On("CancelarEntrada", ctx, 1).Return(errors.New("erro generico"))

		err := service.DeletarEntradas(ctx, 1)
		require.Error(t, err)
		require.NotNil(t, err)
 
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))

		mockRepo.AssertExpectations(t)

	})

}
