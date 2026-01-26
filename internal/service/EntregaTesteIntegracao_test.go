package service

import (
	"context"
	"sync"

	"fmt"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestEntrega(t *testing.T) {

	db := SetupTestDB(t)
	var queries *repository.Queries
	defer db.Close()
	ctx := context.Background()

	repo := repository.NewEntregaRepository(db)
	serv := NewEntregaService(repo, db)

	iduser := CreateUser(t, db)
	iddep := CreateDepartamento(t, db)
	IdFuncao := CreateFuncao(t, db, iddep)
	idtam := CreateTamanho(t, db)
	idprotec := CreateProtecao(t, db)
	idepi := CreateEpi(t, db, idprotec)
	idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao)
	idfuncionario2 := CreateFuncionario(t, db, iddep, IdFuncao)
	_ = CreateEntradaEpi(t, db, idfuncionario, idepi, idprotec, idtam, iduser)
	_ = CreateEntradaEpi(t, db, idfuncionario, idepi, idprotec, idtam, iduser)
	entregas := []model.EntregaParaInserir{

		{
			ID_funcionario:     idfuncionario,
			Id_user:            int(iduser),
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "teste.pop",
			Itens: []model.ItemParaInserir{
				{
					ID_epi:         idepi,
					ID_tamanho:     idtam,
					Quantidade:     10,
				
				},
			},
		},
	}

	t.Run("sucesso ao realizar todas as etapas de uma entrega de epi", func(t *testing.T) {

		tx, err := db.Begin(ctx)
		require.NoError(t, err)

		qtx := queries.WithTx(tx)

		args := repository.AddEntregaEpiParams{
			Idfuncionario:    int32(entregas[0].ID_funcionario),
			DataEntrega:      pgtype.Date{Time: entregas[0].Data_entrega.Time(), Valid: true},
			Assinatura:       entregas[0].Assinatura_Digital,
			TokenValidacao:   pgtype.Text{String: "testeToken"},
			IDUsuarioEntrega: pgtype.Int4{Int32: int32(entregas[0].Id_user)},
		}
		identrega, err := repo.AdicionarEntrega(ctx, qtx, args)
		require.NoError(t, err)

		for _, item := range entregas[0].Itens {

			quantidadeNescessaria := item.Quantidade

			lotes := repository.ListarLotesParaConsumoParams{
				Idepi:     int32(item.ID_epi),
				Idtamanho: int32(item.ID_tamanho),
			}
			/*lista todas as entradas com quantidadeAtual maior que 0 e que tenha os idepie e idtamanhos iguais as passado nos parametros*/
			entradaLotes, err := repo.ListarEntregasDisponiveis(ctx, qtx, lotes)
			require.NoError(t, err)

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
					Identrada:     entradaLote.ID,
				}

				_, err = repo.AbaterEstoqueEntrada(ctx, qtx, repository.AbaterEstoqueLoteParams{
					Quantidadeatual: quantidadeAbater,
					ID:              entradaLote.ID,
				})

				_, err := repo.AdicionarEntregaItem(ctx, qtx, itemAdd)
				require.NoError(t, err)
				require.NoError(t, err)

				quantidadeNescessaria -= int(quantidadeAbater)

				var quantidadeAtual int32

				query := `select quantidadeAtual from entrada_epi where id = $1`
				err = tx.QueryRow(ctx, query, entradaLote.ID).Scan(&quantidadeAtual)

				esperado := entradaLote.Quantidadeatual - quantidadeAbater
				require.Equal(t, esperado, quantidadeAtual)
				fmt.Println(quantidadeAtual)

				fmt.Printf("Lote ID: %d | Abatido: %d | Sobrou: %d\n", entradaLote.ID, quantidadeAbater, quantidadeAtual)
			}

		}

		tx.Commit(ctx)
	})

	t.Run("ERRO - tentar entregar mais do que tem no estoque", func(t *testing.T) {

		entregaErro := model.EntregaParaInserir{
			ID_funcionario:     idfuncionario2,
			Id_user:            int(iduser),
			Data_entrega:       *configs.NewDataBrPtr(time.Now()),
			Assinatura_Digital: "teste.pop",
			IdTroca: nil,
			Itens: []model.ItemParaInserir{
				{
					ID_epi:         idepi,
					ID_tamanho:     idtam,
					Quantidade:     300,
					
				},
			},
		}

		err := serv.Salvar(context.Background(), entregaErro)
		require.Error(t, err)
		//require.True(t, errors.Is(err, helper.ErrEstoqueInsuficiente))
		fmt.Println(err)

		var count int
		db.QueryRow(ctx, "SELECT count(*) FROM entrega_epi WHERE IdFuncionario = $1", entregaErro.ID_funcionario).Scan(&count)
		require.Equal(t, 0, count, "A entrega não deveria ter sido salva no banco")

	})

	t.Run("teste de concorrencia (nao deixar 2 usuarios fazer uma entrega do mesmo lote de uma vez)", func(t *testing.T) {

		db := SetupTestDB(t)
		defer db.Close()
		ctx := context.Background()

		repo := repository.NewEntregaRepository(db)
		serv := NewEntregaService(repo, db)

		iduser := CreateUser(t, db)
		iddep := CreateDepartamento(t, db)
		IdFuncao := CreateFuncao(t, db, iddep)
		idtam := CreateTamanho(t, db)
		idprotec := CreateProtecao(t, db)
		idepi := CreateEpi(t, db, idprotec)
		idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao)

		// Criar lote com APENAS 1 unidade
		// Modifique sua função helper ou faça o insert manual aqui

		identrada1 := CreateEntradaEpi1(t, db, idfuncionario, idepi, idprotec,
			idtam, iduser)
		// 2. PREPARAR AS GOROUTINES
		var wg sync.WaitGroup
		numRequisicoes := 2
		wg.Add(numRequisicoes)

		// Canal para capturar os resultados das goroutines
		errs := make(chan error, numRequisicoes)

		entrega := model.EntregaParaInserir{
			ID_funcionario: idfuncionario,
			Id_user:        int(iduser),
			Data_entrega:   configs.DataBr(time.Now()),
			IdTroca: nil,

			Itens: []model.ItemParaInserir{
				{ID_epi: idepi, ID_tamanho: idtam, Quantidade: 1},
			},
		}

		// 3. DISPARAR CONCORRÊNCIA
		for range numRequisicoes {
			go func() {
				defer wg.Done()
				err := serv.Salvar(ctx, entrega)
				if err != nil {
					fmt.Printf("Falha na goroutine: %v\n", err)
				} else {
					fmt.Println("Sucesso na goroutine")
				}
				errs <- err
			}()
		}

		wg.Wait()
		close(errs)

		// 4. VALIDAR RESULTADOS
		sucessos := 0
		falhas := 0

		for err := range errs {
			if err == nil {
				sucessos++
			} else {
				falhas++
				fmt.Printf("Falha esperada detectada: %v\n", err)
			}
		}

		// ASSERTIONS:
		// Em um ambiente concorrente com 1 item, apenas 1 deve conseguir
		require.Equal(t, 1, sucessos, "Apenas uma entrega deveria ter tido sucesso")
		require.Equal(t, 1, falhas, "Uma entrega deveria ter falhado por falta de estoque")

		// Verificar se o estoque não ficou negativo (deve estar em 0)
		var qtdFinal int32
		db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", identrada1).Scan(&qtdFinal)
		fmt.Println(qtdFinal)
		require.Equal(t, int32(0), qtdFinal, "O estoque final deve ser zero e nunca negativo")
	})

	t.Run("testando sucesso ao cancelar uma entrega de epi", func(t *testing.T) {

		db := SetupTestDB(t)
		defer db.Close()
		ctx := context.Background()

		repo := repository.NewEntregaRepository(db)
		serv := NewEntregaService(repo, db)

		iduser := CreateUser(t, db)
		iddep := CreateDepartamento(t, db)
		IdFuncao := CreateFuncao(t, db, iddep)
		idtam := CreateTamanho(t, db)
		idprotec := CreateProtecao(t, db)
		idepi := CreateEpi(t, db, idprotec)
		idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao)

		idEntrada2 := CreateEntradaEpi(t, db, idfuncionario, idepi, idprotec,
			idtam, iduser)

		for i := range 4 {

			err := serv.Salvar(ctx, entregas[0])
			require.NoError(t, err, "A entrega %d deveria ter funcionado", i+1)
		}

		var qtdTotal int64
		db.QueryRow(ctx, "SELECT count(*) from entrega_epi").Scan(&qtdTotal)
		fmt.Printf("Total de entregas no banco: %d\n", qtdTotal)

		var q int64
		query := `SELECT quantidadeAtual FROM entrada_epi WHERE id = $1`
		err := db.QueryRow(ctx, query, idEntrada2).Scan(&q)
		require.NoError(t, err)

		fmt.Printf("Estoque atual do lote antes de cancelar as entregas %d: %d\n", idEntrada2, q)

		for y := range 4 {

			err := serv.CancelarEntrega(ctx, y+1, int(iduser))
			require.NoError(t, err, "o cancelamento %d deveria ter funcionado", y+1)

		}

		var q1 int64
		query = `SELECT quantidadeAtual FROM entrada_epi WHERE id = $1`
		err = db.QueryRow(ctx, query, idEntrada2).Scan(&q1)
		require.NoError(t, err)

		fmt.Printf("Estoque atual do lote depois de cancelar as entregas %d: %d\n", idEntrada2, q1)

	})
}
