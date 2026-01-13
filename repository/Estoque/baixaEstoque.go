package estoque

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/shopspring/decimal"
)

type BaixaEstoque interface {
	ListarLotesParaConsumo(ctx context.Context, tx *sql.Tx, epiID, tamanhoID int64) ([]LoteDisponivel, error)
	AbaterEstoqueLote(ctx context.Context, tx *sql.Tx, loteID int64, qtdParaRemover int) error
	RegistrarItemEntrega(ctx context.Context, tx *sql.Tx, idEpi, iDTamanho int64, quantidade int, idEntrega int64,
		idEntrada int64, valorUnitario decimal.Decimal) error
}

type EstoqueRepository struct {
	Db *sql.DB
}

func NewEstoqueRepository(repo *sql.DB) *EstoqueRepository {

	return &EstoqueRepository{

		Db: repo,
	}
}

type LoteDisponivel struct {
	ID            int64
	Quantidade    int
	DataValidade  *configs.DataBr
	ValorUnitario decimal.Decimal
}

// 1. Busca os lotes ordenados por validade e TRAVA as linhas (Lock)
func (r *EstoqueRepository) ListarLotesParaConsumo(ctx context.Context, tx *sql.Tx, epiID, tamanhoID int64) ([]LoteDisponivel, error) {
	// A dica de ouro do SQL Server: WITH (UPDLOCK, ROWLOCK)
	// Isso impede que outra pessoa tente pegar esses mesmos itens enquanto você decide
	query := `
        SELECT id, quantidade, data_validade, valor_unitario 
        FROM entrada_epi WITH (UPDLOCK, ROWLOCK) 
        WHERE IdEpi = @p1 
			AND IdTamanho = @p2 
			AND quantidadeAtual > 0 
			AND data_validade >= CAST(GETDATE() AS DATE)
        ORDER BY data_validade ASC
    `

	rows, err := tx.QueryContext(ctx, query, sql.Named("p1", epiID), sql.Named("p2", tamanhoID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lotes []LoteDisponivel
	for rows.Next() {
		var l LoteDisponivel
		if err := rows.Scan(&l.ID, &l.Quantidade, &l.DataValidade, &l.ValorUnitario); err != nil {
			return nil, err
		}
		lotes = append(lotes, l)
	}
	return lotes, nil
}

// 2. Atualiza a quantidade de um lote específico
func (r *EstoqueRepository) AbaterEstoqueLote(ctx context.Context, tx *sql.Tx, entradaId int64, qtdParaRemover int) error {
	query := `UPDATE entrada_epi SET quantidadeAtual = quantidadeAtual - @p1 WHERE id = @p2`
	result, err := tx.ExecContext(ctx, query, sql.Named("p1", qtdParaRemover), sql.Named("p2", entradaId))
	if err != nil {
		
		 return fmt.Errorf("erro interno, %w", Errors.ErrInternal)

	}

	linha, err :=result.RowsAffected()
	if err != nil {

		return 	fmt.Errorf("erro ao verificar o numero de linhas")
	}

	if linha == 0 {

		return fmt.Errorf("id nao existe no sistema, %w", Errors.ErrDadoIncompativel)
	}
	
	return  nil
}

func (r *EstoqueRepository) RegistrarItemEntrega(ctx context.Context, tx *sql.Tx, idEpi, iDTamanho int64, quantidade int, idEntrega int64,idEntrada int64, valorUnitario decimal.Decimal) error {
 	queryItens := `

			insert into epis_entregues(IdEpi,IdTamanho, quantidade,IdEntrega,IdEntrada ,valor_unitario) values (@id_epi, @id_tamanho, @quantidade, @id_entrega,@id_entrada ,@valorUnitario)
		`

	_, err:= tx.ExecContext(ctx, queryItens,
		sql.Named("id_epi", idEpi),
		sql.Named("id_tamanho", iDTamanho),
		sql.Named("quantidade", quantidade),
		sql.Named("id_entrega", idEntrega),
		sql.Named("id_entrada", idEntrada),
		sql.Named("valorUnitario", valorUnitario))

	if err != nil {
		return fmt.Errorf("erro ao inserir epi nas entregas, %w", err)
	}

	return nil
}
