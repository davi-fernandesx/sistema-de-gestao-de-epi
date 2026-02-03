package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/middleware"
	"github.com/gin-gonic/gin"
)

type FuncaoService interface {
	SalvarFuncao(ctx context.Context, model model.Funcao, tenantid int32) error
	ListarFuncao(ctx context.Context, id int, tenantid int32) (model.FuncaoDto, error)
	ListasTodasFuncao(ctx context.Context, tenantId int32) ([]model.FuncaoDto, error)
	DeletarFuncao(ctx context.Context, id int, tenantId int32) error
	AtualizarFuncao(ctx context.Context, id int, funcao string, tenantId int32) error
}

type FuncaoController struct {
	service FuncaoService
}

func NewFuncaoController(service FuncaoService) *FuncaoController {

	return &FuncaoController{service: service}
}

// RegistraFuncao godoc
// @Summary      Criar uma funcao
// @Description  Cadastra uma nova funcao no sistema
// @Tags         funcao
// @Accept       json
// @Produce      json
// @Param        funcao body model.Funcao true "Dados da funcao"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  helper.HTTPError "Dados inválidos"
// @Failure      409  {object}  helper.HTTPError "funcao já existe"
// @Failure      409  {object}  helper.HTTPError "id de departamento nao existe no sistema"
// @Failure      500  {object}  helper.HTTPError "Erro interno"
// @Router       /cadastro-funcao [post]
// @Security     BearerAuth
func (f *FuncaoController) RegistraFuncao() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var input model.Funcao

		if err := ctx.ShouldBindJSON(&input); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{

				"error": err.Error(),
			})
			return
		}

		novaFuncao := model.Funcao{
			Funcao:         input.Funcao,
			IdDepartamento: input.IdDepartamento,
		}
		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		err := f.service.SalvarFuncao(ctx, novaFuncao, tenantID)
		if err != nil {

			if errors.Is(err, helper.ErrDadoDuplicado) {
				ctx.JSON(http.StatusConflict, gin.H{

					"error": err.Error(),
				})
				return
			}

			if errors.Is(err, helper.ErrConflitoIntegridade) {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{

				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{

			"mensagem": "função cadastrada",
		})

	}
}

func (f *FuncaoController) ListarFuncoes() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		funcoes, err := f.service.ListasTodasFuncao(ctx, tenantID)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao listar funcoes",
			})
			return
		}

		ctx.JSON(http.StatusOK, funcoes)
	}
}

func (f *FuncaoController) ListarFuncaoId() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		idString := ctx.Param("id")

		id, err := strconv.Atoi(idString)
		if err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id deve ser um numero",
			})
			return
		}

		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		funcao, err := f.service.ListarFuncao(ctx, id, tenantID)
		if err != nil {

			if errors.Is(err, helper.ErrNaoEncontrado) {

				ctx.JSON(http.StatusNotFound, gin.H{

					"error": "funcao nao encontrada",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{

				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, funcao)
	}
}

func (f *FuncaoController) DeletarFuncao() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		idString := ctx.Param("id")

		id, err := strconv.Atoi(idString)
		if err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id deve ser um numero",
			})
			return
		}

		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		err = f.service.DeletarFuncao(ctx, id, tenantID)
		if err != nil {

			if errors.Is(err, helper.ErrNaoEncontrado) {

				ctx.JSON(http.StatusNotFound, gin.H{

					"error": "funcao nao encontrada",
				})

				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{

				"error": err.Error(),
			})
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func (f *FuncaoController) AtualizarFuncao() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"erro": err.Error(),
			})
			return
		}

		tenantID, ok := middleware.GetTenantID(ctx)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Erro interno de tenant"})
			return
		}

		var input model.Funcao

		if err := ctx.ShouldBindJSON(&input); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{

				"error": err.Error(),
			})
			return
		}

		err = f.service.AtualizarFuncao(ctx, id, input.Funcao, tenantID) 
		if err != nil {

			
			if errors.Is(err, helper.ErrNomeCurto) {

				ctx.JSON(http.StatusNotFound, gin.H{

					"error": "nome da funcao tem que possui 2 ou mais letras",
				})
				return
			}
			if errors.Is(err, helper.ErrNaoEncontrado) {

				ctx.JSON(http.StatusNotFound, gin.H{

					"error": "funcao nao encontrado para atualizar",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{

				"erro": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{"sucesso": "funcao atualizado"})

	}
}
