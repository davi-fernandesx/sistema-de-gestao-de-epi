-- 1. TABELA MÃE: EMPRESAS (TENANTS)
-- Esta tabela gerencia quem são seus clientes (o frigorífico, a construtora, etc.)
CREATE TABLE empresas (
    id SERIAL PRIMARY KEY,
    nome_fantasia VARCHAR(100) NOT NULL,
    razao_social VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) UNIQUE NOT NULL,
    subdominio VARCHAR(50) UNIQUE NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- TABELA: departamento
CREATE TABLE departamento (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    nome VARCHAR(100) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id)
);

-- TABELA: funcao
CREATE TABLE funcao (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    nome VARCHAR(100) NOT NULL,
    IdDepartamento INT NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdDepartamento) REFERENCES departamento(id)
);

-- TABELA: tipo_protecao
CREATE TABLE tipo_protecao (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    nome VARCHAR(100) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id)
);

-- TABELA: tamanho
CREATE TABLE tamanho (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    tamanho VARCHAR(50) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id)
);

-- TABELA: epi
CREATE TABLE epi (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    nome VARCHAR(100) NOT NULL,
    fabricante VARCHAR(100) NOT NULL,
    CA VARCHAR(20) NOT NULL, -- Removi UNIQUE global, pois duas empresas podem cadastrar o mesmo CA. O ideal é UNIQUE(tenant_id, CA)
    descricao TEXT NOT NULL,
    validade_CA DATE NOT NULL,
    IdTipoProtecao INT NOT NULL,
    alerta_minimo INT NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdTipoProtecao) REFERENCES tipo_protecao(id),
    UNIQUE (tenant_id, CA) -- O CA é único apenas DENTRO da empresa
);

-- TABELA: tamanhos_epis
CREATE TABLE tamanhos_epis (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    IdEpi INT NOT NULL,
    IdTamanho INT NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdEpi) REFERENCES epi(id),
    FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
);

-- TABELA: funcionario
CREATE TABLE funcionario (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    nome VARCHAR(100) NOT NULL,
    matricula VARCHAR(20) NOT NULL, -- Aumentei para 20 e removi UNIQUE global
    IdFuncao INT NOT NULL,
    IdDepartamento INT NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdFuncao) REFERENCES funcao(id),
    FOREIGN KEY (IdDepartamento) REFERENCES departamento(id),
    UNIQUE (tenant_id, matricula) -- A matrícula só precisa ser única dentro daquela empresa
);

-- TABELA: entrada_epi
CREATE TABLE entrada_epi (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
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
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdEpi) REFERENCES epi(id),
    FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
);

-- TABELA: entrega_epi
CREATE TABLE entrega_epi (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    IdFuncionario INT NOT NULL,
    data_entrega DATE NOT NULL,
    assinatura TEXT NOT NULL, 
    IdTroca INT NULL,
    cancelada_em TIMESTAMP NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdFuncionario) REFERENCES funcionario(id)
);

-- TABELA: epis_entregues
CREATE TABLE epis_entregues (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    IdEntrega INT NOT NULL,
    IdEntrada INT NOT NULL, -- Importante para saber de qual lote (entrada) saiu
    IdEpi INT NOT NULL,
    IdTamanho INT NOT NULL,
    quantidade INT NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdEntrega) REFERENCES entrega_epi(id),
    FOREIGN KEY (IdEpi) REFERENCES epi(id),
    FOREIGN KEY (IdEntrada) REFERENCES entrada_epi(id)
);

-- TABELA: motivo_devolucao
CREATE TABLE motivo_devolucao (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    motivo VARCHAR(50) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    deletado_em TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id)
);

-- TABELA: devolucao
CREATE TABLE devolucao (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
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
    assinatura_digital TEXT NOT NULL, 
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    FOREIGN KEY (IdEpi) REFERENCES epi(id),
    FOREIGN KEY (IdFuncionario) REFERENCES funcionario(id),
    FOREIGN KEY (IdMotivo) REFERENCES motivo_devolucao(id),
    FOREIGN KEY (IdEpiNovo) REFERENCES epi(id),
    FOREIGN KEY (IdTamanhoNovo) REFERENCES tamanho(id),
    FOREIGN KEY (IdTamanho) REFERENCES tamanho(id)
);