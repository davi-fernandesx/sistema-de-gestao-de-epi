package funcao

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/require"
)


func Test_Funcao_add(t *testing.T){

	ctx := context.Background()
	db, mock, err:= sqlmock.New()
	require.NoError(t, err)

	repo:= NewfuncaoRepository(db)

	funcao:= model.Funcao{
		ID: 1,
		Funcao: "faqueiro",
	}

	query:= regexp.QuoteMeta("insert into funcao (funcao) values (@funcao)")

	t.Run("sucesso ao add  funcao", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(funcao.Funcao).WillReturnResult(sqlmock.NewResult(0,1))
		errFuncao:= repo.AddFuncao(ctx, &funcao)

		require.NoError(t,errFuncao)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}



