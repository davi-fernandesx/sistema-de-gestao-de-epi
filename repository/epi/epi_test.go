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
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"  // Seus models
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
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(5, 0, 0),
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
			epiToInsert.Nome, epiToInsert.Fabricante, epiToInsert.CA, epiToInsert.Descricao,
			epiToInsert.DataFabricacao, epiToInsert.DataValidade, epiToInsert.DataValidadeCa,
			epiToInsert.IDprotecao, epiToInsert.AlertaMinimo,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simula o retorno do ID 1

	// 3. Espera que a query de INSERT para a tabela de associação seja preparada
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanho_epi"))

	// 4. Espera a execução para o primeiro tamanho (ID 1)
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanho_epi")).
		WithArgs(1, int64(1)). // id_tamanho, id_epi
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 5. Espera a execução para o segundo tamanho (ID 2)
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanho_epi")).
		WithArgs(2, int64(1)). // id_tamanho, id_epi
		WillReturnResult(sqlmock.NewResult(2, 1))

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
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanho_epi"))
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanho_epi")).WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Simula um erro no Commit
	s.mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
	
	// A transação deve ser desfeita (Rollback)
	s.mock.ExpectRollback()

	err := s.repo.AddEpi(context.Background(), epiToInsert)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "erro ao commitar transação")

}

// --- Testes para BuscarEpi ---

func (s *EpiRepositorySuite) TestBuscarEpi_Success() {
	epiID := 1
	
	// Mock da primeira query (buscar dados do EPI)
	epiRows := sqlmock.NewRows([]string{"id", 
	"nome", 
	"fabricante", 
	"CA", 
	"descricao",
	"data_fabricacao",
	"data_validade", 
	"validade_CA",   
	"alerta_minimo", 
	"id_tipo_protecao", 
	"nome_protecao"}).
		AddRow(epiID, "Capacete", "MSA", "123", "Desc", time.Now(), time.Now(), time.Now(), 5, 10, "Proteção Craniana")
	
	s.mock.ExpectQuery(regexp.QuoteMeta(`select e.id, e.nome`)).
		WithArgs(epiID).
		WillReturnRows(epiRows)

	// Mock da segunda query (buscar os tamanhos)
	tamanhoRows := sqlmock.NewRows([]string{"id", "tamanho"}).
		AddRow(1, "P").
		AddRow(2, "M")

	// CORREÇÃO: Os nomes das tabelas e colunas precisam ser consistentes
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			select 
			t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhosEpis te on t.id = te.id_tamanho
		where
			te.epiId = @epiId`)).
		WithArgs(epiID).
		WillReturnRows(tamanhoRows)

	epi, err := s.repo.BuscarEpi(context.Background(), epiID)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), epi)
	require.Equal(s.T(), epiID, epi.ID)
	require.Equal(s.T(), "Capacete", epi.Nome)
	require.Len(s.T(), epi.Tamanhos, 2) // Verifica se os 2 tamanhos foram adicionados
	require.Equal(s.T(), "P", epi.Tamanhos[0].Tamanho)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}


func (s *EpiRepositorySuite) TestBuscarEpi_NotFound() {
	epiID := 99
	
	s.mock.ExpectQuery(regexp.QuoteMeta(`select e.id, e.nome`)).
		WithArgs(epiID).
		WillReturnError(sql.ErrNoRows) // Simula que o banco não encontrou nada

	epi, err := s.repo.BuscarEpi(context.Background(), epiID)

	require.Error(s.T(), err)
	require.Nil(s.T(), epi)
	require.ErrorIs(s.T(), err, Errors.ErrNaoEncontrado) // Verifica se o erro customizado foi retornado
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}


// --- Testes para BuscarTodosEpi ---

func (s *EpiRepositorySuite) TestBuscarTodosEpi_Success() {
	// Mock da primeira query (buscar todos os EPIs)
	epiRows := sqlmock.NewRows([]string{"id", "nome", "fabricante", "CA", "descricao", "data_fabricacao", "data_validade", "validade_CA", "alerta_minimo", "id_tipo_protecao", "nome_protecao"}).
		AddRow(1, "Capacete", "MSA", "123", "Desc1", time.Now(), time.Now(), time.Now(), 5, 10, "Proteção Craniana").
		AddRow(2, "Luva", "3M", "456", "Desc2", time.Now(), time.Now(), time.Now(), 10, 11, "Proteção Mãos")

	s.mock.ExpectQuery(regexp.QuoteMeta(`select e.id, e.nome`)).
		WillReturnRows(epiRows)

	// Mock da segunda query (buscar todas as associações de tamanho)
	tamanhoRows := sqlmock.NewRows([]string{"id_epi", "id", "tamanho"}).
		AddRow(1, 1, "P").      // Tamanho P para o Capacete
		AddRow(1, 2, "M").      // Tamanho M para o Capacete
		AddRow(2, 3, "Único") // Tamanho Único para a Luva

	// CORREÇÃO: Corrigido o "inner joinn" e nomes de colunas
	s.mock.ExpectQuery(regexp.QuoteMeta(`select te.id_epi, t.id, t.tamanho from tamanhos t inner join tamanho_epi te`)).
		WillReturnRows(tamanhoRows)

	epis, err := s.repo.BuscarTodosEpi(context.Background())

	require.NoError(s.T(), err)
	require.Len(s.T(), epis, 2)

	// Valida o primeiro EPI
	require.Equal(s.T(), 1, epis[0].ID)
	require.Len(s.T(), epis[0].Tamanhos, 2)

	// Valida o segundo EPI
	require.Equal(s.T(), 2, epis[1].ID)
	require.Len(s.T(), epis[1].Tamanhos, 1)
	require.Equal(s.T(), "Único", epis[1].Tamanhos[0].Tamanho)

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// --- Testes para DeletarEpi ---

func (s *EpiRepositorySuite) TestDeletarEpi_Success() {
	epiID := 1

	s.mock.ExpectBegin()

	// Espera o delete na tabela de associação
	s.mock.ExpectExec(regexp.QuoteMeta("delete from tamanho_epi where id_epi = @id")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 2)) // Simula que 2 associações foram deletadas

	// Espera o delete na tabela principal
	s.mock.ExpectExec(regexp.QuoteMeta("delete from epi where id = @id")).
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
	s.mock.ExpectExec(regexp.QuoteMeta("delete from tamanho_epi")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Nenhuma associação encontrada

	s.mock.ExpectExec(regexp.QuoteMeta("delete from epi")).
		WithArgs(epiID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // Nenhum EPI encontrado

	s.mock.ExpectRollback() // A transação deve ser desfeita

	err := s.repo.DeletarEpi(context.Background(), epiID)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, Errors.ErrNaoEncontrado)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}