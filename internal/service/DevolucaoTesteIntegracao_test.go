package service

import (
	"context"

	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/stretchr/testify/require"
)


func TestSalvarDevolucao(t *testing.T) {
	// 1. Setup do Banco de Teste
	db := SetupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// 2. Inicialização de Repositories e Services
	repo := repository.NewDevolucaoRepository(db)
	repoEntregaImpl := repository.NewEntregaRepository(db)
	servEntrega := NewEntregaService(repoEntregaImpl, db)
	servDevolucao := NewDevolucaoService(repo, db, *servEntrega)

	// 3. Criação dos Dados Auxiliares
	iduser := CreateUser(t, db)
	iddep := CreateDepartamento(t, db)
	IdFuncao := CreateFuncao(t, db, iddep)

	// Tamanhos e Proteções
	idtamAntigo := CreateTamanho(t, db)
	idtamNovo := CreateTamanho(t, db)
	idprotec := CreateProtecao(t, db)

	// EPIs
	idEpiAntigo := CreateEpi(t, db, idprotec)
	idEpiNovo := CreateEpi(t, db, idprotec)
	idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao)

	// Estoque (Entrada) 
	// Entrada Antiga (Vai aumentar +1)
	idEntradaAntiga := CreateEntradaEpi(t, db, idfuncionario, idEpiAntigo, idprotec, idtamAntigo, iduser)
	
	// Entrada Nova (Vai diminuir -1, pois é o item da troca)
	idEntradaNova := CreateEntradaEpi(t, db, idfuncionario, idEpiNovo, idprotec, idtamNovo, iduser)

	// Entregas anteriores
	idEntregaAntiga := CreateEntregaEpi(t, db, idfuncionario, iduser)
	_ = CreateEpiEntregues(t, db, idEntregaAntiga, idEntradaAntiga, idEpiAntigo, idtamAntigo)

	// Motivos
	_ = CreateMotivoDevolucao(t, db, "Desgaste Natural")
	_ = CreateMotivoDevolucao(t, db, "Dano")
	_ = CreateMotivoDevolucao(t, db, "Vencimento")
	_ = CreateMotivoDevolucao(t, db, "Tamanho Errado") // ID 4

	t.Run("Deve realizar uma troca: Devolver item antigo ao estoque e Retirar item novo", func(t *testing.T) {
		
		// --- FUNÇÃO AUXILIAR DE LOG (DEBUG) ---
		logEstadoBanco := func(momento string) {
			t.Logf("\n====== ESTADO DO BANCO: %s ======", momento)
			
			// 1. Checar Estoque do Item ANTIGO (Tem que SUBIR)
			var qtdEstoqueAntigo int
			_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaAntiga).Scan(&qtdEstoqueAntigo)
			t.Logf("[ESTOQUE ANTIGO - ENTROU] ID: %d | Qtd Atual: %d", idEntradaAntiga, qtdEstoqueAntigo)

			// 2. Checar Estoque do Item NOVO (Tem que DESCER)
			var qtdEstoqueNovo int
			_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdEstoqueNovo)
			t.Logf("[ESTOQUE NOVO - SAIU]     ID: %d | Qtd Atual: %d", idEntradaNova, qtdEstoqueNovo)

			t.Log("==========================================\n")
		}

		// --- ARRANGE ---
		qtdDevolver := 1
		qtdNova := 1
		idMotivoTeste := 4 // Tamanho Errado (Devolve ao estoque)

		// Ponteiros
		idEpiNovoInt := int(idEpiNovo)
		idTamNovoInt := int(idtamNovo)
		qtdNovaInt := qtdNova

		dadosDevolucao := model.DevolucaoInserir{
			Troca:               true,
			IdEpiNovo:           &idEpiNovoInt,
			IdTamanhoNovo:       &idTamNovoInt,
			NovaQuantidade:      &qtdNovaInt,
			IdFuncionario:       int(idfuncionario),
			IdEpi:               int(idEpiAntigo),
			IdMotivo:            idMotivoTeste,
			DataDevolucao:       *configs.NewDataBrPtr(time.Now()),
			IdTamanho:           int(idtamAntigo),
			QuantidadeADevolver: qtdDevolver,
			AssinaturaDigital:   "assinatura_base64_teste",
			IdUser:              int(iduser),
		}

		// Captura estado INICIAL para validação matemática
		var qtdAntigoAntes, qtdNovoAntes int
		err := db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaAntiga).Scan(&qtdAntigoAntes)
		require.NoError(t, err)
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdNovoAntes)
		require.NoError(t, err)

		// LOG DO ESTADO ANTES
		logEstadoBanco("ANTES DA EXECUÇÃO")

		// --- ACT ---
		err = servDevolucao.SalvarDevolucao(ctx, dadosDevolucao)

		// LOG DO ESTADO DEPOIS
		logEstadoBanco("DEPOIS DA EXECUÇÃO")

		// --- ASSERT ---
		require.NoError(t, err, "A função SalvarDevolucao retornou erro: %v", err)

		// Validações Finais
		var qtdAntigoDepois, qtdNovoDepois int
		_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaAntiga).Scan(&qtdAntigoDepois)
		_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdNovoDepois)
		
		// O antigo deve ter AUMENTADO
		require.Equal(t, qtdAntigoAntes+qtdDevolver, qtdAntigoDepois, 
			"ERRO: O estoque do item devolvido deveria ter aumentado.")

		// O novo deve ter DIMINUÍDO
		require.Equal(t, qtdNovoAntes-qtdNova, qtdNovoDepois, 
			"ERRO: O estoque do item novo (troca) deveria ter diminuído.")
	})
}



func TestCancelarDevolucao(t *testing.T) {
	// 1. Setup do Banco e Services
	db := SetupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	repo := repository.NewDevolucaoRepository(db)
	repoEntregaImpl := repository.NewEntregaRepository(db)
	servEntrega := NewEntregaService(repoEntregaImpl, db)
	servDevolucao := NewDevolucaoService(repo, db, *servEntrega)

	// 2. Helpers (Criando o cenário)
	iduser := CreateUser(t, db)
	iddep := CreateDepartamento(t, db)
	IdFuncao := CreateFuncao(t, db, iddep)
	idprotec := CreateProtecao(t, db)
	idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao)

	// EPIs e Tamanhos
	idEpiVelho := CreateEpi(t, db, idprotec)
	idTamVelho := CreateTamanho(t, db)
	
	idEpiNovo := CreateEpi(t, db, idprotec)
	idTamNovo := CreateTamanho(t, db)

	// 3. Estoque (Entradas)
	// Item Velho (que o funcionário já tem)
	idEntradaVelha := CreateEntradaEpi(t, db, idfuncionario, idEpiVelho, idprotec, idTamVelho, iduser)
	
	// Item Novo (que será entregue na troca e depois reposto no cancelamento)
	// Vamos fixar uma quantidade alta para facilitar a conta (ex: 100)
	// NOTA: Certifique-se que seu helper CreateEntradaEpi define uma quantidade conhecida
	idEntradaNova := CreateEntradaEpi(t, db, idfuncionario, idEpiNovo, idprotec, idTamNovo, iduser)

	// Simula que o funcionário já tinha o item velho
	idEntregaAntiga := CreateEntregaEpi(t, db, idfuncionario, iduser)
	_ = CreateEpiEntregues(t, db, idEntregaAntiga, idEntradaVelha, idEpiVelho, idTamVelho)

	// Motivo (Dano - Gera troca)// 
	_ = CreateMotivoDevolucao(t, db, "Desgaste Natural")
	idmotivo2:= CreateMotivoDevolucao(t, db, "Dano")//ID 2 (exemplo)
	// 4. PREPARAÇÃO: Executar uma Troca REAL primeiro para ter o que cancelar
	t.Log("--- PREPARANDO CENÁRIO: CRIANDO UMA TROCA ---")
	
	idEpiNovoInt := int(idEpiNovo)
	idTamNovoInt := int(idTamNovo)
	qtd := 1
	
	dadosTroca := model.DevolucaoInserir{
		Troca:               true,
		IdEpiNovo:           &idEpiNovoInt,
		IdTamanhoNovo:       &idTamNovoInt,
		NovaQuantidade:      &qtd,
		IdFuncionario:       int(idfuncionario),
		IdEpi:               int(idEpiVelho),
		IdMotivo:            (int(idmotivo2)), // Dano
		DataDevolucao:       *configs.NewDataBrPtr(time.Now()),
		IdTamanho:           int(idTamVelho),
		QuantidadeADevolver: 1,
		AssinaturaDigital:   "assinatura_teste",
		IdUser:              int(iduser),
	}

	err := servDevolucao.SalvarDevolucao(ctx, dadosTroca)
	require.NoError(t, err, "Falha ao criar o cenário de troca inicial")

	// Descobrir o ID da Troca que acabou de ser criada (Select Max ID)
	var idTrocaCriada int
	// Ajuste 'troca_items' para o nome real da sua tabela de devolução/troca
	err = db.QueryRow(ctx, "SELECT MAX(id) FROM devolucao ").Scan(&idTrocaCriada)
	require.NoError(t, err, "Não foi possível recuperar o ID da troca criada")

	t.Run("Deve Cancelar uma Devolucao/Troca e Repor o Estoque do Item Novo", func(t *testing.T) {

		// --- FUNÇÃO AUXILIAR DE LOG (DEBUG) ---
		logEstadoBanco := func(momento string) {
			t.Logf("\n====== ESTADO DO BANCO: %s ======", momento)

			// Monitoramos o estoque da ENTRADA NOVA (o item que foi dado e agora deve voltar)
			var qtdAtual int
			err := db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdAtual)
			if err == nil {
				t.Logf("[ESTOQUE DO ITEM NOVO] ID Entrada: %d | Qtd Atual: %d", idEntradaNova, qtdAtual)
			} else {
				t.Logf("[ERRO] Ao ler estoque: %v", err)
			}
			
			// Verificar se a entrega está cancelada (opcional, depende da sua tabela)
			
			// Exemplo: SELECT cancelado FROM entregas WHERE id_troca = ...
			// Ajuste a query conforme sua estrutura
			
			t.Log("==========================================\n")
		}

		// Captura estado ANTES de cancelar (mas DEPOIS de ter feito a troca)
		var qtdAntesCancelamento int
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdAntesCancelamento)
		require.NoError(t, err)

		logEstadoBanco("ANTES DE CANCELAR (Já com a troca feita)")

		// --- ACT (AÇÃO) ---
		err = servDevolucao.CancelarDevolucao(ctx, idTrocaCriada, int(iduser))

		// --- ASSERT ---
		require.NoError(t, err, "A função CancelarDevolucao retornou erro: %v", err)

		logEstadoBanco("DEPOIS DE CANCELAR")

		// Validação Final
		var qtdDepoisCancelamento int
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1", idEntradaNova).Scan(&qtdDepoisCancelamento)
		require.NoError(t, err)

		// Lógica:
		// Se eu tinha 100, fiz a troca (-1), fiquei com 99.
		// Ao cancelar, o item volta (+1), volto para 100.
		// Logo: Depois > Antes
		
		require.Equal(t, qtdAntesCancelamento + 1, qtdDepoisCancelamento, 
			"O estoque deveria ter sido reposto. Antes Cancelar: %d, Depois Cancelar: %d", 
			qtdAntesCancelamento, qtdDepoisCancelamento)
	})
}