package gitops

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	argoclient "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type GitOpsClient struct {
	ArgoCDclient     argoclient.Interface
	KubernetesClient kubernetes.Interface
}

var kubeconfig string = ""
var argocdNamespace string = "argocd"

/*
FIXME: hardcoded config for now, interferes with GoConvey flags during testing
func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "path to Kubernetes config file")
	flag.StringVar(&argocdNamespace, "argocdNamespace", "argocd", "argocd Namespace")
	flag.Parse()
}
*/

func NewGitOpsClient() *GitOpsClient {
	return &GitOpsClient{}
}

func (g *GitOpsClient) WithArgoCDClient(client argoclient.Interface) *GitOpsClient {
	g.ArgoCDclient = client
	return g
}

func (g *GitOpsClient) WithKubernetesClient(client kubernetes.Interface) *GitOpsClient {
	g.KubernetesClient = client
	return g
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

func (g *GitOpsClient) CreateNamespace(name string) error {

	// check if namespace already exists, if true exit return without errors
	namespaces, err := g.KubernetesClient.CoreV1().Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	})
	if err != nil {
		return err
	}
	if len(namespaces.Items) > 0 {
		return nil
	}

	namespaceSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_, err = g.KubernetesClient.CoreV1().Namespaces().Create(namespaceSpec)

	return err
}

func getDefaultApplication() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name: "defaultApplication",

			Namespace: argocdNamespace,
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

//NewApplication returns an application with fields based on input parameters
func NewApplication(email, repoURL, path string) (*v1alpha1.Application, error) {
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

// CreateApplication returns an ArgoCD Application relative to email, repoURL and path
func (g *GitOpsClient) CreateApplication(email, repoURL, path string) (*v1alpha1.Application, error) {
	newApplication, err := NewApplication(email, repoURL, path)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	err = g.CreateNamespace(newApplication.GetObjectMeta().GetName())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	return g.ArgoCDclient.ArgoprojV1alpha1().Applications(argocdNamespace).Create(newApplication)
}

//SetApplicationAutomatedSyncPolicy sets the sync policy for the application to Automated
func (g *GitOpsClient) SetApplicationAutomatedSyncPolicy(application *v1alpha1.Application) error {
	application.Spec.SyncPolicy = &v1alpha1.SyncPolicy{
		Automated: &v1alpha1.SyncPolicyAutomated{
			Prune:    true,
			SelfHeal: true,
		},
		SyncOptions: v1alpha1.SyncOptions{},
	}
	_, err := g.ArgoCDclient.ArgoprojV1alpha1().Applications(argocdNamespace).Update(application)
	return err
}

// ListApplications returns a list of ArgoCD applications
func (g *GitOpsClient) ListApplications(email, repoURL, path string) (*v1alpha1.ApplicationList, error) {
	// See comment in newApplication
	selector := fmt.Sprintf("parvaeres.io/email=%s, parvaeres.io/repoURL=%s, parvaeres.io/path=%s",
		sha1String(email), sha1String(repoURL), sha1String(path))

	apps, err := g.ArgoCDclient.ArgoprojV1alpha1().Applications(argocdNamespace).List(
		metav1.ListOptions{
			LabelSelector: selector,
		})

	return apps, errors.Wrap(err, "Unable to list applications")
}

//GetApplicationByName returns an ArgoCD application with the corresponding name
func (g *GitOpsClient) GetApplicationByName(name string) (*v1alpha1.Application, error) {

	selector := fmt.Sprintf("metadata.name=%s", name)

	apps, err := g.ArgoCDclient.ArgoprojV1alpha1().Applications(argocdNamespace).List(
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
