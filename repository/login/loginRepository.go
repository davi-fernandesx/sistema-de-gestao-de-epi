package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
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
func (s *SqlServerLogin) AddLogin(ctx context.Context, model *model.Login) (*model.Login, error) {
	
	query:= `
			INSERT INTO login (usuario, senha) OUTPUT INSERTED.id values (@p1, @p2);

	` 

	err:= s.db.QueryRowContext(ctx,query, model.Nome, model.Senha).Scan(&model.ID)
	if err != nil {

		var errSql mssql.Error //erro especifico do sqlServer
		if errors.As(err, &errSql) && errSql.Number == 2627{ /*verificando se o erro atual, faz parte dos conjuntos de erro do sqlServer e, 
			verificando se o erro do sqlserver é igual ao numero 2627, que é o erro da constraint UNIQUE*/
			
			return nil, ErrusuarioJaExistente
		}

		return nil, err
	}
	
	return  model, nil
}

// DeletarLogin implements loginRepository.
func (s *SqlServerLogin) DeletarLogin(ctx context.Context, id int) error {

	query:= `

	delete from login where id = @id
	`

	result, err:= s.db.Exec(query, sql.Named("id", id))
	if err != nil {

		return err
	}
	
	row, err:= result.RowsAffected()
	if err != nil{

		return ErrLinhasAfetadas
	}
	
	if row == 0{
		return fmt.Errorf("nenhum login encontrado com o id: %d", id)
	}


	return nil

}

// Login implements loginRepository.
func (s *SqlServerLogin) Login(ctx context.Context, nome string) (*model.Login, error) {
	
	query:= `
		SELECT usuario, senha FROM login WHERE usuario = @nome
	`

   login:= &model.Login{}

	err:= s.db.QueryRow(query, sql.Named("nome", nome)).Scan(&login.Nome,
		&login.Senha,)
		
	if err != nil {
		return nil, fmt.Errorf("erro ao encontrar login: %v", err)
	}

	return login, nil
}

