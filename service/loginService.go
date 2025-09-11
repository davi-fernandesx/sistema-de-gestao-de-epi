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