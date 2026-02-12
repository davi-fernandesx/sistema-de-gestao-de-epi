package routers

import (
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/controller"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	_ "github.com/davi-fernandesx/sistema-de-gestao-de-epi/docs"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/service"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Container struct {
	Usuario      controller.LoginController
	Departamento controller.DepartamentoController
	Funcao       controller.FuncaoController
	Funcionario  controller.FuncionarioController
	Tamanho      controller.TamanhoController
	Protecao     controller.TipoProtecaoController
	Epi          controller.EpiController
	Entrada      controller.EntradaController
}

func NewContainer(db *pgxpool.Pool) *Container {

	repoUsuario := repository.NewUsuarioRepository(db)
	repoDepartamento := repository.NewDepartamentoRepository(db)
	repoFuncao := repository.NewFuncaoRepository(db)
	repoFuncionario := repository.NewFuncionarioRepository(db)
	repoTamanho := repository.NewTamanhoRepository(db)
	repoTipoProtecao := repository.NewProtecaoRepository(db)
	repoEpi := repository.NewEpiRepository(db)
	repoEntrada := repository.NewEntradaRepository(db)

	serviceUsuario := service.NewUsuarioService(repoUsuario)
	departamentoService := service.NewDepartamentoService(repoDepartamento)
	funcaoService := service.NewFuncaoService(repoFuncao)
	funcionarioService := service.NewFuncionarioService(repoFuncionario, db)
	tamanhoService := service.NewTamanhoService(repoTamanho)
	TipoProtecaoService := service.NewProtecaoService(repoTipoProtecao)
	epiService := service.NewEpiService(repoEpi, db)
	entradaService := service.NewEntradaService(repoEntrada)

	return &Container{
		Usuario:      *controller.NewLoginController(serviceUsuario),
		Departamento: *controller.NewDepartamentoController(departamentoService),
		Funcao:       *controller.NewFuncaoController(funcaoService),
		Funcionario:  *controller.NewFuncionarioController(funcionarioService),
		Tamanho:      *controller.NewTamanhoControle(tamanhoService),
		Protecao:     *controller.NewTipoProtecaoController(TipoProtecaoService),
		Epi:          *controller.NewEpiController(epiService),
		Entrada:      *controller.NewEntradaController(entradaService),
	}
}
func ConfigurarRotas(r *gin.Engine, c *Container, db *pgxpool.Pool) {

	queries := repository.New(db)
	// --- GRUPO 1: Rotas Públicas (Aberta) ---
	// Qualquer um acessa sem token

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	// --- GRUPO 2: Rotas que precisam do tenentId (SaaS) ---
	// Precisa do tenant Id para passar
	api.Use(middleware.TenantMiddleware(queries))
	{

		api.POST("/cadastro", c.Usuario.Registrar())
		api.POST("/login", c.Usuario.Login())
	}

	// --- GRUPO 3: Rotas Protegidas (SaaS) ---
	// Precisa do Token JWT para passar
	api.Use(middleware.AutenticacaoJWT(), middleware.LoggerComUsuario())
	{

		api.GET("/me", c.Usuario.VerPerfil())
		//departamentos
		api.POST("/cadastro-departamento", c.Departamento.RegistraDepartamento())
		api.GET("/departamentos", c.Departamento.ListarDepartamentos())
		api.GET("/departamentos/:id", c.Departamento.ListarDepartamentoId())
		api.DELETE("/departamento/:id", c.Departamento.DeletarDepartamento())
		api.PUT("/departamento/:id", c.Departamento.AtualizarDepartamento())

		//funcao
		api.POST("cadastro-funcao", c.Funcao.RegistraFuncao())
		api.GET("/funcoes", c.Funcao.ListarFuncoes())
		api.GET("/funcao/:id", c.Funcao.ListarFuncaoId())
		api.DELETE("/funcao/:id", c.Funcao.DeletarFuncao())
		api.PUT("/funcao/:id", c.Funcao.AtualizarFuncao())

		//funcionario
		api.POST("/cadastro-funcionario", c.Funcionario.Adicionar())
		api.GET("/funcionarios", c.Funcionario.ListarFuncionarios())
		api.GET("/funcionario/:matricula", c.Funcionario.ListarFuncionarioPorMatricula())
		api.DELETE("/funcionario/:id", c.Funcionario.DeletarFuncionaioId())
		api.PATCH("/funcionario/:id", c.Funcionario.AtualizaFuncionario())

		//tamanhos disponiveis para vincular a um epi
		api.POST("/cadastro-tamanho", c.Tamanho.Adicionar())
		api.GET("/tamanhos", c.Tamanho.ListarTodosTamanhos())
		api.GET("/tamanho/:id", c.Tamanho.ListarTamanhoPorId())
		api.DELETE("/tamanho/:id", c.Tamanho.DeletarTamanho())

		//proteções dedicada a cada epi
		api.POST("/cadastro-protecao", c.Protecao.AdicionarProtecao())
		api.GET("/protecoes", c.Protecao.ListarProtecoes())
		api.GET("/protecao/:id", c.Protecao.ListarProtecaoPorId())
		api.DELETE("/protecao/:id", c.Protecao.DeletarProtecao())

		//Epi´s
		api.POST("/cadastro-epi", c.Epi.AdicionarEpi())
		api.GET("/epis", c.Epi.ListarEpis())
		api.GET("/epi/:id", c.Epi.ListarEpiPorId())
		api.DELETE("/epi/:id", c.Epi.DeletarEpi())
		api.PATCH("/epi/:id", c.Epi.AtualizaEpi())

		//entradas
		api.POST("/cadastrar-entrada", c.Entrada.AdicionarEntrada())
		api.GET("/entradas", c.Entrada.ListarEntradas())
		api.DELETE("/entrada/:id", c.Entrada.CancelarEntrada())
	}

}
