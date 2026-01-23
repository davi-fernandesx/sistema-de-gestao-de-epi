-- name: CriaDepartamento :exec
insert into departamento (nome) 
values ($1);

-- name: BuscarDepartamento :one
SELECT id, nome 
FROM departamento 
WHERE id = $1 AND ativo = TRUE 
LIMIT 1;

-- name: BuscarTodosDepartamentos :many
SELECT id, nome 
FROM departamento 
WHERE ativo = TRUE;


-- name: DeletarDepartamento :execrows
UPDATE departamento
SET ativo = FALSE,
    deletado_em = NOW()
WHERE id = $1 AND ativo = TRUE;

-- name: UpdateDepartamento :execrows
UPDATE departamento
SET nome = $2
WHERE id = $1 AND ativo = TRUE;