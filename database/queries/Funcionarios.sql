-- name: AddFuncionario :exec
INSERT INTO funcionario (tenant_id, nome, matricula, IdDepartamento, IdFuncao) 
VALUES ($1, $2, $3, $4, $5);

-- name: BuscaFuncionario :one
SELECT 
    fn.id, 
    fn.nome, 
    fn.matricula, 
    fn.IdDepartamento, 
    d.nome as departamento_nome,
    fn.IdFuncao, 
    f.nome as funcao_nome
FROM funcionario fn
INNER JOIN departamento d ON fn.IdDepartamento = d.id
INNER JOIN funcao f ON fn.IdFuncao = f.id
WHERE fn.matricula = $1 
  AND fn.tenant_id = $2 -- IMPORTANTE: Matrícula só é única dentro do tenant
  AND fn.ativo = TRUE;

-- name: BuscarTodosFuncionarios :many
SELECT 
    fn.id, 
    fn.nome, 
    fn.matricula, 
    fn.IdDepartamento, 
    d.nome as departamento_nome,
    fn.IdFuncao, 
    f.nome as funcao_nome
FROM funcionario fn
INNER JOIN departamento d ON fn.IdDepartamento = d.id
INNER JOIN funcao f ON fn.IdFuncao = f.id
WHERE fn.tenant_id = $1 -- SEGURANÇA: Só lista funcionários da empresa atual
  AND fn.ativo = TRUE;

-- name: DeletarFuncionario :execrows
UPDATE funcionario
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 
  AND tenant_id = $2 -- SEGURANÇA
  AND ativo = TRUE;

-- name: UpdateFuncionarioNome :execrows
UPDATE funcionario
SET nome = $2
WHERE id = $1 
  AND tenant_id = $3 -- SEGURANÇA
  AND ativo = TRUE;

-- name: UpdateFuncionarioDepartamento :execrows
UPDATE funcionario
SET IdDepartamento = $2
WHERE id = $1 
  AND tenant_id = $3 -- SEGURANÇA
  AND ativo = TRUE;

-- name: UpdateFuncionarioFuncao :execrows
UPDATE funcionario
SET IdFuncao = $2
WHERE id = $1 
  AND tenant_id = $3 -- SEGURANÇA
  AND ativo = TRUE;

-- name: BuscaFuncionarioPorId :one
SELECT 
    fn.id, 
    fn.nome, 
    fn.matricula, 
    fn.IdDepartamento, 
    d.nome as departamento_nome,
    fn.IdFuncao, 
    f.nome as funcao_nome
FROM funcionario fn
INNER JOIN departamento d ON fn.IdDepartamento = d.id
INNER JOIN funcao f ON fn.IdFuncao = f.id
WHERE fn.id = $1 
  AND fn.tenant_id = $2 -- SEGURANÇA
  AND fn.ativo = TRUE;