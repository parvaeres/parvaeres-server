package parvaeres

import (
	"github.com/pkg/errors"
	"github.com/riccardomc/parvaeres/pkg/gitops"
)

/*CreateDeployment is the flow in response to a Deployment creation request
*
* FIXME: currently we use only repoURL and email, but we probably need a more flexible
* input object
 */
func CreateDeployment(request CreateDeploymentRequest) error {
	err := ValidateCreateDeploymentRequest(request)
	if err != nil {
		return errors.Wrap(err, "CreateDeploymentRequest was invalid")
	}
	existingApplications, err := gitops.ListApplications(request.Email, request.Repository)
	if err != nil {
		return errors.Wrap(err, "CreateDeployment failed")
	}
	if len(existingApplications.Items) > 0 {
		return errors.Errorf("application exists")
	}
	_, err = gitops.CreateApplication(request.Email, request.Repository)
	return err
}

/*ValidateCreateDeploymentRequest FIXME: is not implemented yet
*
 */
func ValidateCreateDeploymentRequest(request CreateDeploymentRequest) error {
	return nil
}
