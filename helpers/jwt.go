package helpers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ClaimsOperation interface {
	CreateToken(string) (string, error)
	GetToken(*http.Request) (string, error)
	VerifyToken(string, string) (*Claims, error)
}

//type assertion
var _ ClaimsOperation = (*Claims)(nil)

type Claims struct {
	Id       uint
	Username string
	Role     string
	jwt.StandardClaims
}

// NewClaim is the constructor of claim ...
func NewClaim(id uint, username, role string) *Claims {
	return &Claims{
		Id:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3000 * time.Minute).Unix(),
		},
	}
}

func (claims *Claims) CreateToken(tokenKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenEncodedString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}

	return tokenEncodedString, nil
}

func (claims *Claims) GetToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("no request token")
	}
	if len(token) > 7 {
		splitToken := strings.Split(token, "Bearer ")
		token = splitToken[1]
		return token, nil
	} else {
		return "", errors.New("could not get token string")
	}
}

func (claims *Claims) VerifyToken(userToken, tokenKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		userToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetVerifiedToken(tokenKey string, r *http.Request) (*Claims, error) {
	claims := &Claims{}
	requestToken, err := claims.GetToken(r)
	if err != nil {
		return nil, err
	}

	//verify token
	verifiedToken, err := claims.VerifyToken(requestToken, tokenKey)
	if err != nil {
		return nil, err
	}

	return verifiedToken, nil
}
