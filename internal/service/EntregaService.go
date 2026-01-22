package service

import (
	"context"
	"fmt"
	"math"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type EntregaRepository interface {
	AdicionarEntrega(ctx context.Context, qtx *repository.Queries, args repository.AddEntregaEpiParams) (int32, error)
	AdicionarEntregaItem(ctx context.Context, qtx *repository.Queries, arg repository.AddItemEntregueParams) (int32, error)
	ListarEntregas(ctx context.Context, args repository.ListarEntregasParams) ([]repository.ListarEntregasRow, error)
	Cancelar(ctx context.Context,qtx *repository.Queries,args repository.CancelarEntregaParams) (int32, error)
	CancelarEntregaItem(ctx context.Context, qtx *repository.Queries, id int32) error
	AbaterEstoqueEntrada(ctx context.Context, qtx *repository.Queries, args repository.AbaterEstoqueLoteParams) (int64, error)
	ReporEstoqueEntrada(ctx context.Context, qtx *repository.Queries, args repository.ReporEstoqueLoteParams) (int64, error)
	ListarEntregasDisponiveis(ctx context.Context, qtx *repository.Queries, args repository.ListarLotesParaConsumoParams) ([]repository.ListarLotesParaConsumoRow, error)
	ListarEpisEntreguesCancelados(ctx context.Context,qtx *repository.Queries ,id int32) ([]repository.ListarItensEntregueCanceladosRow, error)
}

type EntregaService struct {
	repo    EntregaRepository
	db      *pgxpool.Pool
	queries *repository.Queries
}

func NewEntregaService(r EntregaRepository, pool *pgxpool.Pool) *EntregaService {

	return &EntregaService{
		repo: r,
		db: pool,
		queries: repository.New(pool),
	}
}

func (e *EntregaService) Salvar(ctx context.Context, model model.EntregaParaInserir) error {

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	funcionario, err := e.queries.BuscaFuncionarioPorId(ctx, int32(model.ID_funcionario))
	if err != nil {

		return err
	}
	token := helper.GerarTokenAuditoria(funcionario.Nome, funcionario.FuncaoNome, funcionario.DepartamentoNome, model.Data_entrega.Time())

	qtx := e.queries.WithTx(tx)

	args := repository.AddEntregaEpiParams{

		Idfuncionario:  int32(model.ID_funcionario),
		DataEntrega:    pgtype.Date{Time: model.Data_entrega.Time(), Valid: true},
		Assinatura:     model.Assinatura_Digital,
		TokenValidacao: pgtype.Text{String: token},
		IDUsuarioEntrega: pgtype.Int4{Int32: int32(model.Id_user)},
	}

	identrega, err := e.repo.AdicionarEntrega(ctx, qtx, args) //salva o "cabeçalho"
	if err != nil {

		return err
	}

	//percorre todos os item da lista de itens
	for _, item := range model.Itens {

		quantidadeNescessaria := item.Quantidade

		lotes := repository.ListarLotesParaConsumoParams{
			Idepi:     int32(item.ID_epi),
			Idtamanho: int32(item.ID_tamanho),
		}
		/*lista todas as entradas com quantidadeAtual maior que 0 e que tenha os idepie e idtamanhos iguais as passado nos parametros*/
		entradaLotes, err := e.repo.ListarEntregasDisponiveis(ctx, qtx, lotes)
		if err != nil {
			return err
		}

		if len(entradaLotes) == 0 {
             return fmt.Errorf("estoque insuficiente para o EPI ID %d", item.ID_epi)
		}

		/*percorre todas as entradas achadas*/
		for _, entradaLote := range entradaLotes {

			if quantidadeNescessaria <= 0 {
				break
			}

			//escolhe o menor valor entre esses parametros
			quantidadeAbater := min(entradaLote.Quantidadeatual, int32(quantidadeNescessaria))

			itemAdd := repository.AddItemEntregueParams{
				Identrega:     identrega,
				Idepi:         int32(item.ID_epi),
				Idtamanho:     int32(item.ID_tamanho),
				Quantidade:    quantidadeAbater,
				ValorUnitario: entradaLote.ValorUnitario,
				Identrada:     entradaLote.ID,
			}

			_, err := e.repo.AdicionarEntregaItem(ctx, qtx, itemAdd)
			if err != nil {
				return err
			}

			_, err = e.repo.AbaterEstoqueEntrada(ctx, qtx, repository.AbaterEstoqueLoteParams{
				Quantidadeatual: quantidadeAbater,
				ID:              entradaLote.ID,
			})
			if err != nil {
				return err
			}

			quantidadeNescessaria -= int(quantidadeAbater)
		}

		if quantidadeNescessaria > 0 {
			// Se sobrou quantidade, significa que percorremos todos os lotes
			// e ainda não deu o total. Rollback automático pelo defer!
			return fmt.Errorf("estoque insuficiente para o EPI ID %d (faltam %d unidades)",
				item.ID_epi, quantidadeNescessaria)
		}
	}

	return tx.Commit(ctx)
}

type FiltroEntregas struct {
	Canceladas    bool
	EpiID         int32
	EntregaID     int32
	FuncionarioId int32
	DataInicio    configs.DataBr
	DataFim       configs.DataBr
	Pagina        int32
	Quantidade    int32
}

type EntregaPaginada struct {
	Entradas    []model.EntregaDto `json:"entregas"`
	Total       int64              `json:"total"`
	Pagina      int32              `json:"pagina"`
	PaginaFinal int32              `json:"pagina_final"`
}

func (e *EntregaService) ListaEntregas(ctx context.Context, f FiltroEntregas) (EntregaPaginada, error) {

	limit := f.Quantidade
	if limit <= 0 {
		limit = 1
	}
	paginaAtual := f.Pagina
	if paginaAtual <= 0 {
		paginaAtual = 1
	}

	offset := max((paginaAtual-1)*limit, 0)

	filtro := repository.ListarEntregasParams{
		Limit:         limit,
		Offset:        offset,
		Canceladas:    f.Canceladas,
		IDEntrega:     pgtype.Int4{Int32: f.EntregaID, Valid: f.EntregaID > 0},
		IDFuncionario: pgtype.Int4{Int32: f.FuncionarioId, Valid: f.FuncionarioId > 0},
		DataInicio:    pgtype.Date{Time: f.DataInicio.Time(), Valid: !f.DataInicio.IsZero()},
		DataFim:       pgtype.Date{Time: f.DataFim.Time(), Valid: !f.DataFim.IsZero()},
	}

	entregas, err := e.repo.ListarEntregas(ctx, filtro)
	if err != nil {

		return EntregaPaginada{}, err
	}

	todosTamanhos, err := e.queries.BuscarTodosTamanhosAgrupados(ctx)
	if err != nil {

		return EntregaPaginada{}, err
	}

	tamanhosMap := make(map[int32][]model.TamanhoDto)
	for _, t := range todosTamanhos {

		tamanhosMap[t.Idepi] = append(tamanhosMap[t.Idepi], model.TamanhoDto{
			ID:      int(t.ID),
			Tamanho: t.Tamanho,
		}) 
	}
	todosItens, err := e.queries.BuscarTodosItensEntrega(ctx)
	if err != nil {
		return EntregaPaginada{}, err
	}

	itensMap := make(map[int32][]model.ItemEntregueDto)
	for _, I := range todosItens {

		var valorDecimal decimal.Decimal
		if fVal, err := I.ValorUnitario.Float64Value(); err == nil {
			valorDecimal = decimal.NewFromFloat(fVal.Float64)
		}
		itensMap[I.EntregaID] = append(itensMap[I.EntregaID], model.ItemEntregueDto{
			Id: int64(I.ItemID),
			Epi: model.EpiDto{
				Id:             int(I.EpiID),
				Nome:           I.EpiNome,
				Fabricante:     I.Fabricante,
				CA:             I.Ca,
				Tamanho:        tamanhosMap[I.EpiID],
				Descricao:      I.EpiDesc,
				DataValidadeCa: configs.DataBr(I.ValidadeCa.Time),
				Protecao: model.TipoProtecaoDto{
					ID:   int64(I.TpID),
					Nome: I.TpNome,
				},
			},
			Quantidade:    int(I.Quantidade),
			ValorUnitario: valorDecimal,

		})
	}
	dto := make([]model.EntregaDto, 0, len(entregas))

	for _, entrega := range entregas {

		e := model.EntregaDto{
			Id: int64(entrega.EntregaID),
			Funcionario: model.Funcionario_Dto{
				ID:        int(entrega.FuncID),
				Nome:      entrega.FuncNome,
				Matricula: entrega.Matricula,
				Funcao: model.FuncaoDto{
					ID:     int(entrega.FuncaoID),
					Funcao: entrega.FuncaoNome,
					Departamento: model.DepartamentoDto{
						ID:           int(entrega.DepID),
						Departamento: entrega.DepNome,
					},
				},
			},
			Data_entrega:       configs.DataBr(entrega.DataEntrega.Time),
			Assinatura_Digital: entrega.Assinatura,
			Itens:              itensMap[entrega.EntregaID],
			Id_user: int(entrega.IDUsuarioEntrega.Int32),
		}

		dto = append(dto, e)
	}

	var total int64
	if len(entregas) > 0 {
		total = entregas[0].TotalGeral
	}

	//numero da ultima pagina
	ultimaPagina := int32(math.Ceil(float64(total) / float64(limit)))

	return EntregaPaginada{
		Entradas:    dto,
		Total:       total,
		Pagina:      paginaAtual,
		PaginaFinal: ultimaPagina,
	}, nil
}

func (e *EntregaService) CancelarEntrega(ctx context.Context, id int, iduser int) (int64, error) {

	if id <= 0 {

		return 0, helper.ErrId
	}

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	arg:= repository.CancelarEntregaParams{
		ID: int32(id),
		IDUsuarioEntregaCancelamento: pgtype.Int4{Int32: int32(iduser)},
	}

	qtx := e.queries.WithTx(tx)
	identrega, err := e.repo.Cancelar(ctx, qtx, arg)
	if err != nil {

		return 0, err
	}

	if identrega == 0 {

		return 0, helper.ErrNaoEncontrado
	}

	err = e.repo.CancelarEntregaItem(ctx, qtx, identrega)
	if err != nil {
		return 0, err
	}

	cancelados, err:= e.repo.ListarEpisEntreguesCancelados(ctx, qtx, identrega)
	if err != nil {
		return 0, err
	}

	for _, cancelado := range cancelados {

		args:= repository.ReporEstoqueLoteParams{
			Quantidadeatual: cancelado.Quantidade,
			ID: cancelado.Identrada,
		}
		linhasAfetadas, err:= e.repo.ReporEstoqueEntrada(ctx, qtx,args)
		if err != nil {

			return 0, err
		}

		if linhasAfetadas == 0 {

			return 0, fmt.Errorf("lote de entrada %d não encontrado para reposição", cancelado.Identrada)
		}
	}

	if err:= tx.Commit(ctx); err != nil {

		return 0, err
	}
	return int64(identrega), nil
}
