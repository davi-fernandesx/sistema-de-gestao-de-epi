package login

import (
	"context"
	"database/sql"

	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

// interface contendo os metados
type LoginRepository interface {
	AddLogin(ctx context.Context, model *model.Login) error
	DeletarLogin(ctx context.Context, id int) error
	BuscaPorNome(ctx context.Context, nome string) (*model.Login, error)
}

type ConnDB struct {
	db *sql.DB
}

// construtor
func NewLogin(DB *sql.DB) LoginRepository {

	return &ConnDB{
		db: DB,
	}
}

// AddLogin implements loginRepository.
// função para adicionar um login no sistema
func (s *ConnDB) AddLogin(ctx context.Context, model *model.Login) error {

	query := `
			INSERT INTO login (usuario, senha) values (@p1, @p2);

	`
	_, err := s.db.ExecContext(ctx, query, model.Nome, model.Senha)
	if err != nil {
		if helper.IsUniqueViolation(err){

			return fmt.Errorf("usuario com nome: %s ja existente. %w", model.Nome, Errors.ErrSalvar)
		}

		return fmt.Errorf("erro inesperado ao salvar Login: %w", Errors.ErrInternal)
	}

	return nil
}

// DeletarLogin implements loginRepository.
func (s *ConnDB) DeletarLogin(ctx context.Context, id int) error {

	query := `

	delete from login where id = @id
	`

	result, err := s.db.Exec(query, sql.Named("id", id))
	if err != nil {

		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", Errors.ErrLinhasAfetadas)

	}

	if row == 0 {
		return fmt.Errorf("usuario com id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil

}

// busca o usuario pelo nome
func (s *ConnDB) BuscaPorNome(ctx context.Context, nome string) (*model.Login, error) {

	query := `
		select usuario, senha from login
		where usuario = @usuario ;

	`

	var usuario model.Login

	err := s.db.QueryRow(query, sql.Named("usuario", nome)).Scan(
		&usuario.Nome,
		&usuario.Senha,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario com nome %s, não encontrado! %w", nome, Errors.ErrNaoEncontrado)
		}

		return nil, fmt.Errorf("erro ao escanecar dados!, %w", Errors.ErrFalhaAoEscanearDados)
	}

	return &usuario, nil

}

// Login implements loginRepository.
