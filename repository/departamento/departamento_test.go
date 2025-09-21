package departamento

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/require"
)



func Test_Departamento_add(t *testing.T){


	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewDepartamentoRepository(db)

	departamento:= model.Departamento{
			ID: 1,
			Departamento: "adm",
	}

	query:=regexp.QuoteMeta("insert into departamento (departamento) values (@departamento)")
	

	t.Run("sucesso ao add departamento", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(departamento.Departamento).WillReturnResult(sqlmock.NewResult(0,1))
		
		errDepartamento:= repo.AddDepartamento(ctx, &departamento)
		
		require.NoError(t,errDepartamento)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro de departamento ja existente", func(t *testing.T) {
		
		
		mock.ExpectExec(query).WithArgs(departamento.Departamento).WillReturnError(&mssql.Error{Number: 2627})
		errDepartamento:= repo.AddDepartamento(ctx, &departamento)
		require.Error(t,errDepartamento)
		require.Equal(t, repository.ErrDepartamentoJaExistente, errDepartamento)
		require.NoError(t, mock.ExpectationsWereMet())
		

	})

	t.Run("erros genericos", func(t *testing.T) {

		ErroGenericoDb:= errors.New("erro ao se conectar com o banco")
		mock.ExpectExec(query).WithArgs(departamento.Departamento).WillReturnError(ErroGenericoDb)

		err:= repo.AddDepartamento(ctx, &departamento)
		require.Error(t, err)
		require.Equal(t, ErroGenericoDb, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}


func  Test_DepartamentoRepository_delete(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewDepartamentoRepository(db)

	departamento:= model.Departamento{
			ID: 1,
			Departamento: "adm",
	}

	query:= regexp.QuoteMeta("delete from departamento where id = @id")

	t.Run("sucesso ao deletar um departamento", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(departamento.ID).WillReturnResult(sqlmock.NewResult(0,1))

		err := repo.DeletarDepartamento(ctx, departamento.ID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao deletar um departamento", func(t *testing.T) {
		
		mock.ExpectExec(query).WithArgs(departamento.ID). WillReturnResult(sqlmock.NewResult(0,0))
		err:= repo.DeletarDepartamento(ctx, departamento.ID)
		require.Error(t, err)
		require.Equal(t, repository.ErrDepartamentoNaoEncontrado, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

		t.Run("erro ao buscar linhas afetadas", func(t *testing.T) {

		driveErro:= errors.New("driver: RowsAffected not supported")

		mock.ExpectExec(query).WithArgs(departamento.ID).WillReturnResult(sqlmock.NewErrorResult(driveErro))

		err:= repo.DeletarDepartamento(ctx, departamento.ID)
		require.Error(t, err)
		require.Equal(t, repository.ErrLinhasAfetadas,err )
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_DepartamentoRepository_Buscar(t *testing.T) {



	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewDepartamentoRepository(db)

	departamento:= model.Departamento{
			ID: 1,
			Departamento: "adm",
	}
	
	query:= regexp.QuoteMeta("select departamento from departamento where id = @id")

	t.Run("testando o sucesso de de buscar um departamento", func(t *testing.T) {

		linhas:= sqlmock.NewRows([]string{"id", "departamento"}).AddRow(departamento.ID, departamento.Departamento)
		mock.ExpectQuery(query).WithArgs(departamento.ID).WillReturnRows(linhas)

		departamentodb, err:= repo.BuscarDepartamento(ctx, departamento.ID)
		require.NoError(t, err)
		require.NotNil(t, departamentodb)
		require.Equal(t, departamento.ID, departamentodb.ID)
		require.Equal(t, departamento.Departamento, departamentodb.Departamento)

		require.NoError(t, mock.ExpectationsWereMet())


	})

	t.Run("erro de departamento inexistente", func(t *testing.T) {

		departamentoInexistenteId:= 2

		mock.ExpectQuery(query).WithArgs(departamentoInexistenteId).WillReturnError(sql.ErrNoRows)

		Departamentopdb, err:= repo.BuscarDepartamento(ctx, departamentoInexistenteId)
		require.Error(t, err)
		require.Equal(t, repository.ErrDepartamentoNaoEncontrado, err)
		require.Nil(t, Departamentopdb)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de scan", func(t *testing.T) {

		linhas:= sqlmock.NewRows([]string{"id", "departamento"})
		linhas.AddRow(departamento.ID, departamento.Departamento).RowError(0, repository.ErrFalhaAoEscanearDados)

		mock.ExpectQuery(query).WithArgs(departamento.ID).WillReturnRows(linhas)

		departamentodb, err:= repo.BuscarDepartamento(ctx, departamento.ID)
		require.Error(t,err)
		require.Equal(t, repository.ErrFalhaAoEscanearDados, err)
		require.Nil(t, departamentodb)

		require.NoError(t, mock.ExpectationsWereMet())
		 
	})

}

func Test_DepartamentoRepository_buscarTodos(t *testing.T){


	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewDepartamentoRepository(db)

	departamento:= model.Departamento{
			ID: 1,
			Departamento: "adm",
	}

	query:= regexp.QuoteMeta("select id, departamento from departamento")

	t.Run(" sucesso ao buscar muitos departamentos", func(t *testing.T) {


		linhas:= sqlmock.NewRows([]string{"id", "departamento"}).
			AddRow(departamento.ID, departamento.Departamento).
			AddRow(2, "rh").
			AddRow(3, "ti")

		mock.ExpectQuery(query).WillReturnRows(linhas)

		departamentodb, err:= repo.BuscarTodosDepartamentos(ctx)
		require.NoError(t,err)
		require.NotNil(t, departamentodb)
		require.Len(t, *departamentodb, 3)
		require.Equal(t, departamento.Departamento, (*departamentodb)[0].Departamento)
		require.Equal(t, "rh", (*departamentodb)[1].Departamento)
		require.Equal(t, "ti", (*departamentodb)[2].Departamento)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("erro na consulta do banco de dados", func(t *testing.T) {

		mock.ExpectQuery(query).WillReturnError(repository.ErrConexaoDb)

		departamentodb, err:= repo.BuscarTodosDepartamentos(ctx)
		require.Error(t,err)
		require.Nil(t, departamentodb)
		require.Equal(t, repository.ErrBuscarTodosDepartamentos, err)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("erro ao iterar sobre os dados", func(t *testing.T) {

		linhas:= sqlmock.NewRows([]string{"id", "departamento"}).
		AddRow(departamento.ID, departamento.Departamento).
		AddRow(2, " ti").
		CloseError(repository.ErrDadoIncompativel)

		mock.ExpectQuery(query).WillReturnRows(linhas)

		departamentodb, err:= repo.BuscarTodosDepartamentos(ctx)
		require.Error(t, err)
		require.Nil(t, departamentodb)
		require.Equal(t, repository.ErrIterarSobreDepartamentos, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})


}