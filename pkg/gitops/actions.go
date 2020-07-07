package gitops

import "github.com/pkg/errors"

/*PostApplication implements the flow in response to an Application creation request
*
* FIXME: currently we use only repoURL and email, but we probably need a more flexible
* input object
 */
func PostApplication(repoURL string, email string) error {
	existingApplications, err := listApplications(email, repoURL)
	if err != nil {
		return errors.Wrap(err, "PostApplication failed")
	}
	if existingApplications.Size() > 0 {
		return errors.Errorf("application exists")
	}
	newApplication, err := newApplication(email, repoURL)
	if err != nil {
		return errors.Wrap(err, "PostApplication failed")
	}
	return createApplication(newApplication)
}
