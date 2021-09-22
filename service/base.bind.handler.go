package service

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

func bindHandler(s diary.IPage, timeout time.Duration, topic string, extract func(r *http.Request) (map[string]string, error), convertRequest func(data []byte) (interface{}, error), convertResponse func(data []byte) ([]byte, error), permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Scope(fmt.Sprintf("%s.%s", AppName, topic), func(p diary.IPage) {
			defer func() {
				if e := recover(); e != nil {
					err := fmt.Errorf("%v", r)
					if assertedErr, ok := e.(error); ok {
						err = assertedErr
					}
					p.Warning("bind.catch", "failed to complete api call", diary.M{
						"path":     r.URL.RequestURI(),
						"errorMsg": err.Error(),
						"error":    err,
					})
					w.Header().Set("Message", "Something went wrong")
					w.WriteHeader(500)
				}
			}()

			var err error = nil
			authorized := false
			mineOnly := false
			filterMine := false

			var userId *string = nil
			var userType *string = nil
			token := ""
			claims := jwt.MapClaims{}
			if len(permissions) > 0 {
				token = extractToken(r)
				claims = verifyToken(p, token, r)
				if len(permissions) > 0 && claims == nil {
					// A notice instead of a warning because otherwise system health alert will be triggered
					p.Notice("api.call.auth.unauthorized", diary.M{
						"path":   r.URL.RequestURI(),
						"method": r.Method,
					})
					w.Header().Set("Message", "Please log in before attempting to access this resource.")
					w.WriteHeader(401)
					return
				}
				authorized, mineOnly = checkPermissions(claims, permissions...)
			} else {
				authorized = true
				mineOnly = false
			}

			if !authorized {
				p.Warning("api.call.auth.denied", "access has been denied", diary.M{
					"claims":      claims,
					"permissions": permissions,
					"path":        r.URL.RequestURI(),
					"method":      r.Method,
				})
				w.Header().Set("Message", "Access has been denied, please contact your system administrator.")
				w.WriteHeader(403)
				return
			}

			emptyParameters := make([]string, 10)
			var parameters map[string]string = nil
			if extract != nil {
				// todo: extract should validate all parameters
				parameters, err = extract(r)
				if err != nil {
					p.Warning("bind.extract", "failed to extract", diary.M{
						"path":     r.URL.RequestURI(),
						"errorMsg": err.Error(),
						"error":    err,
					})
					w.WriteHeader(400)
					w.Write([]byte(err.Error()))
					return
				}
			}

			if claims != nil {
				if val, exists := claims["id"]; exists {
					if valStr, ok := val.(string); ok {
						userId = &valStr
					}
				}
				if val, exists := claims["type"]; exists {
					if valStr, ok := val.(string); ok {
						userType = &valStr
					}
				}
			}

			p.Notice("api.call", diary.M{
				"claims":      claims,
				"permissions": permissions,
				"path":        r.URL.RequestURI(),
				"method":      r.Method,
				"parameters":  parameters,
			})

			linkKeys := make([]string, 0)
			links := map[string][]string{}
			linkedEntity := ""
			if userId != nil {
				linkKeys = append(linkKeys, "ids")
				links["ids"] = []string{*userId}

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

					// todo: implement fieldKeyToLinkKey in fluid
					//linkKey := fieldKeyToLinkKey(key)
					linkKey := ""
					if contains(linkKeys, key+"s", false) {
						linkKey = key + "s"
					}

					if value == "me" {
						parameters[key] = *userId
						filterMine = true
					} else if strings.HasPrefix(value, "my.") {
						if linkIds, ok := links[strings.TrimPrefix(value, "my.")]; ok {
							parameters[key] = strings.Join(linkIds, ",")
							filterMine = true
						}
					} else if key == "id" && linkedEntity != "" {
						filterMine = contains(links[linkedEntity], value, false)
					} else if linkKey != "" {
						for _, singleValue := range strings.Split(value, ",") {
							if !contains(links[linkKey], singleValue, false) {
								allMine = false
							}
						}
					} else if key == "-filter-to-allowable" {
						// todo: implement fieldKeyToLinkKey in fluid
						//filterLinkKey := fieldKeyToLinkKey(key)
						filterLinkKey := ""
						if contains(linkKeys, value+"s", false) {
							filterLinkKey = value + "s"
						}
						if filterLinkKey != "" {
							if filterLinkKey == linkedEntity {
								parameters["id"] = strings.Join(links[filterLinkKey], ",")
							} else {
								parameters[value] = strings.Join(links[filterLinkKey], ",")
							}
							filterMine = true
						}
						emptyParameters = append(emptyParameters, "-filter-to-allowable")
					}
				}
				if allMine {
					filterMine = true
				}
			}

			var model interface{} = nil
			if !contains([]string{http.MethodGet}, r.Method, false) {
				if convertRequest != nil {
					data, _ := ioutil.ReadAll(r.Body)
					if err := r.Body.Close(); err != nil {
						p.Warning("bind.ioutil.read-all", "failed to read body", diary.M{
							"path":     r.URL.RequestURI(),
							"errorMsg": err.Error(),
							"error":    err,
						})
						w.WriteHeader(400)
						w.Write([]byte(err.Error()))
						return
					}

					var dataMap map[string]interface{}
					if err := json.Unmarshal(data, &dataMap); err != nil {
						p.Warning("bind.convert", "failed to convert body", diary.M{
							"path":     r.URL.RequestURI(),
							"errorMsg": err.Error(),
							"error":    err,
						})
						w.WriteHeader(400)
						w.Write([]byte(err.Error()))
						return
					}

					// todo: convert should validate model
					if userId != nil {
						allMine := true
						for key, valueInterface := range dataMap {
							value := fmt.Sprint(valueInterface)
							// todo: implement fieldKeyToLinkKey in fluid
							//linkKey := fieldKeyToLinkKey(key)
							linkKey := ""
							if contains(linkKeys, key+"s", false) {
								linkKey = key + "s"
							}

							if value == "me" {
								dataMap[key] = *userId
								filterMine = true
							} else if strings.HasPrefix(value, "my.") {
								if linkIds, ok := links[strings.TrimPrefix(value, "my.")]; ok {
									dataMap[key] = strings.Join(linkIds, ",")
									filterMine = true
								}
							} else if linkKey != "" {
								for _, singleValue := range strings.Split(value, ",") {
									if !contains(links[linkKey], singleValue, false) {
										allMine = false
									}
								}
							}
						}
						if allMine {
							filterMine = true
						}
					}

					for key, value := range dataMap {
						if strValue, ok := value.(string); ok {
							if strings.Contains(strValue, ";base64,") {
								dataMap[key] = strValue[strings.Index(strValue, ";base64,")+8:]
							}
							date, err := time.Parse("2006-01-02", strValue)
							if err == nil {
								dataMap[key] = date.Format(time.RFC3339)
							}
						}
					}

					data, err = json.Marshal(dataMap)
					if err != nil {
						p.Warning("bind.convert", "failed to convert body", diary.M{
							"path":  r.URL.RequestURI(),
							"error": err.Error(),
						})
						w.WriteHeader(400)
						w.Write([]byte(err.Error()))
						return
					}

					model, err = convertRequest(data)
					if err != nil {
						p.Warning("bind.convert", "failed to convert body", diary.M{
							"path":  r.URL.RequestURI(),
							"error": err.Error(),
						})
						w.WriteHeader(400)
						w.Write([]byte(err.Error()))
						return
					}
				} else {
					var data []byte

					var requestModel map[string]interface{}

					if r.Header.Get("Content-Type") == "application/json" {
						bodyData, _ := ioutil.ReadAll(r.Body)
						if err := r.Body.Close(); err != nil {
							p.Warning("bind.ioutil.read-all", "failed to read body", diary.M{
								"path":  r.URL.RequestURI(),
								"error": err.Error(),
							})
							w.WriteHeader(500)
							w.Write([]byte(err.Error()))
							return
						}

						if bodyData != nil && len(bodyData) > 0 {
							if err := json.Unmarshal(bodyData, &requestModel); err != nil {
								p.Warning("bind.unmarshal", "failed to read body", diary.M{
									"path":  r.URL.RequestURI(),
									"error": err.Error(),
									"body":  string(bodyData),
								})
							} else {
								model = requestModel
							}
						}
					} else {
						// Parse our multipart form, 10 << 20 specifies a maximum
						// upload of 10 MB files.
						if err := r.ParseMultipartForm(10 << 20); err != nil {
							panic(err)
						}
						// FormFile returns the first file for the given key `file`
						// it also returns the FileHeader so we can get the Filename,
						// the Header and the size of the file
						file, _, err := r.FormFile("file")
						if err == nil && file != nil {
							data, err = ioutil.ReadAll(file)
							if err != nil {
								p.Warning("bind.ioutil.read-all", "failed to read body", diary.M{
									"path":  r.URL.RequestURI(),
									"error": err.Error(),
								})
								w.WriteHeader(500)
								w.Write([]byte(err.Error()))
								return
							}

							if err := file.Close(); err != nil {
								p.Warning("FormFile.Close", "Failed to close FormFile steam", diary.M{
									"path":  r.URL.RequestURI(),
									"error": err.Error(),
								})
							}

							requestModel["file"] = base64.StdEncoding.EncodeToString(data)
							model = requestModel
						}
					}

					if file, exists := requestModel["file"]; exists {
						if strFile, ok := file.(string); ok {
							if strings.Contains(strFile, "data:application/vnd.ms-excel;base64,") || strings.Contains(strFile, "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,") {
								if strings.Contains(strFile, ";base64,") {
									strFile = strFile[strings.Index(strFile, ";base64,")+8:]
								}

								data, err := base64.StdEncoding.DecodeString(strFile)
								if err != nil {
									panic(err)
								}

								temp, err := excelToCsv(data)
								if err == nil {
									requestModel["file"] = base64.StdEncoding.EncodeToString(temp)
									model = requestModel
								} else {
									p.Notice("convert.excel.failed", diary.M{
										"path":  r.URL.RequestURI(),
										"error": err.Error(),
									})
								}
							}
						}
					}
				}
			}

			context := map[string]interface{}{
				"user-type":  userType,
				"user-id":    userId,
				"token":      token,
			}
			p.Info("bind.context", diary.M{
				"user-type":  userType,
				"user-id":    userId,
				"token":      token,
			})

			if mineOnly && !filterMine {
				allow := false

				if claims != nil {
					var response interface{}
					if err := c.Request(p, fmt.Sprintf("%s.permissions.check", AppName), timeout, uniform.Request{
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
							"error": err.Error(),
							"topic": fmt.Sprintf("%s.permissions.check", AppName),
						})
						w.Header().Set("Message", "Something went wrong")
						w.WriteHeader(500)
						w.Write([]byte(err.Error()))
						return
					}
					if boolValue, ok := response.(bool); ok {
						allow = boolValue
					}
				}

				if !allow {
					p.Warning("api.call.auth.denied", "access to foreign data has been denied", diary.M{
						"claims":      claims,
						"permissions": permissions,
						"path":        r.URL.RequestURI(),
						"method":      r.Method,
					})
					w.Header().Set("Message", "Access data which is not your own has been denied, please contact your system administrator.")
					w.WriteHeader(403)
					return
				}
			}

			for key, value := range parameters {
				if value == "" {
					emptyParameters = append(emptyParameters, key)
				}
			}

			for _, key := range emptyParameters {
				delete(parameters, key)
			}

			var responseParams uniform.P
			var response interface{}
			if err := c.Request(p, topic, timeout, uniform.Request{
				Parameters: parameters,
				Model: model,
				Context: context,
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(response)
				responseParams = r.Parameters()
			}); err != nil {
				p.Warning("bind.request", "failed to execute request", diary.M{
					"topic": topic,
					"path":  r.URL.RequestURI(),
					"error": err.Error(),
				})
				w.Header().Set("Message", "Something went wrong")
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			if convertResponse != nil {
				// todo: convertBack should validate model
				data, err := json.Marshal(response)
				if err != nil {
					p.Warning("bind.convert-back", "failed to convert response back", diary.M{
						"path":  r.URL.RequestURI(),
						"error": err.Error(),
					})
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}

				responseBody, err := convertResponse(data)
				if err != nil {
					p.Warning("bind.convert-back", "failed to convert response back", diary.M{
						"path":  r.URL.RequestURI(),
						"error": err.Error(),
					})
					w.WriteHeader(400)
					w.Write([]byte(err.Error()))
					return
				}

				for key, value := range responseParams {
					w.Header().Set(key, value)
				}

				p.Notice("api.call.success", diary.M{
					"index":               1,
					"claims":              claims,
					"permissions":         permissions,
					"path":                r.URL.RequestURI(),
					"method":              r.Method,
					"parameters":          parameters,
					"response-parameters": responseParams,
				})

				w.Header().Set("content-type", "application/json")
				w.WriteHeader(200)
				w.Write(responseBody)
				return
			}

			if responseParams["-encoding"] == "base64" {
				delete(responseParams, "-encoding")

				if str, ok := response.(string); ok {
					data, err := base64.StdEncoding.DecodeString(str)
					if err != nil {
						p.Warning("bind.base64-decode", "failed to convert response back", diary.M{
							"path":  r.URL.RequestURI(),
							"error": err.Error(),
						})
						w.WriteHeader(500)
						w.Write([]byte(err.Error()))
						return
					}

					if strings.HasSuffix(responseParams["Content-Disposition"], ".csv") {
						switch strings.ToLower(r.Header.Get("Accept")) {
						default:
							data = csvToExcel(data)
							responseParams["Content-Type"] = "application/octet-stream"
							responseParams["Content-Disposition"] = fmt.Sprintf("%s.xlsx", strings.TrimSuffix(responseParams["Content-Disposition"], ".csv"))
							break
						case "text/csv":
							break
						case "application/json":
							data = csvToJson(data)
							responseParams["Content-Type"] = "application/json"
							responseParams["Content-Disposition"] = fmt.Sprintf("%s.json", strings.TrimSuffix(responseParams["Content-Disposition"], ".csv"))
							break
						}

						responseParams["Content-Length"] = fmt.Sprint(len(data))
					}

					for key, value := range responseParams {
						w.Header().Set(key, value)
					}

					p.Notice("api.call.success", diary.M{
						"index":               2,
						"claims":              claims,
						"permissions":         permissions,
						"path":                r.URL.RequestURI(),
						"method":              r.Method,
						"parameters":          parameters,
						"response-parameters": responseParams,
					})

					w.WriteHeader(200)
					w.Write(data)
					return
				}
			}

			for key, value := range responseParams {
				w.Header().Set(key, value)
			}

			responseBody, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				p.Warning("bind.respond", "failed to convert response back", diary.M{
					"path":     r.URL.RequestURI(),
					"errorMsg": err.Error(),
					"error":    err,
				})
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			p.Notice("api.call.success", diary.M{
				"index":               3,
				"claims":              claims,
				"permissions":         permissions,
				"path":                r.URL.RequestURI(),
				"method":              r.Method,
				"parameters":          parameters,
				"response-parameters": responseParams,
			})

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(200)
			w.Write(responseBody)
		}); err != nil {
			panic(err)
		}
	}
}