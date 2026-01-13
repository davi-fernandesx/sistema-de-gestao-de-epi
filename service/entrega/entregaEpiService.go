package entrega

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"strings"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	estoque "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Estoque"
)

//go:generate mockery --name=EntregaInterface --output=mocks --outpkg=mocks
type EntregaInterface interface {
	Addentrega(ctx context.Context, tx *sql.Tx,model model.EntregaParaInserir) (int64, error)
	BuscaEntrega(ctx context.Context, id int) (*model.EntregaDto, error)
	BuscaEntregaPorIdFuncionario(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error)
	BuscaEntregaPorIdFuncionarioCanceladas(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error)
	BuscaTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error)
	BuscaTodasEntregasCanceladas(ctx context.Context) ([]*model.EntregaDto, error)
	BuscaEntregaCancelada(ctx context.Context, id int) (*model.EntregaDto, error)
	CancelarEntrega(ctx context.Context, id int) error
}

type EntregaService struct {

	db *sql.DB
	EntregaRepo EntregaInterface
	baixaEstoque estoque.BaixaEstoque

}

func NewEntregaService(Db *sql.DB, repo EntregaInterface, es estoque.BaixaEstoque) *EntregaService {

	return &EntregaService{
		db: Db,
		EntregaRepo: repo,
		baixaEstoque: es,

	}
}

var (
	errDataMenor           = errors.New("A data de entrada não pode ser menor que hoje")
	ErrFalhaNoBancoDeDados = errors.New("falha no banco de dados")
	ErrId                  = errors.New("id invalido")
	ErrNaoEncontrado       = errors.New("entrega não encontrada")
	ErrInterno          = errors.New("erro interno do sistema")
)

func (e *EntregaService) SalvarEntrega(ctx context.Context, model model.EntregaParaInserir) error {

	model.Assinatura_Digital = strings.TrimSpace(model.Assinatura_Digital)

	hoje := time.Now().Truncate(24 * time.Hour)

	if model.Data_entrega.Time().Truncate(24 * time.Hour).Before(hoje) {

		return errDataMenor
	}

	tx, err:= e.db.BeginTx(ctx, nil)
	if err != nil {

		return err
	}

	defer tx.Rollback()

	entregaId, err:= e.EntregaRepo.Addentrega(ctx, tx, model)
	if err != nil {
		return  err
	}

	for _, item := range model.Itens{

		if item.Quantidade <= 0 {

			continue
		}

		entradas, err:= e.baixaEstoque.ListarLotesParaConsumo(ctx, tx, int64(item.ID_epi), int64(item.ID_tamanho))
		if err != nil {

			return fmt.Errorf("erro ao buscar lotes para o EPI %d: %v", item.ID_epi, err)
		}
		 
	    quantidadeRestante :=  item.Quantidade

		for _, entrada := range entradas{

			if quantidadeRestante == 0 {
				break
			}

			var quantidadeParaAbater int

			if entrada.Quantidade >= quantidadeRestante{
				quantidadeParaAbater = quantidadeRestante
				quantidadeRestante = 0 
			}else {

				quantidadeParaAbater = entrada.Quantidade
				quantidadeRestante = quantidadeRestante - entrada.Quantidade
			}

			err := e.baixaEstoque.AbaterEstoqueLote(ctx, tx, entrada.ID, quantidadeParaAbater)
			if err != nil {

				return fmt.Errorf("erro ao baixar estoque do lote %d: %v", entrada.ID, err)
			}

			err = e.baixaEstoque.RegistrarItemEntrega(ctx, tx, int64(item.ID_epi), int64(item.ID_tamanho), quantidadeParaAbater, 
			entregaId, entrada.ID,entrada.ValorUnitario)
			if err != nil {
				return fmt.Errorf("erro ao registrar item da entrega: %v", err)

			}

		}

			if quantidadeRestante > 0 {

				return fmt.Errorf("estoque insuficiente para o EPI ID %d (Tamanho %d). Faltam %d unidades", item.ID_epi, item.ID_tamanho, quantidadeRestante)
			}
	}

	if err := tx.Commit(); err != nil {
        return fmt.Errorf("erro ao fazer commit da transação: %v", err)
    }

    return nil
}

func (e *EntregaService) ListaEntrega(ctx context.Context, id int) (model.EntregaDto, error) {

	if id <= 0 {

		return model.EntregaDto{}, ErrId
	}

	entrega, err := e.EntregaRepo.BuscaEntrega(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrBuscarTodos) {

			return model.EntregaDto{}, ErrNaoEncontrado
		}

		return model.EntregaDto{}, fmt.Errorf("falha interna: %w", ErrInterno)
	}

	if entrega == nil {

		return model.EntregaDto{}, ErrInterno
	}

	return *entrega, nil
}

func (e *EntregaService) ListarTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error) {

	entregas, err := e.EntregaRepo.BuscaTodasEntregas(ctx)
	if err != nil {

		if errors.Is(err, Errors.ErrBuscarTodos) {

			return []*model.EntregaDto{}, nil
		}

		return []*model.EntregaDto{}, fmt.Errorf("falha interna: %w", err)
	}

	if entregas == nil {
		return []*model.EntregaDto{}, nil
	}

	return entregas, nil

}

func (e *EntregaService) CancelarEntrega(ctx context.Context, id int) error {

	if id <= 0 {

		return ErrId
	}

	err := e.EntregaRepo.CancelarEntrega(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrInternal) {

			return ErrId
		}

		return fmt.Errorf("erro crítico ao buscar entrada ID %d: %v", id, err)
	}

	return nil
}
