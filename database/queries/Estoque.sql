-- name: ListarLotesParaConsumo :many
-- O PostgreSQL usa FOR UPDATE para o que o SQL Server chama de UPDLOCK.
SELECT id, quantidadeAtual, data_validade, valor_unitario 
FROM entrada_epi 
WHERE IdEpi = $1 
  AND IdTamanho = $2 
  AND quantidadeAtual > 0 
  AND data_validade >= CURRENT_DATE
  AND ativo = TRUE
ORDER BY data_validade ASC
FOR UPDATE;

-- name: AbaterEstoqueLote :execrows
UPDATE entrada_epi 
SET quantidadeAtual = quantidadeAtual - $1 
WHERE id = $2 
  AND ativo = TRUE
  AND quantidadeAtual >= $1;


-- name: ReporEstoqueLote :execrows
UPDATE entrada_epi 
SET quantidadeAtual = quantidadeAtual + $1 
WHERE id = $2 
  AND ativo = TRUE
  AND quantidadeAtual >= $1;

-- name: RegistrarItemEntrega :exec
INSERT INTO epis_entregues (IdEpi, IdTamanho, quantidade, IdEntrega, IdEntrada, valor_unitario) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DevolverItemAoEstoque :exec
UPDATE entrada_epi eed
SET eed.quantidadeAtual = eed.quantidadeAtual + $3
WHERE id = (
    -- Subselect para achar o ID da entrada mais recente
    SELECT ee.id
    FROM entrada_epi ee
    WHERE ee.IdEpi = $1 
      AND ee.IdTamanho = $2
    ORDER BY ee.data_entrada DESC -- Ordena pela data (mais nova primeiro)
    LIMIT 1 -- Pega apenas uma
);  