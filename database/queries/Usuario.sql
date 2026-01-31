-- name: CreateUser :exec 
insert into usuarios (nome, email, senha_hash) values ($1, $2, $3);


-- name: BuscarPorIdUsuario :one
SELECT id, nome, email, ativo
FROM usuarios
WHERE id = $1 AND ativo = TRUE
LIMIT 1;

-- name: BuscarTodosUsuarios :many
SELECT id, nome, email, ativo
FROM usuarios
WHERE ativo = TRUE;

-- name: DeletarUsuario :execrows
UPDATE usuarios
SET ativo = FALSE
WHERE id = $1 AND ativo = TRUE; 

-- name: BuscarUsuarioPorEmail :one
SELECT id,nome, email, senha_hash
FROM usuarios
WHERE email = $1 AND ativo = TRUE
LIMIT 1;