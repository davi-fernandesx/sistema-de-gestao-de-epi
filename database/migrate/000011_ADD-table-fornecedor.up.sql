CREATE TABLE fornecedores (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    razao_social VARCHAR(100) NOT NULL,
    nome_fantasia VARCHAR(100) NOT NULL,
    cnpj VARCHAR(14) NOT NULL, -- Apenas números
    inscricao_estadual VARCHAR(50) NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    cancelado_em TIMESTAMP DEFAULT NULL   
);

-- Índices e Constraints (A parte mais importante!)

-- 1. Garante busca rápida pelo ID e Tenant (segurança)
CREATE INDEX idx_fornecedores_tenant ON fornecedores(tenant_id);

-- 2. Garante que NÃO existam dois fornecedores com o mesmo CNPJ na MESMA empresa (tenant)
-- Mas permite o mesmo CNPJ em tenants diferentes.
ALTER TABLE fornecedores
ADD CONSTRAINT unique_cnpj_por_tenant UNIQUE (tenant_id, cnpj);