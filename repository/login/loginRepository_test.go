package login

import (
	"context"
	"database/sql"
	"errors"
	"regexp"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func Test_LoginRepository_Add(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewLogin(db)

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
		assert.True(t, errors.Is(err,  Errors.ErrSalvar), "erro deveria ser erro ao salvar")
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
		assert.True(t, errors.Is(err, Errors.ErrInternal), "erro tem que ser do tipo interno")

		require.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_LoginRepositoryDelete(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewLogin(db)

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

		mock.ExpectExec(query).WithArgs(login.ID).WillReturnResult(sqlmock.NewResult(0,0))

		err:= repo.DeletarLogin(ctx, 1)
		assert.True(t, errors.Is(err, Errors.ErrNaoEncontrado), "erro tem que ser do tipo nao encontrado")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("erro ao buscar linhas afetadas", func(t *testing.T) {


		query:= regexp.QuoteMeta("delete from login where id = @id")

		

		mock.ExpectExec(query).WithArgs(login.ID).WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))

		err:= repo.DeletarLogin(ctx, login.ID)
		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrLinhasAfetadas), "erro tem que ser do tipo linhas afetadas")
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_loginRepositoryBuscarUsuario(t *testing.T){

	ctx:= context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo:= NewLogin(db)

	login:= model.Login{
		ID: 1,
		Nome: "rada",
		Senha: "1234",

	}

	t.Run("testando o sucesso da busca por nome do usuario", func(t *testing.T) {


		query:= `select usuario, senha from login
		where usuario = @usuario ;`

		linhas:= sqlmock.NewRows([]string{"usuario", "senha"}).AddRow(login.Nome, login.Senha)
		mock.ExpectQuery(query).WithArgs(login.Nome).WillReturnRows(linhas)

		usuario,err:=  repo.BuscaPorNome(ctx, login.Nome)
		require.NoError(t,err)
		require.NotNil(t, usuario)
		require.Equal(t, login.Nome, usuario.Nome)
		require.Equal(t, login.Senha,usuario.Senha)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de usuario inexistente", func(t *testing.T) {

		query:= regexp.QuoteMeta("select usuario, senha from login where usuario = @usuario ;")

		usuarioInexistente:= "teste"

		mock.ExpectQuery(query).WithArgs(usuarioInexistente).WillReturnError(sql.ErrNoRows)

		usuario, err:= repo.BuscaPorNome(ctx, usuarioInexistente )
		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrNaoEncontrado), "erro deveria ser do tipo nao encotrado")
		require.Nil(t, usuario)

		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de scan", func(t *testing.T) {

		query:= regexp.QuoteMeta("select usuario, senha from login where usuario = @usuario ;")

		linhas:= sqlmock.NewRows([]string{"usuario", "senha"})
		
		linhas.AddRow(login.Nome, login.Senha).RowError(0, Errors.ErrFalhaAoEscanearDados)

		mock.ExpectQuery(query).WithArgs(login.Nome).WillReturnRows(linhas)

		usuario, err:= repo.BuscaPorNome(ctx, login.Nome)

		require.Error(t, err)
		assert.True(t, errors.Is(err, Errors.ErrFalhaAoEscanearDados), "erro deveria ser do  tipo escanear dados")
		require.Nil(t, usuario)

		require.NoError(t, mock.ExpectationsWereMet())



	})


}