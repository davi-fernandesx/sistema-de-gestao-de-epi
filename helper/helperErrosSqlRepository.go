package helper

import (
	"strings"

	
)

// IsUniqueViolation verifica se o erro é de chave duplicada (SQL Server 2627 ou 2601)
func IsUniqueViolation(err error) bool {
    if err == nil {
        return false
    }

    // Usando interface anônima: qualquer coisa que tenha Number() int32 vai entrar aqui
    if se, ok := err.(interface{ Number() int32 }); ok {
        return se.Number() == 2627 || se.Number() == 2601
    }

    msg := err.Error()
    return strings.Contains(msg, "2627") || strings.Contains(msg, "2601")
}

// IsForeignKeyViolation verifica erro 547
func IsForeignKeyViolation(err error) bool {
	if se, ok := err.(interface{Number() int32}); ok {
		return se.Number() == 547
	}

	msg := err.Error()
	return strings.Contains(msg, "547")
}
