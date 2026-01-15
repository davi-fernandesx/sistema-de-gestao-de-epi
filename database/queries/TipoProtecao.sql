-- name: AddProtecao :exec
INSERT INTO tipo_protecao (nome) 
VALUES ($1);

-- name: BuscarProtecao :one
SELECT id, nome 
FROM tipo_protecao 
WHERE id = $1 AND ativo = TRUE 
LIMIT 1;

-- name: BuscarTodasProtecoes :many
SELECT id, nome 
FROM tipo_protecao 
WHERE ativo = TRUE
ORDER BY nome ASC;

-- name: DeletarProtecao :execrows
UPDATE tipo_protecao
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 AND ativo = TRUE;

-- name: UpdateProtecao :execrows
UPDATE tipo_protecao
SET nome = $2
WHERE id = $1 AND ativo = TRUE;