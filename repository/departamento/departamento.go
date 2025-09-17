package departamento

import (
	"context"
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type DepartamentoInterface interface {
	AddDepartamento(ctx context.Context, departamento *model.Departamento) error
	DeletarDepartamento(ctx context.Context, id int) error
	BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error)
	BuscarTodosDepartamentos(ctx context.Context) (*[]model.Departamento, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewDepartamentoRepository(db *sql.DB) DepartamentoInterface {

	return &NewSqlLogin{
		DB: db,
	}
}


// AddDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) AddDepartamento(ctx context.Context, departamento *model.Departamento) error {
	panic("unimplemented")
}

// BuscarDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error) {
	panic("unimplemented")
}

// BuscarTodosDepartamentos implements DepartamentoInterface.
func (n *NewSqlLogin) BuscarTodosDepartamentos(ctx context.Context) (*[]model.Departamento, error) {
	panic("unimplemented")
}

// DeletarDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) DeletarDepartamento(ctx context.Context, id int) error {
	panic("unimplemented")
}

