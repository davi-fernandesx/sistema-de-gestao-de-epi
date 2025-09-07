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

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb" // Driver do SQL Server
	"github.com/stretchr/testify/require"
)

var realDB *sql.DB // Variável global para ser usada pelos testes neste pacote

// TestMain é a função de setup e teardown para os testes de integração.
func TestMain(m *testing.M) {
	// 1. Obter a string de conexão de uma variável de ambiente.
	// Isso evita deixar senhas no código.
	errr := godotenv.Load("../configs/.env")
	if errr != nil {
		log.Fatal("erro ao abrir arquivo .env")
	}
	db_server:= os.Getenv("DB_SERVER")
	db_port:= os.Getenv("DB_PORT")
	db_database:= os.Getenv("DATABASE")
	db_user:= os.Getenv("DB_USER")
	db_pass:= os.Getenv("SA_PASSWORD")


	connString:= fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", db_user, db_pass, db_server, db_port, db_database)
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
	loginCriacao := &Login{
		Nome:          "Davi",
		Senha: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IKcdjMpklQpKp2wGjwNCoEgoETSmbW",
	}

	// Chamando o método que queremos testar. Isso vai realmente fazer um INSERT no banco.
	createdUser, err := repo.AddLogin(ctx, loginCriacao)

	// Verificações básicas
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotZero(t, createdUser.ID) // O ID deve ter sido preenchido pelo banco.

	// Buscando o usuário no banco para garantir que ele foi salvo corretamente.
	var foundUser Login
	query := "SELECT id, usuario, senha FROM login WHERE id = @p1"
	err = realDB.QueryRowContext(ctx, query, createdUser.ID).Scan(&foundUser.ID, &foundUser.Nome, &foundUser.Senha)
	
	require.NoError(t, err, "Falha ao buscar o usuário recém-criado para verificação")
	require.Equal(t, loginCriacao.Nome, foundUser.Nome)
	require.Equal(t, loginCriacao.Senha, foundUser.Senha)
}

