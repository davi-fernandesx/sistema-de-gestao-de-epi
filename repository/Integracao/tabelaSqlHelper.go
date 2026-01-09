package integracao

import (
	"database/sql"
	"testing"
)

func criarTabelasDoUsuario(t *testing.T, db *sql.DB) {
	// Copiei exatamente o seu esquema SQL
	schema := `
	CREATE TABLE departamento (
		id INT PRIMARY KEY IDENTITY(1,1),
		nome VARCHAR(100) NOT NULL UNIQUE
	);

	CREATE TABLE funcao (
		id int primary key identity(1,1),
		nome varchar(100) not null unique,
		IdDepartamento int not null,
		foreign key (IdDepartamento) references departamento(id)
	);

	CREATE TABLE funcionario (
		id int primary key identity(1,1),
		nome varchar(100) not null,
		matricula varchar(7) not null unique,
		IdFuncao int not null,
		IdDepartamento int not null,
		foreign key (IdFuncao) references funcao(id),
		foreign key (IdDepartamento) references departamento(id)
	);


		CREATE TABLE tipo_protecao(
		id int primary key identity(1,1),
		nome varchar(100) not null unique
	);

		CREATE TABLE tamanho(
		id int primary key identity(1,1),
		tamanho varchar(50) not null unique
	);

	
		CREATE TABLE epi(
		id int primary key identity(1,1),
		nome varchar(100) not null,
		fabricante varchar(100) not null,
		CA varchar(20) not null unique,
		descricao text not null,
		validade_CA date not null,
		IdTipoProtecao int not null, -- Vírgula adicionada aqui
		alerta_minimo INT NOT NULL,
		foreign key (IdTipoProtecao) references tipo_protecao(id)
	);

	
		CREATE TABLE tamanhos_epis(
		id int primary key identity(1,1),
		IdEpi int not null,
		IdTamanho int not null, -- Vírgula adicionada aqui
		foreign key (IdEpi) references epi(id),
		foreign key (IdTamanho) references tamanho(id)
	);

		CREATE TABLE entrada_epi(
		id int primary key identity(1,1),
		IdEpi int not null,
		IdTamanho int not null,
		data_entrada date not null,
		quantidade int not null,
		data_fabricacao date not null,
		data_validade date not null,
		lote varchar(50) not null,
		fornecedor varchar(100) not null,
		valor_unitario decimal(10,2) not null,
		cancelada_em datetime null, -- Vírgula adicionada aqui
		foreign key (IdEpi) references epi(id), -- Vírgula adicionada aqui entre as FKs
		foreign key (IdTamanho) references tamanho(id)
	);

		CREATE TABLE entrega_epi(
		id int primary key identity(1,1),
		IdFuncionario int not null,
		data_entrega date not null,
		assinatura varbinary(MAX) not null, -- MUDADO: De varbinary para varbinary(MAX)
		IdTroca int null,
		cancelada_em datetime null, -- Vírgula adicionada aqui
		foreign key (IdFuncionario) references funcionario(id)
	);

		-- TABELA: funcionario
		ALTER TABLE funcionario ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE funcionario ADD deletado_em DATETIME NULL;

		-- TABELA: departamento
		ALTER TABLE departamento ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE departamento ADD deletado_em DATETIME NULL;

		-- TABELA: funcao
		ALTER TABLE funcao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE funcao ADD deletado_em DATETIME NULL;

		-- TABELA : epi
		ALTER TABLE epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE epi ADD deletado_em DATETIME NULL;

		-- TABELA : tipo_protecao
		ALTER TABLE tipo_protecao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE tipo_protecao ADD deletado_em DATETIME NULL;

		-- TABELA: tamanho
		ALTER TABLE tamanho ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE tamanho ADD deletado_em DATETIME NULL;

		-- TABELA: tamanhos_epis
		ALTER TABLE tamanhos_epis ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
		ALTER TABLE tamanhos_epis ADD deletado_em DATETIME NULL;

		-- TABELA: entrada_epi
		ALTER TABLE entrada_epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;

		-- TABELA: entrega_epi
		ALTER TABLE entrega_epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;

	`
	// Executa o CREATE TABLE
	_, err := db.Exec(schema)
	if err != nil {
		t.Fatalf("Erro ao criar tabelas: %v", err)
	}
}
