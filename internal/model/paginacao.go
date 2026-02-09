package model

type PaginacaoParams struct {
    Pagina int32 `form:"pagina,default=1" binding:"min=1"`
    Limite int32 `form:"limite,default=10" binding:"min=1,max=100"`
}