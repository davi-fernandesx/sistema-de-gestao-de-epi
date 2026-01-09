package entregaepi

import (
	
	"testing"

	integracao "github.com/davi-fernandesx/sistema-de-gestao-de-epi/repository/Integracao"
)




func TestEntrega(t *testing.T){

	db:= integracao.SetupTestDB(t)
	defer db.Close()


}