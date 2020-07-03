package parvaeres

import (
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	. "github.com/smartystreets/goconvey/convey"
)

//TestApplication tests the application generation
func TestApplication(t *testing.T) {

	expectedApplication := v1alpha1.Application{}
	expectedApplication.Name = "something"

	Convey("Given a url and an email", t, func() {
		inputURL := "http://blabla"
		inputEmail := "my@email.com"
		Convey("When creating an Application", func() {
			//the function to test
			newApplication, err := generateApplication(inputEmail, inputURL)
			Convey("Then the Application is created as expected", func() {
				So(expectedApplication, ShouldEqual, newApplication)
			})
		})
	})

}
