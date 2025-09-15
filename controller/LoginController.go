package controller

import (
	"net/http"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service"
	"github.com/gin-gonic/gin"
)

type ControllerLogin struct {

	ServiceLogin *service.LoginService
}


func NewControllerLogin(serviceLogin *service.LoginService) *ControllerLogin{

	return &ControllerLogin{
		ServiceLogin: serviceLogin,
	}
}


func (Cl *ControllerLogin) SalvarLoginHttp() gin.HandlerFunc {

	return  func(c *gin.Context) {

		var  request model.LoginDto
		ctx:= c
		err:= c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro": "erro ao realizar a decodificação do json"})
			return 
		}


		err = Cl.ServiceLogin.SalvaLogin(ctx, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro:": err.Error()})
			return 
		}


		c.JSON(http.StatusOK, gin.H{
			"message": "usuario salvo com sucesso",
		})
	}
}

func (Cl *ControllerLogin) AceitarLogin() gin.HandlerFunc{

	return  func(c *gin.Context) {

		var request model.LoginDto
		ctx:= c
		err:= c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": err.Error(),
			})
			return 
		}

		

		loginAceito, err:= Cl.ServiceLogin.Login(ctx, request)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return 
		}

		if loginAceito {
			c.JSON(http.StatusOK, gin.H{"Login feito com sucesso": "ok"})
			return 
		}


		c.JSON(http.StatusUnauthorized, gin.H{"erro":"senha ou usuario incorreto"})
		

	}
}