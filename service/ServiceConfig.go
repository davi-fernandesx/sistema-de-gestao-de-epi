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

func VerificaMatricula(ctx context.Context, matricula string)(int, error){

	err:= 	VerificaContext(ctx)
	if err != nil {
		 
		return 0, err
	}

	matriculaLimpa:= VerificaEspaço(matricula)
	if matriculaLimpa == ""{
		return 0, errors.New("matricula nao pode estar em branco")
	}

	matriculaInt, err := strconv.Atoi(matriculaLimpa)
	if err != nil {

		return 0, errors.New("matricula tem que ser um numero")
	}
	

	if len(matriculaLimpa) == 8 {
		return 0, errors.New("matricula tem que conter ate 7 digitos numericos")
	}

	return  matriculaInt, nil

}