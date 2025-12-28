package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//arquivo feito para jogar as variaveis de ambiente dentro de uma struct, assim não fica repetindo "os.getEnv", (caso um dia precise)

type VariaveisDeAmbiente struct {

	DB_SERVER string
	DB_USER string
	DB_PORT string
	DATABASE string
	ACCEPT_EULA string
	SA_PASSWORD string
	MSSQL_PID string
}

var Env = NewVariaveisAmbiente()

func NewVariaveisAmbiente() *VariaveisDeAmbiente {

	err:= godotenv.Load("configs/.env", "../configs/.env", "../../configs/.env")
	if err != nil {

		log.Println("Aviso: arquivo .env não encontrado. Continuando com variáveis de sistema...")
	}


	return &VariaveisDeAmbiente{
		DB_SERVER: os.Getenv("DB_SERVER"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PORT: os.Getenv("DB_PORT"),
		DATABASE: os.Getenv("DATABASE"),
		ACCEPT_EULA: os.Getenv("ACCEPT_EULA"),
		SA_PASSWORD: os.Getenv("SA_PASSWORD"),
		MSSQL_PID: os.Getenv("MSSQL_PID"),

	}
}