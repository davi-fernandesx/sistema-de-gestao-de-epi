package main

import (
	"log"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"



	"github.com/gin-gonic/gin"
)

func main(){

	postgressConnection := configs. ConexaoDbPostgres{}

	init:= configs.Init{
		Conexao: &postgressConnection,
	}
	
	router:= gin.Default()

	db, err:= init.InitAplicattion()
	if err != nil {

		log.Fatal(err)
	}

	err = postgressConnection.RunMigrationPostgress(db)
	if err != nil {

		log.Fatal(err)
	}



	router.Run(":8080")




}