package login

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/auth"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
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
   		 usuario VARCHAR(50) not null,
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


func Test_RepoLogin(t *testing.T){

	db, apagar:= ConectaAoBanco(t)
	defer apagar()

	repo:= NewSqlLogin(db)
	ctx:= context.Background()

	t.Run("adicionando um usuario no banco", func(t *testing.T){

		login:= model.Login{
			ID: 1,
			Nome: "rada",
			Senha: "rada2003",
		}

		t.Log("criptografando a senha !!")
		senhaHash, err:= auth.HashPassword(login.Senha)
		require.NoError(t, err, "erro ao criptografar a senha")

		login.Senha = string(senhaHash)

		err = repo.AddLogin(ctx, &login)
		require.NoError(t, err, "erro ao adicionar login")
		t.Log("usuario adicionado!!")


	})
	
}