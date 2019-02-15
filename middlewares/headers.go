package middlewares

// Middleware based on https://github.com/unrolled/secure

import (
	"github.com/containous/traefik/log"
	"net/http"

	"github.com/containous/traefik/types"
)

// HeaderOptions is a struct for specifying configuration options for the headers middleware.
type HeaderOptions struct {
	// If Custom request headers are set, these will be added to the request
	CustomRequestHeaders map[string]string
	// If Custom response headers are set, these will be added to the ResponseWriter
	CustomResponseHeaders map[string]string
	ServiceMapping        *types.CowMap
	CustomServiceHeaders  map[string]map[string]string
}

// HeaderStruct is a middleware that helps setup a few basic security features. A single headerOptions struct can be
// provided to configure which features should be enabled, and the ability to override a few of the default values.
type HeaderStruct struct {
	// Customize headers with a headerOptions struct.
	opt HeaderOptions
}

// NewHeaderFromStruct constructs a new header instance from supplied frontend header struct.
func NewHeaderFromStruct(front *types.Frontend) *HeaderStruct {
	headers := front.Headers
	if headers == nil || !headers.HasCustomHeadersDefined() {
		return nil
	}

	return &HeaderStruct{
		opt: HeaderOptions{
			CustomRequestHeaders:  headers.CustomRequestHeaders,
			CustomResponseHeaders: headers.CustomResponseHeaders,
			ServiceMapping:        front.ServiceMapping,
			CustomServiceHeaders:  front.Headers.CustomServiceHeaders,
		},
	}
}

func (s *HeaderStruct) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s.ModifyRequestHeaders(r)
	// If there is a next, call it.
	if next != nil {
		next(w, r)
	}
}

// ModifyRequestHeaders set or delete request headers
func (s *HeaderStruct) ModifyRequestHeaders(r *http.Request) {
	// Loop through Custom request headers
	for header, value := range s.opt.CustomRequestHeaders {
		if value == "" {
			r.Header.Del(header)
		} else {
			r.Header.Set(header, value)
		}
	}
	serviceName := s.opt.ServiceMapping.Get(r.URL.Hostname())
	log.Infof("request service name %s url %v", serviceName, r.URL)
	log.Infof("request host name %s", r.URL.Hostname())
	s.opt.ServiceMapping.LogIt()
	for header, value := range s.opt.CustomServiceHeaders[serviceName] {
		if value == "" {
			log.Infof("deleting header key: %s", header)
			r.Header.Del(header)
		} else {
			log.Infof("adding header key: %s - value: %s", header, value)
			r.Header.Set(header, value)
		}
	}
}

// ModifyResponseHeaders set or delete response headers
func (s *HeaderStruct) ModifyResponseHeaders(res *http.Response) error {
	// Loop through Custom response headers
	for header, value := range s.opt.CustomResponseHeaders {
		if value == "" {
			res.Header.Del(header)
		} else {
			res.Header.Set(header, value)
		}
	}
	return nil
}
