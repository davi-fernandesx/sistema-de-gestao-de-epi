// epi_repository_test.go
package epi

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors" // Seus erros customizados
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"         // Seus models
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// --- Estrutura da Suíte de Testes (Boa prática para organizar) ---

type EpiRepositorySuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo EpiInterface
}

// SetupTest é executado antes de cada teste na suíte
func (s *EpiRepositorySuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.repo = NewEpiRepository(s.db)
}

// TearDownTest é executado depois de cada teste
func (s *EpiRepositorySuite) TearDownTest() {
	s.db.Close()
}

// Função para rodar a suíte
func TestEpiRepository(t *testing.T) {
	suite.Run(t, new(EpiRepositorySuite))
}

// --- Testes para AddEpi ---

func (s *EpiRepositorySuite) TestAddEpi_Success() {
	// Dados de entrada para o teste
	epiToInsert := &model.EpiInserir{
		Nome:           "Capacete Pro-Safety",
		Fabricante:     "MSA",
		CA:             "12345",
		Descricao:      "Capacete de segurança para obras",
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		Idtamanho:      []int{1, 2}, // IDs dos tamanhos a serem associados
		IDprotecao:     10,
		AlertaMinimo:   5,
	}

	// sqlmock espera que essas operações aconteçam NA ORDEM
	s.mock.ExpectBegin() // 1. Espera o início de uma transação

	// 2. Espera a query de INSERT na tabela 'epi'.
	// Usamos regexp.QuoteMeta para tratar a query como texto literal, evitando problemas com caracteres especiais.
	s.mock.ExpectQuery(regexp.QuoteMeta(`insert into epi`)).
		WithArgs(
			epiToInsert.Nome, epiToInsert.Fabricante, epiToInsert.CA, epiToInsert.Descricao, epiToInsert.DataValidadeCa,
			epiToInsert.IDprotecao, epiToInsert.AlertaMinimo,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simula o retorno do ID 1

	// 3. Espera que a query de INSERT para a tabela de associação seja preparada
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanhos_epis"))

	// 4. Espera a execução para o primeiro tamanho (ID 1)
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).
		WithArgs(1, int64(1)). // id_tamanho, id_epi
		WillReturnResult(sqlmock.NewResult(1,1))

	// 5. Espera a execução para o segundo tamanho (ID 2)
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).
		WithArgs(1, int64(2)). // id_tamanho, id_epi
		WillReturnResult(sqlmock.NewResult(1,2))

	// 6. Espera o Commit da transação
	s.mock.ExpectCommit()

	// Executa a função
	err := s.repo.AddEpi(context.Background(), epiToInsert)

	// Verifica se não houve erro
	require.NoError(s.T(), err)
	// Verifica se todas as expectativas do mock foram cumpridas
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *EpiRepositorySuite) TestAddEpi_CommitError() {
	epiToInsert := &model.EpiInserir{Idtamanho: []int{1}} // Dados mínimos para o teste

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`insert into epi`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanhos_epis"))
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Simula um erro no Commit
	s.mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
	
	// A transação deve ser desfeita (Rollback)
	s.mock.ExpectRollback()

	err := s.repo.AddEpi(context.Background(), epiToInsert)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "erro ao commitar transação")

}

// --- Testes para BuscarEpi ---

func TestBuscarEpi(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Inicializa o repositório (ajuste o nome do pacote/struct conforme seu projeto)
	repo := NewEpiRepository(db) 
	ctx := context.Background()

	// Dados de teste
	id := 1
	dataValidade := time.Now()
	
	// Colunas esperadas na primeira query (Dados do EPI)
	colsEpi := []string{
		"id", "nome", "fabricante", "CA", "descricao",
		"validade_CA", "alerta_minimo", "IdTipoProtecao", "nome da protecao",
	}

	// Colunas esperadas na segunda query (Tamanhos)
	colsTamanhos := []string{"id", "tamanho"}

	t.Run("sucesso - deve retornar epi com seus tamanhos", func(t *testing.T) {
		// 1. Mock da primeira query (Dados do EPI)
		// Usamos regexp.QuoteMeta ou um regex genérico para ignorar espaços/quebras de linha
		 // Simplificado para matching
		
		rowEpi := sqlmock.NewRows(colsEpi).AddRow(
			id, "Luva Latex", "Master", "12345", "Luva de proteção",
			dataValidade, 10, 2, "Proteção Química",
		)

		// O mock espera a query do EPI
		// Nota: Como sua query usa sql.Named("id", id), o mock deve validar isso
		mock.ExpectQuery("select"). // Usando regex genérico para facilitar
			WithArgs(sql.Named("id", id)).
			WillReturnRows(rowEpi)

		// 2. Mock da segunda query (Tamanhos do EPI)
		// Essa query roda logo após a primeira ter sucesso
		rowTamanhos := sqlmock.NewRows(colsTamanhos).
			AddRow(10, "P").
			AddRow(11, "M").
			AddRow(12, "G")

		mock.ExpectQuery("select.*tamanho t").
			WithArgs(sql.Named("epiId", id)). // O ID veio do resultado da primeira query
			WillReturnRows(rowTamanhos)

		// Execução
		resultado, err := repo.BuscarEpi(ctx, id)

		// Validações
		require.NoError(t, err)
		require.NotNil(t, resultado)
		assert.Equal(t, "Luva Latex", resultado.Nome)
		assert.Equal(t, "Master", resultado.Fabricante)
		
		// Valida se os tamanhos foram preenchidos
		require.Len(t, resultado.Tamanhos, 3)
		assert.Equal(t, "P", resultado.Tamanhos[0].Tamanho)
		assert.Equal(t, "G", resultado.Tamanhos[2].Tamanho)
	})

	t.Run("erro - epi não encontrado", func(t *testing.T) {
		// Mock retorna erro NoRows na primeira query
		mock.ExpectQuery("select").
			WithArgs(sql.Named("id", id)).
			WillReturnError(sql.ErrNoRows)

		resultado, err := repo.BuscarEpi(ctx, id)

		// Validações
		require.Error(t, err)
		assert.Nil(t, resultado)
		assert.Contains(t, err.Error(), "não encontrado")
		
		// Garante que a segunda query (tamanhos) NÃO foi chamada, pois a primeira falhou
	})

	t.Run("erro - falha ao buscar tamanhos", func(t *testing.T) {
		// 1. Primeira query funciona
		rowEpi := sqlmock.NewRows(colsEpi).AddRow(
			id, "Luva", "Fab", "111", "Desc", time.Now(), 5, 1, "Prot",
		)
		mock.ExpectQuery("select").
			WithArgs(sql.Named("id", id)).
			WillReturnRows(rowEpi)

		// 2. Segunda query falha (ex: banco caiu no meio do processo)
		mock.ExpectQuery("select.*tamanho t").
			WithArgs(sql.Named("epiId", id)).
			WillReturnError(errors.New("erro de conexão"))

		resultado, err := repo.BuscarEpi(ctx, id)

		// Validações
		require.Error(t, err)
		assert.Nil(t, resultado)
		assert.Contains(t, err.Error(), "erro ao buscar tamanhos")
	})
}
// --- Testes para DeletarEpi ---

func (s *EpiRepositorySuite) TestDeletarEpi_Success() {
	epiID := 1

	s.mock.ExpectBegin()

	// Espera o delete na tabela de associação
	s.mock.ExpectExec(regexp.QuoteMeta("delete ")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 2)) // Simula que 2 associações foram deletadas

	// Espera o delete na tabela principal
	s.mock.ExpectExec(regexp.QuoteMeta("delete ")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simula que 1 EPI foi deletado

	s.mock.ExpectCommit()

	err := s.repo.DeletarEpi(context.Background(), epiID)

	require.NoError(s.T(), err)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *EpiRepositorySuite) TestDeletarEpi_NotFound() {
	epiID := 99

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("delete ")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Nenhuma associação encontrada

	s.mock.ExpectExec(regexp.QuoteMeta("delete")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Nenhum EPI encontrado

	s.mock.ExpectRollback() // A transação deve ser desfeita

	err := s.repo.DeletarEpi(context.Background(), epiID)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, Errors.ErrNaoEncontrado)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}