package parvaeres

//GitopsProviderInterface represents a GitOps thing, like ArgoCD
type GitopsProviderInterface interface {
	CreateDeployment(CreateDeploymentRequest) (*CreateDeploymentResponse, error)
	DeleteDeployment(deploymentID string) (*GetDeploymentResponse, error)
	EnableDeployment(deploymentID string) error
	GetDeploymentByID(deploymentID string) (*GetDeploymentResponse, error)
	GetDeploymentLogs(deploymentID string) (*GetDeploymentLogsResponse, error)
}
