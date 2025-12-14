package epi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

type EpiInterface interface {
	AddEpi(ctx context.Context, epi *model.EpiInserir) error
	DeletarEpi(ctx context.Context, id int) error
	BuscarEpi(ctx context.Context, id int) (*model.Epi, error)
	BuscarTodosEpi(ctx context.Context) ([]model.Epi, error)
	UpdateEpiNome(ctx context.Context,id int, nome string)error
	UpdateEpiCa(ctx context.Context,id int, ca string)error
	UpdateEpiFabricante(ctx context.Context,id int, fabricante string)error
	UpdateEpiDescricao(ctx context.Context,id int, descricao string)error
	UpdateEpiDataFabricacao(ctx context.Context,id int, dataFabricacao time.Time)error
	UpdateEpiDataValidade(ctx context.Context,id int, dataValidade string)error
	UpdateEpiDataValidadeCa(ctx context.Context,id int, dataValidadeCa string)error
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
func (n *NewSqlLogin) AddEpi(ctx context.Context, epi *model.EpiInserir) error {

	tx, err := n.DB.BeginTx(ctx, nil)
	if err != nil {
        return fmt.Errorf("erro ao iniciar transação: %w", Errors.ErrInternal)
    }
	defer tx.Rollback()
	query := `insert into epi (nome, fabricante, CA, descricao, data_fabricacao, data_validade, validade_CA, id_tipo_protecao, alerta_minimo) 
			OUTPUT INSERTED.id 
			values 
			(@nome, @fabricante, @CA, @descricao,@data_fabricacao, @data_validade, @validade_CA, @id_tipo_protecao, @alerta_minimo )`// quwry para
			//salvar um epi e retornar seu id

		
	var EpiId int64
	 err = tx.QueryRowContext(ctx, query,
		sql.Named("nome", epi.Nome),
		sql.Named("fabricante", epi.Fabricante),
		sql.Named("CA", epi.CA),
		sql.Named("descricao", epi.Descricao),
		sql.Named("data_fabricacao", epi.DataFabricacao),
		sql.Named("data_validade", epi.DataValidade),
		sql.Named("validade_CA", epi.DataValidadeCa),
		sql.Named("id_tipo_protecao", epi.IDprotecao),
		sql.Named("alerta_minimo", epi.AlertaMinimo)).Scan(&EpiId)//escaneado o id

	if err != nil {
		return fmt.Errorf(" Erro interno ao salvar Epi: %w",  Errors.ErrInternal)
	}

	stmt, err:= tx.PrepareContext(ctx, "insert into tamanho_epi  (id_tamanho, id_epi) values (@id_tamanho, @id_epi)") 
	//preparando a query para inserir na tabela de assosiação o id do epi salvar e o id do tamanho(vindo da model do epi)
	if err != nil{
		return fmt.Errorf("erro ao preparar statemnt para tamanho_epi: %w", Errors.ErrInternal)
	}
	defer stmt.Close()

	for _, idTamanho:= range epi.Idtamanho {
		_, err:= stmt.ExecContext(ctx, sql.Named("id_tamanho", idTamanho), sql.Named("id_epi", EpiId))
		//executando a query preparada, adicionando na tabela de associação os id do epi e o id dos tamanhos
		if err != nil {

			return fmt.Errorf("erro ao inserir na tabela epi_tamanhos para o tamanho ID %d: %w", idTamanho, Errors.ErrInternal)
		}
	}

	err = tx.Commit()
	if err != nil {
		 return fmt.Errorf("erro ao commitar transação: %w", Errors.ErrInternal)
	}


	return nil
}

// BuscarEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {

	query := `
			select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id		
			where
				e.id = @id
	`

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
		&epi.AlertaMinimo,
		&epi.IDprotecao,
		&epi.NomeProtecao,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("epi com id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
		}

		return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
	}
	//query usada para buscar os tamanhos por basee  do id do produto passado
	queryTamanhos:=`
	
		select 
			t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhosEpis te on t.id = te.id_tamanho
		where
			te.epiId = @epiId
		`

	linhas, err:= n.DB.QueryContext(ctx, queryTamanhos, sql.Named("epiId", epi.ID))
	if err != nil {
        // Retorna o EPI encontrado, mas avisa sobre o erro nos tamanhos, ou pode retornar o erro direto
        return nil, fmt.Errorf("erro ao buscar tamanhos para o epi %d: %w", epi.ID, Errors.ErrNaoEncontrado)
	}
	defer linhas.Close()

	var tamanhos []model.Tamanhos

	for linhas.Next(){

		var tamanho model.Tamanhos

		err:= linhas.Scan(&tamanho.ID, &tamanho.Tamanho)
		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)

		}

		tamanhos = append(tamanhos, tamanho)
	}

	epi.Tamanhos = tamanhos
	

	return &epi, nil
}

// BuscarTodosEpi implements EpiInterface.
func (n *NewSqlLogin) BuscarTodosEpi(ctx context.Context) ([]model.Epi, error) {

	query := `
				select
					e.id, e.nome, e.fabricante,e.CA, e.descricao, e.data_fabricacao, e.data_validade, 
					e.validade_CA, e.alerta_minimo, e.id_tipo_protecao, tp.nome
			from
				epi e
			inner join
				tipo_protecao tp on	e.id_tipo_protecao = tp.id`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Epi{}, fmt.Errorf("erro ao buscar todos os epi!, %w", Errors.ErrBuscarTodos)
	}
	defer linhas.Close()

	 EpisMap:= make(map[int]*model.Epi) //usando map ao inves de um slic por motivos de performace

	 for linhas.Next(){

		var epi model.Epi
		err:= linhas.Scan(
		&epi.ID,
		&epi.Nome,
		&epi.Fabricante,
		&epi.CA,
		&epi.Descricao,
		&epi.DataFabricacao,
		&epi.DataValidade,
		&epi.DataValidadeCa,
		&epi.AlertaMinimo,
		&epi.IDprotecao,
		&epi.NomeProtecao,
		)
		if err!= nil {
			return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		EpisMap[epi.ID] = &epi
	 }

	 if len(EpisMap) == 0 {

		return []model.Epi{}, nil
	 }

	 //segunda query

	 queryTamanhos:= ` 
	 		select
	 			 te.id_epi, t.id, t.tamanho
			from tamanhos t
			inner join tamanho_epi te on t.id = te.id_tamanho
	 `// query que retorna o id do epi, id do tamanho e o tamanho

	 linhasTamanhos, err:= n.DB.QueryContext(ctx, queryTamanhos)
	 if err != nil{
			return  nil, fmt.Errorf("erro ao buscar associacao de tamanho, %w", Errors.ErrInternal)		
	 }
	 defer linhasTamanhos.Close()

	 for linhasTamanhos.Next(){

		var epiID int64
		var t model.Tamanhos

		err:= linhasTamanhos.Scan(&epiID, &t.ID, &t.Tamanho)
		if err != nil {
			return nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)

		}

		if epi, ok:= EpisMap[int(epiID)]; ok {
			epi.Tamanhos = append(epi.Tamanhos, t)
		}

	 }
	 epis:= make([]model.Epi, 0, len(EpisMap))
	 for _, epi:= range EpisMap {
		epis = append(epis, *epi)
	 }



	return epis, nil

}

// DeletarEpi implements EpiInterface.
func (n *NewSqlLogin) DeletarEpi(ctx context.Context, id int) error {


	tx, err:= n.DB.BeginTx(ctx, nil)
	if err!= nil {

		return  fmt.Errorf("erro ao começar transaction, %w", Errors.ErrInternal)
	}
	defer tx.Rollback()

	queryTamanhoEpi := ` delete from  tamanho_epi where id_epi = @id`
	_, err = tx.ExecContext(ctx, queryTamanhoEpi, sql.Named("id", id))
	if err != nil {
		return err
	}
		
	query:=  `delete from epi where id = @id`
	result, err := tx.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  err
	}

	linhas, err:= result.RowsAffected()
	if err != nil{

			return fmt.Errorf("erro ao verificar linha afetadas, %w", Errors.ErrLinhasAfetadas)
		
	}

	if linhas == 0 {

		return  fmt.Errorf("epi com o id %d não encontrado!, %w", id, Errors.ErrNaoEncontrado)
	}

	return  tx.Commit()
}

func (n  *NewSqlLogin) UpdateEpiNome(ctx context.Context,id int, nome string)error {

	query:= `update epi
			set nome = @nome
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("nome", nome), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *NewSqlLogin) UpdateEpiFabricante(ctx context.Context,id int, fabricante string)error {

	query:= `update epi
			set fabricante = @fabricante
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("fabricante", fabricante), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}



func (n  *NewSqlLogin) UpdateEpiCa(ctx context.Context,id int, ca string)error {

	query:= `update epi
			set CA = @ca
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("ca", ca), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *NewSqlLogin) UpdateEpiDescricao(ctx context.Context,id int, descricao string)error {

	query:= `update epi
			set descricao = @descricao
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("descricao", descricao), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}


func (n  *NewSqlLogin) UpdateEpiDataFabricacao(ctx context.Context,id int, dataFabricacao time.Time)error {

	query:= `update epi
			set dataFabricacao = @dataFabricacao
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("ca", dataFabricacao), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *NewSqlLogin) UpdateEpiDataValidade(ctx context.Context,id int, dataValidade string)error {

	query:= `update epi
			set dataValidade = @dataValidade
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("dataValidada", dataValidade), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *NewSqlLogin) UpdateEpiDataValidadeCa(ctx context.Context,id int, dataValidadeCa string)error {

	query:= `update epi
			set validadeCa = @dataValidadeCa
			where id = @id`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("dataValidade", dataValidadeCa), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}
