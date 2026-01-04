package departamento

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type DepartamentoInterface interface {
	AddDepartamento(ctx context.Context, departamento *model.Departamento) error
	DeletarDepartamento(ctx context.Context, id int) error
	BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error)
	BuscarTodosDepartamentos(ctx context.Context) ([]model.Departamento, error)
	UpdateDepartamento(ctx context.Context, id int, departamento string)(int64,error)
	PossuiFuncoesVinculadas(ctx context.Context, id int) (bool, error)
}

type DepartamentoRepository struct {
	DB *sql.DB
}

func NewDepartamentoRepository(db *sql.DB) *DepartamentoRepository {

	return &DepartamentoRepository{
		DB: db,
	}
}

// AddDepartamento implements DepartamentoInterface.
func (n *DepartamentoRepository) AddDepartamento(ctx context.Context, departamento *model.Departamento) error {

	query := `insert into departamento (departamento) values (@departamento)`

	_, err := n.DB.ExecContext(ctx, query, sql.Named("departamento", departamento.Departamento))
	if err != nil {

		if helper.IsUniqueViolation(err){
			return fmt.Errorf("departamento %s ja existe no sistema, %w", departamento.Departamento ,Errors.ErrSalvar)
		}
		return fmt.Errorf(" Erro interno ao salvar departamento, %w", Errors.ErrInternal)
	}
	return nil
}

// BuscarDepartamento implements DepartamentoInterface.
func (n *DepartamentoRepository) BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error) {

	query := `select departamento from departamento where id = @id and ativo = 1`

	var departamento model.Departamento

	err := n.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&departamento.ID,
		&departamento.Departamento,
	)

	if err != nil {
		if err == sql.ErrNoRows {
		return  nil, fmt.Errorf("usuario com id %d, não encontrado! %w",id,  Errors.ErrNaoEncontrado)
		}

		return nil,  fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &departamento, nil
}

// BuscarTodosDepartamentos implements DepartamentoInterface.
func (n *DepartamentoRepository) BuscarTodosDepartamentos(ctx context.Context) ([]model.Departamento, error) {

	query := `select id, nome from departamento where ativo = 1`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao procurar todos os departamentos, %w", Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	var departamentos []model.Departamento

	for linhas.Next() {
		var departamento model.Departamento

		if err := linhas.Scan(&departamento.ID, &departamento.Departamento); err != nil {
			return nil,  fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		departamentos = append(departamentos, departamento)
	}

	if err := linhas.Err(); err != nil {

		return nil, fmt.Errorf("erro ao iterar sobre os departamentos, %w", Errors.ErrAoIterar)
	}

	return departamentos, nil

}

func (s *DepartamentoRepository) PossuiFuncoesVinculadas(ctx context.Context, id int) (bool, error) {
    var total int
    query := `SELECT COUNT(1) FROM funcao WHERE IdDepartamento = @id AND ativo = 1`
    
    err := s.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&total)
    return total > 0, err
}
// DeletarDepartamento implements DepartamentoInterface.
func (n *DepartamentoRepository) DeletarDepartamento(ctx context.Context, id int) error {

	query := `update departamento
				set ativo = 0,
					deletado_em = getdate()
				where id = @id and ativo = 1`

	result, err := n.DB.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return err
	}

	linhas, err := result.RowsAffected()
	if err != nil {
			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
	
	}

	if linhas == 0 {
		return fmt.Errorf("departamento com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil
}

func (n *DepartamentoRepository) UpdateDepartamento(ctx context.Context, id int, departamento string)(int64,error){

	query:= `update departamento
			set nome = @departamento
			where id = @id and ativo = 1`

		linha, err:= n.DB.ExecContext(ctx, query, sql.Named("departamento", departamento), sql.Named("id", id))
		if err != nil {
			
			return 0, fmt.Errorf("erro ao atualizar o nome do departamento, %w", Errors.ErrInternal)
		}

	return linha.RowsAffected()
			
}