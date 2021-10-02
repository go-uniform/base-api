package _base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"io/ioutil"
	"net/http"
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
				uniform.Alert(500, "Failed to unmarshal request body")
			}

			for key, value := range requestBody {
				if strValue, ok := value.(string); ok {
					if strings.Contains(strValue, ";base64,") {
						requestBody[key] = strValue[strings.Index(strValue, ";base64,")+8:]
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

var checkBasicPermissions = func(r *http.Request, claims jwt.MapClaims, parameters P, requestBody M) bool {
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
						if !uniform.Contains(links[linkKey], singleValue, false) {
							allMine = false
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

// extract: extract and validate path and query parameters using custom extract routine
func bindHandler(s diary.IPage, timeout time.Duration, topic string, extract func(r *http.Request) P, validateRequest func(data M) M, convertResponse func(response interface{}) []byte, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Scope(local(topic), func(p diary.IPage) {
			defer func() {
				if e := recover(); e != nil {
					// todo: handle custom uniform alert errors, or write generic error message and panic upstream
					w.Header().Set("Message", "Something went wrong")
					w.WriteHeader(500)
					panic(r)
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
			if mineOnly && !filterMine {
				allow := false

				var response interface{}
				if err := c.Request(p, local("permissions.check"), timeout, uniform.Request{
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
						"topic": local("permissions.check"),
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
			if err := c.Request(p, topic, timeout, uniform.Request{
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