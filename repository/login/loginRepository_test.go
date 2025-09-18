package login

import (
	"context"
	"errors"
	"regexp"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func Test_LoginRepository_Add(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewSqlLogin(db)

	t.Run("sucesso ao add o login ao banco de dados", func(t *testing.T) {

	 //injetando o mock no repository

		query:= regexp.QuoteMeta("INSERT INTO login (usuario, senha) values (@p1, @p2)")
		login:= model.Login{
			ID: 1,
			Nome: "rada",
			Senha: "1234",

		}

		mock.ExpectExec(query).WithArgs(login.Nome, login.Senha).WillReturnResult(sqlmock.NewResult(0,1))

		err = repo.AddLogin(ctx, &login)
		assert.NoError(t,err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("erro de usuario ja existente", func(t *testing.T) {

		query:= regexp.QuoteMeta("INSERT INTO login (usuario, senha) values (@p1, @p2)")
		login:= model.Login{
			ID: 1,
			Nome: "rada",
			Senha: "1234",

		}

		sqlErr:= mssql.Error{
			Number: 2627,
			Message: "Violation of UNIQUE KEY constraint.",
		}

		mock.ExpectExec(query).WithArgs(login.Nome, login.Senha).WillReturnError(&sqlErr)

		err = repo.AddLogin(ctx, &login)
		require.Error(t,err)
		require.Equal(t, repository.ErrusuarioJaExistente, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erros genericos", func(t *testing.T) {
		query:= regexp.QuoteMeta("INSERT INTO login (usuario, senha) values (@p1, @p2);")

		login:= model.Login{
			ID: 1,
			Nome: "rada",
			Senha: "123445",
		}

		ErroGenericoDb:= errors.New("erro ao se conectar com o banco")

		mock.ExpectExec(query).WithArgs(login.Nome, login.Senha).WillReturnError(ErroGenericoDb)

		err:=repo.AddLogin(ctx, &login)
		require.Error(t, err)
		require.Equal(t, ErroGenericoDb, err)

		require.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_LoginRepositoryDelete(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewSqlLogin(db)

	login:= model.Login{
		ID: 1,
		Nome: "rada",
		Senha: "1234",

	}

	t.Run("sucesso ao deletar um usuario", func(t *testing.T) {

		query:= regexp.QuoteMeta("delete from login where id = @id")

		mock.ExpectExec(query).WithArgs(login.ID).WillReturnResult(sqlmock.NewResult(0,1))

		err:= repo.DeletarLogin(ctx, 1)
		require.NoError(t,err)
		require.NoError(t, mock.ExpectationsWereMet())

		
	})

	t.Run("Erro ao deletar um usuario", func(t *testing.T) {

		query:= regexp.QuoteMeta("delete from login where id = @id")

		mock.ExpectExec(query).WithArgs(login.ID).WillReturnError(repository.ErrLinhasAfetadas)

		err:= repo.DeletarLogin(ctx, 1)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Usuario não encontrado", func(t *testing.T) {

		query:= regexp.QuoteMeta("delete from login where id = @id")

		mock.ExpectExec(query).WithArgs(login.ID).WillReturnError(repository.ErrUsuarioNaoEncontrado)

		err:= repo.DeletarLogin(ctx, 1)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

		

	})

}
