CREATE TABLE departamento (
    id INT PRIMARY KEY IDENTITY(1,1),
    nome VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE funcao (
    id int primary key identity(1,1),
    nome varchar(100) not null unique,
    IdDepartamento int not null, -- Vírgula adicionada aqui
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

CREATE TABLE funcionario (
    id int primary key identity(1,1),
    nome varchar(100) not null,
    matricula varchar(7) not null unique,
    IdFuncao int not null,
    IdDepartamento int not null, -- Vírgula adicionada aqui
    foreign key (IdFuncao) references funcao(id),
    foreign key (IdDepartamento) references departamento(id)
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

CREATE TABLE epis_entregues(
    id int primary key identity(1,1),
    IdEntrega int not null,
    IdEntrada int not null,
    IdEpi int not null,
    IdTamanho int not null,
    quantidade int not null,
    valor_unitario decimal(10,2) not null, -- Vírgula adicionada aqui
    foreign key (IdEntrega) references entrega_epi(id),
    foreign key (IdEpi) references epi(id),
    foreign key (IdEntrada) references entrada_epi(id)
);

CREATE TABLE motivo_devolucao(
    id int primary key identity(1,1),
    motivo varchar(50) not null unique
);

CREATE TABLE devolucao(
    id int primary key identity(1,1),
    IdEpi int not null,
    IdFuncionario int not null,
    IdMotivo int not null,
    data_devolucao date not null,
    IdTamanho int not null,
    quantidadeAdevolver int not null,
    IdEpiNovo int null,
    IdTamanhoNovo int null,
    quantidadeNova int null,
    cancelada_em datetime null,
    assinatura_digital varbinary(MAX) not null, -- MUDADO: De varbinary para varbinary(MAX) e adicionado vírgula antes
    foreign key (IdEpi) references epi(id),
    foreign key (IdFuncionario) references funcionario(id),
    foreign key (IdMotivo) references motivo_devolucao(id),
    foreign key (IdEpiNovo) references epi(id),
    foreign key (IdTamanhoNovo) references tamanho(id),
    foreign key (IdTamanho) references tamanho(id)
);