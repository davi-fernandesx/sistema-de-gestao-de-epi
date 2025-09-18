package login

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	"github.com/microsoft/go-mssqldb"
)

//interface contendo os metados
type LoginRepository interface{

	 AddLogin( ctx context.Context, model *model.Login) ( error)
	 DeletarLogin(ctx context.Context, id int) error
	 BuscaPorNome(ctx context.Context,  nome string) (*model.Login, error)


}


type SqlServerLogin struct {
	db *sql.DB
}

//construtor
func NewSqlLogin(DB *sql.DB) LoginRepository {

	return &SqlServerLogin{
	db: DB,
}
}
// AddLogin implements loginRepository.
//função para adicionar um login no sistema
func (s *SqlServerLogin) AddLogin( ctx context.Context, model *model.Login) ( error) {
	
	query:= `
			INSERT INTO login (usuario, senha) values (@p1, @p2);

	` 
	_, err:= s.db.ExecContext(ctx, query, model.Nome, model.Senha)
	if err != nil {
		var errSql *mssql.Error
		if errors.As(err, &errSql) && errSql.Number == 2627{
			return  repository.ErrusuarioJaExistente
		}

		return  err
	}
	
	return  nil
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

		return repository.ErrLinhasAfetadas
	}
	
	if row == 0{
		return repository.ErrUsuarioNaoEncontrado
	}


	return nil

}


//busca o usuario pelo nome
func (s *SqlServerLogin) BuscaPorNome(ctx context.Context,  nome string) (*model.Login, error){

	query:= `
		select usuario, senha from login
		where usuario = @usuario ;

	`

	var usuario model.Login

	err:= s.db.QueryRow(query, sql.Named("usuario", nome)).Scan(
		&usuario.Nome,
		&usuario.Senha,

	)

	if err != nil {

		if err == sql.ErrNoRows {
			return  nil, repository.ErrLinhasAfetadas
		}

		return  nil, fmt.Errorf("erro ao buscar usuario; %w", err)
	}


		return  &usuario, nil
	
}

// Login implements loginRepository.


