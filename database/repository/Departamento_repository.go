package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)


type DepartamentoRepository struct {

	q *Queries
	db *pgxpool.Pool
}

func NewDepartamentoRepository(pool *pgxpool.Pool) *DepartamentoRepository {

	return &DepartamentoRepository{
		q: New(pool),
		db: pool,
	}
}

func (d *DepartamentoRepository) Adicionar(ctx context.Context, nome string) error{

	return d.q.CriaDepartamento(ctx, nome)
}

func (d *DepartamentoRepository) ListarDepartamento(ctx context.Context, id int32 ) (BuscarDepartamentoRow, error){

	return d.q.BuscarDepartamento(ctx, id)
}

func (d *DepartamentoRepository) ListarDepartamentos(ctx context.Context)([]BuscarTodosDepartamentosRow, error){

	return  d.q.BuscarTodosDepartamentos(ctx)
}

func (d *DepartamentoRepository) CancelarDepartamento(ctx context.Context, id int32) (int64, error){

	linhasAfetadas, err:= d.q.DeletarDepartamento(ctx, id)

	if err != nil {

		return 0, err
	}

	return  linhasAfetadas, err
}

func (d *DepartamentoRepository) AtualizarDepartamento(ctx context.Context, arg UpdateDepartamentoParams) (int64, error){

	return d.q.UpdateDepartamento(ctx, arg)
}