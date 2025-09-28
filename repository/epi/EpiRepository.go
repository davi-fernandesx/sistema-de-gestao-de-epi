package epi

import (
	"context"
	"database/sql"


	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository"
)

type EpiInterface interface {
	AddEpi(ctx context.Context, epi *model.Epi) error
	DeletarEpi(ctx context.Context, id int) error
	BuscarEpi(ctx context.Context, id int) (*model.Epi, error)
	BuscarTodosEpi(ctx context.Context) ([]model.Epi, error)
}

type NewSqlLogin struct {
	DB *sql.DB
}

func NewEpiRepository(db *sql.DB) EpiInterface {

	return &NewSqlLogin{
		DB: db,
	}
}

// AddEpi implements EpiInterface.
func (n *NewSqlLogin) AddEpi(ctx context.Context, epi *model.Epi) error {

	query := `insert into epi (nome, fabricante, CA, descricao, data_fabricacao, data_validade, validade_CA, id_tipo_protecao, alerta_minimo) values (
			@nome, @fabricante, @CA, @descricao,@data_fabricacao, @data_validade, @validade_CA, @id_tipo_protecao, @alerta_minimo )`

	_, err := n.DB.ExecContext(ctx, query,
		sql.Named("nome", epi.Nome),
		sql.Named("fabricantte", epi.Fabricante),
		sql.Named("CA", epi.CA),
		sql.Named("descricao", epi.Descricao),
		sql.Named("data_fabricacao", epi.DataFabricacao),
		sql.Named("data_validade", epi.DataValidade),
		sql.Named("validade_CA", epi.DataValidadeCa),
		sql.Named("id_tipo_protecao", epi.IDprotecao),
		sql.Named("alerta_minimo", epi.AlertaMinimo))

	if err != nil {
		return repository.ErrEpiAoAdicionarEpi
	}

	return nil
}

// BuscarEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {

	query := `select id, nome, fabricante, CA, descricao, data_fabricacao, data_validade, validade_CA, id_tipo_protecao, alerta_minimo
			from epi where id = @id`

	var epi model.Epi
	err := n.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&epi.ID,
		&epi.Nome,
		&epi.Fabricante,
		&epi.CA,
		&epi.Descricao,
		&epi.DataFabricacao,
		&epi.DataValidade,
		&epi.DataValidadeCa,
		&epi.IDprotecao,
		&epi.AlertaMinimo,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, repository.ErrAoProcurarEpi
		}

		return nil, repository.ErrFalhaAoEscanearDados
	}

	return &epi, nil
}

// BuscarTodosEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarTodosEpi(ctx context.Context) ([]model.Epi, error) {

	query := `select id, nome, fabricante, CA, descricao,
	 		data_fabricacao, data_validade, validade_CA, 
	 		id_tipo_protecao, alerta_minimo
			from epi`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Epi{}, repository.ErrAoBuscarTodosOsEpis
	}
	defer linhas.Close()

	var epis []model.Epi

	for linhas.Next() {

		var epi model.Epi

		if err := linhas.Scan(&epi.ID, &epi.Nome, &epi.Fabricante, &epi.CA, &epi.Descricao, &epi.DataFabricacao, &epi.DataValidade, &epi.DataValidadeCa,
			&epi.IDprotecao, &epi.AlertaMinimo); err != nil {

			return nil, repository.ErrFalhaAoEscanearDados
		}

		epis = append(epis, epi)
	}

	if err := linhas.Err(); err != nil {

		return nil, repository.ErrAoInterarSobreEpis
	}

	return epis, nil

}

// DeletarEpi implements EpiInterface.
func (n *NewSqlLogin) DeletarEpi(ctx context.Context, id int) error {
	
	query:=  `delete from epi where id = @id`

	result, err:= n.DB.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil{

		return  repository.ErrLinhasAfetadas
	}

	if linhas == 0 {

		return  repository.ErrEpiNaoEncontrado
	}

	return  nil
}
