package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func criarTabelasPostgres(t *testing.T, pool *pgxpool.Pool) {

	
	schema := `
			-- TABELA: departamento
		CREATE TABLE departamento (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL UNIQUE,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL
		);

		-- TABELA: funcao
		CREATE TABLE funcao (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL UNIQUE,
			IdDepartamento INT NOT NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL,
			FOREIGN KEY (IdDepartamento) REFERENCES departamento(id)
		);

		-- TABELA: tipo_protecao
		CREATE TABLE tipo_protecao (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL UNIQUE,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL
		);

		-- TABELA: tamanho
		CREATE TABLE tamanho (
			id SERIAL PRIMARY KEY,
			tamanho VARCHAR(50) NOT NULL UNIQUE,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL
		);

		-- TABELA: epi
		CREATE TABLE epi (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			fabricante VARCHAR(100) NOT NULL,
			CA VARCHAR(20) NOT NULL UNIQUE,
			descricao TEXT NOT NULL,
			validade_CA DATE NOT NULL,
			IdTipoProtecao INT NOT NULL,
			alerta_minimo INT NOT NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL,
			FOREIGN KEY (IdTipoProtecao) REFERENCES tipo_protecao(id)
		);

		-- TABELA: tamanhos_epis
		CREATE TABLE tamanhos_epis (
			id SERIAL PRIMARY KEY,
			IdEpi INT NOT NULL,
			IdTamanho INT NOT NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL,
			FOREIGN KEY (IdEpi) REFERENCES epi(id),
			FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
		);

		-- TABELA: funcionario
		CREATE TABLE funcionario (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			matricula VARCHAR(7) NOT NULL UNIQUE,
			IdFuncao INT NOT NULL,
			IdDepartamento INT NOT NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL,
			FOREIGN KEY (IdFuncao) REFERENCES funcao(id),
			FOREIGN KEY (IdDepartamento) REFERENCES departamento(id)
		);

		-- TABELA: entrada_epi
		CREATE TABLE entrada_epi (
			id SERIAL PRIMARY KEY,
			IdEpi INT NOT NULL,
			IdTamanho INT NOT NULL,
			data_entrada DATE NOT NULL,
			quantidade INT NOT NULL,
			quantidadeAtual INT NOT NULL,
			data_fabricacao DATE NOT NULL,
			data_validade DATE NOT NULL,
			lote VARCHAR(50) NOT NULL,
			fornecedor VARCHAR(100) NOT NULL,
			valor_unitario DECIMAL(10,2) NOT NULL,
			cancelada_em TIMESTAMP NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			FOREIGN KEY (IdEpi) REFERENCES epi(id),
			FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
		);

		-- TABELA: entrega_epi
		CREATE TABLE entrega_epi (
			id SERIAL PRIMARY KEY,
			IdFuncionario INT NOT NULL,
			data_entrega DATE NOT NULL,
			assinatura TEXT NOT NULL, 
			IdTroca INT NULL,
			cancelada_em TIMESTAMP NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			FOREIGN KEY (IdFuncionario) REFERENCES funcionario(id)
		);

		-- TABELA: epis_entregues
		CREATE TABLE epis_entregues (
			id SERIAL PRIMARY KEY,
			IdEntrega INT NOT NULL,
			IdEntrada INT NOT NULL,
			IdEpi INT NOT NULL,
			IdTamanho INT NOT NULL,
			quantidade INT NOT NULL,
			valor_unitario DECIMAL(10,2) NOT NULL,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL,
			FOREIGN KEY (IdEntrega) REFERENCES entrega_epi(id),
			FOREIGN KEY (IdEpi) REFERENCES epi(id),
			FOREIGN KEY (IdEntrada) REFERENCES entrada_epi(id)
		);

		-- TABELA: motivo_devolucao
		CREATE TABLE motivo_devolucao (
			id SERIAL PRIMARY KEY,
			motivo VARCHAR(50) NOT NULL UNIQUE,
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			deletado_em TIMESTAMP NULL
		);

		-- TABELA: devolucao
		CREATE TABLE devolucao (
			id SERIAL PRIMARY KEY,
			IdEpi INT NOT NULL,
			IdFuncionario INT NOT NULL,
			IdMotivo INT NOT NULL,
			data_devolucao DATE NOT NULL,
			IdTamanho INT NOT NULL,
			quantidadeAdevolver INT NOT NULL,
			IdEpiNovo INT NULL,
			IdTamanhoNovo INT NULL,
			quantidadeNova INT NULL,
			cancelada_em TIMESTAMP NULL,
			-- Aqui usei TEXT pensando na URL da assinatura do Supabase
			assinatura_digital TEXT NOT NULL, 
			ativo BOOLEAN NOT NULL DEFAULT TRUE,
			FOREIGN KEY (IdEpi) REFERENCES epi(id),
			FOREIGN KEY (IdFuncionario) REFERENCES funcionario(id),
			FOREIGN KEY (IdMotivo) REFERENCES motivo_devolucao(id),
			FOREIGN KEY (IdEpiNovo) REFERENCES epi(id),
			FOREIGN KEY (IdTamanhoNovo) REFERENCES tamanho(id),
			FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
		);

		ALTER TABLE entrada_epi 
		ADD COLUMN nota_fiscal_numero VARCHAR(50) NOT NULL,
		ADD COLUMN nota_fiscal_serie  VARCHAR(50) DEFAULT '1';

		-- Criando a regra de que não pode repetir a mesma NF para o mesmo Fornecedor
		ALTER TABLE entrada_epi 
		ADD CONSTRAINT uk_entrada_nf_fornecedor 
		UNIQUE (nota_fiscal_numero, nota_fiscal_serie, fornecedor);

		ALTER TABLE entrega_epi ADD COLUMN token_validacao TEXT;
	
		CREATE TABLE usuarios (
		id SERIAL PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		senha_hash TEXT NOT NULL,
		ativo BOOLEAN DEFAULT TRUE
		);

		-- 1. Rastrear quem deu entrada nos lotes/EPIs
		ALTER TABLE entrada_epi 
		ADD COLUMN id_usuario_criacao INTEGER REFERENCES usuarios(id);

		-- 2. Rastrear quem realizou a entrega para o funcionário
		ALTER TABLE entrega_epi 
		ADD COLUMN id_usuario_entrega INTEGER REFERENCES usuarios(id);

		-- 3. Rastrear quem realizou o cancelamento (estorno)
		ALTER TABLE devolucao 
		ADD COLUMN id_usuario_cancelamento INTEGER REFERENCES usuarios(id);


		ALTER TABLE entrada_epi 
		ADD COLUMN id_usuario_criacao_cancelamento INTEGER REFERENCES usuarios(id);

		-- 2. Rastrear quem realizou a entrega para o funcionário
		ALTER TABLE entrega_epi 
		ADD COLUMN id_usuario_entrega_cancelamento INTEGER REFERENCES usuarios(id);

		-- 3. Rastrear quem realizou o cancelamento (estorno)
		ALTER TABLE devolucao 
		ADD COLUMN id_usuario_devolucao_cancelamento INTEGER REFERENCES usuarios(id);
	`

	_, err:= pool.Exec(context.Background(), schema)
	if err != nil {

		t.Fatalf("erro  na criação da tabelas, %v", err)
	}

}
