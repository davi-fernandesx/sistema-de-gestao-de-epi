package middleware

import "github.com/gin-gonic/gin"


func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Adiciona os headers em todas as respostas
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// HSTS (Só ative se estiver usando HTTPS, senão pode bloquear seu localhost)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// No main:
// r.Use(SecurityHeaders())