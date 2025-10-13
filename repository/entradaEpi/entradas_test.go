package entradaepi

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func mock(t *testing.T)(*sql.DB, sqlmock.Sqlmock,context.Context ,error ){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)

	


	return db, mock,  ctx,err

}

func Test_EntradaEpi(t *testing.T){

	db, mock, ctx, err:= mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo:= NewEntradaRepository(db)

	entradaInserir:= model.EntradaEpiInserir{
		ID_epi: 1,
		Data_entrada: time.Now(),
		Quantidade: 10,
		Lote: "xyz",
		Fornecedor: "teste1",

	}

	query:= regexp.QuoteMeta(`
		insert into Entrada (id_epi, data_entrada, quantidade, lote, fornecedor)
		values (@id_epi, @data_entrada, @quantidade, @lote, @fornecedor)`)


	t.Run("testando o sucesso ao adicionar uma entrada", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(entradaInserir.ID_epi, entradaInserir.Data_entrada, 
		entradaInserir.Quantidade, entradaInserir.Lote, 
		entradaInserir.Fornecedor).WillReturnResult(sqlmock.NewResult(0,1))
		

		err:= repo.AddEntradaEpi(ctx, &entradaInserir)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro ao adicionar uma entrada", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(entradaInserir.ID_epi, entradaInserir.Data_entrada, 
		entradaInserir.Quantidade, entradaInserir.Lote, 
		entradaInserir.Fornecedor).WillReturnError(Errors.ErrSalvar)

		err:= repo.AddEntradaEpi(ctx, &entradaInserir)
		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrSalvar), "erro tem que ser do tipo salvar")
		require.NoError(t, mock.ExpectationsWereMet())

	})
}


func Test_BuscarEntrada(t *testing.T) {
	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	
	repo := NewEntradaRepository(db)

	// Query SQL esperada, escapada para uso com regexp
	query := regexp.QuoteMeta(`
            select
                ee.id, ee.id_epi,e.nome,e.fabricante, e.CA,e.descricao,
                e.dataFabricacao, e.dataValidade, e.dataValidadeCa, 
                e.id_protecao, e.nomeProtecao, ee.lote, ee.fornecedor 
            from 
                entrada ee
            inner join
                epi e on ee.id_epi = e.id 
            where 
                ee.id = @id
    `)

	idParaBuscar := 1
	agora := time.Now()
	entradaMock := model.EntradaEpi{
		ID:             1,
		ID_epi:         10,
		Nome:           "Capacete de Segurança",
		Fabricante:     "Marca Segura",
		CA:             "12345",
		Descricao:      "Capacete para proteção contra impactos.",
		DataFabricacao: agora.AddDate(0, -6, 0), // 6 meses atrás
		DataValidade:   agora.AddDate(2, 0, 0),  // Daqui a 2 anos
		DataValidadeCa: agora.AddDate(1, 0, 0),  // Daqui a 1 ano
		IDprotecao:     1,
		NomeProtecao:   "Cabeça",
		Lote:           "LOTE-2025-A1",
		Fornecedor:     "Fornecedor Principal",
	}


	t.Run("sucesso ao buscar entrada", func(t *testing.T) {
		// Define as colunas que a query retorna, na ordem correta
		rows := sqlmock.NewRows([]string{
			"id", "id_epi", "nome", "fabricante", "CA", "descricao",
			"dataFabricacao", "dataValidade", "dataValidadeCa", 
			"id_protecao", "nomeProtecao", "lote", "fornecedor",
		}).AddRow( // Adiciona uma linha com os dados do nosso mock
			entradaMock.ID,
			entradaMock.ID_epi,
			entradaMock.Nome,
			entradaMock.Fabricante,
			entradaMock.CA,
			entradaMock.Descricao,
			entradaMock.DataFabricacao,
			entradaMock.DataValidade,
			entradaMock.DataValidadeCa,
			entradaMock.IDprotecao,
			entradaMock.NomeProtecao,
			entradaMock.Lote,
			entradaMock.Fornecedor,
		)

		mock.ExpectQuery(query).
			WithArgs(entradaMock.ID). // Com o argumento correto
			WillReturnRows(rows)                    

	
		entradaRetornada, err := repo.BuscarEntrada(ctx, idParaBuscar)

		// Asserções
		require.NoError(t, err)
		require.NotNil(t, entradaRetornada)// Verifica se o resultado é o esperado
		require.NoError(t, mock.ExpectationsWereMet())  
	})

	t.Run("erro - entrada nao encontrada", func(t *testing.T) {
		// Esperamos uma query, mas desta vez ela retornará um erro `sql.ErrNoRows`
		mock.ExpectQuery(query).
			WithArgs(sql.Named("id", idParaBuscar)).
			WillReturnError(sql.ErrNoRows)

		// Executa a função
		entradaRetornada, err := repo.BuscarEntrada(ctx, idParaBuscar)

		// Asserções
		require.Error(t, err)
		require.Nil(t, entradaRetornada) // O objeto de retorno deve ser nulo
		// Verifica se o erro retornado é o erro customizado correto
		assert.True(t, errors.Is(err, Errors.ErrNaoEncontrado), "erro tem que ser do tipo nao encontrado")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha ao escanear dados", func(t *testing.T) {
		// Para simular um erro de scan, podemos retornar colunas com tipos errados.
		// Por exemplo, uma string onde deveria ser um int.
		rows := sqlmock.NewRows([]string{
			"id", "id_epi", "nome", "fabricante", "CA", "descricao",
			"dataFabricacao", "dataValidade", "dataValidadeCa", 
			"id_protecao", "nomeProtecao", "lote", "fornecedor",
		}).AddRow( // O primeiro campo `id` deveria ser `int`, mas passamos uma `string`
			"id_invalido", 
			entradaMock.ID_epi,
			entradaMock.Nome,
			entradaMock.Fabricante,
			entradaMock.CA,
			entradaMock.Descricao,
			entradaMock.DataFabricacao,
			entradaMock.DataValidade,
			entradaMock.DataValidadeCa,
			entradaMock.IDprotecao,
			entradaMock.NomeProtecao,
			entradaMock.Lote,
			entradaMock.Fornecedor,
		)

		// Esperamos a query...
		mock.ExpectQuery(query).
			WithArgs(sql.Named("id", idParaBuscar)).
			WillReturnRows(rows) // ... que retornará os dados "corrompidos"

		// Executa a função
		entradaRetornada, err := repo.BuscarEntrada(ctx, idParaBuscar)

		// Asserções
		require.Error(t, err)
		require.Nil(t, entradaRetornada)
		// Verifica se o erro retornado é o erro customizado para falha de scan
		assert.True(t, errors.Is(err, Errors.ErrFalhaAoEscanearDados), "erro tem que ser do tipo escanear")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_BuscarTodasEntradas(t *testing.T) {
	// --- Setup ---
	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewEntradaRepository(db)

	// Query SQL esperada para buscar TODAS as entradas (sem a cláusula WHERE)
	query := regexp.QuoteMeta(`
			select
				ee.id, ee.id_epi,e.nome,e.fabricante, e.CA,e.descricao,
				e.dataFabricacao, e.dataValidade, e.dataValidadeCa, 
				e.id_protecao, e.nomeProtecao, ee.lote, ee.fornecedor 
			from 
				entrada ee
			inner join
				epi e on ee.id_epi = e.id 
    `) // Removido o "where ee.id = @id" para alinhar com o nome da função

	// Dados de exemplo que esperamos que o banco retorne
	agora := time.Now()
	entradasMock := []model.EntradaEpi{
		{
			ID:             1,
			ID_epi:         10,
			Nome:           "Capacete de Segurança",
			Fabricante:     "Marca Segura",
			CA:             "12345",
			Descricao:      "Capacete para proteção contra impactos.",
			DataFabricacao: agora.AddDate(0, -6, 0),
			DataValidade:   agora.AddDate(2, 0, 0),
			DataValidadeCa: agora.AddDate(1, 0, 0),
			IDprotecao:     1,
			NomeProtecao:   "Cabeça",
			Lote:           "LOTE-2025-A1",
			Fornecedor:     "Fornecedor Principal",
		},
		{
			ID:             2,
			ID_epi:         20,
			Nome:           "Luva de Proteção",
			Fabricante:     "Marca Tátil",
			CA:             "67890",
			Descricao:      "Luva para proteção química.",
			DataFabricacao: agora.AddDate(0, -3, 0),
			DataValidade:   agora.AddDate(1, 6, 0),
			DataValidadeCa: agora.AddDate(1, 0, 0),
			IDprotecao:     2,
			NomeProtecao:   "Mãos",
			Lote:           "LOTE-2024-B2",
			Fornecedor:     "Fornecedor Secundário",
		},
	}

	// --- Test Cases ---

	t.Run("sucesso ao buscar todas as entradas", func(t *testing.T) {
		// Define as colunas que a query retorna, na ordem correta
		rows := sqlmock.NewRows([]string{
			"id", "id_epi", "nome", "fabricante", "CA", "descricao",
			"dataFabricacao", "dataValidade", "dataValidadeCa",
			"id_protecao", "nomeProtecao", "lote", "fornecedor",
		})

		// Adiciona os dados do nosso mock nas linhas
		for _, entrada := range entradasMock {
			rows.AddRow(
				entrada.ID, entrada.ID_epi, entrada.Nome, entrada.Fabricante, entrada.CA,
				entrada.Descricao, entrada.DataFabricacao, entrada.DataValidade,
				entrada.DataValidadeCa, entrada.IDprotecao, entrada.NomeProtecao,
				entrada.Lote, entrada.Fornecedor,
			)
		}

		mock.ExpectQuery(query).WillReturnRows(rows)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.NoError(t, err)
		require.NotNil(t, resultado)
		require.Len(t, resultado, 2) // Verifica se retornou 2 entradas
		require.Equal(t, entradasMock, resultado)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("sucesso - nenhuma entrada encontrada", func(t *testing.T) {
		// Retorna as colunas, mas nenhuma linha
		rows := sqlmock.NewRows([]string{
			"id", "id_epi", "nome", "fabricante", "CA", "descricao",
			"dataFabricacao", "dataValidade", "dataValidadeCa",
			"id_protecao", "nomeProtecao", "lote", "fornecedor",
		})

		mock.ExpectQuery(query).WillReturnRows(rows)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.NoError(t, err)
		require.NotNil(t, resultado)      // O slice não deve ser nulo
		require.Len(t, resultado, 0)      // O slice deve estar vazio
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha na execucao da query", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(Errors.ErrBuscarTodos)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrBuscarTodos), "erro tem que ser do tipo buscar todos")
		require.Len(t, resultado, 0) // Garante que o slice retornado é vazio
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - falha ao escanear dados de uma linha", func(t *testing.T) {
		// Retorna uma linha com um tipo de dado incompatível (e.g., string para ID)
		rows := sqlmock.NewRows([]string{"id", "id_epi", "nome"}).
			AddRow("id_invalido", 10, "Capacete")

		mock.ExpectQuery(query).WillReturnRows(rows)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrFalhaAoEscanearDados), "erro tem que ser do tipo escanear")
		require.Nil(t, resultado) // A função retorna nil em caso de erro no scan
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro - ao iterar sobre as linhas", func(t *testing.T) {
	
		entradasMock :=model.EntradaEpi{
			ID:             1,
			ID_epi:         10,
			Nome:           "Capacete de Segurança",
			Fabricante:     "Marca Segura",
			CA:             "12345",
			Descricao:      "Capacete para proteção contra impactos.",
			DataFabricacao: agora.AddDate(0, -6, 0),
			DataValidade:   agora.AddDate(2, 0, 0),
			DataValidadeCa: agora.AddDate(1, 0, 0),
			IDprotecao:     1,
			NomeProtecao:   "Cabeça",
			Lote:           "LOTE-2025-A1",
			Fornecedor:     "Fornecedor Principal",
		}


		entradasMock2:=model.EntradaEpi{
			ID:             2,
			ID_epi:         20,
			Nome:           "Luva de Proteção",
			Fabricante:     "Marca Tátil",
			CA:             "67890",
			Descricao:      "Luva para proteção química.",
			DataFabricacao: agora.AddDate(0, -3, 0),
			DataValidade:   agora.AddDate(1, 6, 0),
			DataValidadeCa: agora.AddDate(1, 0, 0),
			IDprotecao:     2,
			NomeProtecao:   "Mãos",
			Lote:           "LOTE-2024-B2",
			Fornecedor:     "Fornecedor Secundário",
		}
		linhas:= sqlmock.NewRows([]string{
			"id", 
			"id_epi",
		    "nome",
			"fabricante", 
			"CA", 
			"descricao",
			"dataFabricacao", 
			"dataValidade", 
			"dataValidadeCa",
			"id_protecao", 
			"nomeProtecao", 
			"lote", 
			"fornecedor",
		}).AddRow(
				entradasMock.ID, 
				entradasMock.ID_epi, 
				entradasMock.Nome, 
				entradasMock.Fabricante, 
				entradasMock.CA,
				entradasMock.Descricao, 
				entradasMock.DataFabricacao, 
				entradasMock.DataValidade,
				entradasMock.DataValidadeCa, 
				entradasMock.IDprotecao, 
				entradasMock.NomeProtecao,
				entradasMock.Lote, 
				entradasMock.Fornecedor,

		).AddRow(

				entradasMock2.ID, 
				entradasMock2.ID_epi, 
				entradasMock2.Nome, 
				entradasMock2.Fabricante, 
				entradasMock2.CA,
				entradasMock2.Descricao, 
				entradasMock2.DataFabricacao, 
				entradasMock2.DataValidade,
				entradasMock2.DataValidadeCa, 
				entradasMock2.IDprotecao, 
				entradasMock2.NomeProtecao,
				entradasMock2.Lote, 
				entradasMock2.Fornecedor,

		).CloseError(Errors.ErrDadoIncompativel)

		// Simulamos um erro que ocorre após a leitura bem-sucedida das linhas (verificado por `linhas.Err()`)
		mock.ExpectQuery(query).WillReturnRows(linhas)

		resultado, err := repo.BuscarTodasEntradas(ctx)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrAoIterar), "erro tem que ser do tipo iterar")
		require.Nil(t, resultado)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}


func TestDeleteEntrada(t *testing.T){

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewEntradaRepository(db)

	entradasMock :=model.EntradaEpi{
			ID:             1,
			ID_epi:         10,
			Nome:           "Capacete de Segurança",
			Fabricante:     "Marca Segura",
			CA:             "12345",
			Descricao:      "Capacete para proteção contra impactos.",
			DataFabricacao: time.Now().AddDate(0, -6, 0),
			DataValidade:   time.Now().AddDate(2, 0, 0),
			DataValidadeCa: time.Now().AddDate(1, 0, 0),
			IDprotecao:     1,
			NomeProtecao:   "Cabeça",
			Lote:           "LOTE-2025-A1",
			Fornecedor:     "Fornecedor Principal",
		}

	query := regexp.QuoteMeta(`delete from  entrada where id = @id`)

t.Run("testando o sucesso ao deletar um epi da base de dados", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewResult(0, 1))

		errEpi := repo.DeletarEntrada(ctx, entradasMock.ID)
		require.NoError(t, errEpi)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao deletar uma entradas", func(t *testing.T) {

		ErroGenericoDb := errors.New("erro ao se conectar com o banco")
		mock.ExpectExec(query).WithArgs(entradasMock.ID).WillReturnError(ErroGenericoDb)

		errEpi := repo.DeletarEntrada(ctx, entradasMock.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrInternal), "erro tem que ser do tipo internal")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de linhas afetadas", func(t *testing.T) {


		mock.ExpectExec(query).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		errEpi := repo.DeletarEntrada(ctx, entradasMock.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrLinhasAfetadas), "erro tem que ser do tipo linhas afetadas")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("epi nao encontrado", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(entradasMock.ID).WillReturnResult(sqlmock.NewResult(0, 0))

		errEpi := repo.DeletarEntrada(ctx, entradasMock.ID)
		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrNaoEncontrado), "erro tem que ser do tipo nao encontrado")

		require.NoError(t, mock.ExpectationsWereMet())
	})
}