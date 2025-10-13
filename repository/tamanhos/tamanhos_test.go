package tamanhos

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
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
		require.ErrorIs(t, err, Errors.ErrSalvar, "erro tem que ser tipo salvar")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha generica ao adicionar", func(t *testing.T) {

		mock.ExpectExec(query).
			WithArgs(tamanho.Tamanho).
			WillReturnError(Errors.ErrSalvar)

		err := repo.AddTamanhos(ctx, &tamanho)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrSalvar, "erro tem que ser do tipo salvar")
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
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado, "erro tem que ser do tipo não encontrado")
		require.Nil(t, tamanhoDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear os dados", func(t *testing.T) {
		// Simulando erro de scan retornando a coluna errada (apenas uma coluna)
		rows := sqlmock.NewRows([]string{"id"}).AddRow(tamanho.ID)

		mock.ExpectQuery(query).WithArgs(tamanho.ID).WillReturnRows(rows)

		tamanhoDB, err := repo.BuscarTamanhos(ctx, tamanho.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados, "erro tem que ser do tipo escanear")
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
		require.ErrorIs(t, err, Errors.ErrBuscarTodos)
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
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados)
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
		require.ErrorIs(t, err, Errors.ErrAoIterar)
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
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado, "erro tem que ser do tipo nao encontrado")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro generico do banco de dados ao deletar", func(t *testing.T) {
		

		mock.ExpectExec(query).WithArgs(idParaDeletar).WillReturnError(Errors.ErrInternal)

		err := repo.DeletarTamanhos(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrInternal) // O erro é retornado diretamente
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao obter linhas afetadas", func(t *testing.T) {
		

		mock.ExpectExec(query).
			WithArgs(idParaDeletar).
			WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		err := repo.DeletarTamanhos(ctx, idParaDeletar)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrLinhasAfetadas, "erro tem que ser do tipo linhas afetadas")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}


func Test_BuscarTamanhosPorIdEpi(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Assumindo que seu repositório é instanciado assim
	repo := NewTamanhoRepository(db) // <-- Adapte esta linha para o seu construtor

	// Query exata que a função usa
	query := regexp.QuoteMeta(`
        select 
            t.id, t.tamanho
        from
            tamanho t
        inner join
            tamanhosEpis te on t.id = te.id_tamanho
        where
            te.epiId = @epiId
    `)

	// Dados de exemplo para o caso de sucesso
	tamanhosEsperados := []model.Tamanhos{
		{ID: 2, Tamanho: "M"},
		{ID: 3, Tamanho: "G"},
	}

	t.Run("sucesso ao buscar os tamanhos de um epi", func(t *testing.T) {
		epiId := 1

		// Prepara as linhas que o mock deve retornar
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(tamanhosEsperados[0].ID, tamanhosEsperados[0].Tamanho).
			AddRow(tamanhosEsperados[1].ID, tamanhosEsperados[1].Tamanho)

		// Define a expectativa: a query será executada com o epiId=1 e retornará as linhas acima
		mock.ExpectQuery(query).WithArgs(epiId).WillReturnRows(rows)

		// Executa a função
		tamanhosDB, err := repo.BuscarTamanhosPorIdEpi(ctx, epiId)

		// Faz as asserções
		require.NoError(t, err)
		require.NotNil(t, tamanhosDB)
		require.Len(t, tamanhosDB, 2)
		require.Equal(t, tamanhosEsperados, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sucesso - epi sem tamanhos associados", func(t *testing.T) {
		epiId := 2

		// Prepara um resultado vazio (apenas as colunas, sem linhas de dados)
		rows := sqlmock.NewRows([]string{"id", "tamanho"})

		mock.ExpectQuery(query).WithArgs(epiId).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTamanhosPorIdEpi(ctx, epiId)

		// Um resultado vazio não é um erro. O slice deve vir vazio.
		require.NoError(t, err)
		require.Empty(t, tamanhosDB)   // Verifica se o slice está vazio
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar a query de busca", func(t *testing.T) {
		epiId := 3
		dbErr := errors.New("falha de conexão")

		// Simula um erro do banco de dados na execução da query
		mock.ExpectQuery(query).WithArgs(epiId).WillReturnError(dbErr)

		tamanhosDB, err := repo.BuscarTamanhosPorIdEpi(ctx, epiId)

		require.Error(t, err)
		// Verifica se o erro retornado é o erro customizado esperado
		require.ErrorIs(t, err, Errors.ErrBuscarTodos, "erro tem que ser do tipo buscar todos")
		require.Nil(t, tamanhosDB) // Em caso de erro, o slice deve ser nulo
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados durante a iteracao", func(t *testing.T) {
		epiId := 4
		// A segunda linha tem um tipo de dado errado (string no ID) para forçar um erro no Scan
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(1, "P").
			AddRow("id-invalido", "M")

		mock.ExpectQuery(query).WithArgs(epiId).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTamanhosPorIdEpi(ctx, epiId)

		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados, "erro tem que ser do tipo escanear")
		require.Nil(t, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro apos a iteracao (linhas.Err)", func(t *testing.T) {
		epiId := 5
		iterErr := errors.New("erro de rede durante a iteracao")

		// Simula um erro que acontece durante a iteração, capturado por `linhas.Err()`
		rows := sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(1, "P").
			CloseError(iterErr)

		mock.ExpectQuery(query).WithArgs(epiId).WillReturnRows(rows)

		tamanhosDB, err := repo.BuscarTamanhosPorIdEpi(ctx, epiId)

		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrAoIterar, "erro tem que ser do tipo iterar")
		require.Nil(t, tamanhosDB)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}