package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type MotivoDevolucaoRepository struct {
	q *Queries
	db *pgxpool.Pool
}


func NewMotivoDevolucaoRepository(pool *pgxpool.Pool) *MotivoDevolucaoRepository {

	return &MotivoDevolucaoRepository{
		q: New(pool),
		db: pool,
	}
}

func (m *MotivoDevolucaoRepository) Adicionar(ctx context.Context, motivo string) error {

	err := m.q.AddMotivoDevolucao(ctx, motivo)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil
}

func (m *MotivoDevolucaoRepository) ListarMotivo(ctx context.Context, id int) (BuscaMotivoDevolucaoRow, error){

	motivo, err:= m.q.BuscaMotivoDevolucao(ctx, int32(id))
	if err != nil {

		 return BuscaMotivoDevolucaoRow{},helper.TraduzErroPostgres(err)
	}

	return  motivo, err
}

func (m *MotivoDevolucaoRepository) ListarMotivos(ctx context.Context) ([]BuscaTodosMotivosDevolucaoRow, error){

	motivos, err:= m.q.BuscaTodosMotivosDevolucao(ctx)
	if err != nil {

		 return []BuscaTodosMotivosDevolucaoRow{},helper.TraduzErroPostgres(err)
	}

	return  motivos, err
}

func (m *MotivoDevolucaoRepository) CancelarMotivoDevolucao(ctx context.Context, id int) (int64, error) {

	linhasAfetadas,err:= m.q.DeleteMotivoDevolucao(ctx, int32(id))
	if err != nil {

		return 0,helper.TraduzErroPostgres(err)
	}

	return linhasAfetadas, nil
}