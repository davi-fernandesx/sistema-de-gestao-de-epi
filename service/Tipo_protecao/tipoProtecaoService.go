package tipoprotecao

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type TipoProtecaoInterface interface {
	AddProtecao(ctx context.Context, protecao *model.TipoProtecao) error
	DeletarProtecao(ctx context.Context, ind int) error
	BuscarProtecao(ctx context.Context, id int) (*model.TipoProtecao, error)
	BuscarTodasProtecao(ctx context.Context) ([]model.TipoProtecao, error)
}
type TipoProtecaoServices struct {
	ProtecaoRepo TipoProtecaoInterface
}
func NewProtecaoServices(repo TipoProtecaoInterface) *TipoProtecaoServices {

	return &TipoProtecaoServices{
		ProtecaoRepo: repo,
	}
}

// SalvarProtecao implements [TipoProtecao].
func (t *TipoProtecaoServices) SalvarProtecao(ctx context.Context, model *model.TipoProtecao) error {
	panic("unimplemented")
}

func (t *TipoProtecaoServices) ListarProtecao(ctx context.Context, id int) (model.TipoProtecaoDto, error) {
	panic("unimplemented")
}

// ListarTodosProteca implements [TipoProtecao].
func (t *TipoProtecaoServices) ListarTodosProteca(ctx context.Context) ([]model.TipoProtecaoDto, error) {
	panic("unimplemented")
}


// AtualizarProtecao implements [TipoProtecao].
func (t *TipoProtecaoServices) AtualizarProtecao(ctx context.Context, id int, NovoTipoProtecao string) error {
	panic("unimplemented")
}

// DeletarProtecao implements [TipoProtecao].
func (t *TipoProtecaoServices) DeletarProtecao(ctx context.Context, id int) error {
	panic("unimplemented")
}

// ListarProtecao implements [TipoProtecao].



