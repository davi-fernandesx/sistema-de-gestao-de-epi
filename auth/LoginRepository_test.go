package auth

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

)


func TestUsuarioCriacao(t *testing.T){

	//crinado o mocke a conexão falsa do banco de dados
	db, mock, err:=  sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))// mudando o comparador de query (desligar o regex)
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)

	}

	defer db.Close() //fechando a conexão falsa

	repo:= NewSqlLogin(db) //injetando a conexãop falsa com o repository

	loginCriacao:= &Login{Nome: "rada", Senha: "1234"} //criando os dados 

	IdEsperado:= 1 

	query:= "INSERT INTO login (usuario, senha) OUTPUT INSERTED.id values (@p1, @p2);"
	rows:= sqlmock.NewRows([]string{"id"}).AddRow(IdEsperado) //linha falsa retonando o id esperado

	mock.ExpectQuery(query). // query esperada
					WithArgs(loginCriacao.Nome, loginCriacao.Senha).//argumentas da query
					WillReturnRows(rows) // se tudo certo, retorna a linha

	ctx, cancelar:= context.WithTimeout(context.Background(), 5*time.Second) //ctx de 5 segundos

	defer cancelar() 

	login,err:= repo.AddLogin(ctx, loginCriacao)
	assert.NoError(t, err)
	assert.NotNil(t, login)
	assert.Equal(t, IdEsperado, login.ID)
	assert.NoError(t, mock.ExpectationsWereMet())


}

