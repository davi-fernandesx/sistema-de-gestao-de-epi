package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func CorsConfig() gin.HandlerFunc {

	return cors.New(cors.Config{

		//da onde meu back-end aceita requisições
		AllowOrigins: []string{
			"https://sgepi-front-end.vercel.app", //homologação
			"http://localhost:3000", //local
			"https://radaptech.com.br",//homologação
			"https://www.radaptech.com.br",
		},

		AllowMethods: []string{
			"POST","PUT","GET","PATCH", "DELETE", "OPTIONS",
		},//metados permitidos

		AllowHeaders: []string{

			"Origin",
			"Content-type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Tenant-ID",
		},//cabeçalhos que o front pode enviar

		ExposeHeaders: []string{

			"Content-Length",
		}, //informações adicionais que o front pode ler

		AllowCredentials: true,
		MaxAge: 9 * time.Hour, // ate 9 horas para o crs fazer a verificação de novo
	})
}