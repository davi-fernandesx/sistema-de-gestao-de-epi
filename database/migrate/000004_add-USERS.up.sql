-- 2. TABELA: USUARIOS (Agora vinculada a uma empresa)
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL, -- De qual empresa esse usuário é?
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL, -- Email pode repetir se for sistemas diferentes, mas o par (email, tenant) deve ser único
    senha_hash TEXT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (tenant_id) REFERENCES empresas(id),
    UNIQUE (tenant_id, email) -- Garante que o mesmo email não existe duplicado DENTRO da mesma empresa
);