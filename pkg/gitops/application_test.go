package parvaeres

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

//TestApplication tests the application generation
func TestGenerateApplication(t *testing.T) {
	Convey("Given a url and an email", t, func() {
		inputURL := "http://blabla"
		inputEmail := "my@email.com"
		Convey("When creating an Application", func() {
			newApplication, err := GenerateApplication(inputEmail, inputURL)
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
