package service

import (
	"errors"
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


	usuario, err:= Ls.LoginRepo.BuscaPorNome(LoginUsuario.Nome)
	if err != nil {

		if errors.Is(err, login.ErrLinhasAfetadas){

			return false, nil
		}

		return  false, err

	}

	senhaLogin, err:= auth.HashCompare([]byte(usuario.Senha), []byte(LoginUsuario.Senha))
	if err != nil {

		return false, fmt.Errorf("erro ao comparar senhas")
	}




	return senhaLogin, nil
}