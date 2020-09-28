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
	"errors"
	"fmt"
	"log"

	"github.com/riccardomc/parvaeres/pkg/email"
	"github.com/riccardomc/parvaeres/pkg/gitops"
)

// DefaultApiService is a service that implents the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
	Gitops           *gitops.GitOpsClient
	EmailProvider    email.EmailProviderInterface
	FeatureFlagEmail bool
	PublicURL        string
}

// NewDefaultApiService creates a default api service
func NewDefaultApiService() DefaultApiServicer {
	return &DefaultApiService{}
}

// DeploymentDeploymentIdGet - Get the deployment with id deploymentId
func (s *DefaultApiService) DeploymentDeploymentIdGet(ctx context.Context, deploymentId string) (interface{}, error) {
	log.Printf("DeploymentDeploymentIdGet: %v", deploymentId)
	response, err := GetDeploymentByID(deploymentId, s.Gitops)
	if err == nil && len(response.Items) > 0 {
		response.Items[0].LogsURL = fmt.Sprintf("%s/v1/deployment/%s/logs", s.PublicURL, deploymentId)
	}
	log.Printf("DeploymentDeploymentIdGet: %v", response)
	return response, nil
}

// DeploymentDeploymentIdLogsGet - Get the deployment with id deploymentId
func (s *DefaultApiService) DeploymentDeploymentIdLogsGet(ctx context.Context, deploymentId string) (interface{}, error) {
	log.Printf("DeploymentDeploymentIdLogsGet: %v", deploymentId)
	response, _ := GetDeploymentLogs(deploymentId, s.Gitops)
	log.Printf("DeploymentDeploymentIdLogsGet: done")
	return response, nil
}

// DeploymentGet - Get all deployments
func (s *DefaultApiService) DeploymentGet(ctx context.Context, getDeploymentRequest GetDeploymentRequest) (interface{}, error) {
	// TODO - update DeploymentGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeploymentGet' not implemented")
}

// DeploymentPost - Create a new deployment
func (s *DefaultApiService) DeploymentPost(ctx context.Context, createDeploymentRequest CreateDeploymentRequest) (interface{}, error) {
	log.Printf("DeploymentPost: %v", createDeploymentRequest)
	response, err := CreateDeployment(createDeploymentRequest, s.Gitops)

	// If deployment is created without errors
	if err == nil {
		// If email communication is enabled
		if s.FeatureFlagEmail {
			id := response.Items[0].UUID
			emailResponse, err := s.EmailProvider.Send(ctx, &email.SendEmailRequest{
				Subject:   "Confirm your application",
				Body:      fmt.Sprintf("%s/v1/deployment/%s", s.PublicURL, id),
				Recipient: createDeploymentRequest.Email,
			})
			if err != nil {
				//FIXME: we should retry
				log.Printf("confirmation email failure: %s", err.Error())
			} else {
				log.Printf("confirmation email success: %s", emailResponse.ID)
			}
		}
		// if email is enabled, we strip ID from API response
		response.Items[0].UUID = ""
	}
	log.Printf("DeploymentPost: %v", response)
	return response, nil
}
