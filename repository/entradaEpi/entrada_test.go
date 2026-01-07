package entradaepi

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEntrada(t *testing.T) {

	db := integracao.SetupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	repo := NewEntradaRepository(db)
	idtam := integracao.CreateTamanho(t, db)
	idproc := integracao.CreateProtecao(t, db)
	idEpi := integracao.CreateEpi(t, db, idproc)

	entradas := []model.EntradaEpiInserir{

		{
			ID_epi:         idEpi,
			Id_tamanho:     idtam,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     11,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -30)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(0, 12, 0)),
			Lote:           "tgf-988",
			Fornecedor:     "test",
			ValorUnitario:  decimal.NewFromFloat(34.99),
		},
		{

			ID_epi:         idEpi,
			Id_tamanho:     idtam,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     11,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -10)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(0, 2, 0)),
			Lote:           "uhj-234",
			Fornecedor:     "test546",
			ValorUnitario:  decimal.NewFromFloat(56.99),
		},
		{
			ID_epi:         idEpi,
			Id_tamanho:     idtam,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     11,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -90)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(0, 7, 0)),
			Lote:           "e4dr-56565",
			Fornecedor:     "tyest43",
			ValorUnitario:  decimal.NewFromFloat(88.99),
		},
	}

	t.Run("sucesso ao adicionar uma entrada", func(t *testing.T) {

		err := repo.AddEntradaEpi(ctx, &entradas[0])
		require.NoError(t, err)

		var idEpi int
		var id int

		err = db.QueryRow("SELECT id, IdEpi FROM entrada_epi WHERE id = 1").
			Scan(&id, &idEpi)

		require.NoError(t, err)
		fmt.Println(err)
		require.Equal(t, entradas[0].ID_epi, idEpi)
	})

	t.Run("id epi ou id tamanho nao existe no sistema", func(t *testing.T) {

		entradaFake := model.EntradaEpiInserir{

			ID_epi:         67,
			Id_tamanho:     idtam,
			Data_entrada:   *configs.NewDataBrPtr(time.Now()),
			Quantidade:     11,
			DataFabricacao: *configs.NewDataBrPtr(time.Now().AddDate(0, 0, -90)),
			DataValidade:   *configs.NewDataBrPtr(time.Now().AddDate(0, 7, 0)),
			Lote:           "e4dr-56565",
			Fornecedor:     "tyest43",
			ValorUnitario:  decimal.NewFromFloat(88.99),
		}

		err:= repo.AddEntradaEpi(ctx, &entradaFake)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))
	})
}
