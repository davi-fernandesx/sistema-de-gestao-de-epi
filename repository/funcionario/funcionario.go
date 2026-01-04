package funcionario

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

// interface sql error (para pegar o codigo do erro)


type FuncionarioRepository struct {
	DB *sql.DB
}

func NewFuncionarioRepository(db *sql.DB) *FuncionarioRepository {

	return &FuncionarioRepository{
		DB: db,
	}
}

// AddFuncionario implements FuncionarioInterface.
func (c *FuncionarioRepository) AddFuncionario(ctx context.Context, funcionario *model.FuncionarioINserir) error {

	query := `insert into funcionario(nome, matricula, IdDepartamento, IdFuncao) values( @nome, @matricula, @id_departamento, @id_funcao)`

	_, err := c.DB.ExecContext(ctx, query,
		sql.Named("nome", funcionario.Nome),
		sql.Named("matricula", funcionario.Matricula),
		sql.Named("id_departamento", funcionario.ID_departamento),
		sql.Named("id_funcao", funcionario.ID_funcao),
	)
	if err != nil {
		if helper.IsForeignKeyViolation(err){
			return fmt.Errorf("id departamento ou id funcao não existente no banco de dados, %w", Errors.ErrDadoIncompativel)
		}

		if helper.IsUniqueViolation(err){
			return fmt.Errorf("matricula %s ja existe no sistema, %w", funcionario.Matricula, Errors.ErrSalvar)

		}
		return fmt.Errorf("erro interno ao salvar funcionario, %w", Errors.ErrInternal)
	}

	return nil
}

// BuscaFuncionario implements FuncionarioInterface.
func (c *FuncionarioRepository) BuscaFuncionario(ctx context.Context, matricula int) (*model.Funcionario, error) {

	query := `select fn.id, fn.nome,fn.matricula, fn.IdDepartamento, d.nome as departamento, 
			fn.IdFuncao, f.nome as funcao
			from funcionario fn
			inner join	
					departamento d on fn.IdDepartamento = d.id
			inner join 	
					funcao f on fn.IdFuncao = f.id
				where fn.matricula = @matricula and fn.ativo = 1`

	var funcionario model.Funcionario
	err := c.DB.QueryRowContext(ctx, query, sql.Named("matricula", matricula)).Scan(
		&funcionario.Id,
		&funcionario.Nome,
		&funcionario.Matricula,
		&funcionario.ID_departamento,
		&funcionario.Departamento,
		&funcionario.ID_funcao,
		&funcionario.Funcao,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("funcionario com a matricula %d não encontrado: %w", funcionario.Matricula, Errors.ErrNaoEncontrado)
		}
		return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &funcionario, nil
}

// BuscarTodosFuncionarios implements FuncionarioInterface.
func (c *FuncionarioRepository) BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error) {

	query := `select fn.id, fn.nome,fn.matricula, fn.IdDepartamento, d.nome as departamento, 
			fn.IdFuncao, f.nome as funcao
			from funcionario fn
			inner join departamento d on fn.IdDepartamento = d.id
			inner join funcao f on fn.IdFuncao = f.id
			where fn.ativo = 1`

	var funcionarios []model.Funcionario

	results, err := c.DB.QueryContext(ctx, query)
	if err != nil {

		return nil, fmt.Errorf("erro ao buscar todos os funcionarios, %w", Errors.ErrBuscarTodos)
	}

	defer results.Close()

	for results.Next() {

		var funcionario model.Funcionario

		err := results.Scan(&funcionario.Id, &funcionario.Nome,&funcionario.Matricula, &funcionario.ID_departamento, &funcionario.Departamento, &funcionario.ID_funcao, &funcionario.Funcao)
		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		funcionarios = append(funcionarios, funcionario)

	}

	if err := results.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os funcionarios, %w", Errors.ErrAoIterar)
	}

	return funcionarios, nil
}

// DeletarFuncionario implements FuncionarioInterface.
func (c *FuncionarioRepository) DeletarFuncionario(ctx context.Context, matricula int) error {

	query:= `update funcionario
			set ativo = 0,
			deletado_em = getdate()
			where id = @id and ativo = 1`

	  result, err:= c.DB.ExecContext(ctx, query, sql.Named("matricula", matricula))
	  if err != nil {
		return  fmt.Errorf("%w", Errors.ErrInternal)
	  }

	linhas, err:= result.RowsAffected()

	if err != nil {
			return  fmt.Errorf("erro ao verificar linhas afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0 {

		return fmt.Errorf("funcionario com a matricula %d nao  encontrado, %w ", matricula, Errors.ErrNaoEncontrado)
	}

	return  nil
}

func (c *FuncionarioRepository)UpdateFuncionarioNome(ctx context.Context, id int, funcionario string)error{

	query:= `update funcionario
		     set nome = @funcionario
			 where id = @id and ativo = 1`

	_, err:= c.DB.ExecContext(ctx, query, sql.Named("funcionario", funcionario), sql.Named("id", id))
	if err != nil {
		return  fmt.Errorf("erro ao atualizar nome do funcinariom, %w", Errors.ErrInternal)
	}

	return  nil
}

func (c *FuncionarioRepository)UpdateFuncionarioDepartamento(ctx context.Context, id int, idDepartamento string)error{

	query:= `update funcionario
		     set IdDepartamento = @idDepartamento
			 where id = @id and ativo = 1`

	_, err:= c.DB.ExecContext(ctx, query, sql.Named("IdDepartamento", idDepartamento), sql.Named("id", id))
	if err != nil {
		return  fmt.Errorf("erro ao atualizar nome do funcinariom, %w", Errors.ErrInternal)
	}

	return  nil
}

func (c *FuncionarioRepository)UpdateFuncionarioFuncao(ctx context.Context, id int, idFuncao string)error{

	query:= `update funcionario
		     set IdFuncao = @idFuncao
			 where id = @id and ativo = 1`

	_, err:= c.DB.ExecContext(ctx, query, sql.Named("IdFuncao", idFuncao), sql.Named("id", id))
	if err != nil {
		return  fmt.Errorf("erro ao atualizar funcao do funcinario, %w", Errors.ErrInternal)
	}

	return  nil
}