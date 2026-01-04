package entrega

import (
	"context"
	"errors"
	"fmt"
	
	"strings"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EntregaInterface interface {
	Addentrega(ctx context.Context, model model.EntregaParaInserir) error
	BuscaEntrega(ctx context.Context, id int) (*model.EntregaDto, error)
	BuscaEntregaPorIdFuncionario(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error)
	BuscaEntregaPorIdFuncionarioCanceladas(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error)
	BuscaTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error)
	BuscaTodasEntregasCanceladas(ctx context.Context) ([]*model.EntregaDto, error)
	BuscaEntregaCancelada(ctx context.Context, id int) (*model.EntregaDto, error)
	CancelarEntrega(ctx context.Context, id int) error
}

type EntregaService struct {
	EntregaRepo EntregaInterface
}

func NewEntregaService(repo EntregaInterface) *EntregaService {

	return &EntregaService{
		EntregaRepo: repo,
	}
}

var (
	errDataMenor           = errors.New("A data de entrada não pode ser menor que hoje")
	ErrFalhaNoBancoDeDados = errors.New("falha no banco de dados")
	ErrId                  = errors.New("id invalido")
	ErrNaoEncontrado       = errors.New("entrega não encontrada")
	ErrErrInterno          = errors.New("erro interno do sistema")
)

func (e *EntregaService) SalvarEntrega(ctx context.Context, model model.EntregaParaInserir) error {

	model.Assinatura_Digital = strings.TrimSpace(model.Assinatura_Digital)

	hoje := time.Now().Truncate(24 * time.Hour)

	if model.Data_entrega.Time().Truncate(24 * time.Hour).Before(hoje) {

		return errDataMenor
	}

	err := e.EntregaRepo.Addentrega(ctx, model)
	if err != nil {

		return ErrFalhaNoBancoDeDados

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

		return model.EntregaDto{}, fmt.Errorf("falha interna: %w", err)
	}

	if entrega == nil {

		return model.EntregaDto{}, ErrErrInterno
	}

	return *entrega, nil
}

func (e *EntregaService) ListarTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error){

	entregas, err:= e.EntregaRepo.BuscaTodasEntregas(ctx)
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

func (e *EntregaService) CancelarEntrega(ctx context.Context, id int) error{

	if id <= 0 {

		return ErrId
	}

	err:= e.EntregaRepo.CancelarEntrega(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrInternal){

			return ErrId
		}

		return  fmt.Errorf("erro crítico ao buscar entrada ID %d: %v", id, err)
	}

	return  nil
}
