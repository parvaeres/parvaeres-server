package parvaeres

/*
* In this file we are supposed to bridge our gitops provider (ArgoCD) with the Parvaeres API
*
 */

import (
	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/pkg/errors"
	"github.com/riccardomc/parvaeres/pkg/gitops"
)

//CreateDeployment is the flow in response to a Deployment creation request
func CreateDeployment(request CreateDeploymentRequest) (*CreateDeploymentResponse, error) {
	err := CreateDeploymentRequestValidate(request)
	if err != nil {
		err = errors.Wrap(err, "CreateDeploymentRequest is invalid")
		return &CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	existingApplications, err := gitops.ListApplications(request.Email, request.Repository)
	if err != nil {
		err = errors.Wrap(err, "CreateDeployment failed")
		return &CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}
	if len(existingApplications.Items) > 0 {
		err = errors.Errorf("application exists")
		return &CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	application, err := gitops.CreateApplication(request.Email, request.Repository)
	if err != nil {
		err = errors.Wrap(err, "CreateDeployment failed")
		return &CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	status, err := GetDeploymentStatusOfApplication(application)
	if err != nil {
		//FIXME: handle cannot get status: the app is created but status cannot be retrieved
	}

	return &CreateDeploymentResponse{
		Error:   false,
		Message: "CREATED",
		Items:   []DeploymentStatus{*status},
	}, nil
}

//CreateDeploymentRequestValidate FIXME: is not implemented yet
func CreateDeploymentRequestValidate(request CreateDeploymentRequest) error {
	return nil
}

//GetDeploymentByID returns a list of deployments based on the request
func GetDeploymentByID(deploymentID string) (*GetDeploymentResponse, error) {
	err := GetDeploymentByIDRequestValidate(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeploymentRequest is invalid")
		return &GetDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	application, err := gitops.GetApplicationByName(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment failed")
		return &GetDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	status, err := GetDeploymentStatusOfApplication(application)
	if err != nil {
		//FIXME: handle cannot get status: the app is created but status cannot be retrieved
	}

	return &GetDeploymentResponse{
		Error:   false,
		Message: "CREATED",
		Items:   []DeploymentStatus{*status},
	}, nil
}

//GetDeploymentStatusOfApplication returns the DeploymentStatus corresponding to the application
func GetDeploymentStatusOfApplication(application *v1alpha1.Application) (*DeploymentStatus, error) {
	deploymentStatus := &DeploymentStatus{
		UUID:     application.GetName(),
		LiveURLs: []string{},
		Errors:   []string{},
		Status:   "PENDING",
	}
	return deploymentStatus, nil
}

//GetDeploymentByIDRequestValidate FIXME: is not implemented yet
func GetDeploymentByIDRequestValidate(deploymentID string) error {
	return nil
}
