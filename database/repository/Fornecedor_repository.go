package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type FornecedorRepository struct {

	q *Queries
	db *pgxpool.Pool
}


func NewFornecedorRepository(pool *pgxpool.Pool) *FornecedorRepository {

	return  &FornecedorRepository{
		q: New(pool),
		db: pool,
	}
}

func (f *FornecedorRepository) Adicionar(ctx context.Context, args CriarFornecedorParams) error {

	err := f.q.CriarFornecedor(ctx, args)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil
}

func (f *FornecedorRepository) ListarFornecedor(ctx context.Context, args ListarFornecedoresParams) ([]ListarFornecedoresRow, error){

	fornecedores, err:= f.q.ListarFornecedores(ctx, args)
	if err != nil {
		return  []ListarFornecedoresRow{}, helper.TraduzErroPostgres(err)
	}

	return fornecedores, nil
}

func (f *FornecedorRepository) CancelarFornecedor(ctx context.Context, args DeletarFornecedorParams) (int64, error){

	linhasAfetadas,err:= f.q.DeletarFornecedor(ctx, args)
	if err != nil {

		return 0,helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas,nil
}

func (f *FornecedorRepository) AtualizaFornecedores(ctx context.Context, args AtualizarFornecedorParams) (int64, error) {

	linhasAfetadas, err:= f.q.AtualizarFornecedor(ctx, args)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return  linhasAfetadas,nil
}