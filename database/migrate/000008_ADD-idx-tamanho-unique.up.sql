CREATE UNIQUE INDEX idx_tamanho_tenant_ativo 
ON tamanho (tenant_id, tamanho) 
WHERE deletado_em IS NULL;