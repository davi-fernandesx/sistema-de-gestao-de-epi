package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
)


func VerificaEspaço(s string) (string){



	StringSemEspaço := strings.TrimSpace(s)


	return  StringSemEspaço
	
}

func VerificaContext(ctx context.Context) error{

	if err:= ctx.Err(); err != nil {

		return err
	}

	return nil
}

func VerificaMatricula(ctx context.Context, matricula string)(string, error){

	err:= 	VerificaContext(ctx)
	if err != nil {
		 
		return "", err
	}

	matriculaLimpa:= VerificaEspaço(matricula)
	if matriculaLimpa == ""{
		return "", errors.New("matricula nao pode estar em branco")
	}

	_, err  = strconv.Atoi(matriculaLimpa)
	if err != nil {

		return "", errors.New("matricula tem que ser um numero")
	}
	

	if len(matriculaLimpa) == 8 {
		return "", errors.New("matricula tem que conter ate 7 digitos numericos")
	}

	return  matriculaLimpa, nil

}