-- 1. Remove a regra antiga (que estava sem o tenant_id)
ALTER TABLE entrada_epi 
DROP CONSTRAINT uk_entrada_nf_fornecedor;

-- 2. Adiciona a regra nova (AGORA COM O TENANT_ID)
ALTER TABLE entrada_epi
ADD CONSTRAINT unique_entrada_Nf
UNIQUE (tenant_id, fornecedor, nota_fiscal_numero, nota_fiscal_serie);