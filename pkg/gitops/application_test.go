package gitops

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
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
				So(newApplication.ObjectMeta.Labels["parvaeres.io/email"], ShouldEqual, hex.EncodeToString([]byte(inputEmail)))
				So(newApplication.ObjectMeta.Labels["parvaeres.io/repoURL"], ShouldEqual, hex.EncodeToString([]byte(inputURL)))
				So(newApplication.ObjectMeta.Labels["parvaeres.io/path"], ShouldEqual, hex.EncodeToString([]byte(inputPath)))
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
