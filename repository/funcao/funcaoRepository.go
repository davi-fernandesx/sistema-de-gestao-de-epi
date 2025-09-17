package funcao

import (
	"context"
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type FuncaoInterface interface {
	AddFuncao(ctx context.Context, funcao *model.Funcao) error
	DeletarFuncao(ctx context.Context, id int) error
	BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error)
	BuscarTodosCargos(ctx context.Context) (*[]model.Funcao, error)
	
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
	panic("unimplemented")
}

// BuscarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error) {
	panic("unimplemented")
}

// BuscarTodosCargos implements FuncaoInterface.
func (s *SqlServerLogin) BuscarTodosCargos(ctx context.Context) (*[]model.Funcao, error) {
	panic("unimplemented")
}

// DeletarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) DeletarFuncao(ctx context.Context, id int) error {
	panic("unimplemented")
}


