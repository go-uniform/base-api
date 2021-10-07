package _base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/info"
	"strconv"
	"strings"
	"time"
)

// extract: extract and validate path and query parameters using custom extract routine
func BindHandler(s diary.IPage, timeout time.Duration, topic string, extract func(r *http.Request) uniform.P, validateRequest func(data uniform.M) uniform.M, convertResponse func(response interface{}) []byte, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Scope(topic, func(p diary.IPage) {
			defer func() {
				if e := recover(); e != nil {
					errMsg := fmt.Sprint(e)
					if strings.HasPrefix(errMsg, "alert:") {
						packet := strings.TrimLeft(errMsg, "alert:")
						codeStr := packet[0:3]
						if code, err := strconv.Atoi(codeStr); err == nil {
							w.Header().Set("Message", packet[4:])
							w.WriteHeader(code)
							return
						}
					} else if strings.HasPrefix(errMsg, "validation:") {
						w.Header().Set("Message", "One or more fields failed validation")
						w.WriteHeader(400)
						_, _ = w.Write([]byte(strings.TrimLeft(errMsg, "validation:")))
						return
					}
					w.Header().Set("Message", "Something went wrong")
					w.WriteHeader(500)
					p.Error("unknown", errMsg, diary.M{
						"error": e,
					})
					return
				}
			}()

			p.Info("api.call.enter", diary.M{
				"permissions": permissions,
				"path":        r.URL.RequestURI(),
				"method":      r.Method,
			})

			// check that current user is authorized to access endpoint and extract token claims
			token, claims, mineOnly := checkAuth(p, r, permissions)
			if claims == nil {
				claims = jwt.MapClaims{}
			}

			// extract and validate path and query parameters using custom extract routine
			var parameters uniform.P
			if extract != nil {
				parameters = extract(r)
			}

			p.Notice("api.call", diary.M{
				"claims":      claims,
				"permissions": permissions,
				"path":        r.URL.RequestURI(),
				"method":      r.Method,
				"parameters":  parameters,
			})

			// extract request body, if available
			requestBody := getRequestBody(p, r, validateRequest)
			p.Debug("request.body", requestBody)

			// performance basic permissions check with know high level links
			filterMine := checkBasicPermissions(r, claims, parameters, requestBody)

			// clean up empty or irrelevant parameters
			removeParameters := make([]string, 0)
			for key, value := range parameters {
				if value == "" || value == "-filter-to-allowable" {
					removeParameters = append(removeParameters, key)
				}
			}
			for _, key := range removeParameters {
				delete(parameters, key)
			}

			// formulate request context
			context := uniform.M{
				"user-type":  fmt.Sprint(claims["type"]),
				"user-id":    fmt.Sprint(claims["id"]),
				"token":      token,
			}
			p.Debug("context", context)

			// do a deep dive permission check (special permissions)
			if permissions != nil && len(permissions) > 0 && mineOnly && !filterMine {
				allow := false

				var response interface{}
				if err := info.Conn.Request(p, TargetLocal("permissions.check"), timeout, uniform.Request{
					Parameters: parameters,
					Model: map[string]interface{}{
						"method": r.Method,
						"path":   r.URL.RequestURI(),
						"permissions": map[string]interface{}{
							"required":  permissions,
							"inverted":  claims["permissions.inverted"],
							"available": claims["permissions.tags"],
						},
					},
					Context: context,
				}, func(r uniform.IRequest, p diary.IPage) {
					if r.HasError() {
						panic(r.Error())
					}
					r.Read(response)
				}); err != nil {
					p.Warning("bind.permissions.check", "failed to execute deep permissions check", diary.M{
						"path":  r.URL.RequestURI(),
						"errorMsg": err.Error(),
						"error":    err,
						"topic": TargetLocal("permissions.check"),
					})
					uniform.Alert(500, "Failed to execute deep permissions check")
				}

				if boolValue, ok := response.(bool); ok {
					allow = boolValue
				}

				if !allow {
					p.Warning("api.call.auth.denied", "access to foreign data has been denied", diary.M{
						"claims":      claims,
						"permissions": permissions,
						"path":        r.URL.RequestURI(),
						"method":      r.Method,
					})
					uniform.Alert(403, "Access data which is not your own has been denied, please contact your system administrator.")
				}
			}

			// push request into pub/sub backbone for processing
			var responseHeaders uniform.P
			var response interface{}
			if err := info.Conn.Request(p, topic, timeout, uniform.Request{
				Parameters: parameters,
				Model: requestBody,
				Context: context,
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(response)
				// response parameters will act as response headers since HTTP responses can't have parameters
				responseHeaders = r.Parameters()
			}); err != nil {
				panic(err)
			}
			p.Debug("response.model", response)
			p.Debug("response.headers", responseHeaders)

			// convert the response
			var err error
			var responseData []byte
			if convertResponse != nil {
				responseData = convertResponse(response)
			}

			// handle a file download response
			if responseHeaders["-encoding"] == "base64" {
				delete(responseHeaders, "-encoding")

				ok := false
				responseData, ok = response.([]byte)
				if !ok {
					// if response data is not in a raw form then decode it as a base64 encoded string
					responseData, err = base64.StdEncoding.DecodeString(fmt.Sprint(response))
					if err != nil {
						p.Warning("bind.base64-decode", "failed to convert response back", diary.M{
							"path":     r.URL.RequestURI(),
							"errorMsg": err.Error(),
							"error":    err,
						})
						uniform.Alert(500, "Failed to convert response back")
					}
				}

				// set output length for download logic
				responseHeaders["Content-Length"] = fmt.Sprint(len(responseData))
			} else {
				// convert raw response model to json format by default
				responseData, err = json.Marshal(response)
				if err != nil {
					p.Warning("bind.respond", "failed to convert response model into data", diary.M{
						"path":     r.URL.RequestURI(),
						"errorMsg": err.Error(),
						"error":    err,
					})
					uniform.Alert(500, "Failed to convert response model into data")
				}
			}

			// set HTTP response headers
			if _, exists := responseHeaders["Content-Type"]; !exists {
				responseHeaders["Content-Type"] = "application/json"
			}
			for key, value := range responseHeaders {
				w.Header().Set(key, value)
			}

			// write response data back to the requester
			p.Notice("api.call.success", diary.M{
				"claims":              claims,
				"permissions":         permissions,
				"path":                r.URL.RequestURI(),
				"method":              r.Method,
				"parameters":          parameters,
				"response-parameters": responseHeaders,
			})
			p.Debug("response.body", responseData)
			w.WriteHeader(200)
			_, _ = w.Write(responseData)
		}); err != nil {
			panic(err)
		}
	}
}