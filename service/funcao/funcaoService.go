package funcao

import (
	"context"
	"errors"

	"fmt"
	"strings"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type FuncaoInterface interface {
	AddFuncao(ctx context.Context, funcao *model.FuncaoInserir) error
	DeletarFuncao(ctx context.Context, id int) error
	BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error)
	UpdateFuncao(ctx context.Context, id int, funcao string) (int64, error)
	BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error)
	PossuiFuncionariosVinculados(ctx context.Context, id int) (bool, error)
}

type FuncaoService struct {
	FuncaoRepo FuncaoInterface
}

func NewFuncaoService(repo FuncaoInterface) *FuncaoService {

	return &FuncaoService{
		FuncaoRepo: repo,
	}
}

var (

	ErrFuncaoCadastrada = errors.New("funcao ja cadastrada no sistema")
	ErrCtx = errors.New("context encerrado ou invalido")
	ErrId = errors.New("id invalido")
	ErrRegistroNaoEncontrado = errors.New("função não encontrada")
	ErrFuncaoMinCaracteres =  errors.New("funcao deve ter ao minimo 2 caracteres")
	ErrVinculo = errors.New("erro ao checar vinculos com funcionario")
	ErrFuncaoComVinculo = errors.New("nao foi possivel apagar funcao, funcao com vinculo a funcionario")
)
// SalvarFuncao implements [Funcao].
func (f *FuncaoService) SalvarFuncao(ctx context.Context, model *model.FuncaoInserir) error {
	

	model.Funcao = strings.TrimSpace(model.Funcao)

	if len(model.Funcao) < 2 {

			return ErrFuncaoMinCaracteres
	}
	if err:= f.FuncaoRepo.AddFuncao(ctx, model); err != nil {
		  
		if errors.Is(err, Errors.ErrSalvar){
			return  ErrFuncaoCadastrada
		}

		if errors.Is(err, Errors.ErrDadoIncompativel){

			return  fmt.Errorf("id departamento invalido, não existe no sistema %w", ErrId)
		}

		return  fmt.Errorf("erro ao salvar funcao: %w", err)
	}

	return  nil
}

// ListarFuncao implements [Funcao].
func (f *FuncaoService) ListarFuncao(ctx context.Context, id int) (model.FuncaoDto, error) {
	
	if id <= 0 {

		return model.FuncaoDto{}, ErrId
	}
	
	funcao, err:= f.FuncaoRepo.BuscarFuncao(ctx, id)
	if err != nil {

		if errors.Is(err, Errors.ErrNaoEncontrado){

			return  model.FuncaoDto{}, ErrRegistroNaoEncontrado
		}

		return model.FuncaoDto{}, ErrRegistroNaoEncontrado
	}

	return model.FuncaoDto{
		ID: funcao.ID,
		Funcao: funcao.Funcao,
		Departamento: model.DepartamentoDto{
			ID: funcao.IdDepartamento,
			Departamento: funcao.NomeDepartamento,
		},
	}, nil

}

// ListasTodasFuncao implements [Funcao].
func (f *FuncaoService) ListasTodasFuncao(ctx context.Context) ([]model.FuncaoDto, error) {


	funcs, err:= f.FuncaoRepo.BuscarTodasFuncao(ctx)
	if err != nil {

		if errors.Is(err, Errors.ErrBuscarTodos){

			return []model.FuncaoDto{}, nil
		}
		return  []model.FuncaoDto{}, err
	}



	dto:= make([]model.FuncaoDto, 0 ,len(funcs))
	for _, funcao:= range funcs {

		Func:= model.FuncaoDto{
			ID: funcao.ID,
			Funcao: funcao.Funcao,
			Departamento: model.DepartamentoDto{
				ID: funcao.IdDepartamento,
				Departamento: funcao.NomeDepartamento,
			},
		}

		dto = append(dto, Func)

	}

	return  dto, nil
}

// AtualizarFuncao implements [Funcao].
func (f *FuncaoService) AtualizarFuncao(ctx context.Context, id int, funcao string) error {

		
	if id <= 0 {
		return  ErrId
	}

	funcaoLimpa:= strings.TrimSpace(funcao)
	
	if len(funcaoLimpa) < 2 {

		return  ErrFuncaoMinCaracteres
	}
	linha, err:= f.FuncaoRepo.UpdateFuncao(ctx, id, funcaoLimpa)
	if err != nil {

		if errors.Is(err, Errors.ErrSalvar){
			return  ErrFuncaoCadastrada
		}

		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return ErrRegistroNaoEncontrado
	}

	return  nil
}

// DeletarFuncao implements [Funcao].
func (f *FuncaoService) DeletarFuncao(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  ErrId
	}

	vinculo, err:= f.FuncaoRepo.PossuiFuncionariosVinculados(ctx, id)
	if err != nil {

		return ErrVinculo
	}

	if vinculo {

		return  ErrFuncaoComVinculo
	}

	err = f.FuncaoRepo.DeletarFuncao(ctx, id)
	if err != nil {

		return  fmt.Errorf("erro ao deletar a funcao, %w", err)
	}
	
	return  nil
}


