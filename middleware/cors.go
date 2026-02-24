package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsConfig() gin.HandlerFunc {

	return cors.New(cors.Config{

		//da onde meu back-end aceita requisições
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:3000" || origin == "https://sgepi-front-end.vercel.app" {
				return true
			}

			//dominio principal
			if origin == "https://radaptech.com.br/" || origin == "https://www.radaptech.com.br/" {
				//
				return true
			}

			//permitir qualquer subdominio ex: (frigorificosaojoao.radaptech.com.br)

			if strings.HasSuffix(origin, ".radaptech.com.br"){
				return true
			}

			return false
		},

		AllowMethods: []string{
			"POST", "PUT", "GET", "PATCH", "DELETE", "OPTIONS",
		}, //metados permitidos

		AllowHeaders: []string{

			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Tenant-ID",
		}, //cabeçalhos que o front pode enviar

		ExposeHeaders: []string{

			"Content-Length",
		}, //informações adicionais que o front pode ler

		AllowCredentials: true,
		MaxAge:           9 * time.Hour, // ate 9 horas para o crs fazer a verificação de novo
	})
}
