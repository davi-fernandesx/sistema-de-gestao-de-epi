package main

import (
	"log"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/routers"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/middleware"

	"github.com/gin-gonic/gin"
)

// @title           SaaS EPI API
// @version         1.0
// @description     API para gestão de EPIs.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte API
// @contact.url     http://www.seusaas.com.br
// @contact.email   suporte@seusaas.com.br

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main(){

	postgressConnection := configs. ConexaoDbPostgres{}

	init:= configs.Init{Conexao: &postgressConnection,}
	
	router:= gin.Default()

	router.Use(middleware.CorsConfig(),middleware.SecurityHeaders())

	db, err:= init.InitAplicattion()
	if err != nil {

		log.Fatal(err)
	}
	err = postgressConnection.RunMigrationPostgress(db)
	if err != nil {

		log.Fatal(err)
	}

	container:= routers.NewContainer(db)

	routers.ConfigurarRotas(router, container, db)

	router.Run(":8080")

}