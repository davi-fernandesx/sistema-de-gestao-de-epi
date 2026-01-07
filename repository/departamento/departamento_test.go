package departamento

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
	repo := NewDepartamentoRepository(db)

	departamentos := []model.Departamento{
		{

			Departamento: "ti",
		},
		{

			Departamento: "rh",
		},
		{

			Departamento: "adm",
		},
	}

	t.Run("sucesso ao adicionar um departamento", func(t *testing.T) {

		err := repo.AddDepartamento(ctx, &departamentos[0])
		require.NoError(t, err)

		var dep string
		var id int

		err = db.QueryRow("SELECT id, nome FROM departamento WHERE id = 1").
			Scan(&id, &dep)

		require.NoError(t, err)
		fmt.Println(err)
		require.Equal(t, "ti", dep)
	})

	t.Run("erro - adicionar um departamento repetido", func(t *testing.T) {

		dep := model.Departamento{

			Departamento: "ti",
		}

		err := repo.AddDepartamento(ctx, &dep)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrSalvar))
	})

	t.Run("sucesso - buscar um departamento", func(t *testing.T) {

		dep, err:= repo.BuscarDepartamento(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, dep)
		require.NotEmpty(t, dep)
		require.Equal(t, dep.ID, 1)
		  
	})

	t.Run("departamento não encontrado", func(t *testing.T) {

		dep, err:= repo.BuscarDepartamento(ctx, 77)
		require.Error(t, err)
		require.Nil(t,dep)
		require.True(t, errors.Is(err,Errors.ErrNaoEncontrado))
	})

	t.Run("buscar todos os departamentos", func(t *testing.T) {

		_ = repo.AddDepartamento(ctx, &departamentos[1])
		_ = repo.AddDepartamento(ctx, &departamentos[2])


		deps, err:= repo.BuscarTodosDepartamentos(ctx)
		require.NoError(t,err)
		require.Equal(t, departamentos[1].Departamento, deps[1].Departamento)
		require.Equal(t, departamentos[2].Departamento, deps[2].Departamento)

	})

	t.Run("sucesso ao fazer um softdelete no departamento", func(t *testing.T) {

		err:= repo.DeletarDepartamento(ctx, 1)
		require.NoError(t,err)
	})

	t.Run("departamento não encontrado", func(t *testing.T) {

		err:= repo.DeletarDepartamento(ctx, 111)
		require.Error(t, err)
		require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))
	})

}
