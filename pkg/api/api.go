/*
 * Parvaeres API
 *
 * Parvaeres magic deployment API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package parvaeres

import (
	"context"
	"net/http"
)

// DefaultApiRouter defines the required methods for binding the api requests to a responses for the DefaultApi
// The DefaultApiRouter implementation should parse necessary information from the http request,
// pass the data to a DefaultApiServicer to perform the required actions, then write the service results to the http response.
type DefaultApiRouter interface {
	DeploymentDeploymentIdGet(http.ResponseWriter, *http.Request)
	DeploymentGet(http.ResponseWriter, *http.Request)
	DeploymentPost(http.ResponseWriter, *http.Request)
}

// DefaultApiServicer defines the api actions for the DefaultApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DefaultApiServicer interface {
	DeploymentDeploymentIdGet(context.Context, string) (interface{}, error)
	DeploymentGet(context.Context, GetDeploymentRequest) (interface{}, error)
	DeploymentPost(context.Context, CreateDeploymentRequest) (interface{}, error)
}