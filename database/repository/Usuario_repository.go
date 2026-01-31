package repository

import (
	"context"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)


type UsuarioRepository struct {

	q *Queries
	db *pgxpool.Pool
}

func NewUsuarioRepository(pool *pgxpool.Pool) *UsuarioRepository {

	return &UsuarioRepository{q: New(pool), db: pool}
}

func (u *UsuarioRepository) Cadastrar(ctx context.Context, user CreateUserParams) error {

	err:= u.q.CreateUser(ctx, user)
	if err != nil {

		return helper.TraduzErroPostgres(err)
	}

	return nil
}

func (u *UsuarioRepository) Listar(ctx context.Context) ([]BuscarTodosUsuariosRow, error){

	usuarios, err:= u.q.BuscarTodosUsuarios(ctx)
	if err != nil {
		return []BuscarTodosUsuariosRow{}, helper.TraduzErroPostgres(err)
	}

	return usuarios, err
}

func(u *UsuarioRepository) BuscarPorEmail(ctx context.Context, email string) (BuscarUsuarioPorEmailRow, error){

	usuario, err:= u.q.BuscarUsuarioPorEmail(ctx, email)
	if err != nil {
		return BuscarUsuarioPorEmailRow{}, helper.TraduzErroPostgres(err)
	}

	return usuario, nil
}

func (u *UsuarioRepository) BuscarPoId(ctx context.Context, id int) (BuscarPorIdUsuarioRow, error){

	usuario, err:= u.q.BuscarPorIdUsuario(ctx, int32(id))
	if err != nil {

		return  BuscarPorIdUsuarioRow{}, helper.TraduzErroPostgres(err)
	}

	return usuario, nil
}