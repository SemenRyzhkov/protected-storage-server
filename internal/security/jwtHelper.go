package security

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type JwtHelper struct {
	secretKey []byte
}

func New(key string) (*JwtHelper, error) {
	secretKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return &JwtHelper{secretKey}, nil
}

func (j *JwtHelper) GenerateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["userID"] = userID
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "Signing Error", err
	}

	return tokenString, nil
}

func (j *JwtHelper) ExtractClaims(ctx context.Context) (string, error) {
	var tokenString string

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			tokenString = values[0]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return j.secretKey, nil
		})
		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			userID := claims["userID"].(string)
			return userID, nil
		}
	}
	return "unable to extract claims", nil
}
