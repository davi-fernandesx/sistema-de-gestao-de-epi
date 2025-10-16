package configs

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)


type Init struct {

	Conexao Conexao
}


func (I *Init)InitAplicattion()(*sql.DB, error){

	err:= godotenv.Load("configs/.env")
	if err != nil {
		return  nil, fmt.Errorf("erro ao carregar arquivo .env: %v", err)
	}	
	
	db, err:= I.Conexao.Conn()
	if err != nil {

		return  nil, err
	}

	

	log.Println("conexao feita com sucesso!!")
	return db, nil
}