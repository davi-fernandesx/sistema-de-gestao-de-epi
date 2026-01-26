package auth

import (

	"errors"
	
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GerarJWT(id int32) (string, error) {

	err := godotenv.Load("configs/.env", "../configs/.env", "../../configs/.env")
	if err != nil {

		log.Println("Aviso: arquivo .env não encontrado. Continuando com variáveis de sistema...")
	}
	secret := os.Getenv("JWT_SECRET")

	if secret  == ""{
		return  "", errors.New("JWT nao configurado")
	}

	claim:= jwt.MapClaims{

		"sub": id, //id do usuario
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 para o token expirar
		"iat": time.Now().Unix(), // data da criação
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return  token.SignedString([]byte(secret))


}
