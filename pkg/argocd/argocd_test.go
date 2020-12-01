package argocd

import (
	"fmt"
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/google/uuid"
	parvaeres "github.com/parvaeres/parvaeres/pkg/api"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

//TestApplication tests the application generation
func TestNewApplication(t *testing.T) {
	Convey("Given a url and an email", t, func() {
		inputURL := "http://blabla"
		inputEmail := "my@email.com"
		inputPath := "/"
		Convey("When creating an Application", func() {
			newApplication, err := newApplication(inputEmail, inputURL, inputPath)
			Convey("Then the Application fields are populated as expected", func() {
				So(err, ShouldBeNil)
				So(newApplication.Spec.Source.RepoURL, ShouldEqual, inputURL)
				So(newApplication.ObjectMeta.Annotations["parvaeres.io/email"], ShouldEqual, inputEmail)
				So(newApplication.ObjectMeta.Annotations["parvaeres.io/repoURL"], ShouldEqual, inputURL)
				So(newApplication.ObjectMeta.Annotations["parvaeres.io/path"], ShouldEqual, inputPath)
				So(newApplication.ObjectMeta.Labels["parvaeres.io/email"], ShouldEqual, sha1String(inputEmail))
				So(newApplication.ObjectMeta.Labels["parvaeres.io/repoURL"], ShouldEqual, sha1String(inputURL))
				So(newApplication.ObjectMeta.Labels["parvaeres.io/path"], ShouldEqual, sha1String(inputPath))
				Convey("And the name field is a UUID", func() {
					uuid, err := uuid.Parse(newApplication.ObjectMeta.Name)
					So(err, ShouldBeNil)
					So(uuid.Version().String(), ShouldEqual, "VERSION_4")
				})
			})
		})
	})

}

func TestCreateNamespace(t *testing.T) {
	expectedName := "myNamespace"
	Convey("Given a GitOpsClient", t, func() {
		c := ArgoCD{
			KubernetesClient: testclient.NewSimpleClientset(),
		}
		Convey(fmt.Sprintf("When creating a Namespace named '%s'", expectedName), func() {
			err := c.createNamespace(expectedName)
			So(err, ShouldBeNil)
			Convey("Then the Namespace is created as expected", func() {
				namespace, err := c.KubernetesClient.CoreV1().Namespaces().Get(expectedName, v1.GetOptions{})
				So(err, ShouldBeNil)
				So(namespace.Name, ShouldEqual, expectedName)
			})
		})
	})
}

/*TestGetDeploymentStatusTypeOfApplication tests the translation of the state from ArgoCD
* to something we can communicate to the user DeploymentStatusType enum in API spec
 */

type DeploymentStatusTypeTest struct {
	description string
	application *v1alpha1.Application
	status      parvaeres.DeploymentStatusType
}

var tests = []DeploymentStatusTypeTest{
	{
		description: "Given a nil application",
		status:      parvaeres.UNKNOWN,
		application: nil,
	},
	{
		description: "Given an application with empty SyncPolicy",
		status:      parvaeres.PENDING,
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
		},
	},
	{
		description: "Given an application with set SyncPolicy and no Status",
		status:      parvaeres.SYNCING,
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
		description: "Given an application with SyncPolicy and SyncStatusCodeSynced",
		status:      parvaeres.DEPLOYED,
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
				Sync: v1alpha1.SyncStatus{
					Status: v1alpha1.SyncStatusCodeSynced,
				},
			},
		},
	},
	{
		description: "Given an application with SyncPolicy and HealthStatusHealthy",
		status:      parvaeres.SYNCING,
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
		status:      parvaeres.ERROR,
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
	status         *parvaeres.DeploymentStatus
	errorString    string
	errorAssertion func(interface{}, ...interface{}) string
}

var deploymentStatusTests = []deploymentStatusTest{
	{
		description:    "Given a nil application",
		status:         nil,
		application:    nil,
		errorString:    "getDeploymentStatusOfApplication failed: application is nil",
		errorAssertion: ShouldNotBeNil,
	},
	{
		description: "Given an application with empty SyncPolicy",
		status: &parvaeres.DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{},
			Errors:   []string{},
			Status:   parvaeres.PENDING,
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
		status: &parvaeres.DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{"http://app.lol.com"},
			Errors:   []string{},
			Status:   parvaeres.PENDING,
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
	{
		description: "Given an application with Some URLs with a Path",
		status: &parvaeres.DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{"http://app1.lol.com", "http://app2.lol.com"},
			Errors:   []string{},
			Status:   parvaeres.PENDING,
		},
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
			Status: v1alpha1.ApplicationStatus{
				Summary: v1alpha1.ApplicationSummary{
					ExternalURLs: []string{"http://app1.lol.com/path", "http://app2.lol.com/*"},
				},
			},
		},
		errorString:    "",
		errorAssertion: ShouldBeNil,
	},
	{
		description: "Given an application with an error condition",
		status: &parvaeres.DeploymentStatus{
			UUID:     "",
			LiveURLs: []string{},
			Errors:   []string{"There is something wrong with mass consumption"},
			Status:   parvaeres.PENDING,
		},
		application: &v1alpha1.Application{
			Spec: v1alpha1.ApplicationSpec{
				SyncPolicy: &v1alpha1.SyncPolicy{},
			},
			Status: v1alpha1.ApplicationStatus{
				Conditions: []v1alpha1.ApplicationCondition{
					v1alpha1.ApplicationCondition{
						Type:    v1alpha1.ApplicationConditionInvalidSpecError,
						Message: "There is something wrong with mass consumption",
					},
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
				c := ArgoCD{
					KubernetesClient: testclient.NewSimpleClientset(),
				}
				deploymentStatus, err := c.getDeploymentStatusOfApplication(test.application)
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
