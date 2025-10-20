package model

import "time"




type ItemParaEntrega struct {
    ID_epi     int `json:"id_epi"`
    ID_tamanho int `json:"id_tamanho"`
    Quantidade int `json:"quantidade"`
}


type EntregaParaInserir struct {
    ID_funcionario     int               `json:"id_funcionario"`
    Data_entrega       time.Time         `json:"data_entrega"`
    Assinatura_Digital string            `json:"assinatura_digital"`
    Itens              []ItemParaEntrega `json:"itens"`
}




type ItemEntregueDto struct {
    Epi        EpiDto     `json:"epi"`
    Tamanho    TamanhoDto `json:"tamanho"`
    Quantidade int        `json:"quantidade"`
}

type EntregaDto struct {
    Id                 int               `json:"id"`
    Funcionario        Funcionario_Dto   `json:"funcionario"`
    Data_entrega       time.Time         `json:"data_entrega"`
    Assinatura_Digital string            `json:"assinatura_digital"`
    Itens              []ItemEntregueDto `json:"itens"`
}
