package tipoprotecao

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	"github.com/stretchr/testify/require"
)

func Test_AddProtecao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTipoProtecaoRepository(db)
	protecao := model.TipoProtecao{Nome: model.Proteção_da_Cabeça_e_Face}
	query := regexp.QuoteMeta(`insert into protecao values (@protecao)`)

	t.Run("sucesso ao adicionar uma protecao", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(protecao.Nome).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddProtecao(ctx, &protecao)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico ao adicionar uma protecao", func(t *testing.T) {
		dbErr := errors.New("falha de conexão")

		mock.ExpectExec(query).
			WithArgs(protecao.Nome).
			WillReturnError(dbErr)

		err := repo.AddProtecao(ctx, &protecao)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoAdicionarProtecao)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarProtecao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTipoProtecaoRepository(db)
	protecao := model.TipoProtecao{ID: 1, Nome: model.Proteção_das_Mãos_e_Braços}
	// Assumindo que a query foi corrigida para "select id, protecao from protecao where id = @id"
	query := regexp.QuoteMeta(`select id, protecao from protecao where id = @id`)

	t.Run("sucesso ao buscar uma protecao", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "protecao"}).AddRow(protecao.ID, protecao.Nome)

		mock.ExpectQuery(query).WithArgs(protecao.ID).WillReturnRows(rows)

		protecaoDB, err := repo.BuscarProtecao(ctx, protecao.ID)
		require.NoError(t, err)
		require.NotNil(t, protecaoDB)
		require.Equal(t, &protecao, protecaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - protecao nao encontrada", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectQuery(query).WithArgs(idNaoExistente).WillReturnError(sql.ErrNoRows)

		protecaoDB, err := repo.BuscarProtecao(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoProcurarProtecao)
		require.Nil(t, protecaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear os dados", func(t *testing.T) {
		// Simulando erro de scan retornando a coluna com tipo de dado errado
		rows := sqlmock.NewRows([]string{"id", "protecao"}).AddRow("id-invalido", protecao.Nome)

		mock.ExpectQuery(query).WithArgs(protecao.ID).WillReturnRows(rows)

		protecaoDB, err := repo.BuscarProtecao(ctx, protecao.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Nil(t, protecaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarTodasProtecao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTipoProtecaoRepository(db)
	query := regexp.QuoteMeta(`select id, protecao from protecao`)
	protecoesEsperadas := []model.TipoProtecao{
		{ID: 1, Nome: model.Proteção_para_os_Pés_e_Pernas},
		{ID: 2, Nome: model.Proteção_do_Corpo},
	}

	t.Run("sucesso ao buscar todas as protecoes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "protecao"}).
			AddRow(protecoesEsperadas[0].ID, protecoesEsperadas[0].Nome).
			AddRow(protecoesEsperadas[1].ID, protecoesEsperadas[1].Nome)

		mock.ExpectQuery(query).WillReturnRows(rows)

		protecoesDB, err := repo.BuscarTodasProtecao(ctx)
		require.NoError(t, err)
		require.Len(t, protecoesDB, 2)
		require.Equal(t, protecoesEsperadas, protecoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar a query de busca", func(t *testing.T) {
		dbErr := errors.New("falha na consulta")
		mock.ExpectQuery(query).WillReturnError(dbErr)

		protecoesDB, err := repo.BuscarTodasProtecao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoBuscarTodasAsProtecoes)
		require.Empty(t, protecoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados durante a iteracao", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "protecao"}).
			AddRow(protecoesEsperadas[0].ID, protecoesEsperadas[0].Nome).
			AddRow(protecoesEsperadas[1].ID, nil) // Segunda linha com valor nulo para causar erro no scan

		mock.ExpectQuery(query).WillReturnRows(rows)

		protecoesDB, err := repo.BuscarTodasProtecao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Nil(t, protecoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro apos a iteracao (linhas.Err)", func(t *testing.T) {
		iterErr := errors.New("erro durante a iteracao")
		rows := sqlmock.NewRows([]string{"id", "protecao"}).
			AddRow(protecoesEsperadas[0].ID, protecoesEsperadas[0].Nome).
			CloseError(iterErr)

		mock.ExpectQuery(query).WillReturnRows(rows)

		protecoesDB, err := repo.BuscarTodasProtecao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoIterarSobreProtecoes)
		require.Nil(t, protecoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_DeletarProtecao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTipoProtecaoRepository(db)
	idParaDeletar := 1
	query := regexp.QuoteMeta(`delete from protecao where id = @id`)

	t.Run("sucesso ao deletar uma protecao", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeletarProtecao(ctx, idParaDeletar)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})


	t.Run("erro - protecao nao encontrada para deletar", func(t *testing.T) {
		idNaoExistente := 99
		mock.ExpectExec(query).
			WithArgs(idNaoExistente).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeletarProtecao(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrProtecaoNaoEncontrada)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico do banco de dados ao deletar", func(t *testing.T) {
		dbErr := errors.New("erro de execucao")
		mock.ExpectExec(query).WithArgs(idParaDeletar).WillReturnError(dbErr)

		err := repo.DeletarProtecao(ctx, idParaDeletar)
		require.Error(t, err)
		require.Equal(t, dbErr, err) // O erro é retornado diretamente
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao obter linhas afetadas", func(t *testing.T) {
		driverErr := errors.New("driver: RowsAffected not supported")
		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewErrorResult(driverErr))

		err := repo.DeletarProtecao(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrLinhasAfetadas)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}