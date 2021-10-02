package info

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"sync"
)

// storage for jwt public key which we will use to check jwt validity
var jwtPublicKey *rsa.PublicKey

// a mutex used to ensure jwt public key population is thread-safe
var lock = sync.Mutex{}

func GetJwtPublicKey() *rsa.PublicKey {
	// thread safety, only one thread may continue at a time
	lock.Lock()
	defer lock.Unlock()

	if jwtPublicKey == nil {
		data, err := ioutil.ReadFile(fmt.Sprint(Args["jwt"]))
		if err != nil {
			panic(err)
		}
		jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(data)
		if err != nil {
			panic(err)
		}
	}

	return jwtPublicKey
}