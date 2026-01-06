package repositoryIntegracao

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// --- Utilitários ---

// Gera string única para evitar erro de UNIQUE constraint
func RandomString(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// Gera matricula de 7 digitos
func RandomMatricula() string {
	return fmt.Sprintf("%07d", rand.Intn(9999999))
}

// --- Factory Helpers (Nível 0 - Sem dependências) ---

func CreateDepartamento(t *testing.T, db *sql.DB) int {
	var id int
	query := `INSERT INTO departamento (nome) OUTPUT INSERTED.id VALUES (@p1)`
	nome := RandomString("Dep")
	
	err := db.QueryRow(query, nome).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateDepartamento falhou: %v", err)
	}
	return id
}

func CreateTipoProtecao(t *testing.T, db *sql.DB) int {
	var id int
	query := `INSERT INTO tipo_protecao (nome) OUTPUT INSERTED.id VALUES (@p1)`
	nome := RandomString("Tipo")

	err := db.QueryRow(query, nome).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTipoProtecao falhou: %v", err)
	}
	return id
}

func CreateTamanho(t *testing.T, db *sql.DB) int {
	var id int
	query := `INSERT INTO tamanho (tamanho) OUTPUT INSERTED.id VALUES (@p1)`
	// Tamanho tbm é unique, então usamos random
	tamanho := RandomString("T") 

	err := db.QueryRow(query, tamanho).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTamanho falhou: %v", err)
	}
	return id
}

func CreateMotivoDevolucao(t *testing.T, db *sql.DB) int {
	var id int
	query := `INSERT INTO motivo_devolucao (motivo) OUTPUT INSERTED.id VALUES (@p1)`
	motivo := RandomString("Motivo")

	err := db.QueryRow(query, motivo).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateMotivoDevolucao falhou: %v", err)
	}
	return id
}

// --- Factory Helpers (Nível 1 - Dependências Simples) ---

func CreateFuncao(t *testing.T, db *sql.DB, idDepartamento int) int {
	var id int
	query := `INSERT INTO funcao (nome, IdDepartamento) OUTPUT INSERTED.id VALUES (@p1, @p2)`
	nome := RandomString("Func")

	err := db.QueryRow(query, nome, idDepartamento).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateFuncao falhou: %v", err)
	}
	return id
}

func CreateEpi(t *testing.T, db *sql.DB, idTipoProtecao int) int {
	var id int
	query := `
		INSERT INTO epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`
	
	nome := RandomString("EPI")
	ca := RandomString("CA") // CA é unique varchar(20)
	// Corta string se ficar muito grande pro CA
	if len(ca) > 20 {
		ca = ca[:20]
	}

	err := db.QueryRow(query, 
		nome, 
		"Fabricante Teste", 
		ca, 
		"Descrição Teste", 
		time.Now().AddDate(1, 0, 0), // Validade +1 ano
		idTipoProtecao, 
		10, // Alerta minimo
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateEpi falhou: %v", err)
	}
	return id
}

// --- Factory Helpers (Nível 2 - Dependências Compostas) ---

func CreateTamanhosEpis(t *testing.T, db *sql.DB, idEpi, idTamanho int) int {
	var id int
	query := `INSERT INTO tamanhos_epis (IdEpi, IdTamanho) OUTPUT INSERTED.id VALUES (@p1, @p2)`
	
	err := db.QueryRow(query, idEpi, idTamanho).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTamanhosEpis falhou: %v", err)
	}
	return id
}

func CreateFuncionario(t *testing.T, db *sql.DB, idFuncao, idDepartamento int) int {
	var id int
	query := `INSERT INTO funcionario (nome, matricula, IdFuncao, IdDepartamento) OUTPUT INSERTED.id VALUES (@p1, @p2, @p3, @p4)`
	
	nome := RandomString("Funcionario")
	matricula := RandomMatricula() // Unique varchar(7)

	err := db.QueryRow(query, nome, matricula, idFuncao, idDepartamento).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateFuncionario falhou: %v", err)
	}
	return id
}

func CreateEntradaEpi(t *testing.T, db *sql.DB, idEpi, idTamanho int) int {
	var id int
	query := `
		INSERT INTO entrada_epi 
		(IdEpi, IdTamanho, data_entrada, quantidade, data_fabricacao, data_validade, lote, fornecedor, valor_unitario) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`

	err := db.QueryRow(query,
		idEpi,
		idTamanho,
		time.Now(),
		100, // Quantidade
		time.Now().AddDate(-1, 0, 0), // Fabricado ano passado
		time.Now().AddDate(2, 0, 0),  // Validade +2 anos
		"LOTE123",
		"Fornecedor X",
		50.00, // Valor
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateEntradaEpi falhou: %v", err)
	}
	return id
}

// --- Factory Helpers (Nível 3 - Processos de Negócio) ---

func CreateEntregaEpi(t *testing.T, db *sql.DB, idFuncionario int) int {
	var id int
	query := `
		INSERT INTO entrega_epi (IdFuncionario, data_entrega, assinatura) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3)`

	assinaturaFake := []byte("assinatura_mock_binaria")

	err := db.QueryRow(query, idFuncionario, time.Now(), assinaturaFake).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateEntregaEpi falhou: %v", err)
	}
	return id
}

// --- Factory Helpers (Nível 4 - Detalhes e Devoluções) ---

func CreateEpisEntregues(t *testing.T, db *sql.DB, idEntrega, idEntrada, idEpi, idTamanho int) int {
	var id int
	query := `
		INSERT INTO epis_entregues (IdEntrega, IdEntrada, IdEpi, IdTamanho, quantidade, valor_unitario) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`

	err := db.QueryRow(query, idEntrega, idEntrada, idEpi, idTamanho, 1, 50.00).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateEpisEntregues falhou: %v", err)
	}
	return id
}

func CreateDevolucao(t *testing.T, db *sql.DB, idEpi, idFuncionario, idMotivo, idTamanho int) int {
	var id int
	query := `
		INSERT INTO devolucao 
		(IdEpi, IdFuncionario, IdMotivo, data_devolucao, IdTamanho, quantidadeAdevolver, assinatura_digital) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`
	
	assinaturaFake := []byte("assinatura_digital_mock")

	// Nota: Campos opcionais (IdEpiNovo, etc) deixei como NULL padrão
	err := db.QueryRow(query, 
		idEpi, 
		idFuncionario, 
		idMotivo, 
		time.Now(), 
		idTamanho, 
		1, // Quantidade a devolver
		assinaturaFake,
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateDevolucao falhou: %v", err)
	}
	return id
}