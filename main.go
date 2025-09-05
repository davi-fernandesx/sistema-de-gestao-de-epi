package main

import (
	"fmt"
	"log"

	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/auth"
	"github.com/davi-fernandesx/sistema-de-gestao-de-epi/model"
)

func main(){


	
	//router:= gin.Default()

	//_, err:= configs.InitAplicattion()
	//if err != nil {

	//	log.Fatal(err)
    //	}


	//router.Run(":8080")

	a:= model.Login{
		Nome: "davi",
		Senha: "uee4n43jkj",
	}

	fmt.Println(a.Nome, a.Senha)

	h, err:= auth.HashPassword(a.Senha)
	if err != nil {
		log.Fatal(err)
	}

	i:= fmt.Sprintf("%x", h)

	fmt.Println("hash", i)

	


}