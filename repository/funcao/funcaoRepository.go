package funcao

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/microsoft/go-mssqldb"
)

type FuncaoInterface interface {
	AddFuncao(ctx context.Context, funcao *model.Funcao) error
	DeletarFuncao(ctx context.Context, id int) error
	BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error)

	BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error)
}

type SqlServerLogin struct {
	Db *sql.DB
}

func NewfuncaoRepository(db *sql.DB) FuncaoInterface {

	return &SqlServerLogin{
		Db: db,
	}
}

// AddFuncao implements FuncaoInterface.
func (s *SqlServerLogin) AddFuncao(ctx context.Context, funcao *model.Funcao) error {
		query:= `insert into funcao (funcao) values (@funcao)`

		_, err:= s.Db.ExecContext(ctx, query, sql.Named("funcao",funcao.Funcao))
		if err != nil{
			var ErrSql *mssql.Error
			if errors.As(err, &ErrSql) && ErrSql.Number == 2627 {
				return  repository.ErrFuncaoJaExistente
			}

			return  err
		}

		return  nil
}

// BuscarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error) {
	
	query:= `select id, funcao from funcao where id = @id`

	var funcao model.Funcao

	err:= s.Db.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&funcao.ID, &funcao.Funcao)

	if err != nil {
		if err == sql.ErrNoRows {
			return  nil, repository.ErrAoProcurarFuncao
		}

		return  nil, repository.ErrFalhaAoEscanearDados
	}

	return &funcao, nil
}

// BuscarTodosCargos implements FuncaoInterface.
func (s *SqlServerLogin) BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error) {
	query:= `select id, funcao from funcao`

	linhas, err:= s.Db.QueryContext(ctx, query)
	if err != nil {
		return []model.Funcao{}, repository.ErrAoBuscarTodasAsFuncoes
	}

	var Funcoes []model.Funcao
	defer linhas.Close()

	for linhas.Next(){

		var funcao model.Funcao

		if err:= linhas.Scan(&funcao.ID, &funcao.Funcao); err!= nil {

			return  nil, repository.ErrFalhaAoEscanearDados
		}

		Funcoes = append(Funcoes, funcao)
	}
	
	err = linhas.Err()
	if err != nil {
		return  nil, repository.ErrAoIterarSobreFuncoes
	}

	return  Funcoes, nil

}

// DeletarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) DeletarFuncao(ctx context.Context, id int) error {
	
	query:= `delete from funcao where id = @id`

	result, err:= s.Db.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  err
	}

	
	linhas, err:= result.RowsAffected()
	if err != nil {
		return repository.ErrLinhasAfetadas
	}

	if linhas == 0 {

		return repository.ErrAoProcurarFuncao
	}

	return nil

}
