package entradaepi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context, error) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock, ctx, err

}

func Test_EntradaEpi(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewEntradaRepository(db)

	entradaInserir := model.EntradaEpiInserir{
		ID_epi:         1,
		Data_entrada:   *configs.NewDataBrPtr(time.Now()),
		Id_tamanho:     2,
		Quantidade:     10,
		DataFabricacao: *configs.NewDataBrPtr(time.Now()),
		DataValidade:   *configs.NewDataBrPtr(time.Now()),
		Lote:           "xyz",
		Fornecedor:     "teste1",
		ValorUnitario:  decimal.NewFromFloat(12.77),
	}

	t.Run("testando o sucesso ao adicionar uma entrada", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("insert into")).WithArgs(entradaInserir.ID_epi, entradaInserir.Id_tamanho, entradaInserir.Data_entrada,
			entradaInserir.Quantidade, entradaInserir.DataFabricacao, entradaInserir.DataValidade, entradaInserir.Lote,
			entradaInserir.Fornecedor, entradaInserir.ValorUnitario).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddEntradaEpi(ctx, &entradaInserir)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro ao adicionar uma entrada", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("insert into")).WithArgs(entradaInserir.ID_epi, entradaInserir.Id_tamanho, entradaInserir.Data_entrada,
			entradaInserir.Quantidade, entradaInserir.DataFabricacao, entradaInserir.DataValidade, entradaInserir.Lote,
			entradaInserir.Fornecedor, entradaInserir.ValorUnitario).WillReturnError(Errors.ErrSalvar)

		err := repo.AddEntradaEpi(ctx, &entradaInserir)
		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrSalvar), "erro tem que ser do tipo salvar")
		require.NoError(t, mock.ExpectationsWereMet())

	})
}

var colunaEntrada = []string{
	"id",
	"id_epi",
	"nome",
	"fabricante",
	"CA",
	"descricao",
	"dataFabricacao",
	"dataValidade",
	"dataValidadeCa",
	"IdProtecao",
	"NomeProtecao",
	"idTamanho",
	"tamanho",
	"quantidade",
	"data_entrada",
	"lote",
	"fornecedor",
	"valorUnitario",
	}

func TestBuscarEntradaPorId(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewEntradaRepository(db)

	id := 1
	dataValidade := time.Now()
	dataFabricacao := time.Now()
	validadeCa := time.Now()
	dataEntrada := time.Now()

	t.Run("sucesso ao achar entrada por id", func(t *testing.T) {

		row := sqlmock.NewRows(colunaEntrada).AddRow(
			id, 23, "luva","master", "1232", "luva borracha", dataFabricacao, dataValidade, validadeCa,
			3, "maos", 2, "G", 10, dataEntrada, "6464-iu", "mister", 13.99)

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("id", id)).WillReturnRows(row)

		resultado, err := repo.BuscarEntrada(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, resultado)
		require.Equal(t, "luva", resultado.Nome)
		require.Equal(t, "master", resultado.Fabricante)

	})

	t.Run("deve retornar um erro", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("id", id)).WillReturnError(sql.ErrConnDone)

		resultado, err := repo.BuscarEntrada(ctx, id)

		require.Error(t, err)
		require.Equal(t, model.EntradaEpi{}, resultado)

	})

	t.Run("erro de no rows", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("id", id)).WillReturnError(sql.ErrNoRows)

		resultado, err := repo.BuscarEntrada(ctx, id)

		require.Error(t, err)
		require.Equal(t, model.EntradaEpi{}, resultado)

	})

	if err := mock.ExpectationsWereMet(); err != nil {

		t.Errorf("algumas expectativas nao foram atendidas")
	}

}

func TestBuscarTodasEntrada(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewEntradaRepository(db)

	dataValidade := time.Now()
	dataFabricacao := time.Now()
	validadeCa := time.Now()
	dataEntrada:= time.Now()

	t.Run("sucesso ao achar entrada por id", func(t *testing.T) {

		row := sqlmock.NewRows(colunaEntrada).AddRow(

			3, 23, "luva","master", "1232", "luva borracha", dataFabricacao, dataValidade, validadeCa,
			3, "maos", 2, "G", 10, dataEntrada, "6464-iu", "EPI-TEST", 13.99,
		).AddRow(

			5, 23, "luva","master", "1232", "luva borracha", dataFabricacao, dataValidade, validadeCa,
			3, "maos", 2, "G", 10, dataEntrada, "6464-iu", "test", 13.99,
		).AddRow(

			8, 23, "luva","master", "1232", "luva borracha", dataFabricacao, dataValidade, validadeCa,
			3, "maos", 2, "G", 10, dataEntrada, "6464-iu", "epi", 13.99,
		)

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WillReturnRows(row)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.NoError(t, err)
		require.NotNil(t, resultado)
		require.Len(t, resultado, 3)
		require.Equal(t, "EPI-TEST", resultado[0].Fornecedor)
		require.Equal(t, "test", resultado[1].Fornecedor)
		require.Equal(t, "epi", resultado[2].Fornecedor)

	})

	t.Run("deve retornar um erro", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WillReturnError(sql.ErrConnDone)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.Error(t, err)
		require.Equal(t, []model.EntradaEpi{}, resultado)
		fmt.Println(resultado)

	})

	t.Run("erro de no rows", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WillReturnError(sql.ErrNoRows)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.Error(t, err)
		require.Equal(t, []model.EntradaEpi{}, resultado)

	})

	if err := mock.ExpectationsWereMet(); err != nil {

		t.Errorf("algumas expectativas nao foram atendidas")
	}

}

func TestCancelarEntrada(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewEntradaRepository(db)

	entradasMock := model.EntradaEpi{
		ID:             1,
		ID_epi:         10,
		Nome:           "Capacete de Segurança",
		Fabricante:     "Marca Segura",
		CA:             "12345",
		Descricao:      "Capacete para proteção contra impactos.",
		DataFabricacao: time.Now().AddDate(0, -6, 0),
		DataValidade:   time.Now().AddDate(2, 0, 0),
		DataValidadeCa: time.Now().AddDate(1, 0, 0),
		IDprotecao:     1,
		NomeProtecao:   "Cabeça",
		Lote:           "LOTE-2025-A1",
		Fornecedor:     "Fornecedor Principal",
	}



	t.Run("testando o sucesso ao cancelar um epi da base de dados", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewResult(0, 1))

		errEpi := repo.CancelarEntrada(ctx, entradasMock.ID)
		require.NoError(t, errEpi)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao cancelar uma entradas", func(t *testing.T) {

		ErroGenericoDb := errors.New("erro ao se conectar com o banco")
		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(entradasMock.ID).WillReturnError(ErroGenericoDb)

		errEpi := repo.CancelarEntrada(ctx, entradasMock.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrInternal), "erro tem que ser do tipo internal")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de linhas afetadas", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		errEpi := repo.CancelarEntrada(ctx, entradasMock.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrLinhasAfetadas), "erro tem que ser do tipo linhas afetadas")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("epi nao encontrado", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewResult(0, 0))

		errEpi := repo.CancelarEntrada(ctx, entradasMock.ID)
		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrNaoEncontrado), "erro tem que ser do tipo nao encontrado")

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
