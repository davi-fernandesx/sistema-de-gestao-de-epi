package controller

/* Sobre o uso do Tenant id, pra que serve ?

o TenantId é usado em todas as rotas do sistema, com ele é possivel  definir de qual empresa é os dados.
o id vem do banco de dados, da tabela de empresas, aonde armazenar os dados basicos de cada empresa que possui o sistema

para fazer isso, primeiro salvamos os dados da empresa no banco, com isso esses dados recebe um id
um desses dados é o subdominio (ex: frigo)

ao passar a url no navegador, ex: http://frigo.lvh.me:8080/api/tamanhos
uso uma funcao para "fatiar" essa url, pegando apenas o valor de "frigo"

depois disso, chamo uma funcao, que compara o nome pega na url com o subdominio salvo no banco
se forem iguais, me retornar esse id

com esse id, salvo ele na memoria temporaria da requisição, ao salvar nessa memoria, eu "passo" esse id
para as proximas requisiçoes que vierem, com isso o banco de dados sabe exatamente quais dados trazer, nao
misturando dados da empresa X com a empresa Y

detalhes da funcao ,na pasta middleware, arquivo TenantId.go

*/