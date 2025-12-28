package tipoprotecao

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/tipo_protecao"
)

type TipoProtecao interface {
	SalvarProtecao(ctx context.Context, model *model.TipoProtecao) error
	ListarProtecao(ctx context.Context, id int) (model.TipoProtecaoDto, error)
	ListarTodosProteca(ctx context.Context) ([]model.TipoProtecaoDto, error)
	DeletarProtecao(ctx context.Context, id int) error
	AtualizarProtecao(ctx context.Context, id int, NovoTipoProtecao string) error
}

type TipoProtecaoServices struct {
	ProtecaoRepo tipoprotecao.TipoProtecaoInterface
}
func NewProtecaoServices(repo tipoprotecao.TipoProtecaoInterface) TipoProtecao {

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



