package funcao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type FuncaoInterface interface {
	AddFuncao(ctx context.Context, funcao *model.Funcao) error
	DeletarFuncao(ctx context.Context, id int) error 
	BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error)
	UpdateFuncao(ctx context.Context, id int, funcao string)error
	BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error)
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
func (s *SqlServerLogin) AddFuncao(ctx context.Context, funcao *model.Funcao) error {
		query:= `insert into funcao (funcao, idDepartamento) values (@funcao, @idDepartamento)`

		_, err:= s.Db.ExecContext(ctx, query, sql.Named("funcao",funcao.Funcao), sql.Named("idDepartamento", funcao.IdDepartamento))
		if err != nil{
			var ErrSql *mssql.Error
			if errors.As(err, &ErrSql) && ErrSql.Number == 2627 {
				return  fmt.Errorf("função %s ja existe no sistema!, %w", funcao.Funcao, Errors.ErrSalvar)
			}

			return  fmt.Errorf("erro interno ao salvar funcao!, %w", Errors.ErrInternal)
		}

		return  nil
}

// BuscarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) BuscarFuncao(ctx context.Context, id int) (*model.Funcao, error) {
	
	query:= `select id, funcao from funcao where id = @id`

	var funcao model.Funcao

	err:= s.Db.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&funcao.ID, &funcao.Funcao)
	if err != nil {
		if err == sql.ErrNoRows {
			return  nil, fmt.Errorf("funcao com o id %d não encotrada!, %w", id, Errors.ErrNaoEncontrado)
		}

		return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &funcao, nil
}

// BuscarTodosCargos implements FuncaoInterface.
func (s *SqlServerLogin) BuscarTodasFuncao(ctx context.Context) ([]model.Funcao, error) {
	query:= `select id, funcao, IdDepartamento from funcao`

	linhas, err:= s.Db.QueryContext(ctx, query)
	if err != nil {
		return []model.Funcao{}, fmt.Errorf("erro ao procurar todas as funções, %w", Errors.ErrBuscarTodos)
	}

	var Funcoes []model.Funcao
	defer linhas.Close()

	for linhas.Next(){

		var funcao model.Funcao

		if err:= linhas.Scan(&funcao.ID, &funcao.Funcao, &funcao.IdDepartamento); err!= nil {

			return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		Funcoes = append(Funcoes, funcao)
	}
	
	err = linhas.Err()
	if err != nil {
		return  nil, fmt.Errorf("erro ao iterar sobre as funções , %w", Errors.ErrAoIterar)
	}

	return  Funcoes, nil

}

// DeletarFuncao implements FuncaoInterface.
func (s *SqlServerLogin) DeletarFuncao(ctx context.Context, id int) error {
	
	query:= `delete from funcao where id = @id`

	result, err:= s.Db.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  err
	}

	
	linhas, err:= result.RowsAffected()
	if err != nil {

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0 {

		return fmt.Errorf("função com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil

}

func (s *SqlServerLogin)UpdateFuncao(ctx context.Context, id int, funcao string)error{

	query:= `update funcao
			set funcao = @funcao
			where id = @id`

	_, err:= s.Db.ExecContext(ctx, query, sql.Named("funcao", funcao), sql.Named("id", id))

	if err != nil {
		 return  fmt.Errorf("erro ao atualizar funcao, %w", Errors.ErrInternal)
	}

	return  nil
}