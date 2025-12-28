package funcao

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
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

type SqlServerLogin struct {
	Db *sql.DB
}

func NewfuncaoRepository(db *sql.DB) FuncaoInterface {

	return &SqlServerLogin{
		Db: db,
	}
}

// AddFuncao implements FuncaoInterface.
func (s *SqlServerLogin) AddFuncao(ctx context.Context, funcao *model.FuncaoInserir) error {
	query := `insert into funcao (nome, IdDepartamento) values (@funcao, @idDepartamento)`

	_, err := s.Db.ExecContext(ctx, query, sql.Named("funcao", funcao.Funcao), sql.Named("idDepartamento", funcao.IdDepartamento))
	if err != nil {
			if helper.IsForeignKeyViolation(err){

				return fmt.Errorf("departamento não existente no banco de dados, %w", Errors.ErrDadoIncompativel)
			}

			if helper.IsUniqueViolation(err){

				return fmt.Errorf("funcao %s ja existe no sistema, %w", funcao.Funcao, Errors.ErrSalvar)
			}

			return fmt.Errorf("erro interno ao salvar: %w", Errors.ErrInternal)
		}
	return nil
}

// BuscarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error) {

	query := `select f.id, f.nome, f.IdDepartamento, d.nome as departamento
			from funcao f
			inner join 
				departamento d on f.IdDepartamento = d.id
			where f.id = @id and f.ativo = 1`

	var funcao model.Funcao

	err := s.Db.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&funcao.ID, &funcao.Funcao, &funcao.IdDepartamento, &funcao.NomeDepartamento)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("funcao com o id %d não encontrada!, %w", id, Errors.ErrNaoEncontrado)
		}

		return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &funcao, nil
}

// BuscarTodosCargos implements FuncaoInterface.
func (s *SqlServerLogin) BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error) {
	query := `select f.id, f.nome, f.IdDepartamento, d.nome as departamento
			from funcao f
			inner join 
				departamento d on f.IdDepartamento = d.id
			where f.ativo = 1`

	linhas, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		return []model.Funcao{}, fmt.Errorf("erro ao procurar todas as funções, %w", Errors.ErrBuscarTodos)
	}

	var Funcoes []model.Funcao
	defer linhas.Close()

	for linhas.Next() {

		var funcao model.Funcao

		if err := linhas.Scan(&funcao.ID, &funcao.Funcao, &funcao.IdDepartamento, &funcao.NomeDepartamento); err != nil {

			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		Funcoes = append(Funcoes, funcao)
	}

	err = linhas.Err()
	if err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre as funções , %w", Errors.ErrAoIterar)
	}

	return Funcoes, nil

}

func (s *SqlServerLogin) PossuiFuncionariosVinculados(ctx context.Context, id int) (bool, error) {
	var total int
	query := `SELECT COUNT(1) FROM funcionario WHERE IdFuncao = @id AND ativo = 1`

	err := s.Db.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&total)
	return total > 0, err
}

// DeletarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) DeletarFuncao(ctx context.Context, id int) error {

	query := `update funcao
				set ativo = 0,
				deletado_em = getdate()
				where id = @id and ativo = 1`

	result, err := s.Db.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return err
	}

	linhas, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)

	}

	if linhas == 0 {

		return fmt.Errorf("função com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil

}

func (s *SqlServerLogin) UpdateFuncao(ctx context.Context, id int, funcao string) (int64, error) {

	query := `update funcao
			set nome = @funcao
			where id = @id and ativo = 1`

	linhas, err := s.Db.ExecContext(ctx, query, sql.Named("funcao", funcao), sql.Named("id", id))

	if err != nil {

		if helper.IsUniqueViolation(err){
			return 0,fmt.Errorf("funcao %s ja existe no sistema, %w", funcao, Errors.ErrSalvar)

		}
			return 0, fmt.Errorf("erro ao atualizar funcao, %w", Errors.ErrInternal)
	}

	
	return linhas.RowsAffected()
}
