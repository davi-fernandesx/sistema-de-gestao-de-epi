package funcionario

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	mssql "github.com/microsoft/go-mssqldb"

	"github.com/stretchr/testify/require"
)

func mock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context, error) {

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock, ctx, err

}

func Test_addFuncinario(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewFuncionarioRepository(db)
	id_dep := 1
	id_func := 2
	func1 := model.FuncionarioINserir{
		Nome:            "davi",
		Matricula:       "123234",
		ID_departamento: &id_dep,
		ID_funcao:       &id_func,
	}

	query := regexp.QuoteMeta(`insert into funcionario(nome, matricula, id_departamento, id_funcao) values( @nome, @matricula, @id_departamento, @id_funcao)`)

	t.Run("testando o sucesso ao adicionar um funcionario", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(func1.Nome, func1.Matricula, func1.ID_departamento, func1.ID_funcao).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddFuncionario(ctx, &func1)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("erro generico ao salvar um funcionario", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(func1.Nome, func1.Matricula, func1.ID_departamento, func1.ID_funcao).WillReturnError(Errors.ErrInternal)
		err := repo.AddFuncionario(ctx, &func1)

		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrInternal)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro ao adicionar um funcionario que ja existe no sistema", func(t *testing.T) {

		mssqlErr := &mssql.Error{Number: 2627}

		mock.ExpectExec(query).
			WithArgs(func1.Nome,func1.Matricula, func1.ID_departamento,func1.ID_funcao).
			WillReturnError(mssqlErr)

		err := repo.AddFuncionario(ctx, &func1)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrSalvar, "erro tem que ser do tipo salvar")
		require.NoError(t, mock.ExpectationsWereMet())

	})

}

func Test_BuscaFuncionario(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewFuncionarioRepository(db)

	query := regexp.QuoteMeta(`select fn.id, fn.nome,fn.matricula, fn.id_departamento, d.departamento, fn.id_funcao, f.funcao
			from funcionario fn
			inner join departamento d on fn.id_departamento = d.id
			inner jon funcao f on fn.id_funcao = f.funcao
			where fn.matricula = @matricula`)

	func1 := model.Funcionario{
		Id:              1,
		Nome:            "davi",
		Matricula:       4234535,
		ID_departamento: 1,
		Departamento:    "dessosa",
		ID_funcao:       2,
		Funcao:          "dessosador",
	}

	t.Run("testando o sucesso ao buscar um funcionarioo", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		}).AddRow(
			func1.Id,
			func1.Nome,
			func1.Matricula,
			func1.ID_departamento,
			func1.Departamento,
			func1.ID_funcao,
			func1.Funcao,
		)

		mock.ExpectQuery(query).WithArgs(func1.Matricula).WillReturnRows(linhas)

		funcionario, err := repo.BuscaFuncionario(ctx, func1.Matricula)
		require.NoError(t, err)
		require.NotNil(t, funcionario)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de funcionario nao encontrado", func(t *testing.T) {

		mock.ExpectQuery(query).WithArgs(func1.Matricula).WillReturnError(sql.ErrNoRows)

		funcionario, err := repo.BuscaFuncionario(ctx, func1.Matricula)

		require.Error(t, err)
		require.Nil(t, funcionario)
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de escanear dados", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		}).AddRow(
			"func1.Id",
			func1.Nome,
			func1.Matricula,
			func1.ID_departamento,
			func1.Departamento,
			func1.ID_funcao,
			func1.Funcao,
		)

		mock.ExpectQuery(query).WithArgs(func1.Matricula).WillReturnRows(linhas).WillReturnError(Errors.ErrFalhaAoEscanearDados)

		funcionario, err := repo.BuscaFuncionario(ctx, func1.Matricula)
		require.Error(t, err)
		require.Nil(t, funcionario)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func Test_BuscarTodosoFuncionarios(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewFuncionarioRepository(db)

	query := regexp.QuoteMeta(`select fn.id,fn.matricula, fn.nome, fn.id_departamento, d.departamento, fn.id_funcao, f.funcao
			from funcionario fn
			inner join departamento d on fn.id_departamento = d.id
			inner jon funcao f on fn.id_funcao = f.funcao`)

	f1 := model.Funcionario{
		Id:              1,
		Nome:            "teste",
		Matricula:       534564,
		ID_departamento: 1,
		Departamento:    "teste",
		ID_funcao:       3,
		Funcao:          "teste",
	}

	f2 := model.Funcionario{
		Id:              1,
		Nome:            "teste",
		Matricula:       52356,
		ID_departamento: 1,
		Departamento:    "teste",
		ID_funcao:       3,
		Funcao:          "teste",
	}

	t.Run("testando o sucesso ao trazer todos os funcionarios", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		}).AddRow(
			f1.Id,
			f1.Nome,
			f1.Matricula,
			f1.ID_departamento,
			f1.Departamento,
			f1.ID_funcao,
			f1.Funcao,
		).AddRow(
			f2.Id,
			f2.Nome,
			f2.Matricula,
			f2.ID_departamento,
			f2.Departamento,
			f2.ID_funcao,
			f2.Funcao,
		)

		mock.ExpectQuery(query).WillReturnRows(linhas)

		funcs, err := repo.BuscarTodosFuncionarios(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, funcs)
		require.NotNil(t, funcs)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o retorno vazio", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		})

		mock.ExpectQuery(query).WillReturnRows(linhas)

		funcs, err := repo.BuscarTodosFuncionarios(ctx)
		require.NoError(t, err)
		require.Empty(t, funcs)
		require.Len(t, funcs, 0)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("falha na execução da query", func(t *testing.T) {

		mock.ExpectQuery(query).WillReturnError(Errors.ErrBuscarTodos)

		funcs, err := repo.BuscarTodosFuncionarios(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrBuscarTodos)
		require.Empty(t, funcs)
		require.Nil(t, funcs)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de escanear os dados", func(t *testing.T) {
		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		}).AddRow(
			"f1.Id",
			f1.Nome,
			f1.Matricula,
			f1.ID_departamento,
			f1.Departamento,
			f1.ID_funcao,
			f1.Funcao,
		).AddRow(
			"f2.Id",
			f2.Nome,
			f2.Matricula,
			f2.ID_departamento,
			f2.Departamento,
			f2.ID_funcao,
			f2.Funcao,
		)

		mock.ExpectQuery(query).WillReturnRows(linhas)

		funcs, err := repo.BuscarTodosFuncionarios(ctx)
		require.Error(t, err)
		require.Nil(t, funcs)
		require.Empty(t, funcs)
		require.ErrorIs(t, err, Errors.ErrFalhaAoEscanearDados)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("testando o erro de iteracao das linhas", func(t *testing.T) {

		linhas := sqlmock.NewRows([]string{
			"id",
			"nome",
			"matricula",
			"id_departamento",
			"departamento",
			"id_funcao",
			"funcao",
		}).AddRow(
			f1.Id,
			f1.Nome,
			f1.Matricula,
			f1.ID_departamento,
			f1.Departamento,
			f1.ID_funcao,
			f1.Funcao,
		).AddRow(
			f2.Id,
			f2.Nome,
			f2.Matricula,
			f2.ID_departamento,
			f2.Departamento,
			f2.ID_funcao,
			f2.Funcao,
		).CloseError(errors.New("erro simulado"))

		mock.ExpectQuery(query).WillReturnRows(linhas)

		funcs, err := repo.BuscarTodosFuncionarios(ctx)
		require.Error(t, err)
		require.Nil(t, funcs)
		require.Empty(t, funcs)
		require.ErrorIs(t, err, Errors.ErrAoIterar)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_deletarFuncionario(t *testing.T) {

	db, mock, ctx, err := mock(t)
	require.NoError(t, err)
	defer db.Close()

	repo := NewFuncionarioRepository(db)

	matricula := 1
	query := regexp.QuoteMeta(`delete from funcionario where matricula = @matricula`)

	t.Run("sucesso ao deletar funcionario", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeletarFuncionario(ctx, matricula)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("testando o erro de funcionario nao encontrado", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(matricula).WillReturnResult(sqlmock.NewResult(0, 0))
		err := repo.DeletarFuncionario(ctx, matricula)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrNaoEncontrado)

	})

	t.Run("testando erro de verificar linhas", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(matricula).WillReturnResult(sqlmock.NewErrorResult(Errors.ErrLinhasAfetadas))
		err := repo.DeletarFuncionario(ctx, matricula)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrLinhasAfetadas)

	})

	t.Run("testando o erro ao deletar um funcionario", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(matricula).WillReturnError(errors.New("erro generico"))
		err := repo.DeletarFuncionario(ctx, matricula)
		require.Error(t, err)
		require.ErrorIs(t, err, Errors.ErrInternal)
	})

}
