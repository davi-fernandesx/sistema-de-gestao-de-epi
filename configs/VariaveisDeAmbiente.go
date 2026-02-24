package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//arquivo feito para jogar as variaveis de ambiente dentro de uma struct, assim não fica repetindo "os.getEnv", (caso um dia precise)

type VariaveisDeAmbiente struct {
	DB_SERVER   string
	DB_USER     string
	DB_PORT     string
	DATABASE    string
	DB_PASSWORD string
	DB_SSLMODE  string
}

func NewVariaveisAmbiente() *VariaveisDeAmbiente {

	dir, _ := os.Getwd()
	log.Println("👀 ATENÇÃO: O Go está procurando o arquivo .env dentro desta pasta:", dir)
	// Tenta carregar. Se não achar, avisa UMA VEZ e segue.
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Continuando com variáveis de sistema...")
	}

	sslmode:= os.Getenv("DB_SSLMODE")
	if sslmode == ""{

		sslmode = "disable"
	}

	return &VariaveisDeAmbiente{
		DB_SERVER:   os.Getenv("DB_SERVER"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DATABASE:    os.Getenv("DATABASE"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_SSLMODE: sslmode,
	}
}
