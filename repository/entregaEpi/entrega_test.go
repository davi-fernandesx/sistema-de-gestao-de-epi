package entregaepi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	estoque "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Estoque"
	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEntrega(t *testing.T) {

	db := integracao.SetupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	repo := NewEntregaRepository(db)
	repoEstoque := estoque.NewEstoqueRepository(db)
	iddep := integracao.CreateDepartamento(t, db)
	idFuncao := integracao.CreateFuncao(t, db, iddep)
	idtam := integracao.CreateTamanho(t, db)
	idproc := integracao.CreateProtecao(t, db)
	idEpi := integracao.CreateEpi(t, db, idproc)
	idfunc := integracao.CreateFuncionario(t, db, iddep, idFuncao)
	_ = integracao.CreateEntradaEpi(t, db, idfunc, idEpi, idproc, idtam)

	entregas := []model.EntregaParaInserir{

		{
			ID_funcionario:     idfunc,
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "teste123",
			Itens: []model.ItemParaInserir{
				{
					ID_epi:         idEpi,
					ID_tamanho:     idtam,
					Quantidade:     10,
					Valor_unitario: decimal.NewFromFloat(45.99),
				},
			},
			Id_troca: sql.NullInt64{Int64: 0},
		},
	}

	t.Run("sucesso ao realizar todas as etapas de uma entraga de epi", func(t *testing.T) {

		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		idEntrega, err := repo.Addentrega(ctx, tx, entregas[0])
		require.NoError(t, err)

		defer tx.Rollback()
		var id int64
		err = tx.QueryRow("SELECT id FROM entrega_epi WHERE id = 1").
			Scan(&id)

		require.Equal(t, idEntrega, id)

		t.Run("sucesso ao listar os lotes(entradas)", func(t *testing.T) {

			for _, item := range entregas[0].Itens {

				if item.Quantidade <= 0 {

					continue
				}
				entradas, err := repoEstoque.ListarLotesParaConsumo(ctx, tx, item.ID_epi, item.ID_tamanho)
				require.NoError(t, err)
				fmt.Println(err)
				require.NotNil(t, entradas)

				quantidadeRestante := item.Quantidade
				for _, entrada := range entradas {

					if quantidadeRestante == 0 {
						break
					}

					var quantidadeParaAbater int

					if entrada.Quantidade >= quantidadeRestante {
						quantidadeParaAbater = quantidadeRestante
						quantidadeRestante = 0
					} else {

						quantidadeParaAbater = entrada.Quantidade
						quantidadeRestante = quantidadeRestante - entrada.Quantidade
					}

					err := repoEstoque.AbaterEstoqueLote(ctx, tx, entrada.ID, quantidadeParaAbater)
					require.NoError(t, err)

					quantidadeEsperada := entrada.Quantidade - item.Quantidade
					var quantidadeAtual int

					query := `select quantidadeAtual from entrada_epi where id = @p1`
					err = tx.QueryRowContext(ctx, query, sql.Named("p1", entrada.ID)).Scan(&quantidadeAtual)

					require.Equal(t, quantidadeEsperada, quantidadeAtual)

					err = repoEstoque.RegistrarItemEntrega(ctx, tx, item.ID_epi, item.ID_tamanho,
						quantidadeParaAbater, idEntrega, entrada.ID, entrada.ValorUnitario)

					require.NoError(t, err)
				}
			}

			err = tx.Commit()
			require.NoError(t, err)
		})

	})

	t.Run("passando um id de funcionario que nao existe", func(t *testing.T) {

		entregaFake := model.EntregaParaInserir{

			ID_funcionario:     12,
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "vwdyuuwdyu",
			Itens: []model.ItemParaInserir{
				{
					ID_epi:         idEpi,
					ID_tamanho:     idtam,
					Quantidade:     10,
					Valor_unitario: decimal.NewFromFloat(45.99),
				},
			},
			Id_troca: sql.NullInt64{Int64: 0},
		}

		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		idEntrega, err := repo.Addentrega(ctx, tx, entregaFake)
		require.Error(t, err)
		require.Equal(t, idEntrega, int64(0))
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))

	})

	t.Run("erro ao passar um id entrada que não existe", func(t *testing.T) {

		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)

		err = repoEstoque.AbaterEstoqueLote(ctx, tx, 45, 10)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))
	})

	t.Run("buscando uma entrega", func(t *testing.T) {

		entrega, err := repo.BuscaEntrega(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, entrega)
		require.NotEmpty(t, entrega)
		require.Equal(t, entrega.Id, int64(1))
		require.Equal(t, entrega.Assinatura_Digital, entregas[0].Assinatura_Digital)

	})

	t.Run("nao achou nenhum entrega com esse id", func(t *testing.T) {

		entrega, err := repo.BuscaEntrega(ctx, 13)
		require.NoError(t, err)
		require.Nil(t, entrega)
		
	})

	t.Run("buscando todas as entregas", func(t *testing.T) {


		entregasTest, err:= repo.BuscaTodasEntregas(ctx)
		require.NoError(t, err)
		require.Equal(t,entregas[0].Itens[0].ID_epi, int64(entregasTest[0].Itens[0].Epi.Id))

	})

	t.Run("sucesso ao fazer um sofdelete", func(t *testing.T) {


		err := repo.CancelarEntrega(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("entrega nao encontrada pelo id", func(t *testing.T) {

		err := repo.CancelarEntrega(ctx, 17)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))		
	})
}
