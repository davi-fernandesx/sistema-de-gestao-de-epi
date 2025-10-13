package entradaepi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EntradaEpi interface {
	AddEntradaEpi(ctx context.Context, EntradaEpi *model.EntradaEpiInserir) error
	DeletarEntrada(ctx context.Context, id int) error
	BuscarEntrada(ctx context.Context, id int) (*model.EntradaEpi, error)
	BuscarTodasEntradas(ctx context.Context) ([]model.EntradaEpi, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewEntradaRepository(db *sql.DB) EntradaEpi {

	return &NewSqlLogin{
		DB: db,
	}
}


// AddEntradaEpi implements EntradaEpi.
func (n *NewSqlLogin) AddEntradaEpi(ctx context.Context, EntradaEpi *model.EntradaEpiInserir) error {
	
	query:= `
		insert into Entrada (id_epi, data_entrada, quantidade, lote, fornecedor)
		values (@id_epi, @data_entrada, @quantidade, @lote, @fornecedor)

	`

	_, err:= n.DB.ExecContext(ctx, query,
		sql.Named("id_epi", EntradaEpi.ID_epi),
		sql.Named("data_entrada", EntradaEpi.Data_entrada),
		sql.Named("quantidade", EntradaEpi.Quantidade),
		sql.Named("lote", EntradaEpi.Lote),
		sql.Named("fornecedor", EntradaEpi.Fornecedor),
	)

	if err != nil {
		return  fmt.Errorf("erro interno ao salvar entrada, %w", Errors.ErrSalvar)
	}

	return  nil
}

// BuscarEntrada implements EntradaEpi.
func (n *NewSqlLogin) BuscarEntrada(ctx context.Context, id int) (*model.EntradaEpi, error) {
	
	query:= `
			select
				ee.id, ee.id_epi,e.nome,e.fabricante, e.CA,e.descricao,
				e.dataFabricacao, e.dataValidade, e.dataValidadeCa, 
				e.id_protecao, e.nomeProtecao, ee.lote, ee.fornecedor 
			from 
				entrada ee
			inner join
				epi e on ee.id_epi = e.id 
			where 
				ee.id = @id
	`

	var entrada model.EntradaEpi

	err:= n.DB.QueryRowContext(ctx, query,sql.Named("id", id)).Scan(
		&entrada.ID,
		&entrada.ID_epi,
		&entrada.Nome,
		&entrada.Fabricante,
		&entrada.CA,
		&entrada.Descricao,
		&entrada.DataFabricacao,
		&entrada.DataValidade,
		&entrada.DataValidadeCa,
		&entrada.IDprotecao,
		&entrada.NomeProtecao,
		&entrada.Lote,
		&entrada.Fornecedor,

	)	

	if err != nil {

		if err == sql.ErrNoRows{

		return  nil, fmt.Errorf("entrada com id %d, não encontrado! %w",id,  Errors.ErrNaoEncontrado)
		}

		return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &entrada, nil
}

// BuscarTodasEntradas implements EntradaEpi.
func (n *NewSqlLogin) BuscarTodasEntradas(ctx context.Context) ([]model.EntradaEpi, error) {
	query:= `
			select
				ee.id, ee.id_epi,e.nome,e.fabricante, e.CA,e.descricao,
				e.dataFabricacao, e.dataValidade, e.dataValidadeCa, 
				e.id_protecao, e.nomeProtecao, ee.lote, ee.fornecedor 
			from 
				entrada ee
			inner join
				epi e on ee.id_epi = e.id  
	`

	linhas, err:= n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.EntradaEpi{}, fmt.Errorf("erro ao procurar todas as entradas, %w", Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	 entradas:= make([]model.EntradaEpi, 0)

	for linhas.Next(){
		var entrada model.EntradaEpi

		err:= linhas.Scan(
			&entrada.ID,
			&entrada.ID_epi,
			&entrada.Nome,
			&entrada.Fabricante,
			&entrada.CA,
			&entrada.Descricao,
			&entrada.DataFabricacao,
			&entrada.DataValidade,
			&entrada.DataValidadeCa,
			&entrada.IDprotecao,
			&entrada.NomeProtecao,
			&entrada.Lote,
			&entrada.Fornecedor,
		)

		if err != nil {
			return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		entradas = append(entradas, entrada)
	}


	if err := linhas.Err(); err != nil {

		return nil, fmt.Errorf("erro ao iterar sobre as entradas , %w", Errors.ErrAoIterar)
	}

	return entradas, nil
}

// DeletarEntrada implements EntradaEpi.
func (n *NewSqlLogin) DeletarEntrada(ctx context.Context, id int) error {
	
		query:= `delete from  entrada where id = @id`

	result, err:= n.DB.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  fmt.Errorf("erro interno ao deletar uma entrada %w", Errors.ErrInternal)
	}

	linhas, err:= result.RowsAffected()
	if err != nil{
		if errors.Is(err, Errors.ErrLinhasAfetadas){

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	}

	if linhas == 0 {

		return  fmt.Errorf("entrada com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return  nil
}

