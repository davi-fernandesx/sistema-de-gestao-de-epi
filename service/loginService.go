package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/auth"
	Errors "github.com/davi-fernandesx/sistema-de-gestao-de-epi/errors"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/login"
)


type LoginServiceInterface interface {

	SalvaLogin(ctx context.Context, loginUsuario model.LoginDto) error
	Login(ctx context.Context, LoginUsuario model.LoginDto) (bool, error)

}
type LoginService struct {
	LoginRepo login.LoginRepository
}

func NewLoginService(Repo login.LoginRepository) LoginServiceInterface {

	return &LoginService{
		LoginRepo: Repo,
	}

}

func (Ls *LoginService) SalvaLogin(ctx context.Context, LoginUsuario model.LoginDto) error {

	if strings.TrimSpace(LoginUsuario.Nome) == "" {
		return fmt.Errorf("nome de usuario nao pode ser em branco")
	}

	SenhaHash, err := auth.HashPassword(LoginUsuario.Senha)
	if err != nil {
		return fmt.Errorf("erro ao criptografar a senha: %v", err)
	}
	loginModel := model.Login{
		Nome:  LoginUsuario.Nome,
		Senha: string(SenhaHash),
	}

	err = Ls.LoginRepo.AddLogin(ctx, &loginModel)
	if err != nil {

		return err
	}

	return nil
}

func (Ls *LoginService) Login(ctx context.Context, LoginUsuario model.LoginDto) (bool, error) {

	usuario, err := Ls.LoginRepo.BuscaPorNome(ctx, LoginUsuario.Nome)
	if err != nil {

		if errors.Is(err, Errors.ErrLinhasAfetadas) {

			return false, nil
		}

		return false, err

	}

	senhaLogin, err := auth.HashCompare([]byte(usuario.Senha), []byte(LoginUsuario.Senha))
	if err != nil {

		return false, fmt.Errorf("erro ao comparar senhas")
	}

	return senhaLogin, nil
}
