package integracao

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)




func SetupTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	// Senha forte definida em uma variável para garantir que é a mesma em todo lugar
	const dbPassword = "Password123!7645!!!iiJ"

	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server:2022-latest",
		ExposedPorts: []string{"1433/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA":       "Y",
			"MSSQL_SA_PASSWORD": dbPassword,
		},
		WaitingFor: wait.ForLog("SQL Server is now ready for client connections").
			WithStartupTimeout(90 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Erro fatal ao subir Docker: %v", err)
	}

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "1433")

	// --- AQUI ESTÁ A CORREÇÃO ---
	// Usamos url.URL para construir a string sem erros de caracteres especiais
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("sa", dbPassword), // Isso trata o "!" corretamente
		Host:   fmt.Sprintf("%s:%s", host, port.Port()),
	}

	// Adiciona os parâmetros extras
	q := u.Query()
	q.Set("database", "master")
	q.Set("encrypt", "disable")
	u.RawQuery = q.Encode()

	// Conecta usando a URL gerada (ex: sqlserver://sa:Password123%21@localhost:49154...)
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		t.Fatalf("Erro ao abrir conexão: %v", err)
	}

	// Tenta pingar com retentativa (SQL Server as vezes demora 1s extra pra aceitar login)
	if err := pingComRetentativa(db); err != nil {
		t.Fatalf("Banco subiu mas não aceitou login: %v", err)
	}

	// Cria tabelas
	criarTabelasDoUsuario(t, db)

	return db
}

// Função auxiliar para insistir no login nos primeiros segundos
func pingComRetentativa(db *sql.DB) error {
	var err error
	for i := 0; i < 10; i++ { // Tenta 10 vezes
		err = db.Ping()
		if err == nil {
			return nil // Sucesso!
		}
		time.Sleep(500 * time.Millisecond) // Espera meio segundo
	}
	return err
}


