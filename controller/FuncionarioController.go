package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/middleware"
	"github.com/gin-gonic/gin"
)

type FuncionarioService interface {
	SalvarFuncionario(ctx context.Context, model model.FuncionarioINserir, tenantId int32) error
	ListarFuncionario(ctx context.Context, matricula string, tenantId int32) (model.Funcionario_Dto, error)
	ListaTodosFuncionarios(ctx context.Context, tenantId int32) ([]model.Funcionario_Dto, error)
	DeletarFuncionario(ctx context.Context, id int, tenantId int32) error
}

type FuncionarioController struct {
	Service FuncionarioService
}

func NewFuncionarioController(service FuncionarioService) *FuncionarioController {

	return &FuncionarioController{
		Service: service,
	}
}

// RegistraDepartamento godoc
// @Summary      Cadastrar novo funcionarios
// @Description  Cadastra um novo funcionario no sistema
// @Tags         funcionario
// @Accept       json
// @Produce      json
// @Param        funcionario body model.FuncionarioINserir true "Dados do funcionario"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  helper.HTTPError "Dados inválidos"
// @Failure      409  {object}  helper.HTTPError "Departamento já existe"
// @Failure      500  {object}  helper.HTTPError "Erro interno"
// @Router       /cadastro-funcionario [post]
// @Security     BearerAuth
func (f *FuncionarioController) Adicionar() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var input model.FuncionarioINserir

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":    "dados invalidos",
				"detalhes": err.Error(),
			})
			return
		}

		novoFunc := model.FuncionarioINserir{
			Nome:            input.Nome,
			Matricula:       input.Matricula,
			ID_departamento: input.ID_departamento,
			ID_funcao:       input.ID_funcao,
		}
		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		err := f.Service.SalvarFuncionario(ctx, novoFunc, tenantID)
		if err != nil {

			if errors.Is(err, helper.ErrDadoDuplicado) {
				ctx.JSON(http.StatusConflict, gin.H{

					"error": "matricula ja cadastrada",
				})
				return
			}

			if errors.Is(err, helper.ErrConflitoIntegridade) {

				ctx.JSON(http.StatusNotFound, gin.H{

					"error": "departamento ou funcao nao encontrado",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{

				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{

			"mensagem": "funcionario cadastrado",
		})
	}

}

// ListarFuncionarios godoc
// @Summary      Listar todos
// @Description  Retorna uma lista com todos os funcionarios
// @Tags         funcionarios
// @Produce      json
// @Success      200  {array}   model.Funcionario_dto
// @Failure      500  {object}  helper.HTTPError "Erro interno"
// @Router       /funcionarios [get]
// @Security     BearerAuth
func (f *FuncionarioController) ListarFuncionarios() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		funcs, err := f.Service.ListaTodosFuncionarios(ctx, tenantID)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao listar departamentos",
			})
			return
		}

		ctx.JSON(http.StatusOK, funcs)
	}
}
