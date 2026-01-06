package entregaepi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	trocaepi "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/trocaEpi"
)

type EntregaRepository struct {
	Db           *sql.DB
	BaixaEstoque trocaepi.DevolucaoInterfaceRepository
}

func NewEntregaRepository(db *sql.DB, RepoDevolucao trocaepi.DevolucaoInterfaceRepository) *EntregaRepository {

	return &EntregaRepository{
		Db:           db,
		BaixaEstoque: RepoDevolucao,
	}

}

const entregaQueryJoin = `
select
		    ee.id,
			ee.data_Entrega,
			ee.IdFuncionario,
			f.nome, 
			f.IdDepartamento, 
			d.nome, 
			f.IdFuncao, 
			ff.nome, 
			i.id, 
			e.nome, 
			e.fabricante, 
			e.CA,
			e.descricao,  
			e.validade_CA,
			e.IdTipoProtecao,
			tp.id, 
			i.IdTamanho, 
			t.tamanho, 
			i.quantidade,
			ee.assinatura,
			i.valor_unitario
			from entrega_epi ee
			inner join
				funcionario f on ee.IdFuncionario = f.id
			inner join
				departamento d on f.IdDepartamento = d.id
			inner join 
				funcao ff on f.IdFuncao = ff.id
			inner join 
				epis_entregues i on i.IdEntrega = ee.id
			inner join 
				epi e on i.IdEpi = e.id
			inner join
				tipo_protecao tp on e.IdTipoProtecao = tp.id
			inner join 
				tamanho t on i.IdTamanho = t.id
			 

`

func (n *EntregaRepository) buscaEntregas(ctx context.Context, query string, args ...any) ([]*model.EntregaDto, error) {
	linhas, err := n.Db.QueryContext(ctx, query, args...)
	if err != nil {

		return nil, fmt.Errorf("falha ao buscas as entregas, %w", Errors.ErrBuscarTodos)
	}
	defer linhas.Close()

	EntregaMap := make(map[int]*model.EntregaDto)

	for linhas.Next() {
		var item model.ItemEntregueDto
		var entregaID int
		var dataEntrega configs.DataBr
		var assinatura string
		var funcID int
		var funcNome string
		var depID int
		var depNome string
		var funcaoID int
		var funcaoNome string

		err := linhas.Scan(
			&entregaID,
			&dataEntrega,
			&funcID,
			&funcNome,
			&depID,
			&depNome,
			&funcaoID,
			&funcaoNome,
			&item.Epi.Id,
			&item.Epi.Nome,
			&item.Epi.Fabricante,
			&item.Epi.CA,
			&item.Epi.Descricao,
			&item.Epi.DataValidadeCa,
			&item.Epi.Protecao.ID,
			&item.Epi.Protecao.Nome,
			&item.Tamanho.ID,
			&item.Tamanho.Tamanho,
			&item.Quantidade,
			&assinatura,
			&item.ValorUnitario,
		)
		if err != nil {

			return nil, fmt.Errorf("a  %w", Errors.ErrFalhaAoEscanearDados)
		}

		if _, ok := EntregaMap[entregaID]; !ok {

			EntregaMap[entregaID] = &model.EntregaDto{ //verifico pelo id, se o map da entrega ja existe, se não, cria um

				Id: entregaID,
				Funcionario: model.Funcionario_Dto{
					ID:           funcID,
					Nome:         funcNome,
					Funcao: model.FuncaoDto{
						ID:     funcaoID,
						Funcao: funcaoNome,
						Departamento: model.DepartamentoDto{
							ID:           depID,
							Departamento: depNome,
						},
					},
				},
				Data_entrega:       dataEntrega,
				Assinatura_Digital: assinatura,
				Itens:              []model.ItemEntregueDto{},
			}
		}

		EntregaMap[entregaID].Itens = append(EntregaMap[entregaID].Itens, item) //caso ja exista, faço um append dos itens
	}

	if err := linhas.Err(); err != nil {

		return nil, fmt.Errorf("%w", Errors.ErrAoIterar)
	}

	EntregaSlice := make([]*model.EntregaDto, 0, len(EntregaMap))

	for _, entrega := range EntregaMap {
		EntregaSlice = append(EntregaSlice, entrega)
	}

	return EntregaSlice, nil

}

// Addentrega implements EntregaInterface.
func (n *EntregaRepository) Addentrega(ctx context.Context, model model.EntregaParaInserir) error {

	tx, err := n.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar a transação, %w", err)
	}
	defer tx.Rollback()

	//add as entregas
	queryEntrega := `insert into entrega_epi(IdFuncionario, data_entrega, assinatura)
	OUTPUT INSERTED.id
	 values (@idFuncionario, @dataEntrega, CAST( @AssinaturaDigital AS VARBINARY(MAX)))
	`

	var idEntrega int64
	errSql := tx.QueryRowContext(ctx, queryEntrega,
		sql.Named("idFuncionario", model.ID_funcionario),
		sql.Named("dataEntrega", model.Data_entrega),
		sql.Named("AssinaturaDigital", model.Assinatura_Digital),
	).Scan(&idEntrega)
	if errSql != nil {

		return fmt.Errorf("erro interno ao salvar entrega, %w", Errors.ErrInternal)
	}

	for _, item := range model.Itens {

		err := n.BaixaEstoque.BaixaEstoque(ctx, tx, int64(item.ID_epi), int64(item.ID_tamanho), item.Quantidade, idEntrega)
		if err != nil {
			return err
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao comitar transação: %w", Errors.ErrInternal)
	}

	return nil
}

// BuscaEntrega implements EntregaInterface.
func (n *EntregaRepository) BuscaEntrega(ctx context.Context, id int) (*model.EntregaDto, error) {

	query := entregaQueryJoin + " where ee.cancelada_em IS NULL and ee.id = @id"

	entrega, err := n.buscaEntregas(ctx, query, sql.Named("id", id))
	if err != nil {
		return nil, err
	}

	if len(entrega) == 0 {

		return nil, nil
	}

	return entrega[0], nil

}

func (n *EntregaRepository) BuscaEntregaPorIdFuncionario(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error) {

	query := entregaQueryJoin + " where ee.cancelada_em IS NULL and ee.IdFuncionario = @id"

	entrega, err := n.buscaEntregas(ctx, query, sql.Named("id", idFuncionario))
	if err != nil {
		return nil, err
	}

	if len(entrega) == 0 {

		return nil, nil
	}

	return entrega, nil

}

func (n *EntregaRepository) BuscaEntregaPorIdFuncionarioCanceladas(ctx context.Context, idFuncionario int) ([]*model.EntregaDto, error) {

	query := entregaQueryJoin + " where ee.cancelada_em IS NOT NULL and ee.IdFuncionario = @id"

	entrega, err := n.buscaEntregas(ctx, query, sql.Named("id", idFuncionario))
	if err != nil {
		return nil, err
	}

	if len(entrega) == 0 {

		return nil, nil
	}

	return entrega, nil

}

// BuscaTodasEntregas implements EntregaInterface.
func (n *EntregaRepository) BuscaTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error) {

	return n.buscaEntregas(ctx, entregaQueryJoin)

}

func (n *EntregaRepository) BuscaTodasEntregasCanceladas(ctx context.Context) ([]*model.EntregaDto, error) {

	query := entregaQueryJoin + " where ee.cancelada_em IS not NULL"

	entrega, err := n.buscaEntregas(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(entrega) == 0 {

		return nil, nil
	}

	return entrega, nil
}

func (n *EntregaRepository) BuscaEntregaCancelada(ctx context.Context, id int) (*model.EntregaDto, error) {

	query := entregaQueryJoin + " where ee.cancelada_em IS not NULL and ee.id = @id"

	entrega, err := n.buscaEntregas(ctx, query, sql.Named("id", id))
	if err != nil {
		return nil, err
	}

	if len(entrega) == 0 {

		return nil, nil
	}

	return entrega[0], nil

}

// DeletarEntregas implements EntregaInterface.
func (n *EntregaRepository) CancelarEntrega(ctx context.Context, id int) error {

	query := `update entrega
			set cancelada_em  = GETDATE(), ativo = 0
			where id = @id AND cancelada_em IS NULL;`

	result, err := n.Db.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {

		return fmt.Errorf("%w", Errors.ErrInternal)
	}

	linha, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)

	}

	if linha == 0 {

		return fmt.Errorf("entrega nao encontrada, %w", Errors.ErrNaoEncontrado)
	}

	return nil

}
