package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type DevolucaoRepository struct {

	q *Queries
	db *pgxpool.Pool
}


func NewDevolucaoRepository(pool *pgxpool.Pool) *DevolucaoRepository{

	return &DevolucaoRepository{
		q: New(pool),
		db: pool,
	} 
}

  
func (d *DevolucaoRepository) AdicionarDevolucao(ctx context.Context,qtx *Queries ,args AddDevolucaoSimplesParams ) error {

	err:= qtx.AddDevolucaoSimples(ctx, args)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil
}

func (d *Devolucao) AdicionarTroca(ctx context.Context,qtx *Queries  ,arg AddTrocaEpiParams) (int32, error) {

	idDevolucao, err:=qtx.AddTrocaEpi(ctx, arg)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return idDevolucao, nil
}

func (d *DevolucaoRepository) EntregaVinculada(ctx context.Context, qtx *Queries ,arg AddEntregaVinculadaParams) (int32, error){

	identrega, err:= qtx.AddEntregaVinculada(ctx, arg)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return identrega, nil
}

func (d *DevolucaoRepository) Cancelar(ctx context.Context,qtx *Queries , arg CancelarDevolucaoParams) (int64, error) {

	linhasAfetadas, err:= qtx.CancelarDevolucao(ctx, arg)
	if err != nil {
		return 0, helper.TraduzErroPostgres(err)
	}

	return  linhasAfetadas, nil
}

func (d *DevolucaoRepository) Listar(ctx context.Context, qtx *Queries ,args ListarDevolucoesParams) ([]ListarDevolucoesRow, error){

	devolucoes, err:= qtx.ListarDevolucoes(ctx, args)
	if err != nil {

		return []ListarDevolucoesRow{}, helper.TraduzErroPostgres(err)
	}

	return devolucoes, nil
}