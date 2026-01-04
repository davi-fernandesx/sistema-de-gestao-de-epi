package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

type Conexao interface {
	Conn() (*sql.DB, error)
}

type ConexaoDbSqlserverSqlx struct{}
type ConexaoDbSqlserver struct{}

func (C *ConexaoDbSqlserver) Conn() (*sql.DB, error) {

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", Env.DB_USER, Env.SA_PASSWORD, Env.DB_SERVER, Env.DB_PORT, Env.DATABASE)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {

		return nil, fmt.Errorf("erro ao se conectar com o banco de dados: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar se a conexao ainda está ativa: %v", err)
	}

	stats := db.Stats()
	log.Printf("Conexões em uso: %d\n", stats.InUse)
	log.Printf("Conexões ociosas: %d\n", stats.Idle)

	return db, nil

}

func (Cx *ConexaoDbSqlserverSqlx) Conn() (*sqlx.DB, error) {

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", Env.DB_USER, Env.SA_PASSWORD, Env.DB_SERVER, Env.DB_PORT, Env.DATABASE)

	db, err := sqlx.Open("sqlserver", connString)
	if err != nil {

		return nil, fmt.Errorf("erro ao se conectar com o banco de dados: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar se a conexao ainda está ativa: %v", err)
	}

	return db, nil

}

func (C *ConexaoDbSqlserver) RunMigrationSqlserver(db *sql.DB) error {

	driver, err := sqlserver.WithInstance(db, &sqlserver.Config{})
	if err != nil {

		return fmt.Errorf("erro ao iniciar drive da migração, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrate",
		"sqlserver", driver)
	if err != nil {
		return fmt.Errorf("erro ao instanciar migraççao no banco de dados")
	}

	dir, _ := os.Getwd()
	fmt.Println("O programa está rodando na pasta:", dir)
	fmt.Println("Tentando ler migrações de:", dir+"/database/migrate")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao aplicar migrações: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Nenhuma migração nova para aplicar.")
	} else {
		log.Println("Migrações aplicadas com sucesso!")
	}

	log.Println("Migrações aplicadas no banco de dados!....")
	return nil
}

type ConexaoDbMysql struct{}

func (m *ConexaoDbMysql) Conn() (*sql.DB, error) { return nil, nil }
