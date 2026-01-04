package departamento

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type DepartamentoRepo interface {
	AddDepartamento(ctx context.Context, departamento *model.Departamento) error
	DeletarDepartamento(ctx context.Context, id int) error
	BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error)
	BuscarTodosDepartamentos(ctx context.Context) ([]model.Departamento, error)
	UpdateDepartamento(ctx context.Context, id int, departamento string)(int64,error)
	PossuiFuncoesVinculadas(ctx context.Context, id int) (bool, error)
}

type DepartamentoServices struct {
	DepartamentoRepo DepartamentoRepo
}


func NewDepartamentoService(repo DepartamentoRepo) *DepartamentoServices {

	return &DepartamentoServices{
		DepartamentoRepo: repo,
	}
}

var (

	ErrNomeVazio = errors.New("nome do departamento não pode ser vazio")
	ErrCtx = errors.New("context encerrado ou invalido")
	ErrId = errors.New("id invalido")
	ErrDep = errors.New("departamento não encontrado")
	ErrVinculo = errors.New("erro ao checar vinculos com funcoes")
	ErrFuncaoComVinculo = errors.New("nao foi possivel apagar departamento, departamento com vinculo a função(õe)s")
)

// SalvarDepartamento implements [Departamento].
func (d *DepartamentoServices) SalvarDepartamento(ctx context.Context, model *model.Departamento) error {
	
	model.Departamento = strings.TrimSpace(model.Departamento)

	if len(model.Departamento) < 2 {

		return  fmt.Errorf("departamento deve ter ao minimo 2 caracteres")

	}

	if err:= d.DepartamentoRepo.AddDepartamento(ctx, model); err != nil {

		return fmt.Errorf("erro ao salvar departamento: %w", err)
	}

	return  nil
	
}
// ListarDepartamento implements [Departamento].
func (d *DepartamentoServices) ListarDepartamento(ctx context.Context, id int) (model.DepartamentoDto, error) {
	
	if id <= 0 {

		return model.DepartamentoDto{},ErrId
	}

	dep, err:= d.DepartamentoRepo.BuscarDepartamento(ctx, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows){

			return model.DepartamentoDto{}, ErrDep
		}
		return model.DepartamentoDto{}, fmt.Errorf("erro ao buscar departamento, %w", err)
	}

	if dep == nil {

		return model.DepartamentoDto{}, ErrDep
	}
	
	
	return model.DepartamentoDto{

		ID: dep.ID,
		Departamento: dep.Departamento,
	}, nil

}

// ListarTodosDepartamentos implements [Departamento].
func (d *DepartamentoServices) ListarTodosDepartamentos(ctx context.Context) ([]model.DepartamentoDto, error) {

	deps, err:= d.DepartamentoRepo.BuscarTodosDepartamentos(ctx)
	if err != nil {

			return  nil, err
	}

	if deps == nil {

		return []model.DepartamentoDto{}, nil
	}

	dto:=make([]model.DepartamentoDto, 0, len(deps))

	for _, dep := range deps {

		departamento:= model.DepartamentoDto {
			ID: dep.ID,
			Departamento: dep.Departamento,
		}


		dto = append(dto, departamento)
	}

	return  dto, nil
}

// DeletarDepartamento implements [Departamento].
func (d *DepartamentoServices) DeletarDepartamento(ctx context.Context, id int) error {
	
	if id <= 0 {

		return  ErrId
	}

	vinculo, err:= d.DepartamentoRepo.PossuiFuncoesVinculadas(ctx, id)
	if err != nil{

		return  ErrVinculo
	}

	if vinculo {

		return  ErrFuncaoComVinculo
	}

	err = d.DepartamentoRepo.DeletarDepartamento(ctx, id)
	if err != nil {


		return  fmt.Errorf("erro ao deletar um departamento, %w, departamento ja pode estar inativo", err)
	}

	return  nil
}

func (d *DepartamentoServices) AtualizarDepartamento(ctx context.Context, id int, departamento string) error {

	 if id <= 0 {

		return  ErrId
	 }

	 linha,errDep:= d.DepartamentoRepo.UpdateDepartamento(ctx, id, departamento)
	 if errDep != nil {

		return fmt.Errorf("erro tecnico ao realizar o update: %w", errDep)
	 }

	 if linha == 0 {

		return fmt.Errorf("departamento com id %d, ,não encontrado", id)
	 }

	 return  nil

}


