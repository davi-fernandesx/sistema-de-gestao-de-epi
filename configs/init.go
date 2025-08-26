package configs

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)



func InitAplicattion()(*sql.DB, error){

	err:= godotenv.Load("configs/.env")
	if err != nil {
		return  nil, fmt.Errorf("erro ao carregar arquivo .env: %v", err)
	}	
	
	db, err:= Conn()
	if err != nil {

		return  nil, err
	}

	log.Println("conexao feita com sucesso!!")
	return db, nil
}