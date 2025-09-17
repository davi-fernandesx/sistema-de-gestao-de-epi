package login

import (
	"context"
	"regexp"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func Test_LoginRepository(t *testing.T){

	ctx:= context.Background()

	t.Run("sucesso ao add o login ao banco de dados", func(t *testing.T) {

		db, mock, err:= sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo:= NewSqlLogin(db)

		query:= regexp.QuoteMeta("INSERT INTO login (usuario, senha) values (@p1, @p2)")
		login:= model.Login{
			ID: 1,
			Nome: "rada",
			Senha: "1234",

		}

		mock.ExpectQuery(query).WithArgs(login.Nome, login.Senha).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err = repo.AddLogin(ctx, &login)
		assert.NoError(t,err)

		assert.Equal(t, 1, login.ID)

		assert.NoError(t, mock.ExpectationsWereMet())



	})
}
