-- name: AddEntregaEpi :one
INSERT INTO entrega_epi (IdFuncionario, data_entrega, assinatura)
VALUES ($1, $2, $3)
RETURNING id;

-- name: AddItemEntregue :exec
INSERT INTO epis_entregues (IdEntrega, IdEpi, IdTamanho, quantidade, valor_unitario)
VALUES ($1, $2, $3, $4, $5);

-- name: ListarEntregas :many
SELECT 
    ee.id as entrega_id, ee.data_entrega, ee.assinatura,
    f.id as func_id, f.nome as func_nome,
    d.id as dep_id, d.nome as dep_nome,
    ff.id as funcao_id, ff.nome as funcao_nome,
    e.id as epi_id, e.nome as epi_nome, e.fabricante, e.CA, e.descricao as epi_desc, e.validade_CA,
    tp.id as tp_id, tp.nome as tp_nome,
    t.id as tam_id, t.tamanho as tam_nome,
    i.quantidade, i.valor_unitario
FROM entrega_epi ee
INNER JOIN funcionario f ON ee.IdFuncionario = f.id
INNER JOIN departamento d ON f.IdDepartamento = d.id
INNER JOIN funcao ff ON f.IdFuncao = ff.id
INNER JOIN epis_entregues i ON i.IdEntrega = ee.id
INNER JOIN epi e ON i.IdEpi = e.id
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
INNER JOIN tamanho t ON i.IdTamanho = t.id
WHERE 
    ((sqlc.arg('canceladas')::boolean IS FALSE AND ee.cancelada_em IS NULL) OR
     (sqlc.arg('canceladas')::boolean IS TRUE AND ee.cancelada_em IS NOT NULL))
AND (sqlc.narg('id_entrega')::int IS NULL OR ee.id = sqlc.narg('id_entrega'))
AND (sqlc.narg('id_funcionario')::int IS NULL OR ee.IdFuncionario = sqlc.narg('id_funcionario'))
ORDER BY ee.data_entrega DESC;

-- name: CancelarEntrega :execrows
UPDATE entrega_epi
SET cancelada_em = NOW(),
    ativo = FALSE
WHERE id = $1 AND cancelada_em IS NULL;