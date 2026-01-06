package repository


import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	// Driver oficial do SQL Server para Go
	_ "github.com/microsoft/go-mssqldb" 
	
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB sobe um container SQL Server novinho para cada teste
func SetupTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	// 1. Configuração do Container SQL Server
	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server:2022-latest",
		ExposedPorts: []string{"1433/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA":     "Y",
			"MSSQL_SA_PASSWORD": "Password123!", // Senha forte obrigatória
		},
		// O SQL Server demora um pouco para subir, esperamos essa mensagem no log
		WaitingFor: wait.ForLog("SQL Server is now ready for client connections").
			WithStartupTimeout(60 * time.Second),
	}

	// 2. Inicia o Container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Falha ao iniciar container SQL Server: %v", err)
	}

	// 3. Garante que o container morra no fim do teste (Cleanup)
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Falha ao terminar container: %v", err)
		}
	})

	// 4. Pega o Host e a Porta mapeada
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "1433")

	// 5. Monta a string de conexão
	// sqlserver://usuario:senha@host:porta?opcoes
	connString := fmt.Sprintf("sqlserver://sa:Password123!@%s:%s?database=master&encrypt=disable", host, port.Port())

	// 6. Abre a conexão
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		t.Fatalf("Falha ao abrir conexão com banco: %v", err)
	}

	// Testa o ping para ter certeza absoluta
	if err = db.Ping(); err != nil {
		t.Fatalf("Banco não respondeu ao ping: %v", err)
	}

	// 7. CRIA AS TABELAS (O Passo Mágico)
	// Como o banco está vazio, precisamos rodar seus CREATE TABLE agora.
	createTables(t, db)

	return db
}

// Essa função roda o script SQL que você me mandou
func createTables(t *testing.T, db *sql.DB) {
	// Dica: Use 'GO' ou separe os comandos se o driver reclamar de batch, 
    // mas geralmente Exec aceita o bloco se não tiver comandos específicos de batch.
    // Aqui separei por ";" para garantir compatibilidade simples.
	schema := `
	CREATE TABLE departamento (
		id INT PRIMARY KEY IDENTITY(1,1),
		nome VARCHAR(100) NOT NULL UNIQUE
	);

	CREATE TABLE funcao (
		id int primary key identity(1,1),
		nome varchar(100) not null unique,
		IdDepartamento int not null,
		foreign key (IdDepartamento) references departamento(id)
	);

	CREATE TABLE tipo_protecao(
		id int primary key identity(1,1),
		nome varchar(100) not null unique
	);

	CREATE TABLE tamanho(
		id int primary key identity(1,1),
		tamanho varchar(50) not null unique
	);

	CREATE TABLE epi(
		id int primary key identity(1,1),
		nome varchar(100) not null,
		fabricante varchar(100) not null,
		CA varchar(20) not null unique,
		descricao text not null,
		validade_CA date not null,
		IdTipoProtecao int not null,
		alerta_minimo INT NOT NULL,
		foreign key (IdTipoProtecao) references tipo_protecao(id)
	);

	CREATE TABLE tamanhos_epis(
		id int primary key identity(1,1),
		IdEpi int not null,
		IdTamanho int not null,
		foreign key (IdEpi) references epi(id),
		foreign key (IdTamanho) references tamanho(id)
	);

	CREATE TABLE funcionario (
		id int primary key identity(1,1),
		nome varchar(100) not null,
		matricula varchar(7) not null unique,
		IdFuncao int not null,
		IdDepartamento int not null,
		foreign key (IdFuncao) references funcao(id),
		foreign key (IdDepartamento) references departamento(id)
	);

    -- ... Adicione o restante das tabelas aqui (entrada_epi, entrega_epi, etc)
    -- Para o exemplo do AddDepartamento, as acima já bastam, mas idealmente ponha todas.
	`
	
	// Executa a criação
	_, err := db.Exec(schema)
	if err != nil {
		t.Fatalf("Falha ao criar schema do banco: %v", err)
	}
}