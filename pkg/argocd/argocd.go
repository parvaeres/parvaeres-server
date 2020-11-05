package argocd

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"reflect"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	argoclient "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/google/uuid"
	parvaeres "github.com/parvaeres/parvaeres/pkg/api"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//ArgoCD represents an ArgoCD friend
type ArgoCD struct {
	ArgoCDclient     argoclient.Interface
	ArgoCDNamespace  string
	KubernetesClient kubernetes.Interface
}

func GetKubernetesConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		return rest.InClusterConfig()
	}
	log.Printf("using configuration from '%s'", kubeconfig)
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func GetKubernetesClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func GetArgoCDClientSet(config *rest.Config) (*argoclient.Clientset, error) {
	return argoclient.NewForConfig(config)
}

//CreateDeployment is the flow in response to a Deployment creation request
func (a *ArgoCD) CreateDeployment(request parvaeres.CreateDeploymentRequest) (*parvaeres.CreateDeploymentResponse, error) {
	existingApplications, err := a.listApplications(request.Email, request.Repository, request.Path)
	if err != nil {
		err = errors.Wrap(err, "CreateDeployment failed")
		return &parvaeres.CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []parvaeres.DeploymentStatus{},
		}, err
	}
	if len(existingApplications.Items) > 0 {
		err = errors.Errorf("application exists")
		return &parvaeres.CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []parvaeres.DeploymentStatus{},
		}, err
	}

	application, err := a.createApplication(request.Email, request.Repository, request.Path)
	if err != nil {
		err = errors.Wrap(err, "CreateDeployment failed")
		return &parvaeres.CreateDeploymentResponse{
			Error:   true,
			Message: err.Error(),
			Items:   []parvaeres.DeploymentStatus{},
		}, err
	}

	status, err := a.getDeploymentStatusOfApplication(application)
	if err != nil {
		//FIXME: handle cannot get status: the app is created but status cannot be retrieved
	}

	return &parvaeres.CreateDeploymentResponse{
		Error:   false,
		Message: "CREATED",
		Items:   []parvaeres.DeploymentStatus{*status},
	}, nil
}

//EnableDeployment confirms the deployment to be deployed
func (a *ArgoCD) EnableDeployment(deploymentID string) error {
	application, err := a.getApplicationByName(deploymentID)
	if err != nil {
		return err
	}
	currentStatus := getDeploymentStatusTypeOfApplication(application)
	if currentStatus == parvaeres.PENDING {
		return a.setApplicationAutomatedSyncPolicy(application)
	}
	return nil
}

//GetDeploymentByID gives the deployment with deploymentID
func (a *ArgoCD) GetDeploymentByID(deploymentID string) (status *parvaeres.GetDeploymentResponse, err error) {
	status = &parvaeres.GetDeploymentResponse{
		Error:   true,
		Message: "",
		Items:   []parvaeres.DeploymentStatus{},
	}

	application, err := a.getApplicationByName(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment failed, deployment not found")
		status.Error = true
		status.Message = err.Error()
		return
	}

	deploymentStatus, err := a.getDeploymentStatusOfApplication(application)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment deployment created but cannot get status")
		status.Error = true
		status.Message = err.Error()
		return
	}

	err = a.EnableDeployment(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment EnableDeployment failed")
		status.Error = true
		status.Message = err.Error()
		return
	}

	deploymentStatus, err = a.getDeploymentStatusOfApplication(application)
	if err != nil {
		err = errors.Wrap(err, "GetDeployment deployment created but cannot get status")
		status.Error = true
		status.Message = err.Error()
		return
	}

	status.Error = false
	status.Message = "FOUND"
	status.Items = []parvaeres.DeploymentStatus{*deploymentStatus}
	err = nil
	return
}

//GetDeploymentLogs returns the logs of an application
func (a *ArgoCD) GetDeploymentLogs(deploymentID string) (response *parvaeres.GetDeploymentLogsResponse, err error) {
	response = &parvaeres.GetDeploymentLogsResponse{
		Error:   false,
		Message: "",
		Items:   []parvaeres.Logs{},
	}
	application, err := a.getApplicationByName(deploymentID)
	if err != nil {
		err = errors.Wrap(err, "GetDeploymentLogs failed, deployment not found")
		response.Error = true
		response.Message = err.Error()
		return
	}

	// get all pods and containers FIXME: we should check StatefulSets, DaemonSets, etc.
	for _, r := range application.Status.Resources {
		if r.Kind == "Deployment" {
			d, err := a.KubernetesClient.AppsV1().Deployments(r.Namespace).Get(r.Name, metav1.GetOptions{})
			if err != nil {
				log.Printf("cannot get deployment: %s: %v", r.Name, err)
			}
			labelMap, _ := metav1.LabelSelectorAsMap(d.Spec.Selector)
			selector := labels.SelectorFromSet(labelMap).String()
			log.Printf("getting pods for deployment: %s using selector: %s", r.Name, selector)
			pods, err := a.KubernetesClient.CoreV1().Pods(r.Namespace).List(
				metav1.ListOptions{
					LabelSelector: selector,
				})
			if err != nil {
				log.Printf("cannot get pods: %s: %v", r.Name, err)
			}
			log.Printf("fetching logs for %d pods", len(pods.Items))
			for _, p := range pods.Items {
				for _, c := range p.Spec.Containers {
					log.Printf("fetching logs for pod: %s, container: %s", p.Name, c.Name)
					req := a.KubernetesClient.CoreV1().Pods(r.Namespace).GetLogs(
						p.Name,
						&corev1.PodLogOptions{
							Container: c.Name,
						},
					)
					// FIXME: error handling
					logs, err := req.Stream()
					if err != nil {
						log.Printf("error fetching logs: %s", err.Error())
					}
					defer logs.Close()
					buf := new(bytes.Buffer)
					size, err := io.Copy(buf, logs)
					if err != nil {
						log.Printf("error copying logs: %s", err.Error())
					}
					log.Printf("found %d bytes of logs", size)
					response.Items = append(response.Items, parvaeres.Logs{
						Pod:       p.Name,
						Container: c.Name,
						Logs:      buf.String(),
					})
				}
			}
		}
	}
	return
}

func (a *ArgoCD) listApplications(email, repoURL, path string) (*v1alpha1.ApplicationList, error) {
	// See comment in newApplication
	selector := fmt.Sprintf("parvaeres.io/email=%s, parvaeres.io/repoURL=%s, parvaeres.io/path=%s",
		sha1String(email), sha1String(repoURL), sha1String(path))

	apps, err := a.ArgoCDclient.ArgoprojV1alpha1().Applications(a.ArgoCDNamespace).List(
		metav1.ListOptions{
			LabelSelector: selector,
		})

	return apps, errors.Wrap(err, "Unable to list applications")
}

func (a *ArgoCD) getApplicationByName(name string) (*v1alpha1.Application, error) {

	selector := fmt.Sprintf("metadata.name=%s", name)

	apps, err := a.ArgoCDclient.ArgoprojV1alpha1().Applications(a.ArgoCDNamespace).List(
		metav1.ListOptions{
			FieldSelector: selector,
		})

	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Application")
	}

	if len(apps.Items) > 0 {
		return &apps.Items[0], nil
	}

	// return nothing without error if the application is not found
	return nil, fmt.Errorf("Application not found")
}

func (a *ArgoCD) getDeploymentStatusOfApplication(application *v1alpha1.Application) (*parvaeres.DeploymentStatus, error) {
	if application == nil {
		return nil, fmt.Errorf("getDeploymentStatusOfApplication failed: application is nil")
	}
	deploymentStatus := &parvaeres.DeploymentStatus{
		UUID:     application.GetName(),
		Email:    application.Annotations["parvaeres.io/email"],
		RepoURL:  application.Annotations["parvaeres.io/repoURL"],
		Path:     application.Annotations["parvaeres.io/path"],
		LiveURLs: a.getExternalURLsOfApplication(application),
		Errors:   getErrorsOfApplication(application),
		Status:   getDeploymentStatusTypeOfApplication(application),
	}
	return deploymentStatus, nil
}

func (a *ArgoCD) getExternalURLsOfApplication(application *v1alpha1.Application) (urls []string) {
	urls = []string{}
	if application == nil {
		return
	}

	// Check ArgoCD summary - Ingresses only
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

	// See if we have any Service of type LoadBalancer
	for _, r := range application.Status.Resources {
		if r.Kind == "Service" {
			s, err := a.KubernetesClient.CoreV1().Services(r.Namespace).Get(r.Name, metav1.GetOptions{})
			if err != nil {
				log.Printf("cannot get service: %s: %v", r.Name, err)
			}
			if s.Spec.Type == v1.ServiceTypeLoadBalancer {
				for _, i := range s.Status.LoadBalancer.Ingress {
					for _, p := range s.Spec.Ports {
						if i.Hostname != "" {
							urls = append(urls, fmt.Sprintf("http://%s:%d/", i.Hostname, p.Port))
						}
						if i.IP != "" {
							urls = append(urls, fmt.Sprintf("http://%s:%d/", i.IP, p.Port))
						}
					}
				}
			}

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

func getDeploymentStatusTypeOfApplication(application *v1alpha1.Application) parvaeres.DeploymentStatusType {
	// UNKOWN = we don't know much about the application, like when is nil
	if application == nil {
		return parvaeres.UNKNOWN
	}

	// PENDING: if an application is not DEPLOYED or not ERROR
	// then, if it has no SyncPolicy it must be PENDING, i.e. needs user confirmation
	if reflect.DeepEqual(application.Spec.SyncPolicy, &v1alpha1.SyncPolicy{}) {
		return parvaeres.PENDING
	}

	// If Application has a Status it might be DEPLOYED or in ERROR status
	if !reflect.DeepEqual(application.Status, v1alpha1.ApplicationStatus{}) {

		// DEPLOYED: status is Healty and Synced
		if application.Status.Health.Status == health.HealthStatusHealthy ||
			application.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
			return parvaeres.DEPLOYED
		}

		// ERROR: status not Healthy
		if application.Status.Health.Status == health.HealthStatusDegraded ||
			application.Status.Health.Status == health.HealthStatusUnknown ||
			application.Status.Health.Status == health.HealthStatusSuspended ||
			application.Status.Health.Status == health.HealthStatusMissing ||
			applicationHasErrorConditions(application) {
			return parvaeres.ERROR
		}
	}

	// SYNCING: an application has been confirmed by the user and has a SyncPolicy
	if reflect.DeepEqual(application.Spec.SyncPolicy.Automated,
		&v1alpha1.SyncPolicyAutomated{Prune: true, SelfHeal: true}) {
		return parvaeres.SYNCING
	}

	return parvaeres.UNKNOWN
}

func applicationHasErrorConditions(application *v1alpha1.Application) bool {
	for _, c := range application.Status.Conditions {
		if c.IsError() {
			return true
		}
	}
	return false
}

func (a *ArgoCD) setApplicationAutomatedSyncPolicy(application *v1alpha1.Application) error {
	application.Spec.SyncPolicy = &v1alpha1.SyncPolicy{
		Automated: &v1alpha1.SyncPolicyAutomated{
			Prune:    true,
			SelfHeal: true,
		},
		SyncOptions: v1alpha1.SyncOptions{},
	}
	_, err := a.ArgoCDclient.ArgoprojV1alpha1().Applications(a.ArgoCDNamespace).Update(application)
	return err
}

func (a *ArgoCD) createApplication(email, repoURL, path string) (*v1alpha1.Application, error) {
	newApplication, err := newApplication(email, repoURL, path)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	err = a.createNamespace(newApplication.GetObjectMeta().GetName())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	return a.ArgoCDclient.ArgoprojV1alpha1().Applications(a.ArgoCDNamespace).Create(newApplication)
}

func (a *ArgoCD) createNamespace(name string) error {

	// check if namespace already exists, if true exit return without errors
	namespaces, err := a.KubernetesClient.CoreV1().Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	})
	if err != nil {
		return err
	}
	if len(namespaces.Items) > 0 {
		return nil
	}

	namespaceSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_, err = a.KubernetesClient.CoreV1().Namespaces().Create(namespaceSpec)

	return err
}

func newApplication(email, repoURL, path string) (*v1alpha1.Application, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}

	// FIXME: we need validation for email, repoURL and path!

	newApplication := getDefaultApplication()
	newApplication.Spec.Source.RepoURL = repoURL
	newApplication.Spec.Source.Path = path
	newApplication.Spec.Destination.Namespace = id.String()
	newApplication.ObjectMeta.Name = id.String()

	/* Here we set a few labels to be able to search/watch based on them.
	*
	* Labels have a limited allowed character set and length, so we use sha1.
	* However, to keep the information immediately human readable we also set annotations.
	*
	* FIXME: figure out if there's a better encoding for this.
	 */
	newApplication.ObjectMeta.Labels["parvaeres.io/email"] = sha1String(email)
	newApplication.ObjectMeta.Labels["parvaeres.io/repoURL"] = sha1String(repoURL)
	newApplication.ObjectMeta.Labels["parvaeres.io/path"] = sha1String(path)
	newApplication.ObjectMeta.Annotations["parvaeres.io/email"] = email
	newApplication.ObjectMeta.Annotations["parvaeres.io/repoURL"] = repoURL
	newApplication.ObjectMeta.Annotations["parvaeres.io/path"] = path

	return newApplication, nil
}

func getDefaultApplication() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name: "defaultApplication",

			Annotations: map[string]string{
				"parvaeres.io/email":   "",
				"parvaeres.io/repoURL": "",
			},
			Labels: map[string]string{
				"parvaeres.io/email":   "",
				"parvaeres.io/repoURL": "",
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		Spec: v1alpha1.ApplicationSpec{
			Project: "default",
			Source: v1alpha1.ApplicationSource{
				Path:           "",
				TargetRevision: "HEAD",
				RepoURL:        "",
			},
			Destination: v1alpha1.ApplicationDestination{
				Namespace: "default",
				Server:    "https://kubernetes.default.svc",
			},
			SyncPolicy: &v1alpha1.SyncPolicy{},
		},
	}
}

func sha1String(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
