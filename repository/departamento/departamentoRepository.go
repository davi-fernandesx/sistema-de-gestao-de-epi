package departamento

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/denisenkom/go-mssqldb"
)

type DepartamentoInterface interface {
	AddDepartamento(ctx context.Context, departamento *model.Departamento) error
	DeletarDepartamento(ctx context.Context, id int) error
	BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error)
	BuscarTodosDepartamentos(ctx context.Context) (*[]model.Departamento, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewDepartamentoRepository(db *sql.DB) DepartamentoInterface {

	return &NewSqlLogin{
		DB: db,
	}
}


// AddDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) AddDepartamento(ctx context.Context, departamento *model.Departamento) error {
	
	query:= `insert into departamento (departamento) values (@departamento)`

	_, err:= n.DB.ExecContext(ctx, query, departamento.Departamento)
	if err != nil {
		var ErrSql *mssql.Error
		if errors.As(err, &ErrSql) && ErrSql.Number == 2627 {
			return  repository.ErrDepartamentoJaExistente
			
		}
		return  err
	}

	return  nil
}

// BuscarDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) BuscarDepartamento(ctx context.Context, id int) (*model.Departamento, error) {
	
	query:= `select departamento from departamento where id = @id`

	var departamento model.Departamento

	err:= n.DB.QueryRow(query, sql.Named("id", id)).Scan(
		&departamento.ID,
		&departamento.Departamento,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return  nil, repository.ErrDepartamentoNaoEncontrado
		}

		return  nil, repository.ErrFalhaAoEscanearDados
	}

	return  &departamento, nil
}

// BuscarTodosDepartamentos implements DepartamentoInterface.
func (n *NewSqlLogin) BuscarTodosDepartamentos(ctx context.Context) (*[]model.Departamento, error) {
	
	query:= `select id, departamento from departamento`

	linhas, err:= n.DB.Query(query)
	if err != nil {
		return  nil, repository.ErrBuscarTodosDepartamentos
	}

	defer linhas.Close()

	var departamentos []model.Departamento

	for linhas.Next(){
		var departamento model.Departamento


		if err:= linhas.Scan(&departamento.ID, &departamento.Departamento); err != nil {
			return  nil, repository.ErrFalhaAoEscanearDados
		}

		departamentos = append(departamentos, departamento)
	}

	if err:= linhas.Err(); err != nil {

		return  nil,  repository.ErrIterarSobreDepartamentos
	}

	return &departamentos, nil

}

// DeletarDepartamento implements DepartamentoInterface.
func (n *NewSqlLogin) DeletarDepartamento(ctx context.Context, id int) error {
		
	query:= `delete from departamento where id = @id`

	result, err:= n.DB.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil {
		return  repository.ErrLinhasAfetadas
	}

	if linhas == 0 {
		return  repository.ErrDepartamentoNaoEncontrado
	}

	return  nil
}
