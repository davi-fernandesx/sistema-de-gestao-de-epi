package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestUsuarioCriacao(t *testing.T){

	//crinado o mocke a conexão falsa do banco de dados
	db, mock, err:=  sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))// mudando o comparador de query (desligar o regex)
	require.NoError(t, err)

	defer db.Close() //fechando a conexão falsa

	repo:= NewSqlLogin(db) //injetando a conexãop falsa com o repository

	loginCriacao:= &model.Login{Nome: "rada", Senha: "1234"} //criando os dados 

	IdEsperado:= 1 

	query:= "INSERT INTO login (usuario, senha) OUTPUT INSERTED.id values (@p1, @p2);"
	rows:= sqlmock.NewRows([]string{"id"}).AddRow(IdEsperado) //linha falsa retonando o id esperado

	mock.ExpectQuery(query). // query esperada
					WithArgs(loginCriacao.Nome, loginCriacao.Senha).//argumentas da query
					WillReturnRows(rows) // se tudo certo, retorna a linha

	ctx, cancelar:= context.WithTimeout(context.Background(), 5*time.Second) //ctx de 5 segundos

	defer cancelar() 

	_,err = repo.AddLogin(ctx, loginCriacao)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())


}

func TestDeletarLogin(t *testing.T) {
	// --- Setup Inicial Comum a todos os sub-testes ---
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Usamos o nome da sua struct concreta aqui
	repo := &SqlServerLogin{db: db}

	query := `delete from login where id = @id`

	// --- Início dos Sub-Testes ---

	t.Run("Sucesso - Deleta o login e retorna nil", func(t *testing.T) {
		// Arrange (Preparação)
		idParaDeletar := 42

		// Esperamos que o comando `Exec` seja chamado com a nossa query.
		// Para um DELETE, usamos ExpectExec.
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", idParaDeletar)). // Verifica se o argumento correto foi passado
			WillReturnResult(sqlmock.NewResult(0, 1)) // Simula que 1 linha foi afetada

		// Act (Execução)
		err := repo.DeletarLogin(context.Background(), idParaDeletar)

		// Assert (Verificação)
		require.NoError(t, err) // Esperamos que nenhum erro seja retornado
		require.NoError(t, mock.ExpectationsWereMet()) // Garante que a expectativa foi cumprida
	})

	t.Run("Falha - ID não encontrado", func(t *testing.T) {
		// Arrange
		idInexistente := 99

		// A query executa com sucesso, mas o banco informa que 0 linhas foram afetadas.
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", idInexistente)).
			WillReturnResult(sqlmock.NewResult(0, 0)) // <-- AQUI ESTÁ A DIFERENÇA

		// Act
		err := repo.DeletarLogin(context.Background(), idInexistente)

		// Assert
		require.Error(t, err) // Esperamos um erro
		// Verificamos se o erro retornado é exatamente o que a função deve criar
		expectedErrorMsg := fmt.Sprintf("nenhuma categoria encontrada com o id: %d", idInexistente)
		assert.EqualError(t, err, expectedErrorMsg)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Falha - Erro na execução do delete", func(t *testing.T) {
		// Arrange
		idParaDeletar := 10

		// Simulamos um erro genérico do banco de dados durante a execução.
		mock.ExpectExec(query).
			WithArgs(sql.Named("id", idParaDeletar)).
			WillReturnError(errors.New("erro de conexão")) // Simula uma falha

		// Act
		err := repo.DeletarLogin(context.Background(), idParaDeletar)

		// Assert
		require.Error(t, err) // Esperamos um erro
		// Verificamos se o erro retornado é o nosso erro customizado
		assert.Equal(t, erroAoApagarUmLogin, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLogin(t *testing.T) {
	// --- Setup Inicial ---
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Assumindo que sua struct se chama SqlServerLogin e o construtor é New...
	repo := &SqlServerLogin{db: db}

	// Query correta que o teste vai esperar
	query := `SELECT usuario, senha FROM login WHERE usuario = @nome`

	// --- Sub-Testes ---

	t.Run("Sucesso - Encontra o login e retorna os dados", func(t *testing.T) {
		// Arrange (Preparação)
		nomeParaBuscar := "davi"
		loginEsperado := &model.Login{
			Nome:  "davi",
			Senha: "hash_da_senha_123", // O banco sempre deve retornar o hash
		}

		// Criamos as linhas que o mock deve retornar
		rows := sqlmock.NewRows([]string{"usuario", "senha"}).
			AddRow(loginEsperado.Nome, loginEsperado.Senha)

		// Esperamos que uma query (SELECT) seja executada
		mock.ExpectQuery(query).
			WithArgs(sql.Named("nome", nomeParaBuscar)).
			WillReturnRows(rows) // E dizemos ao mock para retornar as linhas que criamos

		// Act (Execução)
		loginRetornado, err := repo.Login(context.Background(), nomeParaBuscar)

		// Assert (Verificação)
		require.NoError(t, err)
		require.NotNil(t, loginRetornado)
		assert.Equal(t, loginEsperado.Nome, loginRetornado.Nome)
		assert.Equal(t, loginEsperado.Senha, loginRetornado.Senha)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Falha - Usuário não encontrado", func(t *testing.T) {
		// Arrange
		nomeInexistente := "usuario_fantasma"

		// Esperamos que a query seja executada, mas desta vez, simulamos o erro
		// que o banco retorna quando não encontra nenhum registro: sql.ErrNoRows.
		mock.ExpectQuery(query).
			WithArgs(sql.Named("nome", nomeInexistente)).
			WillReturnError(sql.ErrNoRows)

		// Act
		loginRetornado, err := repo.Login(context.Background(), nomeInexistente)

		// Assert
		require.Error(t, err) // Esperamos um erro
		assert.Nil(t, loginRetornado) // O login retornado deve ser nulo
		
		// Verificamos se a mensagem de erro é a que nossa função customizou
		expectedErrorMsg := fmt.Sprintf("erro ao encontrar login: %v", sql.ErrNoRows)
		assert.EqualError(t, err, expectedErrorMsg)
		
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
