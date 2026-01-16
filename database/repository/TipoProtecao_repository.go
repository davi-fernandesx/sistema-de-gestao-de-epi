package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type ProtecaoRepository struct {

	q *Queries
	db *pgxpool.Pool
}

func NewProtecaoRepository(pool *pgxpool.Pool) *ProtecaoRepository {

	return &ProtecaoRepository{
		q: New(pool),
		db: pool,
	}
}

func (p *ProtecaoRepository) Adicionar(ctx context.Context, nome string) error {

	err := p.q.AddProtecao(ctx, nome)
	if err != nil {
		return  helper.TraduzErroPostgres(err)
	}

	return  nil
}

func (p *ProtecaoRepository) ListarProtecao(ctx context.Context, id int32) (BuscarProtecaoRow, error){

	protc, err:= p.q.BuscarProtecao(ctx, id)
	if err != nil {

		return BuscarProtecaoRow{}, helper.TraduzErroPostgres(err)
	}

	return protc, nil
}

func (p *ProtecaoRepository) ListarProtecoes(ctx context.Context) ([]BuscarTodasProtecoesRow, error){

	protc, err:= p.q.BuscarTodasProtecoes(ctx)
	if err != nil {

		return []BuscarTodasProtecoesRow{}, helper.TraduzErroPostgres(err)
	}

	return protc, nil
}

func (p *ProtecaoRepository) CancelarProtecao(ctx context.Context, id int32) (int64, error){

	linhasAfetadas,err:= p.q.DeletarProtecao(ctx, id)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}