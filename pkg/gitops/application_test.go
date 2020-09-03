package gitops

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

//TestApplication tests the application generation
func TestNewApplication(t *testing.T) {
	Convey("Given a url and an email", t, func() {
		inputURL := "http://blabla"
		inputEmail := "my@email.com"
		Convey("When creating an Application", func() {
			newApplication, err := newApplication(inputEmail, inputURL)
			Convey("Then the Application fields are populated as expected", func() {
				So(err, ShouldBeNil)
				So(newApplication.Spec.Source.RepoURL, ShouldEqual, inputURL)
				So(newApplication.ObjectMeta.Annotations["parvaeres-email"], ShouldEqual, inputEmail)
				So(newApplication.ObjectMeta.Annotations["parvaeres-repoURL"], ShouldEqual, inputURL)
				Convey("And the name field is a UUID", func() {
					uuid, err := uuid.Parse(newApplication.ObjectMeta.Name)
					So(err, ShouldBeNil)
					So(uuid.Version().String(), ShouldEqual, "VERSION_4")
				})
			})
		})
	})

}

func TestGetArgoCDClient(t *testing.T) {
	Convey("When trying to get a client outside the cluster", t, func() {
		clientset, err := getArgoCDClient()
		Convey("Then it fails", func() {
			So(err, ShouldNotBeNil)
			So(clientset, ShouldBeNil)
			So(err.Error(), ShouldStartWith, "unable to load in-cluster configuration")
		})
	})
}
