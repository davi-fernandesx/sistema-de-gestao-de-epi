package configs

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Init struct {
	Conexao Conexao
}

func (I *Init) InitAplicattion() (*pgxpool.Pool, error) {

	db, err := I.Conexao.Conn()
	if err != nil {

		log.Printf("erro ao carregar o arquivo .env: %v", err)
		log.Printf("aplicação não pode seguir daqui")
		log.Fatal()
	}
	log.Println("---")
	log.Println("Carregando informações do banco de dados.....")
	log.Println("ARQUIVOS .ENV CARREGADOS")
	log.Println("conexao feita com sucesso!!")
	return db, nil
}
