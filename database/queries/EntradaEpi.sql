-- name: AddEntradaEpi :exec
INSERT INTO entrada_epi (
    IdEpi, IdTamanho, data_entrada, quantidade, quantidadeAtual, 
    data_fabricacao, data_validade, lote, fornecedor, valor_unitario,nota_fiscal_numero, nota_fiscal_serie
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: ListarEntradas :many
-- name: ListarEntradas :many
SELECT 
    ee.id, ee.IdEpi, e.nome as epi_nome, e.fabricante, e.CA, e.descricao as epi_descricao,
    ee.data_fabricacao, ee.data_validade, e.validade_CA,
    e.IdTipoProtecao, tp.nome as protecao_nome,
    ee.IdTamanho, t.tamanho as tamanho_nome, 
    ee.quantidade, ee.quantidadeAtual, ee.data_entrada,
    ee.lote, ee.fornecedor, ee.valor_unitario, ee.nota_fiscal_numero, ee.nota_fiscal_serie
FROM entrada_epi ee
INNER JOIN epi e ON ee.IdEpi = e.id
INNER JOIN tipo_protecao tp ON e.IdTipoProtecao = tp.id
INNER JOIN tamanho t ON ee.IdTamanho = t.id
WHERE 
    (
        (@canceladas::boolean IS FALSE AND ee.cancelada_em IS NULL) OR
        (@canceladas::boolean IS TRUE AND ee.cancelada_em IS NOT NULL)
    )
    AND (sqlc.narg('id_epi')::int IS NULL OR ee.IdEpi = sqlc.narg('id_epi'))
    AND (sqlc.narg('id_entrada')::int IS NULL OR ee.id = sqlc.narg('id_entrada'))
    AND (sqlc.narg('data_inicio')::date IS NULL OR ee.data_entrada >= sqlc.narg('data_inicio'))
    AND (sqlc.narg('data_fim')::date IS NULL OR ee.data_entrada <= sqlc.narg('data_fim'))
    AND (sqlc.narg('nota_fiscal')::text IS NULL OR ee.nota_fiscal_numero ILIKE '%' || sqlc.narg('nota_fiscal') || '%') -- Busca por NF
ORDER BY ee.data_entrada DESC
LIMIT $1 OFFSET $2;

-- name: CancelarEntrada :execrows
-- name: CancelarEntrada :execrows
UPDATE entrada_epi 
SET 
    cancelada_em = NOW(), 
    ativo = FALSE 
WHERE id = $1 
  AND cancelada_em IS NULL 
  AND quantidadeAtual = quantidade; -- Garante que nada foi usado

-- name: ContarEntradas :one
-- name: ContarEntradasFiltradas :one
SELECT COUNT(*) 
FROM entrada_epi ee
WHERE 
    (
        (@canceladas::boolean IS FALSE AND ee.cancelada_em IS NULL) OR
        (@canceladas::boolean IS TRUE AND ee.cancelada_em IS NOT NULL)
    )
    AND (sqlc.narg('id_epi')::int IS NULL OR ee.IdEpi = sqlc.narg('id_epi'))
    AND (sqlc.narg('id_entrada')::int IS NULL OR ee.id = sqlc.narg('id_entrada'))
    AND (sqlc.narg('data_inicio')::date IS NULL OR ee.data_entrada >= sqlc.narg('data_inicio'))
    AND (sqlc.narg('data_fim')::date IS NULL OR ee.data_entrada <= sqlc.narg('data_fim'))
    AND (sqlc.narg('nota_fiscal')::text IS NULL OR ee.nota_fiscal_numero ILIKE '%' || sqlc.narg('nota_fiscal') || '%');-- adicione os mesmos filtros aqui