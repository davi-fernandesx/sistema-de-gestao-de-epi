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
	DB_URL      string
}

func NewVariaveisAmbiente() *VariaveisDeAmbiente {
    // Tenta carregar. Se não achar, avisa UMA VEZ e segue.
    err := godotenv.Load(".env", "../.env", "../../.env")
    if err != nil {
        log.Println("Aviso: arquivo .env não encontrado. Continuando com variáveis de sistema...")
    }

	return &VariaveisDeAmbiente{
		DB_SERVER:   os.Getenv("DB_SERVER"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DATABASE:    os.Getenv("DATABASE"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_URL:      os.Getenv("DATABASE_URL"),
	}
}
