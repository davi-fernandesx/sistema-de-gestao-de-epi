package tipoprotecao

import (
	"context"
	"database/sql"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type TipoProtecaoInterface interface {
	AddProtecao(ctx context.Context, protecao *model.TipoProtecao) error
	DeletarProtecao(ctx context.Context, ind int) error
	BuscarProtecao(ctx context.Context) (*model.TipoProtecao, error)
	BuscarTodasProtecao(ctx context.Context) (*[]model.TipoProtecao, error)
}

type SqlServerLogin struct {
	DB *sql.DB
}

func NewTipoProtecaoRepository(db *sql.DB) TipoProtecaoInterface {

	return &SqlServerLogin{
		DB: db,
	}
}

// AddProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) AddProtecao(ctx context.Context, protecao *model.TipoProtecao) error {
	panic("unimplemented")
}

// BuscarProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) BuscarProtecao(ctx context.Context) (*model.TipoProtecao, error) {
	panic("unimplemented")
}

// BuscarTodasProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) BuscarTodasProtecao(ctx context.Context) (*[]model.TipoProtecao, error) {
	panic("unimplemented")
}

// DeletarProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) DeletarProtecao(ctx context.Context, ind int) error {
	panic("unimplemented")
}


