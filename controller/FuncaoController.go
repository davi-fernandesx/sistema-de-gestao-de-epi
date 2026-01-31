package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/helper"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/internal/model"
	"github.com/gin-gonic/gin"
)

type FuncaoService interface {
	SalvarFuncao(ctx context.Context, model model.Funcao) error
	ListarFuncao(ctx context.Context, id int) (model.FuncaoDto, error)
	ListasTodasFuncao(ctx context.Context) ([]model.FuncaoDto, error)
	DeletarFuncao(ctx context.Context, id int) error
	AtualizarFuncao(ctx context.Context, id int, funcao string) error
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

		err := f.service.SalvarFuncao(ctx, novaFuncao)
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
