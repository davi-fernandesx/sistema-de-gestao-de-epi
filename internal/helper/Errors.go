package helper

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// Definição de Erros Amigáveis (Negócio)
var (
	ErrId                  = errors.New("id invalido")
	ErrInternal            = errors.New("Erro Internal do banco de dados")
	ErrNaoEncontrado       = errors.New("registro não encontrado")
	ErrDadoDuplicado       = errors.New("este registro já existe no sistema")
	ErrConflitoIntegridade = errors.New("não é possível excluir: existem outros dados vinculados")
	ErrCampoObrigatorio    = errors.New("campo obrigatório não preenchido")
	ErrEstoqueInsuficiente = errors.New("estoque insuficiente para esta operação")
	ErrSessaoExpirada      = errors.New("sua sessão expirou, faça login novamente")
)

// Códigos de Erro Oficiais do PostgreSQL
const (
	PgUniqueViolation     = "23505" // Chave duplicada (Unique Constraint)
	PgForeignKeyViolation = "23503" // Chave estrangeira (ID que não existe ou está em uso)
	PgNotNullViolation    = "23502" // Tentativa de inserir nulo onde não pode
)

// TraduzErroPostgres converte erros técnicos do DB em erros de negócio do seu SaaS
func TraduzErroPostgres(err error) error {
	if err == nil {
		return nil
	}

	// Verifica se o erro veio do driver do Postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PgUniqueViolation:
			return ErrDadoDuplicado
		case PgForeignKeyViolation:
			return ErrConflitoIntegridade
		case PgNotNullViolation:
			return ErrCampoObrigatorio
		}
	}

	return fmt.Errorf("%w: %v", ErrInternal, err)
}
