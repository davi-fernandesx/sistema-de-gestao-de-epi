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