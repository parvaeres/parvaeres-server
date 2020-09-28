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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A DefaultApiController binds http requests to an api service and writes the service results to the http response
type DefaultApiController struct {
	service DefaultApiServicer
}

// NewDefaultApiController creates a default api controller
func NewDefaultApiController(s DefaultApiServicer) Router {
	return &DefaultApiController{service: s}
}

// Routes returns all of the api route for the DefaultApiController
func (c *DefaultApiController) Routes() Routes {
	return Routes{
		{
			"DeploymentDeploymentIdDelete",
			strings.ToUpper("Delete"),
			"/v1/deployment/{deploymentId}",
			c.DeploymentDeploymentIdDelete,
		},
		{
			"DeploymentDeploymentIdGet",
			strings.ToUpper("Get"),
			"/v1/deployment/{deploymentId}",
			c.DeploymentDeploymentIdGet,
		},
		{
			"DeploymentGet",
			strings.ToUpper("Get"),
			"/v1/deployment",
			c.DeploymentGet,
		},
		{
			"DeploymentPost",
			strings.ToUpper("Post"),
			"/v1/deployment",
			c.DeploymentPost,
		},
	}
}

// DeploymentDeploymentIdDelete - Delete the deployment with id deploymentId
func (c *DefaultApiController) DeploymentDeploymentIdDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentId := params["deploymentId"]
	result, err := c.service.DeploymentDeploymentIdDelete(r.Context(), deploymentId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// DeploymentDeploymentIdGet - Get the deployment with id deploymentId
func (c *DefaultApiController) DeploymentDeploymentIdGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentId := params["deploymentId"]
	result, err := c.service.DeploymentDeploymentIdGet(r.Context(), deploymentId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// DeploymentGet - Get all deployments
func (c *DefaultApiController) DeploymentGet(w http.ResponseWriter, r *http.Request) {
	getDeploymentRequest := &GetDeploymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&getDeploymentRequest); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.DeploymentGet(r.Context(), *getDeploymentRequest)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// DeploymentPost - Create a new deployment
func (c *DefaultApiController) DeploymentPost(w http.ResponseWriter, r *http.Request) {
	createDeploymentRequest := &CreateDeploymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createDeploymentRequest); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.DeploymentPost(r.Context(), *createDeploymentRequest)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
