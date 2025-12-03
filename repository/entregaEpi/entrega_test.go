package entregaepi

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
func (m *MockEstoqueRepo) BuscaDevoluvao(ctx context.Context, id int) (*model.Devolucao, error) {
	panic("unimplemented")
}

// BuscaTodasDevolucoe implements trocaepi.DevolucaoInterfaceRepository.
func (m *MockEstoqueRepo) BuscaTodasDevolucoe(ctx context.Context) ([]model.Devolucao, error) {
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
	mockSQL.ExpectQuery(regexp.QuoteMeta(`insert into entrega (id_funcionario, data_entrega, AssinaturaDigital)`)).
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

// --- 3. TESTE: BUSCA ENTREGA (POR ID) ---

func TestBuscaEntrega_Sucesso(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	repo := NewEntregaRepository(db, nil) // Não precisa do mock de estoque aqui

	// Colunas retornadas pela query gigante
	cols := []string{
		"id", "data_Entrega", "id_funcionario", "nome", "id_departamento", "departamento",
		"id_funcao", "funcao", "id_epi", "nome_epi", "fabricante", "CA", "descricao",
		"data_fabricacao", "data_validade", "data_validadeCa", "id_tipo_protecao", "protecao",
		"id_tamanho", "tamanho", "quantidade", "AssinaturaDigital", "valorUnitario",
	}

	// Simula 2 linhas retornadas (1 entrega com 2 itens)
	rows := sqlmock.NewRows(cols).
		AddRow(1, time.Now(), 10, "Joao", 1, "TI", 2, "Dev", 50, "Luva", "FabX", "123", "Desc", time.Now(), time.Now(), time.Now(), 1, "Mao", 1, "G", 2, "ass123", 10.0).
		AddRow(1, time.Now(), 10, "Joao", 1, "TI", 2, "Dev", 51, "Bota", "FabY", "456", "Desc", time.Now(), time.Now(), time.Now(), 2, "Pe", 2, "40", 1, "ass123", 50.0)

	mockSQL.ExpectQuery(regexp.QuoteMeta(`select ee.id, ee.data_Entrega`)). // Match parcial da query
										WithArgs(sql.Named("id", 1)).
										WillReturnRows(rows)

	dto, err := repo.BuscaEntrega(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, 1, dto.Id)
	assert.Equal(t, "Joao", dto.Funcionario.Nome)
	assert.Len(t, dto.Itens, 2) // Deve ter agrupado os 2 itens
	assert.Equal(t, "Luva", dto.Itens[0].Epi.Nome)
	assert.Equal(t, "Bota", dto.Itens[1].Epi.Nome)
}

func TestBuscaEntrega_NaoEncontrado(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	repo := NewEntregaRepository(db, nil)

	// Retorna 0 linhas
	mockSQL.ExpectQuery(regexp.QuoteMeta(`select ee.id`)).
		WithArgs(sql.Named("id", 999)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	dto, err := repo.BuscaEntrega(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, dto)
	assert.Equal(t, Errors.ErrNaoEncontrado, err)
}

// --- 4. TESTE: BUSCA TODAS AS ENTREGAS ---

func TestBuscaTodasEntregas_Sucesso(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()
	defer db.Close()
	repo := NewEntregaRepository(db, nil)

	cols := []string{
		"id", "data_Entrega", "id_funcionario", "nome", "id_departamento", "departamento",
		"id_funcao", "funcao", "id_epi", "nome_epi", "fabricante", "CA", "descricao",
		"data_fabricacao", "data_validade", "data_validadeCa", "id_tipo_protecao", "protecao",
		"id_tamanho", "tamanho", "quantidade", "AssinaturaDigital", "valorUnitario",
	}

	// Simula 3 linhas: Entrega 1 (2 itens) e Entrega 2 (1 item)
	rows := sqlmock.NewRows(cols).
		AddRow(1, time.Now(), 10, "Joao", 1, "TI", 2, "Dev", 50, "Luva", "FabX", "123", "Desc", time.Now(), time.Now(), time.Now(), 1, "Mao", 1, "G", 2, "ass1", 10.0).
		AddRow(1, time.Now(), 10, "Joao", 1, "TI", 2, "Dev", 51, "Bota", "FabY", "456", "Desc", time.Now(), time.Now(), time.Now(), 2, "Pe", 2, "40", 1, "ass1", 50.0).
		AddRow(2, time.Now(), 20, "Maria", 1, "RH", 3, "Gestor", 60, "Oculos", "FabZ", "789", "Desc", time.Now(), time.Now(), time.Now(), 3, "Olho", 3, "U", 1, "ass2", 20.0)

	mockSQL.ExpectQuery(regexp.QuoteMeta(`select ee.id`)).
		WillReturnRows(rows)

	lista, err := repo.BuscaTodasEntregas(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, lista)
	assert.Len(t, lista, 2) // Deve ter agrupado em 2 entregas distintas

	// Verificando se os itens foram distribuidos corretamente
	// Nota: Como o mapa em Go não garante ordem, temos que achar o ID correto
	var entrega1, entrega2 *model.EntregaDto
	for _, e := range lista {
		switch e.Id {
case 1:
			entrega1 = e
		case 2:
			entrega2 = e
		}
	}

	assert.NotNil(t, entrega1)
	assert.Len(t, entrega1.Itens, 2)
	assert.NotNil(t, entrega2)
	assert.Len(t, entrega2.Itens, 1)
}

// --- 5. TESTE: CANCELAR ENTREGA ---

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
