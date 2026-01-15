-- name: AddEpi :one
INSERT INTO epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: AddEpiTamanho :exec
INSERT INTO tamanhos_epis (IdEpi, IdTamanho) 
VALUES ($1, $2);

-- name: BuscarEpi :one
SELECT 
    e.id, e.nome, e.fabricante, e.CA, e.descricao,
    e.validade_CA, e.alerta_minimo, e.IdTipoProtecao, 
    tp.nome as tipo_protecao_nome
FROM epi e
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
WHERE e.id = $1 AND e.ativo = TRUE;

-- name: BuscarTamanhosPorEpi :many
SELECT t.id, t.tamanho
FROM tamanho t
INNER JOIN tamanhos_epis te ON t.id = te.IdTamanho
WHERE te.IdEpi = $1 AND te.ativo = TRUE;

-- name: BuscarTodosEpis :many
SELECT 
    e.id, e.nome, e.fabricante, e.CA, e.descricao,
    e.validade_CA, e.alerta_minimo, e.IdTipoProtecao, 
    tp.nome as tipo_protecao_nome
FROM epi e
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
WHERE e.ativo = TRUE;

-- name: BuscarTodosTamanhosAgrupados :many
SELECT te.IdEpi, t.id, t.tamanho
FROM tamanho t
INNER JOIN tamanhos_epis te ON t.id = te.IdTamanho
WHERE te.ativo = TRUE;

-- name: DeletarEpi :execrows
UPDATE epi SET ativo = FALSE, deletado_em = NOW() WHERE id = $1 AND ativo = TRUE;

-- name: DeletarTamanhosPorEpi :exec
UPDATE tamanhos_epis SET ativo = FALSE, deletado_em = NOW() WHERE IdEpi = $1 AND ativo = TRUE;

-- name: UpdateEpiCampo :execrows
UPDATE epi 
SET nome = COALESCE(sqlc.narg('nome'), nome),
    fabricante = COALESCE(sqlc.narg('fabricante'), fabricante),
    CA = COALESCE(sqlc.narg('ca'), CA),
    descricao = COALESCE(sqlc.narg('descricao'), descricao),
    validade_CA = COALESCE(sqlc.narg('validade_ca'), validade_CA)
WHERE id = sqlc.arg('id') AND ativo = TRUE;