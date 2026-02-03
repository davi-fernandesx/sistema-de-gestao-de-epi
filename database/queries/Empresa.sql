-- name: GetTenantBySubdomain :one
SELECT id, nome_fantasia 
FROM empresas 
WHERE subdominio = $1 AND ativo = TRUE;