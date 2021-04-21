package Service

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService interface {
	GenerateToken(Uid primitive.ObjectID) string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Uid primitive.ObjectID `json:"Uid"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "ujjwalgautam",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (jwtSrv *jwtService) GenerateToken(userid primitive.ObjectID) string {

	// Set custom and standard claims
	claims := &jwtCustomClaims{
		userid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    jwtSrv.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(jwtSrv.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(jwtSrv.secretKey), nil
	})
}
