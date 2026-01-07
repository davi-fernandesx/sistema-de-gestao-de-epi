package funcao

import (
	"context"
	"errors"
	"fmt"
	"testing"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
	"github.com/stretchr/testify/require"
)

func TestDepartamento(t *testing.T) {

	db := integracao.SetupTestDB(t)
	// O defer garante que o DB feche, o cleanup do setupTestDB mata o container
	defer db.Close()
	ctx := context.Background()
	repo := NewfuncaoRepository(db)
	idDepto := integracao.CreateDepartamento(t, db)

	funcoes := []model.FuncaoInserir{
		{

			Funcao:         "analista",
			IdDepartamento: idDepto,
		},
		{
			Funcao:         "dev",
			IdDepartamento: idDepto,
		},
		{

			Funcao:         "gestor",
			IdDepartamento: idDepto,
		},
	}

	t.Run("sucesso ao adicionar um departamento", func(t *testing.T) {

		err := repo.AddFuncao(ctx, &funcoes[0])
		require.NoError(t, err)

		var funcao string
		var id int

		err = db.QueryRow("SELECT id, nome FROM funcao WHERE id = 1").
			Scan(&id, &funcao)

		require.NoError(t, err)
		fmt.Println(err)
		require.Equal(t, "analista", funcao)
	})

	t.Run("erro - adicionar um departamento repetido", func(t *testing.T) {

		funcFalse := model.FuncaoInserir{
			Funcao:         "analista",
			IdDepartamento: idDepto,
		}

		err := repo.AddFuncao(ctx, &funcFalse)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrSalvar))
	})

	t.Run("id departamento invalido", func(t *testing.T) {
		funcFalse := model.FuncaoInserir{
			Funcao:         "dono",
			IdDepartamento: 55,
		}

		err := repo.AddFuncao(ctx, &funcFalse)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrDadoIncompativel))

	})

	t.Run("sucesso - buscar uma funcao", func(t *testing.T) {

		funcao, err := repo.BuscarFuncao(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, funcao)
		require.NotEmpty(t, funcao)
		require.Equal(t, funcao.ID, 1)

	})

	t.Run("funcao não encontrado", func(t *testing.T) {

		dep, err := repo.BuscarFuncao(ctx, 88)
		require.Error(t, err)
		require.Nil(t, dep)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))
	})

	t.Run("buscar todos os departamentos", func(t *testing.T) {

		_ = repo.AddFuncao(ctx, &funcoes[1])
		_ = repo.AddFuncao(ctx, &funcoes[2])

		funcc, err := repo.BuscarTodasFuncao(ctx)
		require.NoError(t, err)
		require.Equal(t, funcc[1].Funcao, funcoes[1].Funcao)
		require.Equal(t, funcc[2].Funcao, funcoes[2].Funcao)

	})

	t.Run("sucesso ao fazer um softdelete no departamento", func(t *testing.T) {

		err := repo.DeletarFuncao(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("departamento não encontrado", func(t *testing.T) {

		err := repo.DeletarFuncao(ctx, 111)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))
	})

}
