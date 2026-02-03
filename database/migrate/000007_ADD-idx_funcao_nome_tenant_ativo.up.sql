CREATE UNIQUE INDEX idx_funcao_nome_tenant_ativo
ON funcao (tenant_id, nome)
WHERE deletado_em IS NULL;