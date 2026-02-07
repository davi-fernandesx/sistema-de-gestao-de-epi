CREATE UNIQUE INDEX idx_protecao_tenant_ativo 
ON tipo_protecao (tenant_id, nome) 
WHERE deletado_em IS NULL;