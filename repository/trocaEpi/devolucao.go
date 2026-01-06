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
	AddDevolucaoEpi(ctx context.Context, devolucao model.DevolucaoInserir) error
	DeleteDevolucao(ctx context.Context, id int) error
	BuscaDevolucaoPorMatricula(ctx context.Context, matricula int) ([]model.Devolucao, error)
	BuscaTodasDevolucoes(ctx context.Context) ([]model.Devolucao, error)
	BaixaEstoque(ctx context.Context, tx *sql.Tx, idEpi, iDTamanho int64, quantidade int, idEntrega int64) error
	BuscaDevolucaoPorId(ctx context.Context, id int)([]model.Devolucao,error)
	BuscaDevolucaoPorIdCancelada(ctx context.Context, id int)([]model.Devolucao, error)
	BuscaDevolucaoPorMatriculaCancelada(ctx context.Context, matricula int) ([]model.Devolucao, error)
	BuscaTodasDevolucoesCancelada(ctx context.Context) ([]model.Devolucao, error)
}

type DevolucaoRepository struct {
	db *sql.DB
}

func NewDevolucaoRepository(db *sql.DB) DevolucaoInterfaceRepository {

	return &DevolucaoRepository{
		db: db,
	}
}

const queryBuscaDevolucao = `
select 
		d.id,
		d.idFuncionario,
		f.nome,
		f.idDepartamento,
		dd.nome,
		f.idFuncao,
		ff.nome,
		d.IdEpi, 
		e.nome, 
		e.fabricante, 
		e.CA,
		d.IdTamanho as id_tamanhoAntigo,
		t.tamanho as tamanhoAntigo,
		d.quantidadeAdevolver,
		d.IdMotivo,
		d.IdEpiNovo as epiNovo, 
		en.nome as NomeEpiNovo, 
		en.fabricante as FabricanteEpiNovo, 
		en.CA as CAEpiNovo,
		d.quantidadeNova as QuantidadeEpiNovo,
		d.IdTamanhoNovo as TamanhoEpiNovo,
		tn.tamanho as TamanhoNovo,
		d.assinatura_digital,
		d.data_devolucao
		from devolucao d
		inner join
			epi e on  d.IdEpi = e.id
		left join
			epi en on d.IdEpiNovo = en.id
		inner join
			funcionario f on d.idFuncionario = f.id	
		inner join 
			departamento dd on f.idDepartamento = dd.id
		inner join 
			funcao ff on f.idFuncao = ff.id
		left join
			tamanho tn on d.IdTamanhoNovo = tn.id
		inner join
			tamanho t on d.IdTamanho = t.id
		where d.cancelada_em is null

`

func (d *DevolucaoRepository) executaBusca(ctx context.Context, query string, args ...any) ([]model.Devolucao, error) {
	linhas, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar devolucões do colaborador, %w", Errors.ErrBuscarTodos)
	}

	var Devolucao []model.Devolucao
	defer linhas.Close()

	for linhas.Next() {

		var devolucao model.Devolucao

		err := linhas.Scan(
			&devolucao.Id,
			&devolucao.Id_funcionario,
			&devolucao.NomeFuncionario,
			&devolucao.Id_departamento,
			&devolucao.Departamento,
			&devolucao.Id_funcao,
			&devolucao.Funcao,
			&devolucao.ID_epiTroca,
			&devolucao.NomeEpiTroca,
			&devolucao.FabricanteTroca,
			&devolucao.CAtroca,
			&devolucao.IdTamanho,
			&devolucao.Tamanho,
			&devolucao.QuantidadeADevolver,
			&devolucao.Motivo,
			&devolucao.ID_epiNovo,
			&devolucao.NomeEpiNovo,
			&devolucao.FabricanteNovo,
			&devolucao.CANovo,
			&devolucao.NovaQuantidade,
			&devolucao.Id_tamanhoNovo,
			&devolucao.TamanhoNovo,
			&devolucao.AssinaturaDigital,
			&devolucao.DataEntrega,
		)

		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		Devolucao = append(Devolucao, devolucao)
	}

	if err := linhas.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre devolucoes, %w", Errors.ErrAoIterar)
	}

	return Devolucao, nil

}

func (d *DevolucaoRepository) BaixaEstoque(ctx context.Context, tx *sql.Tx, idEpi, iDTamanho int64, quantidade int, idEntrega int64) error {

	//buscando o lote mais antigo e que tenha o saldo de epi nescessario
	buscaLote := `

	select  top 1 id as entrada, valor_unitario ,lote
	from entrada_epi
	where quantidade >= 1 and ativo = 1
	order by data_entrada asc

		`

	var idEntrada int64
	var valorUnitario decimal.Decimal
	var saldoLote int

	err := tx.QueryRowContext(ctx, buscaLote,
		sql.Named("idEpi", idEpi),
		sql.Named("id_tamanho", iDTamanho),
		sql.Named("quantidade", quantidade)).Scan(&idEntrada, &valorUnitario, &saldoLote) //adicionando os valores nessas variaveis

	if err == sql.ErrNoRows {
		return fmt.Errorf("estoque  zero para o epi %d (tamanho %d), %w", idEpi, iDTamanho, Errors.ErrEstoqueInsuficiente)
	}

	if err != nil {

		return fmt.Errorf("erro ao dar baixa no estoque, %w", Errors.ErrInternal)
	}

	//atualizando o saldo
	queryBaixa := `

			update entrada_epi 
				set quantidade = quantidade - @qtd 
					where id = @idEntrada and ativo = 1
		`

	_, err = tx.ExecContext(ctx, queryBaixa,
		sql.Named("qtd", quantidade),
		sql.Named("idEntrada", idEntrada))
	if err != nil {

		return fmt.Errorf("erro ao atualizar estoque da entrada %d, %w", idEntrada, Errors.ErrInternal)
	}

	//dando entrada dos epi na tabela auxiliar epi_entrega
	queryItens := `

			insert into epis_entregues(IdEpi,IdTamanho, quantidade,IdEntrega,IdEntrada ,valor_unitario) values (@id_epi, @id_tamanho, @quantidade, @id_entrega,@id_entrada ,@valorUnitario)
		`

	_, err = tx.ExecContext(ctx, queryItens,
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

// AddDevolucaoEpi implements DevolucaoInterfaceRepository.
func (d *DevolucaoRepository) AddDevolucaoEpi(ctx context.Context, devolucao model.DevolucaoInserir) error {

	query := `
		insert into devolucao (IdFuncionario, IdEpi, IdMotivo ,data_devolucao, IdTamanho, quantidadeAdevolver, idEpiNovo, IdTamanhoNovo,quantidade_nova,assinatura_digital)
	 	values (@idFuncionario, @idEpi, @motivo ,@dataDevolucao,@idTamanho ,@quantidadeDevolucao, null,null, null, cast( @assinaturaDigital as varbinary(max)))

	`
	_, err := d.db.ExecContext(ctx, query, sql.Named("idFuncionario", devolucao.IdFuncionario),
		sql.Named("idEpi", devolucao.IdEpi),
		sql.Named("motivo", devolucao.IdMotivo),
		sql.Named("dataDevolucao", devolucao.DataDevolucao),
		sql.Named("idTamanho",devolucao.IdTamanho),
		sql.Named("quantidadeDevolucao", devolucao.QuantidadeADevolver),
		sql.Named("assinaturaDigital", devolucao.AssinaturaDigital))
	if err != nil {

		return fmt.Errorf("erro interno ao salvar devolucao do epi, %w", Errors.ErrInternal)
	}

	return nil
}

// AddDevolucao implements DevolucaoInterfaceRepository.
func (d *DevolucaoRepository) AddTrocaEPI(ctx context.Context, devolucao model.DevolucaoInserir) error {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transaction: %w", err)
	}

	defer tx.Rollback()
	queryInsertDevolucao := `insert into devolucao (IdFuncionario, IdEpi, IdMotivo ,data_devolucao, IdTamanho, quantidadeAdevolver, idEpiNovo, IdTamanhoNovo,quantidade_nova,assinatura_digital)
		OUTPUT INSERTED.id
	 	values (@idFuncionario, @idEpi, @motivo ,@dataDevolucao,@idTamanho ,@quantidadeDevolucao, @idEpiNovo,@IdtamanhoEpiNovo, @quantidadeNova, cast( @assinaturaDigital as varbinary(max)))
	 `

	var idDevolucao int64 //resgando o id da tabela devolucao
	errSqlDevolucao := tx.QueryRowContext(ctx, queryInsertDevolucao,
		sql.Named("idFuncionario", devolucao.IdFuncionario),
		sql.Named("idEpi", devolucao.IdEpi),
		sql.Named("motivo", devolucao.IdMotivo),
		sql.Named("dataDevolucao", devolucao.DataDevolucao),
		sql.Named("idTamanho", devolucao.IdTamanho),
		sql.Named("quantidadeDevolucao", devolucao.QuantidadeADevolver),
		sql.Named("idEpiNovo", devolucao.IdEpiNovo),
		sql.Named("IdtamanhoEpiNovo", devolucao.IdTamanhoNovo),
		sql.Named("quantidadeNova", devolucao.NovaQuantidade),
		sql.Named("assinaturaDigital", devolucao.AssinaturaDigital)).Scan(&idDevolucao)
	if errSqlDevolucao != nil {

		return fmt.Errorf("erro interno ao salvar devolucao, %w", errSqlDevolucao)
	}

	//adicionando o epi na tabela de entregas e pegando seu id
	var idEntrega int64
	queryEntrega := `

		insert into entrega_epi(IdFuncionario, data_entrega, assinatura, IdTroca)
		OUTPUT INSERTED.id
		values (@idFuncionario, @dataEntrega, cast(@assinaturaDigital as varbinary(max)), @idTroca)
	`

	errSql := tx.QueryRowContext(ctx, queryEntrega,
		sql.Named("idFuncionario", devolucao.IdFuncionario),
		sql.Named("dataEntrega", devolucao.DataDevolucao),
		sql.Named("assinaturaDigital", devolucao.AssinaturaDigital),
		sql.Named("idTroca", idDevolucao)).Scan(&idEntrega)
	if errSql != nil {

		return fmt.Errorf("erro interno ao salvar entrega de um epi novo ao devolver epi antigo, %w", Errors.ErrInternal)
	}

	err = d.BaixaEstoque(ctx, tx, int64(*devolucao.IdEpiNovo), int64(*devolucao.IdTamanhoNovo), *devolucao.NovaQuantidade, idEntrega)
	if err != nil {

		return err
	}

	err = tx.Commit()
	if err != nil {

		return fmt.Errorf("erro ao realizar o commit da transação, %w", Errors.ErrInternal)
	}

	return nil
}

//******************************************************************************

// executaBusca implements DevolucaoInterfaceRepository.

// BuscaDevoluvao implements DevolucaoInterfaceRepository.
func (d *DevolucaoRepository) BuscaDevolucaoPorMatricula(ctx context.Context, matricula int) ([]model.Devolucao, error) {

	query := queryBuscaDevolucao + " where f.matricula = @matricula and d.cancelada_em is null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query, sql.Named("matricula", matricula))
}

// BuscaTodasDevolucoe implements DevolucaoInterfaceRepository.
func (d *DevolucaoRepository) BuscaTodasDevolucoes(ctx context.Context) ([]model.Devolucao, error) {

	query := queryBuscaDevolucao + " d.cancelada_em is null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query)
}

func (d *DevolucaoRepository) BuscaDevolucaoPorId(ctx context.Context, id int)([]model.Devolucao,error) {

	query := queryBuscaDevolucao + " where d.id = @id and d.cancelada_em is null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query, sql.Named("id", id))

}

func (d *DevolucaoRepository) BuscaDevolucaoPorIdCancelada(ctx context.Context, id int)([]model.Devolucao, error){

	query := queryBuscaDevolucao + " where d.id = @id and d.cancelada_em is not null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query, sql.Named("id", id))

}

func (d *DevolucaoRepository) BuscaDevolucaoPorMatriculaCancelada(ctx context.Context, matricula int) ([]model.Devolucao, error) {

	query := queryBuscaDevolucao + " where f.matricula = @matricula and d.cancelada_em is not null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query, sql.Named("matricula", matricula))
}

func (d *DevolucaoRepository) BuscaTodasDevolucoesCancelada(ctx context.Context) ([]model.Devolucao, error) {

	query := queryBuscaDevolucao + " d.cancelada_em is not null order by d.dataTroca DESC"

	return d.executaBusca(ctx, query)
}



//***************************************************************************************************

// DeleteDevolucao implements DevolucaoInterfaceRepository.
func (d *DevolucaoRepository) DeleteDevolucao(ctx context.Context, id int) error {
	
	query := `
	
		update devolucao
		set cancelada_em = GETDATE(), ativo = 0
		where id = @id and cancelada_em is null
	`

	result, err:= d.db.ExecContext(ctx, query,  sql.Named("id", id))

	if err != nil {

		return fmt.Errorf("erro interno ao cancelar devolucao, %w", Errors.ErrInternal)
	}

	linhas, err:= result.RowsAffected()
	if err != nil {

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0{

		return fmt.Errorf("entrada com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil

}
