package helper

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"go-crud-api/models"
)

func loadSecretKey() []byte {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("SECRET_KEY")

	return []byte(secret)
}

func GenerateToken(user models.User) (string, error) {

	secret := loadSecretKey()

	claim := jwt.MapClaims{}
	claim["id"] = user.Id

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, err
}
