package entregaepi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	// Ajuste os imports conforme o nome do seu módulo
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	// Importe onde está a interface DevolucaoInterfaceRepository
	// "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/trocaEpi"
)

// --- 1. MOCK DA DEPENDÊNCIA DE ESTOQUE ---
// Precisamos simular o comportamento da BaixaEstoque sem chamar o código real
type MockEstoqueRepo struct {
	mock.Mock
}

// BuscaDevoluvao implements trocaepi.DevolucaoInterfaceRepository.
func (m *MockEstoqueRepo) BuscaDevolucao(ctx context.Context, id int) ([]model.Devolucao, error) {
	panic("unimplemented")
}

// BuscaTodasDevolucoe implements trocaepi.DevolucaoInterfaceRepository.
func (m *MockEstoqueRepo) BuscaTodasDevolucoes(ctx context.Context) ([]model.Devolucao, error) {
	panic("unimplemented")
}

// DeleteDevolucao implements trocaepi.DevolucaoInterfaceRepository.
func (m *MockEstoqueRepo) DeleteDevolucao(ctx context.Context, id int) error {
	panic("unimplemented")
}

// Simulando a função BaixaEstoque
func (m *MockEstoqueRepo) BaixaEstoque(ctx context.Context, tx *sql.Tx, idEpi, idTamanho int64, quantidade int, idEntrega int64) error {
	args := m.Called(ctx, tx, idEpi, idTamanho, quantidade, idEntrega)
	return args.Error(0)
}

// Métodos que não são usados neste teste, mas a interface exige (stubs)
func (m *MockEstoqueRepo) AddDevolucaoEpi(ctx context.Context, devolucao model.DevolucaoInserir) error {
	return nil
}
func (m *MockEstoqueRepo) AddTrocaEPI(ctx context.Context, devolucao model.DevolucaoInserir) error {
	return nil
}

// --- 2. TESTE: ADD ENTREGA ---

func TestAddEntrega_Sucesso(t *testing.T) {
	db, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockEstoque := new(MockEstoqueRepo)
	repo := NewEntregaRepository(db, mockEstoque)

	// Dados de Entrada
	entregaInput := model.EntregaParaInserir{
		ID_funcionario:     1,
		Data_entrega:       time.Now(),
		Assinatura_Digital: "hash_assinatura",
		Itens: []model.ItemParaInserir{
			{ID_epi: 10, ID_tamanho: 2, Quantidade: 5},
		},
	}

	// 1. Expectativa: Iniciar Transação
	mockSQL.ExpectBegin()

	// 2. Expectativa: Insert na Tabela Entrega
	// Usamos QuoteMeta para escapar caracteres especiais do SQL e AnyArg para data
	mockSQL.ExpectQuery(regexp.QuoteMeta("insert ")).
		WithArgs(entregaInput.ID_funcionario, sqlmock.AnyArg(), entregaInput.Assinatura_Digital).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100)) // Retorna ID 100

	// 3. Expectativa: Chamar BaixaEstoque (Mock da Interface)
	// Esperamos que o AddEntrega chame o BaixaEstoque passando o ID 100 gerado acima
	mockEstoque.On("BaixaEstoque", mock.Anything, mock.Anything, int64(10), int64(2), 5, int64(100)).
		Return(nil) // Retorna sucesso

	// 4. Expectativa: Commit
	mockSQL.ExpectCommit()

	// Ação
	err = repo.Addentrega(context.Background(), entregaInput)

	// Validação
	assert.NoError(t, err)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockEstoque.AssertExpectations(t)
}

func TestAddEntrega_ErroNoInsert(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	mockEstoque := new(MockEstoqueRepo)
	repo := NewEntregaRepository(db, mockEstoque)

	entregaInput := model.EntregaParaInserir{ID_funcionario: 1}

	mockSQL.ExpectBegin()
	// Simula erro no banco ao inserir
	mockSQL.ExpectQuery(regexp.QuoteMeta(`insert into entrega`)).
		WillReturnError(errors.New("erro de conexao"))
	mockSQL.ExpectRollback()

	err := repo.Addentrega(context.Background(), entregaInput)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro interno") // Verifica se a mensagem de erro bate com Errors.ErrInternal
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockEstoque.AssertNotCalled(t, "BaixaEstoque") // Garante que não tentou baixar estoque se o insert falhou
}

func TestAddEntrega_ErroNaBaixaDeEstoque(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	mockEstoque := new(MockEstoqueRepo)
	repo := NewEntregaRepository(db, mockEstoque)

	entregaInput := model.EntregaParaInserir{
		ID_funcionario: 1,
		Itens:          []model.ItemParaInserir{{ID_epi: 1, Quantidade: 1}},
	}

	mockSQL.ExpectBegin()
	mockSQL.ExpectQuery(regexp.QuoteMeta(`insert into entrega`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	// Simula erro na interface de estoque
	mockEstoque.On("BaixaEstoque", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("estoque insuficiente"))

	// Deve fazer Rollback se o estoque falhar
	mockSQL.ExpectRollback()

	err := repo.Addentrega(context.Background(), entregaInput)

	assert.Error(t, err)
	assert.Equal(t, "estoque insuficiente", err.Error())
	assert.NoError(t, mockSQL.ExpectationsWereMet())
}

var entregaColunas = []string{

	"id", "dataEntrega", "id_funcionario", "nome", "id_departamento", "departamento", "id_funcao", "funcao",
	"id_epi", "nome", "fabricante", "CA", "descricao",  "data_validadeCA",
	"id_tipo_protecao", "protecao", "id_tamanho", "tamanho", "quantidade", "assinatura_digital", "valorUnitario",
}

func TestBuscaEntregaPorId(t *testing.T) {

	ctx := context.Background()
	mockEstoque := new(MockEstoqueRepo)
	db, mock, err := sqlmock.New()
	if err != nil {

		t.Fatal(err)
	}

	dataEntrega := time.Now()
	dataValidadeCa := time.Now()

	id1 := 1
	id2 := 2
	repo := NewEntregaRepository(db, mockEstoque)

	row := sqlmock.NewRows(entregaColunas).AddRow(

		id1, dataEntrega, 3, "davi", 4, "ti", 3, "dev", 1, "luva", "master", "64556", " luvas de borracha",
		 dataValidadeCa, 2, "maos", 2, "G", 2, "hash", 12.99,
	)

	row2 := sqlmock.NewRows(entregaColunas).AddRow(

		id2, dataEntrega, 4, "rada", 4, "ti", 3, "dev", 1, "luva", "master", "64556", " luvas de borracha",
		 dataValidadeCa, 2, "maos", 2, "G", 2, "hash", 12.99,
	)

	t.Run("sucesso ao retorna entrega por id", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("id", id1)).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("id", id2)).WillReturnRows(row2)

		resultado1, err := repo.BuscaEntrega(ctx, id1)

		require.NoError(t, err)
		require.NotNil(t, resultado1)
		require.Equal(t, "davi", resultado1.Funcionario.Nome)
		require.Equal(t, id1, resultado1.Id)

		resultado2, err := repo.BuscaEntrega(ctx, id2)
		require.NoError(t, err)
		require.NotNil(t, resultado2)
		require.Equal(t, "rada", resultado2.Funcionario.Nome)
		require.Equal(t, id2, resultado2.Id)

	})

	t.Run("erro no banco de dados e em resultados nao encontrados", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select")).WithArgs(sql.Named("id", id1)).WillReturnError(sql.ErrConnDone)

		mock.ExpectQuery(regexp.QuoteMeta("select")).WithArgs(sql.Named("id", id2)).WillReturnError(errors.New("erro inesperado do banco"))

		id3 := 3
		rowVazia := sqlmock.NewRows(entregaColunas)
		mock.ExpectQuery(regexp.QuoteMeta("select")).WithArgs(sql.Named("id", id3)).WillReturnRows(rowVazia)

		resultado1, err1 := repo.BuscaEntrega(ctx, id1)

		require.Error(t, err1)
		require.Nil(t, resultado1)
		require.True(t, errors.Is(err1, Errors.ErrBuscarTodos))

		resultado2, err2 := repo.BuscaEntrega(ctx, id2)

		require.Error(t, err2)
		fmt.Println(err2)
		require.Nil(t, resultado2)
		require.True(t, errors.Is(err2, Errors.ErrBuscarTodos))

		resultado3, err3 := repo.BuscaEntrega(ctx, id3)

		if err3 != nil {
			t.Logf("\n🔴 ERRO NO CENÁRIO 3: %v\n", err3)
		}
		require.NoError(t, err3)
		require.Empty(t, resultado3)

	})

}

func TestCancelarEntrega_Sucesso(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	repo := NewEntregaRepository(db, nil)

	mockSQL.ExpectExec(regexp.QuoteMeta(`update entrega set cancelada_em = GETDATE()`)).
		WithArgs(sql.Named("id", 10)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 Linha afetada

	err := repo.CancelarEntrega(context.Background(), 10)

	assert.NoError(t, err)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
}

func TestCancelarEntrega_NaoEncontrada(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	repo := NewEntregaRepository(db, nil)

	mockSQL.ExpectExec(regexp.QuoteMeta(`update entrega`)).
		WithArgs(sql.Named("id", 999)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 Linhas afetadas

	err := repo.CancelarEntrega(context.Background(), 999)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, Errors.ErrNaoEncontrado)) // Verifica se o erro "wrapper" contém o erro alvo
}
