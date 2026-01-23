package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type EntradaRepository interface {
	Adicionar(ctx context.Context, args repository.AddEntradaEpiParams) error
	ListarEntradas(ctx context.Context, args repository.ListarEntradasParams) ([]repository.ListarEntradasRow, error)
	CancelarEntrada(ctx context.Context, args repository.CancelarEntradaParams) (int64, error)
	TotalEntradas(ctx context.Context, args repository.ContarEntradasParams) (int64, error)
}

type EntradaService struct {
	repo EntradaRepository
}

func NewEntradaService(e EntradaRepository) *EntradaService {

	return &EntradaService{repo: e}
}

func (e *EntradaService) Adicionar(ctx context.Context, model model.EntradaEpiInserir) error {

	//data de entrada menor que a atual
	hoje := time.Now().Truncate(24 * time.Hour)
	if model.Data_entrada.Time().Truncate(24 * time.Hour).Before(hoje) {

		return helper.ErrDataMenor
	}
	//data de validade igual a de fabricacao
	if model.DataValidade.Time().Equal(model.DataFabricacao.Time()) {

		return helper.ErrDataIgual
	}
	//data de validade menor a de fabricacao
	if model.DataValidade.Time().Before(model.DataFabricacao.Time()) {
		return helper.ErrDataMenorValidade
	}

	model.Lote = strings.TrimSpace(model.Lote)
	model.Fornecedor = strings.TrimSpace(model.Fornecedor)
	model.Nota_fiscal_numero = strings.TrimSpace(model.Nota_fiscal_numero)
	model.Nota_fiscal_serie = strings.TrimSpace(model.Nota_fiscal_serie)

	var vm pgtype.Numeric
	err := vm.Scan(model.ValorUnitario.String())
	if err != nil {
		return err
	}

	err = e.repo.Adicionar(ctx, repository.AddEntradaEpiParams{
		Idepi:            int32(model.ID_epi),
		Idtamanho:        int32(model.Id_tamanho),
		DataEntrada:      pgtype.Date{Time: model.Data_entrada.Time(), Valid: true},
		Quantidade:       int32(model.Quantidade),
		Quantidadeatual:  int32(model.Quantidade_Atual),
		DataFabricacao:   pgtype.Date{Time: model.DataFabricacao.Time(), Valid: true},
		DataValidade:     pgtype.Date{Time: model.DataValidade.Time(), Valid: true},
		Fornecedor:       model.Fornecedor,
		Lote:             model.Lote,
		ValorUnitario:    vm,
		NotaFiscalNumero: model.Nota_fiscal_numero,
		NotaFiscalSerie:  pgtype.Text{String: model.Nota_fiscal_serie},
		IDUsuarioCriacao: pgtype.Int4{Int32: int32(model.Id_user), Valid: true},
	})

	return nil
}

type FiltroEntradas struct {
	Canceladas bool
	EpiID      int32
	EntradaID  int32
	DataInicio configs.DataBr
	DataFim    configs.DataBr
	NotaFiscal string
	Pagina     int32
	Quantidade int32
}

type EntradaPaginada struct {
	Entradas    []model.EntradaEpiDto `json:"entradas"`
	Total       int64                 `json:"total"`
	Pagina      int32                 `json:"pagina"`
	PaginaFinal int32                 `json:"pagina_final"`
}

func (e *EntradaService) ListarEntradas(ctx context.Context, f FiltroEntradas) (EntradaPaginada, error) {

	limit := f.Quantidade
	if limit <= 0 {
		limit = 1
	}
	paginaAtual := f.Pagina
	if paginaAtual <= 0 {
		paginaAtual = 1
	}
	offset := max((paginaAtual-1)*limit, 0)

	filtro := repository.ListarEntradasParams{
		Canceladas: f.Canceladas,
		IDEpi:      pgtype.Int4{Int32: f.EpiID, Valid: f.EpiID > 0},
		IDEntrada:  pgtype.Int4{Int32: f.EntradaID, Valid: f.EntradaID > 0},
		DataInicio: pgtype.Date{Time: f.DataInicio.Time(), Valid: !f.DataInicio.IsZero()},
		DataFim:    pgtype.Date{Time: f.DataFim.Time(), Valid: !f.DataFim.IsZero()},
		NotaFiscal: pgtype.Text{String: f.NotaFiscal, Valid: f.NotaFiscal != ""},
		Limit:      limit,
		Offset:     offset,
	}

	entradas, err := e.repo.ListarEntradas(ctx, filtro)
	if err != nil {

		return EntradaPaginada{}, err
	}

	dto := make([]model.EntradaEpiDto, 0, len(entradas))

	for _, entrada := range entradas {

		var valorDecimal decimal.Decimal
		if fVal, err := entrada.ValorUnitario.Float64Value(); err == nil {
			valorDecimal = decimal.NewFromFloat(fVal.Float64)
		}
		var idUsuario int
		if entrada.IDUsuarioCriacao.Valid {
			idUsuario = int(entrada.IDUsuarioCriacao.Int32)
		} else {
			idUsuario = 0 // ou algum valor padrão
		}
		e := model.EntradaEpiDto{
			ID: int(entrada.ID),
			Epi: model.EpiDto{
				Id:         int(entrada.Idepi),
				Nome:       entrada.EpiNome,
				Fabricante: entrada.Fabricante,
				CA:         entrada.Ca,
				Tamanho: []model.TamanhoDto{
					{
						ID:      int(entrada.Idtamanho),
						Tamanho: entrada.TamanhoNome,
					},
				},
				Descricao:      entrada.EpiDescricao,
				DataValidadeCa: configs.DataBr(entrada.ValidadeCa.Time),
				Protecao: model.TipoProtecaoDto{
					ID:   int64(entrada.Idtipoprotecao),
					Nome: entrada.ProtecaoNome,
				},
			},
			Data_entrada:       *configs.NewDataBrPtr(entrada.DataEntrada.Time),
			Quantidade:         int(entrada.Quantidade),
			Quantidade_Atual:   int(entrada.Quantidadeatual),
			Lote:               entrada.Lote,
			Fornecedor:         entrada.Fornecedor,
			Nota_fiscal_serie:  entrada.NotaFiscalSerie.String,
			Nota_fiscal_numero: entrada.NotaFiscalNumero,
			ValorUnitario:      valorDecimal,
			Id_user:            idUsuario,
		}

		dto = append(dto, e)
	}

	total, err := e.repo.TotalEntradas(ctx, repository.ContarEntradasParams{
		Canceladas: filtro.Canceladas,
		IDEpi:      filtro.IDEpi,
		IDEntrada:  filtro.IDEntrada,
		DataInicio: filtro.DataInicio,
		DataFim:    filtro.DataFim,
		NotaFiscal: filtro.NotaFiscal,
	})
	if err != nil {
		return EntradaPaginada{}, err
	}

	paginaFinal := int32(math.Ceil(float64(total) / float64(limit)))

	return EntradaPaginada{
		Entradas:    dto,
		Total:       total,
		Pagina:      paginaAtual,
		PaginaFinal: paginaFinal,
	}, nil
}

func (e *EntradaService) CancelarEntrada(ctx context.Context, id int, idUser int) (int64, error) {

	if id <= 0 {

		return 0, helper.ErrId
	}

	arg:= repository.CancelarEntradaParams{
		ID: int32(id),
		IDUsuarioCriacaoCancelamento: pgtype.Int4{Int32: int32(idUser)},
	}
	linhasAfetadas, err := e.repo.CancelarEntrada(ctx, arg)
	if err != nil {

		return 0, fmt.Errorf("erro técnico ao cancelar: %w", err)
	}

	if linhasAfetadas == 0 {

		return 0, helper.ErrNaoEncontrado
	}

	return linhasAfetadas, nil
}
