package trocaepi

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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

var colunaDevolucao = []string{

	"id", "idFucionario", "nome", "idDepartamernto", "departamento", "idFuncao", "funcao", "id_epi",
	"nome", "fabricante", "CA", "id_tamanho", "tamanho", "quantidadeAdevolver", "motivo", "id_epiNovo",
	"EpiNovo", "fabricanteEpiNovo", "CAEpiNovo", "quantidadeNova", "id_tamanhoNovo", "tamanhoNovo", "assinaturaDigital",
	"dataTroca",
}

func mock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context, error) {

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {

		t.Fatal("Erro ao iniciar mock")
	}

	return db, mock, ctx, err

}

func TestBuscaDevoluicaoPorMatricula(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewDevolucaoRepository(db)

	matricula := 1234
	dataTroca := time.Now()

	t.Run("sucesso da devolucao", func(t *testing.T) {

		rows := sqlmock.NewRows(colunaDevolucao).AddRow(

			1, 12, "rada", 1, "ti", 2, "programador", 1, "luva", "master", "23445", 2, "G",
			1, "desgaste", sql.NullInt64{Int64: 12, Valid: true},
			"luva", "lupo", "12312", 1, 2, "M", "hash", dataTroca,
		)

		mock.ExpectQuery(regexp.QuoteMeta("select")).WithArgs(sql.Named("matricula", matricula)).WillReturnRows(rows)

		resultado, err := repo.BuscaDevolucaoPorMatricula(ctx, matricula)

		require.NoError(t, err)
		require.NotNil(t, resultado)
		require.Len(t, resultado, 1)
		require.Equal(t, "rada", resultado[0].NomeFuncionario)
		require.Equal(t, "luva", resultado[0].NomeEpiTroca)

	})

	t.Run("deve retornar um erro", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("select ")).WithArgs(sql.Named("matricula", matricula)).WillReturnError(sql.ErrConnDone)

		resultado, err := repo.BuscaDevolucaoPorMatricula(ctx, matricula)

		require.Error(t, err)
		require.Nil(t, resultado)

		if err := mock.ExpectationsWereMet(); err != nil {

			t.Errorf("algumas expectativas nao foram atendidas")
		}
	})
}

func TestBuscaTodasDevolucao(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {

		t.Fatal(err)
	}
	defer db.Close()
	repo := NewDevolucaoRepository(db)

	dataTroca := time.Now()

	t.Run("sucesso ao buscar todas as devolucoes", func(t *testing.T) {
		rows := sqlmock.NewRows(colunaDevolucao).AddRow(

			1, 12, "rada", 1, "ti", 2, "programador", 1, "luva", "master", "23445", 2, "G",
			1, "desgaste", sql.NullInt64{Int64: 12, Valid: true},
			"luva", "lupo", "12312", 1, 2, "M", "hash", dataTroca,
		).AddRow(
			3, 1, "davi", 14, "acougue", 3, "dessosa", 1, "luva de aço", "adidas", "23345", 4, "p",
			3, "perdeu", sql.NullInt64{Int64: 9, Valid: false},
			sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false},
			sql.NullString{String: "", Valid: false},
			sql.NullInt64{Int64: 0, Valid: false}, sql.NullInt64{Int64: 0, Valid: false},
			sql.NullString{String: "", Valid: false}, "hash", dataTroca,
		)

		mock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(rows)

		resultado, err := repo.BuscaTodasDevolucoes(ctx)

		require.NoError(t, err)
		require.Len(t, resultado, 2)
		require.NotNil(t, resultado)
		require.Equal(t, "rada", resultado[0].NomeFuncionario)
		require.Equal(t, "davi", resultado[1].NomeFuncionario)

		t.Run("erro ao buscar todas as devoluçoes", func(t *testing.T) {

			mock.ExpectQuery(regexp.QuoteMeta("select ")).WillReturnError(sql.ErrNoRows)

			resultado, err := repo.BuscaTodasDevolucoes(ctx)

			require.Nil(t, resultado)
			require.Error(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {

				t.Errorf("algumas expectativas não foram atendidas")
			}

		})

	})
}
func TestDevolucaoEpi(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {

		t.Fatal(err)
	}

	defer db.Close()

	repo := NewDevolucaoRepository(db)

	testModel := model.DevolucaoInserir{
		IdFuncionario:       1,
		IdEpi:               1,
		IdMotivo:            2,
		IdTamanho: 			 1,
		DataDevolucao:       time.Now(),	
		QuantidadeADevolver: 1,
		NovaQuantidade:      nil,
		IdEpiNovo:           nil,
		IdTamanhoNovo:       nil,
		AssinaturaDigital:   "hash",
	}

	t.Run("sucesso ao adicionar uma devolucao", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("insert into")).WithArgs(
			sql.Named("idFuncionario", testModel.IdFuncionario),
			sql.Named("idEpi", testModel.IdEpi),
			sql.Named("motivo", testModel.IdMotivo),
			sql.Named("dataDevolucao", testModel.DataDevolucao),
			sql.Named("idTamanho", testModel.IdTamanho),
			sql.Named("quantidadeDevolucao", testModel.QuantidadeADevolver),
			sql.Named("assinaturaDigital", testModel.AssinaturaDigital),
		).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddDevolucaoEpi(ctx, testModel)

		require.NoError(t, err)

	})

	t.Run("erro ao salvar devolucao de epi - erro de conexao", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta("insert into")).WithArgs(
			sql.Named("idFuncionario", testModel.IdFuncionario),
			sql.Named("idEpi", testModel.IdEpi),
			sql.Named("motivo", testModel.IdMotivo),
			sql.Named("dataDevolucao", testModel.DataDevolucao),
			sql.Named("idTamanho", testModel.IdTamanho),
			sql.Named("quantidadeDevolucao", testModel.QuantidadeADevolver),
			sql.Named("assinaturaDigital", testModel.AssinaturaDigital),
		).WillReturnError(sql.ErrConnDone)

		err := repo.AddDevolucaoEpi(ctx, testModel)
		require.Error(t, err)

		require.True(t, errors.Is(err, Errors.ErrInternal))

	})

	if err := mock.ExpectationsWereMet(); err != nil {

		t.Fatal("algumas expectativas não foram atendidas")
	}

}

func TestBaixaEstoque(t *testing.T) {
	db, mock, ctx, err := mock(t)
	if err != nil {

		t.Fatal(err)
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewDevolucaoRepository(db)

	idEpi := 2
	idTamanho := 2
	quantidade := 5
	idEntrega := 1

	//dados de retorno
	idEntrada := 50

	valorUnitario := decimal.NewFromFloat(40.99)
	saldoLote := 10

	t.Run("sucesso da funcao baixa estoque", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "valorUnitario", "quantidade"}).AddRow(idEntrada, valorUnitario, saldoLote)

		mock.ExpectQuery(regexp.QuoteMeta("select top 1")).WithArgs(
			sql.Named("idEpi", idEpi),
			sql.Named("id_tamanho", idTamanho),
			sql.Named("quantidade", quantidade),
		).WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("update ")).WithArgs(
			sql.Named("qtd", quantidade),
			sql.Named("idEntrada", idEntrada),
		).WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta("insert into ")).WithArgs(
			sql.Named("id_epi", idEpi),
			sql.Named("id_tamanho", idTamanho),
			sql.Named("quantidade", quantidade),
			sql.Named("id_entrega", idEntrega),
			sql.Named("id_entrada", idEntrada),
			sql.Named("valorUnitario", valorUnitario),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.BaixaEstoque(ctx, tx, int64(idEpi), int64(idTamanho), quantidade, int64(idEntrega))

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("estoque insuficiente", func(t *testing.T) {

		mock.ExpectBegin()
		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery(regexp.QuoteMeta("select top 1")).WithArgs(
			sql.Named("idEpi", idEpi),
			sql.Named("id_tamanho", idTamanho),
			sql.Named("quantidade", quantidade),
		).WillReturnError(sql.ErrNoRows)

		err2 := repo.BaixaEstoque(ctx, tx, int64(idEpi), int64(idTamanho), quantidade, int64(idEntrega))

		require.Error(t, err2)
		require.True(t, errors.Is(err2, Errors.ErrEstoqueInsuficiente))

	})

}

func TestAddTrocaEPI(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewDevolucaoRepository(db)
	ctx := context.Background()

	// Preparando dados (Ponteiros não podem ser nulos aqui pois o código usa *var)
	var qtdNova int = 2
	var idEpiNovo int = 5
	var idTamNovo int = 3

	devolucao := model.DevolucaoInserir{
		IdFuncionario:       1,
		IdEpi:               10,
		IdMotivo:            2,
		DataDevolucao:       time.Now(),
		QuantidadeADevolver: 1,
		AssinaturaDigital:   "hash_teste",
		// Campos da Troca (Ponteiros)
		NovaQuantidade: &qtdNova,
		IdEpiNovo:      &idEpiNovo,
		IdTamanhoNovo:  &idTamNovo,
	}

	// IDs simulados que o banco retornaria nos OUTPUT INSERTED.id
	idDevolucaoGerado := int64(999)
	idEntregaGerado := int64(888)

	// Dados simulados para o BaixaEstoque
	idLoteEstoque := int64(50)
	valorUnitario := decimal.NewFromFloat(15.99)
	qtdNoLote := 10

	t.Run("sucesso ao realizar troca completa", func(t *testing.T) {

		// 1. Iniciar Transação
		mock.ExpectBegin()

		// 2. Insert Devolucao (Usa QueryRow por causa do OUTPUT INSERTED.id)
		// O mock deve retornar uma LINHA contendo o ID gerado
		rowsDevolucao := sqlmock.NewRows([]string{"id"}).AddRow(idDevolucaoGerado)

		mock.ExpectQuery("(?i)insert into devolucao.*").
			WithArgs(
				sql.Named("idFuncionario", devolucao.IdFuncionario),
				sql.Named("idEpi", devolucao.IdEpi),
				sql.Named("motivo", devolucao.IdMotivo),
				sql.Named("dataDevolucao", devolucao.DataDevolucao),
				sql.Named("idTamanho", devolucao.IdTamanho),
				sql.Named("quantidadeDevolucao", devolucao.QuantidadeADevolver),
				sql.Named("idEpiNovo", devolucao.IdEpiNovo),
				sql.Named("IdtamanhoEpiNovo", devolucao.IdTamanhoNovo),
				sql.Named("quantidadeNova", devolucao.NovaQuantidade),
				sql.Named("assinaturaDigital", devolucao.AssinaturaDigital),
			).
			WillReturnRows(rowsDevolucao)

		// 3. Insert Entrega (Usa QueryRow por causa do OUTPUT INSERTED.id)
		rowsEntrega := sqlmock.NewRows([]string{"id"}).AddRow(idEntregaGerado)

		mock.ExpectQuery("(?i)insert into entrega.*").
			WithArgs(
				sql.Named("idFuncionario", devolucao.IdFuncionario),
				sql.Named("dataEntrega", devolucao.DataDevolucao), // Note que vc usa DataDevolucao aqui no codigo
				sql.Named("assinaturaDigital", devolucao.AssinaturaDigital),
				sql.Named("idTroca", idDevolucaoGerado), // O ID que veio do passo anterior!
			).
			WillReturnRows(rowsEntrega)

		// --- AGORA COMEÇA A SIMULAÇÃO DO 'BaixaEstoque' ---

		// 4. Select do Lote (Dentro do BaixaEstoque)
		rowsLote := sqlmock.NewRows([]string{"id", "valorUnitario", "quantidade"}).
			AddRow(idLoteEstoque, valorUnitario, qtdNoLote)

		mock.ExpectQuery("(?i)select top 1.*").
			WithArgs(
				sql.Named("idEpi", int64(*devolucao.IdEpiNovo)),
				sql.Named("id_tamanho", int64(*devolucao.IdTamanhoNovo)),
				sql.Named("quantidade", *devolucao.NovaQuantidade),
			).
			WillReturnRows(rowsLote)

		// 5. Update do Estoque (Dentro do BaixaEstoque)
		mock.ExpectExec("(?i)update entrada.*").
			WithArgs(
				sql.Named("qtd", *devolucao.NovaQuantidade),
				sql.Named("idEntrada", idLoteEstoque),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 6. Insert na Tabela epi_entregas (Dentro do BaixaEstoque)
		mock.ExpectExec("insert into epis_entregues ").
			WithArgs(
				sql.Named("id_epi", int64(*devolucao.IdEpiNovo)),
				sql.Named("id_tamanho", int64(*devolucao.IdTamanhoNovo)),
				sql.Named("quantidade", *devolucao.NovaQuantidade),
				sql.Named("id_entrega", idEntregaGerado), // ID da entrega gerado no passo 3
				sql.Named("id_entrada", idLoteEstoque),
				sql.Named("valorUnitario", valorUnitario),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// --- FIM DO BaixaEstoque ---

		// 7. Commit Final
		mock.ExpectCommit()

		// Execução
		err := repo.AddTrocaEPI(ctx, devolucao)

		// Validação
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteDevolucao(t *testing.T) {

	db, mock, ctx, err := mock(t)
	if err != nil {

		t.Fatal(err)
	}

	defer db.Close()

	repo := NewDevolucaoRepository(db)

	t.Run("sucesso ao fazer soft delete", func(t *testing.T) {

		id := 1
		mock.ExpectExec(regexp.QuoteMeta("update")).WithArgs(sql.Named("id", id)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteDevolucao(ctx, id)

		require.NoError(t, err)
	})

	t.Run("id ja deletado (não encontrado)", func(t *testing.T) {

		id := 44
		mock.ExpectExec(regexp.QuoteMeta("update")).WithArgs(sql.Named("id", id)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteDevolucao(ctx, id)

		require.Error(t, err)
	})

	t.Run("problema no banco de dados", func(t *testing.T) {

				id := 44
		mock.ExpectExec(regexp.QuoteMeta("update")).WithArgs(sql.Named("id", id)).
			WillReturnError(sql.ErrConnDone)

		err := repo.DeleteDevolucao(ctx, id)

		require.Error(t, err)

	})
}
