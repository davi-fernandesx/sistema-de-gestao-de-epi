package motivodevolucao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type MotivoDevolucaoInterface interface {
	AddMotivoDevolucao(ctx context.Context, model model.DevolucaoEpi) error
	BuscaMotivoDevolucao()
	BuscaTodosMotivosDevolucao()
	DeleteMotivoDevolucao()
}

type SqlNewLogin struct {
	db *sql.DB
}

func NewMotivoRepository(db *sql.DB) MotivoDevolucaoInterface {

	return &SqlNewLogin{
		db: db,
	}
}

// AddMotivoDevolucao implements [MotivoDevolucaoInterface].
func (s *SqlNewLogin) AddMotivoDevolucao(ctx context.Context, model model.DevolucaoEpi) error {
	
	query:= "insert into motivo_devolucao(motivo) values (@motivo)"

	_, err:= s.db.ExecContext(ctx, query, sql.Named("motivo", model.Motivo))
	if err != nil {

		return fmt.Errorf("erro interno ao salvar motivo de devolucao, %w", err)
	}

	return  nil
}

// BuscaMotivoDevolucao implements [MotivoDevolucaoInterface].
func (s *SqlNewLogin) BuscaMotivoDevolucao() {
	panic("unimplemented")
}

// BuscaTodosMotivosDevolucao implements [MotivoDevolucaoInterface].
func (s *SqlNewLogin) BuscaTodosMotivosDevolucao() {
	panic("unimplemented")
}

// DeleteMotivoDevolucao implements [MotivoDevolucaoInterface].
func (s *SqlNewLogin) DeleteMotivoDevolucao() {
	panic("unimplemented")
}


