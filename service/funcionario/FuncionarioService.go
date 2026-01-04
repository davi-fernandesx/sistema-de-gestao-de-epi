package funcionario

import (
	"context"
	"database/sql"
	"errors"
	"fmt"


	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/service"
)

type FuncionarioInterface interface {
	AddFuncionario(ctx context.Context, funcionario *model.FuncionarioINserir) error
	BuscaFuncionario(ctx context.Context, matricula int) (*model.Funcionario, error)
	BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error)
	DeletarFuncionario(ctx context.Context, matricula int) error
	UpdateFuncionarioNome(ctx context.Context, id int, funcionario string)error
	UpdateFuncionarioFuncao(ctx context.Context, id int, idFuncao string)error
	UpdateFuncionarioDepartamento(ctx context.Context, id int, idDepartamento string)error
}


type FuncionarioService struct {
	FuncionarioRepo FuncionarioInterface
}

func NewFuncionarioService(repo FuncionarioInterface) *FuncionarioService {

	return &FuncionarioService{
		FuncionarioRepo: repo,
	}
}

var (

		ErrMatricula = errors.New("matricula ja cadastrada")
		ErrFuncionarios = errors.New("funcionario não encontrado")

)
// SalvarFuncinario implements Funcionario.
func (f *FuncionarioService) SalvarFuncionario(ctx context.Context, funcionario model.FuncionarioINserir) error {


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

		return ErrMatricula
	}


	funcSalvar := model.FuncionarioINserir{
		Nome:            funcionario.Nome,
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

	if funcionario == nil {

		return  nil, ErrFuncionarios
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

	if funcionariosDto == nil {

		return []*model.Funcionario_Dto{}, nil
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

			return fmt.Errorf("erro ao deletar funcionario, %w, funcionario ja pode estar inativo", err)
	}

	return nil

}
