package epi

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)


// interface sql error (para pegar o codigo do erro)

type EpiRepository struct {
	DB *sql.DB
}

func NewEpiRepository(db *sql.DB) *EpiRepository {

	return &EpiRepository{
		DB: db,
	}
}

// AddEpi implements EpiInterface.
func (n *EpiRepository) AddEpi(ctx context.Context, epi *model.EpiInserir) error {

	tx, err := n.DB.BeginTx(ctx, nil)
	if err != nil {
        return fmt.Errorf("erro ao iniciar transação: %w", Errors.ErrInternal)
    }
	defer tx.Rollback()
	query := `insert into epi (nome, fabricante, CA, descricao, validade_CA, IdTipoProtecao, alerta_minimo) 
			OUTPUT INSERTED.id 
			values 
			(@nome, @fabricante, @CA, @descricao,@validade_CA, @id_tipo_protecao, @alerta_minimo)`// quwry para
			//salvar um epi e retornar seu id

		
	var EpiId int64
	 err = tx.QueryRowContext(ctx, query,
		sql.Named("nome", epi.Nome),
		sql.Named("fabricante", epi.Fabricante),
		sql.Named("CA", epi.CA),
		sql.Named("descricao", epi.Descricao),
		sql.Named("validade_CA", epi.DataValidadeCa),
		sql.Named("id_tipo_protecao", epi.IDprotecao),
		sql.Named("alerta_minimo", epi.AlertaMinimo)).Scan(&EpiId)//escaneado o id

	if err != nil {
		if helper.IsUniqueViolation(err){

			return fmt.Errorf("CA %s ja existe no sistema, %w",epi.CA, Errors.ErrSalvar)
		}

		if helper.IsForeignKeyViolation(err){

			return fmt.Errorf("id protecao não existente no banco de dados, %w", Errors.ErrDadoIncompativel)
		}
		return fmt.Errorf(" Erro interno ao salvar Epi: %w",  Errors.ErrConexaoDb)
	}

	stmt, err:= tx.PrepareContext(ctx, "insert into tamanhos_epis(IdEpi, IdTamanho) values (@id_epi, @id_tamanho)") 
	//preparando a query para inserir na tabela de assosiação o id do epi salvar e o id do tamanho(vindo da model do epi)
	if err != nil{
		return fmt.Errorf("erro ao preparar statemnt para tamanho_epi: %w", Errors.ErrInternal)
	}
	defer stmt.Close()

	for _, idTamanho:= range epi.Idtamanho {
		_, err:= stmt.ExecContext(ctx, sql.Named("id_epi", EpiId), sql.Named("id_tamanho", idTamanho))
		//executando a query preparada, adicionando na tabela de associação os id do epi e o id dos tamanhos
		if err != nil {

				if helper.IsForeignKeyViolation(err){

					return fmt.Errorf("id tamanho nao existe no banco de dados, %w", Errors.ErrDadoIncompativel)
				}
				return fmt.Errorf("erro ao inserir na tabela tamanhos_epis para o tamanho ID %d: %w", idTamanho, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		 return fmt.Errorf("erro ao commitar transação: %w", err)
	}


	return nil
}

// BuscarEpi implements EpiInterface.
func (n *EpiRepository) BuscarEpi(ctx context.Context, id int) (*model.Epi, error) {

	query := `
		select
			e.id, e.nome, e.fabricante,e.CA, e.descricao,
			e.validade_CA, e.alerta_minimo, e.IdTipoProtecao, tp.nome as 'nome da protecao'
			from
				epi e
			inner join
				tipo_protecao tp on	e.IdTipoProtecao = tp.id		
			where
				e.id = @id and e.ativo = 1
	`

	var epi model.Epi
	err := n.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&epi.ID,
		&epi.Nome,
		&epi.Fabricante,
		&epi.CA,
		&epi.Descricao,
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
	
	select 	t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhos_epis te on t.id = te.IdTamanho
		where te.IdEpi= @epiId and te.ativo = 1

		`

	linhas, err:= n.DB.QueryContext(ctx, queryTamanhos, sql.Named("epiId", epi.ID))
	if err != nil {
        // Retorna o EPI encontrado, mas avisa sobre o erro nos tamanhos, ou pode retornar o erro direto
        return nil, fmt.Errorf("falha ao buscar tamanhos %w ", Errors.ErrFalhaAoEscanearDados)
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

	if len(tamanhos) == 0 {

		return nil, fmt.Errorf("tamanhos nao encontrados, %w", Errors.ErrNaoEncontrado)
	}
	epi.Tamanhos = tamanhos
	

	return &epi, nil
}

// BuscarTodosEpi implements EpiInterface.
func (n *EpiRepository) BuscarTodosEpi(ctx context.Context) ([]model.Epi, error) {

	query := `
		select e.id, e.nome, e.fabricante,e.CA, e.descricao,
		e.validade_CA, e.alerta_minimo, e.IdTipoProtecao, tp.nome as 'nome da protecao'
			from
				epi e
			inner join
				tipo_protecao tp on	e.IdTipoProtecao = tp.id
			where e.ativo = 1`

	linhas, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.Epi{}, fmt.Errorf("erro ao buscar todos os epi!, %w", Errors.ErrBuscarTodos)
	}
	defer linhas.Close()

	 EpisMap:= make(map[int]*model.Epi)
	 var ids []string //usando map ao inves de um slic por motivos de performace

	 for linhas.Next(){

		var epi model.Epi
		err:= linhas.Scan(
		&epi.ID,
		&epi.Nome,
		&epi.Fabricante,
		&epi.CA,
		&epi.Descricao,
		&epi.DataValidadeCa,
		&epi.AlertaMinimo,
		&epi.IDprotecao,
		&epi.NomeProtecao,
		)
		if err!= nil {
			return  nil, fmt.Errorf("%w", Errors.ErrFalhaAoEscanearDados)
		}

		EpisMap[epi.ID] = &epi
		ids = append(ids, strconv.Itoa(epi.ID))
	 }

	 if len(EpisMap) == 0 {

		return []model.Epi{}, nil
	 }

	 //segunda query

	 queryTamanhos:= fmt.Sprintf(` 
		select 
			te.IdEpi ,t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhos_epis te on t.id = te.IdTamanho
		where te.IdEpi IN (%s) and te.ativo = 1`, strings.Join(ids, ","))// query que retorna o id do epi, id do tamanho e o tamanho

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
func (n *EpiRepository) DeletarEpi(ctx context.Context, id int) error {


	tx, err:= n.DB.BeginTx(ctx, nil)
	if err!= nil {

		return  fmt.Errorf("erro ao começar transaction, %w", Errors.ErrInternal)
	}
	defer tx.Rollback()

	queryTamanhoEpi := ` update tamanhos_epis 
							set ativo = 0, 
							deletado_em = getdate()
							where IdEpi = @id and ativo = 1`
	_, err = tx.ExecContext(ctx, queryTamanhoEpi, sql.Named("id", id))
	if err != nil {
		return  fmt.Errorf(" erro ao apagar tamanhos dos epis, %w",Errors.ErrAoapagar)
	}
		
	query:=  `update epi
				set ativo = 0,
				deletado_em = getdate()
			where id = @id and ativo = 1`
	result, err := tx.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return  fmt.Errorf("erro ao apagar epi, %w", Errors.ErrAoapagar)
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

func (n  *EpiRepository) UpdateEpiNome(ctx context.Context,id int, nome string)error {

	query:= `update epi
			set nome = @nome
			where id = @id and ativo = 1`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("nome", nome), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *EpiRepository) UpdateEpiFabricante(ctx context.Context,id int, fabricante string)error {

	query:= `update epi
			set fabricante = @fabricante
			where id = @id and ativo = 1`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("fabricante", fabricante), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}


func (n  *EpiRepository) UpdateEpiCa(ctx context.Context,id int, ca string)error {

	query:= `update epi
			set CA = @ca
			where id = @id and ativo = 1`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("ca", ca), sql.Named("id", id))
	if err != nil {
		if helper.IsUniqueViolation(err){

			return fmt.Errorf("CA ja existente no sistema, %w", Errors.ErrSalvar)
		}

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	
	}

	return nil
}

func (n  *EpiRepository) UpdateEpiDescricao(ctx context.Context,id int, descricao string)error {

	query:= `update epi
			set descricao = @descricao
			where id = @id and ativo = 1`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("descricao", descricao), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}

func (n  *EpiRepository) UpdateEpiDataValidadeCa(ctx context.Context,id int, dataValidadeCa time.Time)error {

	query:= `update epi
			set validadeCa = @dataValidadeCa
			where id = @id and ativo = 1`
	
	_, err:=n.DB.ExecContext(ctx, query, sql.Named("dataValidade", dataValidadeCa), sql.Named("id", id))
	if err != nil {

		return  fmt.Errorf("erro ao atualizar ca do epi, %w",Errors.ErrInternal)
	}

	return nil
}
