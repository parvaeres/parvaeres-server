package parvaeres

import (
	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDefaultApplication() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "defaultApplication",
			Namespace: "default",
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

// GenerateApplication returns an ArgoCD Application relative to email and repoURL
func GenerateApplication(email string, repoURL string) (*v1alpha1.Application, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	newApplication := getDefaultApplication()
	newApplication.Spec.Source.RepoURL = repoURL
	newApplication.ObjectMeta.Name = id.String()
	newApplication.ObjectMeta.Annotations["parvaeres-email"] = email
	newApplication.ObjectMeta.Annotations["parvaeres-repoURL"] = repoURL

	return newApplication, nil
}
