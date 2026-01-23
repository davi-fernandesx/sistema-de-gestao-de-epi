package service

import (
	"context"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
)

type TamanhoRepository interface {
	Adicionar(ctx context.Context, tamanho string) error
	ListarTamanho(ctx context.Context, id int32) (repository.BuscarTamanhoRow, error)
	ListarTamanhos(ctx context.Context) ([]repository.BuscarTodosTamanhosRow, error)
	CancelarTamanho(ctx context.Context, id int) (int64, error)
}

type TamanhoService struct {

	repo TamanhoRepository
}

func NewTamanhoService(t TamanhoRepository) *TamanhoService {

	return &TamanhoService{repo: t}
}

func (t *TamanhoService) SalvarTamanho(ctx context.Context, model model.Tamanhos) error {

	model.Tamanho = strings.TrimSpace(model.Tamanho)

	if err:=t.repo.Adicionar(ctx, model.Tamanho); err != nil {

		return err
	}

	return  nil
}

func (t *TamanhoService) ListarTamanho(ctx context.Context, id int)(model.TamanhoDto, error) {

	if id <= 0 {
		return  model.TamanhoDto{},helper.ErrId
	}

	tamanho, err:= t.repo.ListarTamanho(ctx, int32(id))
	if err != nil {

		return model.TamanhoDto{}, err
	}

	return model.TamanhoDto{
		
		ID: int(tamanho.ID),
		Tamanho: tamanho.Tamanho,
	}, nil
}

func (t *TamanhoService) ListarTodosTamanhos(ctx context.Context) ([]model.TamanhoDto, error){

	tamanhos, err:= t.repo.ListarTamanhos(ctx)
	if err != nil {

		return  []model.TamanhoDto{},err
	}

	tamDto:= make([]model.TamanhoDto, 0, len(tamanhos))

	for _, tamanho:= range tamanhos {

		tam:= model.TamanhoDto{
			ID: int(tamanho.ID),
			Tamanho: tamanho.Tamanho,
		}

		tamDto = append(tamDto, tam)
	}

	if tamDto == nil {

		return []model.TamanhoDto{}, nil
	}

	return  tamDto, nil
}


func (t *TamanhoService) CancelarTamanho(ctx context.Context, id int) error {

	linhas, err := t.repo.CancelarTamanho(ctx, id)
	if err != nil {

		return err
	}

	if linhas == 0 {

		return helper.ErrNaoEncontrado
	}

	return  nil
}
