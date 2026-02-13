-- 1. Transforma a coluna de Texto para Inteiro (e já limpa qualquer resquício)
ALTER TABLE entrada_epi 
ALTER COLUMN fornecedor TYPE INTEGER USING NULL;

-- 2. Renomeia para ficar no padrão de chave estrangeira (Opcional, mas recomendado)
ALTER TABLE entrada_epi 
RENAME COLUMN fornecedor TO Idfornecedor;

-- 3. Cria a ligação (Foreign Key) com a tabela de fornecedores
ALTER TABLE entrada_epi 
ADD CONSTRAINT fk_entrada_fornecedor 
FOREIGN KEY (Idfornecedor) REFERENCES fornecedores(id);