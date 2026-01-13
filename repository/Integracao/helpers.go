package integracao

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/shopspring/decimal"
)

func randomString(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func randomInt() *rand.Rand{

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	return r
}


func CreateDepartamento(t *testing.T, db *sql.DB) int64 {
	var id int64
	// SQL Server usa OUTPUT INSERTED.id para retornar o ID gerado
	query := `INSERT INTO departamento (nome) OUTPUT INSERTED.id VALUES (@p1)`

	err := db.QueryRow(query, randomString("Dep")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateDepartamento falhou: %v", err)
	}
	return id
}

func CreateFuncao(t *testing.T, db *sql.DB, idDepartamento int64) int64 {
	var id int64
	query := `INSERT INTO funcao (nome, IdDepartamento) OUTPUT INSERTED.id VALUES (@p1, @p2)`

	err := db.QueryRow(query, randomString("Func"), idDepartamento).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateFuncao falhou: %v", err)
	}
	return id
}

func CreateProtecao(t *testing.T, db *sql.DB) int64 {

	var id int64

	query := `insert into tipo_protecao(nome) OUTPUT INSERTED.id values (@p1)`
	err := db.QueryRow(query, randomString("protec")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateProtecao falhou: %v", err)
	}

	return id

}

func CreateTamanho(t *testing.T, db *sql.DB) int64 {

	var id int64
	query := `insert into tamanho(tamanho) OUTPUT INSERTED.id values (@p1)`
	err := db.QueryRow(query, randomString("tam")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTamanho falhou: %v", err)
	}
	return id

}

func CreateFuncionario(t *testing.T, db *sql.DB,IdDepartamento,IdFuncao int64) int64 {


	var id int64
	query:= `insert into funcionario(nome, matricula, IdDepartamento, IdFuncao) OUTPUT INSERTED.id values( @p1, @p2, @p3, @p4)`

	r:= randomInt()

	nome:= randomString("rada")
	matricula:=fmt.Sprintf("%d", r.Intn(9999999))

	err:= db.QueryRow(query,
		nome,
		matricula,
		IdDepartamento,
		IdFuncao,
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateFuncionario falhou: %v", err)
	}

	return id

}

func CreateEpi(t *testing.T, db *sql.DB, idTipoProtecao int64) int64 {
	var id int64

	// Usamos @p1, @p2... para o driver mapear automaticamente os argumentos na ordem
	query := `
		INSERT INTO epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
		OUTPUT INSERTED.id 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
	`
	// 1. Gera seed baseada no tempo para garantir aleatoriedade a cada execução

	r:= randomInt()
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

func CreateEntradaEpi(t  *testing.T, db *sql.DB,IdFuncionario, idEpi, IdTipoProtecao, IdTamanho int64) int64 {

	var id int64

	query:= `insert into entrada_epi(IdEpi,IdTamanho, data_entrada, quantidade,quantidadeAtual,data_fabricacao,
			 data_validade, lote, fornecedor, valor_unitario) OUTPUT INSERTED.id 
				values (@p1,@p2, @p3, @p4 ,
						@p5, @p6, @p7, @p8,@p9, @p10 )`

	
	lote:= randomString("lote")
	fabricante := randomString("fabr")
	err:= db.QueryRow(query,
		idEpi,
		IdTamanho,
		configs.NewDataBrPtr(time.Now()),
		100,
		100,
		configs.NewDataBrPtr(time.Now().AddDate(-1,0,0)),
		configs.NewDataBrPtr(time.Now().AddDate(1,0,0)),
		lote,
		fabricante,
		decimal.NewFromFloat(23.99),
	).Scan(&id)
	 
	if err != nil {
		t.Fatalf("Helper CreateEntradaEpi falhou: %v", err)
	}

	return id
}
