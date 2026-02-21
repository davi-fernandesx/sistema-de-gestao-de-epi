package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestIsCNPJ(t *testing.T) {
	// 1. Nossa "Tabela" de casos de teste
	testCases := []struct {
		nome     string
		cnpj     string
		esperado bool
	}{
		// --- CENÁRIOS DE SUCESSO (expected: true) ---
		{"Válido formatado", "12.345.678/0001-95", true},
		{"Válido sem formatação", "12345678000195", true},
		{"Válido outro exemplo formatado", "92.661.872/0001-19", true},
		{"Válido outro exemplo limpo", "92661872000119", true},

		// --- CENÁRIOS DE FALHA (expected: false) ---
		{"Inválido - Primeiro dígito errado", "12.345.678/0001-85", false},
		{"Inválido - Segundo dígito errado", "12.345.678/0001-96", false},
		{"Inválido - Faltando números (menor que 14)", "1234567800019", false},
		{"Inválido - Sobrando números (maior que 14)", "123456780001950", false},
		{"Inválido - String vazia", "", false},
		{"Inválido - Contém letras", "12.ABC.678/0001-95", false},

		// --- CENÁRIOS DA "LISTA NEGRA" (expected: false) ---
		{"Inválido - Tudo Zero", "00.000.000/0000-00", false},
		{"Inválido - Tudo Zero limpo", "00000000000000", false},
		{"Inválido - Tudo Um", "11111111111111", false},
		{"Inválido - Tudo Nove", "99999999999999", false},
	}

	// 2. O Loop que roda todos os testes mágicamente
	for _, tc := range testCases {
		t.Run(tc.nome, func(t *testing.T) {
			
			// Chama a função pura
			resultado := IsCNPJ(tc.cnpj)
			
			// Compara se o resultado é igual ao que esperávamos
			assert.Equal(t, tc.esperado, resultado, "Falhou no cenário: %s", tc.nome)
		})
	}
}