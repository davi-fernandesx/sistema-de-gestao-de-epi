CREATE UNIQUE INDEX uq_epi_ca_ativo 
ON epi(nome, CA) 
WHERE ativo = 1;