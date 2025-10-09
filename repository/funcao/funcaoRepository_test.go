package funcao

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

func Test_AddFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	funcao := model.Funcao{Funcao: "Desenvolvedor"}
	query := regexp.QuoteMeta(`insert into funcao (funcao) values (@funcao)`)

	t.Run("sucesso ao adicionar uma funcao", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(funcao.Funcao).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddFuncao(ctx, &funcao)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - funcao ja existente", func(t *testing.T) {
		// Simula o erro específico do MS SQL Server para violação de chave única
		mssqlErr := &mssql.Error{Number: 2627, Message: "Violation of UNIQUE KEY constraint"}

		mock.ExpectExec(query).
			WithArgs(funcao.Funcao).
			WillReturnError(mssqlErr)

		err := repo.AddFuncao(ctx, &funcao)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFuncaoJaExistente)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha generica do banco de dados", func(t *testing.T) {
		dbErr := errors.New("falha de conexão")

		mock.ExpectExec(query).
			WithArgs(funcao.Funcao).
			WillReturnError(dbErr)

		err := repo.AddFuncao(ctx, &funcao)
		require.Error(t, err)
		require.Equal(t, dbErr, err) // Neste caso, o erro é repassado diretamente
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	funcao := model.Funcao{ID: 1, Funcao: "Analista de QA"}
	query := regexp.QuoteMeta(`select id, funcao from funcao where id = @id`)

	t.Run("sucesso ao buscar uma funcao", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "funcao"}).AddRow(funcao.ID, funcao.Funcao)

		mock.ExpectQuery(query).WithArgs(funcao.ID).WillReturnRows(rows)

		funcaoDB, err := repo.BuscarFuncao(ctx, funcao.ID)
		require.NoError(t, err)
		require.NotNil(t, funcaoDB)
		require.Equal(t, &funcao, funcaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - funcao nao encontrada", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectQuery(query).WithArgs(idNaoExistente).WillReturnError(sql.ErrNoRows)

		funcaoDB, err := repo.BuscarFuncao(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoProcurarFuncao)
		require.Nil(t, funcaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear os dados da funcao", func(t *testing.T) {
		// Simulando erro de scan retornando colunas erradas (faltando 'funcao')
		rows := sqlmock.NewRows([]string{"id"}).AddRow(funcao.ID)

		mock.ExpectQuery(query).WithArgs(funcao.ID).WillReturnRows(rows)

		funcaoDB, err := repo.BuscarFuncao(ctx, funcao.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Nil(t, funcaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarTodasFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	query := regexp.QuoteMeta(`select id, funcao from funcao`)
	funcoesEsperadas := []model.Funcao{
		{ID: 1, Funcao: "Gerente"},
		{ID: 2, Funcao: "Coordenador"},
	}

	t.Run("sucesso ao buscar todas as funcoes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "funcao"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao).
			AddRow(funcoesEsperadas[1].ID, funcoesEsperadas[1].Funcao)

		mock.ExpectQuery(query).WillReturnRows(rows)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.NoError(t, err)
		require.NotNil(t, funcoesDB)
		require.Len(t, funcoesDB, 2)
		require.Equal(t, funcoesEsperadas, funcoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})


	t.Run("erro ao executar a query de busca", func(t *testing.T) {
		dbErr := errors.New("falha na consulta")

		mock.ExpectQuery(query).WillReturnError(dbErr)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoBuscarTodasAsFuncoes)
		require.Empty(t, funcoesDB) // ou require.Len(t, funcoesDB, 0)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados durante a iteracao", func(t *testing.T) {
		// A segunda linha tem um tipo de dado errado para causar falha no Scan
		rows := sqlmock.NewRows([]string{"id", "funcao"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao).
			AddRow("id-invalido", funcoesEsperadas[1].Funcao)

		mock.ExpectQuery(query).WillReturnRows(rows)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrFalhaAoEscanearDados)
		require.Nil(t, funcoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("erro apos a iteracao (linhas.Err)", func(t *testing.T) {
		iterErr := errors.New("erro durante a iteracao")
		
		rows := sqlmock.NewRows([]string{"id", "funcao"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao).
			CloseError(iterErr) // Simula um erro ao fechar as linhas

		mock.ExpectQuery(query).WillReturnRows(rows)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoIterarSobreFuncoes)
		require.Nil(t, funcoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_DeletarFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	idParaDeletar := 1
	query := regexp.QuoteMeta(`delete from funcao where id = @id`)

	t.Run("sucesso ao deletar uma funcao", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 0 para lastInsertId, 1 para rowsAffected

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - funcao nao encontrada para deletar", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectExec(query).
			WithArgs(idNaoExistente).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 linhas afetadas

		err := repo.DeletarFuncao(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAoProcurarFuncao)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico do banco de dados ao deletar", func(t *testing.T) {
		dbErr := errors.New("erro de execucao")

		mock.ExpectExec(query).WithArgs(idParaDeletar).WillReturnError(dbErr)

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.Error(t, err)
		require.Equal(t, dbErr, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao obter linhas afetadas", func(t *testing.T) {
		// Alguns drivers podem não suportar RowsAffected e retornar um erro
		driverErr := errors.New("driver: RowsAffected not supported")
		
		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewErrorResult(driverErr))

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrLinhasAfetadas)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}