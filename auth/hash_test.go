package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestHashPassword_Success(t *testing.T) {
	// Arrange (Preparação)
	senha := "minhaSenhaSuperSegura123"

	//  (Execução)
	hash, err := HashPassword(senha)

	
	// Verificando se não houve erro na execução
	require.NoError(t, err)
	// retorno não está vazio
	require.NotEmpty(t, hash)

	// 3. A verificação mais importante: o hash NÃO É a senha original
	assert.NotEqual(t, []byte(senha), hash)

 /* comparando a senha com o hash gerado, se der certo, o teste passa*/
	err = HashCompare(hash, []byte(senha))
	assert.NoError(t, err)
}