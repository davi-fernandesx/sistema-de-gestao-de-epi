
create database sistema_epi;


create TABLE login (

    id int PRIMARY key IDENTITY(1,1),
    usuario VARCHAR(50) not null,
    senha VARCHAR(256) not NULL,

    

);
