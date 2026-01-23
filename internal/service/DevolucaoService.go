package service

import (
	"context"
	"fmt"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevolucaoRepository interface {
	AdicionarDevolucao(ctx context.Context, qtx *repository.Queries, args repository.AddDevolucaoSimplesParams)
	AdicionarTroca(ctx context.Context, qtx *repository.Queries, arg repository.AddTrocaEpiParams) (int32, error)
	EntregaVinculada(ctx context.Context, qtx *repository.Queries, arg repository.AddEntregaVinculadaParams) (int32, error)
	Cancelar(ctx context.Context, qtx *repository.Queries, arg repository.CancelarDevolucaoParams) (int64, error)
	Listar(ctx context.Context, qtx *repository.Queries, args repository.ListarDevolucoesParams) ([]repository.ListarDevolucoesRow, error)
}

type DevolucaoService struct {
	repo        DevolucaoRepository
	db          *pgxpool.Pool
	queries     *repository.Queries
	repoEntrega EntregaService
}

func NewDevolucaoService(d DevolucaoRepository, db *pgxpool.Pool, repoEntrega EntregaService) *DevolucaoService {

	return &DevolucaoService{

		repo:    d,
		db:      db,
		queries: repository.New(db),
	}
}

func (d *DevolucaoService) SalvarDevolucao(ctx context.Context, modelDevolucao model.DevolucaoInserir) error {
	/*caso tenha um desses 3 motivos(vencimento ca, desgate natural ou danos)
	apenas faz a troca, ao apenas a devolucao, mas no epi trocado não volta pro estoque*/
	if modelDevolucao.IdMotivo == 1 || modelDevolucao.IdMotivo == 2 || modelDevolucao.IdMotivo == 3 {

	
		tx, err := d.db.Begin(ctx)
		if err != nil {
			return err
		}

		defer tx.Rollback(ctx)
		qtx := d.queries.WithTx(tx)

		/*verificando se os ponteiros sao nulos*/
		var epiNovo, tamanhoNovo, quantidadeNova pgtype.Int4

		if modelDevolucao.IdEpiNovo != nil {
			epiNovo = pgtype.Int4{Int32: int32(*modelDevolucao.IdEpiNovo), Valid: true}
		}
		if modelDevolucao.IdTamanhoNovo != nil {
			tamanhoNovo = pgtype.Int4{Int32: int32(*modelDevolucao.IdTamanhoNovo), Valid: true}
		}
		if modelDevolucao.NovaQuantidade != nil {
			quantidadeNova = pgtype.Int4{Int32: int32(*modelDevolucao.NovaQuantidade), Valid: true}
		}

		arg := repository.AddTrocaEpiParams{

			Idfuncionario:         int32(modelDevolucao.IdFuncionario),
			Idepi:                 int32(modelDevolucao.IdEpi),
			Idmotivo:              int32(modelDevolucao.IdMotivo),
			DataDevolucao:         pgtype.Date{Time: time.Time(modelDevolucao.DataDevolucao), Valid: true},
			Idtamanho:             int32(modelDevolucao.IdTamanho),
			Quantidadeadevolver:   int32(modelDevolucao.QuantidadeADevolver),
			Idepinovo:             epiNovo,
			Idtamanhonovo:         tamanhoNovo,
			Quantidadenova:        quantidadeNova,
			AssinaturaDigital:     modelDevolucao.AssinaturaDigital,
			IDUsuarioCancelamento: pgtype.Int4{Int32: int32(modelDevolucao.IdUser), Valid: true},
		}

		iddevolucao, err := d.repo.AdicionarTroca(ctx, qtx, arg)
		if err != nil {
			return err
		}

		//caso a variavel troca seja marcada como verdadeira, a devolucao do epi, vem junto com a entrada
		//de outro epi
		if modelDevolucao.Troca {

			if modelDevolucao.IdEpiNovo == nil || modelDevolucao.IdTamanhoNovo == nil || modelDevolucao.NovaQuantidade == nil {

				return fmt.Errorf("dados do epi povo sao necessario para a troca")
			}

			 iddevolucaoConvertido := int(iddevolucao)

			 modelEntrega:= model.EntregaParaInserir{

				ID_funcionario: int64(arg.Idfuncionario),
				Id_user: modelDevolucao.IdUser,
				Data_entrega: modelDevolucao.DataDevolucao,
				IdTroca: &iddevolucaoConvertido,
				Assinatura_Digital: arg.AssinaturaDigital,
				Itens: []model.ItemParaInserir{

					{
						ID_epi: int64(*modelDevolucao.IdEpiNovo),
						ID_tamanho: int64(*modelDevolucao.IdTamanhoNovo),
						Quantidade: *modelDevolucao.NovaQuantidade,
					},
				},
			 }

			 /*chamando o service de entregas pos ja tem toda a logica necessario para concluir uma entrega com sucesso*/
			 err = d.repoEntrega.RegistrarEntrega(ctx, qtx, modelEntrega)
			 if err != nil {
				return err
			 }
		}
		return tx.Commit(ctx)
	}
	return nil
}
