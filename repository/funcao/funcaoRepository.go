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

	BuscarTodasFuncao(ctx context.Context) (*[]model.Funcao, error)
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
	panic("unimplemented")
}

// BuscarTodosCargos implements FuncaoInterface.
func (s *SqlServerLogin) BuscarTodasFuncao(ctx context.Context) (*[]model.Funcao, error) {
	panic("unimplemented")
}

// DeletarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) DeletarFuncao(ctx context.Context, id int) error {
	panic("unimplemented")
}
