package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type EntradaRepository struct {

	q *	Queries
	db *pgxpool.Pool
}

func NewEntradaRepository(pool *pgxpool.Pool) *EntradaRepository {

	return &EntradaRepository{
		q: New(pool),
		db: pool,
	}
}

func (e *EntradaRepository) Adicionar(ctx context.Context, args AddEntradaEpiParams) error {

	err:= e.q.AddEntradaEpi(ctx, args)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return err
}



func (e *EntradaRepository) ListarEntradas(ctx context.Context, args ListarEntradasParams) ([]ListarEntradasRow, error) {

	entradas,err:= e.q.ListarEntradas(ctx, args)
	if err != nil {

		return []ListarEntradasRow{}, helper.TraduzErroPostgres(err)
	}

	return entradas, nil

}

func (e *EntradaRepository) CancelarEntrada(ctx context.Context, id int) (int64, error) {


	linhasAfetadas,err:= e.q.CancelarEntrada(ctx, int32(id))
	if err != nil {

		return 0 ,err
	}

	return linhasAfetadas, nil
}

func (e *EntradaRepository) TotalEntradas(ctx context.Context, args ContarEntradasParams) (int64, error){

	total, err:= e.q.ContarEntradas(ctx, args)
	if err != nil {
		return 0, helper.TraduzErroPostgres(err)
	}

	return  total, nil
}