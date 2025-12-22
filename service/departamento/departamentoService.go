package departamento

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/departamento"
)

type Departamento interface {
	SalvarDepartamento(ctx context.Context, model *model.Departamento) error
	ListarDepartamento(ctx context.Context, id int) (model.DepartamentoDto, error)
	ListarTodosDepartamentos(ctx context.Context) ([]model.DepartamentoDto, error)
	DeletarDepartamento(ctx context.Context, id int) error
	AtualizarDepartamento(ctx context.Context, id int, departamento string) error
}

type DepartamentoServices struct {
	DepartamentoRepo departamento.DepartamentoInterface
}


func NewDepartamentoService(repo departamento.DepartamentoInterface) Departamento {

	return &DepartamentoServices{
		DepartamentoRepo: repo,
	}
}

var (

	ErrNomeVazio = errors.New("nome do departamento não pode ser vazio")
	ErrCtx = errors.New("context encerrado ou invalido")
	ErrId = errors.New("id invalido")
	ErrDep = errors.New("departamento não encontrado")
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

	err:= d.DepartamentoRepo.DeletarDepartamento(ctx, id)
	if err != nil {

		if strings.Contains(err.Error(), "547 "){

			return  fmt.Errorf("não é possivel excluir, departamento possui dependencia com funcionario e funcao")
		}

		return  fmt.Errorf("erro ao deletar um departamento, %w", err)
	}

	return  nil
}

func (d *DepartamentoServices) AtualizarDepartamento(ctx context.Context, id int, departamento string) error {

	 if id <= 0 {

		return  ErrId
	 }

	 linha,errDep:= d.DepartamentoRepo.UpdateDepartamento(ctx, id, departamento)
	 if errDep != nil {

		if strings.Contains(errDep.Error(), "2627") || strings.Contains(errDep.Error(), "2601"){

			return fmt.Errorf("erro, constraint UNIQUE sendo violado, departamento ja existente com esse nome %s", departamento)
		}

		return fmt.Errorf("erro tecnico ao realizar o update: %w", errDep)
	 }

	 if linha == 0 {

		return fmt.Errorf("departamento com id %d, ,não existe", id)
	 }

	 return  nil

}


