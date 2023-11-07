/*
 * Microsoft 365 App
 *
 * API to access and configure the Microsoft 365 App
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ProxyAPIController binds http requests to an api service and writes the service results to the http response
type ProxyAPIController struct {
	service      ProxyAPIServicer
	errorHandler ErrorHandler
}

// ProxyAPIOption for how the controller is set up.
type ProxyAPIOption func(*ProxyAPIController)

// WithProxyAPIErrorHandler inject ErrorHandler into controller
func WithProxyAPIErrorHandler(h ErrorHandler) ProxyAPIOption {
	return func(c *ProxyAPIController) {
		c.errorHandler = h
	}
}

// NewProxyAPIController creates a default api controller
func NewProxyAPIController(s ProxyAPIServicer, opts ...ProxyAPIOption) Router {
	controller := &ProxyAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ProxyAPIController
func (c *ProxyAPIController) Routes() Routes {
	return Routes{
		"MsproxyMsGraphPathDelete": Route{
			strings.ToUpper("Delete"),
			"/v1/msproxy/{ms-graph-path}",
			c.MsproxyMsGraphPathDelete,
		},
		"MsproxyMsGraphPathGet": Route{
			strings.ToUpper("Get"),
			"/v1/msproxy/{ms-graph-path}",
			c.MsproxyMsGraphPathGet,
		},
		"MsproxyMsGraphPathPost": Route{
			strings.ToUpper("Post"),
			"/v1/msproxy/{ms-graph-path}",
			c.MsproxyMsGraphPathPost,
		},
		"MsproxyMsGraphPathPut": Route{
			strings.ToUpper("Put"),
			"/v1/msproxy/{ms-graph-path}",
			c.MsproxyMsGraphPathPut,
		},
	}
}

// MsproxyMsGraphPathDelete - A proxy server that passes requests to the Microsoft Graph API
func (c *ProxyAPIController) MsproxyMsGraphPathDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	msGraphPathParam := params["ms-graph-path"]
	elionaProjectIdParam := r.Header.Get("eliona-project-id")
	result, err := c.service.MsproxyMsGraphPathDelete(r.Context(), msGraphPathParam, elionaProjectIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// MsproxyMsGraphPathGet - A proxy server that passes requests to the Microsoft Graph API
func (c *ProxyAPIController) MsproxyMsGraphPathGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	msGraphPathParam := params["ms-graph-path"]
	elionaProjectIdParam := r.Header.Get("eliona-project-id")
	result, err := c.service.MsproxyMsGraphPathGet(r.Context(), msGraphPathParam, elionaProjectIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// MsproxyMsGraphPathPost - A proxy server that passes requests to the Microsoft Graph API
func (c *ProxyAPIController) MsproxyMsGraphPathPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	msGraphPathParam := params["ms-graph-path"]
	elionaProjectIdParam := r.Header.Get("eliona-project-id")
	result, err := c.service.MsproxyMsGraphPathPost(r.Context(), msGraphPathParam, elionaProjectIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// MsproxyMsGraphPathPut - A proxy server that passes requests to the Microsoft Graph API
func (c *ProxyAPIController) MsproxyMsGraphPathPut(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	msGraphPathParam := params["ms-graph-path"]
	elionaProjectIdParam := r.Header.Get("eliona-project-id")
	result, err := c.service.MsproxyMsGraphPathPut(r.Context(), msGraphPathParam, elionaProjectIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
