package tamanhos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
	mssql "github.com/microsoft/go-mssqldb"
)

type TamanhsoInterface interface {
	AddTamanhos(ctx context.Context, tamanhos *model.Tamanhos) error
	DeletarTamanhos(ctx context.Context, id int) error
	BuscarTamanhos(ctx context.Context, id int) (*model.Tamanhos, error)
	BuscarTodosTamanhos(ctx context.Context) ([]model.Tamanhos, error)
}

type SqlServerLogin struct {
	DB *sql.DB
}

func NewTamanhoRepository(db *sql.DB) TamanhsoInterface {

	return &SqlServerLogin{
		DB: db,
	}
}

// AddTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) AddTamanhos(ctx context.Context, tamanhos *model.Tamanhos) error {
	query:= `insert into tamanho values (@tamanho)`

	_, err:= s.DB.ExecContext(ctx, query,sql.Named("tamanho", tamanhos.Tamanho))

	if err != nil {
		var ErrSql *mssql.Error
		if errors.As(err, &ErrSql) && ErrSql.Number == 2627{
			return  repository.ErrTamanhoJaExistente
		}
		return repository.ErrAoAdicionarTamanho;
	}

	return  nil
}

// BuscarTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) BuscarTamanhos(ctx context.Context, id int) (*model.Tamanhos, error) {

	query:= "select id, tamanho from tamanho where id = @id"

	var tamanho model.Tamanhos

	err:= s.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&tamanho.ID, &tamanho.Tamanho)

	if err != nil {
		if err == sql.ErrNoRows {
			return  nil, repository.ErrAoProcurarTamanho
		}

		return  nil, repository.ErrFalhaAoEscanearDados
	}

	return  &tamanho, nil
}

// BuscarTodosTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) BuscarTodosTamanhos(ctx context.Context) ([]model.Tamanhos, error) {
	
	query:= "select id, tamanho from tamanho"

	linhas, err:= s.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Tamanhos{}, repository.ErrAoBuscarTodosOsTamanhos
	}

	defer linhas.Close()

	var tamanhos []model.Tamanhos

	for linhas.Next(){
		
		var tamanho model.Tamanhos
		err:= linhas.Scan(&tamanho.ID, &tamanho.Tamanho)
		if err != nil {
			return  nil, repository.ErrFalhaAoEscanearDados
		}

		tamanhos = append(tamanhos, tamanho)
	}

	err = linhas.Err()
	if err != nil {

		return  nil, repository.ErrAoIterarSobreTamanhos
	}

	return  tamanhos, nil
}

// DeletarTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) DeletarTamanhos(ctx context.Context, id int) error {
	
	query:= `delete from tamanho where id = @id`

	result, err:= s.DB.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil {
		return  repository.ErrLinhasAfetadas
	}

	if linhas == 0{
		return repository.ErrTamanhoNaoEncontrado
	}

	return  nil
}
