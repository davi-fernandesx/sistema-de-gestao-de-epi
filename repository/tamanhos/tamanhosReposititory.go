package tamanhos

import (
	"context"
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type TamanhsoInterface interface {
	AddTamanhos(ctx context.Context, tamanhos *model.Tamanho) error
	DeletarTamanhos(ctx context.Context, id int) error
	BuscarTamanhos(ctx context.Context) (*model.Tamanho, error)
	BuscarTodosTamanhos(ctx context.Context) (*[]model.Tamanho, error)
}

type SqlServerLogin struct {
	DB *sql.DB
}

func NewTamanhoRepository(db *sql.DB) TamanhsoInterface {

	return &SqlServerLogin{
		DB: db,
	}
}

// AddTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) AddTamanhos(ctx context.Context, tamanhos *model.Tamanho) error {
	panic("unimplemented")
}

// BuscarTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) BuscarTamanhos(ctx context.Context) (*model.Tamanho, error) {
	panic("unimplemented")
}

// BuscarTodosTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) BuscarTodosTamanhos(ctx context.Context) (*[]model.Tamanho, error) {
	panic("unimplemented")
}

// DeletarTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) DeletarTamanhos(ctx context.Context, id int) error {
	panic("unimplemented")
}
