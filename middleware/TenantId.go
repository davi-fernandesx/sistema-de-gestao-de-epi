package middleware

import (
	"database/sql"
	"net"
	"net/http"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/gin-gonic/gin"
)
const TenantId = "tenantId"


func TenantMiddleware(querie *repository.Queries)  gin.HandlerFunc {

	return func(c *gin.Context) {

		host:= c.Request.Host

		// Remove a porta se existir (ex: localhost:8080 -> localhost)
		if strings.Contains(host, ":"){
			var err error
			host, _, err = net.SplitHostPort(host)
			if err != nil {

				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Host inválido"})
				return
			}
		}

		parts:= strings.Split(host, ".")
		// Lógica simples: Pega a primeira parte como subdomínio.
		// Ex: "frigorifico.radap.com.br" -> parts[0] = "frigorifico"
		// CUIDADO: Em localhost puro ("localhost"), parts tem tamanho 1.
		// Se usar lvh.me ("teste.lvh.me"), parts tem tamanho 3.

		var subdominio string

		if len(parts) >= 3 {
			subdominio = parts[0]
		} else {
			// Caso esteja rodando local sem subdomínio ou ip direto, bloqueia ou define um padrão
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Subdomínio não identificado"})
			return
		}

		// Ignorar 'www' ou 'api' se forem reservados
		if subdominio == "www" || subdominio == "api" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Subdomínio reservado"})
			return
		}


		// Busca no banco (usando o Context da request para cancelamento/timeout)
		empresa,err:= querie.GetTenantBySubdomain(c.Request.Context(), subdominio)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Empresa não encontrada"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao validar empresa"})
			return
		}
		// SUCESSO: Salva o ID dentro do contexto do GIN
		// O Gin tem um mapa chave/valor interno otimizado (c.Set/c.Get)
		c.Set(TenantId, empresa.ID)

		// Passa para o próximo handler
		c.Next()
	}
}

// Helper para pegar o ID dentro dos Controllers de forma tipada
func GetTenantID(c *gin.Context) (int32, bool) {
	val, exists := c.Get(TenantId)
	if !exists {
		return 0, false
	}
	// Faz o cast para int32 (o tipo que o sqlc geralmente usa para IDs)
	id, ok := val.(int32) 
	return id, ok
}