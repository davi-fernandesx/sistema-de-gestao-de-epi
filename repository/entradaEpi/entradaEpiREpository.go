package entradaepi

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EntradaRepositorySQL struct {
	DB *sql.DB
}

func NewEntradaRepository(db *sql.DB) *EntradaRepositorySQL {

	return &EntradaRepositorySQL{
		DB: db,
	}
}



const entradaQueryJoin = `  
select ee.id, ee.IdEpi,  e.nome as epi,  e.fabricante, e.CA, e.descricao,ee.data_fabricacao, ee.data_validade, e.validade_CA,
		e.IdTipoProtecao, tp.nome as 'protecao para',
	   	ee.IdTamanho,t.tamanho as tamanho, ee.quantidade,ee.quantidadeAtual ,ee.data_entrada,
	   ee.lote, ee.fornecedor, ee.valor_unitario
from entrada_epi ee
inner join
	epi e on ee.IdEpi = e.id
inner join
	tipo_protecao tp on e.IdTipoProtecao = tp.id
inner join
	tamanho t on ee.IdTamanho = t.id
		`

func (n *EntradaRepositorySQL) buscaEntradas(ctx context.Context, query string, args ...any)([]model.EntradaEpi, error){

	linhas, err := n.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return []model.EntradaEpi{}, fmt.Errorf("erro ao procurar todas as entradas, %w", Errors.ErrInternal)
	}

	defer linhas.Close()

	entradas := make([]model.EntradaEpi, 0)

	for linhas.Next() {
		var entrada model.EntradaEpi

		err := linhas.Scan(
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
			&entrada.Id_Tamanho,
			&entrada.TamanhoDescricao,
			&entrada.Quantidade,
			&entrada.Quantidade_Atual,
			&entrada.Data_entrada,
			&entrada.Lote,
			&entrada.Fornecedor,
			&entrada.ValorUnitario,
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
// AddEntradaEpi implements EntradaEpi.
func (n *EntradaRepositorySQL) AddEntradaEpi(ctx context.Context, EntradaEpi *model.EntradaEpiInserir) error {

	query := `
		insert into entrada_epi(IdEpi,IdTamanho, data_entrada, quantidade,quantidadeAtual,data_fabricacao, data_validade, lote, fornecedor, valor_unitario)
			values (@id_epi,@id_tamanho, @data_entrada, @quantidade,@QuantidadeAtual, @dataFabricacao, @dataValidade, @lote, @fornecedor,@valorUnitario )
	`

	_, err := n.DB.ExecContext(ctx, query,
		sql.Named("id_epi", EntradaEpi.ID_epi),
		sql.Named("id_tamanho", EntradaEpi.Id_tamanho),
		sql.Named("data_entrada", EntradaEpi.Data_entrada),
		sql.Named("quantidade", EntradaEpi.Quantidade),
		sql.Named("quantidadeAtual", EntradaEpi.Quantidade_Atual),
		sql.Named("dataFabricacao", EntradaEpi.DataFabricacao),
		sql.Named("dataValidade", EntradaEpi.DataValidade),
		sql.Named("lote", EntradaEpi.Lote),
		sql.Named("fornecedor", EntradaEpi.Fornecedor),
		sql.Named("valorUnitario", EntradaEpi.ValorUnitario),
	)

	if err != nil {

		if helper.IsForeignKeyViolation(err){

			return fmt.Errorf("epi ou tamanho não existe no sistema, verifique os dados, %w", Errors.ErrDadoIncompativel)
		}
		return fmt.Errorf("erro interno ao salvar entrada, %w", err)
	}

	return nil
}

func (n *EntradaRepositorySQL) BuscarEntradaPorIdEPI(ctx context.Context, idEpi int)([]model.EntradaEpi, error){
	query:= entradaQueryJoin + " where ee.cancelada_em is null AND ee.IdEp = @id"

	entrada, err:= n.buscaEntradas(ctx, query, sql.Named("id", idEpi))
	if err != nil {

		return []model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  []model.EntradaEpi{}, nil
	}

	return  entrada, nil


}
// BuscarEntrada implements EntradaEpi.
func (n *EntradaRepositorySQL) BuscarEntrada(ctx context.Context, id int) (model.EntradaEpi, error) {

	query:= entradaQueryJoin + " where ee.cancelada_em is null AND ee.id = @id"

	entrada, err:= n.buscaEntradas(ctx, query, sql.Named("id", id))
	if err != nil {

		return model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  model.EntradaEpi{}, Errors.ErrBuscarTodos
	}

	return  entrada[0], nil

}

// BuscarTodasEntradas implements EntradaEpi.
func (n *EntradaRepositorySQL) BuscarTodasEntradas(ctx context.Context) ([]model.EntradaEpi, error) {


	query:= entradaQueryJoin + " where ee.cancelada_em is null"

	entrada, err:= n.buscaEntradas(ctx, query)
	if err != nil {

		return []model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  []model.EntradaEpi{}, nil
	}

	return  entrada, nil
}

func (n *EntradaRepositorySQL) BuscaTodasEntradasCanceladas(ctx context.Context) ([]model.EntradaEpi, error){

	query:= entradaQueryJoin + " where ee.cancelada_em is not null"

	entrada, err:= n.buscaEntradas(ctx, query)
	if err != nil {

		return []model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  []model.EntradaEpi{}, nil
	}

	return  entrada, nil

}

func (n *EntradaRepositorySQL) BuscaEntradasCanceladas(ctx context.Context, id int) (model.EntradaEpi, error){
	
	query:= entradaQueryJoin + " where ee.cancelada_em is not null AND ee.id = @id"

	entrada, err:= n.buscaEntradas(ctx, query, sql.Named("id", id))
	if err != nil {

		return model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  model.EntradaEpi{}, nil
	}

	return  entrada[0], nil

}

func (n *EntradaRepositorySQL) BuscaEntradasCanceladasPorIdEpi(ctx context.Context, idEpi int) ([]model.EntradaEpi, error){

	query:= entradaQueryJoin + " where ee.cancelada_em is not null AND ee.IdEp = @id"

	entrada, err:= n.buscaEntradas(ctx, query, sql.Named("id", idEpi))
	if err != nil {

		return []model.EntradaEpi{},err
	}

	if len(entrada) == 0 {

		return  []model.EntradaEpi{}, nil
	}

	return  entrada, nil

}
// DeletarEntrada implements EntradaEpi.
func (n *EntradaRepositorySQL) CancelarEntrada(ctx context.Context, id int) error {

	query := `update entrada_epi
			set cancelada_em = GETDATE(),
				ativo = 0
			where id = @id AND ativo = 1`

	result, err := n.DB.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("erro interno ao cancelar uma entrada %w", err)
	}

	linhas, err := result.RowsAffected()
	if err != nil {

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0 {

		return fmt.Errorf("entrada com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil
}
