package epi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EpiInterface interface {
	AddEpi(ctx context.Context, epi *model.EpiInserir) error
	DeletarEpi(ctx context.Context, id int) error
	BuscarEpi(ctx context.Context, id int) (*model.Epi, error)
	BuscarTodosEpi(ctx context.Context) ([]model.Epi, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewEpiRepository(db *sql.DB) EpiInterface {

	return &NewSqlLogin{
		DB: db,
	}
}

// AddEpi implements EpiInterface.
func (n *NewSqlLogin) AddEpi(ctx context.Context, epi *model.EpiInserir) error {

	query := `insert into epi (nome, fabricante, CA, descricao, data_fabricacao, data_validade, validade_CA, id_tipo_protecao, alerta_minimo) values (
			@nome, @fabricante, @CA, @descricao,@data_fabricacao, @data_validade, @validade_CA, @id_tipo_protecao, @alerta_minimo )`

	_, err := n.DB.ExecContext(ctx, query,
		sql.Named("nome", epi.Nome),
		sql.Named("fabricantte", epi.Fabricante),
		sql.Named("CA", epi.CA),
		sql.Named("descricao", epi.Descricao),
		sql.Named("data_fabricacao", epi.DataFabricacao),
		sql.Named("data_validade", epi.DataValidade),
		sql.Named("validade_CA", epi.DataValidadeCa),
		sql.Named("id_tipo_protecao", epi.IDprotecao),
		sql.Named("alerta_minimo", epi.AlertaMinimo))

	if err != nil {
		return fmt.Errorf(" Erro interno ao salvar Epi, %w", Errors.ErrInternal)
	}

	return nil
}

// BuscarEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {

	query := `
			select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id		
			where
				e.id = @id
	`

	var epi model.Epi
	err := n.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&epi.ID,
		&epi.Nome,
		&epi.Fabricante,
		&epi.CA,
		&epi.Descricao,
		&epi.DataFabricacao,
		&epi.DataValidade,
		&epi.DataValidadeCa,
		&epi.AlertaMinimo,
		&epi.IDprotecao,
		&epi.NomeProtecao,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("epi com id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
		}

		return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &epi, nil
}

// BuscarTodosEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarTodosEpi(ctx context.Context) ([]model.Epi, error) {

	query := `
				select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Epi{}, fmt.Errorf("erro ao buscar todos os epi!, %w", Errors.ErrBuscarTodos)
	}
	defer linhas.Close()

	var epis []model.Epi

	for linhas.Next() {

		var epi model.Epi

		err := linhas.Scan(
			&epi.ID,
			&epi.Nome,
			&epi.Fabricante,
			&epi.CA,
			&epi.Descricao,
			&epi.DataFabricacao,
			&epi.DataValidade,
			&epi.DataValidadeCa,
			&epi.AlertaMinimo,
			&epi.IDprotecao,
			&epi.NomeProtecao,)
		 if err != nil {

			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		epis = append(epis, epi)
	}

	if err := linhas.Err(); err != nil {

		return nil, fmt.Errorf("erro ao iterar sobre os epis!, %w", Errors.ErrAoIterar)
	}

	return epis, nil

}

// DeletarEpi implements EpiInterface.
func (n *NewSqlLogin) DeletarEpi(ctx context.Context, id int) error {
	
	query:=  `delete from epi where id = @id`

	result, err:= n.DB.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil{
			if errors.Is(err, Errors.ErrLinhasAfetadas){

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	}

	if linhas == 0 {

		return  fmt.Errorf("epi com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return  nil
}
