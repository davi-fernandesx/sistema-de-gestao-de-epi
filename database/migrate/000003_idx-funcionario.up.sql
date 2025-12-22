CREATE UNIQUE INDEX uq_funcionario_ativo 
ON funcionario(matricula) 
WHERE ativo = 1;