// Package plugindemo a demo plugin.
package plugin

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Apidocs map[string]string `json:"apidocs,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	os.Stdout.WriteString("hello create config")
	return &Config{
		Apidocs: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	apidocs  map[string]string
	name     string
	template *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Demo{
		apidocs:  config.Apidocs,
		next:     next,
		name:     name,
		template: template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	log.Println("leck mich ")
	rw.Write([]byte("hello"))
}
