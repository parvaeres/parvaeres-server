package parvaeres

import (
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/smartystreets/goconvey/convey"
)

/*TestGetDeploymentStatusTypeOfApplication tests the translation of the state from ArgoCD
* to something we can communicate to the user DeploymentStatusType enum in API spec
 */

type DeploymentStatusTypeTest struct {
	description string
	application *v1alpha1.Application
	status      DeploymentStatusType
}

var tests = []DeploymentStatusTypeTest{
	{
		description: "Given a nil application",
		status:      UNKNOWN,
		application: nil,
	},
	{
		description: "Given an application with empty SyncPolicy",
		status:      PENDING,
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
		},
	},
	{
		description: "Given an application with set SyncPolicy and no Status",
		status:      SYNCING,
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{
					Automated: &v1alpha1.SyncPolicyAutomated{
						Prune:    true,
						SelfHeal: true,
					},
				},
			},
		},
	},
	{
		description: "Given an application with SyncPolicy and HealthStatusHealthy",
		status:      DEPLOYED,
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{
					Automated: &v1alpha1.SyncPolicyAutomated{
						Prune:    true,
						SelfHeal: true,
					},
				},
			},
			Status: v1alpha1.ApplicationStatus{
				Health: v1alpha1.HealthStatus{
					Status: health.HealthStatusHealthy,
				},
			},
		},
	},
	{
		description: "Given an application with SyncPolicy and HealthStatusMissing",
		status:      ERROR,
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{
					Automated: &v1alpha1.SyncPolicyAutomated{
						Prune:    true,
						SelfHeal: true,
					},
				},
			},
			Status: v1alpha1.ApplicationStatus{
				Health: v1alpha1.HealthStatus{
					Status: health.HealthStatusMissing,
				},
			},
		},
	},
}

func TestGetDeploymentStatusTypeOfApplication(t *testing.T) {

	for _, test := range tests {
		Convey(test.description, t, func() {
			Convey("When getting the corresponding DeploymentStatusType", func() {
				deploymentState := getDeploymentStatusTypeOfApplication(test.application)
				Convey("Then the status value is as expected", func() {
					So(deploymentState, ShouldEqual, test.status)
				})
			})

		})
	}
}
