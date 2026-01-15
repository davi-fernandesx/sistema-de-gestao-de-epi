-- name: AddFuncao :exec
INSERT INTO funcao (nome, IdDepartamento) 
VALUES ($1, $2);

-- name: BuscarFuncao :one
SELECT 
    f.id, 
    f.nome, 
    f.IdDepartamento, 
    d.nome as departamento_nome
FROM funcao f
INNER JOIN departamento d ON f.IdDepartamento = d.id
WHERE f.id = $1 AND f.ativo = TRUE;

-- name: BuscarTodasFuncoes :many
SELECT 
    f.id, 
    f.nome, 
    f.IdDepartamento, 
    d.nome as departamento_nome
FROM funcao f
INNER JOIN departamento d ON f.IdDepartamento = d.id
WHERE f.ativo = TRUE;

-- name: PossuiFuncionariosVinculados :one
SELECT EXISTS (
    SELECT 1 FROM funcionario 
    WHERE IdFuncao = $1 AND ativo = TRUE
);

-- name: DeletarFuncao :execrows
UPDATE funcao
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 AND ativo = TRUE;

-- name: UpdateFuncao :execrows
UPDATE funcao
SET nome = $2
WHERE id = $1 AND ativo = TRUE;