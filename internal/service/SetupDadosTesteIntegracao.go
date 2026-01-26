package service

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func randomString(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func randomInt() *rand.Rand {

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	return r
}

func CreateUser(t *testing.T, db *pgxpool.Pool) int64 {

	var id int64

	query := `insert into usuarios (nome, email, senha_hash) values ($1, $2, $3)
	RETURNING id;`

	nome := randomString("rada")
	email := randomString("rafa@gmail")
	senha := randomString("teste")
	err := db.QueryRow(context.Background(), query,
		nome,
		email,
		senha,
	).Scan(&id)
	if err != nil {

		t.Fatalf("erro ao criar usuario,%v", err)
	}

	return id
}
func CreateDepartamento(t *testing.T, db *pgxpool.Pool) int64 {

	var id int64
	query := `insert into departamento (nome) 
			values ($1)
			RETURNING id;`

	err := db.QueryRow(context.Background(), query, randomString("dep")).Scan(&id)
	if err != nil {
		t.Fatalf("dados de deparatmento falhou durante sua criacao: %v", err)
	}

	return id
}

func CreateFuncao(t *testing.T, db *pgxpool.Pool, idDep int64) int64 {

	var id int64
	query := `
		INSERT INTO funcao (nome, IdDepartamento) 
			VALUES ($1, $2) 
			RETURNING id;
	`

	err := db.QueryRow(context.Background(), query, randomString("func"), idDep).Scan(&id)
	if err != nil {
		t.Fatalf("dados de funcao falhou durante sua criação: %v", err)
	}

	return id
}

func CreateProtecao(t *testing.T, db *pgxpool.Pool) int64 {

	var id int64

	query := `INSERT INTO tipo_protecao (nome) 
		VALUES ($1)
		RETURNING id;
		`
	err := db.QueryRow(context.Background(), query, randomString("protec")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateProtecao falhou: %v", err)
	}

	return id

}

func CreateTamanho(t *testing.T, db *pgxpool.Pool) int64 {

	var id int64
	query := `INSERT INTO tamanho (tamanho) 
		VALUES ($1)
		RETURNING id;`
	err := db.QueryRow(context.Background(), query, randomString("tam")).Scan(&id)
	if err != nil {
		t.Fatalf("Helper CreateTamanho falhou: %v", err)
	}
	return id
}

func CreateFuncionario(t *testing.T, db *pgxpool.Pool, IdDepartamento, IdFuncao int64) int64 {

	var id int64
	query := `INSERT INTO funcionario (nome, matricula, IdDepartamento, IdFuncao) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`

	r := randomInt()

	nome := randomString("rada")
	matricula := fmt.Sprintf("%d", r.Intn(9999999))

	err := db.QueryRow(context.Background(), query,
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

func CreateEpi(t *testing.T, db *pgxpool.Pool, idTipoProtecao int64) int64 {
	var id int64

	// Usamos @p1, @p2... para o driver mapear automaticamente os argumentos na ordem
	query := `
		INSERT INTO epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`
	// 1. Gera seed baseada no tempo para garantir aleatoriedade a cada execução

	r := randomInt()
	// Gerando dados aleatórios para garantir unicidade
	nome := randomString("Luva")
	fabricante := randomString("Fabr")
	ca := fmt.Sprintf("%d", r.Intn(9999999))
	descricao := "Descrição de teste gerada automaticamente"
	validade := time.Now().AddDate(1, 0, 0) // Validade para daqui a 1 ano
	alertaMinimo := 10

	// A ordem aqui deve bater com @p1, @p2, etc.
	err := db.QueryRow(context.Background(), query,
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

func CreateEntradaEpi(t *testing.T, db *pgxpool.Pool, IdFuncionario, idEpi, IdTipoProtecao, IdTamanho, iduser int64) int64 {

	var id int64

	query := `INSERT INTO entrada_epi (
    IdEpi, IdTamanho, data_entrada, quantidade, quantidadeAtual, 
    data_fabricacao, data_validade, lote, fornecedor, valor_unitario,nota_fiscal_numero, nota_fiscal_serie,id_usuario_criacao
	) 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id;`

	lote := randomString("lote")
	fabricante := randomString("fabr")
	notaFiscal := randomString("notaFiscal")
	notaFiscalNumero := randomString("notaFiscalNumero")
	err := db.QueryRow(context.Background(), query,
		idEpi,
		IdTamanho,
		configs.NewDataBrPtr(time.Now()),
		100,
		100,
		configs.NewDataBrPtr(time.Now().AddDate(-1, 0, 0)),
		configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
		lote,
		fabricante,
		decimal.NewFromFloat(23.99),
		notaFiscal,
		notaFiscalNumero,
		iduser,
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateEntradaEpi falhou: %v", err)
	}

	return id
}

func CreateEntradaEpi1(t *testing.T, db *pgxpool.Pool, IdFuncionario, idEpi, IdTipoProtecao, IdTamanho, iduser int64) int64 {

	var id int64

	query := `INSERT INTO entrada_epi (
    IdEpi, IdTamanho, data_entrada, quantidade, quantidadeAtual, 
    data_fabricacao, data_validade, lote, fornecedor, valor_unitario,nota_fiscal_numero, nota_fiscal_serie,id_usuario_criacao
	) 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id;`

	lote := randomString("lote")
	fabricante := randomString("fabr")
	notaFiscal := randomString("notaFiscal")
	notaFiscalNumero := randomString("notaFiscalNumero")
	err := db.QueryRow(context.Background(), query,
		idEpi,
		IdTamanho,
		configs.NewDataBrPtr(time.Now()),
		1,
		1,
		configs.NewDataBrPtr(time.Now().AddDate(-1, 0, 0)),
		configs.NewDataBrPtr(time.Now().AddDate(1, 0, 0)),
		lote,
		fabricante,
		decimal.NewFromFloat(23.99),
		notaFiscal,
		notaFiscalNumero,
		iduser,
	).Scan(&id)

	if err != nil {
		t.Fatalf("Helper CreateEntradaEpi falhou: %v", err)
	}

	return id
}

func CreateEntregaEpi(t *testing.T, db *pgxpool.Pool, idFuncionario, idUser int64) int64 {
	var id int64
	query := `
		INSERT INTO entrega_epi (IdFuncionario, data_entrega, assinatura, token_validacao,id_usuario_entrega)
				VALUES ($1, $2, $3, $4, $5)
					RETURNING id;
	`
	dataEntrega:= configs.NewDataBrPtr(time.Now())
	assinatura:= randomString("teste")
	token:= randomString("oq")
	
	err:= db.QueryRow(context.Background(), query,

			idFuncionario,
			dataEntrega,
			assinatura,
			token,
			idUser,
		).Scan(&id)
	if err != nil{

		t.Fatalf("Helper CreateEntregaEpi falhou: %v", err)

	}
	return id
}


func CreateEpiEntregues(t *testing.T, db *pgxpool.Pool,idEntrega, idEntrada, IdEpi,IdTamanho int64) int64 {


	var id int64
	query:= `
		INSERT INTO epis_entregues (IdEntrega, IdEntrada ,IdEpi, IdTamanho, quantidade)
			VALUES ($1, $2, $3, $4, $5)
				RETURNING id;
	`

	quantidade:= 10

	err:= db.QueryRow(context.Background(), query,

		idEntrega,
		idEntrada,
		IdEpi,
		IdTamanho,
		quantidade,
	).Scan(&id)
	if err != nil {

		t.Fatalf("erro ao criar CreateEpiEntregues: %v", err)
	}

	return id
}

func CreateMotivoDevolucao(t *testing.T, db *pgxpool.Pool, motivo string) int64{

	var id int64

	query:= `	
		INSERT INTO motivo_devolucao (motivo) 
			VALUES ($1)
			returning id;
	`

	err:= db.QueryRow(context.Background(), query,
			motivo,
		).Scan(&id)
	if err != nil {

		t.Fatalf("erro ao criar CreateMotivoDevolucao: %v", err)
	}

	return  id
}