package jwt

import (
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type Claims interface {
	jwtgo.Claims
	ParseToken(token *jwtgo.Token) error
}

var _ jwtgo.StandardClaims

type Jwt struct {
	secret []byte
}

func (j *Jwt) Encode(claims jwtgo.Claims) (string, error) {
	return encode(j.secret, claims)
}

func (j *Jwt) Decode(jwt string, claims Claims) error {
	token, err := decode(j.secret, jwt, claims)
	if err != nil {
		return err
	}

	return claims.ParseToken(token)
}

func New(secretKey string) *Jwt {
	return &Jwt{
		secret: []byte(secretKey),
	}
}

func encode(secretKey []byte, claims jwtgo.Claims) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	tok, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func decode(secretKey []byte, jwtValue string, claims jwtgo.Claims) (*jwtgo.Token, error) {
	return jwtgo.ParseWithClaims(jwtValue, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
}
