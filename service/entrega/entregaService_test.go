package entrega

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service/entrega/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSalvarEntrega(t *testing.T) {

	ctx := context.Background()

	entrega := model.EntregaParaInserir{

		ID_funcionario:     2,
		Data_entrega:       *configs.NewDataBrPtr(time.Now()),
		Assinatura_Digital: "yufvhudfvjfh",
		Itens: []model.ItemParaInserir{

			{
				ID_epi:         2,
				ID_tamanho:     3,
				Quantidade:     2,
				Valor_unitario: decimal.NewFromFloat(12.99),
			},
			{
				ID_epi:         3,
				ID_tamanho:     2,
				Quantidade:     3,
				Valor_unitario: decimal.NewFromFloat(19.99),
			},
		},
	}

	t.Run("sucesso ao realizar uma entrega", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("Addentrega", ctx, mock.Anything).Return(nil)

		err := service.SalvarEntrega(ctx, entrega)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro  -  adicionar uma data menor que a de hoje", func(t *testing.T) {

		entrega := model.EntregaParaInserir{

			ID_funcionario:     2,
			Data_entrega:       *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -1)),
			Assinatura_Digital: "yufvhudfvjfh",
			Itens: []model.ItemParaInserir{

				{
					ID_epi:         2,
					ID_tamanho:     3,
					Quantidade:     2,
					Valor_unitario: decimal.NewFromFloat(12.99),
				},
				{
					ID_epi:         3,
					ID_tamanho:     2,
					Quantidade:     3,
					Valor_unitario: decimal.NewFromFloat(19.99),
				},
			},
		}
		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		err := service.SalvarEntrega(ctx, entrega)
		require.Error(t, err)
		require.True(t, errors.Is(err, errDataMenor))
		mockRepo.AssertExpectations(t)

	})

	t.Run("falha do banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("Addentrega", ctx, mock.Anything).Return(Errors.ErrInternal)

		err := service.SalvarEntrega(ctx, entrega)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrFalhaNoBancoDeDados))
		mockRepo.AssertExpectations(t)

	})
}

func TestListarEntrega(t *testing.T) {

	ctx := context.Background()

	entrega := &model.EntregaDto{

		Id: 1,
		Funcionario: model.Funcionario_Dto{
			ID:        1,
			Nome:      "rada",
			Matricula: "5443",
			Funcao: model.FuncaoDto{
				ID:     1,
				Funcao: "analista",
				Departamento: model.DepartamentoDto{
					ID:           1,
					Departamento: "ti",
				},
			},
		},
		Data_entrega:       *configs.NewDataBrPtr(time.Now()),
		Assinatura_Digital: "fgsdghdgh",
		Itens: []model.ItemEntregueDto{

			{
				Id: 2,
				Epi: model.EpiDto{
					Id:         1,
					Nome:       "luva",
					Fabricante: "test",
					CA:         "43524",
					Tamanho: []model.TamanhoDto{
						{
							ID:      3,
							Tamanho: "M",
						},
					},
					Descricao:      "luca anti-estatica",
					DataValidadeCa: *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
					Protecao: model.TipoProtecaoDto{
						ID:   2,
						Nome: model.Proteção_das_Mãos_e_Braços,
					},
				},
			},
		},
	}

	t.Run("sucesso ao buscar uma entrega", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaEntrega", ctx, entrega.Id).Return(entrega, nil)

		test, err := service.ListaEntrega(ctx, entrega.Id)
		require.NoError(t, err)
		require.NotNil(t, test)
		require.Equal(t, test.Id, entrega.Id)
		require.Equal(t, test.Funcionario.Nome, test.Funcionario.Nome)

		mockRepo.AssertExpectations(t)

	})

	t.Run("id negativo", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		test, err := service.ListaEntrega(ctx, -1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

	t.Run("entrega nao encontrada", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaEntrega", ctx, entrega.Id).Return(&model.EntregaDto{}, Errors.ErrBuscarTodos)

		test, err := service.ListaEntrega(ctx, entrega.Id)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNaoEncontrado))
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

	t.Run("falha do banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaEntrega", ctx, entrega.Id).Return(&model.EntregaDto{}, errors.New("erro generico bd"))

		test, err := service.ListaEntrega(ctx, entrega.Id)
		require.Error(t, err)
		fmt.Println(err)
		require.True(t, errors.Is(err, ErrInterno))
		require.Empty(t, test)

		mockRepo.AssertExpectations(t)

	})

	t.Run("caso a entrega e o erro seja nill", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaEntrega", ctx, entrega.Id).Return(nil, nil)

		test, err := service.ListaEntrega(ctx, entrega.Id)
		require.Error(t, err)
		require.Empty(t, test)
		require.True(t, errors.Is(err, ErrInterno))

		mockRepo.AssertExpectations(t)

	})
}

func TestBuscarTodasEntregas(t *testing.T) {

	ctx := context.Background()

	entradas := []*model.EntregaDto{

		{
			Id: 1,
			Funcionario: model.Funcionario_Dto{
				ID:        1,
				Nome:      "rada",
				Matricula: "5443",
				Funcao: model.FuncaoDto{
					ID:     1,
					Funcao: "analista",
					Departamento: model.DepartamentoDto{
						ID:           1,
						Departamento: "ti",
					},
				},
			},
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "fgsdghdgh",
			Itens: []model.ItemEntregueDto{

				{
					Id: 2,
					Epi: model.EpiDto{
						Id:         1,
						Nome:       "luva",
						Fabricante: "test",
						CA:         "43524",
						Tamanho: []model.TamanhoDto{
							{
								ID:      3,
								Tamanho: "M",
							},
						},
						Descricao:      "luca anti-estatica",
						DataValidadeCa: *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
						Protecao: model.TipoProtecaoDto{
							ID:   2,
							Nome: model.Proteção_das_Mãos_e_Braços,
						},
					},
				},
			},
		},
		{
			Id: 2,
			Funcionario: model.Funcionario_Dto{
				ID:        31,
				Nome:      "davi",
				Matricula: "5333",
				Funcao: model.FuncaoDto{
					ID:     2,
					Funcao: "analista rh",
					Departamento: model.DepartamentoDto{
						ID:           1,
						Departamento: "administrativo",
					},
				},
			},
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "fgsdghdgh",
			Itens: []model.ItemEntregueDto{

				{
					Id: 2,
					Epi: model.EpiDto{
						Id:         1,
						Nome:       "luva",
						Fabricante: "test",
						CA:         "43524",
						Tamanho: []model.TamanhoDto{
							{
								ID:      1,
								Tamanho: "P",
							},
						},
						Descricao:      "luca anti-estatica",
						DataValidadeCa: *configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
						Protecao: model.TipoProtecaoDto{
							ID:   2,
							Nome: model.Proteção_das_Mãos_e_Braços,
						},
					},
				},
			},
		},
	}

	t.Run("sucesso ao buscar todas as entradas", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaTodasEntregas", ctx).Return(entradas, nil)

		test, err := service.ListarTodasEntregas(ctx)
		require.NoError(t, err)
		require.NotNil(t, test)
		require.Equal(t, test[0].Id, entradas[0].Id)
		require.Equal(t, test[1].Id, entradas[1].Id)

		mockRepo.AssertExpectations(t)

	})

	t.Run("Erro ao buscar os dados no banco", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaTodasEntregas", ctx).Return(nil, Errors.ErrBuscarTodos)

		test, err := service.ListarTodasEntregas(ctx)
		require.Empty(t, test)
		require.Nil(t, err)
		mockRepo.AssertExpectations(t)

	})

	t.Run("outros erros do banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaTodasEntregas", ctx).Return(nil, Errors.ErrAoIterar)

		test, err := service.ListarTodasEntregas(ctx)
		require.Error(t, err)
		require.Empty(t, test)
		mockRepo.AssertExpectations(t)

	})

	t.Run("caso as entregas sejam nil", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("BuscaTodasEntregas", ctx).Return(nil, nil)

		test, err := service.ListarTodasEntregas(ctx)
		require.Nil(t, err)
		require.Empty(t, test)
		mockRepo.AssertExpectations(t)

	})

}

func TestCancelarEntrega(t *testing.T) {

	ctx := context.Background()

	t.Run("sucesso ao cancelar uma entrega", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("CancelarEntrega", ctx, 1).Return(nil)

		err := service.CancelarEntrega(ctx, 1)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("id negativo", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		err := service.CancelarEntrega(ctx, -1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))
		mockRepo.AssertExpectations(t)

	})

	t.Run("outros erros do banco de dados", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("CancelarEntrega", ctx, 1).Return(Errors.ErrInternal)

		err := service.CancelarEntrega(ctx, 1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrId))
		mockRepo.AssertExpectations(t)
	})

	t.Run("id nao encontrado", func(t *testing.T) {

		mockRepo := mocks.NewEntregaInterface(t)
		service := NewEntregaService(mockRepo)

		mockRepo.On("CancelarEntrega", ctx, 1).Return(Errors.ErrNaoEncontrado)

		err := service.CancelarEntrega(ctx, 1)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

}
