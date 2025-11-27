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
	CancelarEntrada(ctx context.Context, id int) error
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

	query := `
		insert into Entrada (id_epi,id_tamanho, data_entrada, quantidade, lote, fornecedor, valorUnitario)
		values (@id_epi,@id_tamanho, @data_entrada, @quantidade, @lote, @fornecedor, @valorUnitario)
	`

	_, err := n.DB.ExecContext(ctx, query,
		sql.Named("id_epi", EntradaEpi.ID_epi),
		sql.Named("id_tamanho", EntradaEpi.Id_tamanho),
		sql.Named("data_entrada", EntradaEpi.Data_entrada),
		sql.Named("quantidade", EntradaEpi.Quantidade),
		sql.Named("lote", EntradaEpi.Lote),
		sql.Named("fornecedor", EntradaEpi.Fornecedor),
		sql.Named("valorUnitario", EntradaEpi.ValorUnitario),
	)

	if err != nil {
		return fmt.Errorf("erro interno ao salvar entrada, %w", Errors.ErrSalvar)
	}

	return nil
}

// BuscarEntrada implements EntradaEpi.
func (n *NewSqlLogin) BuscarEntrada(ctx context.Context, id int) (*model.EntradaEpi, error) {

	query := `
     SELECT
            ee.id, ee.id_epi, ee.quantidade, ee.lote, ee.fornecedor, -- Campos da tabela de entrada
            e.nome, e.fabricante, e.CA, e.descricao,ee.valorUnitario,
            e.data_fabricacao, e.data_validade, e.validade_CA, -- Campos do EPI
            tp.id as id_protecao, tp.protecao as nome_protecao, -- Campos do Tipo de Proteção
            t.id as id_tamanho, t.tamanho as tamanho_descricao -- Campos do Tamanho
        FROM 
            entradas_epi ee
        INNER JOIN
            epi e ON ee.id_epi = e.id 
        INNER JOIN
            tipo_protecao tp ON e.id_tipo_protecao = tp.id
        INNER JOIN
            tamanhos t ON ee.id_tamanho = t.id
        WHERE 
            ee.cancelada_em IS NULL AND  ee.id = @id;
	`

	var entrada model.EntradaEpi

	err := n.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&entrada.ID,
		&entrada.ID_epi,
		&entrada.Quantidade,
		&entrada.Lote,
		&entrada.Fornecedor,
		&entrada.Nome,
		&entrada.Fabricante,
		&entrada.CA,
		&entrada.Descricao,
		&entrada.ValorUnitario,
		&entrada.DataFabricacao,
		&entrada.DataValidade,
		&entrada.DataValidadeCa,
		&entrada.IDprotecao,
		&entrada.NomeProtecao,
		&entrada.Id_Tamanho,
		&entrada.TamanhoDescricao,
		
	)

	if err != nil {

		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("entrada com id %d, não encontrado! %w", id, Errors.ErrNaoEncontrado)
		}

		return nil, fmt.Errorf("%w", err)
	}

	return &entrada, nil
}

// BuscarTodasEntradas implements EntradaEpi.
func (n *NewSqlLogin) BuscarTodasEntradas(ctx context.Context) ([]model.EntradaEpi, error) {
	query := `
     SELECT
            ee.id, ee.id_epi, ee.quantidade, ee.lote, ee.fornecedor, -- Campos da tabela de entrada
            e.nome, e.fabricante, e.CA, e.descricao,ee.valorUnitario,
            e.data_fabricacao, e.data_validade, e.validade_CA, -- Campos do EPI
            tp.id as id_protecao, tp.protecao as nome_protecao, -- Campos do Tipo de Proteção
            t.id as id_tamanho, t.tamanho as tamanho_descricao -- Campos do Tamanho
        FROM 
            entradas_epi ee
        INNER JOIN
            epi e ON ee.id_epi = e.id 
        INNER JOIN
            tipo_protecao tp ON e.id_tipo_protecao = tp.id
        INNER JOIN
            tamanhos t ON ee.id_tamanho = t.id
		where ee.cancelada_em IS NULL
		`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.EntradaEpi{}, fmt.Errorf("erro ao procurar todas as entradas, %w", Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	entradas := make([]model.EntradaEpi, 0)

	for linhas.Next() {
		var entrada model.EntradaEpi

		err := linhas.Scan(
			&entrada.ID,
			&entrada.ID_epi,
			&entrada.Quantidade,
			&entrada.Lote,
			&entrada.Fornecedor,
			&entrada.Nome,
			&entrada.Fabricante,
			&entrada.CA,
			&entrada.Descricao,
			&entrada.ValorUnitario,
			&entrada.DataFabricacao,
			&entrada.DataValidade,
			&entrada.DataValidadeCa,
			&entrada.IDprotecao,
			&entrada.NomeProtecao,
			&entrada.Id_Tamanho,
			&entrada.TamanhoDescricao,
		)

		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		entradas = append(entradas, entrada)
	}

	if err := linhas.Err(); err != nil {

		return nil, fmt.Errorf("erro ao iterar sobre as entradas , %w", Errors.ErrAoIterar)
	}

	return entradas, nil
}

// DeletarEntrada implements EntradaEpi.
func (n *NewSqlLogin) CancelarEntrada(ctx context.Context, id int) error {

	query := `update entrada
			set cancelada_em = GETDATE()
			where id = @id AND cancelada_em IS NULL`

	result, err := n.DB.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("erro interno ao cancelar uma entrada %w", Errors.ErrInternal)
	}

	linhas, err := result.RowsAffected()
	if err != nil {
		if errors.Is(err, Errors.ErrLinhasAfetadas) {

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	}

	if linhas == 0 {

		return fmt.Errorf("entrada com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil
}
