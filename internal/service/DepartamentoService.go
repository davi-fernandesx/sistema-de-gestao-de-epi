package service

import (
	"context"
	"errors"
	
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=DepartamentoRepository --output=../../mocks --outpkg=mocks --with-expecter
type DepartamentoRepository interface {
	Adicionar(ctx context.Context, departamento string) error
	ListarDepartamento(ctx context.Context, id int32) (repository.BuscarDepartamentoRow, error)
	ListarDepartamentos(ctx context.Context) ([]repository.BuscarTodosDepartamentosRow, error)
	CancelarDepartamento(ctx context.Context, id int32) (int64, error)
	AtualizarDepartamento(ctx context.Context, arg repository.UpdateDepartamentoParams) (int64, error)
}

type DepartamentoService struct {

	repo DepartamentoRepository
}

func NewDepartamentoService(r DepartamentoRepository) *DepartamentoService {
	return &DepartamentoService{repo: r}
}


func (d *DepartamentoService) SalvarDepartamento(ctx context.Context, model model.Departamento) error {

	model.Departamento = strings.TrimSpace(model.Departamento)

	if len(model.Departamento) < 2 {

		return  helper.ErrNomeCurto

	}

	if err := d.repo.Adicionar(ctx, model.Departamento); err != nil {

	  if errors.Is(err, helper.ErrDadoDuplicado){

		return  err
	  }

	  return err
	}

	return  nil
}

func (d *DepartamentoService) ListarDepartamento(ctx context.Context, id int32) (model.DepartamentoDto, error) {
	
	if id <= 0 {

		return model.DepartamentoDto{},helper.ErrId
	}

	dep, err:= d.repo.ListarDepartamento(ctx, id)
	if err != nil {

		if err == pgx.ErrNoRows {

			return model.DepartamentoDto{},helper.ErrNaoEncontrado
		}
		return  model.DepartamentoDto{}, err
	}
	
	return model.DepartamentoDto{

		ID: int(dep.ID),
		Departamento: dep.Nome,
	}, nil

}

func (d *DepartamentoService) ListarTodosDepartamentos(ctx context.Context) ([]model.DepartamentoDto, error) {

	deps, err:= d.repo.ListarDepartamentos(ctx)
	if err != nil {

			return  nil, err
	}

	if deps == nil {

		return []model.DepartamentoDto{}, nil
	}

	dto:=make([]model.DepartamentoDto, 0, len(deps))

	for _, dep := range deps {

		departamento:= model.DepartamentoDto {
			ID: int(dep.ID),
			Departamento: dep.Nome,
		}


		dto = append(dto, departamento)
	}

	return  dto, nil
}

func (d *DepartamentoService) DeletarDepartamento(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  helper.ErrId
	}


	idDep, err := d.repo.CancelarDepartamento(ctx, int32(id))
	if err != nil {

		return  helper.ErrInternal
	}

	if idDep == 0 {

		return  helper.ErrNaoEncontrado
	}

	return  nil
}

func (d *DepartamentoService) AtualizarDepartamento(ctx context.Context, id int32, novoNome string) error {

	novoNome = strings.TrimSpace(novoNome)
	if len(novoNome) < 2 {
        return helper.ErrNomeCurto
    }

	arg := repository.UpdateDepartamentoParams{
        ID:   id,
        Nome: novoNome,
    }

	 linha,errDep:= d.repo.AtualizarDepartamento(ctx, arg)
	 if errDep != nil {

		return  errDep
	 }

	 if linha == 0 {

		return helper.ErrNaoEncontrado
	 }

	 return  nil

}