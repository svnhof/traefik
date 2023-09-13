// Package plugindemo a demo plugin.
package plugin

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
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
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, _ := loader.LoadFromFile("plugins-local/src/github.com/traefik/plugin/petstore-v1.yaml")
	router, _ := gorillamux.NewRouter(doc)
	route, pathParams, _ := router.FindRoute(req)
	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
	}
	_ := openapi3filter.ValidateRequest(ctx, requestValidationInput)

	os.Stdout.WriteString("hello sdfsdfs")
	rw.Write([]byte("hello"))
}
