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

	// 3. Criação dos Dados Auxiliares (SaaS: O Tenant vem primeiro)
	idEmpresa := CreateEmpresa(t, db) // Novo Helper Mestre

	// Agora passamos idEmpresa para TUDO
	iduser := CreateUser(t, db, idEmpresa)
	iddep := CreateDepartamento(t, db, idEmpresa)
	IdFuncao := CreateFuncao(t, db, iddep, idEmpresa)

	// Tamanhos e Proteções
	idtamAntigo := CreateTamanho(t, db, idEmpresa)
	idtamNovo := CreateTamanho(t, db, idEmpresa)
	idprotec := CreateProtecao(t, db, idEmpresa)

	// EPIs
	idEpiAntigo := CreateEpi(t, db, idprotec, idEmpresa)
	idEpiNovo := CreateEpi(t, db, idprotec, idEmpresa)
	idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao, idEmpresa)

	// Estoque (Entrada)
	// Nota: Passamos iduser (quem criou) e idEmpresa (dono do dado)

	//fornecedores
	Idfornecedor := CreateFornecedor(t, db, idEmpresa)
	// Entrada Antiga (Vai aumentar +1 na devolução)
	idEntradaAntiga := CreateEntradaEpi(t, db, idfuncionario, idEpiAntigo, idprotec, idtamAntigo, iduser, Idfornecedor, idEmpresa)

	// Entrada Nova (Vai diminuir -1, pois é o item da troca)
	idEntradaNova := CreateEntradaEpi(t, db, idfuncionario, idEpiNovo, idprotec, idtamNovo, iduser, Idfornecedor, idEmpresa)

	// Entregas anteriores
	idEntregaAntiga := CreateEntregaEpi(t, db, idfuncionario, iduser, idEmpresa)
	_ = CreateEpiEntregues(t, db, idEntregaAntiga, idEntradaAntiga, idEpiAntigo, idtamAntigo, idEmpresa)

	// Motivos
	_ = CreateMotivoDevolucao(t, db, "Desgaste Natural", idEmpresa)
	_ = CreateMotivoDevolucao(t, db, "Dano", idEmpresa)
	_ = CreateMotivoDevolucao(t, db, "Vencimento", idEmpresa)
	_ = CreateMotivoDevolucao(t, db, "Tamanho Errado", idEmpresa) // ID 4 (Suposição de ordem)
	idMotivoTeste := 4                                            // Ajuste conforme o ID retornado se necessário, ou capture o ID criado acima

	t.Run("Deve realizar uma troca: Devolver item antigo ao estoque e Retirar item novo", func(t *testing.T) {

		// --- FUNÇÃO AUXILIAR DE LOG (DEBUG) ---
		logEstadoBanco := func(momento string) {
			t.Logf("\n====== ESTADO DO BANCO (Tenant: %d): %s ======", idEmpresa, momento)

			// Validamos com TenantID para garantir isolamento
			var qtdEstoqueAntigo int
			err := db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaAntiga, idEmpresa).Scan(&qtdEstoqueAntigo)
			if err == nil {
				t.Logf("[ESTOQUE ANTIGO - ENTROU] ID: %d | Qtd Atual: %d", idEntradaAntiga, qtdEstoqueAntigo)
			}

			var qtdEstoqueNovo int
			err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaNova, idEmpresa).Scan(&qtdEstoqueNovo)
			if err == nil {
				t.Logf("[ESTOQUE NOVO - SAIU]     ID: %d | Qtd Atual: %d", idEntradaNova, qtdEstoqueNovo)
			}
			t.Log("==========================================\n")
		}

		// --- ARRANGE ---
		qtdDevolver := 1
		qtdNova := 1

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
			// Dica: Se seu Service precisar validar o Tenant, injete no Context ou na Struct aqui
		}

		// Captura estado INICIAL
		var qtdAntigoAntes, qtdNovoAntes int
		err := db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaAntiga, idEmpresa).Scan(&qtdAntigoAntes)
		require.NoError(t, err)
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaNova, idEmpresa).Scan(&qtdNovoAntes)
		require.NoError(t, err)

		logEstadoBanco("ANTES DA EXECUÇÃO")

		// --- ACT ---
		// Assumindo que o service extrai o tenant ou usa o IdUser para validar
		err = servDevolucao.SalvarDevolucao(ctx, dadosDevolucao, int32(idEmpresa))

		logEstadoBanco("DEPOIS DA EXECUÇÃO")

		// --- ASSERT ---
		require.NoError(t, err, "A função SalvarDevolucao retornou erro: %v", err)

		// Validações Finais
		var qtdAntigoDepois, qtdNovoDepois int
		_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaAntiga, idEmpresa).Scan(&qtdAntigoDepois)
		_ = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaNova, idEmpresa).Scan(&qtdNovoDepois)

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

	// 2. Helpers (Cenário SaaS Completo)
	idEmpresa := CreateEmpresa(t, db) // Tenant Isolation

	iduser := CreateUser(t, db, idEmpresa)
	iddep := CreateDepartamento(t, db, idEmpresa)
	IdFuncao := CreateFuncao(t, db, iddep, idEmpresa)
	idprotec := CreateProtecao(t, db, idEmpresa)
	idfuncionario := CreateFuncionario(t, db, iddep, IdFuncao, idEmpresa)

	// EPIs e Tamanhos
	idEpiVelho := CreateEpi(t, db, idprotec, idEmpresa)
	idTamVelho := CreateTamanho(t, db, idEmpresa)

	idEpiNovo := CreateEpi(t, db, idprotec, idEmpresa)
	idTamNovo := CreateTamanho(t, db, idEmpresa)

	//fornecedores
	Idfornecedor := CreateFornecedor(t, db, idEmpresa)
	// 3. Estoque (Entradas)
	// Item Velho (já está com funcionário)
	idEntradaVelha := CreateEntradaEpi(t, db, idfuncionario, idEpiVelho, idprotec, idTamVelho, iduser, Idfornecedor,idEmpresa)

	// Item Novo (será a troca)
	idEntradaNova := CreateEntradaEpi(t, db, idfuncionario, idEpiNovo, idprotec, idTamNovo, iduser, Idfornecedor,idEmpresa)

	// Simula entrega anterior
	idEntregaAntiga := CreateEntregaEpi(t, db, idfuncionario, iduser, idEmpresa)
	_ = CreateEpiEntregues(t, db, idEntregaAntiga, idEntradaVelha, idEpiVelho, idTamVelho, idEmpresa)

	// Motivos
	_ = CreateMotivoDevolucao(t, db, "Desgaste Natural", idEmpresa)
	idmotivo2 := CreateMotivoDevolucao(t, db, "Dano", idEmpresa)

	// 4. PREPARAÇÃO: Executar uma Troca REAL primeiro
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
		IdMotivo:            int(idmotivo2), // Dano
		DataDevolucao:       *configs.NewDataBrPtr(time.Now()),
		IdTamanho:           int(idTamVelho),
		QuantidadeADevolver: 1,
		AssinaturaDigital:   "assinatura_teste",
		IdUser:              int(iduser),
	}

	err := servDevolucao.SalvarDevolucao(ctx, dadosTroca, int32(idEmpresa))
	require.NoError(t, err, "Falha ao criar o cenário de troca inicial")

	// Descobrir o ID da Troca (Filtrando pelo Tenant para segurança)
	var idTrocaCriada int
	err = db.QueryRow(ctx, "SELECT MAX(id) FROM devolucao WHERE tenant_id = $1", idEmpresa).Scan(&idTrocaCriada)
	require.NoError(t, err, "Não foi possível recuperar o ID da troca criada")

	t.Run("Deve Cancelar uma Devolucao/Troca e Repor o Estoque do Item Novo", func(t *testing.T) {

		// Captura estado ANTES de cancelar
		var qtdAntesCancelamento int
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaNova, idEmpresa).Scan(&qtdAntesCancelamento)
		require.NoError(t, err)

		t.Logf("Estoque Antes Cancelar (Item Novo): %d", qtdAntesCancelamento)

		// --- ACT (AÇÃO) ---
		// Passamos o ID da troca e o ID do usuário que está cancelando
		err = servDevolucao.CancelarDevolucao(ctx, idTrocaCriada, int(iduser), int(idEmpresa))

		// --- ASSERT ---
		require.NoError(t, err, "A função CancelarDevolucao retornou erro: %v", err)

		// Validação Final
		var qtdDepoisCancelamento int
		err = db.QueryRow(ctx, "SELECT quantidadeAtual FROM entrada_epi WHERE id = $1 AND tenant_id = $2", idEntradaNova, idEmpresa).Scan(&qtdDepoisCancelamento)
		require.NoError(t, err)

		// Lógica: Tinha 99 (pós troca), cancelei, deve voltar para 100.
		require.Equal(t, qtdAntesCancelamento+1, qtdDepoisCancelamento,
			"O estoque deveria ter sido reposto. Antes: %d, Depois: %d",
			qtdAntesCancelamento, qtdDepoisCancelamento)
	})
}
