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
	Id     int64  `json:"id"`
	Uid    int64  `json:"uid"`
	Mask   int64  `json:"mask"`
	Issuer string `json:"issuer"`
	Expire int64  `json:"expire"`
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

func GenerateJwt(uid, mask, expire int64) (int64, string, error) {
	claims := &Claims{
		Id:     eg.SId().GetRegionId(),
		Uid:    uid,
		Mask:   mask,
		Expire: expire,
		Issuer: _Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(_JwtKey)
	return claims.Id, str, err
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
