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


	
	router:= gin.Default()

	db, err:= configs.InitAplicattion()
	if err != nil {

	log.Fatal(err)
	}

	repoLogin:= login.NewSqlLogin(db)
	ServiceLogin:= service.NewLoginService(repoLogin)
	ControllerLogin:= controller.NewControllerLogin(ServiceLogin)

	router.POST("/Salvar-login", ControllerLogin.SalvarLoginHttp())
	router.GET("/Login", ControllerLogin.AceitarLogin())



	router.Run(":8080")





}