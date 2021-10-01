//+build generate
//go:generate go run generator.go

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"encoding/json"
)

const (
	blobFileName string = "service/_base/meta.go"
	embedFolder  string = "resources/"
)

var conv = map[string]interface{}{"conv": fmtByteSlice}

var tmpl = template.Must(template.New("").Funcs(conv).Parse(`package _base
// Code generated by go generate; DO NOT EDIT.

import "fmt"

type embedBox struct {
	storage map[string][]byte
}

// Create new box for embed files
func newEmbedBox() *embedBox {
	return &embedBox{storage: make(map[string][]byte)}
}

// Add a file to box
func (e *embedBox) Add(file string, content []byte) {
	e.storage[file] = content
}

// Get file's content
// Always use / for looking up
// For example: /init/README.md is actually configs/init/README.md
func (e *embedBox) Get(file string) []byte {
	if f, ok := e.storage[file]; ok {
		return f
	}
	return nil
}

// Find for a file
func (e *embedBox) Has(file string) bool {
	if _, ok := e.storage[file]; ok {
		return true
	}
	return false
}

// Embed box expose
var box = newEmbedBox()

// Add a file content to box
func Add(file string, content []byte) {
	box.Add(file, content)
}

// Get a file from box
func Get(file string) []byte {
	return box.Get(file)
}

// Has a file in box
func Has(file string) bool {
	return box.Has(file)
}

const (
  AppName="{{.Name}}"
  AppDescription="{{.Description}}"
  AppVersion="{{.Version}}"
  AppCommit="{{.Commit}}"
  AppRepository="{{.Repository}}"
)

var MustAsset = func(file string) []byte {
	var data = box.Get(file)
	if data == nil {
		panic(fmt.Sprintf("resource '%s' not found", file))
	}
	return data
}

func init() {
	{{- range $name, $file := .Files }}
    	box.Add("{{ $name }}", []byte{ {{ conv $file }} })
	{{- end }}
}`),
)

func fmtByteSlice(s []byte) string {
	builder := strings.Builder{}

	for _, v := range s {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}

	return builder.String()
}

type TemplateModel struct {
	Name string
	Description string
	Version string
	Commit string
	Repository string
	Files map[string][]byte
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var name = path.Base(root)
	var description = ""
	var version = "alpha.dev"
	var commit = ""
	var repository = ""

	tmp, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err == nil && tmp != nil && strings.TrimSpace(string(tmp)) != "" {
		repository = strings.TrimSpace(string(tmp))
	}

	if strings.HasPrefix(repository, "git@github.com:") {
		githubRepoInfoUrl := fmt.Sprintf("https://api.github.com/repos/%s", strings.TrimPrefix(strings.TrimSuffix(repository, ".git"), "git@github.com:"))
		request, err := http.NewRequest("GET", githubRepoInfoUrl, nil)
		client := &http.Client{}
		resp, err := client.Do(request)
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				var object map[string]interface{}
				if json.Unmarshal(body, &object) == nil {
					if value, exists := object["description"]; exists {
						description = fmt.Sprint(value)
					}
				}
			}
		}
	}

	tmp, err = ioutil.ReadFile(".description")
	if err == nil {
		description = string(tmp)
	}

	hostname, err := os.Hostname()
	if err == nil {
		me, err := user.Current()
		if err == nil && me != nil {
			commit = fmt.Sprintf("%s@%s", me.Username, hostname)
		}
	}

	tmp, err = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err == nil && tmp != nil {
		commit = strings.TrimSpace(string(tmp))
		version = fmt.Sprintf("%s.dev", commit)

		tmp, err = exec.Command("git", "tag", "--points-at", "HEAD").Output()
		if err == nil && tmp != nil && strings.TrimSpace(string(tmp)) != "" {
			version = strings.TrimSpace(string(tmp))
		}
	}

	// Checking directory with files
	if _, err := os.Stat(embedFolder); os.IsNotExist(err) {
		log.Fatal("Static directory does not exists!")
	}

	// Create map for filenames
	configs := TemplateModel{
		Name: name,
		Description: description,
		Version: version,
		Commit: commit,
		Repository: repository,
		Files: make(map[string][]byte),
	}

	// Walking through embed directory
	err = filepath.Walk(embedFolder, func(path string, info os.FileInfo, err error) error {
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, embedFolder))

		if info.IsDir() {
			// Skip directories
			log.Println(path, "is a directory, skipping...")
			return nil
		} else {
			// If element is a simple file, embed
			log.Println(path, "is a file, packing in...")

			b, err := ioutil.ReadFile(path)
			if err != nil {
				// If file not reading
				log.Printf("Error reading %s: %s", path, err)
				return err
			}

			// Add file name to map
			configs.Files[relativePath] = b
		}

		return nil
	})
	if err != nil {
		log.Fatal("Error walking through embed directory:", err)
	}

	makeBlobFile(configs)
}

func makeBlobFile(configs TemplateModel) {
	// Create blob file
	f, err := os.Create(blobFileName)
	if err != nil {
		log.Fatal("Error creating blob file:", err)
	}
	defer f.Close()

	// Create buffer
	builder := &bytes.Buffer{}

	// Execute template
	if err = tmpl.Execute(builder, configs); err != nil {
		log.Fatal("Error executing template", err)
	}

	// Formatting generated code
	data, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("Error formatting generated code", err)
	}

	// Writing blob file
	if err = ioutil.WriteFile(blobFileName, data, os.ModePerm); err != nil {
		log.Fatal("Error writing blob file", err)
	}
}