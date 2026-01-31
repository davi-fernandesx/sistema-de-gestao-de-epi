package configs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Conexao interface {
	Conn(configEnv *VariaveisDeAmbiente) (*pgxpool.Pool, error)
}

type ConexaoDbSqlserverSqlx struct{}
type ConexaoDbSqlserver struct{}
type ConexaoDbPostgres struct {
	Pool *pgxpool.Pool
}

func (p *ConexaoDbPostgres) Conn(configEnv *VariaveisDeAmbiente) (*pgxpool.Pool, error) {

	// Formato: postgres://usuario:senha@host:porta/database
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configEnv.DB_SERVER,
		configEnv.DB_PORT,
		configEnv.DB_USER,
		configEnv.DB_PASSWORD,
		configEnv.DATABASE,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {

		return nil, fmt.Errorf("erro ao configurar pool de banco de dados: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("erro ao se conectar com o banco de dados: %v", err)
	}

	// 4. Ping para verificar se está ativo
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar se a conexao ainda está ativa: %v", err)
	}

	s := pool.Stat()
	log.Printf("Conexões totais: %d | Em uso: %d | Ociosas: %d\n",
		s.TotalConns(), s.AcquiredConns(), s.IdleConns())

	p.Pool = pool
	return pool, nil
}

func (p *ConexaoDbPostgres) RunMigrationPostgress(db *pgxpool.Pool) error {

	sqlDB := stdlib.OpenDBFromPool(db)
	driver, err := pgx.WithInstance(sqlDB, &pgx.Config{})
	if err != nil {

		return fmt.Errorf("erro ao iniciar drive da migração, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrate",
		"postgres", driver)
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
