package service

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"io/ioutil"
)

var rsaPublic *rsa.PublicKey

func RunBefore(p diary.IPage) {
	data, err := ioutil.ReadFile(fmt.Sprint(args["jwt"]))
	if err != nil {
		panic(err)
	}
	rsaPublic, err = jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
}