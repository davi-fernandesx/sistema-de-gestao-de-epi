package departamento

import (
	"context"
	"testing"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/repositoryIntegracao"

	"github.com/stretchr/testify/assert"
	// Importe seu pacote de models e repository aqui
)

func TestAddDepartamento_Integration(t *testing.T) {
	// 1. SETUP: Pega a conexão do Testcontainers (ou banco local)
	// Supondo que você tem uma func 'setupTestDB' que retorna *sql.DB limpo
	db := setupTestDB(t) 
	defer db.Close()

	// Inicializa o Repositório
	repo := NewDepartamentoRepository(db)

	// Usamos sub-testes (t.Run) para organizar os cenários
	
	// --- CENÁRIO 1: Inserção com Sucesso ---
	t.Run("Deve salvar um novo departamento com sucesso", func(t *testing.T) {
		// Arrange
		nomeDep := "Financeiro_" + randomString("Unique") // Garante nome único
		input := &model.Departamento{
			Departamento: nomeDep,
		}

		// Act
		err := repo.AddDepartamento(context.Background(), input)

		// Assert
		assert.NoError(t, err)

		// Verificação Prova Real (Vai no banco ver se gravou mesmo)
		// Nota: Ajuste o nome da coluna 'departamento' ou 'nome' conforme sua tabela real
		var idGerado int
		queryValidacao := "SELECT id FROM departamento WHERE departamento = @p1"
		err = db.QueryRow(queryValidacao, nomeDep).Scan(&idGerado)

		assert.NoError(t, err, "Deveria encontrar o registro no banco")
		assert.Greater(t, idGerado, 0, "O ID deveria ter sido gerado")
	})

	// --- CENÁRIO 2: Erro de Duplicidade (Unique Violation) ---
	t.Run("Deve retornar erro formatado ao tentar duplicar departamento", func(t *testing.T) {
		// Arrange
		nomeDuplicado := "RH_" + randomString("Dup")
		input := &model.Departamento{
			Departamento: nomeDuplicado,
		}

		// Passo 1: Insere a primeira vez (deve funcionar)
		err := repo.AddDepartamento(context.Background(), input)
		assert.NoError(t, err)

		// Passo 2: Tenta inserir DE NOVO o mesmo objeto
		err = repo.AddDepartamento(context.Background(), input)

		// Assert
		assert.Error(t, err) // Tem que dar erro
		
		// Verifica se a sua mensagem customizada está lá
		// Você escreveu: "departamento %s ja existe no sistema"
		assert.Contains(t, err.Error(), "ja existe no sistema")
		
		// Opcional: Verifica se o erro original (Wrapped) está lá
		assert.ErrorIs(t, err, Errors.ErrSalvar)
	})
}