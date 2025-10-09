package tamanhos

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func Test_AddTamanhos(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTamanhoRepository(db)
	tamanho := model.Tamanhos{Tamanho: "M"}
	query := regexp.QuoteMeta(`insert into tamanho values (@tamanho)`)

	t.Run("sucesso ao adicionar um tamanho", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(tamanho.Tamanho).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddTamanhos(ctx, &tamanho)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - tamanho ja existente", func(t *testing.T) {
		mssqlErr := &mssql.Error{Number: 2627}

		mock.ExpectExec(query).
			WithArgs(tamanho.Tamanho).
			WillReturnError(mssqlErr)

		err := repo.AddTamanhos(ctx, &tamanho)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrTamanhoJaExistente)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha generica ao adicionar", func(t *testing.T) {
		dbErr := errors.New("falha de conexão")

		mock.ExpectExec(query).
			WithArgs(tamanho.Tamanho).
			WillReturnError(dbErr)

		err := repo.AddTamanhos(ctx, &tamanho)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoAdicionarTamanho)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarTamanhos(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTamanhoRepository(db)
	tamanho := model.Tamanhos{ID: 1, Tamanho: "G"}
	query := regexp.QuoteMeta("select id, tamanho from tamanho where id = @id")

	t.Run("sucesso ao buscar um tamanho", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).AddRow(tamanho.ID, tamanho.Tamanho)

		mock.ExpectQuery(query).WithArgs(tamanho.ID).WillReturnRows(rows)

		tamanhoDB, err := repo.BuscarTamanhos(ctx, tamanho.ID)
		require.NoError(t, err)
		require.NotNil(t, tamanhoDB)
		require.Equal(t, &tamanho, tamanhoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - tamanho nao encontrado", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectQuery(query).WithArgs(idNaoExistente).WillReturnError(sql.ErrNoRows)

		tamanhoDB, err := repo.BuscarTamanhos(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoProcurarTamanho)
		require.Nil(t, tamanhoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear os dados", func(t *testing.T) {
		// Simulando erro de scan retornando a coluna errada (apenas uma coluna)
		rows := sqlmock.NewRows([]string{"id"}).AddRow(tamanho.ID)

		mock.ExpectQuery(query).WithArgs(tamanho.ID).WillReturnRows(rows)

		tamanhoDB, err := repo.BuscarTamanhos(ctx, tamanho.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Nil(t, tamanhoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarTodosTamanhos(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTamanhoRepository(db)
	query := regexp.QuoteMeta("select id, tamanho from tamanho")
	tamanhosEsperados := []model.Tamanhos{
		{ID: 1, Tamanho: "P"},
		{ID: 2, Tamanho: "M"},
	}

	t.Run("sucesso ao buscar todos os tamanhos", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(tamanhosEsperados[0].ID, tamanhosEsperados[0].Tamanho).
			AddRow(tamanhosEsperados[1].ID, tamanhosEsperados[1].Tamanho)

		mock.ExpectQuery(query).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTodosTamanhos(ctx)
		require.NoError(t, err)
		require.Len(t, tamanhosDB, 2)
		require.Equal(t, tamanhosEsperados, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar a query de busca", func(t *testing.T) {
		dbErr := errors.New("falha na consulta")

		mock.ExpectQuery(query).WillReturnError(dbErr)

		tamanhosDB, err := repo.BuscarTodosTamanhos(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoBuscarTodosOsTamanhos)
		require.Equal(t, repository.ErrAoBuscarTodosOsTamanhos, err)
		require.Empty(t, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados durante a iteracao", func(t *testing.T) {
		// A segunda linha tem um tipo de dado errado para causar falha no Scan
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(tamanhosEsperados[0].ID, tamanhosEsperados[0].Tamanho).
			AddRow("id-invalido", tamanhosEsperados[1].Tamanho)

		mock.ExpectQuery(query).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTodosTamanhos(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Equal(t, repository.ErrFalhaAoEscanearDados, err)
		require.Nil(t, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro apos a iteracao (linhas.Err)", func(t *testing.T) {
		iterErr := errors.New("erro durante a iteracao")

		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(tamanhosEsperados[0].ID, tamanhosEsperados[0].Tamanho).
			CloseError(iterErr) // Simula um erro ao fechar as linhas

		mock.ExpectQuery(query).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTodosTamanhos(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoIterarSobreTamanhos)
		require.Nil(t, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_DeletarTamanhos(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTamanhoRepository(db)
	idParaDeletar := 1
	query := regexp.QuoteMeta(`delete from tamanho where id = @id`)

	t.Run("sucesso ao deletar um tamanho", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeletarTamanhos(ctx, idParaDeletar)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - tamanho nao encontrado para deletar", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectExec(query).
			WithArgs(idNaoExistente).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 linhas afetadas

		err := repo.DeletarTamanhos(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrTamanhoNaoEncontrado)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico do banco de dados ao deletar", func(t *testing.T) {
		dbErr := errors.New("erro de execucao")

		mock.ExpectExec(query).WithArgs(idParaDeletar).WillReturnError(dbErr)

		err := repo.DeletarTamanhos(ctx, idParaDeletar)
		require.Error(t, err)
		require.Equal(t, dbErr, err) // O erro é retornado diretamente
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao obter linhas afetadas", func(t *testing.T) {
		driverErr := errors.New("driver: RowsAffected not supported")

		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewErrorResult(driverErr))

		err := repo.DeletarTamanhos(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrLinhasAfetadas)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}