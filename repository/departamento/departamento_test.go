package departamento

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)



func ConectaAoBanco(t *testing.T)(*sql.DB, func()){


	err:= godotenv.Load("../../configs/.env")
	require.NoError( t, err, "erro ao carregar .env")

	db_server:= os.Getenv("DB_SERVER")
	db_port:= os.Getenv("DB_PORT")
	db_database:= os.Getenv("DATABASE_TESTE")
	db_user:= os.Getenv("DB_USER")
	db_pass:= os.Getenv("SA_PASSWORD")

	if db_database == "" || db_pass == "" || db_port == "" || db_server == "" || db_user == "" {
		t.Skip("pulando teste, variaveis de ambiente nao carregadas")
	}

	connString:= fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", db_user, db_pass, db_server, db_port, db_database)

	db, err:= sql.Open("sqlserver", connString)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err, "erro ao se conectar com o banco de dados")

	ctx:= context.Background()

	t.Log("criando a tabela....")
	tabelaLogin:=`

		use testDb;
		create TABLE login (

    	 id int PRIMARY key IDENTITY(1,1),
   		 usuario VARCHAR(50) unique  not null,
   		 senha VARCHAR(255) not NUll,
);
	` 
	db.ExecContext(ctx, tabelaLogin)

	apagar:= func(){

		ApagaTabelaLogin:= `drop table login; `

		db.ExecContext(ctx, ApagaTabelaLogin)
			require.NoError(t, err, "erro ao apagar a tabela login")

		db.Close()
	
		t.Log("apagando a tabela....")
	}

	

	return  db, apagar
}