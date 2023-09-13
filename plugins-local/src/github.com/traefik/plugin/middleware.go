// Package plugindemo a demo plugin.
package plugin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	os.Stdout.WriteString("hello create config")
	return &Config{
		Headers: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	name     string
	template *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	os.Stdout.WriteString("hello constructor")

	if len(config.Headers) == 0 {
		os.Stdout.WriteString("config headers 0")
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &Demo{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	log.Println("leck mich ")
	rw.Write([]byte("hello"))
}
