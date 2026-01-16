package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type TamanhosRepository struct {

	q *Queries
	db *pgxpool.Pool
}

func NewTamanhoRepository(pool *pgxpool.Pool) *TamanhosRepository {

	return &TamanhosRepository{q: New(pool), db: pool}
}

func (t *TamanhosRepository) Adicionar(ctx context.Context, tamanho string) error {

	err := t.q.AddTamanho(ctx, tamanho)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil
}

func (t *TamanhosRepository) ListarTamanho(ctx context.Context, id int32) (BuscarTamanhoRow, error){

	tamanho, err:= t.q.BuscarTamanho(ctx, id)
	if err != nil {

		return  BuscarTamanhoRow{}, helper.TraduzErroPostgres(err)
	}

	return tamanho, nil

}

func (t *TamanhosRepository) ListarTamanhos(ctx context.Context) ([]BuscarTodosTamanhosRow, error){

	tamanhos, err:= t.q.BuscarTodosTamanhos(ctx)
	if err != nil {

		return  []BuscarTodosTamanhosRow{}, helper.TraduzErroPostgres(err)
	}

	return tamanhos, nil

}

func (t *TamanhosRepository) CancelarTamanho(ctx context.Context, id int) (int64, error) {

	linhasAfetadas, err:= t.q.DeletarTamanho(ctx, int32(id))
		if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}