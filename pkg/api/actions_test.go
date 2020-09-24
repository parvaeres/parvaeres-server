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

type deploymentStatusTest struct {
	description    string
	application    *v1alpha1.Application
	status         *DeploymentStatus
	errorString    string
	errorAssertion func(interface{}, ...interface{}) string
}

var deploymentStatusTests = []deploymentStatusTest{
	{
		description:    "Given a nil application",
		status:         nil,
		application:    nil,
		errorString:    "GetDeploymentStatusOfApplication failed: application is nil",
		errorAssertion: ShouldNotBeNil,
	},
	{
		description: "Given an application with empty SyncPolicy",
		status: &DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{},
			Errors:   []string{},
			Status:   PENDING,
		},
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
		},
		errorString:    "",
		errorAssertion: ShouldBeNil,
	},
	{
		description: "Given an application with Some URLs",
		status: &DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{"http://app.lol.com"},
			Errors:   []string{},
			Status:   PENDING,
		},
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
			Status: v1alpha1.ApplicationStatus{
				Summary: v1alpha1.ApplicationSummary{
					ExternalURLs: []string{"http://app.lol.com"},
				},
			},
		},
		errorString:    "",
		errorAssertion: ShouldBeNil,
	},
}

func TestGetDeploymentStatusOfApplication(t *testing.T) {

	for _, test := range deploymentStatusTests {
		Convey(test.description, t, func() {
			Convey("When getting the corresponding DeploymentStatus", func() {
				deploymentStatus, err := GetDeploymentStatusOfApplication(test.application)
				Convey("Then the status value is as expected", func() {
					So(err, test.errorAssertion)
					if err != nil {
						So(err.Error(), ShouldEqual, test.errorString)
					}
					So(deploymentStatus, ShouldResemble, test.status)
				})
			})

		})
	}
}
