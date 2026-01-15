-- name: AddDevolucaoSimples :exec
INSERT INTO devolucao (
    IdFuncionario, IdEpi, IdMotivo, data_devolucao, IdTamanho, 
    quantidadeAdevolver, assinatura_digital
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: AddTrocaEpi :one
INSERT INTO devolucao (
    IdFuncionario, IdEpi, IdMotivo, data_devolucao, IdTamanho, 
    quantidadeAdevolver, IdEpiNovo, IdTamanhoNovo, quantidadeNova, assinatura_digital
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: AddEntregaVinculada :one
INSERT INTO entrega_epi (IdFuncionario, data_entrega, assinatura, IdTroca)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListarDevolucoes :many
SELECT 
    d.id, d.IdFuncionario, f.nome as func_nome, f.matricula,
    f.IdDepartamento, dd.nome as dep_nome,
    f.IdFuncao, ff.nome as funcao_nome,
    d.IdEpi, e.nome as epi_antigo_nome, e.fabricante as epi_antigo_fab, e.CA as epi_antigo_ca,
    d.IdTamanho as tam_antigo_id, t.tamanho as tam_antigo_nome,
    d.quantidadeAdevolver, d.IdMotivo, m.motivo as motivo_nome,
    d.IdEpiNovo, en.nome as epi_novo_nome, en.fabricante as epi_novo_fab, en.CA as epi_novo_ca,
    d.quantidadeNova, d.IdTamanhoNovo, tn.tamanho as tam_novo_nome,
    d.assinatura_digital, d.data_devolucao
FROM devolucao d
INNER JOIN epi e ON d.IdEpi = e.id
INNER JOIN funcionario f ON d.IdFuncionario = f.id	
INNER JOIN departamento dd ON f.IdDepartamento = dd.id
INNER JOIN funcao ff ON f.IdFuncao = ff.id
INNER JOIN tamanho t ON d.IdTamanho = t.id
INNER JOIN motivo_devolucao m ON d.IdMotivo = m.id
LEFT JOIN epi en ON d.IdEpiNovo = en.id
LEFT JOIN tamanho tn ON d.IdTamanhoNovo = tn.id
WHERE 
    ((sqlc.arg('canceladas')::boolean IS FALSE AND d.cancelada_em IS NULL) OR
     (sqlc.arg('canceladas')::boolean IS TRUE AND d.cancelada_em IS NOT NULL))
AND (sqlc.narg('id')::int IS NULL OR d.id = sqlc.narg('id'))
AND (sqlc.narg('matricula')::text IS NULL OR f.matricula = sqlc.narg('matricula'))
ORDER BY d.data_devolucao DESC;

-- name: CancelarDevolucao :execrows
UPDATE devolucao
SET cancelada_em = NOW(),
    ativo = FALSE
WHERE id = $1 AND cancelada_em IS NULL;