package configs

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)


func TestConexao(t *testing.T) {

	_ = godotenv.Load()

	db_server:= os.Getenv("DB_SERVER")
	db_port:= os.Getenv("DB_PORT")
	db_database:= os.Getenv("DATABASE")
	db_user:= os.Getenv("DB_USER")
	db_pass:= os.Getenv("SA_PASSWORD")


	connString:= fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", db_user, db_pass, db_server, db_port, db_database)

	db,err:= sql.Open("sqlserver", connString)
	if err != nil {

		t.Fatalf("erro ao abrir conexao: %v",err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {

		t.Fatalf("erro ao se conectar com o banco de dados: %v",err)

	}

	t.Log("conexao feita")
}