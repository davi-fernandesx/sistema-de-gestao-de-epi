package integracao

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func randomString(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func CreateDepartamento(t *testing.T, db *sql.DB) int {
	var id int
	// SQL Server usa OUTPUT INSERTED.id para retornar o ID gerado
	query := `INSERT INTO departamento (nome) OUTPUT INSERTED.id VALUES (@p1)`

	err := db.QueryRow(query, randomString("Dep")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateDepartamento falhou: %v", err)
	}
	return id
}

func CreateFuncao(t *testing.T, db *sql.DB, idDepartamento int) int {
	var id int
	query := `INSERT INTO funcao (nome, IdDepartamento) OUTPUT INSERTED.id VALUES (@p1, @p2)`

	err := db.QueryRow(query, randomString("Func"), idDepartamento).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateFuncao falhou: %v", err)
	}
	return id
}

func CreateProtecao(t *testing.T, db *sql.DB) int {

	var id int

	query := `insert into tipo_protecao(nome) OUTPUT INSERTED.id values (@p1)`
	err := db.QueryRow(query, randomString("protec")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateProtecao falhou: %v", err)
	}

	return id

}

func CreateTamanho(t *testing.T, db *sql.DB) int {

	var id int
	query := `insert into tamanho(tamanho) OUTPUT INSERTED.id values (@p1)`
	err := db.QueryRow(query, randomString("tam")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTamanho falhou: %v", err)
	}
	return id

}

func CreateEpi(t *testing.T, db *sql.DB, idTipoProtecao int) int {
	var id int

	// Usamos @p1, @p2... para o driver mapear automaticamente os argumentos na ordem
	query := `
		INSERT INTO epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
	`
	// 1. Gera seed baseada no tempo para garantir aleatoriedade a cada execução
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// Gerando dados aleatórios para garantir unicidade
	nome := randomString("Luva")
	fabricante := randomString("Fabr")
	ca := fmt.Sprintf("%d", r.Intn(9999999))
	descricao := "Descrição de teste gerada automaticamente"
	validade := time.Now().AddDate(1, 0, 0) // Validade para daqui a 1 ano
	alertaMinimo := 10

	// A ordem aqui deve bater com @p1, @p2, etc.
	err := db.QueryRow(query,
		nome,
		fabricante,
		ca,
		descricao,
		validade,
		idTipoProtecao, // FK que recebemos como argumento
		alertaMinimo,
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateEpi falhou: %v", err)
	}

	return id
}
