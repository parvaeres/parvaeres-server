package gitops

import (
	"github.com/pkg/errors"
	parvaeres "github.com/riccardomc/parvaeres/pkg/api"
)

/*CreateDeployment is the flow in response to a Deployment creation request
*
* FIXME: currently we use only repoURL and email, but we probably need a more flexible
* input object
 */
func CreateDeployment(request parvaeres.CreateDeploymentRequest) error {
	err := ValidateCreateDeploymentRequest(request)
	if err != nil {
		return errors.Wrap(err, "CreateDeploymentRequest was invalid")
	}
	existingApplications, err := listApplications(request.Email, request.Repository)
	if err != nil {
		return errors.Wrap(err, "CreateDeployment failed")
	}
	if len(existingApplications.Items) > 0 {
		return errors.Errorf("application exists")
	}
	newApplication, err := newApplication(request.Email, request.Repository)
	if err != nil {
		return errors.Wrap(err, "CreateDeployment failed")
	}
	return createApplication(newApplication)
}

/*ValidateCreateDeploymentRequest FIXME: is not implemented yet
*
 */
func ValidateCreateDeploymentRequest(request parvaeres.CreateDeploymentRequest) error {
	return nil
}
