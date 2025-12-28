package funcao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/require"
)

func Test_AddFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	funcao := model.FuncaoInserir{Funcao: "Desenvolvedor"}
	query := regexp.QuoteMeta("insert into")

	t.Run("sucesso ao adicionar uma funcao", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(funcao.Funcao, funcao.IdDepartamento).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddFuncao(ctx, &funcao)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - funcao ja existente", func(t *testing.T) {
		// Simula o erro específico do MS SQL Server para violação de chave única

		errMock := fmt.Errorf("mssql: number 2627, message: duplicate key")
		mock.ExpectExec(query).
			WithArgs(funcao.Funcao, funcao.IdDepartamento).
			WillReturnError(errMock)

		err := repo.AddFuncao(ctx, &funcao)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrSalvar, "erro tem que ser do tipo salvar")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha generica do banco de dados", func(t *testing.T) {
		dbErr := errors.New("falha de conexão")

		mock.ExpectExec(query).
			WithArgs(funcao.Funcao, funcao.IdDepartamento).
			WillReturnError(dbErr)

		err := repo.AddFuncao(ctx, &funcao)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrInternal, "erro tem que ser do tipo internal") // Neste caso, o erro é repassado diretamente
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarFuncao(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewfuncaoRepository(db)
	funcao := model.Funcao{ID: 1, Funcao: "Analista de QA", IdDepartamento: 1, NomeDepartamento: "rh"}
	query := regexp.QuoteMeta(`select `)

	t.Run("sucesso ao buscar uma funcao", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "funcao", "idDepartamento", "nomeDepartamento"}).AddRow(funcao.ID, funcao.Funcao, funcao.IdDepartamento, funcao.NomeDepartamento)

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
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado, "erro tem que ser do tipo nao encontrado")
		require.Nil(t, funcaoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear os dados da funcao", func(t *testing.T) {
		// Simulando erro de scan retornando colunas erradas (faltando 'funcao')
		rows := sqlmock.NewRows([]string{"id"}).AddRow(funcao.ID)

		mock.ExpectQuery(query).WithArgs(funcao.ID).WillReturnRows(rows)

		funcaoDB, err := repo.BuscarFuncao(ctx, funcao.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados, "erro tem que ser do tipo escanear")
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
	query := regexp.QuoteMeta(`select `)
	funcoesEsperadas := []model.Funcao{
		{ID: 1, Funcao: "Gerente", IdDepartamento: 2, NomeDepartamento: "dev"},
		{ID: 2, Funcao: "Coordenador", IdDepartamento: 2, NomeDepartamento: "dev"},
	}

	t.Run("sucesso ao buscar todas as funcoes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "funcao", "idDepartamento", "nomeDepartamenrto"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao, funcoesEsperadas[0].IdDepartamento, funcoesEsperadas[0].NomeDepartamento).
			AddRow(funcoesEsperadas[1].ID, funcoesEsperadas[1].Funcao, funcoesEsperadas[1].IdDepartamento, funcoesEsperadas[1].NomeDepartamento)

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
		require.ErrorIs(t, err, Errors.ErrBuscarTodos)
		require.Empty(t, funcoesDB) // ou require.Len(t, funcoesDB, 0)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados durante a iteracao", func(t *testing.T) {
		// A segunda linha tem um tipo de dado errado para causar falha no Scan
		rows := sqlmock.NewRows([]string{"id", "funcao", "idDepartamento"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao, funcoesEsperadas[0].IdDepartamento).
			AddRow("id-invalido", funcoesEsperadas[1].Funcao, funcoesEsperadas[1].IdDepartamento)

		mock.ExpectQuery(query).WillReturnRows(rows)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados)
		require.Nil(t, funcoesDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro apos a iteracao (linhas.Err)", func(t *testing.T) {
		iterErr := errors.New("erro durante a iteracao")

		rows := sqlmock.NewRows([]string{"id", "funcao", "idDepartamento", "nomeDepartamento"}).
			AddRow(funcoesEsperadas[0].ID, funcoesEsperadas[0].Funcao, funcoesEsperadas[1].IdDepartamento, funcoesEsperadas[1].NomeDepartamento).
			CloseError(iterErr) // Simula um erro ao fechar as linhas

		mock.ExpectQuery(query).WillReturnRows(rows)

		funcoesDB, err := repo.BuscarTodasFuncao(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrAoIterar)
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

	t.Run("sucesso ao deletar uma funcao", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("update ")).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 0 para lastInsertId, 1 para rowsAffected

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - funcao nao encontrada para deletar", func(t *testing.T) {
		idNaoExistente := 99

		mock.ExpectExec(regexp.QuoteMeta("update")).
			WithArgs(idNaoExistente).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 linhas afetadas

		err := repo.DeletarFuncao(ctx, idNaoExistente)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico do banco de dados ao deletar", func(t *testing.T) {
		dbErr := errors.New("erro de execucao")

		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(idParaDeletar).WillReturnError(dbErr)

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.Error(t, err)
		require.Equal(t, dbErr, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao obter linhas afetadas", func(t *testing.T) {
		// Alguns drivers podem não suportar RowsAffected e retornar um erro

		mock.ExpectExec(regexp.QuoteMeta("update ")).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		err := repo.DeletarFuncao(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrLinhasAfetadas)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
