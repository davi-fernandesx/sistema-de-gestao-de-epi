/*secao departamento*/

insert into departamento (nome) values ('ti');
select nome from departamento where id = 1;
select id, nome from departamento;
update departamento
			set nome = 'rh'
			where id = 1


/*secao funcao*/

insert into funcao (nome, IdDepartamento) values ('analista', 1)
select id, nome from funcao where id = 1
select id, nome, IdDepartamento from funcao

select f.id, f.nome, f.IdDepartamento, d.nome as departamento
from funcao f
inner join 
	departamento d on f.IdDepartamento = d.id

update funcao
			set nome = 'secretaria'
			where id = 1

/*secao funcionario*/

insert into funcionario(nome, matricula, IdDepartamento, IdFuncao) values( 'davi', '757478', 1, 1)

select fn.id, fn.nome,fn.matricula, fn.IdDepartamento, d.nome as departamento, 
fn.IdFuncao, f.nome as funcao
			from funcionario fn
			inner join departamento d on fn.IdDepartamento = d.id
			inner join funcao f on fn.IdFuncao = f.id

select  * from funcionario
update funcionario
		     set IdDepartamento = 2
			 where id =  1

/*secao tipo_protecao*/

insert into tipo_protecao(nome) values ('protecao para as maos')
insert into tipo_protecao(nome) values ('protecao para os pes')
insert into tipo_protecao(nome) values ('protecao para a cabeca')
insert into tipo_protecao(nome) values ('protecao para o corpo')

select id, nome from tipo_protecao
delete from tipo_protecao where id = 1
/*secao tamanho*/

insert into tamanho values ('G')
insert into tamanho values ('M')
insert into tamanho values ('P')

select id, tamanho from tamanho where id = 2
select id, tamanho from tamanho
/*secao epi*/

insert into epi (nome, fabricante, CA, descricao, validade_CA,IdTipoProtecao, alerta_minimo) 
			OUTPUT INSERTED.id 
			values 
			('bota', 'mister', '2122', 'bota de borracha',GETDATE(),3, 10)

insert into tamanhos_epis(IdEpi, IdTamanho) values (1003,3);
insert into tamanhos_epis(IdEpi, IdTamanho) values (3,2);
insert into tamanhos_epis(IdEpi, IdTamanho) values (3,19);

select * from tamanhos_epis;

select
		e.id, e.nome, e.fabricante,e.CA, e.descricao,
		e.validade_CA, e.alerta_minimo, e.IdTipoProtecao, tp.nome as 'nome da protecao'
			from
				epi e
			inner join
				tipo_protecao tp on	e.IdTipoProtecao = tp.id	
select 
		t.id, t.tamanho
		from
			tamanho t
		inner join
			tamanhos_epis te on t.id = te.IdTamanho
		where te.IdEpi = 1003
delete epi from epi
select * from epi
select * from schema_migrations
update schema_migrations
	set dirty = 0

/*  secao de entrada */


insert into entrada_epi(IdEpi,IdTamanho, data_entrada, quantidade,data_fabricacao, data_validade, lote, fornecedor, valor_unitario)
		values (3,2, GETDATE()-40, 12 ,GETDATE(), GETDATE(), 'fff-ey', 'teste',22.99 )

select ee.id, ee.IdEpi,  e.nome as epi,  e.fabricante, e.CA, e.descricao,ee.data_fabricacao, ee.data_validade, e.validade_CA,
		e.IdTipoProtecao, tp.nome as 'protecao para',
	   	ee.IdTamanho,t.tamanho as tamanho, ee.quantidade, ee.data_entrada,
	   ee.lote, ee.fornecedor, ee.valor_unitario
from entrada_epi ee
inner join
	epi e on ee.IdEpi = e.id
inner join
	tipo_protecao tp on e.IdTipoProtecao = tp.id
inner join
	tamanho t on ee.IdTamanho = t.id
where ee.cancelada_em is not null 

select * from entrada_epi
update entrada_epi
			set cancelada_em = GETDATE()
			where id = 4 AND cancelada_em IS NULL
			
/*secao de entrega*/

insert into entrega_epi(IdFuncionario, data_entrega, assinatura)
	OUTPUT INSERTED.id
	 values (1, GETDATE(), CAST('We will store this string as varbinary' AS VARBINARY(MAX)))

select * from entrega_epi


select  top 1 id as entrada, valor_unitario , quantidade ,lote, data_entrada, IdTamanho, IdEpi
from entrada_epi
where IdEpi = 1003 AND IdTamanho = 2 and quantidade >= 1
order by data_entrada asc

		update entrada_epi 
				set quantidade = quantidade - 1
					where id = 1003

select  id as entrada, valor_unitario , quantidade ,lote, data_entrada, IdTamanho
from entrada_epi
select * from entrada_epi



select * from epis_entregues

select
		    ee.id,
			ee.data_Entrega,
			ee.IdFuncionario,
			f.nome, 
			f.IdDepartamento, 
			d.nome, 
			f.IdFuncao, 
			ff.nome, 
			i.id, 
			e.nome, 
			e.fabricante, 
			e.CA,
			e.descricao,  
			e.validade_CA,
			e.IdTipoProtecao,
			tp.id, 
			i.IdTamanho, 
			t.tamanho, 
			i.quantidade,
			ee.assinatura,
			i.valor_unitario,
			ee.Idtroca
			from entrega_epi ee
			inner join
				funcionario f on ee.IdFuncionario = f.id
			inner join
				departamento d on f.IdDepartamento = d.id
			inner join 
				funcao ff on f.IdFuncao = ff.id
			inner join 
				epis_entregues i on i.IdEntrega = ee.id
			inner join 
				epi e on i.IdEpi = e.id
			inner join
				tipo_protecao tp on e.IdTipoProtecao = tp.id
			inner join 
				tamanho t on i.IdTamanho = t.id
where ee.cancelada_em IS NULL

/*  secao devolucao e motivo devolucao*/

	insert into devolucao (IdFuncionario, IdEpi, IdMotivo ,data_devolucao, IdTamanho, quantidadeAdevolver, idEpiNovo, IdTamanhoNovo,quantidade_nova,assinatura_digital)
	OUTPUT INSERTED.id
	 	values (1, 1003, 3 ,GETDATE(), 1, 1, 1003,3, 1, cast('teste' as varbinary(max)))


	insert into entrega_epi(IdFuncionario, data_entrega, assinatura, IdTroca)
		OUTPUT INSERTED.id
		values (1, GETDATE(), cast('ok' as varbinary(max)), 5)

	insert into epis_entregues(IdEpi,IdTamanho, quantidade,IdEntrega,IdEntrada ,valor_unitario) values (3, 1, 1,6,1005 ,22.99)

		select * from entrega_epi
		

select * from devolucao

insert into motivo_devolucao(motivo) values ('Numeração ou tamanho errado')
insert into motivo_devolucao(motivo) values ('Substituição por Desgaste ou Dano')
insert into motivo_devolucao(motivo) values ('Vencimento da validade ou do CA')
insert into motivo_devolucao(motivo) values ('Mudança de Função ou Setor')
insert into motivo_devolucao(motivo) values ('Demissão')

select * from motivo_devolucao