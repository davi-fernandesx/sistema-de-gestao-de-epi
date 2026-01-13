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
	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
	"github.com/stretchr/testify/require"
)

func TestEpiAdd(t *testing.T) {

	db := integracao.SetupTestDB(t)
	// O defer garante que o DB feche, o cleanup do setupTestDB mata o container
	defer db.Close()
	ctx := context.Background()
	repo := NewEpiRepository(db)
	idTam := integracao.CreateTamanho(t, db)
	idProtec := integracao.CreateProtecao(t, db)

	epis := []model.EpiInserir{

		{
			Nome:           "luva",
			Fabricante:     "test1",
			CA:             "12345",
			Descricao:      "luva de borracha",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{int(idTam)},
			IDprotecao:     int(idProtec),
			AlertaMinimo:   10,
		},
		{
			Nome:           "bota",
			Fabricante:     "test1",
			CA:             "23456",
			Descricao:      "bopta de borracha",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{int(idTam)},
			IDprotecao:      int(idProtec),
			AlertaMinimo:   3,
		},
		{
			Nome:           "mascara",
			Fabricante:     "test1",
			CA:             "34567",
			Descricao:      "mascara de pano",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{int(idTam)},
			IDprotecao:      int(idProtec),
			AlertaMinimo:   4,
		},
	}

	t.Run("sucesso ao adicionar um epi", func(t *testing.T) {

		err := repo.AddEpi(ctx, &epis[0])
		require.NoError(t, err)

		var nome string
		var id int

		err = db.QueryRow("SELECT id, nome FROM epi WHERE id = 1").
			Scan(&id, &nome)

		require.NoError(t, err)
		fmt.Println(err)
		require.Equal(t, "luva", nome)

	})

	t.Run("erro - adicionar um CA repetido", func(t *testing.T) {

		epiFalse := model.EpiInserir{

			Nome:           "repetido",
			Fabricante:     "test1",
			CA:             "12345",
			Descricao:      "repetiddo de borracha",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{int(idTam)},
			IDprotecao:     int(idProtec),
			AlertaMinimo:   10,
		}

		err := repo.AddEpi(ctx, &epiFalse)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrSalvar))
	})

	t.Run("erro - id protecao nao existe no sistema", func(t *testing.T) {

		epiFalse := model.EpiInserir{

			Nome:           "repetido1",
			Fabricante:     "test1",
			CA:             "9999",
			Descricao:      "repetiddo de borracha",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{int(idTam)},
			IDprotecao:     44,
			AlertaMinimo:   1,
		}

		err := repo.AddEpi(ctx, &epiFalse)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))
	})

	t.Run("erro - id tamanho nao existe no sistema", func(t *testing.T) {

		epiFalse := model.EpiInserir{

			Nome:           "repetido2",
			Fabricante:     "test3",
			CA:             "8888",
			Descricao:      "repetido de borracha2",
			DataValidadeCa: *configs.NewDataBrPtr(time.Now()),
			Idtamanho:      []int{8},
			IDprotecao:     int(idProtec),
			AlertaMinimo:   12,
		}

		err := repo.AddEpi(ctx, &epiFalse)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))
	})

	t.Run("sucesso ao buscar um epi", func(t *testing.T) {

		epi, err := repo.BuscarEpi(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, epi)
		require.NotEmpty(t, epi)
		require.Equal(t, epi.ID, 1)
		require.NotEmpty(t, epi.Tamanhos)
		require.NotNil(t, epi.Tamanhos)

	})

	t.Run("epi não encontrado", func(t *testing.T) {

		epi, err := repo.BuscarEpi(ctx, 99)
		require.Error(t, err)
		require.Nil(t, epi)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))
	})

	t.Run("buscar todos os epis", func(t *testing.T) {

		_ = repo.AddEpi(ctx, &epis[1])

		_ = repo.AddEpi(ctx, &epis[2])

		episTest, err := repo.BuscarTodosEpi(ctx)
		require.NoError(t, err)
		// Lista do que esperamos encontrar
		esperados := []string{epis[0].CA,epis[1].CA, epis[2].CA}

		// Lista do que veio do banco
		var encontrados []string
		for _, e := range episTest {
			encontrados = append(encontrados, e.CA)
		}

		// Verifica se os elementos de uma lista estão na outra, sem ligar pra ordem
		require.ElementsMatch(t, esperados, encontrados)
	})

	t.Run("sucesso ao fazer um softdelete", func(t *testing.T) {

		err := repo.DeletarEpi(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("epi nao encontrado para deletar", func(t *testing.T) {

		err := repo.DeletarEpi(ctx, 111)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))
	})

}
