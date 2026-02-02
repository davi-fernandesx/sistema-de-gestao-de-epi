package routers

import (
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/controller"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/database/repository"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/service"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_"github.com/davi-fernandesx/sistema-de-gestao-de-epi/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


type Container struct {

	Usuario  controller.LoginController
	Departamento controller.DepartamentoController
	Funcao 	controller.FuncaoController
}


func NewContainer(db *pgxpool.Pool) *Container {

	repoUsuario:= repository.NewUsuarioRepository(db)
	repoDepartamento:= repository.NewDepartamentoRepository(db)
	repoFuncao:= repository.NewFuncaoRepository(db)



	serviceUsuario:= service.NewUsuarioService(repoUsuario)
	departamentoService:= service.NewDepartamentoService(repoDepartamento)
	funcaoService:= service.NewFuncaoService(repoFuncao)


	return &Container{
		Usuario: *controller.NewLoginController(serviceUsuario),
		Departamento: *controller.NewDepartamentoController(departamentoService),
		Funcao: *controller.NewFuncaoController(funcaoService),
	}
}
func ConfigurarRotas(r *gin.Engine, c *Container) {


	// --- GRUPO 1: Rotas Públicas (Aberta) ---
    // Qualquer um acessa sem token

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/cadastro", c.Usuario.Registrar())
	r.POST("/login", c.Usuario.Login())

	// --- GRUPO 2: Rotas Protegidas (SaaS) ---
    // Precisa do Token JWT para passar

	api:= r.Group("/api")

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
	}

}

