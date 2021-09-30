package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func extractToken(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		log.Println("Authorization header is empty")
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		log.Println("Authorization header is not set to bearer")
		return ""
	}
	return authorization[7:]
}

func extractIdPathParameter(r *http.Request) P {
	return P{ "id": mux.Vars(r)["id"] }
}

func verifyToken(p diary.IPage, token string, r *http.Request) jwt.MapClaims {
	webToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublic, nil
	})
	if err != nil {
		p.Warning("verify.jwt.parse", "unable to parse JWT", diary.M{
			"path": r.URL.RequestURI(),
			"token": token,
		})
		log.Println(err)
		return nil
	}

	claims, ok := webToken.Claims.(jwt.MapClaims)
	if !ok {
		p.Warning("verify.jwt.claims", "JWT claims not ok", diary.M{
			"path": r.URL.RequestURI(),
			"webToken": webToken,
		})
		return nil
	}

	// todo: validate issuer is from a list of trusted issuers

	if !claims.VerifyAudience(r.Host, true) {
		p.Info("verify.jwt.aud", diary.M{
			"path": r.URL.RequestURI(),
			"webToken": webToken,
			"host": r.Host,
			"aud": claims["aud"],
		})
		return nil
	}

	if !claims.VerifyNotBefore(time.Now().Unix(), false) {
		p.Info("verify.jwt.before", diary.M{
			"webToken": webToken,
		})
		return nil
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		p.Info("verify.jwt.expired", diary.M{
			"webToken": webToken,
		})
		return nil
	}

	return claims
}

// checks permissions and returns if you are authorized and if you may only access your own data
func checkPermissions(claims jwt.MapClaims, permissions ...string) (authorized bool, mineOnly bool) {
	if len(permissions) == 0 {
		return true, false
	}

	if claims == nil || len(claims) == 0 {
		return false, false
	}

	inverted := false
	if assertedValue, ok := claims["permissions.inverted"].(bool); ok {
		inverted = assertedValue
	} else {
		panic("unable to retrieve permissions.inverted key from JWT")
	}

	userPermissions := make([]string, 0)
	if assertedValue, ok := claims["permissions.tags"].([]interface{}); ok {
		data := make([]string, len(assertedValue))
		for i, d := range assertedValue {
			data[i] = fmt.Sprint(d)
		}
		userPermissions = data
	} else if assertedValue, ok := claims["permissions.tags"].([]string); ok {
		userPermissions = assertedValue
	}

	for _, permission := range permissions {
		if strings.HasPrefix(permission, "my.") {
			continue
		}

		if contains(userPermissions, permission, false) {
			return !inverted, false
		}
	}

	for _, permission := range permissions {
		if !strings.HasPrefix(permission, "my.") {
			continue
		}

		if containsWildcard(userPermissions, permission) {
			return !inverted, true
		}
	}

	return inverted, false
}

func containsWildcard(s []string, e string) bool {
	for _, a := range s {
		t := e
		if len(a) < len(e) {
			t = e[:len(a)]
		}
		if strings.ToLower(a) == strings.ToLower(t) {
			return true
		}
	}
	return false
}
