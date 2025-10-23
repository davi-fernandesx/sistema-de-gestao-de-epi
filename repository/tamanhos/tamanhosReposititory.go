package tamanhos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type TamanhsoInterface interface {
	AddTamanhos(ctx context.Context, tamanhos *model.Tamanhos) error
	DeletarTamanhos(ctx context.Context, id int) error
	BuscarTamanhos(ctx context.Context, id int) (*model.Tamanhos, error)
	BuscarTodosTamanhos(ctx context.Context) ([]model.Tamanhos, error)
	BuscarTamanhosPorIdEpi(ctx context.Context, epiId int)([]model.Tamanhos, error)
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
			return  fmt.Errorf("tamanho %s ja existente!, %w",tamanhos.Tamanho, Errors.ErrSalvar)
		}
		return fmt.Errorf("erro interno ao salvar tamanho. %w", Errors.ErrSalvar)
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
			return  nil, fmt.Errorf("tamanho com id %d, não encontrado! %w",id,  Errors.ErrNaoEncontrado)
		}

			return  nil,  fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		

	return  &tamanho, nil
}

//funcao para trazer do banco de dados todos os tamanhos de um unico epi(por id)
func (s *SqlServerLogin) BuscarTamanhosPorIdEpi(ctx context.Context, epiId int)([]model.Tamanhos, error){

	query := `

		select 
			t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhosEpis te on t.id = te.id_tamanho
		where
			te.epiId = @epiId
	`

	linhas, err:= s.DB.QueryContext(ctx, query, sql.Named("epiId", epiId))
	if err != nil {

		return  nil, fmt.Errorf("erro ao procurar todos os tamanhos com o id %d, %w",epiId,  Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	var tamanhos []model.Tamanhos

	for linhas.Next() {

		var t model.Tamanhos
		err:= linhas.Scan(&t.ID, &t.Tamanho)
		if err != nil {

			return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		tamanhos = append(tamanhos, t)
	}

	err = linhas.Err()
	if err != nil {

		return  nil, fmt.Errorf("erro ao iterar sobre os tamanhos , %w", Errors.ErrAoIterar)
	}

	return tamanhos, nil
}

// BuscarTodosTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) BuscarTodosTamanhos(ctx context.Context) ([]model.Tamanhos, error) {
	
	query:= "select id, tamanho from tamanho"

	linhas, err:= s.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Tamanhos{},  fmt.Errorf("erro ao procurar todos os tamanhos, %w",  Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	var tamanhos []model.Tamanhos

	for linhas.Next(){
		
		var tamanho model.Tamanhos
		err:= linhas.Scan(&tamanho.ID, &tamanho.Tamanho)
		if err != nil {
			return  nil,  fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		tamanhos = append(tamanhos, tamanho)
	}

	err = linhas.Err()
	if err != nil {

		return  nil,  fmt.Errorf("erro ao iterar sobre os tamanhos , %w", Errors.ErrAoIterar)
	}

	return  tamanhos, nil
}

// DeletarTamanhos implements TamanhsoInterface.
func (s *SqlServerLogin) DeletarTamanhos(ctx context.Context, id int) error {
	
	query:= `delete from tamanho where id = @id`

	result, err:= s.DB.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return  fmt.Errorf("erro interno ao deletar tamanho, %w", Errors.ErrInternal)
	}

	linhas, err:= result.RowsAffected()
	if err != nil {
		if errors.Is(err, Errors.ErrLinhasAfetadas){

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	}

	if linhas == 0{
		return  fmt.Errorf("tamanho com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return  nil
}


func (s *SqlServerLogin) UpdateIdEpi(ctx context.Context, id int, idAtualizado int)error{

	query:= `update tamanho_epi
			set id_epi = @idAtualizado
			where id_epi = @id`

	_, err:= s.DB.ExecContext(ctx, query, sql.Named("id_epi", idAtualizado), sql.Named("id_epi", id))
	if err != nil {
		return  fmt.Errorf("erro ao atualizar id do epi, %w", Errors.ErrInternal)
	}

	return  nil
}