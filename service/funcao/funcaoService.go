package funcao

import (
	"context"
	"database/sql"
	"errors"

	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/funcao"
)

type Funcao interface {
	SalvarFuncao(ctx context.Context, model *model.FuncaoInserir) error
	ListarFuncao(ctx context.Context, id int) (model.FuncaoDto, error)
	ListasTodasFuncao(ctx context.Context) ([]model.FuncaoDto, error)
	DeletarFuncao(ctx context.Context, id int) error
	AtualizarFuncao(ctx context.Context, id int, funcao string) error
}

type FuncaoService struct {
	FuncaoRepo funcao.FuncaoInterface
}

func NewFuncaoService(repo funcao.FuncaoInterface) Funcao {

	return &FuncaoService{
		FuncaoRepo: repo,
	}
}

var (

	ErrFuncaoCadastrada = errors.New("funcao ja cadastrada no sistema")
	ErrCtx = errors.New("context encerrado ou invalido")
	ErrId = errors.New("id invalido")
	ErrRegistroNaoEncontrado = errors.New("função não encontrada")
)
// SalvarFuncao implements [Funcao].
func (f *FuncaoService) SalvarFuncao(ctx context.Context, model *model.FuncaoInserir) error {
	

	model.Funcao = strings.TrimSpace(model.Funcao)

	if len(model.Funcao) < 2 {

			return  fmt.Errorf("funcao deve ter ao minimo 2 caracteres")
	}
	if err:= f.FuncaoRepo.AddFuncao(ctx, model); err != nil {
		  
		if strings.Contains(err.Error(), "2627"){
			return  ErrFuncaoCadastrada
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

		if errors.Is(err, sql.ErrNoRows){

			return  model.FuncaoDto{}, err
		}

		return model.FuncaoDto{}, ErrRegistroNaoEncontrado
	}

	if funcao == nil {
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

		return  nil, err
	}

	if funcs == nil {

		return []model.FuncaoDto{}, nil
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

		return  fmt.Errorf("funcao tem que ter no minimo 2 caracteres")
	}
	linha, err:= f.FuncaoRepo.UpdateFuncao(ctx, id, funcaoLimpa)
	if err != nil {

		if strings.Contains(err.Error(), "2627") || strings.Contains(err.Error(), "2601"){

			return fmt.Errorf("erro, constraint UNIQUE sendo violado, funcao ja existente com esse nome: %s", funcao)
		}

		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return fmt.Errorf("funcao com id %d, ,não existe", id)
	}

	return  nil
}

// DeletarFuncao implements [Funcao].
func (f *FuncaoService) DeletarFuncao(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  ErrId
	}

	err:=f.FuncaoRepo.DeletarFuncao(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "547"){
			return  fmt.Errorf("não é possivel excluir, funcao possui dependencia com funcionario")
		}

		return  fmt.Errorf("erro ao deletar a funcao, %w", err)
	}
	
	return  nil
}


