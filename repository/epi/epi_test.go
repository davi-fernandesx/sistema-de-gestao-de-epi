package epi

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
	"github.com/stretchr/testify/suite"
)

// Struct para simular erros do MSSQL sem importar o driver
type mssqlErrorMock struct {
	number int32
}

func (e mssqlErrorMock) Number() int32 { return e.number }
func (e mssqlErrorMock) Error() string { return fmt.Sprintf("mssql: number %d", e.number) }

type EpiRepositorySuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
	repo *EpiRepository
}

func (s *EpiRepositorySuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.NoError(err)
	s.db = db
	s.mock = mock
	s.repo = NewEpiRepository(db)
}

func (s *EpiRepositorySuite) TearDownTest() {
	s.db.Close()
}

func TestEpiRepository(t *testing.T) {
	suite.Run(t, new(EpiRepositorySuite))
}

// --- TESTES DO ADDEPI ---

func (s *EpiRepositorySuite) TestAddEpi_Success() {
	ctx := context.Background()
	epi := &model.EpiInserir{
		Nome:           "Luva Nitrílica",
		Fabricante:     "Danny",
		CA:             "12345",
		Descricao:      "Luva de proteção",
		DataValidadeCa:  *configs.NewDataBrPtr(time.Now()),
		IDprotecao:     1,
		AlertaMinimo:   10,
		Idtamanho:      []int{1, 2},
	}

	s.mock.ExpectBegin()

	// Mock do Insert Principal (EPI)
	s.mock.ExpectQuery(regexp.QuoteMeta(`insert into epi`)).
		WithArgs(
			sqlmock.AnyArg(), // nome
			sqlmock.AnyArg(), // fabricante
			sqlmock.AnyArg(), // CA
			sqlmock.AnyArg(), // descricao
			sqlmock.AnyArg(), // validade_CA
			sqlmock.AnyArg(), // id_tipo_protecao
			sqlmock.AnyArg(), // alerta_minimo
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	// Mock do Prepare e Execução dos Tamanhos
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanhos_epis"))
	
	// Espera dois inserts (um para cada ID no slice Idtamanho)
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).
		WithArgs(int64(100), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).
		WithArgs(int64(100), 2).WillReturnResult(sqlmock.NewResult(2, 1))

	s.mock.ExpectCommit()

	err := s.repo.AddEpi(ctx, epi)
	s.NoError(err)
}

func (s *EpiRepositorySuite) TestAddEpi_CA_Duplicado() {
    ctx := context.Background()
    epi := &model.EpiInserir{CA: "12345"}

    s.mock.ExpectBegin()
    
    errMock := mssqlErrorMock{number: 2627}
    
    // Use AnyArg() para garantir que o mock capture a query 
    // independente dos valores nulos/vazios da struct
    s.mock.ExpectQuery(regexp.QuoteMeta(`insert into epi`)).
        WithArgs(
            sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
            sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
            sqlmock.AnyArg(),
        ).WillReturnError(errMock)
        
    s.mock.ExpectRollback()

    err := s.repo.AddEpi(ctx, epi)

    s.Error(err)
	
    s.True(errors.Is(err, Errors.ErrSalvar))
   // s.Contains(err.Error(), "CA 12345 ja existe no sistema")
}

func (s *EpiRepositorySuite) TestAddEpi_CommitError() {
	epi := &model.EpiInserir{Idtamanho: []int{1}}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`insert into epi`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectPrepare(regexp.QuoteMeta("insert into tamanhos_epis"))
	s.mock.ExpectExec(regexp.QuoteMeta("insert into tamanhos_epis")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Simula falha no commit
	s.mock.ExpectCommit().WillReturnError(errors.New("falha no commit"))
	s.mock.ExpectRollback() // Chamado pelo defer

	err := s.repo.AddEpi(context.Background(), epi)

	s.Error(err)
	s.Contains(err.Error(), "erro ao commitar transação")
}

// --- TESTES DO BUSCAREPI ---

func (s *EpiRepositorySuite) TestBuscarEpi_Success() {
	id := 1
	s.mock.ExpectQuery(regexp.QuoteMeta("select e.id, e.nome")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "fabricante", "CA", "descricao", "validade_CA", "alerta_minimo", "IdTipoProtecao", "nome da protecao"}).
			AddRow(1, "EPI Teste", "Fab", "123", "Desc", time.Now(), 5, 1, "Proteção Mãos"))

	// Mock da busca de tamanhos
	s.mock.ExpectQuery(regexp.QuoteMeta("select t.id, t.tamanho")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tamanho"}).
			AddRow(1, "G").
			AddRow(2, "M"))

	res, err := s.repo.BuscarEpi(context.Background(), id)

	s.NoError(err)
	s.NotNil(res)
	s.Equal("EPI Teste", res.Nome)
	s.Len(res.Tamanhos, 2)
}

func (s *EpiRepositorySuite) TestBuscarEpi_NotFound() {
	s.mock.ExpectQuery("").WillReturnError(sql.ErrNoRows)

	res, err := s.repo.BuscarEpi(context.Background(), 99)

	s.Error(err)
	s.Nil(res)
	s.True(errors.Is(err, Errors.ErrNaoEncontrado))
}

// --- TESTES DE DELETAR ---

func (s *EpiRepositorySuite) TestDeletarEpi_Success() {
	id := 10
	s.mock.ExpectBegin()
	s.mock.ExpectExec("update tamanhos_epis").WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 2))
	s.mock.ExpectExec("update epi").WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada
	s.mock.ExpectCommit()

	err := s.repo.DeletarEpi(context.Background(), id)
	s.NoError(err)
}

func (s *EpiRepositorySuite) TestDeletarEpi_NotFound() {
	id := 99
	s.mock.ExpectBegin()
	s.mock.ExpectExec("update tamanhos_epis").WillReturnResult(sqlmock.NewResult(0, 0))
	s.mock.ExpectExec("update epi").WillReturnResult(sqlmock.NewResult(0, 0)) // 0 linhas
	s.mock.ExpectRollback()

	err := s.repo.DeletarEpi(context.Background(), id)

	s.Error(err)
	s.True(errors.Is(err, Errors.ErrNaoEncontrado))
}