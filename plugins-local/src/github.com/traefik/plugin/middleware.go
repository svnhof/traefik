// Package plugindemo a demo plugin.
package plugin

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"log"
	"net/http"
	"os"
	"strings"
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
	log.Printf("Start validating request %s", req.URL.Path)
	var apidoc *string
	for k, v := range a.apidocs {
		if strings.HasPrefix(req.URL.Path, k) {
			apidoc = &v
		}
	}

	if apidoc != nil {
		log.Printf("Found apidoc %s", *apidoc)
		ctx := req.Context()
		loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
		doc, _ := loader.LoadFromFile(*apidoc)
		router, _ := gorillamux.NewRouter(doc)
		route, pathParams, _ := router.FindRoute(req)
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
		}
		err := openapi3filter.ValidateRequest(ctx, requestValidationInput)
		if err != nil {
			switch typedErr := err.(type) {
			case *openapi3filter.RequestError:
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
			case *openapi3filter.SecurityRequirementsError:
				status := http.StatusUnauthorized
				if len(typedErr.Errors) > 0 && typedErr.Errors[0] == openapi3filter.ErrAuthenticationServiceMissing {
					status = http.StatusInternalServerError
				}
				rw.WriteHeader(status)
				rw.Write([]byte(err.Error()))

			default:
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(err.Error()))
			}
			return
		}
		a.next.ServeHTTP(rw, req)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("API Not Found"))
	}

}
