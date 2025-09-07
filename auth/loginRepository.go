package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/microsoft/go-mssqldb"
)

type SqlServerLogin struct {
	db *sql.DB
}


func NewSqlLogin(DB *sql.DB) loginRepository {

	return &SqlServerLogin{
		db: DB,
	}
}
// AddLogin implements loginRepository.
//função para adicionar um login no sistema
func (s *SqlServerLogin) AddLogin(ctx context.Context, model *Login) (*Login, error) {
	
	query:= `
			INSERT INTO login (usuario, senha) OUTPUT INSERTED.id values (@p1, @p2);

	` 

	err:= s.db.QueryRowContext(ctx,query, model.Nome, model.Senha).Scan(&model.ID)
	if err != nil {

		var errSql mssql.Error //erro especifico do sqlServer
		if errors.As(err, &errSql) && errSql.Number == 2627{ /*verificando se o erro atual, faz parte dos conjuntos de erro do sqlServer e, 
			verificando se o erro do sqlserver é igual ao numero 2627, que é o erro da constraint UNIQUE*/
			
			return nil, UsuarioJaExistente
		}

		return nil, err
	}
	
	return  model, nil
}

// DeletarLogin implements loginRepository.
func (s *SqlServerLogin) DeletarLogin(ctx context.Context, model Login) error {
	panic("unimplemented")
}

// Login implements loginRepository.
func (s *SqlServerLogin) Login(ctx context.Context, model Login) bool {
	panic("unimplemented")
}

