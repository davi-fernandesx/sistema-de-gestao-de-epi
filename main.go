package main

import (
	"log"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/gin-gonic/gin"
)

func main(){


	
	router:= gin.Default()

	_, err:= configs.InitAplicattion()
	if err != nil {

	log.Fatal(err)
	}


	router.Run(":8080")


	


}