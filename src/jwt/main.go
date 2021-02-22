package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

//GenerateJWT is to create aceess token.
func GenerateJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
