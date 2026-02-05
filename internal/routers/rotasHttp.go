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
}

func NewContainer(db *pgxpool.Pool) *Container {

	repoUsuario := repository.NewUsuarioRepository(db)
	repoDepartamento := repository.NewDepartamentoRepository(db)
	repoFuncao := repository.NewFuncaoRepository(db)
	repoFuncionario := repository.NewFuncionarioRepository(db)

	serviceUsuario := service.NewUsuarioService(repoUsuario)
	departamentoService := service.NewDepartamentoService(repoDepartamento)
	funcaoService := service.NewFuncaoService(repoFuncao)
	funcionarioService := service.NewFuncionarioService(repoFuncionario, db)

	return &Container{
		Usuario:      *controller.NewLoginController(serviceUsuario),
		Departamento: *controller.NewDepartamentoController(departamentoService),
		Funcao:       *controller.NewFuncaoController(funcaoService),
		Funcionario:  *controller.NewFuncionarioController(funcionarioService),
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
	}

}
