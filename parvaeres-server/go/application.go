package parvaeres

import (
	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/google/uuid"
)

func GenerateApplication(email string, repoURL string) (*v1alpha1.Application, error) {

	// generate UUID
	// create Deployment{email, repository, UUID, created_at_timestamp, status}
	// store Deployment in data store
	// create corresponding ArgoCD Application object
	// write it in k8s

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	newApplication := v1alpha1.Application{}
	newApplication.Name = id.String()
	newApplication.Annotations["parvaeres-email"] = email
	newApplication.Annotations["parvaeres-repoURL"] = repoURL

	return &newApplication, nil
}
