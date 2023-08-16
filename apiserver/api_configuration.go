/*
 * Microsoft 365 App
 *
 * API to access and configure the Microsoft 365 App
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ConfigurationApiController binds http requests to an api service and writes the service results to the http response
type ConfigurationApiController struct {
	service      ConfigurationApiServicer
	errorHandler ErrorHandler
}

// ConfigurationApiOption for how the controller is set up.
type ConfigurationApiOption func(*ConfigurationApiController)

// WithConfigurationApiErrorHandler inject ErrorHandler into controller
func WithConfigurationApiErrorHandler(h ErrorHandler) ConfigurationApiOption {
	return func(c *ConfigurationApiController) {
		c.errorHandler = h
	}
}

// NewConfigurationApiController creates a default api controller
func NewConfigurationApiController(s ConfigurationApiServicer, opts ...ConfigurationApiOption) Router {
	controller := &ConfigurationApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ConfigurationApiController
func (c *ConfigurationApiController) Routes() Routes {
	return Routes{
		{
			"DeleteConfigurationById",
			strings.ToUpper("Delete"),
			"/v1/configs/{config-id}",
			c.DeleteConfigurationById,
		},
		{
			"GetConfigurationById",
			strings.ToUpper("Get"),
			"/v1/configs/{config-id}",
			c.GetConfigurationById,
		},
		{
			"GetConfigurations",
			strings.ToUpper("Get"),
			"/v1/configs",
			c.GetConfigurations,
		},
		{
			"PostConfiguration",
			strings.ToUpper("Post"),
			"/v1/configs",
			c.PostConfiguration,
		},
		{
			"PutConfigurationById",
			strings.ToUpper("Put"),
			"/v1/configs/{config-id}",
			c.PutConfigurationById,
		},
	}
}

// DeleteConfigurationById - Deletes a configuration
func (c *ConfigurationApiController) DeleteConfigurationById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	configIdParam, err := parseInt64Parameter(params["config-id"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.DeleteConfigurationById(r.Context(), configIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetConfigurationById - Get configuration
func (c *ConfigurationApiController) GetConfigurationById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	configIdParam, err := parseInt64Parameter(params["config-id"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetConfigurationById(r.Context(), configIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetConfigurations - Get configurations
func (c *ConfigurationApiController) GetConfigurations(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetConfigurations(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PostConfiguration - Creates a configuration
func (c *ConfigurationApiController) PostConfiguration(w http.ResponseWriter, r *http.Request) {
	configurationParam := Configuration{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&configurationParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertConfigurationRequired(configurationParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.PostConfiguration(r.Context(), configurationParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PutConfigurationById - Updates a configuration
func (c *ConfigurationApiController) PutConfigurationById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	configIdParam, err := parseInt64Parameter(params["config-id"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	configurationParam := Configuration{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&configurationParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertConfigurationRequired(configurationParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.PutConfigurationById(r.Context(), configIdParam, configurationParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
