package entrada

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

//go:generate mockery --name=EntradaRepository --output=mocks --outpkg=mocks
type EntradaRepository interface {
	AddEntradaEpi(ctx context.Context, model *model.EntradaEpiInserir) error
	BuscarEntrada(ctx context.Context, id int) (model.EntradaEpi, error)
	BuscarTodasEntradas(ctx context.Context) ([]model.EntradaEpi, error)
	CancelarEntrada(ctx context.Context, id int) error
}

type EntradaService struct {
	EntradaRepo EntradaRepository
}

func NewEntradaService(repo EntradaRepository) *EntradaService {

	return &EntradaService{

		EntradaRepo: repo,
	}
}

var (
	errDataMenor           = errors.New("A data de entrada não pode ser menor que hoje")
	errDataMenorValidade   = errors.New("A data de validade não pode ser menor que a data de fabricação")
	ErrNaoCadastrado       = errors.New("epi ou tamanho não esta cadastrado no sistema")
	ErrFalhaNoBancoDeDados = errors.New("falha no banco de dados")
	ErrId                  = errors.New("id invalido")
	ErrNaoEncontrado       = errors.New("entrada não encontrada")
	ErrDataIgual           = errors.New("data de fabricacao e validade não podem ser iguais")
)

// SalvarEntrada implements [entrada].
func (e *EntradaService) SalvarEntrada(ctx context.Context, model *model.EntradaEpiInserir) error {

	model.Lote = strings.TrimSpace(model.Lote)
	model.Fornecedor = strings.TrimSpace(model.Fornecedor)

	//verificando se a data é menor do que hoje
	hoje := time.Now().Truncate(24 * time.Hour)
	if model.Data_entrada.Time().Truncate(24 * time.Hour).Before(hoje) {

		return errDataMenor
	}

	if model.DataValidade.Time().Equal(model.DataFabricacao.Time()) {

		return ErrDataIgual
	}

	if model.DataValidade.Time().Before(model.DataFabricacao.Time()) {
		return errDataMenorValidade
	}

	err := e.EntradaRepo.AddEntradaEpi(ctx, model)
	if err != nil {

		log.Printf("erro ao salvar entrada: %v", err)
		if errors.Is(err, Errors.ErrDadoIncompativel) {

			return ErrNaoCadastrado
		}

		return ErrFalhaNoBancoDeDados
	}

	return nil
}

// ListarEntrada implements [entrada].
func (e *EntradaService) ListarEntrada(ctx context.Context, id int) (model.EntradaEpiDto, error) {

	if id <= 0 {

		return model.EntradaEpiDto{}, ErrId
	}

	entrada, err := e.EntradaRepo.BuscarEntrada(ctx, id)
	if err != nil {

		log.Printf("erro ao buscar entrada: %v", err)
		if errors.Is(err, Errors.ErrBuscarTodos) {

			return model.EntradaEpiDto{}, nil
		}

		log.Printf("erro crítico ao buscar entrada ID %d: %v", id, err)
		return model.EntradaEpiDto{}, ErrFalhaNoBancoDeDados
	}

	return model.EntradaEpiDto{
		ID: entrada.ID,
		Epi: model.EpiDto{
			Id:         entrada.ID_epi,
			Nome:       entrada.Nome,
			Fabricante: entrada.Fabricante,
			CA:         entrada.CA,
			Tamanho: []model.TamanhoDto{
				{
					ID:      entrada.Id_Tamanho,
					Tamanho: entrada.TamanhoDescricao,
				},
			},
			Descricao:      entrada.Descricao,
			DataValidadeCa: entrada.DataValidadeCa,
			Protecao: model.TipoProtecaoDto{
				ID:   entrada.IDprotecao,
				Nome: model.Protecao(entrada.NomeProtecao),
			},
		},
		Data_entrada:  entrada.Data_entrada,
		Quantidade:    entrada.Quantidade,
		Lote:          entrada.Lote,
		Fornecedor:    entrada.Fornecedor,
		ValorUnitario: entrada.ValorUnitario,
	}, nil
}

// ListasTodasEntradas implements [entrada].
func (e *EntradaService) ListasTodasEntradas(ctx context.Context) ([]model.EntradaEpiDto, error) {

	entradas, err := e.EntradaRepo.BuscarTodasEntradas(ctx)
	if err != nil {

		if err != nil {

			if errors.Is(err, Errors.ErrBuscarTodos) {

			return []model.EntradaEpiDto{}, nil
			}

			return []model.EntradaEpiDto{}, fmt.Errorf("falha interna: %w", err)
		}
	}



	dto := make([]model.EntradaEpiDto, 0, len(entradas))

	for _, entrada := range entradas {

		e := model.EntradaEpiDto{
			ID: entrada.ID,
			Epi: model.EpiDto{
				Id:         entrada.ID_epi,
				Nome:       entrada.Nome,
				Fabricante: entrada.Fabricante,
				CA:         entrada.CA,
				Tamanho: []model.TamanhoDto{
					{
						ID:      entrada.Id_Tamanho,
						Tamanho: entrada.TamanhoDescricao,
					},
				},
				Descricao:      entrada.Descricao,
				DataValidadeCa: entrada.DataValidadeCa,
				Protecao: model.TipoProtecaoDto{
					ID:   entrada.IDprotecao,
					Nome: model.Protecao(entrada.NomeProtecao)},
			},
			Data_entrada:  entrada.Data_entrada,
			Quantidade:    entrada.Quantidade,
			Lote:          entrada.Lote,
			Fornecedor:    entrada.Fornecedor,
			ValorUnitario: entrada.ValorUnitario,
		}

		dto = append(dto, e)
	}

	return dto, nil
}

// DeletarEntradas implements [entrada].
func (e *EntradaService) DeletarEntradas(ctx context.Context, id int) error {

	if id <= 0 {

		return ErrId
	}

	err := e.EntradaRepo.CancelarEntrada(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrNaoEncontrado) {

			return ErrId
		}

		log.Printf("erro crítico ao buscar entrada ID %d: %v", id, err)
		return ErrFalhaNoBancoDeDados

	}

	return nil

}
