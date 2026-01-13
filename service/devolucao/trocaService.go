package devolucao

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)


type DevolucaoInterfaceRepository interface {
	AddTrocaEPI(ctx context.Context, devolucao model.DevolucaoInserir) error
	AddDevolucaoEpi(ctx context.Context, devolucao model.DevolucaoInserir) error
	DeleteDevolucao(ctx context.Context, id int) error
	BuscaDevolucaoPorMatricula(ctx context.Context, matricula int) ([]model.Devolucao, error)
	BuscaTodasDevolucoes(ctx context.Context) ([]model.Devolucao, error)
	//BaixaEstoque(ctx context.Context, tx *sql.Tx, idEpi, iDTamanho int64, quantidade int, idEntrega int64) error
	BuscaDevolucaoPorId(ctx context.Context, id int)([]model.Devolucao,error)
	BuscaDevolucaoPorIdCancelada(ctx context.Context, id int)([]model.Devolucao, error)
	BuscaDevolucaoPorMatriculaCancelada(ctx context.Context, matricula int) ([]model.Devolucao, error)
	BuscaTodasDevolucoesCancelada(ctx context.Context) ([]model.Devolucao, error)
}


type TrocaService struct {

	TrocaRepo DevolucaoInterfaceRepository
}

func NewTrocaService(repo DevolucaoInterfaceRepository) *TrocaService {

	return &TrocaService{
		TrocaRepo: repo,
	}
}