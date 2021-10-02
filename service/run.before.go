package service

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"io/ioutil"
	"service/service/info"
	"sync"
)

var rsaPublic *rsa.PublicKey

func RunBefore(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	data, err := ioutil.ReadFile(fmt.Sprint(info.Args["jwt"]))
	if err != nil {
		panic(err)
	}
	rsaPublic, err = jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
}