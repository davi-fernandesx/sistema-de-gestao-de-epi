-- name: AddEntregaEpi :one
INSERT INTO entrega_epi (IdFuncionario, data_entrega, assinatura, token_validacao,id_usuario_entrega)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: AddItemEntregue :one
INSERT INTO epis_entregues (IdEntrega, IdEntrada ,IdEpi, IdTamanho, quantidade, valor_unitario)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING IdEntrega;

-- name: CancelaItemEntregue :exec
UPDATE epis_entregues
set ativo = FALSE, deletado_em = NOW()
WHERE IdEntrega = $1 ; 

-- name: ListarEntregas :many
SELECT 
    ee.id as entrega_id, ee.data_entrega, ee.assinatura,ee.token_validacao,ee.id_usuario_entrega,   
    f.id as func_id, f.nome as func_nome, f.matricula,
    d.id as dep_id, d.nome as dep_nome,
    ff.id as funcao_id, ff.nome as funcao_nome,
    e.id as epi_id, e.nome as epi_nome, e.fabricante, e.CA, e.descricao as epi_desc, e.validade_CA,
    tp.id as tp_id, tp.nome as tp_nome,
    t.id as tam_id, t.tamanho as tam_nome,
    i.quantidade, i.valor_unitario,
    COUNT(*) OVER() as total_geral
FROM entrega_epi ee
INNER JOIN funcionario f ON ee.IdFuncionario = f.id
INNER JOIN departamento d ON f.IdDepartamento = d.id
INNER JOIN funcao ff ON f.IdFuncao = ff.id
INNER JOIN epis_entregues i ON i.IdEntrega = ee.id
INNER JOIN epi e ON i.IdEpi = e.id
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
INNER JOIN tamanho t ON i.IdTamanho = t.id

WHERE 
    (   (@canceladas::boolean IS FALSE AND ee.cancelada_em IS NULL) OR
        (@canceladas::boolean IS TRUE AND ee.cancelada_em IS NOT NULL))
AND (sqlc.narg('id_entrega')::int IS NULL OR ee.id = sqlc.narg('id_entrega'))
AND (sqlc.narg('id_funcionario')::int IS NULL OR ee.IdFuncionario = sqlc.narg('id_funcionario'))
and (sqlc.narg('data_inicio')::date IS NULL OR ee.data_entrega >= sqlc.narg('data_inicio'))
and (sqlc.narg('data_fim')::date IS NULL OR ee.data_entrega <= sqlc.narg('data_fim'))
ORDER BY ee.data_entrega DESC
limit $1 offset $2;

-- name: CancelarEntrega :one
UPDATE entrega_epi
SET cancelada_em = NOW(),
    ativo = FALSE,
    id_usuario_entrega_cancelamento = $2
WHERE id = $1 AND cancelada_em IS NULL
RETURNING id;

-- name: BuscarTodosItensEntrega :many
SELECT 
    i.IdEntrega as entrega_id,i.id as item_id , i.quantidade, i.valor_unitario,
    e.id as epi_id, e.nome as epi_nome, e.fabricante, e.CA, e.descricao as epi_desc, e.validade_CA,
    tp.id as tp_id, tp.nome as tp_nome,
    t.id as tam_id, t.tamanho as tam_nome
FROM  epis_entregues i
INNER JOIN epi e ON i.IdEpi = e.id
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
INNER JOIN tamanho t ON i.IdTamanho = t.id
WHERE  i.ativo = TRUE;

-- name: ListarItensEntregueCancelados :many
select quantidade,IdEntrada
from epis_entregues
where IdEntrega = $1 and ativo = FALSE and deletado_em is not null;