package entregaepi


import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Ajuste o caminho de import para o seu projeto
	
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

// newMock cria um mock de DB para os testes
func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ctx := context.Background()
	return db, mock, ctx
}

func TestAddentrega(t *testing.T) {
	db, mock, ctx := newMock(t)
	defer db.Close()

	repo := NewEntregaRepository(db)

	entregaModel := model.EntregaParaInserir{
		ID_funcionario:     1,
		Data_entrega:       time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC),
		Assinatura_Digital: "base64-string",
		Itens: []model.ItemParaEntrega{
			{ID_epi: 10, ID_tamanho: 2, Quantidade: 1},
			{ID_epi: 11, ID_tamanho: 3, Quantidade: 2},
		},
	}
	const newEntregaID int64 = 1

	// Query CORRIGIDA (com @AssinaturaDigital)
	queryEntrega := regexp.QuoteMeta(`insert into entrega (id_funcionario, data_entrega, AssinaturaDigital)
     values (@idFuncionario, @dataEntrega, @AssinaturaDigital)
     OUTPUT INSERTED.id`)

	queryItem := regexp.QuoteMeta(`insert into epi_entregas(id_epi,id_tamanho, quantidade, id_entrega) values (@id_epi, @id_tamanho, @quantidade, @id_entrega)`)

	t.Run("sucesso ao adicionar entrega com 2 itens", func(t *testing.T) {
		mock.ExpectBegin() // 1. Espera a transação começar

		// 2. Espera a query de inserção da entrega "pai"
		mock.ExpectQuery(queryEntrega).
			WithArgs(
				sql.Named("idFuncionario", entregaModel.ID_funcionario),
				sql.Named("dataEntrega", entregaModel.Data_entrega),
				sql.Named("AssinaturaDigital", entregaModel.Assinatura_Digital),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEntregaID))

		// 3. Espera o "Prepare" da query dos itens
		mock.ExpectPrepare(queryItem)

		// 4. Espera a execução para o Item 1
		mock.ExpectExec(queryItem). // O nome do "prepare" no SQL Server
						WithArgs(
				sql.Named("id_epi", entregaModel.Itens[0].ID_epi),
				sql.Named("id_tamanho", entregaModel.Itens[0].ID_tamanho),
				sql.Named("quantidade", entregaModel.Itens[0].Quantidade),
				sql.Named("id_entrega", newEntregaID),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 5. Espera a execução para o Item 2
		mock.ExpectExec(queryItem).
			WithArgs(
				sql.Named("id_epi", entregaModel.Itens[1].ID_epi),
				sql.Named("id_tamanho", entregaModel.Itens[1].ID_tamanho),
				sql.Named("quantidade", entregaModel.Itens[1].Quantidade),
				sql.Named("id_entrega", newEntregaID),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 6. Espera o Commit final
		mock.ExpectCommit()

		err := repo.Addentrega(ctx, entregaModel)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao iniciar transacao", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("db error"))

		err := repo.Addentrega(ctx, entregaModel)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao inserir entrega (pai)", func(t *testing.T) {
		dbError := errors.New("insert error")
		mock.ExpectBegin()
		mock.ExpectQuery(queryEntrega).WillReturnError(dbError)
		mock.ExpectRollback() // Espera o rollback automático do defer

		err := repo.Addentrega(ctx, entregaModel)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrInternal) // Verifica o erro customizado
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao preparar query de itens", func(t *testing.T) {
		dbError := errors.New("prepare error")
		mock.ExpectBegin()
		mock.ExpectQuery(queryEntrega).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEntregaID))
		mock.ExpectPrepare(queryItem).WillReturnError(dbError)
		mock.ExpectRollback()

		err := repo.Addentrega(ctx, entregaModel)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrInternal)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao inserir um item (filho)", func(t *testing.T) {
		dbError := errors.New("item exec error")
		mock.ExpectBegin()
		mock.ExpectQuery(queryEntrega).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEntregaID))
		mock.ExpectPrepare(queryItem)
		mock.ExpectExec(queryItem).WillReturnError(dbError) // Falha no primeiro item
		mock.ExpectRollback()

		err := repo.Addentrega(ctx, entregaModel)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrInternal)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao comitar transacao", func(t *testing.T) {
		dbError := errors.New("commit error")
		mock.ExpectBegin()
		mock.ExpectQuery(queryEntrega).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEntregaID))
		mock.ExpectPrepare(queryItem)
		mock.ExpectExec(queryItem). // Item 1 OK
						WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec(queryItem). // Item 2 OK
						WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit().WillReturnError(dbError) // Commit falha
		

		err := repo.Addentrega(ctx, entregaModel)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrInternal)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

// Colunas esperadas pelas queries de busca
var colunasBusca = []string{
	"id", "data_Entrega", "id_funcionario", "nome", "id_departamento",
	"departamento", "id_funcao", "funcao", "id_epi", "nome_epi",
	"fabricante", "CA", "descricao", "data_fabricacao", "data_validade",
	"data_validadeCa", "id_tipo_protecao", "protecao", "id_tamanho",
	"tamanho", "quantidade", "AssinaturaDigital",
}

// Dados de mock para uma linha de item
var mockItem1 = []driver.Value{
	1, time.Now(), 10, "Funcionario Teste", 2, "Produção", 3, "Operador",
	101, "Capacete", "Marca X", "CA123", "Desc", time.Now(), time.Now(), time.Now(),
	5, "Cabeça", 7, "Único", 1, "sig-base64",
}
var mockItem2 = []driver.Value{
	1, time.Now(), 10, "Funcionario Teste", 2, "Produção", 3, "Operador",
	102, "Luva", "Marca Y", "CA456", "Desc Luva", time.Now(), time.Now(), time.Now(),
	6, "Mãos", 8, "M", 2, "sig-base64",
}
var mockItem3_Entrega2 = []driver.Value{
	2, time.Now(), 11, "Outro Func", 2, "Produção", 4, "Ajudante",
	103, "Bota", "Marca Z", "CA789", "Desc Bota", time.Now(), time.Now(), time.Now(),
	7, "Pés", 9, "42", 1, "sig-base64-2",
}

func TestBuscaEntrega(t *testing.T) {
	db, mock, ctx := newMock(t)
	defer db.Close()

	repo := NewEntregaRepository(db)

	// Query exata do seu código
	query := regexp.QuoteMeta(`select
            ee.id,
            ee.data_Entrega,
            ee.id_funcionario,
            f.nome, 
            f.id_departamento, 
            d.departamento, 
            f.id_funcao, 
            ff.funcao, 
            i.id_epi, 
            e.nome, 
            e.fabricante, 
            e.CA,
            e.descricao, 
            e.data_fabricacao, 
            e.data_validade, 
            e.data_validadeCa,
            e.id_tipo_protecao,
            tp.protecao, 
            i.id_tamanho, 
            t.tamanho, 
            i.quantidade,
            ee.AssinaturaDigital
            from entrega ee
            inner join
                funcionario f on ee.id_funcionario = f.id
            inner join
                departamentos d on f.id_departamento = d.id
            inner join 
                funcao ff on f.id_funcao = ff.id
            inner join 
                epi_entregues i on i.id_entrega = ee.id
            inner join 
                epi e on i.id_epi = e.id
            inner join
                tipo_protecao tp on e.id_tipo_protecao = tp.id
            inner join 
                tamanho t on i.id_tamanho = t.id
            where ee.cancelada_em IS NULL and ee.id = @id`)

	t.Run("sucesso ao buscar entrega com 2 itens", func(t *testing.T) {
		rows := sqlmock.NewRows(colunasBusca).
			AddRow(mockItem1...).
			AddRow(mockItem2...)

		mock.ExpectQuery(query).WithArgs(sql.Named("id", 1)).WillReturnRows(rows)

		entrega, err := repo.BuscaEntrega(ctx, 1)

		require.NoError(t, err)
		require.NotNil(t, entrega)
		assert.Equal(t, 1, entrega.Id)
		assert.Equal(t, "Funcionario Teste", entrega.Funcionario.Nome)
		assert.Equal(t, 2, len(entrega.Itens))
		assert.Equal(t, "Capacete", entrega.Itens[0].Epi.Nome)
		assert.Equal(t, "Luva", entrega.Itens[1].Epi.Nome)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - entrega nao encontrada", func(t *testing.T) {
		rows := sqlmock.NewRows(colunasBusca) // Retorna 0 linhas

		mock.ExpectQuery(query).WithArgs(sql.Named("id", 1)).WillReturnRows(rows)

		entrega, err := repo.BuscaEntrega(ctx, 1)

		// Agora que seu código retorna um erro, o teste espera por ele
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrNaoEncontrado, "Erro deve ser 'Não Encontrado'")
		assert.Nil(t, entrega, "Entrega deve ser nil")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar query no banco", func(t *testing.T) {
		dbError := errors.New("db error")
		mock.ExpectQuery(query).WithArgs(sql.Named("id", 1)).WillReturnError(dbError)

		entrega, err := repo.BuscaEntrega(ctx, 1)

		require.Error(t, err)
		assert.Nil(t, entrega)
		assert.ErrorIs(t, err, Errors.ErrBuscarTodos)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao escanear dados (scan error)", func(t *testing.T) {
		// Retorna ID como string, o que vai quebrar o .Scan()
		rows := sqlmock.NewRows(colunasBusca).
			AddRow("ID-INVALIDO", time.Now(), 10, "Func Teste", 2, "Prod", 3, "Op",
				101, "Capacete", "Marca X", "CA123", "Desc", time.Now(), time.Now(), time.Now(),
				5, "Cabeça", 7, "Único", 1, "sig-base64")

		mock.ExpectQuery(query).WithArgs(sql.Named("id", 1)).WillReturnRows(rows)

		entrega, err := repo.BuscaEntrega(ctx, 1)

		require.Error(t, err)
		assert.Nil(t, entrega)
		assert.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuscaTodasEntregas(t *testing.T) {
	db, mock, ctx := newMock(t)
	defer db.Close()

	repo := NewEntregaRepository(db)

	// Query CORRIGIDA (com 'where' no lugar certo)
	query := regexp.QuoteMeta(`select
            ee.id,
            ee.data_Entrega,
            ee.id_funcionario,
            f.nome, 
            f.id_departamento, 
            d.departamento, 
            f.id_funcao, 
            ff.funcao, 
            i.id_epi, 
            e.nome, 
            e.fabricante, 
            e.CA,
            e.descricao, 
            e.data_fabricacao, 
            e.data_validade, 
            e.data_validadeCa,
            e.id_tipo_protecao,
            tp.protecao, 
            i.id_tamanho, 
            t.tamanho, 
            i.quantidade,
            ee.AssinaturaDigital
            from entrega ee
            inner join
                funcionario f on ee.id_funcionario = f.id
            inner join
                departamentos d on f.id_departamento = d.id
            inner join 
                funcao ff on f.id_funcao = ff.id
            inner join 
                epi_entregues i on i.id_entrega = ee.id
            inner join 
                epi e on i.id_epi = e.id
            inner join
                tipo_protecao tp on e.id_tipo_protecao = tp.id
            inner join 
                tamanho t on i.id_tamanho = t.id
            where
                 ee.cancelada_em IS NULL
            ORDER BY ee.id`)

	t.Run("sucesso ao buscar todas as entregas (2 entregas, 3 itens)", func(t *testing.T) {
		rows := sqlmock.NewRows(colunasBusca).
			AddRow(mockItem1...).        // Entrega 1, Item 1
			AddRow(mockItem2...).        // Entrega 1, Item 2
			AddRow(mockItem3_Entrega2...) // Entrega 2, Item 1

		mock.ExpectQuery(query).WillReturnRows(rows)

		entregas, err := repo.BuscaTodasEntregas(ctx)

		require.NoError(t, err)
		require.NotNil(t, entregas)
		assert.Equal(t, 2, len(entregas), "Deveria haver 2 entregas no total")

		// Como o resultado final vem de um map, a ordem do slice pode variar
		// Precisamos encontrar a entrega correta para testar
		entrega1 := findEntregaByID(entregas, 1)
		entrega2 := findEntregaByID(entregas, 2)

		require.NotNil(t, entrega1)
		require.NotNil(t, entrega2)

		assert.Equal(t, 2, len(entrega1.Itens), "Entrega 1 deveria ter 2 itens")
		assert.Equal(t, 1, len(entrega2.Itens), "Entrega 2 deveria ter 1 item")
		assert.Equal(t, "Funcionario Teste", entrega1.Funcionario.Nome)
		assert.Equal(t, "Outro Func", entrega2.Funcionario.Nome)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sucesso - nenhuma entrega encontrada", func(t *testing.T) {
		rows := sqlmock.NewRows(colunasBusca) // 0 linhas
		mock.ExpectQuery(query).WillReturnRows(rows)

		entregas, err := repo.BuscaTodasEntregas(ctx)

		require.NoError(t, err)
		require.NotNil(t, entregas, "O slice não deve ser nil, deve ser vazio")
		assert.Equal(t, 0, len(entregas), "O slice de entregas deve estar vazio")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar query no banco", func(t *testing.T) {
		dbError := errors.New("db error")
		mock.ExpectQuery(query).WillReturnError(dbError)

		entregas, err := repo.BuscaTodasEntregas(ctx)

		require.Error(t, err)
		assert.Nil(t, entregas)
		assert.ErrorIs(t, err, Errors.ErrBuscarTodos)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

// Helper para TestBuscaTodasEntregas
func findEntregaByID(entregas []*model.EntregaDto, id int) *model.EntregaDto {
	for _, e := range entregas {
		if e.Id == id {
			return e
		}
	}
	return nil
}


func TestCancelarEntrega(t *testing.T) {
	db, mock, ctx := newMock(t)
	defer db.Close()

	repo :=NewEntregaRepository(db)
	const mockEntregaID = 1

	// Query CORRIGIDA
	query := regexp.QuoteMeta(`update entrega
            set cancelada_em  = GETDATE() 
            where id = @id AND cancelada_em IS NULL;`)

	t.Run("sucesso ao cancelar entrega", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", mockEntregaID)).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada

		err := repo.CancelarEntrega(ctx, mockEntregaID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - entrega nao encontrada (0 linhas afetadas)", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", mockEntregaID)).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 linhas afetadas

		err := repo.CancelarEntrega(ctx, mockEntregaID)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrNaoEncontrado)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao executar query no banco", func(t *testing.T) {
		dbError := errors.New("db error")
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", mockEntregaID)).
			WillReturnError(dbError)

		err := repo.CancelarEntrega(ctx, mockEntregaID)
		require.Error(t, err)
		assert.ErrorIs(t, err, Errors.ErrInternal)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao checar linhas afetadas", func(t *testing.T) {
		dbError := Errors.ErrLinhasAfetadas // Simula um erro específico
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", mockEntregaID)).
			WillReturnResult(sqlmock.NewErrorResult(dbError)) // Erro ao chamar RowsAffected()

		err := repo.CancelarEntrega(ctx, mockEntregaID)
		require.Error(t, err)
		// Seu código checa por 'ErrLinhasAfetadas'
		assert.ErrorIs(t, err, Errors.ErrLinhasAfetadas)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}