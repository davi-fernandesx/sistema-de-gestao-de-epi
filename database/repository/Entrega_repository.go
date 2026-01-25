package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type EntregaRepository struct {

	q *Queries
	db *pgxpool.Pool
}


func NewEntregaRepository(pool *pgxpool.Pool) *EntregaRepository {

	return &EntregaRepository{
		 q: New(pool),
		 db: pool,
	}

}


func (e *EntregaRepository) AdicionarEntrega(ctx context.Context, qtx *Queries ,args AddEntregaEpiParams) (int32, error) {

	id,err:= qtx.AddEntregaEpi(ctx, args)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return id, nil
}

func (e *EntregaRepository) AdicionarEntregaItem(ctx context.Context, qtx *Queries, arg AddItemEntregueParams) (int32, error) {


	idEntrega,err:= qtx.AddItemEntregue(ctx, arg)
	if err != nil {

		return  0,helper.TraduzErroPostgres(err)
	}

	return idEntrega, nil
}

func (e *EntregaRepository) ListarEntregas(ctx context.Context, args ListarEntregasParams) ([]ListarEntregasRow, error){


	entregas,err:= e.q.ListarEntregas(ctx, args)
	if err != nil {
		return []ListarEntregasRow{}, helper.TraduzErroPostgres(err)
	}
	return entregas, nil
}

func (e *EntregaRepository) Cancelar(ctx context.Context,qtx *Queries,args CancelarEntregaParams) (int32, error) {

	id,err:= qtx.CancelarEntrega(ctx,args)
	if err != nil {
		return 0, helper.TraduzErroPostgres(err)
	}

	return id, nil
}

func (e *EntregaRepository) CancelarEntregaItem(ctx context.Context, qtx *Queries,id int32) ([]CancelaItemEntregueRow, error) {

	itemsCancelados, err:= qtx.CancelaItemEntregue(ctx, id)
	if err != nil {

		return []CancelaItemEntregueRow{},helper.TraduzErroPostgres(err)
	}

	return  itemsCancelados,nil
}

func (e *EntregaRepository) AbaterEstoqueEntrada(ctx context.Context, qtx *Queries, args AbaterEstoqueLoteParams) (int64, error) {

	linhasAfetadas, err:= qtx.AbaterEstoqueLote(ctx, args)
	if err  != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return  linhasAfetadas, nil
}

func (e *EntregaRepository) ReporEstoqueEntrada(ctx context.Context, qtx *Queries, args ReporEstoqueLoteParams) (int64, error) {

	linhasAfetadas, err := qtx.ReporEstoqueLote(ctx, args)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}
func (e *EntregaRepository) ListarEntregasDisponiveis(ctx context.Context, qtx *Queries, args ListarLotesParaConsumoParams) ([]ListarLotesParaConsumoRow, error){

	lotes, err:= qtx.ListarLotesParaConsumo(ctx, args)
	if err != nil {

		return  []ListarLotesParaConsumoRow{}, helper.TraduzErroPostgres(err)
	}

	return  lotes, nil
}

func (e *EntregaRepository) ListarEpisEntreguesCancelados(ctx context.Context,qtx *Queries ,id int32) ([]ListarItensEntregueCanceladosRow, error){

	cancelados, err:= qtx.ListarItensEntregueCancelados(ctx, id)
	if err != nil {

		return []ListarItensEntregueCanceladosRow{}, err
	}

	return cancelados, nil
}