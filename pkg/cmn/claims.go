package cmn

import (
	"errors"
	"github.com/95eh/eg"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	_JwtKey = []byte("95eh.com")
	_Issuer = "easy game"
)

type Claims struct {
	Uid    string `json:"uid"`
	Issuer string `json:"issuer"`
}

func (c Claims) Valid() error {
	if c.Issuer != _Issuer {
		return errors.New("wrong issuer")
	}
	return nil
}

func SetJwtKey(key string) {
	_JwtKey = []byte(key)
}

func SetJwtIssuer(issuer string) {
	_Issuer = issuer
}

func GenerateJwt(uid string) (str string, err error) {
	claims := &Claims{
		Uid:    uid,
		Issuer: _Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err = token.SignedString(_JwtKey)
	return
}

func ParseJwt(str string) (tkn *jwt.Token, c Claims, err eg.IErr) {
	var e error
	tkn, e = jwt.ParseWithClaims(str, &c, func(token *jwt.Token) (interface{}, error) {
		return _JwtKey, nil
	})
	if e != nil {
		err = eg.WrapErr(eg.EcUnAuth, e)
	}
	return
}
