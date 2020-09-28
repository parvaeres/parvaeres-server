package parvaeres

/*
* In this file we are supposed to bridge our gitops provider (ArgoCD) with the Parvaeres API
*
 */

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/pkg/errors"
	"github.com/riccardomc/parvaeres/pkg/gitops"
)

//CreateDeployment is the flow in response to a Deployment creation request
func CreateDeployment(request CreateDeploymentRequest, gitops *gitops.GitOpsClient) (*CreateDeploymentResponse, error) {
	err := CreateDeploymentRequestValidate(request)
	if err != nil {
		err = errors.Wrap(err, "CreateDeploymentRequest is invalid")
		return &CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []DeploymentStatus{},
		}, err
	}

	existingApplications, err := gitops.ListApplications(request.Email, request.Repository, request.Path)
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

	application, err := gitops.CreateApplication(request.Email, request.Repository, request.Path)
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
func GetDeploymentByID(deploymentID string, gitops *gitops.GitOpsClient) (status *GetDeploymentResponse, err error) {
	status = &GetDeploymentResponse{
		Error:   true,
		Message: "",
		Items:   []DeploymentStatus{},
	}

	err = GetDeploymentByIDRequestValidate(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeploymentRequest is invalid")
		status.Error = true
		status.Message = err.Error()
		return
	}

	application, err := gitops.GetApplicationByName(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment failed, deployment not found")
		status.Error = true
		status.Message = err.Error()
		return
	}

	deploymentStatus, err := GetDeploymentStatusOfApplication(application)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment deployment created but cannot get status")
		status.Error = true
		status.Message = err.Error()
		return
	}

	err = EnableDeployment(deploymentID, gitops)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment EnableDeployment failed")
		status.Error = true
		status.Message = err.Error()
		return
	}

	deploymentStatus, err = GetDeploymentStatusOfApplication(application)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment deployment created but cannot get status")
		status.Error = true
		status.Message = err.Error()
		return
	}

	status.Error = false
	status.Message = "FOUND"
	status.Items = []DeploymentStatus{*deploymentStatus}
	err = nil
	return
}

//GetDeploymentStatusOfApplication returns the DeploymentStatus corresponding to the application
func GetDeploymentStatusOfApplication(application *v1alpha1.Application) (*DeploymentStatus, error) {
	if application == nil {
		return nil, fmt.Errorf("GetDeploymentStatusOfApplication failed: application is nil")
	}
	deploymentStatus := &DeploymentStatus{
		UUID:     application.GetName(),
		LiveURLs: getExternalURLsOfApplication(application),
		Errors:   getErrorsOfApplication(application),
		Status:   getDeploymentStatusTypeOfApplication(application),
	}
	return deploymentStatus, nil
}

//Return parsed urls, without path
func getExternalURLsOfApplication(application *v1alpha1.Application) (urls []string) {
	urls = []string{}
	if application == nil {
		return
	}

	for _, rawurl := range application.Status.Summary.ExternalURLs {
		u, err := url.Parse(rawurl)
		// FIXME: silently drop malformed URLs, we should log them?
		if err == nil {
			u.Path = ""
			u.RawQuery = ""
			u.User = nil
			urls = append(urls, u.String())
		}
	}

	return
}

func getErrorsOfApplication(application *v1alpha1.Application) (errors []string) {
	errors = []string{}
	if application != nil {
		for _, c := range application.Status.Conditions {
			if c.IsError() {
				errors = append(errors, c.Message)
			}
		}
	}
	return
}

func applicationHasErrorConditions(application *v1alpha1.Application) bool {
	for _, c := range application.Status.Conditions {
		if c.IsError() {
			return true
		}
	}
	return false
}

func getDeploymentStatusTypeOfApplication(application *v1alpha1.Application) DeploymentStatusType {
	// UNKOWN = we don't know much about the application, like when is nil
	if application == nil {
		return UNKNOWN
	}

	// PENDING: if an application is not DEPLOYED or not ERROR
	// then, if it has no SyncPolicy it must be PENDING, i.e. needs user confirmation
	if reflect.DeepEqual(application.Spec.SyncPolicy, &v1alpha1.SyncPolicy{}) {
		return PENDING
	}

	// If Application has a Status it might be DEPLOYED or in ERROR status
	if !reflect.DeepEqual(application.Status, v1alpha1.ApplicationStatus{}) {

		// DEPLOYED: status is Healty and Synced
		if application.Status.Health.Status == health.HealthStatusHealthy ||
			application.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
			return DEPLOYED
		}

		// ERROR: status not Healthy
		if application.Status.Health.Status == health.HealthStatusDegraded ||
			application.Status.Health.Status == health.HealthStatusUnknown ||
			application.Status.Health.Status == health.HealthStatusSuspended ||
			application.Status.Health.Status == health.HealthStatusMissing ||
			applicationHasErrorConditions(application) {
			return ERROR
		}
	}

	// SYNCING: an application has been confirmed by the user and has a SyncPolicy
	if reflect.DeepEqual(application.Spec.SyncPolicy.Automated,
		&v1alpha1.SyncPolicyAutomated{Prune: true, SelfHeal: true}) {
		return SYNCING
	}

	return UNKNOWN
}

//EnableDeployment confirms the deployment to be deployed
func EnableDeployment(deploymentID string, gitops *gitops.GitOpsClient) error {
	application, err := gitops.GetApplicationByName(deploymentID)
	if err != nil {
		return err
	}
	currentStatus := getDeploymentStatusTypeOfApplication(application)
	if currentStatus == PENDING {
		return gitops.SetApplicationAutomatedSyncPolicy(application)
	}
	return nil
}

//GetDeploymentByIDRequestValidate FIXME: is not implemented yet
func GetDeploymentByIDRequestValidate(deploymentID string) error {
	return nil
}
