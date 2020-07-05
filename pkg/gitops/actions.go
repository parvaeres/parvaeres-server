package parvaeres

/* PostApplication implements the flow in response to an Application creation request
*
* FIXME: currently we use only repoURL and email, but we probably need a more flexible
* input object
 */
func PostApplication(email string, repoURL string) error {
	/*existingApplications, err := listApplications(email, repoURL)
	* if err != nil {
	*	return errors.Wrap(err, "PostApplication failed")
	* }
	* if len(existingApplications) != 0 {
	*	return errors.Errorf("application exists")
	* }
	* newApplication, err := newApplication(email, repoURL)
	* if err != nil {
	*	return errors.Wrap(err, "PostApplication failed")
	* }
	* return createApplication(newApplication)
	 */
	return nil
}
