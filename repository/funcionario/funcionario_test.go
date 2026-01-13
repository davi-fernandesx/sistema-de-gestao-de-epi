package funcionario

import (
	"context"
	"errors"

	"fmt"
	"testing"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCriarFuncionario_Integracao(t *testing.T) {
	// 1. SETUP: Sobe o Docker e cria tabelas (Demora uns 5s na primeira vez)
	db := integracao.SetupTestDB(t)
	// O defer garante que o DB feche, o cleanup do setupTestDB mata o container
	defer db.Close()
	ctx := context.Background()
	repo := NewFuncionarioRepository(db)

	// 2. ARRANGE (Preparação dos dados)
	// Precisamos de um Depto e uma Funcao válidos, senão o SQL Server rejeita o Funcionário.
	idDepto := integracao.CreateDepartamento(t, db)
	idFuncao := integracao.CreateFuncao(t, db, idDepto)

	funcionarios := []model.FuncionarioINserir{

		{
			Nome:            "rada",
			Matricula:       "42343",
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncao,
		},
		{
			Nome:            "davi",
			Matricula:       "3345",
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncao,
		},
		{
			Nome:            "felipe",
			Matricula:       "12345",
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncao,
		},
		{
			Nome:            "rada",
			Matricula:       "4234",
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncao,
		},
	}

	t.Run("sucesso ao adicionar um funcionario no banco de dados", func(t *testing.T) {

		// 3. ACT (Execução)
		// Tenta salvar no banco de verdade
		err := repo.AddFuncionario(ctx, &funcionarios[0])

		// 4. ASSERT (Verificação)
		assert.NoError(t, err)

		// Prova real: Vamos buscar no banco pra ver se gravou mesmo
		var nomeGravado string
		var IdGravado int

		err = db.QueryRow("SELECT id, nome FROM funcionario WHERE matricula = '42343'").
			Scan(&IdGravado, &nomeGravado)

		assert.NoError(t, err)
		fmt.Println(err)
		assert.Equal(t, "rada", nomeGravado)

	})

	t.Run("erro de fk, id funcao ou departamento nao existe", func(t *testing.T) {

		idDepto := int64(9999)
		idFuncaoNaoExiste := int64(9999)

		funcComErro := &model.FuncionarioINserir{
			Nome:            "davi",
			Matricula:       "54545",
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncaoNaoExiste,
		}

		err := repo.AddFuncionario(ctx, funcComErro)
		// O Teste deve passar se der ERRO. Se não der erro, o banco está quebrado.
		assert.Error(t, err)
		fmt.Println(err)
		// O erro deve mencionar Foreign Key constraint
		assert.ErrorIs(t, err, Errors.ErrDadoIncompativel)

	})

	t.Run("matricula duplicada", func(t *testing.T) {

		idDepto := integracao.CreateDepartamento(t, db)
		idFuncao := integracao.CreateFuncao(t, db, idDepto)
		matriculaDuplicada := "42343"

		// Agora criamos o segundo com a MESMA matrícula
		funcionario2 := &model.FuncionarioINserir{
			Nome:            "Clone",
			Matricula:       matriculaDuplicada, // <--- O Culpado! Já existe.
			ID_departamento: &idDepto,
			ID_funcao:       &idFuncao,
		}

		// Act:
		err := repo.AddFuncionario(context.Background(), funcionario2)

		// Assert:
		assert.Error(t, err)

		fmt.Println(err)
		// 2. Verifica o tipo do erro
		assert.ErrorIs(t, err, Errors.ErrSalvar)
	})

	t.Run("sucesso ao buscar funcionario pela sua matricula", func(t *testing.T) {

		funcionario, err := repo.BuscaFuncionario(ctx, funcionarios[0].Matricula)
		require.NoError(t, err)
		require.NotNil(t, funcionario)
		require.NotEmpty(t, funcionario)
		require.Equal(t, &funcionario.ID_departamento, funcionarios[0].ID_departamento)

		t.Run("erro - matricula nao encontrada", func(t *testing.T) {

			funcionario, err := repo.BuscaFuncionario(ctx, "2133")
			require.Error(t, err)
			require.Nil(t, funcionario)
			require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))

		})

	})

	t.Run("sucesso ao buscar todos os funcionarios", func(t *testing.T) {

		_ = repo.AddFuncionario(ctx, &funcionarios[1])
		_ = repo.AddFuncionario(ctx, &funcionarios[2])
		_ = repo.AddFuncionario(ctx, &funcionarios[3])

		tests, err := repo.BuscarTodosFuncionarios(ctx)
		require.NoError(t, err)
		fmt.Println(err)
		require.Equal(t, funcionarios[1].Nome, tests[1].Nome)
		require.Equal(t, funcionarios[2].Nome, tests[2].Nome)
		require.Equal(t, funcionarios[3].Nome, tests[3].Nome)

	})

	t.Run("sucesso ao fazer um softdelete no funcionario", func(t *testing.T) {


		var IdGravado int

		err:= db.QueryRow("SELECT id  FROM funcionario WHERE matricula = '42343'").
			Scan(&IdGravado)

		assert.NoError(t, err)
		fmt.Println(err)
	
		err  = repo.DeletarFuncionario(ctx, IdGravado)
		require.NoError(t, err)

		t.Run("funcionario nao encontrado para deletar", func(t *testing.T) {

			err = repo.DeletarFuncionario(ctx, 75677)
			require.Error(t, err)
			require.True(t, errors.Is(err, Errors.ErrNaoEncontrado))

		})
	})

	t.Run("update do nome, departamento e funcao do funcionario", func(t *testing.T) {

		funcionario := model.Funcionario{
			Id:              1,
		
			ID_departamento: 1,
		
			ID_funcao:       1,
	
		}

		err := repo.UpdateFuncionarioNome(ctx, funcionario.Id, "juze")
		require.NoError(t, err)

		err = repo.UpdateFuncionarioDepartamento(ctx, funcionario.ID_departamento, "rh")
		require.NoError(t, err)

		err = repo.UpdateFuncionarioFuncao(ctx, funcionario.ID_funcao, "gestor")

	})

}
