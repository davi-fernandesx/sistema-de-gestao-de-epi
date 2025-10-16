package configs

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/microsoft/go-mssqldb"
)

type ConexaoDbSqlserver struct {}
type Conexao interface {

	Conn()( *sql.DB, error)
}

func (C *ConexaoDbSqlserver) Conn()(*sql.DB, error) {

	db_server:= os.Getenv("DB_SERVER")
	db_port:= os.Getenv("DB_PORT")
	db_database:= os.Getenv("DATABASE")
	db_user:= os.Getenv("DB_USER")
	db_pass:= os.Getenv("SA_PASSWORD")


	connString:= fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", db_user, db_pass, db_server, db_port, db_database)

	db, err:= sql.Open("sqlserver",connString)
	if err != nil {

		return  nil, fmt.Errorf("erro ao se conectar com o banco de dados: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return  nil, fmt.Errorf("erro ao verificar se a conexao ainda está ativa: %v", err)
	}

	return  db, nil


}


type ConexaoDbMysql struct{}

func (m *ConexaoDbMysql) Conn()(*sql.DB, error){return  nil, nil}