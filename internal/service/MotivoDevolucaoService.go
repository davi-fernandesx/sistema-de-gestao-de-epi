package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
)

type MotivoDevolucaoRepository interface {
	Adicionar(ctx context.Context, motivo string) error
	ListarMotivo(ctx context.Context, id int) (repository.BuscaMotivoDevolucaoRow, error)
	ListarMotivos(ctx context.Context) ([]repository.BuscaTodosMotivosDevolucaoRow, error)
	CancelarMotivoDevolucao(ctx context.Context, id int) (int64, error)
}

type MotivoDevolucaoService struct {
	repo MotivoDevolucaoRepository
}

func NewMotivoDevolucaoRepositoryServe(m MotivoDevolucaoRepository) *MotivoDevolucaoService {

	return &MotivoDevolucaoService{repo: m}
}

func (m *MotivoDevolucaoService) Salvar(ctx context.Context, model model.MotivoDevolucao) error {

	model.Motivo = strings.TrimSpace(model.Motivo)

	err := m.repo.Adicionar(ctx, model.Motivo)
	if err != nil {

		return err
	}

	return nil
}

func (m *MotivoDevolucaoService) ListarMotivo(ctx context.Context, id int) (model.MotivoDevolucaoEpiDto, error) {

	if id <= 0 {

		return model.MotivoDevolucaoEpiDto{}, helper.ErrId
	}

	motivo, err:= m.repo.ListarMotivo(ctx, id)
	if err != nil {

		return model.MotivoDevolucaoEpiDto{}, err
	}

	return model.MotivoDevolucaoEpiDto{

		Id: int(motivo.ID),
		Motivo: motivo.Motivo,
	}, nil
}

func (m *MotivoDevolucaoService) ListarMotivos(ctx context.Context) ([]model.MotivoDevolucaoEpiDto, error){

	motivos, err:= m.repo.ListarMotivos(ctx)
	if err != nil {

		return []model.MotivoDevolucaoEpiDto{},err
	}

	dto:= make([]model.MotivoDevolucaoEpiDto, 0, len(motivos))

	for _, mot := range motivos {

		M := model.MotivoDevolucaoEpiDto{
			Id: int(mot.ID),
			Motivo: mot.Motivo,
		}
		dto = append(dto, M)
	}

	return dto, nil
}

func (m *MotivoDevolucaoService) DeletarMotivo(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  helper.ErrId
	}

	linha,err := m.repo.CancelarMotivoDevolucao(ctx,id)
	if err != nil {

		return  fmt.Errorf("erro ao deletar a funcao, %w", err)
	}


	if linha == 0 {

		return helper.ErrNaoEncontrado
	}
	
	return  nil
}