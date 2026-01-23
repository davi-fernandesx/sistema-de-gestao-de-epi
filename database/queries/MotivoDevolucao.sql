-- name: AddMotivoDevolucao :exec
INSERT INTO motivo_devolucao (motivo) 
VALUES ($1);

-- name: BuscaMotivoDevolucao :one
SELECT id, motivo 
FROM motivo_devolucao 
WHERE id = $1 AND ativo = TRUE 
LIMIT 1;

-- name: BuscaTodosMotivosDevolucao :many
SELECT id, motivo 
FROM motivo_devolucao 
WHERE ativo = TRUE
ORDER BY motivo ASC;

-- name: DeleteMotivoDevolucao :execrows
UPDATE motivo_devolucao
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 AND ativo = TRUE;