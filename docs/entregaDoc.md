## Fluxo de Registro de Entrega de EPI

Este diagrama apresenta o fluxo completo do processo de registrar uma entrega de Equipamento de Proteção Individual (EPI) com controle transacional e gerenciamento de estoque.

### Componentes Principais

**Inicialização:**
- Inicia uma transação no banco de dados
- Cria um contexto de transação (qtx) para garantir consistência
- Busca os dados do funcionário destinatário
- Gera um token de auditoria para rastreamento

**Cabeçalho da Entrega:**
- Salva o registro principal da entrega na tabela `entrega_epi`
- Valida integridade antes de prosseguir com itens

**Loop Externo - Processamento de Itens:**
- Itera sobre cada item solicitado na entrega
- Para cada item, busca lotes disponíveis ordenados por data/validade (FIFO)
- Valida disponibilidade de estoque antes de processar

**Loop Interno - Consumo de Lotes (FIFO):**
- Implementa abatimento de estoque por lote usando estratégia First-In-First-Out
- Calcula quantidade a abater (menor entre disponível no lote e necessário)
- Salva registro de item entregue vinculado ao ID de entrada
- Atualiza estoque decrementando quantidade do lote
- Continua até satisfazer quantidade necessária do item

**Tratamento de Erros:**
- Valida existência de estoque em cada etapa
- Realiza rollback em caso de erro (garante transação atômica)
- Retorna mensagem de erro específico (estoque insuficiente)

**Finalização:**
- Commit da transação se todos os itens forem processados com sucesso
- Retorna confirmação da entrega completada
```mermaid
flowchart TD
    Start([Início]) --> Tx[Iniciar Transação]
    Tx --> Qtx[Criar qtx com Transação]
    
    %% Início da Lógica de RegistrarEntrega
    Qtx --> BuscaFunc[Busca Funcionário]
    BuscaFunc --> Token[Gerar Token Auditoria]
    
    %% AQUI ESTAVA O ERRO: Adicionei aspas
    Token --> SaveHeader["Salvar Cabeçalho da Entrega(Tabela entrega_epi)"]
    
    SaveHeader -- Erro --> Rollback
    SaveHeader -- Sucesso --> LoopItens{Há Itens na Lista?}

    %% LOOP EXTERNO: ITENS SOLICITADOS
    subgraph LoopItensGraph [Loop: Processar Itens]
        LoopItens -- Sim --> GetLotes["Listar Lotes Disponíveis\n(Ordenado por Data/Validade)"]
        
        GetLotes --> CheckLotes{Existem Lotes?}
        CheckLotes -- Não --> ErroEstoque["Erro: Estoque Insuficiente"]
        
        CheckLotes -- Sim --> LoopLotes{Qtd Necessária > 0?}
        
        %% LOOP INTERNO: CONSUMO DE LOTES (FIFO)
        subgraph LoopLotesGraph [Loop Interno: Abater Lotes]
            LoopLotes -- Sim --> Calc["Calc. Qtd a Abater\n(Mínimo entre Lote e Necessário)"]
            Calc --> SaveItem["Salvar Item Entregue\n(Vincula ao ID da Entrada)"]
            
            SaveItem -- Erro --> Rollback
            SaveItem -- Sucesso --> UpdateEstoque[Abater Estoque da Entrada]
            
            UpdateEstoque -- Erro --> Rollback
            UpdateEstoque -- Sucesso --> Decrement[Subtrair Qtd Necessária]
            Decrement --> LoopLotes
        end
        
        LoopLotes -- "Não (Satisfeito)" --> CheckFinalItem{Sobrou Necessidade?}
        CheckFinalItem -- "Sim (Acabaram os lotes)" --> ErroEstoque
        CheckFinalItem -- Não --> LoopItens
    end

    %% Finalização
    ErroEstoque --> Rollback["Rollback & Retornar Erro"]
    LoopItens -- "Não (Todos Processados)" --> Commit[Commit Transação]
    
    Commit --> Fim(["Fim (Sucesso)"])

```


## RegistrarCancelamento - Fluxo de Processo

Este diagrama descreve o fluxo de cancelamento de uma entrega de EPI, incluindo validações, transações de banco de dados e reposição de estoque.

### Fluxo Principal

1. **Validação Inicial**: Verifica se o ID da entrega é válido (> 0)
2. **Transação**: Inicia uma transação de banco de dados para garantir consistência
3. **Cancelamento do Cabeçalho**: Atualiza o status da entrega como cancelada
4. **Cancelamento dos Itens**: Marca todos os itens da entrega como cancelados
5. **Reposição de Estoque**: Loop que repõe o estoque de cada lote associado aos itens cancelados
6. **Confirmação**: Commit da transação se todos os passos forem bem-sucedidos

### Tratamento de Erros

- **ID Inválido**: Retorna erro se ID ≤ 0
- **Entrega Não Encontrada**: Retorna erro se nenhum registro for atualizado
- **Erro de Banco de Dados**: Rollback da transação em caso de falha em qualquer operação de atualização
- **Lote Não Encontrado**: Rollback se o lote de reposição não existir

### Garantias

- **ACID**: Transação garante atomicidade - sucesso total ou falha total (rollback)
- **Estorno Completo**: Todos os itens cancelados têm seu estoque reposto
```mermaid
flowchart TD
    Start([Início]) --> CheckID{ID > 0?}
    CheckID -- "Não" --> ErrID[Retornar Erro ID]
    
    CheckID -- "Sim" --> Tx[Iniciar Transação]
    Tx --> Qtx[Criar qtx com Transação]

    %% Início da Lógica de RegistrarCancelamento
    Qtx --> CancelHeader["Cancelar Cabeçalho\n(Update entrega_epi)"]
    
    CancelHeader -- Erro --> Rollback
    CancelHeader -- Sucesso --> CheckFound{Encontrou\nEntrega?}
    
    CheckFound -- "Não (Zero Linhas)" --> ErrNotFound["Retornar Erro\nNão Encontrado"]
    ErrNotFound --> Rollback

    CheckFound -- Sim --> CancelItems["Cancelar Itens da Entrega\n(Update epis_entregues)"]
    
    CancelItems -- Erro --> Rollback
    CancelItems -- Sucesso --> ListItems["Listar Itens Cancelados\n(Busca dados para estorno)"]

    %% Loop de Reposição
    ListItems -- Sucesso --> Loop{Há Itens?}
    
    subgraph LoopReposicao [Loop: Estorno de Estoque]
        Loop -- Sim --> Repor["Repor Estoque do Lote\n(Update entrada_epi)"]
        
        Repor -- "Erro BD" --> Rollback
        Repor -- "Erro (Lote não achado)" --> Rollback
        
        Repor -- Sucesso --> Loop
    end

    %% Finalização
    Loop -- "Não (Todos Processados)" --> Commit[Commit Transação]
    
    Commit --> Fim(["Fim (Sucesso)"])
    
    %% AQUI ESTAVA O ERRO (Adicionei aspas):
    Rollback["Rollback (Desfaz Tudo)"] --> FimErro(["Fim (Erro)"])
```