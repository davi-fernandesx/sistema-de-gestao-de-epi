package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
)



type ProtecaoRepository interface {

	Adicionar(ctx context.Context, nome string) error
	ListarProtecao(ctx context.Context, id int32) (repository.BuscarProtecaoRow, error)
	ListarProtecoes(ctx context.Context) ([]repository.BuscarTodasProtecoesRow, error)
	CancelarProtecao(ctx context.Context, id int32) (int64, error)
}

type ProtecaoService struct {

	repo ProtecaoRepository
}

func NewProtecaoService(p ProtecaoRepository) *ProtecaoService {

	return &ProtecaoService{repo: p}
}

func (p *ProtecaoService) SalvarProtecao(ctx context.Context, model model.TipoProtecao) error {

	model.Nome = strings.TrimSpace(model.Nome)

	err:= p.repo.Adicionar(ctx, model.Nome)
	if err != nil {
		return  err
	}

	return nil
}

func (p *ProtecaoService) ListarProtecao(ctx context.Context, id int) (model.TipoProtecaoDto, error){

	if id <= 0 {
		return  model.TipoProtecaoDto{},helper.ErrId
	}

	protecao, err:= p.repo.ListarProtecao(ctx, int32(id))
	if err != nil {

		return model.TipoProtecaoDto{}, err
	}

	return model.TipoProtecaoDto{
		ID: int64(protecao.ID),
		Nome: protecao.Nome,
	}, nil
}

func (p *ProtecaoService) ListarProtecoes(ctx context.Context) ([]model.TipoProtecaoDto, error) {

	protec, err := p.repo.ListarProtecoes(ctx)
	if err != nil {

		return  []model.TipoProtecaoDto{}, err
	}

	protecDto := make([]model.TipoProtecaoDto, 0, len(protec))

	for _, prot := range protecDto {

		pro := model.TipoProtecaoDto{
			ID: prot.ID,
			Nome: prot.Nome,
		}

		protecDto = append(protecDto, pro)

	}

	if protec == nil {

		return []model.TipoProtecaoDto{}, nil
	}

	return protecDto, nil

}

func (p *ProtecaoService) DeletarProtecao(ctx context.Context, id int) error {


	linhas, err := p.repo.CancelarProtecao(ctx, int32(id))
	if err != nil {

		return fmt.Errorf("erro ao deletar funcionario, %w, funcionario ja pode estar inativo", err)
	}

	if linhas == 0 {

		return helper.ErrNaoEncontrado
	}

	return nil

}