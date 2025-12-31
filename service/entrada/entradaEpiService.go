package entrada

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	entradaepi "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/entradaEpi"
)

type entrada interface {
	SalvarEntrada(ctx context.Context, model *model.EntradaEpiInserir) error
	ListarEntrada(ctx context.Context, id int) (model.EntradaEpiDto, error)
	ListasTodasEntradas(ctx context.Context) ([]model.EntradaEpiDto, error)
	DeletarEntradas(ctx context.Context, id int) error
}

type EntradaService struct {
	EntradaRepo entradaepi.EntradaEpi
}

func NewEntradaService(repo entradaepi.EntradaEpi) entrada {

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

	if !model.DataValidade.Time().After(model.DataFabricacao.Time()) {

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

			return model.EntradaEpiDto{},ErrNaoEncontrado
		}

		if errors.Is(err, Errors.ErrFalhaAoEscanearDados){

			return model.EntradaEpiDto{}, ErrFalhaNoBancoDeDados
		}

		if errors.Is(err, Errors.ErrAoIterar){

			return  model.EntradaEpiDto{}, ErrFalhaNoBancoDeDados
		}

		return model.EntradaEpiDto{},err
	}

	return model.EntradaEpiDto{
		ID: entrada.ID,
		Epi: model.EpiDto{
			Id: entrada.ID_epi,
			Nome: entrada.Nome,
			Fabricante: entrada.Fabricante,
			CA: entrada.CA,
			Tamanho:[]model.TamanhoDto{
				{
					ID: entrada.Id_Tamanho,
					Tamanho: entrada.TamanhoDescricao,
				},
			},
			Descricao: entrada.Descricao,
			DataValidadeCa: entrada.DataValidadeCa,
			Protecao: model.TipoProtecaoDto{
				ID: entrada.IDprotecao,
				Nome: model.Protecao(entrada.NomeProtecao),
			},			
		},
		Data_entrada: entrada.Data_entrada,
		Quantidade: entrada.Quantidade,
		Lote: entrada.Lote,
		Fornecedor: entrada.Fornecedor,
		ValorUnitario: entrada.ValorUnitario,
	}, nil
}

// ListasTodasEntradas implements [entrada].
func (e *EntradaService) ListasTodasEntradas(ctx context.Context) ([]model.EntradaEpiDto, error) {
	panic("unimplemented")
}

// DeletarEntradas implements [entrada].
func (e *EntradaService) DeletarEntradas(ctx context.Context, id int) error {
	panic("unimplemented")
}
