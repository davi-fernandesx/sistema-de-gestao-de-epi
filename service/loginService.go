package service

import (

	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/auth"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/login"
)

type LoginService struct {

	LoginRepo  *login.SqlServerLogin
}

func NewLoginService(loginRepo  *login.SqlServerLogin) *LoginService{

	return  &LoginService{
		LoginRepo: loginRepo,
	}
}


func (Ls *LoginService)SalvaLogin(LoginUsuario model.LoginDto) error {

	if strings.TrimSpace(LoginUsuario.Nome) == " "{
		return fmt.Errorf("nome de usuario nao pode ser em branco")
	}


	SenhaHash, err:= auth.HashPassword(LoginUsuario.Senha)
	if err != nil {
		return fmt.Errorf("erro ao criptografar a senha: %v", err)
	}
	loginModel:= model.Login{
		Nome: LoginUsuario.Nome,
		Senha: string(SenhaHash),
	}

	err =  Ls.LoginRepo.AddLogin( &loginModel)
	if err != nil {

		return  err
	}

	return  nil
}

func (Ls *LoginService) Login(LoginUsuario model.LoginDto) (bool,error) {

	usuarios, err:= Ls.LoginRepo.RetornaLogin();
	if err != nil {
		return  false, err
	}
	modelUsuario:= model.Login{
		Nome: LoginUsuario.Nome,
		Senha: LoginUsuario.Senha,
	}

	for _, usuario:= range *usuarios {

		if usuario.Nome == modelUsuario.Nome {

			LoginAceito, err:=auth.HashCompare([]byte(usuario.Senha), []byte(modelUsuario.Senha))
			if err != nil {
				return false, nil
			}

			if LoginAceito {
				return  true, nil
			}

		}
	}

	return false, fmt.Errorf("fim da funcao login")

}