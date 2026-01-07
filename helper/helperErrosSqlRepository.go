package helper

import (
	"errors"

	mssql "github.com/microsoft/go-mssqldb"
)

// IsUniqueViolation verifica se o erro é de chave duplicada (SQL Server 2627 ou 2601)
// IsForeignKeyViolation verifica se o erro é o código 547 do SQL Server
func IsForeignKeyViolation(err error) bool {
	var sqlErr mssql.Error
	
	// errors.As tenta encontrar um erro do tipo mssql.Error dentro da cadeia de erros
	if errors.As(err, &sqlErr) {
		// 547 = The INSERT statement conflicted with the FOREIGN KEY constraint
		return sqlErr.Number == 547
	}
	
	return false
}

// Aproveite e corrija o de Unique também (códigos 2601 e 2627)
func IsUniqueViolation(err error) bool {
	var sqlErr mssql.Error
	
	if errors.As(err, &sqlErr) {
		// 2601 = Cannot insert duplicate key row ... with unique index
		// 2627 = Violation of UNIQUE KEY constraint
		return sqlErr.Number == 2601 || sqlErr.Number == 2627
	}
	
	return false
}