package helper

import (
	"regexp"
	"strconv"
	"github.com/go-playground/validator/v10" // Importante: v10
)

func IsCNPJ(cnpj string) bool {
   
		// 1. Remove caracteres não numéricos
	reg := regexp.MustCompile("[^0-9]")
	cnpj = reg.ReplaceAllString(cnpj, "")

	// 2. Valida tamanho
	if len(cnpj) != 14 {
		return false
	}

	// 3. Elimina inválidos conhecidos
	if cnpj == "00000000000000" || cnpj == "11111111111111" ||
		cnpj == "22222222222222" || cnpj == "33333333333333" ||
		cnpj == "44444444444444" || cnpj == "55555555555555" ||
		cnpj == "66666666666666" || cnpj == "77777777777777" ||
		cnpj == "88888888888888" || cnpj == "99999999999999" {
		return false
	}

	// 4. Valida Dígitos (Algoritmo padrão)
	tamanho := len(cnpj) - 2
	numeros := cnpj[0:tamanho]
	digitos := cnpj[tamanho:]
	soma := 0
	pos := tamanho - 7

	for i := tamanho; i >= 1; i-- {
		num, _ := strconv.Atoi(string(numeros[tamanho-i]))
		soma += num * pos
		pos--
		if pos < 2 {
			pos = 9
		}
	}

	resultado := soma % 11
	if resultado < 2 {
		resultado = 0
	} else {
		resultado = 11 - resultado
	}

	digito1, _ := strconv.Atoi(string(digitos[0]))
	if resultado != digito1 {
		return false
	}

	tamanho = tamanho + 1
	numeros = cnpj[0:tamanho]
	soma = 0
	pos = tamanho - 7

	for i := tamanho; i >= 1; i-- {
		num, _ := strconv.Atoi(string(numeros[tamanho-i]))
		soma += num * pos
		pos--
		if pos < 2 {
			pos = 9
		}
	}

	resultado = soma % 11
	if resultado < 2 {
		resultado = 0
	} else {
		resultado = 11 - resultado
	}

	digito2, _ := strconv.Atoi(string(digitos[1]))
	return resultado == digito2
}


// ValidateCNPJ é a função que o Gin vai chamar
func ValidateCNPJ(fl validator.FieldLevel) bool {
	
	return  IsCNPJ(fl.Field().String())
}


