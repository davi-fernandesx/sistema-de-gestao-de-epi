package entregaepi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EntregaInterface interface {
	Addentrega(ctx context.Context, model model.EntregaParaInserir) error
	BuscaEntrega(ctx context.Context, id int) (*model.EntregaDto, error)
	BuscaTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error)
	CancelarEntrega(ctx context.Context, id int) error
}

type NewsqlLogin struct {
	Db *sql.DB
}

func NewEntregaRepository(db *sql.DB) EntregaInterface {

	return &NewsqlLogin{
		Db: db,
	}

}

// Addentrega implements EntregaInterface.
func (n *NewsqlLogin) Addentrega(ctx context.Context, model model.EntregaParaInserir) error {
	
	tx, err:= n.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar a transação, %w", err)
	}
	defer tx.Rollback()

	queryEntrega:=`insert into entrega (id_funcionario, data_entrega, AssinaturaDigital)
	 values (@idFuncionario, @dataEntrega, @AssinaturaDigital)
	 OUTPUT INSERTED.id`

	 var id int64
	 errSql:= tx.QueryRowContext(ctx, queryEntrega,
		sql.Named("idFuncionario", model.ID_funcionario),
		sql.Named("dataEntrega", model.Data_entrega),
		sql.Named("AssinaturaDigital", model.Assinatura_Digital),

	).Scan(&id)
	if errSql != nil {

		return fmt.Errorf("erro interno ao salvar entrega, %w", Errors.ErrInternal)
	}

	queryItem:= `insert into epi_entregas(id_epi,id_tamanho, quantidade, id_entrega) values (@id_epi, @id_tamanho, @quantidade, @id_entrega)`
	itens, errStmt:= tx.PrepareContext(ctx, queryItem)
	if errStmt != nil {
		return fmt.Errorf("erro interno ao preparar itens de entrega, %w", Errors.ErrInternal )
	}
	defer itens.Close()

	for _, item:= range model.Itens{

		_, err:= itens.ExecContext(ctx, sql.Named("id_epi", item.ID_epi), sql.Named("id_tamanho", 
		item.ID_tamanho),sql.Named("quantidade", item.Quantidade), sql.Named("id_entrega", id))
		if err != nil {
				return  fmt.Errorf("erro interno ao salvar itens para entrega, %w", Errors.ErrInternal)
		}
	}

	if err:= tx.Commit(); err != nil {
		return  fmt.Errorf("erro ao comitar transação: %w", Errors.ErrInternal)
	}

	return nil
}

// BuscaEntrega implements EntregaInterface.
func (n *NewsqlLogin) BuscaEntrega(ctx context.Context, id int) (*model.EntregaDto, error) {
	
	query:=`select
		    ee.id,
			ee.data_Entrega,
			ee.id_funcionario,
			f.nome, 
			f.id_departamento, 
			d.departamento, 
			f.id_funcao, 
			ff.funcao, 
			i.id_epi, 
			e.nome, 
			e.fabricante, 
			e.CA,
			e.descricao, 
			e.data_fabricacao, 
			e.data_validade, 
			e.data_validadeCa,
			e.id_tipo_protecao,
			tp.protecao, 
			i.id_tamanho, 
			t.tamanho, 
			i.quantidade,
			ee.AssinaturaDigital
			from entrega ee
			inner join
				funcionario f on ee.id_funcionario = f.id
			inner join
				departamentos d on f.id_departamento = d.id
			inner join 
				funcao ff on f.id_funcao = ff.id
			inner join 
				epi_entregues i on i.id_entrega = ee.id
			inner join 
				epi e on i.id_epi = e.id
			inner join
				tipo_protecao tp on e.id_tipo_protecao = tp.id
			inner join 
				tamanho t on i.id_tamanho = t.id
			where ee.cancelada_em IS NULL and ee.id = @id
		`

	var entrega *model.EntregaDto

	linhas, err:= n.Db.QueryContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return  nil, fmt.Errorf("erro ao buscar entregas, %w", Errors.ErrBuscarTodos)
	}
	defer linhas.Close()

	for linhas.Next() {

		var item model.ItemEntregueDto
		var entregaID int
        var dataEntrega time.Time
        var assinatura string
        var funcID int
        var funcNome string
        var depID int
        var depNome string
        var funcaoID int
        var funcaoNome string

		err:= linhas.Scan(
			&entregaID,
            &dataEntrega,
            &funcID,
            &funcNome,
            &depID,
            &depNome,
            &funcaoID,
            &funcaoNome,
            &item.Epi.Id, // Assumindo que EpiDto tem esses campos
            &item.Epi.Nome,
            &item.Epi.Fabricante,
            &item.Epi.CA,
            &item.Epi.Descricao,
            &item.Epi.DataFabricacao,
            &item.Epi.DataValidade,
            &item.Epi.DataValidadeCa,
            &item.Epi.Protecao.ID, // Assumindo estrutura aninhada
            &item.Epi.Protecao.Nome,
            &item.Tamanho.ID, // Assumindo que TamanhoDto tem Id e Nome
            &item.Tamanho.Tamanho,
            &item.Quantidade,
			&assinatura,
		); if err!= nil {

			return  nil, fmt.Errorf(" %w", Errors.ErrFalhaAoEscanearDados)
		}

		if entrega == nil {

			entrega = &model.EntregaDto{
					Id: entregaID,
					Funcionario: model.Funcionario_Dto{
						ID: funcID,
						Nome: funcNome,
						Departamento: model.DepartamentoDto{
							ID: depID,
							Departamento: depNome,
						},
						Funcao: model.FuncaoDto{
							ID: funcID,
							Funcao: funcaoNome,
						},	
					},
					Data_entrega: dataEntrega,
					Assinatura_Digital: assinatura,
					Itens: []model.ItemEntregueDto{},
			}
		}

		entrega.Itens = append(entrega.Itens, item )
	}

	if err:= linhas.Err(); err != nil {

		return nil, fmt.Errorf("%w", Errors.ErrAoIterar)
	}
	if entrega == nil {
    return nil, Errors.ErrNaoEncontrado // Retorna o erro específico
}

	return  entrega, nil
	
}

// BuscaTodasEntregas implements EntregaInterface.
func (n *NewsqlLogin) BuscaTodasEntregas(ctx context.Context) ([]*model.EntregaDto, error) {
	query:= `select
		    ee.id,
			ee.data_Entrega,
			ee.id_funcionario,
			f.nome, 
			f.id_departamento, 
			d.departamento, 
			f.id_funcao, 
			ff.funcao, 
			i.id_epi, 
			e.nome, 
			e.fabricante, 
			e.CA,
			e.descricao, 
			e.data_fabricacao, 
			e.data_validade, 
			e.data_validadeCa,
			e.id_tipo_protecao,
			tp.protecao, 
			i.id_tamanho, 
			t.tamanho, 
			i.quantidade,
			ee.AssinaturaDigital
			from entrega ee
			inner join
				funcionario f on ee.id_funcionario = f.id
			inner join
				departamentos d on f.id_departamento = d.id
			inner join 
				funcao ff on f.id_funcao = ff.id
			inner join 
				epi_entregues i on i.id_entrega = ee.id
			inner join 
				epi e on i.id_epi = e.id
			inner join
				tipo_protecao tp on e.id_tipo_protecao = tp.id
			inner join 
				tamanho t on i.id_tamanho = t.id
			where
				 ee.cancelada_em IS NULL
			ORDER BY ee.id`

		linhas, err:= n.Db.QueryContext(ctx, query)
		if err != nil {

			return  nil, fmt.Errorf("falha ao buscas as entregas, %w", Errors.ErrBuscarTodos)
		}
		defer linhas.Close()

		EntregaMap := make(map[int]*model.EntregaDto)

		for linhas.Next(){
		var item model.ItemEntregueDto
		var entregaID int
        var dataEntrega time.Time
        var assinatura string
        var funcID int
        var funcNome string
        var depID int
        var depNome string
        var funcaoID int
        var funcaoNome string

		err:= linhas.Scan(
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
            &item.Epi.DataFabricacao,
            &item.Epi.DataValidade,
            &item.Epi.DataValidadeCa,
            &item.Epi.Protecao.ID, 
            &item.Epi.Protecao.Nome,
            &item.Tamanho.ID,
            &item.Tamanho.Tamanho,
            &item.Quantidade,
			  &assinatura,
		); if err!= nil {

			return  nil, fmt.Errorf(" %w", Errors.ErrFalhaAoEscanearDados)
		}

		if _ , ok := EntregaMap[entregaID]; !ok {

			EntregaMap[entregaID] = &model.EntregaDto{ //verifico pelo id, se o map da entrega ja existe, se não, cria um

					Id: entregaID,
					Funcionario: model.Funcionario_Dto{
						ID: funcID,
						Nome: funcNome,
						Departamento: model.DepartamentoDto{
							ID: depID,
							Departamento: depNome,
						},
						Funcao: model.FuncaoDto{
							ID: funcID,
							Funcao: funcaoNome,
						},	
					},
					Data_entrega: dataEntrega,
					Assinatura_Digital: assinatura,
					Itens: []model.ItemEntregueDto{},
			
			    }
		}

		EntregaMap[entregaID].Itens = append(EntregaMap[entregaID].Itens, item) //caso ja exista, faço um append dos itens
	}

	if err:= linhas.Err(); err != nil {

		return  nil, fmt.Errorf("%w", Errors.ErrAoIterar)
	}

	EntregaSlice := make([]*model.EntregaDto, 0, len(EntregaMap))

	for _, entrega:= range EntregaMap {
		EntregaSlice = append(EntregaSlice, entrega)
	}

	if err == sql.ErrNoRows {

		return nil, fmt.Errorf("nenhuma entrega encontrada, %w", Errors.ErrBuscarTodos)
	}


	return EntregaSlice, nil
}

// DeletarEntregas implements EntregaInterface.
func (n *NewsqlLogin) CancelarEntrega(ctx context.Context, id int) error {
	
	query:= `update entrega
			set cancelada_em  = GETDATE() 
			where id = @id AND cancelada_em IS NULL;`

	result, err:= n.Db.ExecContext(ctx, query, sql.Named("id",id))
	if err != nil {

		return  fmt.Errorf("%w", Errors.ErrInternal)
	}

	 linha, err:= result.RowsAffected()
	 if err != nil {

			if errors.Is(err, Errors.ErrLinhasAfetadas){

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		}
	 }

	 if linha == 0 {

		return fmt.Errorf("entrega nao encontrada, %w", Errors.ErrNaoEncontrado)
	 }


	 return  nil
	
}


