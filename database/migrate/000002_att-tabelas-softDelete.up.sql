-- TABELA: epi
ALTER TABLE epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE epi ADD deletado_em DATETIME NULL;


-- TABELA: funcionario
ALTER TABLE funcionario ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE funcionario ADD deletado_em DATETIME NULL;


-- TABELA: entrada_epi
ALTER TABLE entrada_epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;


-- TABELA: entrega_epi
ALTER TABLE entrega_epi ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;


-- TABELA: devolucao
ALTER TABLE devolucao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;


-- TABELA: departamento
ALTER TABLE departamento ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE departamento ADD deletado_em DATETIME NULL;


-- TABELA: funcao
ALTER TABLE funcao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE funcao ADD deletado_em DATETIME NULL;


-- TABELA: tipo_protecao
ALTER TABLE tipo_protecao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE tipo_protecao ADD deletado_em DATETIME NULL;


-- TABELA: tamanho
ALTER TABLE tamanho ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE tamanho ADD deletado_em DATETIME NULL;


-- TABELA: tamanhos_epis
ALTER TABLE tamanhos_epis ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE tamanhos_epis ADD deletado_em DATETIME NULL;


-- TABELA: epis_entregues (Cuidado: Verifique se essa tabela já não é a 'entrega_epi')
ALTER TABLE epis_entregues ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE epis_entregues ADD deletado_em DATETIME NULL;


-- TABELA: motivo_devolucao
ALTER TABLE motivo_devolucao ADD ativo BIT NOT NULL DEFAULT 1 WITH VALUES;
ALTER TABLE motivo_devolucao ADD deletado_em DATETIME NULL;


