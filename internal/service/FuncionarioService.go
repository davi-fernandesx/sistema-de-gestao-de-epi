package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
)

type FuncionarioRepository interface {
	Adicionar(ctx context.Context, args repository.AddFuncionarioParams) error
	ListarFuncionario(ctx context.Context, matricula string) (repository.BuscaFuncionarioRow, error)
	ListarFuncionarios(ctx context.Context) ([]repository.BuscarTodosFuncionariosRow, error)
	CancelarFuncionario(ctx context.Context, id int32) (int64, error)
	AtualizarFuncionarioNome(ctx context.Context, arg repository.UpdateFuncionarioNomeParams) (int64, error)
	AtualizarFuncionarioDepartamento(ctx context.Context, arg repository.UpdateFuncionarioDepartamentoParams) (int64, error)
	AtualizarFuncionarioFuncao(ctx context.Context, arg repository.UpdateFuncionarioFuncaoParams) (int64, error)
}

type FuncionarioService struct {

	repo FuncionarioRepository
}

func NewFuncionarioService(f FuncionarioRepository) *FuncionarioService {
	return &FuncionarioService{repo: f}
}


func (f *FuncionarioService) SalvarFuncionario(ctx context.Context, model model.FuncionarioINserir) error {

	model.Nome = strings.TrimSpace(model.Nome)

	args:= repository.AddFuncionarioParams{
		Nome: model.Nome,
		Matricula: model.Matricula,
		Iddepartamento: int32(model.ID_departamento),
		Idfuncao: int32(model.ID_funcao),
	}
	err:= f.repo.Adicionar(ctx, args)
	if err != nil {

		return err
	}

	return nil
}

func (f *FuncionarioService) ListarFuncionario(ctx context.Context, matricula string) (model.Funcionario_Dto, error) {

	
	if matricula <= "" {

		return model.Funcionario_Dto{}, helper.ErrId
	}
	funcionario, err:= f.repo.ListarFuncionario(ctx, matricula)
	if err != nil {

		return model.Funcionario_Dto{},err
	}

	funcDto := model.Funcionario_Dto{
		ID:        int(funcionario.ID),
		Nome:      funcionario.Nome,
		Matricula: funcionario.Matricula,
		Funcao: model.FuncaoDto{
			ID:     int(funcionario.Idfuncao),
			Funcao: funcionario.FuncaoNome,
			Departamento: model.DepartamentoDto{
				ID:           int(funcionario.Iddepartamento),
				Departamento: funcionario.DepartamentoNome,
			},
		},
	}

	return funcDto, nil

}

func (f *FuncionarioService) ListaTodosFuncionarios(ctx context.Context) ([]model.Funcionario_Dto, error) {

	funcionarios, err := f.repo.ListarFuncionarios(ctx)
	if err != nil {

		return  []model.Funcionario_Dto{}, err
	}

	funcionariosDto := make([]model.Funcionario_Dto, 0, len(funcionarios))

	for _, funcionario := range funcionarios {

		funcDto := model.Funcionario_Dto{
			ID:           int(funcionario.ID),
			Nome:         funcionario.Nome,
			Matricula:    funcionario.Matricula,
			Funcao: model.FuncaoDto{
				ID:     int(funcionario.Idfuncao),
				Funcao: funcionario.FuncaoNome,
				Departamento: model.DepartamentoDto{
					ID:           int(funcionario.Iddepartamento),
					Departamento: funcionario.DepartamentoNome,
				},
			},
		}

		funcionariosDto = append(funcionariosDto, funcDto)

	}

	if funcionariosDto == nil {

		return []model.Funcionario_Dto{}, nil
	}

	return funcionariosDto, nil

}

func (f *FuncionarioService) DeletarFuncionario(ctx context.Context, id int) error {


	linhas, err := f.repo.CancelarFuncionario(ctx, int32(id))
	if err != nil {

		return fmt.Errorf("erro ao deletar funcionario, %w, funcionario ja pode estar inativo", err)
	}

	if linhas == 0 {

		return helper.ErrNaoEncontrado
	}

	return nil

}

func (f *FuncionarioService) AtualizaNomeFuncionario(ctx context.Context, id int, nome string) error {

	if id <= 0 {
		return  helper.ErrId
	}

	nomeLimpo := strings.TrimSpace(nome)

	if len(nomeLimpo) < 2 {

		return  helper.ErrNomeCurto
	}
	args:= repository.UpdateFuncionarioNomeParams{
		ID: int32(id),
		Nome: nomeLimpo,
	}

	linha, err:= f.repo.AtualizarFuncionarioNome(ctx, args)
	if err != nil {

		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return helper.ErrNaoEncontrado
	}

	return  nil
}

func (f *FuncionarioService) AtualizaDepartamentoFuncionario(ctx context.Context, id , iddepartamento int) error {
	
	if id <= 0 {
		return  helper.ErrId
	}

	args:= repository.UpdateFuncionarioDepartamentoParams{
		ID: int32(id),
		Iddepartamento: int32(iddepartamento),
	}

	linha, err:= f.repo.AtualizarFuncionarioDepartamento(ctx, args)
	if err != nil {

		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return helper.ErrNaoEncontrado
	}

	return  nil
}

func (f *FuncionarioService) AtualizaFuncaoFuncionario(ctx context.Context, id , idfuncao int) error {
	
	if id <= 0 {
		return  helper.ErrId
	}

	args:= repository.UpdateFuncionarioFuncaoParams{
		ID: int32(id),
		Idfuncao: int32(idfuncao),
	}
	linha, err:= f.repo.AtualizarFuncionarioFuncao(ctx, args)
	if err != nil {

		return fmt.Errorf("erro tecnico ao realizar o update: %w", err) 
	}

	if linha == 0 {
		return helper.ErrNaoEncontrado
	}

	return  nil
}