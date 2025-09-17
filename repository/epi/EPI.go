package epi

import (
	"context"
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EpiInterface interface {
	AddEpi(ctx context.Context, epi *model.Epi) error
	DeletarEpi(ctx context.Context, id int) error
	BuscarEpi(ctx context.Context, id int) (*model.Epi, error)
	BuscarTodosEpi(ctx context.Context) (*[]model.Epi, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewEpiRepository(db *sql.DB) EpiInterface {

	return &NewSqlLogin{
		DB: db,
	}
}


// AddEpi implements EpiInterface.
func (n *NewSqlLogin) AddEpi(ctx context.Context, epi *model.Epi) error {
	panic("unimplemented")
}

// BuscarEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {
	panic("unimplemented")
}

// BuscarTodosEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarTodosEpi(ctx context.Context) (*[]model.Epi, error) {
	panic("unimplemented")
}

// DeletarEpi implements EpiInterface.
func (n *NewSqlLogin) DeletarEpi(ctx context.Context, id int) error {
	panic("unimplemented")
}

