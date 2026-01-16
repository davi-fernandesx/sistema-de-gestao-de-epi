package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)



type FuncionarioRepository struct {

	q *Queries
	db *pgxpool.Pool
}

func NewFuncionarioRepository(pool *pgxpool.Pool) *FuncaoRepository {

	return &FuncaoRepository{
		q: New(pool),
		db: pool,
	}
}

func (f *FuncionarioRepository) Adicionar(ctx context.Context, args AddFuncionarioParams) error {

	err :=f.q.AddFuncionario(ctx, args)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil	
}

func(f *FuncionarioRepository) ListarFuncionario(ctx context.Context, matricula string) (BuscaFuncionarioRow, error){

	funcionario, err:= f.q.BuscaFuncionario(ctx, matricula)
	if err != nil {

		return  BuscaFuncionarioRow{},helper.TraduzErroPostgres(err)
	}

	return funcionario, nil
}

func (f *FuncionarioRepository) ListarFuncionarios(ctx context.Context)([]BuscarTodosFuncionariosRow, error) {

	funcs, err:= f.q.BuscarTodosFuncionarios(ctx)
	if err != nil {

		return []BuscarTodosFuncionariosRow{}, helper.TraduzErroPostgres(err)
	}

	return  funcs, nil
}

func (f *FuncionarioRepository) CancelarFuncionario(ctx context.Context, id int32) (int64, error){

	linhasAfetadas,err:= f.q.DeletarFuncionario(ctx, id)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}

func (f *FuncionarioRepository) AtualizarFuncionarioNome(ctx context.Context, arg UpdateFuncionarioNomeParams) (int64, error){

	linhasAfetadas,err:= f.q.UpdateFuncionarioNome(ctx, arg)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}

func (f *FuncionarioRepository) AtualizarFuncionarioDepartamento(ctx context.Context, arg UpdateFuncionarioDepartamentoParams) (int64, error){

	linhasAfetadas,err:= f.q.UpdateFuncionarioDepartamento(ctx, arg)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}

func (f *FuncionarioRepository) AtualizarFuncionarioFuncao(ctx context.Context, arg UpdateFuncionarioFuncaoParams) (int64, error){

	linhasAfetadas,err:= f.q.UpdateFuncionarioFuncao(ctx, arg)
	if err != nil {

		return 0, helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}