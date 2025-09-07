//go:build integration

package auth

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"os"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb" // Driver do SQL Server
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

var realDB *sql.DB // Variável global para ser usada pelos testes neste pacote

// TestMain é a função de setup e teardown para os testes de integração.
func TestMain(m *testing.M) {
	// 1. Obter a string de conexão de uma variável de ambiente.
	// Isso evita deixar senhas no código.
	errr := godotenv.Load("../../configs/.env")
	if errr != nil {
		log.Fatal("erro ao abrir arquivo .env")
	}
	db_server := os.Getenv("DB_SERVER")
	db_port := os.Getenv("DB_PORT")
	db_database := os.Getenv("DATABASE")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("SA_PASSWORD")

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", db_user, db_pass, db_server, db_port, db_database)
	if connString == "" {
		log.Println("DATABASE_URL_TEST não definida. Pulando testes de integração.")
		// Pula os testes se a variável não estiver configurada.
		os.Exit(0)
	}

	var err error
	// 2. Tentar conectar ao banco de dados em um loop (retry).
	// O container pode demorar alguns segundos para ficar pronto.
	for i := 0; i < 5; i++ {
		realDB, err = sql.Open("sqlserver", connString)
		if err == nil {
			err = realDB.Ping()
			if err == nil {
				log.Println("Conexão com o banco de dados de teste bem-sucedida!")
				break // Conexão bem-sucedida, sai do loop.
			}
		}
		log.Printf("Tentando conectar ao banco de dados... Tentativa %d/5. Erro: %v", i+1, err)
		time.Sleep(3 * time.Second) // Espera 3 segundos antes de tentar novamente.
	}

	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados após várias tentativas: %v", err)
	}
	defer realDB.Close()

	// 3. Limpar o ambiente antes de rodar os testes.
	// Garante que a tabela esteja vazia e pronta para os testes.
	_, err = realDB.ExecContext(context.Background(), "DELETE FROM login;")
	if err != nil {
		log.Fatalf("Falha ao limpar a tabela de login: %v", err)
	}

	// 4. Rodar os testes.
	log.Println("Iniciando a execução dos testes de integração...")
	code := m.Run()

	// 5. Sair com o código de status dos testes.
	os.Exit(code)
}

// (continuação do arquivo login_integration_test.go)

func TestIntegration_CriacaoUsuario_Successo(t *testing.T) {

	repo := NewSqlLogin(realDB) // usando a conexao do realDB!
	ctx := context.Background()
	loginCriacao := &model.Login{
		Nome:  "aaaa",
		Senha: "teste.ok",
	}

	// Chamando o método que queremos testar. Isso vai realmente fazer um INSERT no banco.
	createdUser, err := repo.AddLogin(ctx, loginCriacao)

	// Verificações básicas
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotZero(t, createdUser.ID) // O ID deve ter sido preenchido pelo banco.

	// Buscando o usuário no banco para garantir que ele foi salvo corretamente.
	var foundUser model.Login
	query := "SELECT id, usuario, senha FROM login WHERE id = @p1"
	err = realDB.QueryRowContext(ctx, query, createdUser.ID).Scan(&foundUser.ID, &foundUser.Nome, &foundUser.Senha)

	require.NoError(t, err, "Falha ao buscar o usuário recém-criado para verificação")
	require.Equal(t, loginCriacao.Nome, foundUser.Nome)
	require.Equal(t, loginCriacao.Senha, foundUser.Senha)
}

func TestIntegration_DeletarLogin(t *testing.T) {
	// A conexão `realDB` já está pronta e a tabela já foi limpa pelo TestMain.
	repo := NewSqlLogin(realDB) // Usamos a conexão REAL!
	ctx := context.Background()

	// --- Cenário 1: Deletar um login que existe (caminho feliz) ---
	t.Run("Sucesso - Deleta um login existente", func(t *testing.T) {
		// ARInserir um registro no banco para poder deletar
		var idParaDeletar int
		insertQuery := "INSERT INTO login (usuario, senha) OUTPUT INSERTED.id VALUES (@p1, @p2)"
		err := realDB.QueryRowContext(ctx, insertQuery, "radinha", "teste.ok").Scan(&idParaDeletar)
		require.NoError(t, err, "Falha ao inserir o registro de setup para o teste de delete")
		require.NotZero(t, idParaDeletar, "O ID retornado pelo insert de setup não pode ser zero")

		// ACT: Executar a função que quero testar.
		err = repo.DeletarLogin(ctx, idParaDeletar)

		
		require.NoError(t, err) // A função de delete não deve retornar erro.

		// A VERIFICAÇÃO FINAL: Vamos ao banco e tentamos encontrar o registro. Ele não deve existir.
		var count int
		verifyQuery := "SELECT COUNT(*) FROM login WHERE id = @p1"
		err = realDB.QueryRowContext(ctx, verifyQuery, idParaDeletar).Scan(&count)
		
		require.NoError(t, err, "Erro ao executar a query de verificação do delete")
		require.Equal(t, 0, count, "O registro não foi deletado do banco, pois a contagem não é zero")
	})

	// --- Cenário 2: Tentar deletar um login que NÃO existe ---
	t.Run("Falha - ID não encontrado", func(t *testing.T) {
		// ARRANGE
		idInexistente := 999999 // Um ID que com certeza não existe.

		// ACT
		err := repo.DeletarLogin(ctx, idInexistente)

		// ASSERT
		require.Error(t, err) // Esperamos que a função retorne um erro.
		
		// Verificamos se a mensagem de erro é a esperada pela nossa regra de negócio.
		expectedErrorMsg := fmt.Sprintf("nenhuma categoria encontrada com o id: %d", idInexistente)
		assert.EqualError(t, err, expectedErrorMsg)
	})
}
func TestIntegration_Login(t *testing.T) {
	// A conexão `realDB` já está pronta e a tabela foi limpa pelo TestMain.
	repo := NewSqlLogin(realDB) // Lembre-se de usar o construtor correto
	ctx := context.Background()

	// --- Cenário 1: Encontrar um usuário que existe (caminho feliz) ---
	t.Run("Sucesso - Encontra um login existente", func(t *testing.T) {
		// ARRANGE: Inserir um usuário no banco para que possamos buscá-lo.
		userParaBuscar := &model.Login{
			Nome:  "Davi",
			Senha: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IKcdjMpklQpKp2wGjwNCoEgoETSmbW",
		}
		insertQuery := "INSERT INTO login (usuario, senha) VALUES (@p1, @p2)"
		_, err := realDB.ExecContext(ctx, insertQuery, userParaBuscar.Nome, userParaBuscar.Senha)
		require.NoError(t, err, "Falha ao inserir o registro de setup para o teste de busca")

		// ACT: Executar a função que queremos testar.
		foundUser, err := repo.Login(ctx, userParaBuscar.Nome)

		// ASSERT: Verificar o resultado.
		require.NoError(t, err) // Não deve haver erro.
		require.NotNil(t, foundUser) // O usuário encontrado não deve ser nulo.
		assert.Equal(t, userParaBuscar.Nome, foundUser.Nome)
		assert.Equal(t, userParaBuscar.Senha, foundUser.Senha)
	})

	// --- Cenário 2: Tentar encontrar um usuário que NÃO existe ---
	t.Run("Falha - Usuário não encontrado", func(t *testing.T) {
		// ARRANGE
		nomeInexistente := "usuario.que.nao.existe"
		// Não precisamos inserir nada, pois o TestMain já limpou a tabela.

		// ACT
		foundUser, err := repo.Login(ctx, nomeInexistente)

		// ASSERT
		require.Error(t, err) // Esperamos que a função retorne um erro.
		assert.Nil(t, foundUser) // Nenhum usuário deve ser retornado.

		// VERIFICAÇÃO CHAVE:
		// A função `Login` deve retornar o erro `sql.ErrNoRows` encapsulado
		// na mensagem de erro customizada. Verificamos se o erro original está presente.
		assert.Contains(t, err.Error(), sql.ErrNoRows.Error(), "A mensagem de erro deveria indicar que o registro não foi encontrado")
	})
}
