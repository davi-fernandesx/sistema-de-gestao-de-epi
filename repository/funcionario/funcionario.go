package funcionario

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type FuncionarioInterface interface {
	AddFuncionario(ctx context.Context, funcionario *model.FuncionarioINserir) error
	BuscaFuncionario(ctx context.Context, id int) (*model.Funcionario, error)
	BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error)
	DeletarFuncionario(ctx context.Context, id int) error
}

type ConnDB struct {
	DB *sql.DB
}

func NewFuncionarioRepository(db *sql.DB) FuncionarioInterface {

	return &ConnDB{
		DB: db,
	}
}

// AddFuncionario implements FuncionarioInterface.
func (c *ConnDB) AddFuncionario(ctx context.Context, funcionario *model.FuncionarioINserir) error {

	query := `insert into funcionario(nome, id_departamento, id_funcao) values( @nome, @id_departamento, @id_funcao)`

	_, err := c.DB.ExecContext(ctx, query,
		sql.Named("nome", funcionario.Nome),
		sql.Named("id_departamento", funcionario.ID_departamento),
		sql.Named("id_funcao", funcionario.ID_funcao),
	)
	if err != nil {

		var Errsql *mssql.Error
		if errors.As(err, &Errsql) && Errsql.Number == 2627 {

			return fmt.Errorf("funcionario: %s, ja existe no sistema, %w", funcionario.Nome, Errors.ErrSalvar)

		}
		return fmt.Errorf("erro interno ao salvar funcionario, %w", Errors.ErrInternal)
	}

	return nil
}

// BuscaFuncionario implements FuncionarioInterface.
func (c *ConnDB) BuscaFuncionario(ctx context.Context, id int) (*model.Funcionario, error) {

	query := `select fn.id, fn.nome, fn.id_departamento, d.departamento, fn.id_funcao, f.funcao
			from funcionario fn
			inner join departamento d on fn.id_departamento = d.id
			inner jon funcao f on fn.id_funcao = f.funcao
			where fn.id = @id`

	var funcionario model.Funcionario
	err := c.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&funcionario.Id,
		&funcionario.Nome,
		&funcionario.ID_departamento,
		&funcionario.Departamento,
		&funcionario.ID_funcao,
		&funcionario.Funcao,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("funcionario com o id %d não encontrado: %w", funcionario.Id, Errors.ErrNaoEncontrado)
		}
		return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &funcionario, nil
}

// BuscarTodosFuncionarios implements FuncionarioInterface.
func (c *ConnDB) BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error) {

	query := `select fn.id, fn.nome, fn.id_departamento, d.departamento, fn.id_funcao, f.funcao
			from funcionario fn
			inner join departamento d on fn.id_departamento = d.id
			inner jon funcao f on fn.id_funcao = f.funcao`

	var funcionarios []model.Funcionario

	results, err := c.DB.QueryContext(ctx, query)
	if err != nil {

		return nil, fmt.Errorf("erro ao buscar todos os funcionarios, %w", Errors.ErrBuscarTodos)
	}

	defer results.Close()

	for results.Next() {

		var funcionario model.Funcionario

		err := results.Scan(&funcionario.Id, &funcionario.Nome, &funcionario.ID_departamento, &funcionario.Departamento, &funcionario.ID_funcao, &funcionario.Funcao)
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
func (c *ConnDB) DeletarFuncionario(ctx context.Context, id int) error {

	query:= `delete from funcionario where id = @id`

	  result, err:= c.DB.ExecContext(ctx, query, sql.Named("id", id))
	  if err != nil {
		return  fmt.Errorf("%w", Errors.ErrInternal)
	  }

	linhas, err:= result.RowsAffected()

	if err != nil {

		if errors.Is(err, Errors.ErrLinhasAfetadas){
			return  fmt.Errorf("erro ao verificar linhas afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	}

	if linhas == 0 {

		return fmt.Errorf("funcionario com o id %d nao  encontrado, %w ", id, Errors.ErrNaoEncontrado)
	}

	return  nil
}
