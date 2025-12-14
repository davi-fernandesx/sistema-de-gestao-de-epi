-- 1. Primeiro, apagamos a regra antiga (use o nome que você achou no passo 1)
ALTER TABLE tamanhos_epis 
DROP CONSTRAINT FK__tamanhos___IdTam__14270015;

-- 2. Agora criamos a nova regra (Dê um nome bonito agora!)
ALTER TABLE tamanhos_epis
ADD CONSTRAINT fk_tamanhos_epis_tamanho -- Nome que você escolheu
FOREIGN KEY (IdTamanho) 
REFERENCES tamanho(id)
ON DELETE CASCADE; -- A mudança que você queria fazer