package gitops

import (
	"encoding/hex"
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
				Path:           "/",
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

func newApplication(email string, repoURL string) (*v1alpha1.Application, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}

	newApplication := getDefaultApplication()
	newApplication.Spec.Source.RepoURL = repoURL
	newApplication.ObjectMeta.Name = id.String()
	/* Labels have a limited allowed character set, so we encode the data in hex
	* maybe there's a more efficient encoding. Labels are useful for search/watch operations.
	* I prefer using hex instead of invertible hashing so data can be retrieved.
	* However, to keep the information human readable we also set annotations.
	*
	* FIXME: figure out if there's a better encoding for this.
	 */
	newApplication.ObjectMeta.Labels["parvaeres.io/email"] = hex.EncodeToString([]byte(email))
	newApplication.ObjectMeta.Labels["parvaeres.io/repoURL"] = hex.EncodeToString([]byte(repoURL))
	newApplication.ObjectMeta.Annotations["parvaeres.io/email"] = email
	newApplication.ObjectMeta.Annotations["parvaeres.io/repoURL"] = repoURL

	return newApplication, nil
}

// CreateApplication returns an ArgoCD Application relative to email and repoURL
func CreateApplication(email string, repoURL string) (*v1alpha1.Application, error) {
	client, err := getArgoCDClient()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	newApplication, err := newApplication(email, repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create application")
	}
	return client.ArgoprojV1alpha1().Applications(argocdNamespace).Create(newApplication)
}

// ListApplications returns a list of ArgoCD applications
func ListApplications(email string, repoURL string) (*v1alpha1.ApplicationList, error) {
	client, err := getArgoCDClient()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list Applications")
	}

	// See comment in newApplication
	selector := fmt.Sprintf("parvaeres.io/email=%s, parvaeres.io/repoURL=%s",
		hex.EncodeToString([]byte(email)), hex.EncodeToString([]byte(repoURL)))

	apps, err := client.ArgoprojV1alpha1().Applications(argocdNamespace).List(
		metav1.ListOptions{
			LabelSelector: selector,
		})

	return apps, errors.Wrap(err, "Unable to list applications")
}

//GetApplicationByName returns an ArgoCD application with the corresponding name
func GetApplicationByName(name string) (*v1alpha1.Application, error) {
	client, err := getArgoCDClient()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Application")
	}

	selector := fmt.Sprintf("metadata.name=%s", name)

	apps, err := client.ArgoprojV1alpha1().Applications(argocdNamespace).List(
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
