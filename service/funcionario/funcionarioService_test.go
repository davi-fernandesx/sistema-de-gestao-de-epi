package funcionario

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockFuncionarioRepo struct {
	mock.Mock
}

func (m *MockFuncionarioRepo) AddFuncionario(ctx context.Context, Funcionario *model.FuncionarioINserir) error {

	args := m.Called(ctx, Funcionario)

	return args.Error(0)
}

func (m *MockFuncionarioRepo) BuscaFuncionario(ctx context.Context, matricula int) (*model.Funcionario, error) {

	args := m.Called(ctx, matricula) //passando os argumentos

	var f *model.Funcionario
	if args.Get(0) != nil {
		f = args.Get(0).(*model.Funcionario)
	}

	return f, args.Error(1)
}

func (m *MockFuncionarioRepo) BuscarTodosFuncionarios(ctx context.Context) ([]model.Funcionario, error) {

	args := m.Called(ctx)

	var funcs []model.Funcionario

	if args.Get(0) != nil {

		funcs = args.Get(0).([]model.Funcionario)
	}
	return funcs, args.Error(1)
}

func (m *MockFuncionarioRepo) DeletarFuncionario(ctx context.Context, matricula int) error {
	args := m.Called(ctx, matricula)

	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioNome(ctx context.Context, id int, funcionario string) error {

	args := m.Called(ctx, id, funcionario)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioDepartamento(ctx context.Context, id int, idDepartamento string) error {

	args := m.Called(ctx, id, idDepartamento)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}

func (m *MockFuncionarioRepo) UpdateFuncionarioFuncao(ctx context.Context, id int, idFuncao string) error {

	args := m.Called(ctx, id, idFuncao)

	// Retorna o primeiro valor (índice 0) como um erro
	return args.Error(0)
}
func TestSalvarFuncionario(t *testing.T) {

	idf, idd := 3, 4
	funcionarioInputTeste := model.FuncionarioINserir{
		Nome:            " rahil ", //espaço a mais para testar o errro
		Matricula:       "546373",
		ID_departamento: &idd,
		ID_funcao:       &idf,
	}

	funcionarioInput := model.FuncionarioINserir{
		Nome:            "rahil",
		Matricula:       "546373",
		ID_departamento: &idd,
		ID_funcao:       &idf,
	}

	CasoDeTestes := []struct {
		nome         string
		ctx          context.Context
		modelAtestar model.FuncionarioINserir
		mock         func(mockRepo *MockFuncionarioRepo) //funcao do mock
		erro         error
	}{
		{
			nome:         "funcionario salvo",
			ctx:          context.Background(),
			modelAtestar: funcionarioInputTeste,
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscaFuncionario", mock.Anything, 546373). //erro esparado
											Return(nil, sql.ErrNoRows).
											Once()

				mockRepo.On("AddFuncionario", mock.Anything, &funcionarioInput). //salvando o funcionario
													Return(nil).
													Once()

			},
			erro: nil,
		},
		{
			nome:         "matricula ja cadastrada",
			ctx:          context.Background(),
			modelAtestar: funcionarioInputTeste,
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscaFuncionario", mock.Anything, 546373).
					Return(&model.Funcionario{}, nil).
					Once()
			},
			erro: errors.New("matricula ja cadastrada"),
		},
		{

			nome: "context cancelado",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			modelAtestar: funcionarioInputTeste,
			mock:         func(mockRepo *MockFuncionarioRepo) {},
			erro:         context.Canceled,
		},
		{
			nome:         "falha ao buscar o funcionario",
			ctx:          context.Background(),
			modelAtestar: funcionarioInputTeste,
			mock: func(mockRepo *MockFuncionarioRepo) {
				dbError := errors.New("falha de conexao")
				mockRepo.On("BuscaFuncionario", mock.Anything, 546373).
					Return(nil, dbError).
					Once()
			},
			erro: errors.New("falha de conexao"),
		},
		{
			nome:         "falha no banco ao adicionar",
			ctx:          context.Background(),
			modelAtestar: funcionarioInputTeste,
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscaFuncionario", mock.Anything, 546373). //erro esparado
											Return(nil, sql.ErrNoRows).
											Once()

				dbError := errors.New("erro ao inserir")
				mockRepo.On("AddFuncionario", mock.Anything, &funcionarioInput).
					Return(dbError).
					Once()
			},
			erro: errors.New("erro ao inserir"),
		},
	}

	for _, tc := range CasoDeTestes {

		t.Run(tc.nome, func(t *testing.T) {

			mockRepo := new(MockFuncionarioRepo)

			tc.mock(mockRepo)

			service := &FuncionarioService{
				FuncionarioRepo: mockRepo,
			}

			err := service.SalvarFuncionario(tc.ctx, tc.modelAtestar)

			if tc.erro != nil {

				require.Error(t, err)
				require.Equal(t, tc.erro.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBuscaFuncionario(t *testing.T) {

	funcionarioInput := model.Funcionario{
		Id:              1,
		Nome:            "davi",
		Matricula:       12323,
		ID_departamento: 2,
		Departamento:    "adm",
		ID_funcao:       1,
		Funcao:          "teste",
	}

	funcionarioDto := model.Funcionario_Dto{
		ID:        funcionarioInput.Id,
		Nome:      funcionarioInput.Nome,
		Matricula: funcionarioInput.Matricula,
		Departamento: model.DepartamentoDto{
			ID:           funcionarioInput.ID_departamento,
			Departamento: funcionarioInput.Departamento},
		Funcao: model.FuncaoDto{
			ID:     funcionarioInput.ID_funcao,
			Funcao: funcionarioInput.Funcao},
	}

	matriculaValidaStr := "12323"
	matriculaValidaINt := 12323
	matriculaINvalida := "str"

	CasoDeTestes := []struct {
		nome           string
		ctx            context.Context
		matriculaINput string
		mock           func(mockRepo *MockFuncionarioRepo)
		dto            *model.Funcionario_Dto
		erro           error
	}{

		{
			nome:           "retorna certo do dto",
			ctx:            context.Background(),
			matriculaINput: matriculaValidaStr,
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscaFuncionario", mock.Anything, matriculaValidaINt).
					Return(&funcionarioInput, nil).Once()
			},
			dto:  &funcionarioDto,
			erro: nil,
		},
		{
			nome:           "funcionario nao encontrado",
			ctx:            context.Background(),
			matriculaINput: matriculaValidaStr,
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscaFuncionario", mock.Anything, matriculaValidaINt).
					Return(nil, sql.ErrNoRows).Once()
			},
			dto:  nil,
			erro: errors.New("funcionario nao encontrado"),
		},
		{
			nome:           "matricula com formato invalido",
			ctx:            context.Background(),
			matriculaINput: matriculaINvalida,
			mock:           func(mockRepo *MockFuncionarioRepo) {},
			dto:            nil,
			erro:           errors.New("matricula tem que ser um numero"),
		},
		{
			nome:           "falha no banco de daddos",
			ctx:            context.Background(),
			matriculaINput: matriculaValidaStr,
			mock: func(mockRepo *MockFuncionarioRepo) {
				dbError := errors.New("erro no banco de dados")
				mockRepo.On("BuscaFuncionario", mock.Anything, matriculaValidaINt).
					Return(nil, dbError).Once()
			},
			dto:  nil,
			erro: errors.New("erro no banco de dados"),
		},
	}

	for _, tc := range CasoDeTestes {
		t.Run(tc.nome, func(t *testing.T) {
			mockRepo := new(MockFuncionarioRepo)
			tc.mock(mockRepo)

			s := FuncionarioService{
				FuncionarioRepo: mockRepo,
			}

			dto, err := s.ListarFuncionarioPorMatricula(tc.ctx, tc.matriculaINput)

			if tc.erro != nil {
				require.Error(t, err)
				require.Equal(t, tc.erro.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.dto, dto)

			mockRepo.AssertExpectations(t)
		})
	}

}

func TestListaTodosFuncinarios(t *testing.T) {

	FuncionarioDb := []model.Funcionario{
		{
			Id:              1,
			Nome:            "rada",
			Matricula:       234345,
			ID_departamento: 3,
			Departamento:    "ti",
			ID_funcao:       12,
			Funcao:          "dev",
		},
		{
			Id:              3,
			Nome:            "davi",
			Matricula:       23232,
			ID_departamento: 34,
			Departamento:    "teste",
			ID_funcao:       12,
			Funcao:          "teste",
		},
	}

	FuncionarioDto := []*model.Funcionario_Dto{
		{
			ID:        1,
			Nome:      "rada",
			Matricula: 234345,
			Departamento: model.DepartamentoDto{
				ID:           3,
				Departamento: "ti",
			},
			Funcao: model.FuncaoDto{
				ID:     12,
				Funcao: "dev",
			},
		},

		{
			ID:        3,
			Nome:      "davi",
			Matricula: 23232,
			Departamento: model.DepartamentoDto{
				ID:           34,
				Departamento: "teste",
			},
			Funcao: model.FuncaoDto{
				ID:     12,
				Funcao: "teste",
			},
		},
	}

	CasoDeTestes := []struct {
		nome string
		ctx  context.Context
		mock func(mockRepo *MockFuncionarioRepo)
		dto  []*model.Funcionario_Dto
		erro error
	}{

		{
			nome: "sucesso ao buscar todos os funcionarios",
			ctx:  context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscarTodosFuncionarios", mock.Anything).
					Return(FuncionarioDb, nil).Once()
			},
			dto:  FuncionarioDto,
			erro: nil,
		},
		{
			nome: "funcionarios nao encontrados",
			ctx:  context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscarTodosFuncionarios", mock.Anything).
					Return(nil, Errors.ErrBuscarTodos).Once()
			},
			dto:  []*model.Funcionario_Dto{},
			erro: nil,
		},
		{
			nome: "erro ao iterar",
			ctx:  context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscarTodosFuncionarios", mock.Anything).
					Return(nil, Errors.ErrAoIterar).Once()
			},
			dto:  []*model.Funcionario_Dto{},
			erro: errors.New("erro inesperado ao processar os dados dos funcionarios: erro ao iterar"),
		},
		{
			nome: "erro ao escanear dados",
			ctx:  context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("BuscarTodosFuncionarios", mock.Anything).
					Return(nil, Errors.ErrFalhaAoEscanearDados).Once()
			},
			dto:  []*model.Funcionario_Dto{},
			erro: errors.New("erro interno ao processar dados dos funcionarios, erro ao escanear os dados"),
		},
		{
			nome: "falha de conexao",
			ctx:  context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				dbErr := errors.New("falha no banco de dados")
				mockRepo.On("BuscarTodosFuncionarios", mock.Anything).
					Return(nil, dbErr).Once()
			},
			dto:  []*model.Funcionario_Dto{},
			erro: errors.New("erro inesperado ao buscar funcionarios, falha no banco de dados"),
		},
		{
			nome: "problema no context",
			ctx: func() context.Context {
				ctx, canceld := context.WithCancel(context.Background())
				canceld()
				return ctx
			}(),
			mock: func(mockRepo *MockFuncionarioRepo) {},
			dto:  nil,
			erro: errors.New("context canceled"),
		},
	}

	for _, tc := range CasoDeTestes {

		t.Run(tc.nome, func(t *testing.T) {
			mockRepo := new(MockFuncionarioRepo)
			tc.mock(mockRepo)
			s := FuncionarioService{
				FuncionarioRepo: mockRepo,
			}

			dto, err:= s.ListaTodosFuncionarios(tc.ctx)

			if tc.erro != nil {

				require.Error(t, err)
				require.Equal(t, tc.erro.Error(), err.Error())
			}else{
				require.NoError(t, err)
			}
			require.Equal(t, tc.dto, dto)

			mockRepo.AssertExpectations(t)

		})

	}
}

func TestDeleteFuncionario(t *testing.T){


	
	
	CasoDeTestes:= []struct{

		Nome string
		Matricula string
		ctx context.Context
		mock  func(mockRepo *MockFuncionarioRepo)
		erro	error

	}{
		{
			Nome: "sucesso ao deletar um funcionario",
			Matricula: "12345",
			ctx: context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("DeletarFuncionario", mock.Anything, 12345).Return(nil).Once()
			},
			erro: nil,
		},
		{
			Nome: "matricula invalida",
			Matricula: "reffgg",
			ctx: context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {},
			erro: errors.New("matricula tem que ser um numero"),
		},
		{
			Nome: "erro do repositorio",
			Matricula: "123456",
			ctx: context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("DeletarFuncionario", mock.Anything, 123456).Return(Errors.ErrInternal).Once()
			},
			erro: errors.New("erro interno ao processar dados, erro interno do repositório"),
		},
		{
			Nome: "erro ao verificar linhas afetadas",
			ctx: context.Background(),
			Matricula: "1234567",
			mock: func(mockRepo *MockFuncionarioRepo) {
				mockRepo.On("DeletarFuncionario", mock.Anything, 1234567).Return(Errors.ErrLinhasAfetadas).Once()
			},
			erro: errors.New("erro ao verificar linha afetada"),
		},
		{
			Nome: "erro inesperado",
			Matricula: "1245678",
			ctx: context.Background(),
			mock: func(mockRepo *MockFuncionarioRepo) {
				ErrDb:= errors.New("erro de conexao")
				mockRepo.On("DeletarFuncionario", mock.Anything, 1245678).Return(ErrDb)
			},
			erro: errors.New("erro inesperado ao deletar funcionario, erro de conexao"),

		},
		{
			Nome: "erro no context",
			Matricula: "112233",
			ctx: func() context.Context{
				ctx, canceld:= context.WithCancel(context.Background())
				canceld()
				return ctx
			}(),
			mock: func(mockRepo *MockFuncionarioRepo) {},
			erro: errors.New("context canceled"),
		},
	}



	for _, tc := range CasoDeTestes {
		t.Run(tc.Nome, func(t *testing.T) {
			mockRepo := new(MockFuncionarioRepo)
			tc.mock(mockRepo)

			s := FuncionarioService{
				FuncionarioRepo: mockRepo,
			}

			 err := s.DeletarFuncionario(tc.ctx, tc.Matricula)

			if tc.erro != nil {
				require.Error(t, err)
				require.Equal(t, tc.erro.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}

		

			mockRepo.AssertExpectations(t)
		})
	}
}
