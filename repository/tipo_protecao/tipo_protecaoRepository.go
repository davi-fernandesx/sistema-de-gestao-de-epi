package tipoprotecao

import (
	"context"
	"database/sql"
	"fmt"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type TipoProtecaoInterface interface {
	AddProtecao(ctx context.Context, protecao *model.TipoProtecao) error
	DeletarProtecao(ctx context.Context, ind int) error
	BuscarProtecao(ctx context.Context, id int) (*model.TipoProtecao, error)
	BuscarTodasProtecao(ctx context.Context) ([]model.TipoProtecao, error)
}

type SqlServerLogin struct {
	DB *sql.DB
}

func NewTipoProtecaoRepository(db *sql.DB) TipoProtecaoInterface {

	return &SqlServerLogin{
		DB: db,
	}
}

// AddProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) AddProtecao(ctx context.Context, protecao *model.TipoProtecao) error {
	
	query:= `insert into tipo_protecao(nome) values (@protecao)`

	_, err:= s.DB.ExecContext(ctx, query, sql.Named("protecao", protecao.Nome))
	if err != nil {
		return fmt.Errorf("erro interno ao salvar protecao, %w", Errors.ErrSalvar)
	}

	return  nil

}

// BuscarProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) BuscarProtecao(ctx context.Context, id int) (*model.TipoProtecao, error) {
	
	query:= `select id, nome from tipo_protecao where id = @id and ativo = 1`

	var protecao model.TipoProtecao

	err:= s.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(&protecao.ID, &protecao.Nome)

	if err != nil {
		if err == sql.ErrNoRows {
			return  nil,  fmt.Errorf("protecao com id %d, não encontrado! %w",id,  Errors.ErrNaoEncontrado)
		}

		return nil,fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}

	return &protecao, nil
}

// BuscarTodasProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) BuscarTodasProtecao(ctx context.Context) ([]model.TipoProtecao, error) {
	
	query:= `select id, nome from tipo_protecao where ativo = 1 `

	linhas, err:= s.DB.QueryContext(ctx, query)
	if err != nil {
		return  []model.TipoProtecao{}, fmt.Errorf("erro ao procurar todas as proteções, %w", Errors.ErrBuscarTodos)
	}

	defer linhas.Close()

	var protecoes []model.TipoProtecao

	for linhas.Next(){
		
		var protecao model.TipoProtecao
		err:= linhas.Scan(&protecao.ID, &protecao.Nome)
		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		protecoes = append(protecoes, protecao)
	}

	err = linhas.Err()
	if err != nil {
		return  nil, fmt.Errorf("erro ao iterar sobre as proteções , %w", Errors.ErrAoIterar)
	}

	return protecoes, nil
}

// DeletarProtecao implements TipoProtecaoInterface.
func (s *SqlServerLogin) DeletarProtecao(ctx context.Context, id int) error {
	
	query:= `update tipo_protecao
			set ativo = 0, deletado_em = getdate()
			where id = @id and ativo = 1`

	result, err:= s.DB.ExecContext(ctx, query, sql.Named("id", id))
	
	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil {
			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0 {

	return  fmt.Errorf("proteção com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return nil
}


