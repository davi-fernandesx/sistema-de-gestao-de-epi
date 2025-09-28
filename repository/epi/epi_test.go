package epi

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	"github.com/stretchr/testify/require"
)

func Test_Epi_add(t *testing.T) {

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
		IDprotecao:     1,
		AlertaMinimo:   10,
	}

	query := regexp.QuoteMeta(`insert into epi (nome, fabricante, CA, descricao,
				data_fabricacao, data_validade, 
				validade_CA, id_tipo_protecao, alerta_minimo) values (
				@nome, @fabricante, @CA, @descricao,@data_fabricacao, @data_validade,
				@validade_CA, @id_tipo_protecao, @alerta_minimo )`)

	t.Run("testando o sucesso ao adicionar um  epi no banco de dados", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi.Nome, Epi.Fabricante, Epi.CA, Epi.Descricao,
			Epi.DataFabricacao,
			Epi.DataValidade, Epi.DataValidadeCa, Epi.IDprotecao, Epi.AlertaMinimo).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ErrEpi := repo.AddEpi(ctx, &Epi)
		require.NoError(t, ErrEpi)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro ao adicionar um epi no banco de dados", func(t *testing.T) {

		errGen := repository.ErrEpiAoAdicionarEpi
		mock.ExpectExec(query).WithArgs(Epi.Nome, Epi.Fabricante, Epi.CA, Epi.Descricao,
			Epi.DataFabricacao,
			Epi.DataValidade, Epi.DataValidadeCa, Epi.IDprotecao, Epi.AlertaMinimo).
			WillReturnError(errGen)

		errEpi := repo.AddEpi(ctx, &Epi)
		require.Error(t, errEpi)
		require.Equal(t, errGen, errEpi)
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
		IDprotecao:     1,
		AlertaMinimo:   10,
	}

	query := regexp.QuoteMeta(`select id, nome, fabricante, CA, descricao, data_fabricacao, data_validade, 
							validade_CA, id_tipo_protecao, alerta_minimo
							from epi where id = @id`)

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
			"id_tipo_protecao",
			"alerta_minimo",
		}).AddRow(
			Epi.ID,
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
		require.NoError(t, err)
		require.Equal(t, Epi.ID, epidb.ID)
		require.Equal(t, Epi.Nome, epidb.Nome)
		require.Equal(t, Epi.Fabricante, epidb.Fabricante)
		require.Equal(t, Epi.CA, epidb.CA)
		require.Equal(t, Epi.Descricao, epidb.Descricao)
		require.Equal(t, Epi.DataFabricacao, epidb.DataFabricacao)
		require.Equal(t, Epi.DataValidade, epidb.DataValidade)
		require.Equal(t, Epi.DataValidadeCa, epidb.DataValidadeCa)
		require.Equal(t, Epi.IDprotecao, epidb.IDprotecao)
		require.Equal(t, Epi.AlertaMinimo, epidb.AlertaMinimo)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro ao buscar um epi", func(t *testing.T) {

		idEpiNaoEXistente := 2

		mock.ExpectQuery(query).WithArgs(idEpiNaoEXistente).WillReturnError(sql.ErrNoRows)

		epiDB, err := repo.BuscarEpi(ctx, idEpiNaoEXistente)
		require.Error(t, err)
		require.Equal(t, repository.ErrAoProcurarEpi, err)
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
		require.Equal(t, repository.ErrFalhaAoEscanearDados, err)
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
		IDprotecao:     3,
		AlertaMinimo:   20,
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
		IDprotecao:     1,
		AlertaMinimo:   10,
	}

	query := regexp.QuoteMeta(`select id, nome, fabricante, CA, descricao,
	 		data_fabricacao, data_validade, validade_CA, 
	 		id_tipo_protecao, alerta_minimo
			from epi`)

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
			"id_tipo_protecao",
			"alerta_minimo",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			Epi1.CA,
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.IDprotecao,
			Epi1.AlertaMinimo,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			Epi2.CA,
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.IDprotecao,
			Epi2.AlertaMinimo,
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

		mock.ExpectQuery(query).WillReturnError(repository.ErrAoBuscarTodosOsEpis)

		epis, err:= repo.BuscarTodosEpi(ctx)
		require.Error(t,err)
		require.Empty(t,epis)
		require.Equal(t, repository.ErrAoBuscarTodosOsEpis, err)

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
			"id_tipo_protecao",
			"alerta_minimo",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			//Epi1.CA, gerando o erro
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.IDprotecao,
			Epi1.AlertaMinimo,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			//Epi2.CA, gerando o erro
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.IDprotecao,
			Epi2.AlertaMinimo,
		)
			

		mock.ExpectQuery(query).WillReturnRows(linhas)

		epiDb, err:= repo.BuscarTodosEpi(ctx)

		require.Error(t, err)
		require.Equal(t, repository.ErrFalhaAoEscanearDados, err)
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
			"id_tipo_protecao",
			"alerta_minimo",
		}).AddRow(
			Epi1.ID,
			Epi1.Nome,
			Epi1.Fabricante,
			Epi1.CA,
			Epi1.Descricao,
			Epi1.DataFabricacao,
			Epi1.DataValidade,
			Epi1.DataValidadeCa,
			Epi1.IDprotecao,
			Epi1.AlertaMinimo,
		).AddRow(
			Epi2.ID,
			Epi2.Nome,
			Epi2.Fabricante,
			Epi2.CA,
			Epi2.Descricao,
			Epi2.DataFabricacao,
			Epi2.DataValidade,
			Epi2.DataValidadeCa,
			Epi2.IDprotecao,
			Epi2.AlertaMinimo,
		).CloseError(repository.ErrDadoIncompativel)


		mock.ExpectQuery(query).WillReturnRows(linhas)

		epis, err:= repo.BuscarTodosEpi(ctx)

		require.Error(t, err)
		require.Equal(t, repository.ErrAoInterarSobreEpis, err)
		require.Nil(t, epis)

		require.NoError(t, mock.ExpectationsWereMet())

	})

}


func Test_deletarEpi(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo:= NewEpiRepository(db)

	Epi1 := model.Epi{
		ID:             1,
		Nome:           "luvas isolante termicas",
		Fabricante:     "maicol",
		CA:             "36022",
		Descricao:      "luvas de borracha",
		DataFabricacao: time.Now(),
		DataValidade:   time.Now().AddDate(3, 12, 0),
		DataValidadeCa: time.Now().AddDate(2, 0, 0),
		IDprotecao:     3,
		AlertaMinimo:   20,
	}

	query:= regexp.QuoteMeta(`delete from epi where id = @id`)

	t.Run("testando o sucesso ao deletar um epi da base de dados", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewResult(0,1))

		errEpi:= repo.DeletarEpi(ctx, Epi1.ID)
		require.NoError(t,errEpi)
		
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro na execução da query no repository", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnError(repository.ErrConexaoDb)

		errEpi:= repo.DeletarEpi(ctx, Epi1.ID)

		require.Error(t, errEpi)
		require.Equal(t, repository.ErrConexaoDb, errEpi)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de linhas afetadas", func(t *testing.T) {

		driveErro := errors.New("driver: RowsAffected not supported")

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewErrorResult(driveErro))

		errEpi:= repo.DeletarEpi(ctx, Epi1.ID)

		require.Error(t, errEpi)
		require.Equal(t, repository.ErrLinhasAfetadas, errEpi)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("epi nao encontrado", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(Epi1.ID).WillReturnResult(sqlmock.NewResult(0,0))

		errEpi:= repo.DeletarEpi(ctx, Epi1.ID)
		require.Error(t,errEpi)
		require.Equal(t, repository.ErrEpiNaoEncontrado, errEpi)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}