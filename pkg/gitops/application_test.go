package gitops

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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
			newApplication, err := NewApplication(inputEmail, inputURL, inputPath)
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
		c := NewGitOpsClient().WithKubernetesClient(testclient.NewSimpleClientset())
		Convey(fmt.Sprintf("When creating a Namespace named '%s'", expectedName), func() {
			err := c.CreateNamespace(expectedName)
			So(err, ShouldBeNil)
			Convey("Then the Namespace is created as expected", func() {
				namespace, err := c.KubernetesClient.CoreV1().Namespaces().Get(expectedName, v1.GetOptions{})
				So(err, ShouldBeNil)
				So(namespace.Name, ShouldEqual, expectedName)
			})
		})
	})
}
