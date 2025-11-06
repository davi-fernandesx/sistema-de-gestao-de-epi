package funcionario

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/funcionario"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service"
	//"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service"
)

type Funcionario interface {
	SalvarFuncionario(ctx context.Context, funcionario model.FuncionarioINserir) error
	ListarFuncionarioPorMatricula(ctx context.Context, matricula string) (*model.Funcionario_Dto, error)
	ListaTodosFuncionarios(ctx context.Context) ([]*model.Funcionario_Dto, error)
	DeletarFuncionario(ctx context.Context, matricula string) error
}

type FuncionarioService struct {
	FuncionarioRepo funcionario.FuncionarioInterface
}

func NewFuncionarioService(repo funcionario.FuncionarioInterface) Funcionario {

	return &FuncionarioService{
		FuncionarioRepo: repo,
	}
}

// SalvarFuncinario implements Funcionario.
func (f *FuncionarioService) SalvarFuncionario(ctx context.Context, funcionario model.FuncionarioINserir) error {

	if strings.TrimSpace(funcionario.Nome) == "" {

		return errors.New("nome não pode ser em branco")
	}

	if funcionario.ID_departamento == nil {

		return  errors.New("id do departamento é obrigatorio")
	}

	if funcionario.ID_funcao == nil {
		return  errors.New("id da funcao é obrigatorio")
	}

	if *funcionario.ID_departamento <= 0 || *funcionario.ID_funcao <= 0 {

			return errors.New("ids tem que ter numero positivo")
	}

	err := service.VerificaContext(ctx)
	if err != nil {
		return err
	}

	matriculaInt, err := service.VerificaMatricula(ctx, funcionario.Matricula)
    if err != nil {
        return err // Retorna erro se a matrícula for inválida (ex: "abc")
    }	
	_, err = f.FuncionarioRepo.BuscaFuncionario(ctx, matriculaInt)

	if err != nil {

		if err != sql.ErrNoRows {

			return err
		}

	} else {

		return errors.New("matricula ja cadastrada")
	}

	nomeSemEspacio := service.VerificaEspaço(funcionario.Nome)

	funcSalvar := model.FuncionarioINserir{
		Nome:            nomeSemEspacio,
		Matricula:       funcionario.Matricula,
		ID_departamento: funcionario.ID_departamento,
		ID_funcao:       funcionario.ID_funcao,
	}

	return f.FuncionarioRepo.AddFuncionario(ctx, &funcSalvar)

}

func (f *FuncionarioService) ListarFuncionarioPorMatricula(ctx context.Context, matricula string) (*model.Funcionario_Dto, error) {

	matriculaInt, err := service.VerificaMatricula(ctx, matricula)
	if err != nil {

		return nil, err
	}

	funcionario, err := f.FuncionarioRepo.BuscaFuncionario(ctx, matriculaInt)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			return nil, errors.New("funcionario nao encontrado")
		}

		return nil, err

	}

	funcDto := model.Funcionario_Dto{
		ID:        funcionario.Id,
		Nome:      funcionario.Nome,
		Matricula: funcionario.Matricula,
		Departamento: model.DepartamentoDto{
			ID:           funcionario.ID_departamento,
			Departamento: funcionario.Departamento,
		},
		Funcao: model.FuncaoDto{
			ID:     funcionario.ID_funcao,
			Funcao: funcionario.Funcao,
		},
	}

	return &funcDto, nil

}

func (f *FuncionarioService) ListaTodosFuncionarios(ctx context.Context) ([]*model.Funcionario_Dto, error) {

	err := service.VerificaContext(ctx)
	if err != nil {
		return nil, err
	}

	funcionarios, err := f.FuncionarioRepo.BuscarTodosFuncionarios(ctx)
	if err != nil {

		if errors.Is(err, Errors.ErrBuscarTodos) {

			return []*model.Funcionario_Dto{}, nil
		}

		if errors.Is(err, Errors.ErrFalhaAoEscanearDados) {
			return []*model.Funcionario_Dto{}, fmt.Errorf("erro interno ao processar dados dos funcionarios, %w", err)
		}

		if errors.Is(err, Errors.ErrAoIterar) {

			return []*model.Funcionario_Dto{}, fmt.Errorf("erro inesperado ao processar os dados dos funcionarios: %w", err)
		}

		return []*model.Funcionario_Dto{}, fmt.Errorf("erro inesperado ao buscar funcionarios, %w", err)
	}

	funcionariosDto := make([]*model.Funcionario_Dto, 0, len(funcionarios))

	for _, funcionario := range funcionarios {

		funcDto := model.Funcionario_Dto{
			ID:        funcionario.Id,
			Nome:      funcionario.Nome,
			Matricula: funcionario.Matricula,
			Departamento: model.DepartamentoDto{
				ID:           funcionario.ID_departamento,
				Departamento: funcionario.Departamento,
			},
			Funcao: model.FuncaoDto{
				ID:     funcionario.ID_funcao,
				Funcao: funcionario.Funcao,
			},
		}

		funcionariosDto = append(funcionariosDto, &funcDto)

	}

	return funcionariosDto, nil

}

func (f *FuncionarioService) DeletarFuncionario(ctx context.Context, matricula string) error {

	matriculaInt, err := service.VerificaMatricula(ctx, matricula)
	if err != nil {
		return err
	}

	err = f.FuncionarioRepo.DeletarFuncionario(ctx, matriculaInt)
	if err != nil {

		if errors.Is(err, Errors.ErrInternal) {
			return fmt.Errorf("erro interno ao processar dados, %w", err)
		}

		if errors.Is(err, Errors.ErrLinhasAfetadas) {

			return fmt.Errorf("erro ao verificar linha afetada")
		}

		return fmt.Errorf("erro inesperado ao deletar funcionario, %w", err)
	}

	return nil

}
