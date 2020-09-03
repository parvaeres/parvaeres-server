package gitops

import (
	"fmt"
	"log"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	argoclient "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func getArgoCDClient() (*argoclient.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return argoclient.NewForConfig(config)
}

func getDefaultApplication() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "defaultApplication",
			Namespace: argocdNamespace,
			Annotations: map[string]string{
				"parvaeres-email":   "",
				"parvaeres-repoURL": "",
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
			SyncPolicy: &v1alpha1.SyncPolicy{
				Automated: &v1alpha1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
				SyncOptions: v1alpha1.SyncOptions{},
			},
		},
	}
}

// CreateApplication returns an ArgoCD Application relative to email and repoURL
func CreateApplication(email string, repoURL string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "Unable to create application")
	}
	client, err := getArgoCDClient()
	if err != nil {
		return errors.Wrap(err, "Unable to create application")
	}

	newApplication := getDefaultApplication()
	newApplication.Spec.Source.RepoURL = repoURL
	newApplication.ObjectMeta.Name = id.String()
	newApplication.ObjectMeta.Annotations["parvaeres-email"] = email
	newApplication.ObjectMeta.Annotations["parvaeres-repoURL"] = repoURL

	_, err = client.ArgoprojV1alpha1().Applications(argocdNamespace).Create(newApplication)

	return err
}

// ListApplications returns a list of ArgoCD applications
func ListApplications(email string, repoURL string) (*v1alpha1.ApplicationList, error) {
	client, err := getArgoCDClient()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list Applications")
	}

	selector := fmt.Sprintf("parvaeres-email=%s, parvaeres-repoURL=%s", email, repoURL)

	apps, err := client.ArgoprojV1alpha1().Applications(argocdNamespace).List(
		metav1.ListOptions{
			LabelSelector: selector,
		})

	return apps, errors.Wrap(err, "Unable to list applications")
}
