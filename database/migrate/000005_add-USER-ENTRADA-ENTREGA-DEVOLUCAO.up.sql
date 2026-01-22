-- 1. Rastrear quem deu entrada nos lotes/EPIs
ALTER TABLE entrada_epi 
ADD COLUMN id_usuario_criacao INTEGER REFERENCES usuarios(id);

-- 2. Rastrear quem realizou a entrega para o funcionário
ALTER TABLE entrega_epi 
ADD COLUMN id_usuario_entrega INTEGER REFERENCES usuarios(id);

-- 3. Rastrear quem realizou o cancelamento (estorno)
ALTER TABLE devolucao 
ADD COLUMN id_usuario_cancelamento INTEGER REFERENCES usuarios(id);



ALTER TABLE entrada_epi 
ADD COLUMN id_usuario_criacao_cancelamento INTEGER REFERENCES usuarios(id);

-- 2. Rastrear quem realizou a entrega para o funcionário
ALTER TABLE entrega_epi 
ADD COLUMN id_usuario_entrega_cancelamento INTEGER REFERENCES usuarios(id);

-- 3. Rastrear quem realizou o cancelamento (estorno)
ALTER TABLE devolucao 
ADD COLUMN id_usuario_devolucao_cancelamento INTEGER REFERENCES usuarios(id);