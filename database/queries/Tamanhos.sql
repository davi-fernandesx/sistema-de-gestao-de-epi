-- name: AddTamanho :exec
INSERT INTO tamanho (tamanho) 
VALUES ($1);

-- name: BuscarTamanho :one
SELECT id, tamanho 
FROM tamanho 
WHERE id = $1 AND ativo = TRUE 
LIMIT 1;

-- name: BuscarTamanhosPorIdEpi :many
SELECT t.id, t.tamanho
FROM tamanho t
INNER JOIN tamanhos_epis te ON t.id = te.IdTamanho
WHERE te.IdEpi = $1 AND te.ativo = TRUE;

-- name: BuscarTodosTamanhos :many
SELECT id, tamanho 
FROM tamanho 
WHERE ativo = TRUE
ORDER BY tamanho ASC;

-- name: DeletarTamanho :execrows
UPDATE tamanho
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 AND ativo = TRUE;

-- name: UpdateEpiNosTamanhos :execrows
-- Esta query atualiza a associação na tabela muitos-para-muitos
UPDATE tamanhos_epis
SET IdEpi = $2
WHERE IdEpi = $1 AND ativo = TRUE;