package epi

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

func Test_Epi_add(t *testing.T) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo := NewEpiRepository(db)

	EpiInserir := model.EpiInserir{
		ID:             1,
		Nome:           "botas isolante antiderapante",
		Fabricante:     "maicol",
		CA:             "36025",
		Descricao:      "Bota de Pvc Cano Médio 28Cm Branca",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		AlertaMinimo:   10,
		IDprotecao:     1,
	}

	query := regexp.QuoteMeta(`insert into epi (nome, fabricante, CA, descricao,
				data_fabricacao, data_validade, 
				validade_CA, id_tipo_protecao, alerta_minimo) values (
				@nome, @fabricante, @CA, @descricao,@data_fabricacao, @data_validade,
				@validade_CA, @id_tipo_protecao, @alerta_minimo )`)

	t.Run("testando o sucesso ao adicionar um  epi no banco de dados", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(EpiInserir.Nome, EpiInserir.Fabricante, EpiInserir.CA, EpiInserir.Descricao,
			EpiInserir.DataFabricacao,
			EpiInserir.DataValidade, EpiInserir.DataValidadeCa, EpiInserir.IDprotecao, EpiInserir.AlertaMinimo).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ErrEpi := repo.AddEpi(ctx, &EpiInserir)
		require.NoError(t, ErrEpi)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro ao adicionar um epi no banco de dados", func(t *testing.T) {

		errGen :=Errors.ErrSalvar
		mock.ExpectExec(query).WithArgs(EpiInserir.Nome, EpiInserir.Fabricante, EpiInserir.CA, EpiInserir.Descricao,
			EpiInserir.DataFabricacao,
			EpiInserir.DataValidade, EpiInserir.DataValidadeCa, EpiInserir.IDprotecao, EpiInserir.AlertaMinimo).
			WillReturnError(errGen)

		errEpi := repo.AddEpi(ctx, &EpiInserir)
		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrInternal), "erro tem que ser do tipo internal")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_buscarEpi(t *testing.T) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo := NewEpiRepository(db)

	Epi := model.Epi{
		ID:             1,
		Nome:           "botas isolante antiderapante",
		Fabricante:     "maicol",
		CA:             "36025",
		Descricao:      "Bota de Pvc Cano Médio 28Cm Branca",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		AlertaMinimo:   10,
		IDprotecao:     1,
		NomeProtecao:   "protecao_para_cabeça",
	}

	query := regexp.QuoteMeta(`
				select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id		
			where
				e.id = @id`)

	t.Run("testando o sucesso ao buscar um epi", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"fabricante",
			"CA",
			"descricao",
			"data_fabricacao",
			"data_validade",
			"validade_CA",
			"alerta_minimo",
			"id_tipo_protecao",
			"nome_protecao",
		}).AddRow(
			Epi.ID,
			Epi.Nome,
			Epi.Fabricante,
			Epi.CA,
			Epi.Descricao,
			Epi.DataFabricacao,
			Epi.DataValidade,
			Epi.DataValidadeCa,
			Epi.AlertaMinimo,
			Epi.IDprotecao,
			Epi.NomeProtecao,
		)

		mock.ExpectQuery(query).WithArgs(Epi.ID).WillReturnRows(linhas)

		epidb, err := repo.BuscarEpi(ctx, Epi.ID)
		require.NoError(t, err)
		require.Equal(t, Epi.ID, epidb.ID)
		require.Equal(t, Epi.Nome, epidb.Nome)
		require.Equal(t, Epi.Fabricante, epidb.Fabricante)
		require.Equal(t, Epi.CA, epidb.CA)
		require.Equal(t, Epi.Descricao, epidb.Descricao)
		require.Equal(t, Epi.DataFabricacao, epidb.DataFabricacao)
		require.Equal(t, Epi.DataValidade, epidb.DataValidade)
		require.Equal(t, Epi.DataValidadeCa, epidb.DataValidadeCa)
		require.Equal(t, Epi.AlertaMinimo, epidb.AlertaMinimo)
		require.Equal(t, Epi.IDprotecao, epidb.IDprotecao)
		require.Equal(t, Epi.NomeProtecao, epidb.NomeProtecao)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro ao buscar um epi", func(t *testing.T) {

		idEpiNaoEXistente := 2

		mock.ExpectQuery(query).WithArgs(idEpiNaoEXistente).WillReturnError(sql.ErrNoRows)

		epiDB, err := repo.BuscarEpi(ctx, idEpiNaoEXistente)
		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrNaoEncontrado), "erro tem que ser do tipo nao encontrado")
		require.Nil(t, epiDB)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando erro  ao escanear os epi do banco de dados", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{

			"nome",
			"fabricante",
			"CA",
			"descricao",
			"data_fabricacao",
			"data_validade",
			"validade_CA",
			"id_tipo_protecao",
			"alerta_minimo",
		}).AddRow(
			//Epi.ID,  nao passando o id, para gerar o erro
			Epi.Nome,
			Epi.Fabricante,
			Epi.CA,
			Epi.Descricao,
			Epi.DataFabricacao,
			Epi.DataValidade,
			Epi.DataValidadeCa,
			Epi.IDprotecao,
			Epi.AlertaMinimo,
		)

		mock.ExpectQuery(query).WithArgs(Epi.ID).WillReturnRows(linhas)

		epidb, err := repo.BuscarEpi(ctx, Epi.ID)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrFalhaAoEscanearDados), "erro tem que ser do tipo escanaear")
		require.Nil(t, epidb)

		require.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_buscarTodosEpis(t *testing.T) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo := NewEpiRepository(db)

	Epi1 := model.Epi{
		ID:             1,
		Nome:           "luvas isolante termicas",
		Fabricante:     "maicol",
		CA:             "36022",
		Descricao:      "luvas de borracha",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		AlertaMinimo:   20,
		IDprotecao:     3,
		NomeProtecao:   "protecao_para_os_pes",
	}

	Epi2 := model.Epi{
		ID:             2,
		Nome:           "botas isolante antiderapante",
		Fabricante:     "maicol",
		CA:             "36025",
		Descricao:      "Bota de Pvc Cano Médio 28Cm Branca",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		AlertaMinimo:   10,
		IDprotecao:     1,
		NomeProtecao:   "protecao_para_os_pes",
	}

	query := regexp.QuoteMeta(`
					select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id`)

	t.Run("testando o sucesso ao buscar todos os epis", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"fabricante",
			"CA",
			"descricao",
			"data_fabricacao",
			"data_validade",
			"validade_CA",
			"alerta_minimo",
			"id_tipo_protecao",
			"nomeProtecao",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			Epi1.CA,
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.AlertaMinimo,
			Epi1.IDprotecao,
			Epi1.NomeProtecao,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			Epi2.CA,
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.AlertaMinimo,
			Epi2.IDprotecao,
			Epi2.NomeProtecao,
		)

		EpisESperados := []model.Epi{Epi1, Epi2}

		mock.ExpectQuery(query).WillReturnRows(linhas)

		epis, err := repo.BuscarTodosEpi(ctx)
		require.NoError(t, err)
		require.NotNil(t, epis)
		require.Len(t, epis, 2)
		require.Equal(t, EpisESperados, epis)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de driver do sql", func(t *testing.T) {

		mock.ExpectQuery(query).WillReturnError(Errors.ErrBuscarTodos)

		epis, err := repo.BuscarTodosEpi(ctx)
		require.Error(t, err)
		require.Empty(t, epis)
		assert.True(t, errors.Is(err, Errors.ErrBuscarTodos))

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro ao buscar todos os epis", func(t *testing.T) {
		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"fabricante",
			//"CA",
			"descricao",
			"data_fabricacao",
			"data_validade",
			"validade_CA",
			"alerta_minimo",
			"id_tipo_protecao",
			"nomeProtecao",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			//Epi1.CA,
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.AlertaMinimo,
			Epi1.IDprotecao,
			Epi1.NomeProtecao,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			//Epi2.CA,
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.AlertaMinimo,
			Epi2.IDprotecao,
			Epi2.NomeProtecao,
		)

		mock.ExpectQuery(query).WillReturnRows(linhas)

		epiDb, err := repo.BuscarTodosEpi(ctx)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrFalhaAoEscanearDados), "erro tem que ser do tipo escanaear")
		require.Nil(t, epiDb)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando erro ao iterar sobre os epis", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"fabricante",
			"CA",
			"descricao",
			"data_fabricacao",
			"data_validade",
			"validade_CA",
			"alerta_minimo",
			"id_tipo_protecao",
			"nomeProtecao",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			Epi1.CA,
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.AlertaMinimo,
			Epi1.IDprotecao,
			Epi1.NomeProtecao,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			Epi2.CA,
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.AlertaMinimo,
			Epi2.IDprotecao,
			Epi2.NomeProtecao,
		).CloseError(Errors.ErrDadoIncompativel)

		mock.ExpectQuery(query).WillReturnRows(linhas)

		epis, err := repo.BuscarTodosEpi(ctx)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrAoIterar), "erro tem que ser do tipo iterar")
		require.Nil(t, epis)

		require.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_deletarEpi(t *testing.T) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo := NewEpiRepository(db)

	Epi1 := model.Epi{
		ID:             1,
		Nome:           "luvas isolante termicas",
		Fabricante:     "maicol",
		CA:             "36022",
		Descricao:      "luvas de borracha",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		AlertaMinimo:   20,
		IDprotecao:     3,
		NomeProtecao: "protecao_para_as_maos",
	}

	query := regexp.QuoteMeta(`delete from epi where id = @id`)

	t.Run("testando o sucesso ao deletar um epi da base de dados", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewResult(0, 1))

		errEpi := repo.DeletarEpi(ctx, Epi1.ID)
		require.NoError(t, errEpi)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro na execução da query no repository", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnError(Errors.ErrConexaoDb)

		errEpi := repo.DeletarEpi(ctx, Epi1.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrConexaoDb))
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de linhas afetadas", func(t *testing.T) {

		

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		errEpi := repo.DeletarEpi(ctx, Epi1.ID)

		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrLinhasAfetadas))

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("epi nao encontrado", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewResult(0, 0))

		errEpi := repo.DeletarEpi(ctx, Epi1.ID)
		require.Error(t, errEpi)
		assert.True(t, errors.Is(errEpi, Errors.ErrNaoEncontrado))

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
