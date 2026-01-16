ALTER TABLE entrada_epi 
ADD COLUMN nota_fiscal_numero VARCHAR(50) NOT NULL,
ADD COLUMN nota_fiscal_serie  VARCHAR(10) DEFAULT '1';

-- Criando a regra de que não pode repetir a mesma NF para o mesmo Fornecedor
ALTER TABLE entrada_epi 
ADD CONSTRAINT uk_entrada_nf_fornecedor 
UNIQUE (nota_fiscal_numero, nota_fiscal_serie, fornecedor);