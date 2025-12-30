package epi

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/epi"
)

type Epi interface {
	SalvarEpi(ctx context.Context, model *model.EpiInserir) error
	ListarEpi(ctx context.Context, id int) (model.EpiDto, error)
	ListasTodosEpis(ctx context.Context) ([]model.EpiDto, error)
	DeletarEpi(ctx context.Context, id int) error
	AtualizarEpiNome(ctx context.Context, id int, nomeNovo string) error
	AtualizarEpiCa(ctx context.Context, id int, CAnovo string) error
	AtualizarEpiFabricante(ctx context.Context, id int, FabricanteNovo string) error
	AtualizaDescricao(ctx context.Context, id int, DescricaoNova string) error

}

type EpiService struct {
	EpiRepo epi.EpiInterface
}



func NewEpiServices(repo epi.EpiInterface) Epi {

	return &EpiService{

		EpiRepo: repo,
	}
}

var (

	ErrCaCadastrado = errors.New("CA ja cadastrada no sistema")
	ErrId = errors.New("id invalido")
	ErrNulo = errors.New("esse campo tem que conter pelo menos 1 tamanho")
	ErrEpiNaoEncontrado = errors.New("epi nao encontrado no sistema")
	ErrFalhaNoBancoDeDados = errors.New("falha no banco de dados")
	ErrCa = errors.New("Ca invalido")
	errDataMenor = errors.New("A data de validade não pode ser menor que hoje")
	ErrDataZero = errors.New("data não pode ser vazia")
)

var reCA = regexp.MustCompile(`^[0-9]{1,6}$`)
// SalvarEpi implements [Epi].
func (e *EpiService) SalvarEpi(ctx context.Context, model *model.EpiInserir) error {

	model.Nome = strings.TrimSpace(model.Nome)
	model.Fabricante = strings.TrimSpace(model.Fabricante)
	model.CA = strings.TrimSpace(model.CA)
	model.Descricao = strings.TrimSpace(model.Descricao)

	if len(model.Idtamanho) == 0 {

		return ErrNulo
	
	}

	err:= e.EpiRepo.AddEpi(ctx, model)
	if err != nil{

		if errors.Is(err, Errors.ErrSalvar){

			return  ErrCaCadastrado
		}

		if errors.Is(err, Errors.ErrDadoIncompativel){

			return fmt.Errorf("verifique os ids de protecao ou tamanhos, %w", ErrId)
		}
		
		return  fmt.Errorf("erro ao cadastrar Epi, %w", err)
	}

	return nil

}


func (e *EpiService) ListasTodosEpis(ctx context.Context) ([]model.EpiDto, error) {
	
	epis, err:= e.EpiRepo.BuscarTodosEpi(ctx)
	if err != nil {

		return nil, ErrFalhaNoBancoDeDados
	}

	if len(epis) == 0 {

		return nil, nil
	}

	dto:= make([]model.EpiDto, 0, len(epis))


	for _, epi:= range epis {

		tamanhosDtos:= make([]model.TamanhoDto, 0, len(epi.Tamanhos))
		
		for _, t := range epi.Tamanhos {

			tamanhosDtos = append(tamanhosDtos, model.TamanhoDto{
				ID: t.ID,
				Tamanho: t.Tamanho,
			})
		}

		e := model.EpiDto{
				Id: epi.ID,
				Nome: epi.Nome,
				Fabricante: epi.Fabricante,
				CA: epi.CA,
				Tamanho: tamanhosDtos,
				Descricao: epi.Descricao,
				DataValidadeCa: epi.DataValidadeCa.Time(),
				Protecao: model.TipoProtecaoDto{
					ID: epi.IDprotecao,
					Nome: model.Protecao(epi.NomeProtecao),
				},
			}
		
			dto = append(dto, e)
	}

	return  dto, nil
}


func (e *EpiService) ListarEpi(ctx context.Context, id int) (model.EpiDto, error) {

	if id <= 0 {

		return  model.EpiDto{}, ErrId
	}

	epi, err:= e.EpiRepo.BuscarEpi(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrNaoEncontrado){

			return model.EpiDto{},ErrEpiNaoEncontrado
		}

		if errors.Is(err, Errors.ErrFalhaAoEscanearDados){
			
			return model.EpiDto{}, ErrFalhaNoBancoDeDados
		}

		return  model.EpiDto{}, ErrFalhaNoBancoDeDados
	}

	dto:= model.EpiDto{
		Id: epi.ID,
		Nome: epi.Nome,
		Fabricante: epi.Fabricante,
		CA: epi.CA,
		Tamanho: make([]model.TamanhoDto, 0, len(epi.Tamanhos)),
		Descricao: epi.Descricao,
		DataValidadeCa: epi.DataValidadeCa.Time(),
		Protecao: model.TipoProtecaoDto{
			ID: epi.IDprotecao,
			Nome: model.Protecao(epi.NomeProtecao),
		},
	}

	
	for _, tamanho:= range epi.Tamanhos{

		tamanhos:= model.TamanhoDto{

			ID: tamanho.ID,
			Tamanho: tamanho.Tamanho,
		}

		dto.Tamanho = append(dto.Tamanho, tamanhos)
	}

	return  dto, nil
}

// ListasTodosEpis implements [Epi].

// AtualizarEpiCa implements [Epi].
func (e *EpiService) AtualizarEpiCa(ctx context.Context, id int, CAnovo string) error {
	
	if id <= 0 {

		return  ErrId
	}
	
	CAnovo = strings.TrimSpace(CAnovo)

	if !reCA.MatchString(CAnovo) {

		return  ErrCa

	}

	err := e.EpiRepo.UpdateEpiCa(ctx, id, CAnovo)
		if err != nil {

			if errors.Is(err, Errors.ErrSalvar){

				return ErrCaCadastrado
			}
				return  ErrFalhaNoBancoDeDados
		}


	return  nil

}

// AtualizarEpiFabricante implements [Epi].
func (e *EpiService) AtualizarEpiFabricante(ctx context.Context, id int, FabricanteNovo string) error {
	
	if id <= 0 {

		return  ErrId
	}
	FabricanteNovo = strings.TrimSpace(FabricanteNovo)

	err := e.EpiRepo.UpdateEpiFabricante(ctx, id, FabricanteNovo)
	if err != nil {

		if errors.Is(err, Errors.ErrInternal){

			return ErrFalhaNoBancoDeDados
		}

		return ErrFalhaNoBancoDeDados
	}

	return  nil
}

// AtualizarEpiNome implements [Epi].
func (e *EpiService) AtualizarEpiNome(ctx context.Context, id int, nomeNovo string) error {
	 
	if id <= 0 {

		return  ErrId
	}
	nomeNovo = strings.TrimSpace(nomeNovo)

	err :=  e.EpiRepo.UpdateEpiNome(ctx, id, nomeNovo)
	if err != nil{

		if errors.Is(err, Errors.ErrInternal){

			return  ErrFalhaNoBancoDeDados
		}

		return  ErrFalhaNoBancoDeDados
	}


	return  nil
}

func (e *EpiService) AtualizaDescricao(ctx context.Context, id int, DescricaoNova string) error {

	if id <= 0 {

		return  ErrId
	}
	DescricaoNova = strings.TrimSpace(DescricaoNova)

	err := e.EpiRepo.UpdateEpiDescricao(ctx, id, DescricaoNova)
	if err != nil {

		if errors.Is(err, Errors.ErrInternal){

			return ErrFalhaNoBancoDeDados
		}

		return  ErrFalhaNoBancoDeDados
	}

	return  nil
}

func (e *EpiService) AtualizaDataValidadeCa(ctx context.Context, id int, dataNova configs.DataBr) error{

	if dataNova.IsZero() {

		return  ErrDataZero
	}
	novaData:= dataNova.Time()
	if id <=0 {

		return ErrId
	}

	hoje:= time.Now().Truncate(24 * time.Hour)

	if novaData.Before(hoje){

		return  errDataMenor

	}

	return  e.EpiRepo.UpdateEpiDataValidadeCa(ctx, id, novaData)
}

// DeletarEpi implements [Epi].
func (e *EpiService) DeletarEpi(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  ErrId
	}

	err:= e.EpiRepo.DeletarEpi(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrAoapagar){

			return fmt.Errorf("verifique o ID passado, %w",ErrId)
		}
	}

	return  nil
}

// ListarEpi implements [Epi].

