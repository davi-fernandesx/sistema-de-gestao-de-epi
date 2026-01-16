package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
)


type FuncaoRepository interface {

	Adicionar(ctx context.Context, args repository.AddFuncaoParams) error
	ListarFuncao(ctx context.Context, id int32) (repository.BuscarFuncaoRow, error)
	ListarFuncoes(ctx context.Context)([]repository.BuscarTodasFuncoesRow, error) 
	CancelarFuncao(ctx context.Context, id int32) (int64, error)
	AtualizarFuncao(ctx context.Context, arg repository.UpdateFuncaoParams) (int64, error)

}

type FuncaoService struct {

	repo FuncaoRepository
}

func NewFuncaoService(f FuncaoRepository) *FuncaoService {
	return &FuncaoService{repo: f}
}


func (f *FuncaoService) SalvarFuncao(ctx context.Context, model model.Funcao) error {

	model.Funcao = strings.TrimSpace(model.Funcao)

	F:= repository.AddFuncaoParams{
		Nome: model.Funcao,
		Iddepartamento: int32(model.IdDepartamento),
	}
	if err:= f.repo.Adicionar(ctx, F); err != nil {

		return  fmt.Errorf("erro ao salvar funcao, %w", err)
	}

	return  nil
}

func (f *FuncaoService) ListarFuncao(ctx context.Context, id int) (model.FuncaoDto, error) {
	
	if id <= 0 {

		return model.FuncaoDto{}, helper.ErrId
	}
	
	funcao, err:= f.repo.ListarFuncao(ctx, int32(id))
	if err != nil {


		return model.FuncaoDto{},fmt.Errorf("erro ao listar funcao, %w", err)
	}



	return model.FuncaoDto{
		ID: int(funcao.ID),
		Funcao: funcao.Nome,
		Departamento: model.DepartamentoDto{
			ID: int(funcao.Iddepartamento),
			Departamento: funcao.DepartamentoNome,
		},
	}, nil

}

func (f *FuncaoService) ListasTodasFuncao(ctx context.Context) ([]model.FuncaoDto, error) {


	funcs, err:= f.repo.ListarFuncoes(ctx)
	if err != nil {

		return []model.FuncaoDto{}, fmt.Errorf("erro ao listar todas funcoes, %w", err)
	}



	dto:= make([]model.FuncaoDto, 0 ,len(funcs))
	for _, funcao:= range funcs {

		Func:= model.FuncaoDto{
			ID: int(funcao.ID),
			Funcao: funcao.Nome,
			Departamento: model.DepartamentoDto{
				ID: int(funcao.Iddepartamento),
				Departamento: funcao.DepartamentoNome,
			},
		}

		dto = append(dto, Func)

	}

	return  dto, nil
}

func (f *FuncaoService) DeletarFuncao(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  helper.ErrId
	}

	linha,err := f.repo.CancelarFuncao(ctx,int32(id))
	if err != nil {

		return  fmt.Errorf("erro ao deletar a funcao, %w", err)
	}


	if linha == 0 {

		return helper.ErrNaoEncontrado
	}
	
	return  nil
}

func (f *FuncaoService) AtualizarFuncao(ctx context.Context, id int, funcao string) error {

		
	if id <= 0 {
		return  helper.ErrId
	}

	funcaoLimpa:= strings.TrimSpace(funcao)
	
	if len(funcaoLimpa) < 2 {

		return  helper.ErrNomeCurto
	}

	arg:= repository.UpdateFuncaoParams{
		ID: int32(id),
		Nome: funcaoLimpa,
	}

	linha, err:= f.repo.AtualizarFuncao(ctx, arg)
	if err != nil {


		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return helper.ErrNaoEncontrado
	}

	return  nil
}
