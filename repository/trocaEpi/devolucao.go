package trocaepi

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/shopspring/decimal"
)

type DevolucaoInterfaceRepository interface {
	AddTrocaEPI(ctx context.Context, devolucao model.DevolucaoInserir) error
	DeleteDevolucao(ctx context.Context, id int) error
	BuscaDevoluvao(ctx context.Context, id int) (*model.Devolucao, error)
	BuscaTodasDevolucoe(ctx context.Context) ([]model.Devolucao, error)
	
}

type DevolucaoRepository struct {
	db *sql.DB
}

func NewDevolucaoRepository(db *sql.DB) DevolucaoInterfaceRepository {

	return DevolucaoRepository{
		db: db,
	}
}

// AddDevolucao implements DevolucaoInterfaceRepository.
func (d DevolucaoRepository) AddTrocaEPI(ctx context.Context, devolucao model.DevolucaoInserir) error {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transaction: %w", err)
	}

	defer tx.Rollback()
	queryInsertDevolucao := `insert into devolucao (idFuncionario, idEpi, motivo ,dataDevolucao, quantidade, idEpiNovo, IdtamanhoEpiNovo, assinaturaDigital)
	 values (@idFuncionario, @idEpi, @motivo ,@dataDevolucao, @quantidade, @idEpiNovo, @IdtamanhoEpiNovo, @assinaturaDigital)
	 OUTPUT INSERTED.id`

	var idDevolucao int64 //resgando o id da tabela devolucao
	errSqlDevolucao := tx.QueryRowContext(ctx, queryInsertDevolucao,
		sql.Named("idFuncionario", devolucao.IdFuncionario),
		sql.Named("idEpi", devolucao.IdEpi),
		sql.Named("motivo", devolucao.IdMotivo),
		sql.Named("dataDevolucao", devolucao.DataDevolucao),
		sql.Named("quantidade", devolucao.Quantidade),
		sql.Named("idEpiNovo", devolucao.IdEpiNovo),
		sql.Named("IdtamanhoEpiNovo", devolucao.IdTamanhoNovo),
		sql.Named("assinaturaDigital", devolucao.AssinaturaDigital)).Scan(&idDevolucao)
	if errSqlDevolucao != nil {

		return fmt.Errorf("erro interno ao salvar devolucao, %w", Errors.ErrInternal)
	}

	//adicionando o epi na tabela de entregas e pegando seu id
	var idEntrega int64
	queryEntrega := `

		insert into entrega (idFuncionario, dataEntrega, assinaturaDigital, idTroca)
		values (@idFuncionario, @dataEntrega, @assinaturaDigital, @idTroca)
		OUTPUT INSERTED.id
	`

	errSql := tx.QueryRowContext(ctx, queryEntrega,
		sql.Named("idFuncionario", devolucao.IdFuncionario),
		sql.Named("dataEntrega", devolucao.DataDevolucao),
		sql.Named("assinaturaDigital", devolucao.AssinaturaDigital),
		sql.Named("idTroca", idDevolucao)).Scan(&idEntrega)
	if errSql != nil {
		
		return fmt.Errorf("erro interno ao salvar entrega de um epi novo ao devolver epi antigo, %w", Errors.ErrInternal)
	}

	//pegando o valor unitario do novo epi por meio da tabela de entrada, usando o valor unitario do epi com a entrada mais antiga
	//pegando tambem o id da entrega
	var valorUnitario decimal.Decimal
	var idEntrada int64
	queryValorUnitario := `
	
		select top 1 id, valorUnitario 
		from Entrada with (updlock) 
		where id_epi = @idEpi 
			and id_tamanho = @idTamanho 
			and quantidade > 0
		order by data_entrada asc
	`

	errQueryValorUnitario := tx.QueryRowContext(ctx, queryValorUnitario, sql.Named("idEpi", devolucao.IdEpiNovo),
		sql.Named("idTamanho", devolucao.IdTamanhoNovo)).Scan(&idEntrada, &valorUnitario)

	if errQueryValorUnitario == sql.ErrNoRows {
	
		return fmt.Errorf("epi com estoque zero (id %d), tamanho %d, %w", devolucao.IdEpiNovo, devolucao.IdTamanhoNovo, Errors.ErrEstoqueInsuficiente)
	}

	if errQueryValorUnitario != nil {
		
		return fmt.Errorf("erro ao buscar entrada prioritaria: %w", Errors.ErrInternal)
	}

	//diminuido o estoque nessa entrada
	baixaEstoque := `

				update entrada set quantidade = quantidade - @quantidade
				where id = @id_entrada
				`
	_, err = tx.ExecContext(ctx, baixaEstoque, sql.Named("id_entrada", idEntrada), sql.Named("quantidade", devolucao.Quantidade))
	if err != nil {
		
		return fmt.Errorf("erro ao dar baixa no estoque da entrada %d: %w", idEntrada, Errors.ErrInternal)
	}

	//Adicionando o epi novo na tabela de epi_entregas
	queryEpiEntregas := `
		insert into epi_entregas(id_epi, id_tamanho, quantidade, id_entrega,id_entrada, valorUnitario)
		values (@idEpi, @idTamanho, @quantidade, @idEntrega, @valorUnitario, @idEntrada)
	`

	_, err = tx.ExecContext(ctx, queryEpiEntregas,
		sql.Named("idEpi", devolucao.IdEpiNovo),
		sql.Named("idTamanho", devolucao.IdTamanhoNovo),
		sql.Named("quantidade", devolucao.Quantidade),
		sql.Named("idEntrega", idEntrega),
		sql.Named("idEntrada", idEntrada),
		sql.Named("valorUnitario", valorUnitario))

		if err != nil {
			
			return fmt.Errorf("erro interno ao salvar dados na tabela epio_entregas, %w", Errors.ErrInternal)

		}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao comitar transação: %w", Errors.ErrInternal)
	}


	return nil
}

// BuscaDevoluvao implements DevolucaoInterfaceRepository.
func (d DevolucaoRepository) BuscaDevoluvao(ctx context.Context, id int) (*model.Devolucao, error) {
	panic("unimplemented")
}

// BuscaTodasDevolucoe implements DevolucaoInterfaceRepository.
func (d DevolucaoRepository) BuscaTodasDevolucoe(ctx context.Context) ([]model.Devolucao, error) {
	panic("unimplemented")
}

// DeleteDevolucao implements DevolucaoInterfaceRepository.
func (d DevolucaoRepository) DeleteDevolucao(ctx context.Context, id int) error {
	panic("unimplemented")
}
