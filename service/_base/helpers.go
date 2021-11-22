package _base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"service/service/info"
	"strings"
	"time"
)

var fieldKeyToLinkKey = func(field string) string {
	// todo: find an extensible way to handle this that will pair well with go-fluid
	if len(field) > 0 && !strings.HasSuffix(field, "s") {
		return field + "s"
	}
	return field
}

var checkAuth = func(p diary.IPage, r *http.Request, permissions []string) (string, jwt.MapClaims, bool) {
	authorized := true
	mineOnly := false
	token := ""
	claims := jwt.MapClaims{}

	if len(permissions) > 0 {
		token := extractToken(r)
		claims = verifyToken(p, token, r)
		if len(permissions) > 0 && claims == nil {
			// A notice instead of a warning because otherwise system health alert will be triggered
			p.Notice("api.call.auth.unauthorized", diary.M{
				"path":   r.URL.RequestURI(),
				"method": r.Method,
			})
			uniform.Alert(401, "Please log in before attempting to access this resource.")
		}
		authorized, mineOnly = checkPermissions(claims, permissions...)
	}

	if !authorized {
		p.Warning("api.call.auth.denied", "access has been denied", diary.M{
			"claims":      claims,
			"permissions": permissions,
			"path":        r.URL.RequestURI(),
			"method":      r.Method,
		})
		uniform.Alert(403, "Access has been denied, please contact your system administrator.")
	}

	return token, claims, mineOnly
}

var getRequestBody = func(p diary.IPage, r *http.Request, validateRequest func(uniform.M) uniform.M) uniform.M {
	requestBody := uniform.M{}

	if !uniform.Contains([]string{http.MethodGet}, r.Method, false) {
		// handle file uploads
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
			// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				panic(err)
			}

			// FormFile returns the first file for the given key `file`
			// it also returns the FileHeader so we can get the Filename,
			// the Header and the size of the file
			file, _, err := r.FormFile("file")
			if err == nil && file != nil {
				data, err := ioutil.ReadAll(file)
				if err != nil {
					p.Warning("bind.body.read", "failed to read body", diary.M{
						"path":     r.URL.RequestURI(),
						"errorMsg": err.Error(),
						"error":    err,
					})
					uniform.Alert(500, "Failed to read body")
				}
				_ = file.Close()

				requestBody = uniform.M{
					"file": base64.StdEncoding.EncodeToString(data),
				}
			}
		} else {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				p.Warning("bind.body.read", "failed to read request body", diary.M{
					"path":     r.URL.RequestURI(),
					"errorMsg": err.Error(),
					"error":    err,
				})
				uniform.Alert(500, "Failed to read request body")
			}
			_ = r.Body.Close()

			if err := json.Unmarshal(data, &requestBody); err != nil {
				p.Warning("bind.body.unmarshal", "failed to read request body", diary.M{
					"path":     r.URL.RequestURI(),
					"errorMsg": err.Error(),
					"error":    err,
				})
				uniform.Alert(500, "Failed to parse request body")
			}

			for key, value := range requestBody {
				if strValue, ok := value.(string); ok {
					if strings.Contains(strValue, ";base64,") {
						requestBody[key] = strValue[strings.Index(strValue, ";base64,")+8:]
					} else {
						requestBody[key] = strings.TrimSpace(strValue)
					}
					date, err := time.Parse("2006-01-02", strValue)
					if err == nil {
						requestBody[key] = date.Format(time.RFC3339)
					}
				}
			}
		}

		if validateRequest != nil {
			requestBody = validateRequest(requestBody)
		}
	}

	return requestBody
}

var checkBasicPermissions = func(r *http.Request, claims jwt.MapClaims, parameters uniform.P, requestBody uniform.M) bool {
	// todo: we must be able to clean this up even more just need to think a bit
	userId := fmt.Sprint(claims["id"])
	filterMine := false
	linkKeys := make([]string, 0)
	links := map[string][]string{}
	linkedEntity := ""
	if userId != "" {
		if val, exists := claims["links"]; exists {
			if tmp, ok := val.(map[string]interface{}); ok {
				for key, linkVal := range tmp {
					if strings.Contains(r.URL.Path, fmt.Sprintf("/%s/", key)) || strings.HasSuffix(r.URL.Path, fmt.Sprintf("/%s", key)) {
						linkedEntity = key
					}
					if linkIds, isSlice := linkVal.([]interface{}); isSlice {
						ids := make([]string, 0, len(linkIds))
						for _, id := range linkIds {
							ids = append(ids, fmt.Sprintf("%v", id))
						}
						links[key] = ids
						linkKeys = append(linkKeys, key)
					} else {
						links[key] = []string{fmt.Sprintf("%v", linkVal)}
						linkKeys = append(linkKeys, key)
					}
				}
			}
		}

		allMine := true
		for key, value := range parameters {
			if key != "id" && value == "" {
				continue
			}

			linkKey := fieldKeyToLinkKey(key)
			if value == "me" {
				parameters[key] = userId
				filterMine = true
			} else if strings.HasPrefix(value, "my.") {
				if linkIds, ok := links[strings.TrimPrefix(value, "my.")]; ok {
					parameters[key] = strings.Join(linkIds, ",")
					filterMine = true
				}
			} else if key == "id" && linkedEntity != "" {
				filterMine = uniform.Contains(links[linkedEntity], value, false)
			} else if linkKey != "" {
				for _, singleValue := range strings.Split(value, ",") {
					if !uniform.Contains(links[linkKey], singleValue, false) {
						allMine = false
					}
				}
			} else if key == "-filter-to-allowable" {
				filterLinkKey := fieldKeyToLinkKey(key)
				if filterLinkKey != "" {
					if filterLinkKey == linkedEntity {
						parameters["id"] = strings.Join(links[filterLinkKey], ",")
					} else {
						parameters[value] = strings.Join(links[filterLinkKey], ",")
					}
					filterMine = true
				}
			}
		}
		if allMine {
			filterMine = true
		}

		if requestBody != nil {
			for key, valueInterface := range requestBody {
				value := fmt.Sprint(valueInterface)
				linkKey := fieldKeyToLinkKey(key)

				if value == "me" {
					requestBody[key] = userId
					filterMine = true
				} else if strings.HasPrefix(value, "my.") {
					if linkIds, ok := links[strings.TrimPrefix(value, "my.")]; ok {
						requestBody[key] = strings.Join(linkIds, ",")
						filterMine = true
					}
				} else if linkKey != "" {
					for _, singleValue := range strings.Split(value, ",") {
						if values, ok := links[linkKey]; ok {
							if !uniform.Contains(values, singleValue, false) {
								allMine = false
							}
						}
					}
				}
			}
		}
		if allMine {
			filterMine = true
		}
	}

	return filterMine
}

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

func ExtractIdPathParameter(r *http.Request) uniform.P {
	return uniform.P{"id": mux.Vars(r)["id"]}
}

func verifyToken(p diary.IPage, token string, r *http.Request) jwt.MapClaims {
	webToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return info.GetJwtPublicKey(), nil
	})
	if err != nil {
		p.Warning("verify.jwt.parse", "unable to parse JWT", diary.M{
			"path":  r.URL.RequestURI(),
			"token": token,
		})
		log.Println(err)
		return nil
	}

	claims, ok := webToken.Claims.(jwt.MapClaims)
	if !ok {
		p.Warning("verify.jwt.claims", "JWT claims not ok", diary.M{
			"path":     r.URL.RequestURI(),
			"webToken": webToken,
		})
		return nil
	}

	// todo: validate issuer is from a list of trusted issuers

	if !claims.VerifyAudience(r.Host, true) {
		p.Info("verify.jwt.aud", diary.M{
			"path":     r.URL.RequestURI(),
			"webToken": webToken,
			"host":     r.Host,
			"aud":      claims["aud"],
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

		if uniform.Contains(userPermissions, permission, false) {
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
