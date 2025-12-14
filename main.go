package main

import (
	"log"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/controller"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/login"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service"
	"github.com/gin-gonic/gin"
)

func main(){

	sqlServerConnection := configs.ConexaoDbSqlserver{}

	init:= configs.Init{
		Conexao: &sqlServerConnection,
	}
	
	router:= gin.Default()

	db, err:= init.InitAplicattion()
	if err != nil {

		log.Fatal(err)
	}

	err = sqlServerConnection.RunMigrationSqlserver(db)
	if err != nil {

		log.Fatal(err)
	}

	repoLogin:= login.NewLogin(db)
	ServiceLogin:= service.NewLoginService(repoLogin)
	ControllerLogin:= controller.NewControllerLogin(ServiceLogin)

	router.POST("/Salvar-login", ControllerLogin.SalvarLoginHttp())
	router.GET("/Login", ControllerLogin.AceitarLogin())



	router.Run(":8080")




}